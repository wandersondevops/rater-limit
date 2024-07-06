package storage

import "time"

type Storage interface {
	Get(key string) (int, error)
	Increment(key string) error
	Block(key string, duration time.Duration) error
}
