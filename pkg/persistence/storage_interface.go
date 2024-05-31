package persistence

import (
	"context"
	"time"
)

type StoragePersistenceClient interface {
	Set(ctx context.Context, key string, ttl time.Duration) error
	Get(ctx context.Context, key string) (int, error)
	Increment(ctx context.Context, key string, ttl time.Duration) (int, error)
	Exists(ctx context.Context, key string) (bool, error)
	IsBlocked(ctx context.Context, key string) (bool, error)
}
