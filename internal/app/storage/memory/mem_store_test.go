package memory

import (
	"context"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/handlers"
	"github.com/korol8484/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemStore_Add(t *testing.T) {
	store := NewMemStore()

	defer func(store handlers.Store) {
		_ = store.Close()
	}(store)

	err := store.Add(context.Background(), &domain.URL{
		URL:   "http://www.ya.ru",
		Alias: "7A2S4z",
	})
	require.NoError(t, err)
}

func TestMemStore_Read(t *testing.T) {
	type want struct {
		alias string
		url   string
		err   error
	}

	store := NewMemStore()

	defer func(store handlers.Store) {
		_ = store.Close()
	}(store)

	err := store.Add(context.Background(), &domain.URL{
		URL:   "http://www.ya.ru",
		Alias: "7A2S4z",
	})
	require.NoError(t, err)

	tests := []struct {
		name string
		want want
	}{
		{
			name: "Success_Read_by_alias",
			want: want{
				alias: "7A2S4z",
				url:   "http://www.ya.ru",
			},
		},
		{
			name: "Alias_not_found",
			want: want{
				err:   storage.ErrNotFound,
				alias: "7A2S",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ent, err := store.Read(context.Background(), test.want.alias)
			if test.want.err == nil {
				require.NoError(t, err)
				assert.Equal(t, test.want.url, ent.URL)
				assert.Equal(t, test.want.alias, ent.Alias)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, test.want.err, err)
			}
		})
	}
}
