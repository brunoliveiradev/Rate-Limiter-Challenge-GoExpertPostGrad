package main

import (
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/configs"
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/limiter"
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/middleware"
	"Rate-Limiter-Challenge-GoExpertPostGrad/pkg/persistence"
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var redisClient *persistence.RedisClient

func TestMain(m *testing.M) {
	// Set up the Redis client
	redisClient = persistence.NewRedisClient("localhost:6379", "", 0)
	// Run tests
	code := m.Run()
	os.Exit(code)
}

func cleanRedis() {
	redisClient.FlushAll(context.Background())
}

func loadConfig() *configs.Config {
	return &configs.Config{
		ServerPort:        8080,
		RedisAddr:         "localhost:6379",
		RedisPassword:     "",
		RedisDB:           0,
		RateLimitByIP:     2,
		RateLimiteByToken: 10,
		Tokens:            []string{"token1", "token2", "token3"},
		Ttl:               5,
	}
}

func TestValidTokenWithinLimit(t *testing.T) {
	cleanRedis()
	conf := loadConfig()
	rateLimiter := limiter.NewLimiter(redisClient, conf)
	mw := middleware.NewRateLimiterMiddleware(rateLimiter)

	for i := 0; i < 10; i++ {
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		req.Header.Set("API_KEY", "token1")
		req.RemoteAddr = "192.168.0.1:12345"

		rr := httptest.NewRecorder()
		handler := mw.Limit(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	}
}

func TestValidTokenExceedingLimit(t *testing.T) {
	cleanRedis()
	conf := loadConfig()
	rateLimiter := limiter.NewLimiter(redisClient, conf)
	mw := middleware.NewRateLimiterMiddleware(rateLimiter)

	for i := 0; i < 11; i++ {
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		req.Header.Set("API_KEY", "token1")
		req.RemoteAddr = "192.168.0.1:12345"

		rr := httptest.NewRecorder()
		handler := mw.Limit(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler.ServeHTTP(rr, req)

		if i < 10 {
			assert.Equal(t, http.StatusOK, rr.Code)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		}
	}
}

func TestInvalidTokenWithinLimit(t *testing.T) {
	cleanRedis()
	conf := loadConfig()
	rateLimiter := limiter.NewLimiter(redisClient, conf)
	mw := middleware.NewRateLimiterMiddleware(rateLimiter)

	for i := 0; i < 2; i++ {
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		req.Header.Set("API_KEY", "invalid-token")
		req.RemoteAddr = "192.168.0.2:12345"

		rr := httptest.NewRecorder()
		handler := mw.Limit(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	}
}

func TestInvalidTokenExceedingLimit(t *testing.T) {
	cleanRedis()
	conf := loadConfig()
	rateLimiter := limiter.NewLimiter(redisClient, conf)
	mw := middleware.NewRateLimiterMiddleware(rateLimiter)

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		req.Header.Set("API_KEY", "invalid-token")
		req.RemoteAddr = "192.168.0.2:12345"

		rr := httptest.NewRecorder()
		handler := mw.Limit(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler.ServeHTTP(rr, req)

		if i < 2 {
			assert.Equal(t, http.StatusOK, rr.Code)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		}
	}
}

func TestNoTokenWithinLimit(t *testing.T) {
	cleanRedis()
	conf := loadConfig()
	rateLimiter := limiter.NewLimiter(redisClient, conf)
	mw := middleware.NewRateLimiterMiddleware(rateLimiter)

	for i := 0; i < 2; i++ {
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		req.RemoteAddr = "192.168.0.3:12345"

		rr := httptest.NewRecorder()
		handler := mw.Limit(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	}
}

func TestNoTokenExceedingLimit(t *testing.T) {
	cleanRedis()
	conf := loadConfig()
	rateLimiter := limiter.NewLimiter(redisClient, conf)
	mw := middleware.NewRateLimiterMiddleware(rateLimiter)

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)
		req.RemoteAddr = "192.168.0.3:12345"

		rr := httptest.NewRecorder()
		handler := mw.Limit(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler.ServeHTTP(rr, req)

		if i < 2 {
			assert.Equal(t, http.StatusOK, rr.Code)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		}
	}
}
