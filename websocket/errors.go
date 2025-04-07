package websocket

import (
	"errors"
)

var (
	ErrConnectionClosed = errors.New("connection closed")

	ErrClientNotFound = errors.New("client id not found")

	ErrClientAlreadyExists = errors.New("client already exists")
)
