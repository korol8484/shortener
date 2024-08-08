package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/handlers"
	"github.com/korol8484/shortener/internal/app/storage"
	"net/url"
	"os"
	"path/filepath"
	"sync"
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

func NewFileStore(config Config, base handlers.Store) (handlers.Store, error) {
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

		if err := f.baseStore.Add(&domain.URL{
			URL:   v.URL,
			Alias: v.Alias,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (f *Store) Add(ent *domain.URL) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	_, err := f.baseStore.Read(ent.Alias)
	if !errors.Is(err, storage.ErrNotFound) {
		return err
	}

	if err = f.save(ent); err != nil {
		return err
	}

	if err = f.baseStore.Add(ent); err != nil {
		return err
	}

	return nil
}

func (f *Store) Read(alias string) (*domain.URL, error) {
	return f.baseStore.Read(alias)
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