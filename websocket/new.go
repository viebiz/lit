package websocket

import (
	"context"
	"sync"
	"time"
)

const (
	MaxConnections = 1000
	PingInterval   = 30 * time.Second
	PongWait       = 35 * time.Second
)

type impl struct {
	clients   sync.Map
	broadcast chan event
	mutex     sync.Mutex
	count     int
	ctx       context.Context
	cancel    context.CancelFunc
}

func New() Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &impl{
		broadcast: make(chan event),
		ctx:       ctx,
		cancel:    cancel,
	}
}
