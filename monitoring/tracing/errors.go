package tracing

import (
	"errors"
)

var (
	ErrMissingExporterURL   = errors.New("missing exporter url")
	ErrInvalidTransportType = errors.New("invalid transport type")
)
