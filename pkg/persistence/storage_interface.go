package persistence

import (
	"context"
)

type StoragePersistenceClient interface {
	Set(ctx context.Context, key string, value string, ttl int) error
	Get(ctx context.Context, key string) (string, error)
	Increment(ctx context.Context, key string, ttl int) (int, error)
	Exists(ctx context.Context, key string) (bool, error)
	IsKeyBlocked(ctx context.Context, key string) (bool, error)
	FlushAll(ctx context.Context) error
}
