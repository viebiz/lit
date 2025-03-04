package lit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/testutil"

	"github.com/viebiz/lit/monitoring"
	"github.com/viebiz/lit/monitoring/tracing/mocktracer"
)

func TestRootMiddleware(t *testing.T) {
	tp := mocktracer.Start()
	defer tp.Stop()

	type handler struct {
		Method string
		Path   string
		Func   ErrHandlerFunc
	}
	tcs := map[string]struct {
		givenReq  *http.Request
		hdl       handler
		expStatus int
		expBody   string
		expLogs   []map[string]string
	}{
		"success - GET method": {
			givenReq: httptest.NewRequest(http.MethodGet, "/ping", nil),
			hdl: handler{
				Method: http.MethodGet,
				Path:   "/ping",
				Func: func(c Context) error {
					c.JSON(http.StatusOK, gin.H{"message": "pong"})
					return nil
				},
			},
			expStatus: http.StatusOK,
			expBody:   `{"message":"pong"}`,
			expLogs: []map[string]string{
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "Sentry DSN not provided. Not using Sentry Error Reporting", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "OTelExporter URL not provided. Not using Distributed Tracing", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "Wrote {\"message\":\"pong\"}", "server.name": "lightning", "environment": "dev", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001", "http.request.method": "GET", "server.address": "example.com", "user_agent": "", "url": "/ping", "network.peer.address": "192.0.2.1:1234", "network.protocol.version": "HTTP/1.1"},
				{"level": "INFO", "ts": "2025-02-23T18:23:26.434+0700", "msg": "http.incoming_request", "server.name": "lightning", "environment": "dev", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001", "http.request.method": "GET", "server.address": "example.com", "user_agent": "", "url": "/ping", "network.peer.address": "192.0.2.1:1234", "network.protocol.version": "HTTP/1.1", "http.response.status": "200", "http.response.size": "18"},
			},
		},
		"success - POST method": {
			givenReq: httptest.NewRequest(http.MethodPost, "/ping", bytes.NewBufferString(`{"message":"Hello lightning"}`)),
			hdl: handler{
				Method: http.MethodPost,
				Path:   "/ping",
				Func: func(c Context) error {
					var msg struct {
						Message string `json:"message"`
					}
					if err := c.Bind(&msg); err != nil {
						return err
					}

					c.JSON(http.StatusOK, msg)
					return nil
				},
			},
			expStatus: http.StatusOK,
			expBody:   `{"message":"Hello lightning"}`,
			expLogs: []map[string]string{
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "Sentry DSN not provided. Not using Sentry Error Reporting", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "OTelExporter URL not provided. Not using Distributed Tracing", "server.name": "lightning", "environment": "dev", "version": "1.0.0"},
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "Wrote {\"message\":\"Hello lightning\"}", "server.name": "lightning", "environment": "dev", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001", "http.request.method": "POST", "http.request.body.size": "29", "server.address": "example.com", "user_agent": "", "url": "/ping", "network.peer.address": "192.0.2.1:1234", "network.protocol.version": "HTTP/1.1"},
				{"level": "INFO", "ts": "2025-02-23T18:23:26.434+0700", "msg": "http.incoming_request", "server.name": "lightning", "environment": "dev", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001", "http.request.method": "POST", "http.request.body.size": "29", "http.request.body": "{\"message\":\"Hello lightning\"}", "server.address": "example.com", "user_agent": "", "url": "/ping", "network.peer.address": "192.0.2.1:1234", "network.protocol.version": "HTTP/1.1", "http.response.status": "200", "http.response.size": "29"},
			},
		},
		"error - Expected error": {
			givenReq: httptest.NewRequest(http.MethodPatch, "/ping", bytes.NewBufferString(`{"message":"pong"}`)),
			hdl: handler{
				Method: http.MethodPatch,
				Path:   "/ping",
				Func: func(c Context) error {
					return HttpError{Status: http.StatusBadRequest, Code: "validation_error", Desc: "Invalid request"}
				},
			},
			expStatus: http.StatusBadRequest,
			expBody:   "{\"error\":\"validation_error\",\"error_description\":\"Invalid request\"}",
			expLogs: []map[string]string{
				{"environment": "dev", "level": "INFO", "msg": "Sentry DSN not provided. Not using Sentry Error Reporting", "server.name": "lightning", "ts": "2025-02-23T18:43:12.5460700", "version": "1.0.0"},
				{"environment": "dev", "level": "INFO", "msg": "OTelExporter URL not provided. Not using Distributed Tracing", "server.name": "lightning", "ts": "2025-02-23T18:43:12.5460700", "version": "1.0.0"},
				{"environment": "dev", "http.request.body.size": "18", "http.request.method": "PATCH", "level": "INFO", "msg": `Wrote {"error":"validation_error","error_description":"Invalid request"}`, "user_agent": "", "url": "/ping", "network.peer.address": "192.0.2.1:1234", "network.protocol.version": "HTTP/1.1", "server.address": "example.com", "server.name": "lightning", "ts": "2025-02-23T18:43:12.5460700", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
				{"environment": "dev", "http.request.body": `{"message":"pong"}`, "http.request.body.size": "18", "http.request.method": "PATCH", "http.response.size": "66", "http.response.status": "400", "level": "INFO", "msg": "http.incoming_request", "user_agent": "", "url": "/ping", "network.peer.address": "192.0.2.1:1234", "network.protocol.version": "HTTP/1.1", "server.address": "example.com", "server.name": "lightning", "ts": "2025-02-23T18:43:12.5460700", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
			},
		},
		"error - PANIC request": {
			givenReq: httptest.NewRequest(http.MethodPatch, "/ping", bytes.NewBufferString(`{"message":"pong"}`)),
			hdl: handler{
				Method: http.MethodPatch,
				Path:   "/ping",
				Func: func(c Context) error {
					panic(errors.New("simulated panic"))
				},
			},
			expStatus: http.StatusInternalServerError,
			expBody:   "{\"error\":\"internal_server_error\",\"error_description\":\"Something went wrong\"}",
			expLogs: []map[string]string{
				{"environment": "dev", "level": "INFO", "msg": "Sentry DSN not provided. Not using Sentry Error Reporting", "server.name": "lightning", "ts": "2025-02-23T18:43:12.5460700", "version": "1.0.0"},
				{"environment": "dev", "level": "INFO", "msg": "OTelExporter URL not provided. Not using Distributed Tracing", "server.name": "lightning", "ts": "2025-02-23T18:43:12.5460700", "version": "1.0.0"},
				{"environment": "dev", "level": "ERROR", "msg": "Caught a panic", "http.request.body.size": "18", "http.request.method": "PATCH", "error.kind": "*errors.errorString", "error.message": "simulated panic", "user_agent": "", "url": "/ping", "network.peer.address": "192.0.2.1:1234", "network.protocol.version": "HTTP/1.1", "server.address": "example.com", "server.name": "lightning", "ts": "2025-02-23T18:43:12.5460700", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
				{"environment": "dev", "level": "INFO", "msg": `Wrote {"error":"internal_server_error","error_description":"Something went wrong"}`, "http.request.body.size": "18", "http.request.method": "PATCH", "user_agent": "", "url": "/ping", "network.peer.address": "192.0.2.1:1234", "network.protocol.version": "HTTP/1.1", "server.address": "example.com", "server.name": "lightning", "ts": "2025-02-23T18:43:12.5460700", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
			},
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()

			// Given
			logBuffer := bytes.NewBuffer(nil)
			m, err := monitoring.New(monitoring.Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: logBuffer})
			require.NoError(t, err)
			appCtx := monitoring.SetInContext(context.Background(), m)

			w := httptest.NewRecorder()
			route, ctx, handleRequest := NewRouterForTest(w)
			route.Use(rootMiddleware(appCtx))
			route.HandleWithErr(tc.hdl.Method, tc.hdl.Path, tc.hdl.Func)

			if slices.Contains([]string{http.MethodPost, http.MethodPut, http.MethodPatch}, tc.givenReq.Method) {
				tc.givenReq.Header.Set("Content-Type", "application/json")
			}
			ctx.SetRequest(tc.givenReq)

			// When
			handleRequest()

			// Then
			require.Equal(t, tc.expStatus, w.Code)
			require.Equal(t, tc.expBody, w.Body.String())
			pasedLogs, err := parseLog(logBuffer.Bytes())
			require.NoError(t, err)
			testutil.Equal(t, tc.expLogs, pasedLogs, testutil.IgnoreSliceMapEntries(func(k string, v string) bool {
				if k == "ts" {
					return true
				}

				if k == "error.stack" {
					return true
				}

				if v == "Caught a panic" {
					return true
				}

				return false
			}))
		})
	}
}

func parseLog(b []byte) ([]map[string]string, error) {
	var result []map[string]string
	for _, s := range strings.Split(string(b), "\n") {
		if s == "" {
			break
		}
		var r map[string]string
		if err := json.Unmarshal([]byte(s), &r); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

func BenchmarkRootMiddleware(b *testing.B) {
	// Given
	// Setup Monitor
	m, err := monitoring.New(monitoring.Config{ServerName: "lightning", Environment: "dev", Version: "1.0.0", Writer: io.Discard})
	require.NoError(b, err)

	// Init trace test
	tp := mocktracer.Start()
	defer tp.Stop()

	appCtx := monitoring.SetInContext(context.Background(), m)

	r, hdl := NewRouter()
	r.Use(rootMiddleware(appCtx)) // Add middleware for benchmark

	// Define a dummy route
	r.Post("/users", func(ctx Context) error {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "your are the best developer",
		})
		return nil
	})

	// Pre-create request to avoid redundant allocations inside the loop
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"username":"the-witcher-knight"}`))
	req.Header.Set("Content-Type", "application/json")

	// When
	b.ReportAllocs() // Track memory allocations
	b.ResetTimer()   // Reset timer to avoid setup overhead
	for i := 0; i < b.N; i++ {
		// Create a response recorder
		w := httptest.NewRecorder()
		hdl.ServeHTTP(w, req)
	}
}
