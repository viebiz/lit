package lightning

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"

	"github.com/viebiz/lit/monitoring"
)

func TestRootMiddleware(t *testing.T) {
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
		expLogs   []map[string]interface{}
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
			expLogs: []map[string]interface{}{
				{
					"level":    "info",
					"msg":      `Wrote {"message":"pong"}`,
					"span_id":  "0000000000000001",
					"trace_id": "00000000000000000000000000000001",
				},
				{
					"http.request.endpoint": "/ping",
					"http.request.method":   "GET",
					"http.response.size":    float64(18),
					"http.response.status":  float64(200),
					"level":                 "info",
					"msg":                   "http.incoming_request",
					"span_id":               "0000000000000001",
					"trace_id":              "00000000000000000000000000000001",
				},
			},
		},
		"success - POST method": {
			givenReq: httptest.NewRequest(http.MethodPost, "/ping", bytes.NewBufferString(`{"message":"pong"}`)),
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
			expBody:   `{"message":"pong"}`,
			expLogs: []map[string]interface{}{
				{
					"level":    "info",
					"msg":      `Wrote {"message":"pong"}`,
					"span_id":  "0000000000000001",
					"trace_id": "00000000000000000000000000000001",
				},
				{
					"http.request.body":     `{"message":"pong"}`,
					"http.request.endpoint": "/ping",
					"http.request.method":   "POST",
					"http.response.size":    float64(18),
					"http.response.status":  float64(200),
					"level":                 "info",
					"msg":                   "http.incoming_request",
					"span_id":               "0000000000000001",
					"trace_id":              "00000000000000000000000000000001",
				},
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
			expBody:   ErrInternalServerError.Error(),
			expLogs: []map[string]interface{}{
				{
					"error": "simulated panic",
					"level": "error",
					"msg": `
						"""
						Caught a panic: goroutine 7 [running]:
						runtime/debug.Stack()
							/Users/locdang/sdk/go1.23.3/src/runtime/debug/stack.go:26 +0x64
						gitlab.com/bizgroup2/lightning.TestRootMiddleware.func4.rootMiddleware.2.1()
							/Users/locdang/IdeaProjects/playground/lightning/middleware_root.go:29 +0xc4
						panic({0x1013bf000?, 0x14000238900?})
							/Users/locdang/sdk/go1.23.3/src/runtime/panic.go:785 +0x124
						gitlab.com/bizgroup2/lightning.TestRootMiddleware.func3({0x1400024adc0?, 0x101ffe978?})
							/Users/locdang/IdeaProjects/playground/lightning/middleware_root_test.go:113 +0x50
						gitlab.com/bizgroup2/lightning.router.Handle.handleHttpError.func1(0x1400020a200)
							/Users/locdang/IdeaProjects/playground/lightning/handler_func.go:20 +0x34
						github.com/gin-gonic/gin.(*Context).Next(...)
							/Users/locdang/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/context.go:185
						gitlab.com/bizgroup2/lightning.TestRootMiddleware.func4.rootMiddleware.2({0x1014eb488, 0x1400020a200})
							/Users/locdang/IdeaProjects/playground/lightning/middleware_root.go:45 +0x14c
						... // 16 elided lines
					`,
					"span_id":  "0000000000000001",
					"trace_id": "00000000000000000000000000000001",
				},
				{
					"level":    "info",
					"msg":      `Wrote {"error":"internal_server_error","error_description":"internal server error"}`,
					"span_id":  "0000000000000001",
					"trace_id": "00000000000000000000000000000001",
				},
			},
		},
	}

	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			// Given
			monitor, endTest := monitoring.NewMonitorTest()
			defer endTest()

			w := httptest.NewRecorder()

			route, ctx, handleRequest := NewRouterForTest(w)
			route.Use(rootMiddleware(monitor.Context()))
			route.Handle(tc.hdl.Method, tc.hdl.Path, tc.hdl.Func)

			if slices.Contains([]string{http.MethodPost, http.MethodPut, http.MethodPatch}, tc.givenReq.Method) {
				tc.givenReq.Header.Set("Content-Type", "application/json")
			}
			ctx.SetRequest(tc.givenReq)

			// When
			handleRequest()

			// Then
			require.Equal(t, tc.expStatus, w.Code)
			require.Equal(t, tc.expBody, w.Body.String())

			if diff := cmp.Diff(tc.expLogs, monitor.GetLogs(t), cmpopts.IgnoreMapEntries(func(key string, value interface{}) bool {
				if key == "msg" && strings.HasPrefix(value.(string), "Caught a panic") {
					return true
				}

				return key == "ts"
			})); diff != "" {
				t.Errorf("unexpected result (-want, got)\n%s", diff)
			}
		})
	}
}

func BenchmarkRootMiddleware(b *testing.B) {
	// Start a new HTTP server for test
	const srvAddr = "localhost:1604"
	go func() {
		logger := monitoring.NewLoggerWithWriter(bytes.NewBuffer(nil))
		appCtx := monitoring.SetInContext(context.Background(), logger)
		srv := NewHttpServer(appCtx, srvAddr, func(r Router) {
			r.Post("/weather", func(c Context) error {
				req := new(map[string]interface{})
				if err := c.Bind(&req); err != nil {
					return err
				}

				c.JSON(http.StatusOK, req)
				return nil
			})
		})

		require.NoError(b, srv.start(context.Background()))
	}()

	b.Run("incoming_http", func(b *testing.B) {
		b.Helper()
		b.ReportAllocs()
		b.ResetTimer()

		for idx := 0; idx < b.N; idx++ {
			reqBody, err := os.ReadFile("testdata/incoming_http_request_body.json")
			require.NoError(b, err)

			if _, err := http.Post("http://"+srvAddr+"/weather", "application/json", bytes.NewBuffer(reqBody)); err != nil {
				b.Error(err)
			}
		}
	})
}
