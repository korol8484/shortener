package memory

import (
	"context"
	"sync"

	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage"
)

type MemStore struct {
	mu        sync.RWMutex
	items     map[string]string
	userItems map[int64][]string
}

func NewMemStore() *MemStore {
	store := &MemStore{
		items:     make(map[string]string),
		userItems: make(map[int64][]string),
	}

	return store
}

func (m *MemStore) Add(ctx context.Context, ent *domain.URL, user *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[ent.Alias] = ent.URL
	m.userItems[user.ID] = append(m.userItems[user.ID], ent.Alias)

	return nil
}

func (m *MemStore) AddBatch(ctx context.Context, batch domain.BatchURL, user *domain.User) error {
	for _, v := range batch {
		if err := m.Add(ctx, v, user); err != nil {
			return err
		}
	}

	return nil
}

func (m *MemStore) ReadUserUrl(ctx context.Context, user *domain.User) (domain.BatchURL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var batch domain.BatchURL

	aliases, ok := m.userItems[user.ID]
	if !ok {
		return batch, nil
	}

	for _, alias := range aliases {
		u := m.items[alias]

		batch = append(batch, &domain.URL{
			URL:   u,
			Alias: alias,
		})
	}

	return batch, nil
}

func (m *MemStore) Read(ctx context.Context, alias string) (*domain.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.hasAlias(alias) {
		return nil, storage.ErrNotFound
	}

	return &domain.URL{Alias: alias, URL: m.items[alias]}, nil
}

func (m *MemStore) ReadByURL(ctx context.Context, URL string) (*domain.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.items {
		if v == URL {
			return &domain.URL{
				URL:   v,
				Alias: k,
			}, nil
		}
	}

	return nil, storage.ErrNotFound
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
