package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/wandersondevops/rater-limit/limiter"
	"github.com/wandersondevops/rater-limit/limiter/storage"
)

func setupRouter(rateLimiter *limiter.RateLimiter) *mux.Router {
	r := mux.NewRouter()
	r.Use(rateLimiter.Middleware)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome!"))
	}).Methods("GET")
	return r
}

func TestRateLimiterByIP(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	store := storage.NewRedisStorage(redisClient)

	rateLimiter := limiter.NewRateLimiter(store, limiter.Config{
		RateLimitIP:    2,
		RateLimitToken: 100,
		BlockTime:      1 * time.Minute,
	})

	router := setupRouter(rateLimiter)

	// Make two requests from the same IP
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 but got %v", w.Code)
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 but got %v", w.Code)
	}

	// Make a third request which should be rate limited
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429 but got %v", w.Code)
	}
}

func TestRateLimiterByToken(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	store := storage.NewRedisStorage(redisClient)

	rateLimiter := limiter.NewRateLimiter(store, limiter.Config{
		RateLimitIP:    100,
		RateLimitToken: 2,
		BlockTime:      1 * time.Minute,
	})

	router := setupRouter(rateLimiter)

	// Make two requests with the same token
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "test_token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 but got %v", w.Code)
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 but got %v", w.Code)
	}

	// Make a third request which should be rate limited
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429 but got %v", w.Code)
	}
}

func TestRateLimiterBlockTime(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	store := storage.NewRedisStorage(redisClient)

	rateLimiter := limiter.NewRateLimiter(store, limiter.Config{
		RateLimitIP:    2,
		RateLimitToken: 100,
		BlockTime:      1 * time.Second,
	})

	router := setupRouter(rateLimiter)

	// Make two requests from the same IP
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 but got %v", w.Code)
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 but got %v", w.Code)
	}

	// Make a third request which should be rate limited
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429 but got %v", w.Code)
	}

	// Wait for block time to expire
	time.Sleep(2 * time.Second)

	// Make another request which should be allowed
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 but got %v", w.Code)
	}
}
