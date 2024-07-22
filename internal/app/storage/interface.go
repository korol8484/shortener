package storage

import "github.com/korol8484/shortener/internal/app/domain"

type Store interface {
	Add(ent *domain.URL) error
	Read(alias string) (*domain.URL, error)
}
