package websocket

import (
	"time"
)

// HubOption represents option for creates websocket hub
type HubOption[TRequest Request, TResponse Response] func(Hub[TRequest, TResponse])

// WithPingInterval sets the client identifier for the hub
func WithPingInterval[TRequest Request, TResponse Response](interval time.Duration) HubOption[TRequest, TResponse] {
	return func(h Hub[TRequest, TResponse]) {
		impl, ok := h.(*hub[TRequest, TResponse])
		if !ok {
			return
		}

		impl.settings.pingInterval = interval
	}
}

// WithPingTimeout sets the ping timeout for the hub
func WithPingTimeout[TRequest Request, TResponse Response](timeout time.Duration) HubOption[TRequest, TResponse] {
	return func(h Hub[TRequest, TResponse]) {
		impl, ok := h.(*hub[TRequest, TResponse])
		if !ok {
			return
		}

		impl.settings.pingTimeout = timeout
	}
}

// WithHeaderBasedClientIdentifier sets the client identifier for the hub
func WithHeaderBasedClientIdentifier[TRequest Request, TResponse Response](headerName string) HubOption[TRequest, TResponse] {
	return func(h Hub[TRequest, TResponse]) {
		impl, ok := h.(*hub[TRequest, TResponse])
		if !ok {
			return
		}

		impl.settings.clientIdentifier = NewHeaderClientIdentifier(headerName)
	}
}

// WithQueryBasedClientIdentifier sets the client identifier for the hub
func WithQueryBasedClientIdentifier[TRequest Request, TResponse Response](paramName string) HubOption[TRequest, TResponse] {
	return func(h Hub[TRequest, TResponse]) {
		impl, ok := h.(*hub[TRequest, TResponse])
		if !ok {
			return
		}

		impl.settings.clientIdentifier = NewParamClientIdentifier(paramName)
	}
}

// WithSerializer sets the serializer for the hub
func WithSerializer[TRequest Request, TResponse Response](serializer Serializer) HubOption[TRequest, TResponse] {
	return func(h Hub[TRequest, TResponse]) {
		impl, ok := h.(*hub[TRequest, TResponse])
		if !ok {
			return
		}

		impl.settings.serializer = serializer
	}
}

// WithReadLimit sets the read limit for the hub
func WithReadLimit[TRequest Request, TResponse Response](limit int64) HubOption[TRequest, TResponse] {
	return func(h Hub[TRequest, TResponse]) {
		impl, ok := h.(*hub[TRequest, TResponse])
		if !ok {
			return
		}

		impl.settings.readLimit = limit
	}
}
