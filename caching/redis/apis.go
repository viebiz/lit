package redis

import (
	"crypto/tls"

	"github.com/redis/go-redis/v9"

	pkgerrors "github.com/pkg/errors"

	"github.com/viebiz/lit/monitoring"
)

// NewClient creates new redis client
//
// Examples:
//
//	client, err := NewClient("redis://user:password@localhost:6379/0?protocol=3")
func NewClient(url string) (Client, error) {
	return NewClientWithTLS(url, nil)
}

// NewClientWithTLS creates new redis client with TLS configuration
func NewClientWithTLS(url string, tlsConfig *tls.Config) (Client, error) {
	// 1. Prepare Redis client configurations
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	// 1.2. Add TLS config to Redis client options
	opts.TLSConfig = tlsConfig

	// 2. Create new Redis client
	rdb := redis.NewClient(opts)

	// 3. Setup Redis tracing hook
	info := monitoring.NewExternalServiceInfo(url)
	rdb.AddHook(newTracingHook(info))

	return redisClient{
		rdb: rdb,
	}, nil
}
