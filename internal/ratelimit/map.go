package ratelimit

import (
	"math"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type mapRateLimiter[K comparable] struct {
	keyMap map[K]*mapLimiter
	mu     sync.RWMutex
	r      rate.Limit
	b      int
}

func NewMapRateLimiter[K comparable](r rate.Limit, b int) *mapRateLimiter[K] {
	return &mapRateLimiter[K]{
		keyMap: make(map[K]*mapLimiter),
		mu:     sync.RWMutex{},
		r:      r,
		b:      b,
	}
}

// mapLimiter represents the limiter object with metadata
type mapLimiter struct {
	rate.Limiter

	// Only used for TTL purposes. Does not need to be precise.
	lastAccessed time.Time
}

func newLimiter(r rate.Limit, b int) *mapLimiter {
	return &mapLimiter{
		Limiter:      *rate.NewLimiter(r, b),
		lastAccessed: time.Now(),
	}
}

// Get safely returns a limiter from key map
func (rl *mapRateLimiter[K]) Get(key K) (*mapLimiter, bool) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	val, ok := rl.keyMap[key]
	return val, ok
}

func (rl *mapRateLimiter[K]) Access(key K) *mapLimiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, ok := rl.keyMap[key]
	if !ok || limiter == nil {
		limiter = newLimiter(rl.r, rl.b)
		rl.keyMap[key] = limiter
		go rl.DelLimiterTTL(key)
	}

	limiter.lastAccessed = time.Now()
	return limiter
}

// Set safely sets a value in key map
func (rl *mapRateLimiter[K]) Set(key K, val *mapLimiter) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.keyMap[key] = val
}

// Del safely deletes a key in key map
func (rl *mapRateLimiter[K]) Del(key K) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.keyMap, key)
}

func (rl *mapRateLimiter[K]) Allow(key K) bool {
	limiter := rl.Access(key)
	return limiter.Allow()
}

func (rl *mapRateLimiter[K]) DelLimiterTTL(key K) {
	// TTL (time to live) - prevent wasting memory as server scales
	clearAfterDuration := time.Second * time.Duration(math.Ceil(1/float64(rl.r)))

	for {
		<-time.After(clearAfterDuration)

		limiter, ok := rl.Get(key)
		if !ok {
			return
		}

		if time.Since(limiter.lastAccessed) >= clearAfterDuration {
			rl.Del(key)
			return
		}
	}
}
