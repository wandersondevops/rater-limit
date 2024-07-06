package limiter

import (
	"net/http"
)

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
