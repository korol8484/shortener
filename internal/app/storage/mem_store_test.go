package storage

import (
	"github.com/korol8484/shortener/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemStore_Add(t *testing.T) {
	store := NewMemStore()
	err := store.Add(&app.Entity{
		URL:   "http://www.ya.ru",
		Alias: "7A2S4z",
	})
	require.NoError(t, err)
}

func TestMemStore_Read(t *testing.T) {
	type want struct {
		alias string
		url   string
		err   bool
	}

	store := NewMemStore()
	err := store.Add(&app.Entity{
		URL:   "http://www.ya.ru",
		Alias: "7A2S4z",
	})
	require.NoError(t, err)

	tests := []struct {
		name string
		want want
	}{
		{
			name: "Success",
			want: want{
				alias: "7A2S4z",
				url:   "http://www.ya.ru",
			},
		},
		{
			name: "Success",
			want: want{
				err:   true,
				alias: "7A2S",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ent, err := store.Read(test.want.alias)
			if !test.want.err {
				require.NoError(t, err)
				assert.Equal(t, test.want.url, ent.URL)
				assert.Equal(t, test.want.alias, ent.Alias)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
