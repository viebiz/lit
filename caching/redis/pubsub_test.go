package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/mocks/mockredis"
)

func TestRedisClient_Publish(t *testing.T) {
	type mockCmdArg struct {
		givenContext context.Context
		givenChannel string
		givenValue   interface{}
		expCmd       *redis.IntCmd
	}

	type args struct {
		givenMockCmdArgFn func() mockCmdArg
		givenMockCmdArg   mockCmdArg
		givenContext      context.Context
		givenChannel      string
		givenValue        any
		expErr            error
	}
	tcs := map[string]args{
		"success": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.IntCmd
				return mockCmdArg{
					givenContext: context.Background(),
					givenChannel: "key",
					givenValue:   "value",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenChannel: "key",
			givenValue:   "value",
		},
		"error: unexpected error": {
			givenMockCmdArgFn: func() mockCmdArg {
				var cmd redis.IntCmd
				cmd.SetErr(errors.New("simulated error"))
				return mockCmdArg{
					givenContext: context.Background(),
					givenChannel: "key",
					givenValue:   "value",
					expCmd:       &cmd,
				}
			},
			givenContext: context.Background(),
			givenChannel: "key",
			givenValue:   "value",
			expErr:       errors.New("simulated error"),
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// Mock
			tc.givenMockCmdArg = tc.givenMockCmdArgFn()
			mockUniversalClient.ExpectedCalls = []*mock.Call{
				mockUniversalClient.On(
					"Publish",
					tc.givenMockCmdArg.givenContext,
					tc.givenMockCmdArg.givenChannel,
					tc.givenMockCmdArg.givenValue,
				).Return(tc.givenMockCmdArg.expCmd),
			}

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			err := instance.Publish(tc.givenContext, tc.givenChannel, tc.givenValue)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRedisClient_Subscribe(t *testing.T) {
	type args struct {
		givenContext  context.Context
		givenChannels []string
		givenHandler  MessageHandler
	}
	tcs := map[string]args{
		"success": {
			givenContext:  context.Background(),
			givenChannels: []string{"key"},
			givenHandler: func(ctx context.Context, msg Message) error {
				return nil
			},
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			sub := instance.Subscribe(tc.givenContext, tc.givenChannels, tc.givenHandler)

			// Then
			require.NotNil(t, sub)
		})
	}
}

func TestSubscriber_handleMessage(t *testing.T) {
	type args struct {
		givenMsg      *redis.Message
		handlerErr    error
		panicInHandle bool
	}
	tcs := map[string]args{
		"success: handler ok": {
			givenMsg: &redis.Message{Channel: "a", Payload: "ok"},
		},
		"error: handler returns error": {
			givenMsg:   &redis.Message{Channel: "b", Payload: "fail"},
			handlerErr: errors.New("failed"),
		},
		"panic: handler panics": {
			givenMsg:      &redis.Message{Channel: "c", Payload: "panic"},
			panicInHandle: true,
		},
	}
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			var called bool
			handler := func(ctx context.Context, m Message) error {
				called = true
				if tc.panicInHandle {
					panic("fake panic")
				}
				return tc.handlerErr
			}
			s := &subscriber{
				handler: handler,
			}
			s.handleMessage(tc.givenMsg)
			require.True(t, called || tc.panicInHandle)
		})
	}
}

func TestSubscriber_Subscribe(t *testing.T) {
	type args struct {
		givenContext  context.Context
		givenChannels []string
		givenHandler  MessageHandler
	}
	tcs := map[string]args{
		"success": {
			givenContext:  context.Background(),
			givenChannels: []string{"test-channel"},
			givenHandler: func(ctx context.Context, msg Message) error {
				return nil
			},
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockUniversalClient := new(mockredis.MockUniversalClient)

			// When
			instance := redisClient{
				rdb: mockUniversalClient,
			}
			sub := instance.Subscribe(tc.givenContext, tc.givenChannels, tc.givenHandler)

			// Then
			require.NotNil(t, sub)
			require.IsType(t, &subscriber{}, sub)
		})
	}
}

func TestSubscriber_SubscribeWithOptions(t *testing.T) {
	type args struct {
		givenContext  context.Context
		givenChannels []string
		givenHandler  MessageHandler
		givenOpts     ChannelOption
		receiveErr    error
		channelClosed bool
		contextDone   bool
	}
	tcs := map[string]args{
		"success: receives message and handles": {
			givenContext: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				go func() {
					time.Sleep(10 * time.Millisecond) // Wait for message to be processed
					cancel()
				}()
				return ctx
			}(),
			givenChannels: []string{"test-channel"},
			givenHandler: func(ctx context.Context, msg Message) error {
				return nil
			},
			givenOpts: ChannelOption{},
		},
		"error: receive fails": {
			givenContext:  context.Background(),
			givenChannels: []string{"test-channel"},
			givenHandler: func(ctx context.Context, msg Message) error {
				return nil
			},
			givenOpts:  ChannelOption{},
			receiveErr: errors.New("receive error"),
		},
		"channel closed": {
			givenContext:  context.Background(),
			givenChannels: []string{"test-channel"},
			givenHandler: func(ctx context.Context, msg Message) error {
				return nil
			},
			givenOpts:     ChannelOption{},
			channelClosed: true,
		},
		"context done": {
			givenContext: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			givenChannels: []string{"test-channel"},
			givenHandler: func(ctx context.Context, msg Message) error {
				return nil
			},
			givenOpts:   ChannelOption{},
			contextDone: true,
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			mockPubSub := new(mockredis.MockPubSub)
			mockPubSub.On("Receive", mock.Anything).Return(nil, tc.receiveErr)
			if tc.receiveErr == nil {
				msgCh := make(chan *redis.Message, 1)
				if tc.channelClosed {
					close(msgCh)
				} else {
					msgCh <- &redis.Message{Channel: "test-channel", Payload: "test"}
				}
				mockPubSub.On("Channel", mock.Anything).Return((<-chan *redis.Message)(msgCh))
			}
			mockPubSub.On("Close").Return(nil)

			s := &subscriber{
				channels: tc.givenChannels,
				handler:  tc.givenHandler,
				subscribeFunc: func(ctx context.Context, channels ...string) RedisPubSub {
					return mockPubSub
				},
			}

			// When
			err := s.SubscribeWithOptions(tc.givenContext, tc.givenOpts)

			// Then
			if tc.receiveErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), "initial subscribe/receive failed")
			} else if tc.channelClosed {
				require.Error(t, err)
				require.Contains(t, err.Error(), "channel was closed by redis SDK")
			} else if tc.contextDone {
				require.NoError(t, err)
			} else {
				// For success case, we need to cancel context to exit the loop
				// But since it's a unit test, we'll assume it runs without error
				require.NoError(t, err)
			}
			mockPubSub.AssertExpectations(t)
		})
	}
}
