package httpclient

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewSharedPool(t *testing.T) {
	p := NewSharedPool()
	require.Equal(t, defaultTimeoutPerTry, p.Timeout)
	tp := p.Transport.(*http.Transport).Clone()
	require.Equal(t, defaultMaxIdleConnsPerHost, tp.MaxIdleConnsPerHost)
	require.False(t, tp.TLSClientConfig.InsecureSkipVerify)
}

func TestNewSharedCustomPool(t *testing.T) {
	p := NewSharedCustomPool()
	require.Equal(t, time.Duration(0), p.Timeout)
	tp := p.Transport.(*http.Transport).Clone()
	require.Equal(t, defaultMaxIdleConnsPerHost, tp.MaxIdleConnsPerHost)
	require.False(t, tp.TLSClientConfig.InsecureSkipVerify)
}
