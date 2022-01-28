package dbu

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

var (
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
	*sql.DB
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
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		du.Host, du.Port, du.User, du.Password, du.Dbname)
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
	file.Close()
	var dbusers DbUsers
	json.Unmarshal(byteValue, &dbusers)
	for i := 0; i < len(dbusers.DbUsers); i++ {
		if system == dbusers.DbUsers[i].System &&
			version == dbusers.DbUsers[i].Version &&
			env == dbusers.DbUsers[i].Env &&
			dbname == dbusers.DbUsers[i].Dbname {
			connInfo := dbusers.DbUsers[i].getConnString()
			db, err := sql.Open("postgres", connInfo)
			if err != nil {
				return nil, ErrCouldNotOpenDB
			}
			err = db.Ping()
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
	err := dbo.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

// Exec performs a query on the db that doesn't return rows.
func (dbo *DBo) Exec(Query string, Args ...interface{}) (sql.Result, error) {
	result, err := dbo.DB.Exec(Query, Args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// QueryReturnId performs a query on the db that doesn't return rows.
func (dbo *DBo) QueryReturnId(Query string, Args ...interface{}) (int64, error) {
	var id int64
	row := dbo.DB.QueryRow(Query, Args...)
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	if row.Err() != nil {
		return -1, row.Err()
	}
	return id, nil
}
