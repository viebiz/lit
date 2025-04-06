package websocket

import (
	"time"
)

const (
	// Default buffer sizes for channels in the client connection
	defaultRequestBufferPerClient       = 10
	defaultResponseBufferPerClient      = 10
	defaultErrorResponseBufferPerClient = 10

	// Default ping interval for the client connection
	// Refer config from https://socket.io/docs/v4/server-options/#pingtimeout
	defaultPingTimeout  = 20 * time.Second
	defaultPingInterval = 25 * time.Second

	//
	defaultClientIDHeader = "X-Client-ID"
)
