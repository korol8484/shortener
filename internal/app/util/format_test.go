package util

import (
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddURLToAlias(t *testing.T) {
	URL := &domain.URL{
		URL:   "1",
		Alias: "1",
	}
	res := URL.FormatAlias(AddURLToAlias("2"))

	assert.Equal(t, res, "2/1")
}
