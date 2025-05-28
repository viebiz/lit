package websocket

import (
	"github.com/coder/websocket"
)

// MessageType represents the type of a WebSocket message.
// See https://tools.ietf.org/html/rfc6455#section-5.6
type MessageType int

const (
	// MessageTypeText is for UTF-8 encoded text messages like JSON.
	MessageTypeText = MessageType(websocket.MessageText)
	// MessageTypeBinary is for binary messages like protobufs.
	MessageTypeBinary = MessageType(websocket.MessageBinary)
)
