package usecase

import (
	"context"
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestGenAlias(t *testing.T) {
	alias := GenAlias(6, "testString")
	if alias != "Jlf8iW" {
		t.Fatal("invalid alias generated")
	}
}

func BenchmarkGenAlias(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenAlias(5, "https://ya.ru")
	}
}

func TestUsecase_Ping(t *testing.T) {
	uCase := NewUsecase(
		&config.App{BaseShortURL: "http://localhost"},
		memory.NewMemStore(),
		NewPingDummy(),
		zap.L(),
	)
	defer uCase.Close()

	require.Equal(t, uCase.Ping(), true)
}

func TestUsecase_CreateURL(t *testing.T) {
	uCase := NewUsecase(
		&config.App{BaseShortURL: "http://localhost"},
		memory.NewMemStore(),
		NewPingDummy(),
		zap.L(),
	)
	defer uCase.Close()

	u, err := uCase.CreateURL(context.Background(), "http://ya.ru", 1)
	require.NoError(t, err)

	assert.Equal(t, u.Alias, "zVFF0J")
	assert.Equal(t, u.URL, "http://ya.ru")

	_, err = uCase.LoadByAlias(context.Background(), "zVFF0J")
	require.NoError(t, err)

	_, err = uCase.CreateURL(context.Background(), "http__://ya.ru", 1)
	require.Error(t, err)
}

func TestUsecase_LoadAllUserURL(t *testing.T) {
	uCase := NewUsecase(
		&config.App{BaseShortURL: "http://localhost"},
		memory.NewMemStore(),
		NewPingDummy(),
		zap.L(),
	)
	defer uCase.Close()

	_, err := uCase.CreateURL(context.Background(), "http://ya.ru", 1)
	require.NoError(t, err)

	batch, err := uCase.LoadAllUserURL(context.Background(), 1)
	require.NoError(t, err)
	assert.Len(t, batch, 1)
}

func TestUsecase_GetStats(t *testing.T) {
	uCase := NewUsecase(
		&config.App{BaseShortURL: "http://localhost"},
		memory.NewMemStore(),
		NewPingDummy(),
		zap.L(),
	)
	defer uCase.Close()

	_, err := uCase.CreateURL(context.Background(), "http://ya.ru", 1)
	require.NoError(t, err)

	st, err := uCase.GetStats(context.Background())
	require.NoError(t, err)

	assert.Equal(t, st.Users, int64(1))
	assert.Equal(t, st.Urls, int64(1))
}

func TestUsecase_AddBatch(t *testing.T) {
	uCase := NewUsecase(
		&config.App{BaseShortURL: "http://localhost"},
		memory.NewMemStore(),
		NewPingDummy(),
		zap.L(),
	)
	defer uCase.Close()

	err := uCase.AddBatch(context.Background(), []*domain.URL{
		{
			URL:   "http://ya.ru",
			Alias: "wfW5C1",
		},
	}, 1)
	require.NoError(t, err)
}

func TestUsecase_FormatAlias(t *testing.T) {
	uCase := NewUsecase(
		&config.App{BaseShortURL: "http://localhost"},
		memory.NewMemStore(),
		NewPingDummy(),
		zap.L(),
	)
	defer uCase.Close()

	u, err := uCase.CreateURL(context.Background(), "http://ya.ru", 1)
	require.NoError(t, err)
	assert.Equal(t, uCase.FormatAlias(u), "http://localhost/zVFF0J")
}

func TestUsecase_AddToDelete(t *testing.T) {
	uCase := NewUsecase(
		&config.App{BaseShortURL: "http://localhost"},
		memory.NewMemStore(),
		NewPingDummy(),
		zap.L(),
	)
	defer uCase.Close()

	u, err := uCase.CreateURL(context.Background(), "http://ya.ru", 1)
	require.NoError(t, err)

	uCase.AddToDelete([]string{u.Alias}, 1)
	time.Sleep(time.Second)
}
