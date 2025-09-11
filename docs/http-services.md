# HTTP Services

This guide shows how to build HTTP services with `lit`. It covers the router, context wrapper, server lifecycle and functional options. Examples also demonstrate liveness/profiling endpoints and graceful shutdown.

## Router

Create a router with `NewRouter` and register handlers using the `Route`/`Group` APIs. Root middleware for logging, tracing and panic recovery is attached automatically and additional middleware can be added with `Use`.

```go
appCtx := context.Background()
r := lit.NewRouter(appCtx,
    lit.WithLivenessEndpoint("/alive"),
    lit.WithProfiling(),
)

r.Use(func(c lit.Context) error {
    // custom root middleware
    c.Next()
    return nil
})

r.Get("/ping", func(c lit.Context) error {
    return c.String(http.StatusOK, "pong")
})
```

## Context

Handlers receive a `lit.Context` which wraps `gin.Context` and adds helpers for request data, parameters and responses.

```go
func hello(c lit.Context) error {
    name := c.ParamWithDefault("name", "world")
    return c.JSON(http.StatusOK, map[string]string{"hello": name})
}
```

## Server and Graceful Shutdown

Use `NewHttpServer` to start an HTTP server. Functional options configure behaviour such as timeouts and shutdown grace period. `Run` listens for termination signals and shuts down gracefully.

```go
srv := lit.NewHttpServer(":8080", r.Handler(),
    lit.ServerShutdownGrace(10*time.Second),
)

if err := srv.Run(); err != nil {
    log.Fatal(err)
}
```

## Option Patterns

`lit` uses functional options. `RouterOption` modifies the router (e.g. `WithLivenessEndpoint`, `WithProfiling`). `ServerOption` customises the server (e.g. `ServerShutdownGrace`, `ServerReadTimeout`).

## Liveness and Profiling Endpoints

`WithLivenessEndpoint` exposes a plain text `OK` response which is ignored by monitoring. `WithProfiling` mounts Go's `net/http/pprof` handlers under `/_/profile`.

## Graceful Shutdown

The server uses a context to listen for `SIGINT`/`SIGTERM` and shuts down with the specified grace period, falling back to a forced close if needed.
