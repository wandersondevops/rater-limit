package limiter

import (
	"net/http"
	"time"

	"github.com/wandersondevops/rater-limit/rater-limit/limiter/storage"
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

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		token := r.Header.Get("API_KEY")

		if token != "" {
			if !rl.allow(token, rl.config.RateLimitToken) {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		} else {
			if !rl.allow(ip, rl.config.RateLimitIP) {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
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
