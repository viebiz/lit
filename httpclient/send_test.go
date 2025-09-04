package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_Send(t *testing.T) {
	// Given:
	ctx := context.Background()

	reqBody := []byte(`{"id": 1, "name": "loc.dang", "legion": "darkangels"}`)
	query := url.Values{
		"key1": []string{"value1"},
		"key2": []string{"value2"},
	}
	vars := map[string]string{
		"var1": "one",
	}
	contentType := "application/json"
	headers := map[string]string{
		"X-Default": "123456789",
	}

	// Mock:
	var called bool
	mockSvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert request reqBody
		rbody, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, string(reqBody), string(rbody))
		// Assert path variable substitution
		require.False(t, strings.Contains(r.URL.Path, ":var1"))
		require.True(t, strings.HasSuffix(r.URL.Path, "/one"))
		// Assert query params
		require.Equal(t, query.Encode(), r.URL.Query().Encode())
		// Assert request headers
		require.Equal(t, contentType, r.Header.Get("Content-Type"))
		require.Equal(t, headers["X-Default"], r.Header.Get("X-Default"))

		called = true

		w.Header().Set("key", "value")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "expected response"}`))
	}))
	sURL := mockSvr.URL + "/:var1"

	c, err := NewUnauthenticated(
		Config{URL: sURL, Method: http.MethodPost, ServiceName: "svc"},
		NewSharedCustomPool(),
		OverrideBaseRequestHeaders(headers),
		OverrideContentType(contentType),
	)
	require.NoError(t, err)

	// When:
	resp, err := c.Send(ctx, Payload{
		Body:        reqBody,
		QueryParams: query,
		PathVars:    vars,
		Header:      headers,
	})

	// Then:
	require.NoError(t, err)
	require.True(t, called)
	require.Equal(t, http.StatusAccepted, resp.Status)
	require.NoError(t, err)
	require.Equal(t, []byte(`{"message": "expected response"}`), resp.Body)
	require.Equal(t, "value", resp.Header.Get("key"))
}

func BenchmarkClient_Send(b *testing.B) {
	// Given:
	ctx := context.Background()

	body := []byte(`{"id": 1, "name": "loc.dang", "legion": "darkangels"}`)
	query := url.Values{
		"key1": []string{"value1"},
		"key2": []string{"value2"},
	}
	vars := map[string]string{
		"var1": "one",
	}
	contentType := "application/json"
	headers := map[string]string{
		"X-Default": "123456789",
	}

	// Mock:
	mockSvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("key", "value")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "expected response"}`))
	}))
	defer mockSvr.Close()

	sURL := mockSvr.URL + "/:var1"

	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		c, err := NewUnauthenticated(
			Config{URL: sURL, Method: http.MethodPost, ServiceName: "svc"},
			NewSharedCustomPool(),
			OverrideBaseRequestHeaders(headers),
			OverrideContentType(contentType),
		)
		require.NoError(b, err)

		for pb.Next() {
			_, _ = c.Send(ctx, Payload{
				Body:        body,
				QueryParams: query,
				PathVars:    vars,
				Header:      headers,
			})
		}
	})
}

func TestClient_constructURL(t *testing.T) {
	c := &Client{url: "http://example.com/:id"}

	// Test path vars
	p := Payload{PathVars: map[string]string{"id": "123"}}
	resultURL := c.constructURL(p)
	require.Equal(t, "http://example.com/123", resultURL)

	// Test query params
	query := url.Values{"key": []string{"value"}}
	p = Payload{QueryParams: query}
	resultURL = c.constructURL(p)
	require.Equal(t, "http://example.com/:id?key=value", resultURL)

	// Test both
	p = Payload{
		PathVars:    map[string]string{"id": "123"},
		QueryParams: url.Values{"key": []string{"value"}},
	}
	resultURL = c.constructURL(p)
	require.Equal(t, "http://example.com/123?key=value", resultURL)

	// Test existing query
	c.url = "http://example.com/:id?existing=1"
	p = Payload{
		PathVars:    map[string]string{"id": "123"},
		QueryParams: url.Values{"key": []string{"value"}},
	}
	resultURL = c.constructURL(p)
	require.Equal(t, "http://example.com/123?existing=1&key=value", resultURL)
}

func TestClient_setHeader(t *testing.T) {
	c := &Client{
		userAgent:   "test-agent",
		contentType: "application/json",
		header:      header{values: map[string]string{"X-Default": "default"}},
	}
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	p := Payload{Header: map[string]string{"X-Custom": "custom", "X-Default": "override"}}

	c.setHeader(req, p)

	require.Equal(t, "test-agent", req.Header.Get("User-Agent"))
	require.Equal(t, "application/json", req.Header.Get("Content-Type"))
	require.Equal(t, "override", req.Header.Get("X-Default"))
	require.Equal(t, "custom", req.Header.Get("X-Custom"))
}

func TestClient_Send_Timeout(t *testing.T) {
	ctx := context.Background()
	mockSvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Simulate delay
		w.WriteHeader(http.StatusOK)
	}))
	defer mockSvr.Close()

	c, err := NewUnauthenticated(
		Config{URL: mockSvr.URL, Method: http.MethodGet, ServiceName: "svc"},
		NewSharedCustomPool(),
		OverrideTimeoutAndRetryOption(1, 50*time.Millisecond, 100*time.Millisecond, true, nil),
	)
	require.NoError(t, err)

	_, err = c.Send(ctx, Payload{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "timeout")
}

func TestClient_Send_RetryOnStatusCode(t *testing.T) {
	callCount := 0
	mockSvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer mockSvr.Close()

	c, err := NewUnauthenticated(
		Config{URL: mockSvr.URL, Method: http.MethodGet, ServiceName: "svc"},
		NewSharedCustomPool(),
		OverrideTimeoutAndRetryOption(2, 100*time.Millisecond, 10*time.Second, false, []int{500}),
	)
	require.NoError(t, err)

	resp, err := c.Send(context.Background(), Payload{})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.Status)
	require.Equal(t, 2, callCount)
}

func TestClient_Send_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockSvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cancel() // Cancel context
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockSvr.Close()

	c, err := NewUnauthenticated(
		Config{URL: mockSvr.URL, Method: http.MethodGet, ServiceName: "svc"},
		NewSharedCustomPool(),
	)
	require.NoError(t, err)

	_, err = c.Send(ctx, Payload{})
	require.Error(t, err)
	require.Equal(t, ErrOperationContextCanceled, err)
}

func TestReadRespBody_Error(t *testing.T) {
	// Mock a reader that returns error
	reader := &errorReader{err: io.ErrUnexpectedEOF}
	_, err := readRespBody(reader)
	require.Error(t, err)
	require.Equal(t, io.ErrUnexpectedEOF, err)
}

type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}
