package usecase

import (
	"context"
	"github.com/korol8484/shortener/internal/app/user/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestJwt_CreateNewToken(t *testing.T) {
	jwt := NewJwt(storage.NewMemoryStore(), zap.L(), "123")

	u, token, err := jwt.CreateNewToken(context.Background())
	require.NoError(t, err)

	cl, err := jwt.LoadClaims(token)
	require.NoError(t, err)

	assert.Equal(t, u.ID, cl.UserID)
}

func TestJwt_CreateNewTokenErr(t *testing.T) {
	jwt := NewJwt(storage.NewMemoryStore(), zap.L(), "123")
	_, err := jwt.LoadClaims("1234")
	require.Error(t, err)
}

func TestJwt_GetTokenName(t *testing.T) {
	jwt := NewJwt(storage.NewMemoryStore(), zap.L(), "123")
	assert.Equal(t, jwt.GetTokenName(), "Authorization")
}
