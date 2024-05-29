package ratelimiter

import (
	"context"
)

type RateLimiter interface {
	IsRateLimited(ctx context.Context, key string, limit int) bool
	IncrementRequestCount(ctx context.Context, key string)
	GetRateLimit(key string) int
	GetIP() int
}
