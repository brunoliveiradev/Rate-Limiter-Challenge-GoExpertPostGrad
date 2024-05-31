package persistence

import (
	"context"
	"errors"
	"fmt"
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
func (r *RedisClient) Set(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Set(ctx, key, true, ttl).Err()
}

// Get retrieves a value from Redis and converts it to an integer.
func (r *RedisClient) Get(ctx context.Context, key string) (int, error) {
	val, err := r.client.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("key %s does not exist", key)
	}
	return val, err
}

// Increment increments a value in Redis and sets the TTL if provided.
func (r *RedisClient) Increment(ctx context.Context, key string, ttl time.Duration) (int, error) {
	pipe := r.client.TxPipeline()

	incr := pipe.Incr(ctx, key)
	exists := pipe.Exists(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	if exists.Val() == 1 && incr.Val() == 1 {
		err = r.client.Expire(ctx, key, ttl*time.Second).Err()
		if err != nil {
			return 0, err
		}
	}

	return int(incr.Val()), nil
}

// Exists checks if a key exists in Redis. Returns 1 (true) if the key exists, 0 (false) if it does not
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	val, err := r.client.Exists(ctx, key).Result()
	return val == 1, err
}

func (r *RedisClient) IsBlocked(ctx context.Context, key string) (bool, error) {
	return r.Exists(ctx, key)
}
