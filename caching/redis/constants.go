package redis

import (
	"github.com/redis/go-redis/v9"
)

const (
	KeepTTL = redis.KeepTTL // Used to retain the existing TTL (time to live) of a key when modifying its value in Redis

	statusOK = "OK"

	setModeNone setMode = ""

	setModeNX setMode = "NX"

	setModeXX setMode = "XX"
)

type setMode string

func (s setMode) String() string {
	return string(s)
}

func (s setMode) IsValid() bool {
	return s == setModeNX || s == setModeXX
}
