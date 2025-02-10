package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/viebiz/lit"
)

func TestRequestIDMiddleware(t *testing.T) {
	const staticRequestID = "mock-uuid-for-test"

	type handler struct {
		Method string
		Path   string
		Func   lightning.ErrHandlerFunc
	}
	type arg struct {
		givenReq  *http.Request
		hdl       handler
		expStatus int
		expBody   string
	}
	tcs := map[string]arg{
		"success - GET method": {
			givenReq: httptest.NewRequest(http.MethodGet, "/ping", nil),
			hdl: handler{
				Method: http.MethodGet,
				Path:   "/ping",
				Func: func(c lightning.Context) error {
					c.JSON(http.StatusOK, gin.H{"message": "pong"})
					return nil
				},
			},
			expStatus: http.StatusOK,
			expBody:   `{"message":"pong"}`,
		},
		"success - POST method": {
			givenReq: httptest.NewRequest(http.MethodPost, "/ping", bytes.NewBufferString(`{"message":"pong"}`)),
			hdl: handler{
				Method: http.MethodPost,
				Path:   "/ping",
				Func: func(c lightning.Context) error {
					var msg map[string]string
					if err := c.Bind(&msg); err != nil {
						return err
					}

					c.JSON(http.StatusOK, msg)
					return nil
				},
			},
			expStatus: http.StatusOK,
			expBody:   `{"message":"pong"}`,
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Mock func
			idFunc = func() string {
				return staticRequestID
			}
			defer func() {
				idFunc = uuid.NewString
			}()

			// Given
			w := httptest.NewRecorder()
			route, c, hdlRequest := lightning.NewRouterForTest(w)
			route.Use(RequestIDMiddleware())
			route.Handle(tc.hdl.Method, tc.hdl.Path, tc.hdl.Func)

			if slices.Contains([]string{http.MethodPost, http.MethodPut, http.MethodPatch}, tc.givenReq.Method) {
				tc.givenReq.Header.Set("Content-Type", "application/json")
			}

			c.SetRequest(tc.givenReq)

			// When
			hdlRequest()

			// Then
			require.Equal(t, tc.expStatus, w.Code)
			require.Equal(t, tc.expBody, w.Body.String())
			require.Equal(t, staticRequestID, w.Header().Get(RequestIDHeaderName))
		})
	}
}
