package db

import (
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewPgDB(t *testing.T) {
	_, err := NewPgDB(&config.App{DBDsn: ""})
	require.NoError(t, err)
}
