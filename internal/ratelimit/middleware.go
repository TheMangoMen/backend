package ratelimit

import (
	"net/http"

	"github.com/TheMangoMen/backend/internal/auth"
)

// Middleware represents the rate limit validator, will allow requests with valid rate limits to pass
func AuthedMiddleware(rl RateLimiter[string]) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := auth.MustFromContext(r.Context())
			key := user.UID

			if !rl.Allow(key) {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
