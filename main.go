package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/wandersondevops/rater-limit/limiter"
	"github.com/wandersondevops/rater-limit/limiter/storage"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
	})

	// Create the limiter storage
	store := storage.NewRedisStorage(redisClient)

	// Initialize rate limiter
	rateLimiter := limiter.NewRateLimiter(store, limiter.Config{
		RateLimitIP:    getEnvAsInt("RATE_LIMIT_IP", 10),
		RateLimitToken: getEnvAsInt("RATE_LIMIT_TOKEN", 100),
		BlockTime:      getEnvAsDuration("BLOCK_TIME", 5*time.Minute),
	})

	// Router
	r := mux.NewRouter()

	// Apply the rate limiter middleware
	r.Use(rateLimiter.Middleware)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome!"))
	}).Methods("GET")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Helper functions to get environment variables
func getEnvAsInt(name string, defaultValue int) int {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(name string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
