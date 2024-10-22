// Package storage errors for storage
package storage

import (
	"errors"
)

var (
	// ErrNotFound - Ошибка что запрошенные данные не найдены
	ErrNotFound = errors.New("can't find requested alias")
	// ErrIssetURL - Ошибка что запрошенные данные уже есть
	ErrIssetURL = errors.New("requested url isset")
)
