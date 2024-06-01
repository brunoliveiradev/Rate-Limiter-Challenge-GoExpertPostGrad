package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type LimiterMock struct {
	mock.Mock
}

func NewLimiterMock() *LimiterMock {
	return &LimiterMock{}
}

func (lm *LimiterMock) RateLimited(ctx context.Context, ip string, token string) error {
	args := lm.Called(ctx, ip, token)
	return args.Error(0)
}
