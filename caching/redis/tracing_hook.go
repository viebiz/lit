package redis

import (
	"context"
	"net"

	"github.com/redis/go-redis/extra/rediscmd/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/viebiz/lit/monitoring"
)

type tracingHook struct {
	info monitoring.ExternalServiceInfo
}

func newTracingHook(info monitoring.ExternalServiceInfo) redis.Hook {
	return tracingHook{
		info: info,
	}
}

func (t tracingHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		var err error
		span := trace.SpanFromContext(ctx)
		span.AddEvent("Dial")
		defer func() {
			if err != nil {
				recordError(span, err)
			}
		}()

		conn, err := next(ctx, network, addr)
		if err != nil {
			return nil, err
		}

		return conn, nil
	}
}

func (t tracingHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		var err error
		span := trace.SpanFromContext(ctx)
		span.AddEvent("Process", trace.WithAttributes(
			attribute.String("Name", cmd.FullName()),
			attribute.String("Statement", rediscmd.CmdString(cmd)),
		))
		defer func() {
			if err != nil {
				recordError(span, err)
			}
		}()

		if err := next(ctx, cmd); err != nil {
			return err
		}

		return nil
	}
}

func (t tracingHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		var err error
		summary, cmdsString := rediscmd.CmdsString(cmds)
		span := trace.SpanFromContext(ctx)
		span.AddEvent("ProcessPipeline", trace.WithAttributes(
			attribute.String("Summary", summary),
			attribute.String("Statement", cmdsString),
		))
		defer func() {
			if err != nil {
				recordError(span, err)
			}
		}()

		if err := next(ctx, cmds); err != nil {
			return err
		}

		return nil
	}
}

func recordError(span trace.Span, err error) {
	span.AddEvent("Redis error", trace.WithAttributes(
		attribute.String("Error", err.Error()),
	))
}
