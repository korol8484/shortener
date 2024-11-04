package storage

import (
	"context"
	"database/sql"
	"sync"

	"github.com/korol8484/shortener/internal/app/domain"
)

// DBStorage User storage
type DBStorage struct {
	mu sync.RWMutex
	db *sql.DB
}

// NewStorage - User storage factory
func NewStorage(db *sql.DB) (*DBStorage, error) {
	storage := &DBStorage{db: db}

	err := storage.migrate(context.Background())
	if err != nil {
		return nil, err
	}

	return storage, nil
}

// NewUser - Create new user in DB
func (d *DBStorage) NewUser(ctx context.Context) (*domain.User, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	r := d.db.QueryRowContext(
		ctx, `INSERT INTO "user" (id) VALUES (DEFAULT) returning id;`,
	)
	if r.Err() != nil {
		return nil, r.Err()
	}

	user := &domain.User{}
	err := r.Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *DBStorage) migrate(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		if err != nil {
			_ = tx.Rollback()
		}
	}(tx)

	_, err = tx.ExecContext(ctx, `
	create table if not exists "user"
	(
		id    bigserial
			constraint user_pk
				primary key
	);`)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
