package limiter

import (
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/configs"
	"Rate-Limiter-Challenge-GoExpertPostGrad/pkg/persistence"
	"context"
	"errors"
)

var ErrNotAllowed = errors.New("access not allowed")

type Limiter struct {
	config        *configs.Envs
	redis         *persistence.RedisClient
	allowedTokens map[string]bool
}

func NewLimiter(redis *persistence.RedisClient, config *configs.Envs) *Limiter {
	// This map will allow for faster lookups when checking if a token is allowed instead of iterating over the slice of tokens.
	allowedTokens := make(map[string]bool)
	for _, token := range config.AllowedTokens {
		allowedTokens[token] = true
	}

	return &Limiter{
		config:        config,
		redis:         redis,
		allowedTokens: allowedTokens,
	}
}

// RateLimited checks if the request is rate limited based on the IP address and API key.
// If the request is rate limited, it returns an error. If the request is not rate limited, it increments the counter and returns nil.
func (rl *Limiter) RateLimited(ctx context.Context, ip string, token string) error {
	key, maxRequests := rl.determineRateLimitKeyAndValue(token, ip)

	if isBlocked, _ := rl.redis.IsKeyBlocked(ctx, key+":blocked"); isBlocked {
		return ErrNotAllowed
	}

	counter, err := rl.redis.Increment(ctx, key, rl.config.TTLSeconds)
	if err != nil {
		return err
	}

	if counter > maxRequests {
		blockKey := key + ":blocked"
		err = rl.redis.Set(ctx, blockKey, "1", rl.config.TTLSeconds)
		if err != nil {
			return err
		}
		return ErrNotAllowed
	}

	return nil
}

func (rl *Limiter) determineRateLimitKeyAndValue(token string, ip string) (string, int) {
	if token != "" && rl.isAllowedToken(token) {
		return token, rl.config.RateLimitByToken
	}
	return ip, rl.config.RateLimitByIP
}

// isAllowedToken checks if the token is allowed using a map of allowed tokens.
func (rl *Limiter) isAllowedToken(token string) bool {
	_, ok := rl.allowedTokens[token]
	return ok
}
