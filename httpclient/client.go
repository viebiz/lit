package httpclient

import (
	"net/http"

	"github.com/viebiz/lit/monitoring"
)

type Client struct {
	underlyingClient *http.Client

	extSvcInfo monitoring.ExternalServiceInfo

	url string

	userAgent string

	contentType string

	header http.Header

	timeoutAndRetryOption timeoutAndRetryOption

	disableReqBodyLogging bool

	disableRespBodyLogging bool
}
