package middleware

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"strings"

	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/usecase"
	"github.com/korol8484/shortener/internal/app/user/util"
)

// Jwt middleware
type Jwt struct {
	uc     *usecase.Jwt
	logger *zap.Logger
}

// NewJwt implements a simple middleware handler for adding JWT auth.
func NewJwt(uc *usecase.Jwt, logger *zap.Logger) *Jwt {
	return &Jwt{
		uc:     uc,
		logger: logger,
	}
}

// HandlerRead returns a middleware that will read and validate JWT token from header or cookie
func (j *Jwt) HandlerRead() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := j.tokenFromHeader(r)
			if token == "" {
				cToken, err := r.Cookie(j.uc.GetTokenName())
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

// HandlerSet returns a middleware that will set JWT token if not exist to header and cookie
func (j *Jwt) HandlerSet() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := j.tokenFromHeader(r)
			if token == "" {
				cToken, err := r.Cookie(j.uc.GetTokenName())
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
				if !errors.Is(err, usecase.ErrInvalidToken) {
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
	user, token, err := j.uc.CreateNewToken(ctx)
	if err != nil {
		return nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     j.uc.GetTokenName(),
		Value:    token,
		Path:     "/",
		Secure:   false,
		HttpOnly: false,
	})

	w.Header().Set(j.uc.GetTokenName(), token)

	return user, nil
}

func (j *Jwt) loadClaims(tokenStr string) (*usecase.Claims, error) {
	return j.uc.LoadClaims(tokenStr)
}

func (j *Jwt) tokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get(j.uc.GetTokenName())
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}
