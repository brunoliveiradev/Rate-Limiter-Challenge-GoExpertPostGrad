package persistence

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisClient{
		client: rdb,
	}
}

// Set sets a key in Redis with a value and TTL.
func (r *RedisClient) Set(ctx context.Context, key string, value string, ttl int) error {
	return r.client.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err()
}

// Get gets a value from Redis by key.
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Increment increments a value in Redis and sets the TTL if provided.
func (r *RedisClient) Increment(ctx context.Context, key string, ttl int) (int, error) {
	txPipeline := r.client.TxPipeline()

	incr := txPipeline.Incr(ctx, key)

	// Check if the key exists before setting the TTL
	exists, err := r.Exists(ctx, key)
	if err != nil {
		return 0, err
	}

	// If the key does not exist, set the TTL
	if !exists {
		txPipeline.Expire(ctx, key, time.Duration(ttl)*time.Second)
	}

	_, err = txPipeline.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return int(incr.Val()), nil
}

// Exists checks if a key exists in Redis. Returns 1 (true) if the key exists, 0 (false) if it does not
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	val, err := r.client.Exists(ctx, key).Result()
	return val == 1, err
}

// IsKeyBlocked checks if a key is blocked in Redis. Returns true if the key is blocked, false if it does not.
func (r *RedisClient) IsKeyBlocked(ctx context.Context, key string) (bool, error) {
	return r.Exists(ctx, key)
}

func (r *RedisClient) FlushAll(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}
