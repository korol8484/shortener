package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	cookieName string
}

type claims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id,omitempty"`
}

var (
	errIvalidToken = errors.New("token not valid")
)

func NewJwt(userRep UserAddRepository) *Jwt {
	return &Jwt{
		secret:     "12345dsdsdtoken",
		expire:     3 * time.Hour,
		userRep:    userRep,
		signMethod: jwt.SigningMethodHS256,
		cookieName: "token",
	}
}

func (j *Jwt) HandlerRead() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(j.cookieName)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claim, err := j.loadClaims(cookie.Value)
			if err != nil {
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
			cookie, err := r.Cookie(j.cookieName)
			if err != nil && !errors.Is(err, http.ErrNoCookie) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if cookie == nil || cookie.Value == "" {
				user, err := j.setNewCookie(r.Context(), w)
				fmt.Println(err)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				next.ServeHTTP(w, r.WithContext(util.SetUserIDToCtx(r.Context(), user.ID)))
				return
			}

			claim, err := j.loadClaims(cookie.Value)
			if err != nil {
				if !errors.Is(err, errIvalidToken) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				user, err := j.setNewCookie(r.Context(), w)
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

func (j *Jwt) setNewCookie(ctx context.Context, w http.ResponseWriter) (*domain.User, error) {
	user, err := j.userRep.NewUser(ctx)
	if err != nil {
		return nil, err
	}

	token, err := j.buildJWTString(user)
	if err != nil {
		return nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     j.cookieName,
		Value:    token,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

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
		return nil, errIvalidToken
	}

	return claims, nil
}
