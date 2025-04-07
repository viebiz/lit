package websocket

import (
	"time"
)

type HubInfo struct {
	ClientCount int
}

type ClientConnInfo struct {
	ID           string
	Subprotocol  string
	RemoteAddr   string
	MessageCount int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
