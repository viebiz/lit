package websocket

import "github.com/coder/websocket"

type event struct {
	Type    websocket.MessageType
	Content []byte
}
