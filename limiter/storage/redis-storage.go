package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

func (s *RedisStorage) Get(key string) (int, error) {
	val, err := s.client.Get(context.Background(), key).Int()
	if err != nil && err != redis.Nil {
		return 0, err
	}
	return val, nil
}

func (s *RedisStorage) Increment(key string) error {
	_, err := s.client.Incr(context.Background(), key).Result()
	return err
}

func (s *RedisStorage) Block(key string, duration time.Duration) error {
	_, err := s.client.Set(context.Background(), key, "blocked", duration).Result()
	return err
}
