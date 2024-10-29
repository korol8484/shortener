package storage

import (
	"context"
	"sync/atomic"

	"github.com/korol8484/shortener/internal/app/domain"
)

// MemoryStore struct
type MemoryStore struct {
	id int64
}

// NewMemoryStore - Factory for user memory storage
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

// NewUser - Create new user in memory
func (m *MemoryStore) NewUser(ctx context.Context) (*domain.User, error) {
	return &domain.User{
		ID: atomic.AddInt64(&m.id, 1),
	}, nil
}
