package httpclient

import (
	"strings"
)

// ClientOption alters behaviour of the Client
type ClientOption func(c *Client)

// OverrideBaseRequestHeaders method sets a map of request header key-value pairs to the request raised from client
func OverrideBaseRequestHeaders(m map[string]string) ClientOption {
	return func(c *Client) {
		if c.header.values == nil {
			c.header.values = map[string]string{}
		}
		for k, v := range m {
			c.header.values[k] = v
		}
	}
}

// OverrideContentType method override default the content type into the request header
func OverrideContentType(contentType string) ClientOption {
	return func(c *Client) {
		c.contentType = strings.TrimSpace(contentType)
	}
}

// OverrideUserAgent method override default the user-agent into the request header
func OverrideUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// DisableRequestBodyLogging method disables the default behaviour of logging request body
func DisableRequestBodyLogging() ClientOption {
	return func(c *Client) {
		c.disableReqBodyLogging = true
	}
}

// DisableResponseBodyLogging method disables the default behaviour of logging response body
func DisableResponseBodyLogging() ClientOption {
	return func(c *Client) {
		c.disableRespBodyLogging = true
	}
}
