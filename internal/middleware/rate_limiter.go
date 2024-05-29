package middleware

import (
	"Rate-Limiter-Challenge-GoExpertPostGrad/pkg/ratelimiter"
	"log"
	"net"
	"net/http"
	"strings"
)

type RateLimiterMiddleware struct {
	limiter ratelimiter.RateLimiter
}

func NewRateLimiterMiddleware(limiter ratelimiter.RateLimiter) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{limiter: limiter}
}

func (rlm *RateLimiterMiddleware) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = "127.0.0.1"
		}

		token := strings.TrimSpace(r.Header.Get("API_KEY"))

		var key string
		var limit int

		if token != "" {
			key = "token:" + token
			limit = rlm.limiter.GetRateLimit(token)
		} else {
			key = "ip:" + ip
			limit = rlm.limiter.GetIP()
		}

		if rlm.limiter.IsRateLimited(ctx, key, limit) {
			log.Printf("Request blocked for key %s with limit %d", key, limit)
			http.Error(w, `{"error": "you have reached the maximum number of requests or actions allowed within a certain time frame"}`, http.StatusTooManyRequests)
			return
		}
		log.Printf("Request allowed for key %s", key)

		rlm.limiter.IncrementRequestCount(ctx, key)
		next.ServeHTTP(w, r)
	})
}
