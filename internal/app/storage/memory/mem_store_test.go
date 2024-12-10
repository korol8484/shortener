package memory

import (
	"context"
	"github.com/korol8484/shortener/internal/app/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage"
)

func TestMemStore_Add(t *testing.T) {
	store := NewMemStore()

	defer func(store usecase.Store) {
		_ = store.Close()
	}(store)

	err := store.Add(context.Background(), &domain.URL{
		URL:   "http://www.ya.ru",
		Alias: "7A2S4z",
	}, &domain.User{ID: 1})
	require.NoError(t, err)

	err = store.Add(context.Background(), &domain.URL{
		URL:   "http://www.ya.ru",
		Alias: "7A2S4z",
	}, &domain.User{ID: 1})
	require.ErrorIs(t, storage.ErrIssetURL, err)

	err = store.AddBatch(context.Background(), domain.BatchURL{
		&domain.URL{
			URL:   "http://www.ya1.ru",
			Alias: "7A1S4z",
		},
	}, &domain.User{ID: 1})
	require.NoError(t, err)

	err = store.AddBatch(context.Background(), domain.BatchURL{
		&domain.URL{
			URL:   "http://www.ya1.ru",
			Alias: "7A1S4z",
		},
	}, &domain.User{ID: 1})
	require.ErrorIs(t, storage.ErrIssetURL, err)
}

func TestMemStore_Read(t *testing.T) {
	type want struct {
		alias string
		url   string
		err   error
	}

	store := NewMemStore()

	defer func(store usecase.Store) {
		_ = store.Close()
	}(store)

	err := store.Add(context.Background(), &domain.URL{
		URL:   "http://www.ya.ru",
		Alias: "7A2S4z",
	}, &domain.User{ID: 1})
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

func TestMemStore_ReadUserURL(t *testing.T) {
	store := NewMemStore()
	defer func(store usecase.Store) {
		_ = store.Close()
	}(store)

	user := &domain.User{ID: 1}

	err := store.Add(context.Background(), &domain.URL{URL: "http://www.ya.ru", Alias: "7A2S4z"}, user)
	require.NoError(t, err)

	userURL, err := store.ReadUserURL(context.Background(), user)
	require.NoError(t, err)
	require.Len(t, userURL, 1)

	assert.Equal(t, "http://www.ya.ru", userURL[0].URL)
	assert.Equal(t, "7A2S4z", userURL[0].Alias)

	userURL, err = store.ReadUserURL(context.Background(), &domain.User{ID: 2})

	require.NoError(t, err)
	require.Len(t, userURL, 0)
}

func TestMemStore_BatchDelete(t *testing.T) {
	store := NewMemStore()
	defer func(store usecase.Store) {
		_ = store.Close()
	}(store)

	user := &domain.User{ID: 1}

	err := store.Add(context.Background(), &domain.URL{URL: "http://www.ya.ru", Alias: "7A2S4z"}, user)
	require.NoError(t, err)

	err = store.BatchDelete(context.Background(), []string{"7A2S4z"}, user.ID)
	require.NoError(t, err)

	userURL, err := store.ReadUserURL(context.Background(), user)
	require.NoError(t, err)
	require.Len(t, userURL, 1)

	assert.Equal(t, "http://www.ya.ru", userURL[0].URL)
	assert.Equal(t, "7A2S4z", userURL[0].Alias)
	assert.Equal(t, true, userURL[0].Deleted)
}
