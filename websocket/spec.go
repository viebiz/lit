package websocket

import (
	"context"
	"io"

	"github.com/coder/websocket"
	"github.com/viebiz/lit"
)

// Request represents a request interface for websocket communication
type Request interface {
	//json.Unmarshaler
	//json.Marshaler
}

// Response represents a response interface for websocket communication
type Response interface {
	//json.Unmarshaler
	//json.Marshaler
}

// Hub represents a websocket Hub that manages client connections and handles messages
type Hub[TRequest Request, TResponse Response] interface {
	On(c lit.Router)

	Emit(ctx context.Context, id ClientIdentifier, message TResponse) error
}

// ClientConn represents a connection interface for a websocket client
type ClientConn[TRequest Request, TResponse Response] interface {
	Conn() Conn
	Request() chan TRequest
	Response() chan TResponse
	Error() chan error
	Ping(ctx context.Context) error
	Discard()
	Disconnect() error
}

// Conn represents a websocket connection interface
type Conn interface {
	Writer(ctx context.Context, typ websocket.MessageType) (io.WriteCloser, error)
	Write(ctx context.Context, typ websocket.MessageType, p []byte) error

	Close(code websocket.StatusCode, reason string) (err error)
	CloseNow() (err error)

	Reader(ctx context.Context) (websocket.MessageType, io.Reader, error)
	Read(ctx context.Context) (websocket.MessageType, []byte, error)
	CloseRead(ctx context.Context) context.Context

	SetReadLimit(n int64)
	Subprotocol() string

	Ping(ctx context.Context) error
}

// Handler represents a generic handler function type for handling websocket requests
type Handler[TRequest Request, TResponse Response] func(ctx context.Context, req TRequest) (TResponse, error)

type ClientIdentifier interface {
	ID() string
}
