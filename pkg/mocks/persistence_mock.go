package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type StoragePersistenceClientMock struct {
	mock.Mock
}

func NewStorageMock() *StoragePersistenceClientMock {
	return &StoragePersistenceClientMock{}
}

func (m *StoragePersistenceClientMock) Set(ctx context.Context, key string, value string, ttl int) error {
	args := m.Called(ctx, key, value, time.Duration(ttl)*time.Second)
	return args.Error(0)
}

func (m *StoragePersistenceClientMock) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), nil
}

func (m *StoragePersistenceClientMock) Increment(ctx context.Context, key string, ttl int) (int, error) {
	args := m.Called(ctx, key, time.Duration(ttl)*time.Second)
	return args.Int(0), args.Error(1)
}

func (m *StoragePersistenceClientMock) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *StoragePersistenceClientMock) IsKeyBlocked(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *StoragePersistenceClientMock) FlushAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
