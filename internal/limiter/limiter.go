package limiter

import (
	"Rate-Limiter-Challenge-GoExpertPostGrad/pkg/persistence"
	"Rate-Limiter-Challenge-GoExpertPostGrad/pkg/ratelimiter"
	"context"
	"log"
)

type Limiter struct {
	client                 persistence.StoragePersistenceClient
	rateLimitIP            int
	rateLimitTokenDefault  int
	rateLimitTokenSpecific map[string]int
	blockTime              int
}

func NewLimiter(client persistence.StoragePersistenceClient, rateLimitIP, rateLimitTokenDefault int, rateLimitTokenSpecific map[string]int, blockTime int) *Limiter {
	return &Limiter{
		client:                 client,
		rateLimitIP:            rateLimitIP,
		rateLimitTokenDefault:  rateLimitTokenDefault,
		rateLimitTokenSpecific: rateLimitTokenSpecific,
		blockTime:              blockTime,
	}
}

func (l *Limiter) IsRateLimited(ctx context.Context, key string, limit int) bool {
	count, err := l.client.Get(ctx, key)
	if err != nil {
		log.Printf("Error getting count for key %s: %v", key, err)
		return false
	}
	return count >= limit
}

func (l *Limiter) IncrementRequestCount(ctx context.Context, key string) {
	_, err := l.client.Increment(ctx, key, l.blockTime)
	if err != nil {
		log.Printf("Error incrementing count for key %s: %v", key, err)
	}
}

func (l *Limiter) GetRateLimit(key string) int {
	if limit, exists := l.rateLimitTokenSpecific[key]; exists {
		return limit
	}
	return l.rateLimitTokenDefault
}

func (l *Limiter) GetIP() int {
	return l.rateLimitIP
}

var _ ratelimiter.RateLimiter = (*Limiter)(nil) // Ensures Limiter implements RateLimiter interface
