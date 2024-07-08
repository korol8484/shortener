package storage

import "github.com/korol8484/shortener/internal/app/domain"

type Store interface {
	Add(ent *domain.Entity) error
	Read(alias string) (*domain.Entity, error)
}
