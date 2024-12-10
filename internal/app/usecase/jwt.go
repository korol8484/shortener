package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/korol8484/shortener/internal/app/domain"
	"go.uber.org/zap"
	"time"
)

var (
	ErrInvalidToken = errors.New("token not valid")
)

// UserAddRepository - jwt user repository contract
type UserAddRepository interface {
	NewUser(ctx context.Context) (*domain.User, error)
}

// Jwt service
type Jwt struct {
	secret     string
	expire     time.Duration
	userRep    UserAddRepository
	signMethod jwt.SigningMethod
	tokenName  string
	logger     *zap.Logger
}

// Claims - jwt claims with userID
type Claims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id,omitempty"`
}

// NewJwt implements a simple middleware handler for adding JWT auth.
func NewJwt(userRep UserAddRepository, logger *zap.Logger, secret string) *Jwt {
	return &Jwt{
		secret:     secret,
		expire:     100 * time.Hour,
		userRep:    userRep,
		signMethod: jwt.SigningMethodHS256,
		tokenName:  "Authorization",
		logger:     logger,
	}
}

// GetTokenName - return token name
func (j *Jwt) GetTokenName() string {
	return j.tokenName
}

// CreateNewToken - create new token with user
func (j *Jwt) CreateNewToken(ctx context.Context) (*domain.User, string, error) {
	user, err := j.userRep.NewUser(ctx)
	if err != nil {
		return nil, "", err
	}

	token, err := j.buildJWTString(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// LoadClaims - load token claims with validate token
func (j *Jwt) LoadClaims(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != j.signMethod.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(j.secret), nil
		})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (j *Jwt) buildJWTString(user *domain.User) (string, error) {
	token := jwt.NewWithClaims(j.signMethod, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expire)),
		},
		UserID: user.ID,
	})

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
