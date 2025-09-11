# Redis Caching

This guide covers the `caching/redis` package, which provides a thin wrapper around the [go-redis v9](https://github.com/redis/go-redis) client with additional tracing and convenience helpers.

## Client Setup

Create a client using the `NewClient` helper. The function accepts a Redis connection URL and automatically registers tracing hooks:

```go
import (
    "context"
    redis "github.com/viebiz/lit/caching/redis"
)

func main() {
    ctx := context.Background()

    // redis://user:password@localhost:6379/0?protocol=3
    client, err := redis.NewClient("redis://localhost:6379/0")
    if err != nil {
        panic(err)
    }
    defer client.Close()

    if err := client.Ping(ctx); err != nil {
        panic(err)
    }
}
```

`NewClientWithTLS` allows passing a `*tls.Config` when TLS is required.

## Common Cache Operations

The client exposes helpers for working with primitive values and hashes:

```go
// Strings and numbers
author := "alice"
_ = client.SetString(ctx, "user:1:name", author, time.Hour)
name, _ := client.GetString(ctx, "user:1:name")

_ = client.SetInt(ctx, "counter", 1, 0)
next, _ := client.IncrementBy(ctx, "counter", 1)

// Hashes
type Profile struct {
    Name string `redis:"name"`
    Age  int    `redis:"age"`
}

_ = client.HashSet(ctx, "user:1", Profile{Name: "alice", Age: 30})
var p Profile
_ = client.HashGetAll(ctx, "user:1", &p)
```

Always pass a `context.Context` with appropriate timeouts or cancellation and handle returned errors.

## Pipelines

Batch multiple commands with `DoInBatch`, which uses Redis pipelines under the hood:

```go
err := client.DoInBatch(ctx, func(cm redis.Commander) error {
    _ = cm.SetString(ctx, "user:1", "alice", time.Hour)
    _ = cm.IncrementBy(ctx, "page:views", 1)
    // Results become available after the function returns.
    return nil
})
```

Use pipelines to minimise round trips when issuing many independent commands. The `Commander` passed to the callback supports most standard operations.

## Pub/Sub

Publish messages with `Publish` and create subscribers using `Subscribe`:

```go
handler := func(ctx context.Context, msg redis.Message) error {
    fmt.Println("got", msg.Channel, msg.Payload)
    return nil
}

sub := client.Subscribe(ctx, []string{"events"}, handler)
go func() {
    if err := sub.Subscribe(ctx); err != nil {
        log.Println("subscribe error", err)
    }
}()

_ = client.Publish(ctx, "events", "hello world")
```

`SubscribeWithOptions` accepts a `ChannelOption` for tuning channel size, health check interval, and send timeout.

## Tracing Hooks & Usage Patterns

`NewClient` automatically attaches OpenTelemetry tracing hooks. Each command, pipeline, and connection attempt is recorded as a span event, and failures are captured as errors. Ensure that the context passed to Redis operations carries the active trace span so the instrumentation can attach data to it.

Recommended practices:

- Reuse a single `Client` instance throughout your application and close it on shutdown.
- Use pipelines for large batches of independent commands.
- Prefer structured contexts and handle cancellation.
- For pub/sub consumers, run `Subscribe` in its own goroutine and respect context cancellation to exit cleanly.
