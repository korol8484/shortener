package app

import (
	"errors"
	"fmt"
)

type MemStore struct {
	items map[string]string
}

func NewMemStore() Store {
	return &MemStore{items: make(map[string]string)}
}

func (m *MemStore) Add(ent *Entity) error {
	m.items[ent.Alias] = ent.Url

	return nil
}

func (m *MemStore) Read(alias string) (*Entity, error) {
	if !m.hasAlias(alias) {
		return nil, errors.New(fmt.Sprintf("can't find requested alias %s", alias))
	}

	return &Entity{Alias: alias, Url: m.items[alias]}, nil
}

func (m *MemStore) hasAlias(alias string) bool {
	if _, ok := m.items[alias]; ok {
		return ok
	}

	return false
}
