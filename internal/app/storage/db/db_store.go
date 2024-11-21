package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage"
)

// Storage - Db storage
type Storage struct {
	mu sync.RWMutex
	db *sql.DB
}

// NewStorage - DB storage Factory
func NewStorage(db *sql.DB) (*Storage, error) {
	st := &Storage{db: db}

	err := st.migrate(context.Background())
	if err != nil {
		return nil, err
	}

	return st, nil
}

// Add save shorten URL
func (s *Storage) Add(ctx context.Context, ent *domain.URL, user *domain.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		if err != nil {
			_ = tx.Rollback()
		}
	}(tx)

	var id int64

	err = tx.QueryRowContext(
		ctx, `INSERT INTO shortener (url, alias) VALUES ($1,$2) ON CONFLICT (url) DO NOTHING RETURNING id`, ent.URL, ent.Alias,
	).Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	var isset bool
	if id < 1 {
		isset = true

		err = s.db.QueryRowContext(ctx, "SELECT t.id FROM shortener t WHERE url = $1", ent.URL).Scan(&id)
		if err != nil {
			return err
		}
	}

	_, err = tx.ExecContext(
		ctx, `INSERT INTO user_url (user_id, url_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, user.ID, id,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	if isset {
		return storage.ErrIssetURL
	}

	return nil
}

// AddBatch save shorten collection URL
func (s *Storage) AddBatch(ctx context.Context, batch domain.BatchURL, user *domain.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		if err != nil {
			_ = tx.Rollback()
		}
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

	insert := fmt.Sprintf("INSERT INTO shortener (url, alias) VALUES %s ON CONFLICT (url) DO UPDATE SET url=EXCLUDED.url RETURNING id", strings.Join(placeholders, ","))
	rows, err := tx.QueryContext(ctx, insert, vals...)
	if err != nil {
		return err
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	defer rows.Close()

	ids := make([]int64, 0, len(batch))
	for rows.Next() {
		var n int64

		if err = rows.Scan(&n); err != nil {
			return err
		}

		ids = append(ids, n)
	}

	var (
		userPlaceholders []string
		userVals         []interface{}
	)

	for i, v := range ids {
		userPlaceholders = append(userPlaceholders, fmt.Sprintf("($%d,$%d)",
			i*2+1,
			i*2+2,
		))

		userVals = append(userVals, user.ID, v)
	}

	insertUser := fmt.Sprintf("INSERT INTO user_url (user_id, url_id) VALUES %s ON CONFLICT DO NOTHING", strings.Join(userPlaceholders, ","))
	_, err = tx.ExecContext(ctx, insertUser, userVals...)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// BatchDelete delete shorten collection URL
func (s *Storage) BatchDelete(ctx context.Context, aliases []string, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		placeholders []string
		vals         []interface{}
	)

	vals = append(vals, userID)
	for i, v := range aliases {
		placeholders = append(placeholders, fmt.Sprintf("$%d",
			i+2,
		))

		vals = append(vals, v)
	}

	up := fmt.Sprintf(
		"UPDATE shortener s SET deleted = true FROM user_url uu WHERE s.id = uu.url_id AND uu.user_id = $1 AND s.alias IN (%s);",
		strings.Join(placeholders, ","),
	)

	_, err := s.db.ExecContext(ctx, up, vals...)
	if err != nil {
		return err
	}

	return nil
}

// ReadUserURL read user shorten URL
func (s *Storage) ReadUserURL(ctx context.Context, user *domain.User) (domain.BatchURL, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.QueryContext(ctx, "SELECT s.url, s.alias FROM shortener s INNER JOIN user_url uu on s.id = uu.url_id WHERE uu.user_id = $1;", user.ID)
	if err != nil {
		return nil, err
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	defer rows.Close()

	var batch domain.BatchURL
	for rows.Next() {
		var u domain.URL
		if err = rows.Scan(&u.URL, &u.Alias); err != nil {
			return nil, err
		}

		batch = append(batch, &domain.URL{
			URL:   u.URL,
			Alias: u.Alias,
		})
	}

	return batch, nil
}

// Read - read shorten URL
func (s *Storage) Read(ctx context.Context, alias string) (*domain.URL, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	row := s.db.QueryRowContext(ctx, "SELECT t.url, t.alias, t.deleted FROM shortener t WHERE alias = $1", alias)
	if row.Err() != nil {
		return nil, row.Err()
	}

	ent := &domain.URL{}
	err := row.Scan(&ent.URL, &ent.Alias, &ent.Deleted)
	if err != nil {
		return nil, err
	}

	return ent, nil
}

// ReadByURL read shorten URL by URL
func (s *Storage) ReadByURL(ctx context.Context, URL string) (*domain.URL, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	row := s.db.QueryRowContext(ctx, "SELECT t.url, t.alias FROM shortener t WHERE url = $1", URL)
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

// Close - close db
func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) migrate(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		if err != nil {
			_ = tx.Rollback()
		}
	}(tx)

	_, err = tx.ExecContext(ctx, `
	create table if not exists shortener
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

	_, err = tx.ExecContext(ctx, `create index if not exists shortener_alias_index on shortener (alias);`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `create unique index if not exists shortener_uidx_url ON shortener (url);`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `alter table shortener add IF NOT EXISTS deleted bool default false not null;`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `create table if not exists user_url
	(
		user_id bigserial
			constraint user_url_user_id_fk
				references "user"
				on delete cascade,
		url_id  bigserial
			constraint user_url_shortener_id_fk
				references shortener
				on delete cascade
	);`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `create unique index if not exists user_url_url_id_user_id_uindex on user_url (url_id, user_id);`)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
