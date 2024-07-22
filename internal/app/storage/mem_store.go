package storage

import (
	"fmt"
	"github.com/korol8484/shortener/internal/app/domain"
)

type MemStore struct {
	items map[string]string
}

func NewMemStore() Store {
	return &MemStore{items: make(map[string]string)}
}

func (m *MemStore) Add(ent *domain.Entity) error {
	m.items[ent.Alias] = ent.URL

	return nil
}

func (m *MemStore) Read(alias string) (*domain.Entity, error) {
	if !m.hasAlias(alias) {
		return nil, fmt.Errorf("can't find requested alias %s", alias)
	}

	return &domain.Entity{Alias: alias, URL: m.items[alias]}, nil
}

func (m *MemStore) hasAlias(alias string) bool {
	if _, ok := m.items[alias]; ok {
		return ok
	}

	return false
}
