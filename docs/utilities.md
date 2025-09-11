# Utilities

This project provides several helper packages that simplify common tasks such as
ID generation and loading resources.  These utilities are designed to be small
and composable so they can be used across services.

## Snowflake ID generation

The [`snowflake/`](../snowflake/) package wraps the [sonyflake](https://github.com/sony/sonyflake)
implementation to produce unique 64‑bit identifiers.  It exposes a `Generator`
that can be configured with options such as a custom epoch (`StartTime`) or a
specific machine identifier (`MachineID`).

```go
import "github.com/viebiz/lit/snowflake"

// Create a generator with default settings
idGen, err := snowflake.New()
if err != nil {
    // handle error
}

// Generate a new ID
id, err := idGen.Generate()
if err != nil {
    // handle error
}
```

## Resource loading helpers

The [`ioutil/`](../ioutil/) package offers a small abstraction for loading files
from a configured resources directory.  `SetResourceDir` is safe to call multiple
times but will only set the directory once, and `ReadFile` resolves paths relative
to that directory.

```go
import "github.com/viebiz/lit/ioutil"

// Configure the resource directory once at startup
ioutil.SetResourceDir("./resources")

// Read a file relative to that directory
data, err := ioutil.ReadFile("config.yaml")
if err != nil {
    // handle error
}
```

## Other helpers

Several additional utility modules appear throughout the codebase:

- `errors.go` defines an `HTTPError` type that standardises HTTP error responses.
- `context.go` extends `gin.Context` with convenience methods for parameter and
  header handling while keeping the standard `context.Context` interface.
- `httpclient` and `grpcclient` provide reusable client factories with options
  for pooling, retries and instrumentation.
- The `guard` package hosts middleware and helpers for authentication and
  permission checks.

These cross‑cutting utilities are leveraged by routers, middleware and tests to
keep service code focused on business logic.

