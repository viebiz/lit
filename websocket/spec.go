package websocket

import (
	"context"
	"io"

	"github.com/coder/websocket"
	"github.com/viebiz/lit"
)

// Request represents a request interface for websocket communication
type Request interface{}

// Response represents a response interface for websocket communication
type Response interface{}

// Hub represents a websocket Hub that manages client connections and handles messages
type Hub[TRequest Request, TResponse Response] interface {
	On(c lit.Router)

	Emit(ctx context.Context, id string, message TResponse, opt EmitOption) error

	Broadcast(ctx context.Context, message TResponse, opt EmitOption) error
}

// ClientConn represents a connection interface for a websocket client
type ClientConn[TRequest Request, TResponse Response] interface {
	ID() string
	Context() context.Context
	Conn() Conn
	Request() chan TRequest
	Response() chan TResponse
	Error() chan error
	Ping(ctx context.Context) error
	Discard()
	Disconnect() error
	Info() ClientConnInfo
}

// Conn represents a websocket connection interface
type Conn interface {
	Writer(ctx context.Context, typ MessageType) (io.WriteCloser, error)
	Write(ctx context.Context, typ MessageType, p []byte) error

	Close(code websocket.StatusCode, reason string) (err error)
	CloseNow() (err error)

	Reader(ctx context.Context) (MessageType, io.Reader, error)
	Read(ctx context.Context) (MessageType, []byte, error)
	CloseRead(ctx context.Context) context.Context

	SetReadLimit(n int64)
	Subprotocol() string

	Ping(ctx context.Context) error
}

// Handler represents a generic handler function type for handling websocket requests
type Handler[TRequest Request, TResponse Response] func(ctx context.Context, req TRequest) (TResponse, error)

// Serializer is an interface for serializing and deserializing messages
type Serializer interface {
	Serialize(o any, w io.Writer) error

	Deserialize(r io.Reader, o any) error
}
