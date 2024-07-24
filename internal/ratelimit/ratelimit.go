package ratelimit

type RateLimiter[K comparable] interface {
	Allow(key K) bool
}
