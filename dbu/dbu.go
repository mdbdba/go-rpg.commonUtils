package dbu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

var (
	// ErrConfigInfo is returned when the database connection config info could not be parsed
	ErrConfigInfo = errors.New("dbu: could not make sense out of the db config info")
	// ErrConnFileNotFound is returned when connection info file cannot be found
	ErrConnFileNotFound = errors.New("dbu: connection info file not found")
	// ErrConnFileNotReadable is returned when connection info file cannot be read
	ErrConnFileNotReadable = errors.New("dbu: connection info file not readable")
	// ErrConnNotFound is returned when connection info cannot be found
	// that matches what was passed.
	ErrConnNotFound = errors.New("dbu: connection info not found")
	// ErrCouldNotOpenDB is returned when the DB could not be opened
	ErrCouldNotOpenDB = errors.New("dbu: could not open DB")
	// ErrDBUnreachable is returned when the DB is unreachable.
	ErrDBUnreachable = errors.New("dbu: DB unreachable")
	// ErrUnexpectedEffectedCnt is returned when an expected number of effected rows is not matched
	ErrUnexpectedEffectedCnt = errors.New("dbu: unexpected effected count")
	// ErrNotFound is returned when a query returns no rows
	ErrNotFound = errors.New("dbu: no records found")
)

type DbUsers struct {
	DbUsers []DbUser `json:"db_users"`
}

type DbUser struct {
	System   string `json:"system"`
	Version  string `json:"version"`
	Env      string `json:"env"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

// DBo struct encapsulates sql.DB to add new functions
type DBo struct {
	*pgx.Conn
	logger  *zap.SugaredLogger
	System  string
	Version string
	Env     string
	User    string
	Dbname  string
}

func (du *DbUser) ToString() string {
	return fmt.Sprintf("System: %s, Version: %s, Env: %s, Host: %s, Port: %d, User: %s, Password: %s, Dbname: %s\n",
		du.System, du.Version, du.Env, du.Host, du.Port, du.User, du.Password, du.Dbname)
}

func (du *DbUser) getConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=require", du.User, du.Password,
		du.Host, du.Port, du.Dbname)
}

// OpenConn creates a db connection to one of the postgres dbs
func OpenConn(logger *zap.SugaredLogger, system string,
	version string, env string, dbname string) (*DBo, error) {
	file, err := os.Open("dbi.json")
	if err != nil {
		return nil, ErrConnFileNotFound
	}
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, ErrConnFileNotReadable
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	var dbusers DbUsers
	err = json.Unmarshal(byteValue, &dbusers)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(dbusers.DbUsers); i++ {
		if system == dbusers.DbUsers[i].System &&
			version == dbusers.DbUsers[i].Version &&
			env == dbusers.DbUsers[i].Env &&
			dbname == dbusers.DbUsers[i].Dbname {
			connInfo := dbusers.DbUsers[i].getConnString()
			config, err := pgx.ParseConfig(connInfo)
			if err != nil {
				return nil, ErrConfigInfo
			}
			//db, err := sql.Open("postgres", connInfo)
			db, err := pgx.ConnectConfig(context.Background(), config)
			if err != nil {
				return nil, ErrCouldNotOpenDB
			}
			err = db.Ping(context.Background())
			if err != nil {
				return nil, ErrDBUnreachable
			}
			return &DBo{
				db,
				logger,
				dbusers.DbUsers[i].System,
				dbusers.DbUsers[i].Version,
				dbusers.DbUsers[i].Env,
				dbusers.DbUsers[i].User,
				dbusers.DbUsers[i].Dbname,
			}, nil

		}
	}
	return nil, ErrConnNotFound
}

// CleanUpAndClose handles any cleanup that needs to happen and closes
// the db connection.
func (dbo *DBo) CleanUpAndClose() error {
	err := dbo.Conn.Close(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func performExec(ctx context.Context, tx pgx.Tx, expected int, query string, args ...interface{}) error {
	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if expected != -1 {
		if result.RowsAffected() != int64(expected) {
			return ErrUnexpectedEffectedCnt
		}
	} else if result.RowsAffected() == int64(0) {
		return ErrNotFound
	}
	return nil
}

// Exec performs a query on the db that doesn't return rows.
func (dbo *DBo) Exec(ctx context.Context, expected int, query string, args ...interface{}) error {
	err := crdbpgx.ExecuteTx(ctx, dbo.Conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return performExec(ctx, tx, expected, query, args...)
	})
	if err != nil {
		return err
	}
	return nil
}

// QueryReturnId performs a query on the db that doesn't return rows.
func (dbo *DBo) QueryReturnId(ctx context.Context, query string, args ...interface{}) (int64, error) {
	var id int64
	row := dbo.Conn.QueryRow(ctx, query, args...)
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (dbo *DBo) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := dbo.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
