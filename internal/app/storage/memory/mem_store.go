package memory

import (
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage"
	"sync"
)

type MemStore struct {
	mu    sync.RWMutex
	items map[string]string
}

func NewMemStore() *MemStore {
	store := &MemStore{
		items: make(map[string]string),
	}

	return store
}

func (m *MemStore) Add(ent *domain.URL) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[ent.Alias] = ent.URL

	return nil
}

func (m *MemStore) Read(alias string) (*domain.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.hasAlias(alias) {
		return nil, storage.ErrNotFound
	}

	return &domain.URL{Alias: alias, URL: m.items[alias]}, nil
}

func (m *MemStore) Close() error {
	return nil
}

func (m *MemStore) hasAlias(alias string) bool {
	if _, ok := m.items[alias]; ok {
		return ok
	}

	return false
}
