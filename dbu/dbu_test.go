package dbu

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestOpenConn(t *testing.T) {
	// Given
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLoggerSugared := zap.New(observedZapCore).Sugar()

	dbo, err := OpenConn(observedLoggerSugared, "dnd", "5e",
		"dev", "webuser")
	assert.Equal(t, nil, err)

	assert.Equal(t, 0, observedLogs.Len())
	err = dbo.CleanUpAndClose()
	assert.Equal(t, nil, err)

}
func TestExec(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLoggerSugared := zap.New(observedZapCore).Sugar()

	dbo, err := OpenConn(observedLoggerSugared, "dnd", "5e",
		"dev", "webuser")
	assert.Equal(t, nil, err)

	err = dbo.Exec(context.Background(), `Insert into public."user"(first_name, email, `+
		`created_at, updated_at) values ($1, $2, now(), now())`,
		"demotestuser", "demotestuser@mycrazydomain.io")
	assert.Equal(t, nil, err)
	// lastId, err := exec.LastInsertId()
	// last Insert Id is not supported by the driver
	// nbrAffected, err := exec.RowsAffected()
	assert.Equal(t, nil, err)
	// assert.Equal(t, int64(1), nbrAffected)
	assert.Equal(t, 0, observedLogs.Len())
	err = dbo.CleanUpAndClose()
	assert.Equal(t, nil, err)
}

func TestQueryReturnId(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLoggerSugared := zap.New(observedZapCore).Sugar()

	dbo, err := OpenConn(observedLoggerSugared, "dnd", "5e",
		"dev", "webuser")
	assert.Equal(t, nil, err)

	id, err := dbo.QueryReturnId(context.Background(), `Insert into public."user"(first_name, email, `+
		`created_at, updated_at) values ($1, $2, now(), now()) returning id`,
		"demotestuser", "demotestuser@mycrazydomain.io")
	assert.Equal(t, nil, err)

	assert.GreaterOrEqual(t, id, int64(1))
	assert.Equal(t, 0, observedLogs.Len())
	err = dbo.CleanUpAndClose()
	assert.Equal(t, nil, err)
}

func TestQuery(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLoggerSugared := zap.New(observedZapCore).Sugar()

	dbo, err := OpenConn(observedLoggerSugared, "dnd", "5e",
		"dev", "webuser")
	assert.Equal(t, nil, err)

	rows, err := dbo.Query(context.Background(), `select id, first_name, email, created_at from public."user" `+
		`where first_name=$1 and email=$2`,
		"demotestuser", "demotestuser@mycrazydomain.io")
	assert.Equal(t, nil, err)
	rowCnt := 0
	for rows.Next() {
		rowCnt++
		var i int64
		var f sql.NullString
		var fs string
		var e sql.NullString
		var es string
		var d sql.NullTime
		err := rows.Scan(&i, &f, &e, &d)
		assert.Equal(t, nil, err)
		if f.Valid {
			fs = f.String
		} else {
			fs = "Unknown"
		}
		if e.Valid {
			es = e.String
		} else {
			es = "Unknown"
		}
		assert.Equal(t, true, d.Valid)
		assert.Equal(t, "demotestuser", fs)
		assert.Equal(t, "demotestuser@mycrazydomain.io", es)
	}
	assert.Equal(t, 2, rowCnt)
	assert.Equal(t, 0, observedLogs.Len())
	rows.Close()
	err = dbo.CleanUpAndClose()
	assert.Equal(t, nil, err)
}

func TestUserCleanup(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLoggerSugared := zap.New(observedZapCore).Sugar()
	dbo, err := OpenConn(observedLoggerSugared, "dnd", "5e",
		"dev", "webuser")
	assert.Equal(t, nil, err)
	err = dbo.Exec(context.Background(), "Delete from public.\"user\" where first_name = 'demotestuser'")
	assert.Equal(t, nil, err)
	// nbrAffected, err := exec.RowsAffected()
	// assert.Equal(t, nil, err)
	// assert.GreaterOrEqual(t, nbrAffected, int64(1))
	assert.Equal(t, 0, observedLogs.Len())
	err = dbo.CleanUpAndClose()
	assert.Equal(t, nil, err)
}
