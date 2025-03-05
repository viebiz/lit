package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

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
