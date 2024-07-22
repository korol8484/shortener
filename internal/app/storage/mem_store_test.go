package storage

import (
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemStore_Add(t *testing.T) {
	store := NewMemStore()
	err := store.Add(&domain.URL{
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
	err := store.Add(&domain.URL{
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
				err:   NotFound,
				alias: "7A2S",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ent, err := store.Read(test.want.alias)
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
