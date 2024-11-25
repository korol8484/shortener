package file

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage"
	"github.com/korol8484/shortener/internal/app/storage/memory"
)

type StoreCfg string

func (s StoreCfg) GetStoragePath() string {
	return string(s)
}

func getStore(t *testing.T) (*Store, string) {
	storePath := path.Join(os.TempDir(), uuid.NewString())

	store, err := NewFileStore(StoreCfg(storePath), memory.NewMemStore())
	require.NoError(t, err)

	_, err = os.Stat(storePath)
	require.NoError(t, err)

	return store, storePath
}

func TestStore_Add(t *testing.T) {
	store, dPath := getStore(t)

	defer func() {
		_ = store.Close()
		_ = os.Remove(dPath)
	}()

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

func TestStore_Read(t *testing.T) {
	store, dPath := getStore(t)

	defer func() {
		_ = store.Close()
		_ = os.Remove(dPath)
	}()

	type want struct {
		alias string
		url   string
		err   error
	}
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

func TestStore_ReadUserURL(t *testing.T) {
	store, dPath := getStore(t)

	defer func() {
		_ = store.Close()
		_ = os.Remove(dPath)
	}()

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

func TestStore_ReadByURL(t *testing.T) {
	store, dPath := getStore(t)

	defer func() {
		_ = store.Close()
		_ = os.Remove(dPath)
	}()

	err := store.Add(context.Background(), &domain.URL{URL: "http://www.ya.ru", Alias: "7A2S4z"}, &domain.User{ID: 1})
	require.NoError(t, err)

	URL, err := store.ReadByURL(context.Background(), "http://www.ya.ru")
	require.NoError(t, err)

	assert.Equal(t, "http://www.ya.ru", URL.URL)
	assert.Equal(t, "7A2S4z", URL.Alias)

	_, err = store.ReadByURL(context.Background(), "http://www.ya1.ru")
	require.ErrorIs(t, err, storage.ErrNotFound)
}

func TestStore_BatchDelete(t *testing.T) {
	store, dPath := getStore(t)

	defer func() {
		_ = store.Close()
		_ = os.Remove(dPath)
	}()

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

func TestStore_save(t *testing.T) {
	store, dPath := getStore(t)

	err := store.save(&domain.URL{URL: "1", Alias: "1"}, &domain.User{ID: 1})
	require.NoError(t, err)

	err = store.Close()
	require.NoError(t, err)
	err = os.Remove(dPath)
	require.NoError(t, err)

	store, dPath = getStore(t)
	_ = os.Remove(dPath)
	_ = store.Close()

	err = store.save(&domain.URL{URL: "1", Alias: "1"}, &domain.User{ID: 1})
	require.Error(t, err)
}

func TestStore_AddBatch(t *testing.T) {
}

func Test_load(t *testing.T) {
	p := path.Join(os.TempDir(), uuid.NewString())
	f, err := os.Create(p)
	require.NoError(t, err)

	defer func(name string) {
		_ = os.Remove(name)
	}(p)

	_, err = f.Write([]byte("{\"uuid\":\"52edec03-edee-4600-999f-f5af452c29f0\",\"short_url\":\"7qfJga\",\"original_url\":\"http://www.ya1111111sdfdsfcccccfsdc11.ru\",\"user_id\":8}"))
	require.NoError(t, err)

	_ = f.Close()

	store, err := NewFileStore(StoreCfg(p), memory.NewMemStore())
	require.NoError(t, err)

	err = store.Close()
	require.NoError(t, err)
}
