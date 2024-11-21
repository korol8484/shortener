package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg := &App{
		Listen:          ":8080",
		BaseShortURL:    "testBase",
		FileStoragePath: "fileStore",
		DBDsn:           "dbDsn",
	}

	assert.Equal(t, "dbDsn", cfg.GetDsn())
	assert.Equal(t, "fileStore", cfg.GetStoragePath())
	assert.Equal(t, "testBase", cfg.GetBaseShortURL())
}
