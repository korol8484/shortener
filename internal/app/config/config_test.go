package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNewConfig(t *testing.T) {
	cfg, err := NewConfig()
	require.NoError(t, err)

	assert.Equal(t, "", cfg.GetDsn())
	assert.NotEmpty(t, cfg.GetStoragePath())
	assert.Equal(t, false, cfg.Https.Enable)
}
