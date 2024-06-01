package ratelimiter

import "context"

type LimiterInterface interface {
	RateLimited(ctx context.Context, ip string, token string) error
}
