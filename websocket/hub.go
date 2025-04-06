package websocket

import (
	"context"
	"sync"
	"time"

	"github.com/coder/websocket"
	pkgerrors "github.com/pkg/errors"
	"github.com/viebiz/lit"
	"github.com/viebiz/lit/monitoring"
)

type hub[TRequest Request, TResponse Response] struct {
	clients   map[string]ClientConn[TRequest, TResponse]
	clientMux sync.RWMutex

	baseURL string
	handler func(ctx context.Context, req TRequest) (TResponse, error)
}

func (h *hub[TRequest, TResponse]) On(router lit.Router) {
	router.Get(h.baseURL, h.on)
}

func (h *hub[TRequest, TResponse]) on(c lit.Context) error {
	// 1. Settings black list/ white list ...

	// 2. Get the client id from the request
	clientID := c.Request().Header.Get(defaultClientIDHeader)
	if clientID == "" {
		return pkgerrors.New("client id is empty")
	}

	// 3. Check if the client is already connected
	if exists := h.existsClient(clientID); exists {
		return pkgerrors.New("client already connected")
	}

	// 4. Upgrades the connection to a WebSocket
	conn, err := websocket.Accept(c.Writer(), c.Request(), &websocket.AcceptOptions{})
	if err != nil {
		return pkgerrors.Wrap(err, "accept websocket connection")
	}

	// 5. Create a new client connection
	ctx, cancel := context.WithCancel(monitoring.NewContext(c))
	defer cancel()

	clc := initClientConn[TRequest, TResponse](ctx, conn, cancel)
	h.setClient(clientID, clc)

	// 6. Ping the client
	go h.pingClient(ctx, clc)

	// 7. Handle incoming messages
	go h.handleRequest(ctx, clc)

	// 8. Process messages
	go h.processMessage(ctx, clc)

	// 9.1. Handle outgoing messages
	go h.handleResponse(ctx, clc)

	// 9.2. Handle errors
	go h.handleErrorResponse(ctx, clc)

	// 10. Handle connection close
	select {
	case <-ctx.Done():
		monitoring.FromContext(ctx).Infof("[websocket] attempt to disconnect client %s", clientID)

		h.removeClient(clientID)
		if err := clc.Disconnect(); err != nil {
			monitoring.FromContext(ctx).Errorf(err, "[websocket] close client connection failed")
		}

		monitoring.FromContext(ctx).Infof("[websocket] client %s disconnected", clientID)
	}

	return nil
}

func (h *hub[TRequest, TResponse]) Emit(ctx context.Context, identifier ClientIdentifier, payload TResponse) error {
	clientID := identifier.ID()

	cl, exist := h.clients[clientID]
	if !exist {
		return pkgerrors.New("client not found")
	}

	cl.Response() <- payload

	return nil
}

func (h *hub[TRequest, TResponse]) handleRequest(ctx context.Context, cl ClientConn[TRequest, TResponse]) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			req, err := readRequest[TRequest](ctx, cl.Conn())
			if err != nil {
				monitoring.FromContext(ctx).Errorf(err, "[websocket] read request failed")

				// TODO: consider skip this request
				cl.Error() <- err
				cl.Discard()
			}

			cl.Request() <- req
		}
	}
}

func (h *hub[TRequest, TResponse]) processMessage(ctx context.Context, cl ClientConn[TRequest, TResponse]) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-cl.Request():
			resp, err := h.handler(ctx, req)
			if err != nil {
				cl.Error() <- err
				return
			}

			cl.Response() <- resp
		}
	}
}

func (h *hub[TRequest, TResponse]) handleResponse(ctx context.Context, cl ClientConn[TRequest, TResponse]) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp := <-cl.Response():
			if err := writeResponse(ctx, cl.Conn(), resp); err != nil {
				monitoring.FromContext(ctx).Errorf(err, "[websocket] write response failed")
				return
			}
		}
	}
}

func (h *hub[TRequest, TResponse]) handleErrorResponse(ctx context.Context, cl ClientConn[TRequest, TResponse]) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-cl.Error():
			if err := writeResponse(ctx, cl.Conn(), err); err != nil {
				monitoring.FromContext(ctx).Errorf(err, "[websocket] write error response failed")
				cl.Discard()
			}
			return
		}
	}
}

func (h *hub[TRequest, TResponse]) pingClient(ctx context.Context, cl ClientConn[TRequest, TResponse]) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(defaultPingInterval):
			err := cl.Ping(ctx)
			if err == nil {
				monitoring.FromContext(ctx).Infof("[websocket] ping client %s", cl.Conn().Subprotocol())
				continue
			}

			monitoring.FromContext(ctx).Errorf(err, "[websocket] ping client failed")

			// If the ping fails, discard the client
			cl.Discard()
		}
	}
}

func (h *hub[TRequest, TResponse]) setClient(clientID string, client ClientConn[TRequest, TResponse]) {
	h.clientMux.Lock()
	defer h.clientMux.Unlock()

	h.clients[clientID] = client
}

func (h *hub[TRequest, TResponse]) getClient(clientID string) (ClientConn[TRequest, TResponse], error) {
	h.clientMux.RLock()
	defer h.clientMux.RUnlock()

	client, exists := h.clients[clientID]
	if !exists {
		return nil, pkgerrors.New("client not found")
	}
	return client, nil
}

func (h *hub[TRequest, TResponse]) existsClient(clientID string) bool {
	h.clientMux.RLock()
	defer h.clientMux.RUnlock()

	_, exists := h.clients[clientID]
	return exists
}

func (h *hub[TRequest, TResponse]) removeClient(clientID string) {
	h.clientMux.Lock()
	defer h.clientMux.Unlock()

	if _, exists := h.clients[clientID]; exists {
		delete(h.clients, clientID)
	}
}
