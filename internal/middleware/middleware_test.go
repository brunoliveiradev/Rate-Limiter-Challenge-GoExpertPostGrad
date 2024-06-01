package middleware

import (
	"Rate-Limiter-Challenge-GoExpertPostGrad/internal/limiter"

	"Rate-Limiter-Challenge-GoExpertPostGrad/pkg/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateLimiterMiddleware_ValidTokenWithinLimit(t *testing.T) {
	limiterMock := new(mocks.LimiterMock)
	mw := NewRateLimiterMiddleware(limiterMock)

	limiterMock.On("RateLimited", mock.Anything, "192.168.0.1", "token1").Return(nil).Once()

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
	limiterMock.AssertExpectations(t)
}

func TestRateLimiterMiddleware_TokenExceedingLimit(t *testing.T) {
	limiterMock := new(mocks.LimiterMock)
	mw := NewRateLimiterMiddleware(limiterMock)

	limiterMock.On("RateLimited", mock.Anything, "192.168.0.1", "token1").Return(limiter.ErrNotAllowed).Once()

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	req.Header.Set("API_KEY", "token1")
	req.RemoteAddr = "192.168.0.1:12345"

	rr := httptest.NewRecorder()
	handler := mw.Limit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
	limiterMock.AssertExpectations(t)
}

func TestRateLimiterMiddleware_ValidIPWithinLimit(t *testing.T) {
	limiterMock := new(mocks.LimiterMock)
	mw := NewRateLimiterMiddleware(limiterMock)

	limiterMock.On("RateLimited", mock.Anything, "192.168.0.2", "").Return(nil).Once()

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	req.RemoteAddr = "192.168.0.2:12345"

	rr := httptest.NewRecorder()
	handler := mw.Limit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	limiterMock.AssertExpectations(t)
}

func TestRateLimiterMiddleware_IPExceedingLimit(t *testing.T) {
	limiterMock := new(mocks.LimiterMock)
	mw := NewRateLimiterMiddleware(limiterMock)

	limiterMock.On("RateLimited", mock.Anything, "192.168.0.2", "").Return(limiter.ErrNotAllowed).Once()

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	req.RemoteAddr = "192.168.0.2:12345"

	rr := httptest.NewRecorder()
	handler := mw.Limit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
	limiterMock.AssertExpectations(t)
}

func TestRateLimiterMiddleware_InvalidToken(t *testing.T) {
	limiterMock := new(mocks.LimiterMock)
	mw := NewRateLimiterMiddleware(limiterMock)

	limiterMock.On("RateLimited", mock.Anything, "192.168.0.3", "invalid-token").Return(nil).Once()

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	req.Header.Set("API_KEY", "invalid-token")
	req.RemoteAddr = "192.168.0.3:12345"

	rr := httptest.NewRecorder()
	handler := mw.Limit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	limiterMock.AssertExpectations(t)
}

func TestRateLimiterMiddleware_ErrorFromRateLimiter(t *testing.T) {
	limiterMock := new(mocks.LimiterMock)
	mw := NewRateLimiterMiddleware(limiterMock)

	limiterMock.On("RateLimited", mock.Anything, "192.168.0.4", "error-token").Return(errors.New("rate limiter error")).Once()

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	req.Header.Set("API_KEY", "error-token")
	req.RemoteAddr = "192.168.0.4:12345"

	rr := httptest.NewRecorder()
	handler := mw.Limit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	limiterMock.AssertExpectations(t)
}

func TestRateLimiterMiddleware_InvalidIP(t *testing.T) {
	limiterMock := new(mocks.LimiterMock)
	mw := NewRateLimiterMiddleware(limiterMock)

	limiterMock.On("RateLimited", mock.Anything, "invalid-ip", "").Return(nil).Once()

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	req.RemoteAddr = "invalid-ip:12345"

	rr := httptest.NewRecorder()
	handler := mw.Limit(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	limiterMock.AssertExpectations(t)
}
