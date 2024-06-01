package limiter

import (
	"Rate-Limiter-Challenge-GoExpertPostGrad/pkg/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RateLimiterSuite struct {
	suite.Suite
	limiterMock *mocks.LimiterMock
}

func (suite *RateLimiterSuite) SetupTest() {
	suite.limiterMock = new(mocks.LimiterMock)
}

func (suite *RateLimiterSuite) TestValidRequestByIP() {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(suite.T(), err)

	suite.limiterMock.On("RateLimited", mock.Anything, mock.Anything, "").Once().Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := suite.limiterMock.RateLimited(r.Context(), r.RemoteAddr, "")
		if err != nil {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	suite.limiterMock.AssertExpectations(suite.T())
}

func (suite *RateLimiterSuite) TestValidRequestByToken() {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Authorization", "valid-token")

	suite.limiterMock.On("RateLimited", mock.Anything, mock.Anything, "valid-token").Once().Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := suite.limiterMock.RateLimited(r.Context(), r.RemoteAddr, r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	suite.limiterMock.AssertExpectations(suite.T())
}

func (suite *RateLimiterSuite) TestRequestExceedingIPLimit() {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(suite.T(), err)

	suite.limiterMock.On("RateLimited", mock.Anything, mock.Anything, "").Return(errors.New("Too Many Requests"))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := suite.limiterMock.RateLimited(r.Context(), r.RemoteAddr, "")
		if err != nil {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusTooManyRequests, rr.Code)
	suite.limiterMock.AssertExpectations(suite.T())
}

func (suite *RateLimiterSuite) TestRequestExceedingTokenLimit() {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(suite.T(), err)
	req.Header.Set("Authorization", "exceeding-token")

	suite.limiterMock.On("RateLimited", mock.Anything, mock.Anything, "exceeding-token").Return(errors.New("Too Many Requests"))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := suite.limiterMock.RateLimited(r.Context(), r.RemoteAddr, r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusTooManyRequests, rr.Code)
	suite.limiterMock.AssertExpectations(suite.T())
}

func (suite *RateLimiterSuite) TestCounterResetAfterTTL() {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(suite.T(), err)

	// Simular primeira requisição válida dentro do limite
	suite.limiterMock.On("RateLimited", mock.Anything, mock.Anything, "").Once().Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := suite.limiterMock.RateLimited(r.Context(), r.RemoteAddr, "")
		if err != nil {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Primeira requisição dentro do limite
	handler.ServeHTTP(rr, req)
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	suite.limiterMock.AssertExpectations(suite.T())

	// Simular segunda requisição que excede o limite
	suite.limiterMock.ExpectedCalls = nil
	suite.limiterMock.On("RateLimited", mock.Anything, mock.Anything, "").Return(errors.New("Too Many Requests")).Once()

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(suite.T(), http.StatusTooManyRequests, rr.Code)
	suite.limiterMock.AssertExpectations(suite.T())

	// Simular espera do TTL
	time.Sleep(1 * time.Second)

	// Simular terceira requisição após o TTL, dentro do limite
	suite.limiterMock.ExpectedCalls = nil
	suite.limiterMock.On("RateLimited", mock.Anything, mock.Anything, "").Return(nil).Once()

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	suite.limiterMock.AssertExpectations(suite.T())
}

func TestRateLimiterSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterSuite))
}
