package middleware

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/user/util"
)

type UserAddRepository interface {
	NewUser(ctx context.Context) (*domain.User, error)
}

type Jwt struct {
	secret     string
	expire     time.Duration
	userRep    UserAddRepository
	signMethod jwt.SigningMethod
	tokenName  string
	logger     *zap.Logger
}

type claims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id,omitempty"`
	logger *zap.Logger
}

var (
	errInvalidToken = errors.New("token not valid")
)

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

func (j *Jwt) HandlerRead() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := j.tokenFromHeader(r)
			if token == "" {
				cToken, err := r.Cookie(j.tokenName)
				if err != nil {
					j.logger.Error("cookie not found")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				token = cToken.Value
			}

			claim, err := j.loadClaims(token)
			if err != nil {
				j.logger.Error("token not valid")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(util.SetUserIDToCtx(r.Context(), claim.UserID)))
		})
	}
}

func (j *Jwt) HandlerSet() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := j.tokenFromHeader(r)
			if token == "" {
				cToken, err := r.Cookie(j.tokenName)
				if err != nil && !errors.Is(err, http.ErrNoCookie) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if cToken != nil {
					token = cToken.Value
				}
			}

			if token == "" {
				user, err := j.setNewToken(r.Context(), w)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				next.ServeHTTP(w, r.WithContext(util.SetUserIDToCtx(r.Context(), user.ID)))
				return
			}

			claim, err := j.loadClaims(token)
			if err != nil {
				if !errors.Is(err, errInvalidToken) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				user, err := j.setNewToken(r.Context(), w)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				next.ServeHTTP(w, r.WithContext(util.SetUserIDToCtx(r.Context(), user.ID)))
				return
			}

			next.ServeHTTP(w, r.WithContext(util.SetUserIDToCtx(r.Context(), claim.UserID)))
		})
	}
}

func (j *Jwt) setNewToken(ctx context.Context, w http.ResponseWriter) (*domain.User, error) {
	user, err := j.userRep.NewUser(ctx)
	if err != nil {
		return nil, err
	}

	token, err := j.buildJWTString(user)
	if err != nil {
		return nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     j.tokenName,
		Value:    token,
		Path:     "/",
		Secure:   false,
		HttpOnly: false,
	})

	w.Header().Add(j.tokenName, token)
	w.Header().Set(j.tokenName, token)

	return user, nil
}

func (j *Jwt) buildJWTString(user *domain.User) (string, error) {
	token := jwt.NewWithClaims(j.signMethod, claims{
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

func (j *Jwt) loadClaims(tokenStr string) (*claims, error) {
	claims := &claims{}
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
		return nil, errInvalidToken
	}

	return claims, nil
}

func (j *Jwt) tokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get(j.tokenName)
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}
