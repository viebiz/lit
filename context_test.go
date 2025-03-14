package lit

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAbortWithError(t *testing.T) {
	type mockWriter struct {
		useMock  bool
		inStatus int
		outErr   error
	}

	tests := map[string]struct {
		inErr          error
		mockWriter     mockWriter
		expectedStatus int
		expectedBody   string
	}{
		"common error": {
			inErr:          errors.New("simulated error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: func() string {
				b, err := json.Marshal(ErrDefaultInternal)
				require.NoError(t, err)
				return string(b)
			}(),
		},
		"expected error": {
			inErr: testError{
				Code: http.StatusBadRequest,
				Msg:  "bad request",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func() string {
				b, err := json.Marshal(testError{Code: http.StatusBadRequest, Msg: "bad request"})
				require.NoError(t, err)
				return string(b)
			}(),
		},
		"error when marshal": {
			inErr:          testErrorMarshal{},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: func() string {
				b, err := json.Marshal(ErrDefaultInternal)
				require.NoError(t, err)
				return string(b)
			}(),
		},
		"error when write": {
			inErr: testError{
				Code: http.StatusBadRequest,
				Msg:  "bad request",
			},
			mockWriter: mockWriter{
				useMock:  true,
				inStatus: http.StatusBadRequest,
				outErr:   errors.New("simulated error"),
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			recorder := httptest.NewRecorder()
			c := CreateTestContext(recorder)

			if tc.mockWriter.useMock {
				mockResponseWriter := NewMockResponseWriter(t)
				mockResponseWriter.On("Header").Return(recorder.Header())
				mockResponseWriter.On("WriteHeader", tc.mockWriter.inStatus).Run(func(args mock.Arguments) {
					recorder.WriteHeader(tc.mockWriter.inStatus)
				})
				mockResponseWriter.EXPECT().WriteHeaderNow()
				mockResponseWriter.EXPECT().Write(mock.AnythingOfType("[]uint8")).Return(0, tc.mockWriter.outErr)

				// Override the response writer
				c.SetWriter(mockResponseWriter)
			}

			// When
			c.AbortWithError(tc.inErr)

			// Then
			require.Equal(t, tc.expectedStatus, recorder.Code)
			require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
			require.Equal(t, tc.expectedBody, recorder.Body.String())
		})
	}
}

type testError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e testError) Error() string {
	return e.Msg
}

func (e testError) StatusCode() int {
	return e.Code
}

type testErrorMarshal struct{}

func (e testErrorMarshal) Error() string {
	return "marshal error"
}

func (e testErrorMarshal) StatusCode() int {
	return http.StatusBadRequest
}

func (e testErrorMarshal) MarshalJSON() ([]byte, error) {
	return nil, errors.New("marshal error")
}
