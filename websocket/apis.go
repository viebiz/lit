package websocket

import (
	"context"
)

func NewPublicHub[TRequest Request, TResponse Response](
	ctx context.Context,
	baseURL string,
	handler Handler[TRequest, TResponse],
) Hub[TRequest, TResponse] {
	h := &hub[TRequest, TResponse]{
		clients: make(map[string]ClientConn[TRequest, TResponse]),
		baseURL: baseURL,
		handler: handler,
	}

	return h
}
