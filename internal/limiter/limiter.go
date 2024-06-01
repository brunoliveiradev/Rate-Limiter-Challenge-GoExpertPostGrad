package limiter

import (
	"Rate-Limiter-Challenge-GoExpertPostGrad/pkg/persistence"
	"context"
	"errors"
)

var ErrNotAllowed = errors.New("access not allowed")

type Limiter struct {
	client           persistence.StoragePersistenceClient
	rateLimitPerIP   int
	rateLimitByToken int
	ttl              int
	allowedTokens    []string
}

func NewLimiter(client persistence.StoragePersistenceClient, rateLimitPerIP int, rateLimitByToken int, ttl int, allowedTokens []string) *Limiter {
	return &Limiter{
		client:           client,
		rateLimitPerIP:   rateLimitPerIP,
		rateLimitByToken: rateLimitByToken,
		ttl:              ttl,
		allowedTokens:    allowedTokens,
	}
}

// RateLimited checks if the request is rate limited based on the IP address and API key.
// If the request is rate limited, it returns an error. If the request is not rate limited, it increments the counter and returns nil.
func (l *Limiter) RateLimited(ctx context.Context, ip string, token string) error {
	key, maxRequests := l.determineRateLimitKeyAndValue(token, ip)

	if isBlocked, _ := l.client.IsKeyBlocked(ctx, key+":blocked"); isBlocked {
		return ErrNotAllowed
	}

	counter, err := l.client.Increment(ctx, key, l.ttl)
	if err != nil {
		return err
	}

	if counter > maxRequests {
		blockKey := key + ":blocked"
		err = l.client.Set(ctx, blockKey, "1", l.ttl)
		if err != nil {
			return err
		}
		return ErrNotAllowed
	}

	return nil
}

// determineRateLimitKeyAndValue decides the appropriate key and rate limit based on the API key and IP address.
func (l *Limiter) determineRateLimitKeyAndValue(apiKey, ip string) (string, int) {
	var key string
	var limit int

	// If an API key is provided, and it is in the list of allowed tokens, use the token as the key and rate limit by token.
	if apiKey != "" && contains(l.allowedTokens, apiKey) {
		key = apiKey
		limit = l.rateLimitByToken
	} else {
		// If no API key is provided or the API key is not in the list of allowed tokens, use the IP address as the key and rate limit per IP.
		key = ip
		limit = l.rateLimitPerIP
	}

	return key, limit
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
