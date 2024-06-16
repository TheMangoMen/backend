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

func (a Auth) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("Authorization")
			if err != nil {
				// No cookie
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}
			const bearerPrefix = "Bearer "
			if len(cookie.Value) >= len(bearerPrefix) && !strings.EqualFold(cookie.Value[:len(bearerPrefix)], bearerPrefix) {
				// Invalid cookie
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}
			token := cookie.Value[len(bearerPrefix):]
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

func FromContext(ctx context.Context) (uID string) {
	return ctx.Value(authKey).(string)
}
