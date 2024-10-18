package storage

import (
	"context"
	"sync/atomic"

	"github.com/korol8484/shortener/internal/app/domain"
)

type MemoryStore struct {
	id int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) NewUser(ctx context.Context) (*domain.User, error) {
	return &domain.User{
		ID: atomic.AddInt64(&m.id, 1),
	}, nil
}
