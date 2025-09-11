# HTTP Middleware

This document describes the HTTP middleware and CORS utilities provided by **lit**.

## Request ID

`RequestIDMiddleware` ensures every request has a unique identifier. It accepts an incoming `X-Request-ID` header or generates one. The ID is injected into the request context for monitoring and returned in the response header.

## Localization

`LocalizationMiddleware` loads translation bundles and injects a localized message provider into the request context. It reads the language from the `Accept-Language` header, falling back to `en`, and sets the chosen language in the `Content-Language` response header.

## Logging

By default, lit logs request and response data. Use `SkipLoggingResponseBodyMiddleware` to prevent the response body from being written to the logs for sensitive endpoints.

## CORS Configuration

The [`cors`](../cors) package exposes a configurable middleware for Cross-Origin Resource Sharing.

```go
import (
    "github.com/viebiz/lit/cors"
)

cfg := cors.New([]string{"https://example.com"})
cfg.SetAllowMethods("GET", "POST")
// cfg.SetAllowHeaders(...)
// cfg.DisableCredentials()
// cfg.SetMaxAge(10 * time.Minute)

r.Use(cors.Middleware(cfg))
```

## Middleware Chaining

Middlewares are executed in the order provided. They can be applied globally or per handler.

```go
r := lit.NewRouter(ctx)

r.Use(http.RequestIDMiddleware())

r.Get("/hello", helloHandler,
    http.LocalizationMiddleware(ctx, http.Config{}),
    http.SkipLoggingResponseBodyMiddleware(),
)
```

The example above assigns a request ID to all requests and applies localization and logging control only to the `/hello` route.
