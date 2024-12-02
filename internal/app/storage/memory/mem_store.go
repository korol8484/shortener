package memory

import (
	"context"
	"sync"

	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage"
)

type item struct {
	URL     string
	deleted bool
}

// MemStore - in memory shorten links storage
type MemStore struct {
	mu           sync.RWMutex
	items        map[string]item
	userItems    map[int64][]string
	deletedItems map[int64]map[string]interface{}
}

// NewMemStore in memory shorten links storage factory
func NewMemStore() *MemStore {
	store := &MemStore{
		items:        make(map[string]item),
		userItems:    make(map[int64][]string),
		deletedItems: make(map[int64]map[string]interface{}),
	}

	return store
}

// Add save shorten URL
func (m *MemStore) Add(ctx context.Context, ent *domain.URL, user *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.items[ent.Alias]; ok {
		return storage.ErrIssetURL
	}

	m.items[ent.Alias] = item{URL: ent.URL}
	m.userItems[user.ID] = append(m.userItems[user.ID], ent.Alias)

	return nil
}

// AddBatch save shorten collection URL
func (m *MemStore) AddBatch(ctx context.Context, batch domain.BatchURL, user *domain.User) error {
	for _, v := range batch {
		if err := m.Add(ctx, v, user); err != nil {
			return err
		}
	}

	return nil
}

// ReadUserURL read user shorten URL
func (m *MemStore) ReadUserURL(ctx context.Context, user *domain.User) (domain.BatchURL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var batch domain.BatchURL

	aliases, ok := m.userItems[user.ID]
	if !ok {
		return batch, nil
	}

	deleted, hasDel := m.deletedItems[user.ID]

	for _, alias := range aliases {
		u := m.items[alias]

		URL := &domain.URL{
			URL:   u.URL,
			Alias: alias,
		}

		if hasDel {
			if _, ok := deleted[alias]; ok {
				URL.Deleted = true
			}
		}

		batch = append(batch, URL)
	}

	return batch, nil
}

// Read - read shorten URL
func (m *MemStore) Read(ctx context.Context, alias string) (*domain.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.hasAlias(alias) {
		return nil, storage.ErrNotFound
	}

	return &domain.URL{Alias: alias, URL: m.items[alias].URL, Deleted: m.items[alias].deleted}, nil
}

// ReadByURL read shorten URL by URL
func (m *MemStore) ReadByURL(ctx context.Context, URL string) (*domain.URL, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.items {
		if v.URL == URL {
			return &domain.URL{
				URL:   v.URL,
				Alias: k,
			}, nil
		}
	}

	return nil, storage.ErrNotFound
}

// BatchDelete delete shorten collection URL
func (m *MemStore) BatchDelete(ctx context.Context, aliases []string, userID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.deletedItems[userID]; !ok {
		m.deletedItems[userID] = make(map[string]interface{}, len(aliases))
	}

	for _, alias := range aliases {
		m.deletedItems[userID][alias] = nil
	}

	return nil
}

// LoadStats load url items, user count
func (m *MemStore) LoadStats(ctx context.Context) (*domain.StatsModel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &domain.StatsModel{
			Urls:  int64(len(m.items)),
			Users: int64(len(m.userItems)),
		}, nil
	}
}

// Close - clear store resources
func (m *MemStore) Close() error {
	return nil
}

func (m *MemStore) hasAlias(alias string) bool {
	if _, ok := m.items[alias]; ok {
		return ok
	}

	return false
}
