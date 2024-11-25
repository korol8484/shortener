package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db    *sql.DB
	mock  sqlmock.Sqlmock
	store *DBStorage
)

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (int, error) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		return -1, err
	}

	mock.ExpectBegin()
	mock.ExpectExec("create table if not exists \"user\"").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store, err = NewStorage(db)
	if err != nil {
		return -1, err
	}

	return m.Run(), nil
}

func TestDBStorage_NewUser(t *testing.T) {
	mock.ExpectQuery("INSERT INTO \"user\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	user, err := store.NewUser(context.Background())

	require.NoError(t, err)
	assert.Equal(t, user.ID, int64(1))

	mock.ExpectQuery("INSERT INTO \"user\"").
		WillReturnError(sql.ErrNoRows)

	_, err = store.NewUser(context.Background())

	require.Error(t, err)
}
