package websocket

import (
	pkgerrors "github.com/pkg/errors"
	"github.com/viebiz/lit"
)

// ClientIdentifier defines how client IDs are extracted from a request context
type ClientIdentifier interface {
	Extract(ctx lit.Context) (string, error)
}

// HeaderClientIdentifier extracts client ID from HTTP headers
type HeaderClientIdentifier struct {
	HeaderName string
}

func NewHeaderClientIdentifier(headerName string) *HeaderClientIdentifier {
	if headerName == "" {
		headerName = defaultClientIDHeaderName
	}

	return &HeaderClientIdentifier{HeaderName: headerName}
}

func (h *HeaderClientIdentifier) Extract(ctx lit.Context) (string, error) {
	clientID := ctx.Request().Header.Get(h.HeaderName)
	if clientID == "" {
		return "", pkgerrors.Errorf("client ID not found in header %s", h.HeaderName)
	}
	return clientID, nil
}

// ParamClientIdentifier extracts client ID from URL parameters
type ParamClientIdentifier struct {
	ParamName string
}

func NewParamClientIdentifier(paramName string) *ParamClientIdentifier {
	if paramName == "" {
		paramName = defaultClientIDParamName
	}
	return &ParamClientIdentifier{ParamName: paramName}
}

func (p *ParamClientIdentifier) Extract(ctx lit.Context) (string, error) {
	clientID := ctx.Param(p.ParamName)
	if clientID == "" {
		return "", pkgerrors.Errorf("client ID not found in param %s", p.ParamName)
	}
	return clientID, nil
}
