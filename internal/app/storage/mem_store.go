package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/korol8484/shortener/internal/app/domain"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

var (
	// ErrNotFound - Ошибка что запрошенные данные не найдены
	ErrNotFound = errors.New("can't find requested alias")
)

type Config interface {
	GetStoragePath() string
}

type storeEntity struct {
	Uuid  string `json:"uuid"`
	Alias string `json:"short_url"`
	URL   string `json:"original_url"`
}

type MemStore struct {
	mu    sync.RWMutex
	items map[string]string
	file  *os.File
}

func NewMemStore(config Config) (Store, error) {
	file, err := create(config.GetStoragePath())
	if err != nil {
		return nil, err
	}

	store := &MemStore{
		items: make(map[string]string),
		file:  file,
	}

	if err = store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

func (m *MemStore) load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	scanner := bufio.NewScanner(m.file)
	for scanner.Scan() {
		v := &storeEntity{}
		if err := json.Unmarshal(scanner.Bytes(), v); err != nil {
			return err
		}

		if _, err := url.Parse(v.URL); err != nil {
			return err
		}

		m.items[v.Alias] = v.URL
	}

	return nil
}

func (m *MemStore) Add(ent *domain.URL) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.hasAlias(ent.Alias) {
		if err := m.save(ent); err != nil {
			return err
		}
	}

	m.items[ent.Alias] = ent.URL

	return nil
}

func (m *MemStore) Read(alias string) (*domain.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.hasAlias(alias) {
		return nil, ErrNotFound
	}

	return &domain.URL{Alias: alias, URL: m.items[alias]}, nil
}

func (m *MemStore) Close() error {
	return m.file.Close()
}

func (m *MemStore) save(ent *domain.URL) error {
	v := &storeEntity{
		Uuid:  uuid.NewString(),
		Alias: ent.Alias,
		URL:   ent.URL,
	}

	b, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	if _, err = m.file.Write(b); err != nil {
		return err
	}

	if _, err = m.file.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

func (m *MemStore) hasAlias(alias string) bool {
	if _, ok := m.items[alias]; ok {
		return ok
	}

	return false
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}

	return os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}
