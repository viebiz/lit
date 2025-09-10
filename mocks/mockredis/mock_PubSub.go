package mockredis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type MockPubSub struct {
	mock.Mock
}

func (m *MockPubSub) Receive(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0), args.Error(1)
}

func (m *MockPubSub) Channel(opts ...redis.ChannelOption) <-chan *redis.Message {
	args := m.Called(opts)
	return args.Get(0).(<-chan *redis.Message)
}

func (m *MockPubSub) Close() error {
	args := m.Called()
	return args.Error(0)
}
