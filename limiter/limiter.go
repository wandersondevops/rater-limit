package limiter

import (
	"time"

	"github.com/wandersondevops/rater-limit/limiter/storage"
)

type Config struct {
	RateLimitIP    int
	RateLimitToken int
	BlockTime      time.Duration
}

type RateLimiter struct {
	store  storage.Storage
	config Config
}

func NewRateLimiter(store storage.Storage, config Config) *RateLimiter {
	return &RateLimiter{store: store, config: config}
}

func (rl *RateLimiter) allow(identifier string, rateLimit int) bool {
	count, err := rl.store.Get(identifier)
	if err != nil {
		return false
	}

	if count >= rateLimit {
		rl.store.Block(identifier, rl.config.BlockTime)
		return false
	}

	rl.store.Increment(identifier)
	return true
}
