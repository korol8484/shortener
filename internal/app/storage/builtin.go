package storage

import (
	"errors"
)

var (
	// ErrNotFound - Ошибка что запрошенные данные не найдены
	ErrNotFound = errors.New("can't find requested alias")
	ErrIssetURL = errors.New("requested url isset")
)
