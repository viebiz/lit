package instrumenthttp

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/viebiz/lit/monitoring"
)

func TestStartOutgoingGroupSegment_NoCurSpan(t *testing.T) {
	_, end := StartOutgoingGroupSegment(
		context.Background(),
		monitoring.ExternalServiceInfo{},
		"svc",
		http.MethodGet,
		"/url",
	)

	end(errors.New("some err"))
}

func TestStartOutgoingGroupSegment_WithCurSpan(t *testing.T) {
	ctx, _ := tracer.Start(context.Background(), "test")

	_, end := StartOutgoingSegment(
		ctx,
		monitoring.ExternalServiceInfo{},
		"svc",
		httptest.NewRequest(http.MethodGet, "/url", nil),
	)

	end(200, errors.New("some err"))
}

func TestStartOutgoingSegment_NoCurSpan(t *testing.T) {
	_, end := StartOutgoingSegment(
		context.Background(),
		monitoring.ExternalServiceInfo{},
		"svc",
		httptest.NewRequest(http.MethodGet, "/url", nil),
	)

	end(200, errors.New("some err"))
}

func TestStartOutgoingSegment_WithCurSpan(t *testing.T) {
	ctx, _ := tracer.Start(context.Background(), "test")

	_, end := StartOutgoingSegment(
		ctx,
		monitoring.ExternalServiceInfo{},
		"svc",
		httptest.NewRequest(http.MethodGet, "/url", nil),
	)

	end(200, errors.New("some err"))
}
