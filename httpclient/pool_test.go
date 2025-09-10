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

func TestOverridePoolMaxIdleConns(t *testing.T) {
	p := NewSharedPool(OverridePoolMaxIdleConns(50))
	tp := p.Transport.(*http.Transport).Clone()
	require.Equal(t, 50, tp.MaxIdleConns)
}

func TestOverridePoolMaxConnsPerHost(t *testing.T) {
	p := NewSharedPool(OverridePoolMaxConnsPerHost(20))
	tp := p.Transport.(*http.Transport).Clone()
	require.Equal(t, 20, tp.MaxConnsPerHost)
}

func TestOverridePoolMaxIdleConnsPerHost(t *testing.T) {
	p := NewSharedPool(OverridePoolMaxIdleConnsPerHost(15))
	tp := p.Transport.(*http.Transport).Clone()
	require.Equal(t, 15, tp.MaxIdleConnsPerHost)
}
