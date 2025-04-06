package websocket

import (
	"context"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	pkgerrors "github.com/pkg/errors"
)

func readRequest[TRequest Request](ctx context.Context, conn Conn) (TRequest, error) {
	var req TRequest
	if conn == nil {
		return req, pkgerrors.New("nil connection")
	}

	wsConn, ok := conn.(*websocket.Conn)
	if !ok {
		return req, pkgerrors.New("invalid connection type")
	}

	if err := wsjson.Read(ctx, wsConn, &req); err != nil {
		return req, pkgerrors.Wrap(err, "read request failed")
	}
	return req, nil
}

func writeResponse[TResponse Response](ctx context.Context, conn Conn, resp TResponse) error {
	if conn == nil {
		return pkgerrors.New("nil connection")
	}

	wsConn, ok := conn.(*websocket.Conn)
	if !ok {
		return pkgerrors.New("invalid connection type")
	}

	if err := wsjson.Write(ctx, wsConn, resp); err != nil {
		return pkgerrors.Wrap(err, "write response failed")
	}
	return nil
}
