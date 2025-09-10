package redis

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/viebiz/lit/monitoring"
)

// Publish posts a new message to server
func (client redisClient) Publish(ctx context.Context, channel string, payload any) error {
	// Use SPUBLISH for better efficiency with sharding (Redis Cluster), If not using sharding it the same PUBLISH
	return client.rdb.SPublish(ctx, channel, payload).Err()
}

func (client redisClient) Subscribe(ctx context.Context, channels []string, handler MessageHandler) Subscriber {
	monitor := monitoring.FromContext(ctx)
	monitor.Infof("Redis subscriber initializing: [%s]", channels)

	monitor = monitor.With(map[string]string{
		"redis.subscriber.channels": strings.Join(channels, ","),
	})

	return &subscriber{
		rdb:      client.rdb,
		channels: channels,
		handler:  handler,
		monitor:  monitor,
		subscribeFunc: func(ctx context.Context, channels ...string) RedisPubSub {
			return client.rdb.SSubscribe(ctx, channels...)
		},
	}
}

type Message struct {
	Channel      string
	Pattern      string
	Payload      string
	PayloadSlice []string
}

func (msg *Message) from(m *redis.Message) {
	msg.Channel = m.Channel
	msg.Pattern = m.Pattern
	msg.Payload = m.Payload
	msg.PayloadSlice = m.PayloadSlice
}

type MessageHandler func(ctx context.Context, msg Message) error

type Subscriber interface {
	Subscribe(ctx context.Context) error

	SubscribeWithOptions(ctx context.Context, opts ChannelOption) error
}

type RedisPubSub interface {
	Receive(ctx context.Context) (interface{}, error)
	Channel(opts ...redis.ChannelOption) <-chan *redis.Message
	Close() error
}

type subscriber struct {
	rdb           redis.UniversalClient
	monitor       *monitoring.Monitor
	channels      []string
	handler       MessageHandler
	subscribeFunc func(ctx context.Context, channels ...string) RedisPubSub
}

func (s *subscriber) Subscribe(ctx context.Context) error {
	return s.SubscribeWithOptions(ctx, ChannelOption{
		HealthCheckInterval: time.Minute,
	})
}

func (s *subscriber) SubscribeWithOptions(ctx context.Context, opts ChannelOption) error {
	ps := s.subscribeFunc(ctx, s.channels...)
	defer func() {
		if err := ps.Close(); err != nil {
			s.monitor.Errorf(err, "[redis_subscriber] Closing subscriber channels: %v", s.channels)
		}
		s.monitor.Infof("[redis_subscriber] closed")
	}()

	if _, err := ps.Receive(ctx); err != nil {
		return pkgerrors.Wrap(err, "initial subscribe/receive failed")
	}

	s.monitor.Infof("[redis_subscriber] subscribed; channels=%v", s.channels)

	chOpts := opts.toRedisOptions()
	msgCh := ps.Channel(chOpts...)
	for {
		select {
		case <-ctx.Done():
			s.monitor.Infof("[redis_subscriber] context done → closing: %v", s.channels)
			return nil
		case msg, ok := <-msgCh:
			if !ok {
				s.monitor.Info("[redis_subscriber] message channel was closed by redis SDK")
				return pkgerrors.New("channel was closed by redis SDK")
			}

			s.handleMessage(msg)
		}
	}
}

func (s *subscriber) handleMessage(msg *redis.Message) {
	ctx := context.Background()

	var err error
	defer func() {
		if rcv := recover(); rcv != nil {
			err = fmt.Errorf("panic: %v", rcv)
			monitoring.FromContext(ctx).Errorf(err, "Caught PANIC. Stack trace: %s", debug.Stack())
		}
	}()

	s.monitor.Infof("[redis_subscriber] Received message: %s", msg.Payload)
	var m Message
	m.from(msg)

	// TODO: Implement retry backoff
	if err := s.handler(ctx, m); err != nil {
		s.monitor.Errorf(err, "[redis_subscriber] Handle message failed")
		return
	}
}
