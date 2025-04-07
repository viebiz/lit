package websocket

import (
	"time"
)

const (
	// Default conn configs
	// Refer https://socket.io/docs/v4/server-options/#pingtimeout
	defaultPingTimeout  = 20 * time.Second
	defaultPingInterval = 25 * time.Second
	defaultReadLimit    = 32 << 10 // 32KB

	defaultRequestBufferPerClient       = 100
	defaultResponseBufferPerClient      = 100
	defaultErrorResponseBufferPerClient = 100

	defaultClientIDHeaderName = "X-Client-ID"
	defaultClientIDParamName  = "client_id"
)
