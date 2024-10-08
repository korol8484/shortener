package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type Config interface {
	GetDsn() string
}

func NewPgDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.GetDsn())
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(time.Minute * 1)

	return db, nil
}
