# Monitoring

This project provides a monitoring layer that combines structured logging, distributed tracing, and error reporting.
The central type is `monitoring.Monitor`, which wraps a Zap logger and an optional Sentry client.
The monitor is stored in context and can be enriched with additional tags for better observability.

## Core components

### Logging and error reporting

`Monitor` wraps `zap.Logger` for structured logs and integrates Sentry for error reporting.  Use `Info` or `Infof` for informational logs.
`Error` and `Errorf` capture errors, attach stack traces, and forward them to Sentry.
```go
m.Info("starting up")
if err != nil {
    m.Errorf(err, "failed to start")
}
```

### Tracing

Tracing is built on OpenTelemetry.  A global tracer records spans and propagates trace information through contexts.
Helper functions like `InjectField` and `StartSegment` attach attributes to spans and ensure trace and span IDs are logged.

## Instrumentation packages

### HTTP

The `instrumenthttp` package instruments both servers and clients.

Server handlers wrap each incoming request with a span and collect metadata:
```go
func handler(m *monitoring.Monitor, w http.ResponseWriter, r *http.Request) {
    ctx, reqMeta, end := instrumenthttp.StartIncomingRequest(m, r, "/users/:id")
    defer end(http.StatusOK, nil)

    m = monitoring.FromContext(ctx)
    m.Infof("request", monitoring.String("method", reqMeta.Method))
}
```

Clients group retries with `StartOutgoingGroupSegment` and instrument individual requests with `StartOutgoingSegment`:
```go
ctx, endGroup := instrumenthttp.StartOutgoingGroupSegment(ctx, svc, "user-service", http.MethodGet, url)
req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
ctx, endReq := instrumenthttp.StartOutgoingSegment(ctx, req)
resp, err := client.Do(req)
endReq(resp.StatusCode, err)
endGroup(err)
```

### gRPC

`instrumentgrpc` instruments unary RPCs.

Server-side usage:
```go
ctx, meta, end := instrumentgrpc.StartUnaryIncomingCall(ctx, m, fullMethod, req)
defer end(nil)
```

Client-side usage:
```go
ctx, end := instrumentgrpc.StartUnaryCallSegment(ctx, svcInfo, fullMethod)
resp, err := client.SomeCall(ctx, req)
end(err)
```

### Kafka

`instrumentkafka` traces Kafka producers.

```go
ctx, end := instrumentkafka.StartSyncPublishSegment(ctx, svc, clientID, msg, key)
partition, offset, err := producer.SendMessage(msg)
end(partition, offset, err)
```

For async producers, `StartAsyncEnqueueSegment` instruments the enqueue operation and returns a `PublishSegment` whose `End` method records partition and offset once the broker acknowledges the message.

### Postgres

Wrap database pools and transactions with `instrumentpg.WithInstrumentation` and `instrumentpg.WithInstrumentationTx` to record query events and errors:
```go
pool := instrumentpg.WithInstrumentation(pgpool)
ctx, end := monitoring.StartSegment(ctx, "insert-user")
defer end()
_, err := pool.ExecContext(ctx, "INSERT INTO users(name) VALUES($1)", name)
```

## Patterns

1. **Tracing** – Start spans at service boundaries using instrumentation helpers.  Spans are propagated through context and annotated with attributes.
2. **Logging** – Retrieve the monitor from context with `monitoring.FromContext` and log messages with contextual tags.
3. **Error reporting** – Use `Error`/`Errorf` so errors are both logged and sent to Sentry.

These patterns ensure consistent observability across HTTP, gRPC, Kafka, and Postgres interactions.

