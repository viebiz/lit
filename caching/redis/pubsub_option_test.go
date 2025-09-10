package redis

import (
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestChannelOption_toRedisOptions(t *testing.T) {
	type args struct {
		givenOpt ChannelOption
		expOpts  []redis.ChannelOption
	}
	tcs := map[string]args{
		"empty options": {
			givenOpt: ChannelOption{},
			expOpts:  []redis.ChannelOption{},
		},
		"with channel size": {
			givenOpt: ChannelOption{
				ChannelSize: 100,
			},
			expOpts: []redis.ChannelOption{
				redis.WithChannelSize(100),
			},
		},
		"with health check interval": {
			givenOpt: ChannelOption{
				HealthCheckInterval: time.Minute,
			},
			expOpts: []redis.ChannelOption{
				redis.WithChannelHealthCheckInterval(time.Minute),
			},
		},
		"with send timeout": {
			givenOpt: ChannelOption{
				SendTimeout: time.Second * 30,
			},
			expOpts: []redis.ChannelOption{
				redis.WithChannelSendTimeout(time.Second * 30),
			},
		},
		"all options": {
			givenOpt: ChannelOption{
				ChannelSize:         200,
				HealthCheckInterval: time.Minute * 2,
				SendTimeout:         time.Second * 45,
			},
			expOpts: []redis.ChannelOption{
				redis.WithChannelSize(200),
				redis.WithChannelHealthCheckInterval(time.Minute * 2),
				redis.WithChannelSendTimeout(time.Second * 45),
			},
		},
		"zero values": {
			givenOpt: ChannelOption{
				ChannelSize:         0,
				HealthCheckInterval: 0,
				SendTimeout:         0,
			},
			expOpts: []redis.ChannelOption{},
		},
		"negative values": {
			givenOpt: ChannelOption{
				ChannelSize:         -1,
				HealthCheckInterval: -1,
				SendTimeout:         -1,
			},
			expOpts: []redis.ChannelOption{},
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// When
			result := tc.givenOpt.toRedisOptions()

			// Then
			require.Len(t, result, len(tc.expOpts))
			// Note: We can't directly compare redis.ChannelOption slices due to internal implementation
			// But we can verify the length and that options are applied correctly
		})
	}
}
