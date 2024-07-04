package storage

import (
	"fmt"
	"github.com/korol8484/shortener/internal/app"
)

type MemStore struct {
	items map[string]string
}

func NewMemStore() app.Store {
	return &MemStore{items: make(map[string]string)}
}

func (m *MemStore) Add(ent *app.Entity) error {
	m.items[ent.Alias] = ent.URL

	return nil
}

func (m *MemStore) Read(alias string) (*app.Entity, error) {
	if !m.hasAlias(alias) {
		return nil, fmt.Errorf("can't find requested alias %s", alias)
	}

	return &app.Entity{Alias: alias, URL: m.items[alias]}, nil
}

func (m *MemStore) hasAlias(alias string) bool {
	if _, ok := m.items[alias]; ok {
		return ok
	}

	return false
}
