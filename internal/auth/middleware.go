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

// MiddlewareOptional routes authenticated access to the authed handler and non-authenticated access notAuthed handler.
// If a user is authenticated, the request context will be populated with a value.
func (a Auth) MiddlewareFork(authed http.Handler, notAuthed http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			// No header
			if notAuthed != nil {
				notAuthed.ServeHTTP(w, r)
			} else {
				http.Error(w, "unauthorized", http.StatusForbidden)
			}
			return
		}
		const bearerPrefix = "Bearer "
		if len(header) < len(bearerPrefix) || !strings.EqualFold(header[:len(bearerPrefix)], bearerPrefix) {
			// Invalid header
			http.Error(w, "unauthorized", http.StatusForbidden)
			return
		}
		token := header[len(bearerPrefix):]
		user, err := a.ParseToken(token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "token expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "unauthorized", http.StatusForbidden)
			}
			return
		}

		ctx := context.WithValue(r.Context(), authKey, user)
		r = r.WithContext(ctx)
		authed.ServeHTTP(w, r)
	})
}

// MiddlewareOptional allows both authenticated and non-authenticated access to the provided route.
// It is recommended pair this with the [FromContext] method.
// See [MiddlewareFork] for more.
func (a Auth) MiddlewareOptional(next http.Handler) http.Handler {
	return a.MiddlewareFork(next, next)
}

func (a Auth) Middleware(next http.Handler) http.Handler {
	return a.MiddlewareFork(next, nil)
}

func MustFromContext(ctx context.Context) (user User) {
	return ctx.Value(authKey).(User)
}

func FromContext(ctx context.Context) (user User, ok bool) {
	if val := ctx.Value(authKey); val != nil {
		return val.(User), true
	}
	return
}
