package persistence

import "context"

type StoragePersistenceClient interface {
	Get(ctx context.Context, key string) (int, error)
	Increment(ctx context.Context, key string, ttl int) (int, error)
	Expire(ctx context.Context, key string, ttl int) error
	Set(ctx context.Context, key string, value int, ttl int) error
	Exists(ctx context.Context, key string) (int, error)
}
