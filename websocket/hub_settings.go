package websocket

import (
	"time"
)

type HubSettings struct {
	clientIdentifier ClientIdentifier
	serializer       Serializer

	pingInterval, pingTimeout time.Duration
	readLimit                 int64
}

func getDefaultHubSettings() HubSettings {
	return HubSettings{
		clientIdentifier: NewParamClientIdentifier(""),
		serializer:       NewJSONSerializer(),
		pingInterval:     defaultPingInterval,
		pingTimeout:      defaultPingTimeout,
		readLimit:        defaultReadLimit,
	}
}
