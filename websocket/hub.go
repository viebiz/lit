package websocket

import (
	"context"
	"errors"
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

	settings HubSettings
	info     HubInfo
}

func (h *hub[TRequest, TResponse]) On(router lit.Router) {
	router.Get(h.baseURL, h.handle)
}

func (h *hub[TRequest, TResponse]) handle(c lit.Context) error {
	// 1. Get the client id from the request
	clientID, err := h.settings.clientIdentifier.Extract(c)
	if err != nil {
		return err
	}

	// 2. Check if the client is already connected
	if exists := h.existsClient(clientID); exists {
		return ErrClientAlreadyExists
	}

	// 3. Upgrades the connection to a WebSocket
	conn, err := websocket.Accept(c.Writer(), c.Request(), &websocket.AcceptOptions{})
	if err != nil {
		return pkgerrors.Wrap(err, "accept websocket connection")
	}

	// 4. Create a new client connection
	ctx, cancel := context.WithCancel(monitoring.NewContext(c))
	defer cancel()

	clientCfg := ClientConfig{
		ID:        clientID,
		ReadLimit: h.settings.readLimit,
	}
	clc := initClientConn[TRequest, TResponse](ctx, clientCfg, cancel, conn)
	h.setClient(clientID, clc)

	// 5. Ping the client
	go h.pingClient(clc)

	// 6. Handle incoming messages
	go h.handleRequest(clc)

	// 7. Process messages
	go h.processMessage(clc)

	// 8.1. Handle outgoing messages
	go h.handleResponse(clc)

	// 8.2. Handle errors
	go h.handleErrorResponse(clc)

	// 9. Handle connection close
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

type EmitOption struct {
	//SessionID string
}

func (h *hub[TRequest, TResponse]) Emit(ctx context.Context, id string, payload TResponse, opts EmitOption) error {
	clc, err := h.getClient(id)
	if err != nil {
		return err
	}

	clc.Response() <- payload

	return nil
}

func (h *hub[TRequest, TResponse]) Broadcast(ctx context.Context, payload TResponse, opts EmitOption) error {
	h.clientMux.RLock()
	defer h.clientMux.RUnlock()

	var sendErr error
	for _, client := range h.clients {
		// Emit to each client, if it fails, log the error and continue
		if err := h.Emit(ctx, client.ID(), payload, opts); err != nil {
			sendErr = errors.Join(sendErr, err)
			continue // Go to the next client
		}
	}
	if sendErr != nil {
		return pkgerrors.WithStack(sendErr)
	}

	return nil
}

func (h *hub[TRequest, TResponse]) handleRequest(clc ClientConn[TRequest, TResponse]) {
	for {
		select {
		case <-clc.Context().Done():
			return
		default:
			ctx := monitoring.NewContext(clc.Context())

			var req TRequest
			if err := h.readMessage(ctx, clc.Conn(), &req); err != nil {
				// Check first if it's specifically a close error
				var closeErr websocket.CloseError
				if errors.As(err, &closeErr) {
					// Attempt tp close the connection
					clc.Discard()

					// Log the close error
					if closeErr.Code != websocket.StatusNormalClosure {
						monitoring.FromContext(clc.Context()).Errorf(err, "[websocket] got unexpected close message")
					}
					return
				}

				clc.Error() <- err
			}

			clc.Request() <- req
		}
	}
}

func (h *hub[TRequest, TResponse]) processMessage(clc ClientConn[TRequest, TResponse]) {
	for {
		select {
		case <-clc.Context().Done():
			return
		case req := <-clc.Request():
			ctx := monitoring.NewContext(clc.Context())

			resp, err := h.handler(ctx, req)
			if err != nil {
				clc.Error() <- err
				continue
			}

			clc.Response() <- resp
		}
	}
}

func (h *hub[TRequest, TResponse]) handleResponse(clc ClientConn[TRequest, TResponse]) {
	for {
		select {
		case <-clc.Context().Done():
			return
		case resp := <-clc.Response():
			ctx := monitoring.NewContext(clc.Context())
			if err := h.writeText(ctx, clc.Conn(), resp); err != nil {
				clc.Error() <- err
			}
		}
	}
}

func (h *hub[TRequest, TResponse]) handleErrorResponse(clc ClientConn[TRequest, TResponse]) {
	for {
		select {
		case <-clc.Context().Done():
			return
		case err := <-clc.Error():
			monitoring.FromContext(clc.Context()).Errorf(err, "[websocket] got error")

			return
		}
	}
}

func (h *hub[TRequest, TResponse]) pingClient(clc ClientConn[TRequest, TResponse]) {
	for {
		select {
		case <-clc.Context().Done():
			return
		case <-time.After(defaultPingInterval):
			ctx, cancel := context.WithTimeout(clc.Context(), defaultPingTimeout)

			err := clc.Ping(ctx)
			if err == nil {
				monitoring.FromContext(ctx).Infof("[websocket] ping client %s successfully", clc.ID())
				cancel()

				continue
			}

			monitoring.FromContext(ctx).Errorf(err, "[websocket] ping client failed")
			cancel()

			// If the ping fails, discard the client
			clc.Discard()
		}
	}
}

func (h *hub[TRequest, TResponse]) readMessage(ctx context.Context, c Conn, v any) error {
	_, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	if err := h.settings.serializer.Deserialize(r, v); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}

func (h *hub[TRequest, TResponse]) writeText(ctx context.Context, c Conn, v any) error {
	wrt, err := c.Writer(ctx, MessageTypeText)
	if err != nil {
		return pkgerrors.WithStack(err)
	}
	defer wrt.Close()

	if err := h.settings.serializer.Serialize(v, wrt); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}

func (h *hub[TRequest, TResponse]) setClient(clientID string, client ClientConn[TRequest, TResponse]) {
	h.clientMux.Lock()
	defer h.clientMux.Unlock()

	h.clients[clientID] = client

	// Update info
	h.info.ClientCount++
}

func (h *hub[TRequest, TResponse]) getClient(clientID string) (ClientConn[TRequest, TResponse], error) {
	h.clientMux.RLock()
	defer h.clientMux.RUnlock()

	client, exists := h.clients[clientID]
	if !exists {
		return nil, pkgerrors.WithStack(ErrClientNotFound)
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

	delete(h.clients, clientID)

	// Update info
	h.info.ClientCount--
}
