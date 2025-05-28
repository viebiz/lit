package websocket

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"github.com/coder/websocket"
	pkgerrors "github.com/pkg/errors"
)

var (
	timeNowWrapper = time.Now
)

type clientConn[TRequest Request, TResponse Response] struct {
	id             string
	underlyingConn Conn
	ctx            context.Context
	cancel         func()

	requestCh  chan TRequest
	responseCh chan TResponse
	errCh      chan error

	info ClientConnInfo
}

type ClientConfig struct {
	ID        string
	ReadLimit int64
}

func initClientConn[TRequest Request, TResponse Response](
	ctx context.Context,
	cfg ClientConfig,
	cancelFunc func(),
	conn *websocket.Conn,
) ClientConn[TRequest, TResponse] {
	now := timeNowWrapper()

	if cfg.ReadLimit > 0 {
		conn.SetReadLimit(cfg.ReadLimit)
	}

	clc := &clientConn[TRequest, TResponse]{
		id:             cfg.ID,
		underlyingConn: websocketConn{Conn: conn},
		ctx:            ctx,
		cancel:         cancelFunc,
		requestCh:      make(chan TRequest, defaultRequestBufferPerClient),
		responseCh:     make(chan TResponse, defaultResponseBufferPerClient),
		errCh:          make(chan error, defaultErrorResponseBufferPerClient),
		info: ClientConnInfo{
			ID:          cfg.ID,
			Subprotocol: conn.Subprotocol(),
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	return clc
}

func (c *clientConn[TRequest, TResponse]) ID() string {
	return c.id
}

func (c *clientConn[TRequest, TResponse]) Context() context.Context {
	return c.ctx
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
	if err := c.underlyingConn.Ping(ctx); err != nil {
		if errors.Is(err, net.ErrClosed) {
			return ErrConnectionClosed
		}

		return pkgerrors.Wrap(err, "ping websocket connection")
	}

	c.info.UpdatedAt = timeNowWrapper()
	return nil
}

func (c *clientConn[TRequest, TResponse]) Disconnect() error {
	// TODO: Graceful disconnection
	close(c.requestCh)
	close(c.responseCh)
	close(c.errCh)

	if err := c.underlyingConn.Close(websocket.StatusNormalClosure, ""); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			return pkgerrors.Wrap(err, "close websocket connection")
		}
		return nil
	}

	return nil
}

func (c *clientConn[TRequest, TResponse]) Discard() {
	c.cancel()
}

func (c *clientConn[TRequest, TResponse]) Info() ClientConnInfo {
	return c.info
}

type websocketConn struct {
	*websocket.Conn
}

func (conn websocketConn) Writer(ctx context.Context, typ MessageType) (io.WriteCloser, error) {
	return conn.Conn.Writer(ctx, websocket.MessageType(typ))
}

func (conn websocketConn) Write(ctx context.Context, typ MessageType, p []byte) error {
	return conn.Conn.Write(ctx, websocket.MessageType(typ), p)
}

func (conn websocketConn) Reader(ctx context.Context) (MessageType, io.Reader, error) {
	msgType, reader, err := conn.Conn.Reader(ctx)
	return MessageType(msgType), reader, err
}

func (conn websocketConn) Read(ctx context.Context) (MessageType, []byte, error) {
	msgType, reader, err := conn.Conn.Read(ctx)
	return MessageType(msgType), reader, err
}
