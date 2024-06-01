package middleware

import (
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/limiter"
	"Rate-Limiter-Challenge-GoExpertPostGrad/pkg/ratelimiter"
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
)

type RateLimiterMiddleware struct {
	limiter ratelimiter.LimiterInterface
}

func NewRateLimiterMiddleware(l ratelimiter.LimiterInterface) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{limiter: l}
}

func (rl *RateLimiterMiddleware) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		token := strings.TrimSpace(r.Header.Get("API_KEY"))
		ip := getClientIP(r)

		err := rl.limiter.RateLimited(ctx, ip, token)

		if errors.Is(err, limiter.ErrNotAllowed) {
			http.Error(w, "You have reached the maximum number of requests or actions allowed within a certain time frame.", http.StatusTooManyRequests)
			return
		}

		if err != nil {
			log.Printf("Limiter Error: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // Use the original value if it can't be split
	}
	return host
}
