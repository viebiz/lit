# gRPC Support

This project provides helpers for building instrumented gRPC servers and clients.

## Server

`GRPCServer` wraps `grpc.Server` and applies sensible defaults, including the built‑in unary interceptor.

```go
srv, err := lit.NewGRPCServer(ctx, ":50051")
if err != nil { /* handle */ }
```

Additional `GRPCOption`s may be supplied when constructing a server. For example, TLS can be enabled:

```go
tlsCfg := &tls.Config{/* ... */}
srv, err := lit.NewGRPCServerWithOptions(ctx, ":50051", lit.WithTLSConfig(tlsCfg))
```

### Interceptors

`WithDefaultInterceptors` adds a unary interceptor that instruments calls, logs request/response bodies, and recovers from panics.

```go
lit.WithDefaultInterceptors(ctx)
```

The interceptor starts a tracing span, logs the incoming request and response, and converts panics into `codes.Internal` errors.

### Service Registration

`GRPCServer` exposes a `Registrar` method that returns a `ServiceRegistrar` interface so that services can be registered with generated code:

```go
pb.RegisterGreeterServer(srv.Registrar(), greeterImpl)
```

## Client

The `grpcclient` package supplies helpers for dialing services. `NewUnauthenticatedConnection` creates a `grpc.ClientConn` wrapped by the `Conn` interface. Calls are instrumented and the request payload is logged by a unary client interceptor.

```go
conn, err := grpcclient.NewUnauthenticatedConnection(ctx, "localhost:50051")
if err != nil { /* handle */ }
client := pb.NewGreeterClient(conn)
resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "lit"})
```

## Example Server and Client

```go
package main

import (
    "context"
    "log"
    "crypto/tls"

    "github.com/viebiz/lit"
    "github.com/viebiz/lit/grpcclient"
    pb "path/to/your/service"
)

func main() {
    ctx := context.Background()

    tlsCfg := &tls.Config{ /* ... */ }
    srv, err := lit.NewGRPCServerWithOptions(ctx, ":50051", lit.WithTLSConfig(tlsCfg))
    if err != nil { log.Fatal(err) }
    pb.RegisterGreeterServer(srv.Registrar(), greeter{})
    go func() {
        if err := srv.Run(); err != nil { log.Fatal(err) }
    }()

    conn, err := grpcclient.NewUnauthenticatedConnection(ctx, ":50051")
    if err != nil { log.Fatal(err) }
    client := pb.NewGreeterClient(conn)
    _, _ = client.SayHello(ctx, &pb.HelloRequest{Name: "world"})
}
```
