# Getting Started

## Project Overview

Lightning is a modular collection of Go libraries for building reliable backends quickly. It reduces boilerplate by providing reusable components for common tasks such as HTTP routing, gRPC service setup, internationalization, messaging, caching, logging, and configuration management.

## Features

- HTTP router with middleware and authentication guards
- gRPC server and client helpers
- Internationalization (i18n) utilities
- Kafka and Redis integrations
- Structured logging and monitoring tools
- Configuration management with environment overrides

## Repository Layout

- `router.go`, `server.go` – core HTTP router and server implementations
- `grpc_server.go`, `grpcclient/` – gRPC server and client helpers
- `middleware/` – reusable middleware components
- `broker/`, `caching/` – messaging and caching integrations (Kafka, Redis)
- `monitoring/` – metrics and tracing support
- `testutil/`, `mocks/` – testing utilities and mocks

## Installation

Ensure you have Go 1.20 or newer and a configured Go environment. Set up module mode if needed:

```bash
go env -w GO111MODULE=on
```

Install Lightning:

```bash
go get github.com/viebiz/lit@latest
```

## Hello World

```go
package main

import (
  "context"
  "net/http"

  "github.com/viebiz/lit"
)

func main() {
  r := lit.NewRouter(context.Background())
  r.Get("/", func(c lit.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
  })

  srv := lit.NewHttpServer(":8080", r.Handler())
  if err := srv.Run(); err != nil {
    panic(err)
  }
}
```

Run the example with `go run main.go` and visit `http://localhost:8080` to see the greeting.

