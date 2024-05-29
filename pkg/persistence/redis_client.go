package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
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

// Get retrieves a value from Redis and converts it to an integer.
func (r *RedisClient) Get(ctx context.Context, key string) (int, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("key does not exist")
	} else if err != nil {
		return 0, err
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return intVal, nil
}

// Increment increments a value in Redis and sets the TTL if provided.
func (r *RedisClient) Increment(ctx context.Context, key string, ttl int) (int, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if ttl > 0 {
		err = r.client.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
		if err != nil {
			return 0, err
		}
	}
	return int(val), nil
}

// Expire sets the TTL for a key in Redis.
func (r *RedisClient) Expire(ctx context.Context, key string, ttl int) error {
	return r.client.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
}

// Set sets a key in Redis with an initial value and TTL.
func (r *RedisClient) Set(ctx context.Context, key string, value int, ttl int) error {
	err := r.client.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

// Exists checks if a key exists in Redis.
func (r *RedisClient) Exists(ctx context.Context, key string) (int, error) {
	val, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(val), nil
}
