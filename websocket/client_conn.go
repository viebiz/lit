package websocket

import (
	"context"

	"github.com/coder/websocket"
	pkgerrors "github.com/pkg/errors"
)

type clientConn[TRequest Request, TResponse Response] struct {
	requestCh      chan TRequest
	responseCh     chan TResponse
	errCh          chan error
	underlyingConn *websocket.Conn
	cancel         func()
}

func initClientConn[TRequest Request, TResponse Response](
	ctx context.Context,
	conn *websocket.Conn,
	cancelFunc func(),
) ClientConn[TRequest, TResponse] {
	return &clientConn[TRequest, TResponse]{
		requestCh:      make(chan TRequest, defaultRequestBufferPerClient),
		responseCh:     make(chan TResponse, defaultResponseBufferPerClient),
		errCh:          make(chan error, defaultErrorResponseBufferPerClient),
		underlyingConn: conn,
		cancel:         cancelFunc,
	}
}

func (c *clientConn[TRequest, TResponse]) Conn() Conn {
	return c.underlyingConn
}

func (c *clientConn[TRequest, TResponse]) Request() chan TRequest {
	return c.requestCh
}

func (c *clientConn[TRequest, TResponse]) Response() chan TResponse {
	return c.responseCh
}

func (c *clientConn[TRequest, TResponse]) Error() chan error {
	return c.errCh
}

func (c *clientConn[TRequest, TResponse]) Ping(ctx context.Context) error {
	if c.underlyingConn == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, defaultPingTimeout)
	defer cancel()

	if err := c.underlyingConn.Ping(ctx); err != nil {
		return pkgerrors.Wrap(err, "ping websocket connection")
	}

	return nil
}

func (c *clientConn[TRequest, TResponse]) Disconnect() error {
	if c.underlyingConn != nil {
		if err := c.underlyingConn.Close(websocket.StatusNormalClosure, ""); err != nil {
			return pkgerrors.Wrap(err, "close websocket connection")
		}
	}

	close(c.requestCh)
	close(c.responseCh)
	close(c.errCh)
	return nil
}

func (c *clientConn[TRequest, TResponse]) Discard() {
	c.cancel()
}
