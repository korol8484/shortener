package file

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"

	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/handlers"
	"github.com/korol8484/shortener/internal/app/storage"
)

// Config - file storage config interface
type Config interface {
	GetStoragePath() string
}

type storeEntity struct {
	UUID   string `json:"uuid"`
	Alias  string `json:"short_url"`
	URL    string `json:"original_url"`
	UserID int64  `json:"user_id,omitempty"`
}

// Store - file storage
type Store struct {
	mu        sync.RWMutex
	file      *os.File
	baseStore handlers.Store
}

// NewFileStore - file storage factory
func NewFileStore(config Config, base handlers.Store) (*Store, error) {
	file, err := create(config.GetStoragePath())
	if err != nil {
		return nil, err
	}

	store := &Store{
		file:      file,
		baseStore: base,
	}

	if err = store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

func (f *Store) load() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		v := &storeEntity{}
		if err := json.Unmarshal(scanner.Bytes(), v); err != nil {
			return err
		}

		if _, err := url.Parse(v.URL); err != nil {
			return err
		}

		if err := f.baseStore.Add(context.Background(), &domain.URL{
			URL:   v.URL,
			Alias: v.Alias,
		}, &domain.User{ID: v.UserID}); err != nil {
			return err
		}
	}

	return nil
}

// Add save shorten URL
func (f *Store) Add(ctx context.Context, ent *domain.URL, user *domain.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	_, err := f.baseStore.Read(ctx, ent.Alias)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return err
	}

	if err = f.save(ent, user); err != nil {
		return err
	}

	if err = f.baseStore.Add(ctx, ent, user); err != nil {
		return err
	}

	return nil
}

// AddBatch save shorten collection URL
func (f *Store) AddBatch(ctx context.Context, batch domain.BatchURL, user *domain.User) error {
	for _, v := range batch {
		if err := f.Add(ctx, v, user); err != nil {
			return err
		}
	}

	return nil
}

// ReadUserURL read user shorten URL
func (f *Store) ReadUserURL(ctx context.Context, user *domain.User) (domain.BatchURL, error) {
	return f.baseStore.ReadUserURL(ctx, user)
}

// Read - read shorten URL
func (f *Store) Read(ctx context.Context, alias string) (*domain.URL, error) {
	return f.baseStore.Read(ctx, alias)
}

// ReadByURL read shorten URL by URL
func (f *Store) ReadByURL(ctx context.Context, URL string) (*domain.URL, error) {
	return f.baseStore.ReadByURL(ctx, URL)
}

// BatchDelete delete shorten collection URL
func (f *Store) BatchDelete(ctx context.Context, aliases []string, userID int64) error {
	return f.baseStore.BatchDelete(ctx, aliases, userID)
}

// Close - close file store
func (f *Store) Close() error {
	return f.file.Close()
}

func (f *Store) save(ent *domain.URL, user *domain.User) error {
	v := &storeEntity{
		UUID:   uuid.NewString(),
		Alias:  ent.Alias,
		URL:    ent.URL,
		UserID: user.ID,
	}

	b, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	if _, err = f.file.Write(b); err != nil {
		return err
	}

	if _, err = f.file.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}

	return os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}
