package db

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (int, error) {
	var err error

	db, err = sql.Open("sqlite3", "file:./test.db?cache=shared&mode=memory")
	if err != nil {
		return -1, fmt.Errorf("could not connect to database: %w", err)
	}

	defer func() {
		_ = db.Close()
	}()

	return m.Run(), nil
}

func TestNewStorage(t *testing.T) {
	_, err := NewStorage(db)
	if err != nil {
		t.Fatal(err)
	}
}
