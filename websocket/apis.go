package websocket

import (
	"context"
	"sync"
)

func NewHub[TRequest Request, TResponse Response](
	ctx context.Context,
	baseURL string,
	handler Handler[TRequest, TResponse],
	opts ...HubOption[TRequest, TResponse],
) Hub[TRequest, TResponse] {
	h := &hub[TRequest, TResponse]{
		baseURL:   baseURL,
		handler:   handler,
		clients:   make(map[string]ClientConn[TRequest, TResponse]),
		clientMux: sync.RWMutex{},
		settings:  getDefaultHubSettings(),
	}
	for _, opt := range opts {
		opt(h)
	}

	return h
}
