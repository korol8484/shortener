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

// &config.App{FileStoragePath: os.TempDir() + "/test"}

type Config interface {
	GetStoragePath() string
}

type storeEntity struct {
	UUID  string `json:"uuid"`
	Alias string `json:"short_url"`
	URL   string `json:"original_url"`
}

type Store struct {
	mu        sync.RWMutex
	file      *os.File
	baseStore handlers.Store
}

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
		}); err != nil {
			return err
		}
	}

	return nil
}

func (f *Store) Add(ctx context.Context, ent *domain.URL) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	_, err := f.baseStore.Read(ctx, ent.Alias)
	if !errors.Is(err, storage.ErrNotFound) {
		return err
	}

	if err = f.save(ent); err != nil {
		return err
	}

	if err = f.baseStore.Add(ctx, ent); err != nil {
		return err
	}

	return nil
}

func (f *Store) AddBatch(ctx context.Context, batch domain.BatchURL) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, v := range batch {
		_, err := f.baseStore.Read(ctx, v.Alias)
		if !errors.Is(err, storage.ErrNotFound) {
			return err
		}

		if err = f.save(v); err != nil {
			return err
		}

		if err = f.baseStore.Add(ctx, v); err != nil {
			return err
		}
	}

	return nil
}

func (f *Store) Read(ctx context.Context, alias string) (*domain.URL, error) {
	return f.baseStore.Read(ctx, alias)
}

func (f *Store) ReadByURL(ctx context.Context, URL string) (*domain.URL, error) {
	return f.baseStore.ReadByURL(ctx, URL)
}

func (f *Store) Close() error {
	return f.file.Close()
}

func (f *Store) save(ent *domain.URL) error {
	v := &storeEntity{
		UUID:  uuid.NewString(),
		Alias: ent.Alias,
		URL:   ent.URL,
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
