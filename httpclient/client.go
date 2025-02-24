package httpclient

import (
	"net/http"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"github.com/viebiz/lit/monitoring"
)

// A Client describes an HTTP endpoint's client. This client is
// mainly used to send HTTP request based on the request configuration set into
// the client and the response handling logic configured within the client.
type Client struct {
	// HTTP client to be used to execute HTTP call.
	underlyingClient *http.Client

	extSvcInfo monitoring.ExternalServiceInfo

	// URL for this Client. This value can contain path
	// variable placeholders for substitution when sending the HTTP request
	url string

	// HTTP method
	method string

	// Name of the service to call. Together with serviceName will form the label in logger fields & error code prefixes.
	serviceName string

	// User agent value (i.e. RFC7231)
	userAgent string

	// Content MIME type
	// Default: application/json
	contentType string

	// Default request header configuration
	header header

	timeoutAndRetryOption timeoutAndRetryOption

	// Disable request body logging
	// Default: false,
	disableReqBodyLogging bool

	// Disable response body logging
	// Default: false,
	disableRespBodyLogging bool

	// Disable request/response log redaction
	// Default: false,
	disableLogRedaction bool
}

// header describes the HTTP request headers
type header struct {
	// Other request header values
	// Note: Can be overridden on Payload
	// Default: nil
	values map[string]string
}

// newClient method creates a new client
func newClient(
	underlyingClient *http.Client,
	url,
	method,
	serviceName string,
	opts ...ClientOption,
) (*Client, error) {
	c := &Client{
		underlyingClient: underlyingClient,
		timeoutAndRetryOption: timeoutAndRetryOption{
			maxRetries:         defaultMaxRetriesOnErrOrTimeout,
			maxWaitPerTry:      defaultTimeoutPerTry,
			maxWaitInclRetries: defaultMaxWaitInclRetries,
			onTimeout:          defaultRetryOnTimeout,
			onStatusCodes:      make(map[int]bool),
		},
		contentType: defaultContentType,
	}

	c.url = strings.TrimSpace(url)
	if c.url == "" {
		return nil, pkgerrors.WithStack(ErrMissingURL)
	}
	c.method = strings.TrimSpace(method)
	if c.method == "" {
		return nil, pkgerrors.WithStack(ErrMissingMethod)
	}
	c.serviceName = strings.TrimSpace(serviceName)
	if c.serviceName == "" {
		return nil, pkgerrors.WithStack(ErrMissingServiceName)
	}

	for _, opt := range opts {
		opt(c)
	}

	if err := c.timeoutAndRetryOption.IsValid(); err != nil {
		return nil, err
	}

	// TODO: Setup user agent based on configs, so use default go client user-agent
	//c.userAgent = fmt.Sprintf(
	//	"%s/%s (%s)",
	//	appCfg.AppName,
	//	appCfg.Version,
	//	appCfg.Server,
	//)
	c.extSvcInfo = monitoring.NewExternalServiceInfo(c.url)

	return c, nil
}
