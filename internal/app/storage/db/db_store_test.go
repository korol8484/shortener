package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	db    *sql.DB
	mock  sqlmock.Sqlmock
	store *Storage
)

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}

	err = store.Close()
	if err == nil {
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
	mock.ExpectExec("create table if not exists shortener .*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("create index .*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("create unique .*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("alter table shortener .*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("create table if not exists user_url.*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("create unique index.*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store, err = NewStorage(db)
	if err != nil {
		return -1, err
	}

	return m.Run(), nil
}

func TestStorage_Add(t *testing.T) {
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO shortener").
		WithArgs("1", "1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO user_url").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := store.Add(context.Background(), &domain.URL{URL: "1", Alias: "1"}, &domain.User{ID: 1})
	require.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO shortener").
		WithArgs("2", "2").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))

	mock.ExpectQuery("SELECT t.id FROM shortener t").
		WithArgs("2").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	mock.ExpectExec("INSERT INTO user_url").
		WithArgs(2, 2).
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	err = store.Add(context.Background(), &domain.URL{URL: "2", Alias: "2"}, &domain.User{ID: 2})
	require.ErrorIs(t, err, storage.ErrIssetURL)
}

func TestStorage_Read(t *testing.T) {
	mock.ExpectQuery("SELECT t.url, t.alias, t.deleted FROM shortener t").
		WithArgs("alias").
		WillReturnRows(
			sqlmock.NewRows([]string{"url", "alias", "deleted"}).
				AddRow("http://ya.ru", "alias", false),
		)

	url, err := store.Read(context.Background(), "alias")
	require.NoError(t, err)

	assert.Equal(t, "http://ya.ru", url.URL)
	assert.Equal(t, "alias", url.Alias)
	assert.False(t, url.Deleted)

	mock.ExpectQuery("SELECT t.url, t.alias, t.deleted FROM shortener t").
		WithArgs("alias").WillReturnError(sql.ErrNoRows)

	_, err = store.Read(context.Background(), "alias")
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestStorage_ReadByURL(t *testing.T) {
	mock.ExpectQuery("SELECT t.url, t.alias FROM shortener t").
		WithArgs("http://ya.ru").
		WillReturnRows(
			sqlmock.NewRows([]string{"url", "alias"}).
				AddRow("http://ya.ru", "alias"),
		)

	url, err := store.ReadByURL(context.Background(), "http://ya.ru")
	require.NoError(t, err)

	assert.Equal(t, "http://ya.ru", url.URL)
	assert.Equal(t, "alias", url.Alias)
	assert.False(t, url.Deleted)
}

func TestStorage_ReadUserURL(t *testing.T) {
	mock.ExpectQuery("SELECT s.url, s.alias FROM shortener s").
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"url", "alias"}).
				AddRow("http://ya.ru", "alias"),
		)

	url, err := store.ReadUserURL(context.Background(), &domain.User{ID: 1})
	require.NoError(t, err)

	assert.Len(t, url, 1)
}

func TestStorage_BatchDelete(t *testing.T) {
	mock.ExpectExec("UPDATE shortener s SET deleted = true").
		WithArgs(1, "alias").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := store.BatchDelete(context.Background(), []string{"alias"}, 1)
	require.NoError(t, err)
}

func TestStorage_AddBatch(t *testing.T) {
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO shortener").
		WithArgs("1", "1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO user_url").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := store.AddBatch(context.Background(), domain.BatchURL{
		&domain.URL{URL: "1", Alias: "1"},
	}, &domain.User{ID: 1})
	assert.NoError(t, err)
}

func TestStorage_Close(t *testing.T) {
	cDB, CMock, err := sqlmock.New()
	require.NoError(t, err)

	CMock.ExpectBegin()
	CMock.ExpectExec("create table if not exists shortener .*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	CMock.ExpectExec("create index .*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	CMock.ExpectExec("create unique .*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	CMock.ExpectExec("alter table shortener .*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	CMock.ExpectExec("create table if not exists user_url.*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	CMock.ExpectExec("create unique index.*").WithoutArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	CMock.ExpectCommit()

	cStore, err := NewStorage(cDB)
	require.NoError(t, err)

	CMock.ExpectClose()
	err = cStore.Close()
	require.NoError(t, err)
}
