package db

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestNewStorage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(fmt.Errorf("an error '%w' was not expected when opening a stub database connection", err))
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("create table if not exists shortener .*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("create index .*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("create unique .*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("alter table shortener.*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("create table if not exists user_url.*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("create unique index.*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	_, err = NewStorage(db)
	if err != nil {
		t.Fatal(err)
	}
}
