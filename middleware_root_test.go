package lit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/viebiz/lit/mocks/mocktracer"
	"github.com/viebiz/lit/testutil"

	"github.com/viebiz/lit/monitoring"
)

func TestRootMiddleware(t *testing.T) {
	tp := mocktracer.NewTestTracerProvider(t)
	defer tp.Shutdown(t)

	type handler struct {
		Method string
		Path   string
		Func   HandlerFunc
	}
	tcs := map[string]struct {
		givenReq  *http.Request
		hdl       handler
		expStatus int
		expBody   string
		expLogs   []map[string]interface{}
	}{
		"success - GET method": {
			givenReq: httptest.NewRequest(http.MethodGet, "/ping", nil),
			hdl: handler{
				Method: http.MethodGet,
				Path:   "/ping",
				Func: func(c Context) error {
					return c.JSON(http.StatusOK, gin.H{"message": "pong"})
				},
			},
			expStatus: http.StatusOK,
			expBody:   "{\"message\":\"pong\"}\n",
			expLogs: []map[string]interface{}{
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "[incoming_request] Wrote response", "http.request.method": "GET", "url.path": "/ping", "http.response.body": map[string]any{"message": "pong"}, "server.name": "lightning", "environment": "dev", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
				{"level": "INFO", "ts": "2025-02-23T18:23:26.434+0700", "msg": "http.incoming_request", "http.request.method": "GET", "url.path": "/ping", "url.query": "", "http.response.body.size": float64(19), "http.response.status_code": float64(200), "server.name": "lightning", "environment": "dev", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
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

					return c.JSON(http.StatusOK, msg)
				},
			},
			expStatus: http.StatusOK,
			expBody:   "{\"message\":\"Hello lightning\"}\n",
			expLogs: []map[string]interface{}{
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "[incoming_request] Wrote response", "http.request.method": "POST", "url.path": "/ping", "http.response.body": map[string]any{"message": "Hello lightning"}, "server.name": "lightning", "environment": "dev", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
				{"level": "INFO", "ts": "2025-02-23T18:23:26.434+0700", "msg": "http.incoming_request", "http.request.method": "POST", "url.path": "/ping", "url.query": "", "http.request.body": map[string]any{"message": "Hello lightning"}, "http.response.body.size": float64(30), "http.response.status_code": float64(200), "server.name": "lightning", "environment": "dev", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
			},
		},
		"error - Expected error": {
			givenReq: httptest.NewRequest(http.MethodPatch, "/ping", bytes.NewBufferString(`{"message":"pong"}`)),
			hdl: handler{
				Method: http.MethodPatch,
				Path:   "/ping",
				Func: func(c Context) error {
					return &HTTPError{Status: http.StatusBadRequest, Code: "validation_error", Desc: "Invalid request"}
				},
			},
			expStatus: http.StatusBadRequest,
			expBody:   "{\"error\":\"validation_error\",\"error_description\":\"Invalid request\"}\n",
			expLogs: []map[string]interface{}{
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "[incoming_request] Wrote response", "http.request.method": "PATCH", "url.path": "/ping", "http.response.body": map[string]any{"error": "validation_error", "error_description": "Invalid request"}, "server.name": "lightning", "environment": "dev", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
				{"level": "INFO", "ts": "2025-02-23T18:23:26.434+0700", "msg": "http.incoming_request", "http.request.method": "PATCH", "url.path": "/ping", "url.query": "", "http.request.body": map[string]any{"message": "pong"}, "http.response.body.size": float64(67), "http.response.status_code": float64(400), "server.name": "lightning", "environment": "dev", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
			},
		},
		"error - Invalid request": {
			givenReq: httptest.NewRequest(http.MethodPatch, "/ping", bytes.NewBufferString(`{"message":"pong"}`)),
			hdl: handler{
				Method: http.MethodPatch,
				Path:   "/ping",
				Func: func(c Context) error {
					return errors.New("simulated server error")
				},
			},
			expStatus: http.StatusInternalServerError,
			expBody:   "{\"error\":\"Internal Server Error\",\"error_description\":\"Internal Server Error\"}\n",
			expLogs: []map[string]interface{}{
				{"level": "ERROR", "ts": "2025-09-11T10:04:33.991+0700", "msg": "Server error. Err: simulated server error", "span_id": "0000000000000001", "version": "1.0.0", "server.name": "lightning", "environment": "dev", "trace_id": "00000000000000000000000000000001", "error.kind": "*errors.errorString", "error.message": "simulated server error"},
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "[incoming_request] Wrote response", "http.request.method": "PATCH", "url.path": "/ping", "http.response.body": map[string]any{"error": "Internal Server Error", "error_description": "Internal Server Error"}, "server.name": "lightning", "environment": "dev", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
				{"level": "INFO", "ts": "2025-09-11T10:04:33.992+0700", "msg": "http.incoming_request", "http.request.method": "PATCH", "url.path": "/ping", "url.query": "", "http.request.body": map[string]any{"message": "pong"}, "http.response.status_code": float64(500), "http.response.body.size": float64(78), "environment": "dev", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001", "version": "1.0.0", "server.name": "lightning"},
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
			expBody:   "{\"error\":\"Internal Server Error\",\"error_description\":\"Internal Server Error\"}\n",
			expLogs: []map[string]interface{}{
				{"environment": "dev", "level": "ERROR", "msg": "Caught a panic", "error.kind": "*errors.errorString", "error.message": "simulated panic", "server.name": "lightning", "ts": "2025-02-23T18:43:12.5460700", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
				{"level": "INFO", "ts": "2025-02-23T18:18:48.186+0700", "msg": "[incoming_request] Wrote response", "http.request.method": "PATCH", "url.path": "/ping", "http.response.body": map[string]any{"error": "Internal Server Error", "error_description": "Internal Server Error"}, "server.name": "lightning", "environment": "dev", "version": "1.0.0", "trace_id": "00000000000000000000000000000001", "span_id": "0000000000000001"},
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
			monitorCtx := monitoring.SetInContext(context.Background(), m)

			w := httptest.NewRecorder()
			r, ctx, handleRequest := NewRouterForTest(w)
			r.Use(rootMiddleware(monitorCtx))
			r.Handle(tc.hdl.Method, tc.hdl.Path, tc.hdl.Func)

			if slices.Contains([]string{http.MethodPost, http.MethodPut, http.MethodPatch}, tc.givenReq.Method) {
				tc.givenReq.Header.Set("Content-Type", "application/json")
			}
			ctx.SetRequest(tc.givenReq)

			// When
			handleRequest()

			// Then
			require.Equal(t, tc.expStatus, w.Code)
			require.Equal(t, tc.expBody, w.Body.String())
			parsedLogs, err := parseLog(logBuffer.Bytes(), 2) // Skip 2 init log
			require.NoError(t, err)
			testutil.Equal(t, tc.expLogs, parsedLogs, testutil.IgnoreSliceMapEntries(func(k string, v interface{}) bool {
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

func parseLog(b []byte, skip int) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	for idx, s := range strings.Split(string(b), "\n") {
		if s == "" {
			break
		}
		if idx < skip {
			continue // Go to next line
		}
		var r map[string]interface{}
		if err := json.Unmarshal([]byte(s), &r); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}
