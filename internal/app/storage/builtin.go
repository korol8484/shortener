package storage

import (
	"errors"
)

var (
	// ErrNotFound - Ошибка что запрошенные данные не найдены
	ErrNotFound = errors.New("can't find requested alias")
	ErrIssetUrl = errors.New("requested url isset")
)
