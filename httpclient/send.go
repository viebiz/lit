package httpclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	pkgerrors "github.com/pkg/errors"
	"github.com/viebiz/lit/monitoring"
	"github.com/viebiz/lit/monitoring/instrumenthttp"
)

var (
	execWithRetryFunc = execWithRetry
	readRespBodyFunc  = readRespBody
)

// Payload is the request payload struct representation
type Payload struct {
	// Request body
	Body []byte
	// QueryParams contains the request/query parameters
	QueryParams url.Values
	// PathVars contains the path variables used to replace placeholders
	// wrapped with {} in Client.URL
	PathVars map[string]string
	// Header contains custom request headers that will be added to the request
	// on http call.
	// The values on this field will override Client.Headers.Values
	Header map[string]string
}

// Response is the result of the http call
type Response struct {
	Status int
	Body   []byte
	Header http.Header
}

// Send executes an HTTP call based on the information and configuration
// in the Client
func (c *Client) Send(ctx context.Context, p Payload) (Response, error) {
	// Endpoint URL
	endpointURL := c.constructURL(p)

	var err error
	ctx, segEnd := instrumenthttp.StartOutgoingGroupSegment(ctx, c.extSvcInfo, c.serviceName, c.method, c.url)
	defer func() { segEnd(err) }()
	monitor := monitoring.FromContext(ctx)

	if !c.disableReqBodyLogging && (c.method == http.MethodPost || c.method == http.MethodPut || c.method == http.MethodPatch) {
		v := p.Body

		monitor.Infof("[ext_http_req] request body:(%s)", string(v))
	}

	// Create context with max timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, c.timeoutAndRetryOption.maxWaitInclRetries)
	defer cancel()

	// HTTP operation
	resp, err := c.execute(ctxTimeout, endpointURL, p)
	if err != nil {
		return Response{}, err
	}

	if !c.disableRespBodyLogging {
		v := resp.Body

		monitor.Infof("[ext_http_req] response body:(%s)", v)
	} else {
		monitor.Infof("[ext_http_req] skipping logging resp body")
	}

	return resp, nil
}

func (c *Client) execute(
	ctx context.Context,
	endpointURL string,
	p Payload,
) (Response, error) {
	var resultResp Response
	var attempts int

	if err := execWithRetryFunc(ctx, c.timeoutAndRetryOption.maxRetries, c.timeoutAndRetryOption.maxWaitInclRetries,
		func() error {
			// start new attempt for the http request
			attempts++

			req, err := c.createHTTPRequest(endpointURL, p.Body) // create HTTP request
			if err != nil {
				return err
			}

			// var err error
			var status int
			reqCtx, segEnd := instrumenthttp.StartOutgoingSegment(ctx, c.extSvcInfo, c.serviceName, req)
			defer func() { segEnd(status, err) }()
			monitor := monitoring.FromContext(ctx)

			reqCtx, cancelReqCtx := context.WithTimeout(reqCtx, c.timeoutAndRetryOption.maxWaitPerTry)
			defer cancelReqCtx()
			req = req.WithContext(reqCtx) // limit each HTTP request timeout option for per try

			c.setHeader(req, p) // set request headers

			// start sending request
			start := time.Now()
			monitor.Infof("[ext_http_req] start new attempt (%d)", attempts)

			resp, err := c.underlyingClient.Do(req)

			monitor = monitor.
				WithTag("ext_http_resp_duration", fmt.Sprintf("%dms", time.Since(start).Milliseconds()))
			if err != nil {
				// handle request error
				monitor.Infof("[ext_http_req] end with error: (%+v), attempt (%d)", err, attempts) // intentionally not using Errorf for err as this is just info and the req can be retried.

				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					return backoff.Permanent(ErrOverflowMaxWait) // stop retry by returning backoff.Permanent error
				}

				// evaluate if err is caused by connection timeout
				uerr, ok := err.(*url.Error)
				if !ok || !uerr.Timeout() {
					if errors.Is(err, context.Canceled) {
						return backoff.Permanent(ErrOperationContextCanceled) // stop retry by returning backoff.Permanent error
					}
					return pkgerrors.WithStack(err)
				}

				// check if we need retry on timeout or not
				if !c.timeoutAndRetryOption.onTimeout {
					return backoff.Permanent(ErrTimeout) // stop retry by returning backoff.Permanent error
				}

				return ErrTimeout
			}

			status = resp.StatusCode
			if _, ok := c.timeoutAndRetryOption.onStatusCodes[resp.StatusCode]; ok {
				// Only returns error for backoff retry function retries the call
				// when attempts count haven't reached max retries
				if uint64(attempts) <= c.timeoutAndRetryOption.maxRetries {
					monitor.Infof("[ext_http_req] retry on status code: (%d), attempt (%d)", resp.StatusCode, attempts)
					return fmt.Errorf("retry on status code %v", resp.StatusCode)
				}
			}

			monitor.Infof("[ext_http_req] end with status code: (%d), attempt (%d)", resp.StatusCode, attempts)

			// attempt to read response body
			// we need to read the response body in the same function where we cancel the retry context
			// otherwise, for big payload, the context will be cancelled while reading the body
			monitor.Infof("[ext_http_req] attempting to read body: attempt (%d)", attempts)
			respBody, err := readRespBodyFunc(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				monitor.Infof("[ext_http_req] err reading from body: (%+v), attempt (%d)", err, attempts)
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					return backoff.Permanent(ErrOverflowMaxWait) // stop retry by returning backoff.Permanent error
				}

				// evaluate if err is caused by connection timeout
				if errors.Is(err, context.Canceled) {
					return backoff.Permanent(ErrOperationContextCanceled) // stop retry by returning backoff.Permanent error
				}
				if errors.Is(err, context.DeadlineExceeded) {
					return ErrTimeout
				}

				return pkgerrors.WithStack(err)
			}

			resultResp.Status = resp.StatusCode
			resultResp.Body = respBody
			resultResp.Header = resp.Header

			return nil
		}); err != nil {
		switch err {
		case context.Canceled:
			return Response{}, ErrOperationContextCanceled
		case context.DeadlineExceeded:
			return Response{}, ErrOverflowMaxWait
		default:
			return Response{}, err
		}
	}

	return resultResp, nil
}

// constructURL returns the full URL with query params and path variables substitution
func (c *Client) constructURL(p Payload) string {
	u := c.url
	// Replace path variables
	for k, v := range p.PathVars {
		u = strings.Replace(u, fmt.Sprintf(":%s", k), v, -1)
	}
	// Add query params
	if q := p.QueryParams.Encode(); q != "" {
		sep := "?"
		if strings.Contains(u, "?") {
			sep = "&"
		}
		u = u + sep + q
	}
	return u
}

// setHeader sets the request headers based on the resource client configuration and payload
func (c *Client) setHeader(r *http.Request, p Payload) {
	// Set default request headers
	if c.userAgent != "" {
		r.Header.Set("User-Agent", c.userAgent)
	}
	if c.contentType != "" {
		r.Header.Set("Content-Type", c.contentType)
	}
	for k, v := range c.header.values { // Resource client default headers
		r.Header.Set(k, v)
	}
	for k, v := range p.Header { // Payload headers
		r.Header.Set(k, v)
	}
}

func (c *Client) createHTTPRequest(endpointURL string, body []byte) (*http.Request, error) {
	var b io.Reader
	if len(body) > 0 {
		b = bytes.NewBuffer(body)
	}

	r, err := http.NewRequest(c.method, endpointURL, b)
	return r, pkgerrors.WithStack(err)
}

func readRespBody(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
