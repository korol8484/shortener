package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/korol8484/shortener/internal/app/domain"
	"strings"
	"sync"
)

type Storage struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewStorage(db *sql.DB) (*Storage, error) {
	storage := &Storage{db: db}

	err := storage.migrate(context.Background())
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) Add(ctx context.Context, ent *domain.URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.ExecContext(
		ctx, `INSERT INTO shortener (url, alias) VALUES ($1,$2)`, ent.URL, ent.Alias,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) AddBatch(ctx context.Context, batch domain.BatchURL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	var (
		placeholders []string
		vals         []interface{}
	)

	for i, v := range batch {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d)",
			i*2+1,
			i*2+2,
		))

		vals = append(vals, v.URL, v.Alias)
	}

	insert := fmt.Sprintf("INSERT INTO shortener (url, alias) VALUES %s", strings.Join(placeholders, ","))
	_, err = tx.ExecContext(ctx, insert, vals...)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Read(ctx context.Context, alias string) (*domain.URL, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	row := s.db.QueryRowContext(ctx, "SELECT t.url, t.alias FROM public.shortener t WHERE alias = $1", alias)
	if row.Err() != nil {
		return nil, row.Err()
	}

	ent := &domain.URL{}
	err := row.Scan(&ent.URL, &ent.Alias)
	if err != nil {
		return nil, err
	}

	return ent, nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) migrate(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	_, err = tx.ExecContext(ctx, `
	create table if not exists public.shortener
	(
		id    bigserial
			constraint shortener_pk
				primary key,
		url   varchar(1000) not null,
		alias varchar(10)   not null
	);`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `create index if not exists shortener_alias_index on public.shortener (alias);`)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
