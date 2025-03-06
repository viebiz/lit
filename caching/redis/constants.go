package redis

import (
	"github.com/viebiz/redis"
)

const (
	KeepTTL = redis.KeepTTL // Used to retain the existing TTL (time to live) of a key when modifying its value in Redis

	statusOK = "OK"
)
