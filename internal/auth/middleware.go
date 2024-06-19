package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type key int

var authKey key

func (a Auth) MiddlewareOptional(alternate http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				// No header
				if alternate != nil {
					alternate.ServeHTTP(w, r)
				} else {
					http.Error(w, "unauthorized", http.StatusForbidden)
				}
				return
			}
			const bearerPrefix = "Bearer "
			if len(header) >= len(bearerPrefix) && !strings.EqualFold(header[:len(bearerPrefix)], bearerPrefix) {
				// Invalid header
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}
			token := header[len(bearerPrefix):]
			uID, err := a.ParseToken(token)
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					http.Error(w, "token expired", http.StatusUnauthorized)
				} else {
					http.Error(w, "unauthorized", http.StatusForbidden)
				}
				return
			}

			ctx := context.WithValue(r.Context(), authKey, uID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func (a Auth) Middleware() func(http.Handler) http.Handler {
	return a.MiddlewareOptional(nil)
}

func MustFromContext(ctx context.Context) (uID string) {
	return ctx.Value(authKey).(string)
}

func FromContext(ctx context.Context) (uID string, ok bool) {
	if val := ctx.Value(authKey); val != nil {
		return val.(string), true
	}
	return
}
