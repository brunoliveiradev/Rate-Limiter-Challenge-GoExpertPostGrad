package persistence

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisClient_Set(t *testing.T) {
	redisClient := NewRedisClient("localhost:6379", "", 0)
	ctx := context.Background()

	redisClient.client.Del(ctx, "testKey")

	err := redisClient.Set(ctx, "testKey", "1", 10)
	assert.NoError(t, err)

	exists, err := redisClient.Exists(ctx, "testKey")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRedisClient_Get(t *testing.T) {
	redisClient := NewRedisClient("localhost:6379", "", 0)
	ctx := context.Background()

	_, err := redisClient.Get(ctx, "nonexistentKey")
	assert.Error(t, err)

	redisClient.Set(ctx, "testKey", "1", 10)
	val, err := redisClient.Get(ctx, "testKey")
	assert.NoError(t, err)
	assert.Equal(t, "1", val)
}

func TestRedisClient_Increment(t *testing.T) {
	redisClient := NewRedisClient("localhost:6379", "", 0)
	ctx := context.Background()

	redisClient.client.Del(ctx, "testKey")

	val, err := redisClient.Increment(ctx, "testKey", 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	val, err = redisClient.Increment(ctx, "testKey", 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, val)
}

func TestRedisClient_Exists(t *testing.T) {
	redisClient := NewRedisClient("localhost:6379", "", 0)
	ctx := context.Background()

	exists, err := redisClient.Exists(ctx, "nonexistentKey")
	assert.NoError(t, err)
	assert.False(t, exists)

	redisClient.Set(ctx, "testKey", "1", 10)
	exists, err = redisClient.Exists(ctx, "testKey")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRedisClient_IsBlocked(t *testing.T) {
	redisClient := NewRedisClient("localhost:6379", "", 0)
	ctx := context.Background()

	isBlocked, err := redisClient.IsKeyBlocked(ctx, "nonexistentKey")
	assert.NoError(t, err)
	assert.False(t, isBlocked)

	redisClient.Set(ctx, "testKey:blocked", "1", 10)
	isBlocked, err = redisClient.IsKeyBlocked(ctx, "testKey")
	assert.NoError(t, err)
	assert.True(t, isBlocked)
}
