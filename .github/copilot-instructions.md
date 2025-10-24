# GitHub Copilot Instructions for Lightning (Lit)

This document provides guidance for AI coding assistants (GitHub Copilot, Claude, Cursor, etc.) working with the **Lightning (Lit)** project - a modular, production-ready collection of Go libraries for building reliable backend services.

---

## Project Overview

**Lightning (Lit)** is a comprehensive Go framework that reduces boilerplate by providing reusable components for:

- HTTP routing with authentication guards and middleware
- gRPC server and client setup with service registration
- Internationalization (i18n) support
- Kafka and Redis integrations
- Structured, context-aware logging with distributed tracing
- Configuration management with environment overrides
- Authorization and authentication (JWT, Casbin)
- Testing utilities and mocks

**Repository:** `github.com/viebiz/lit`  
**Go Version:** 1.24.6  
**License:** Apache-2.0

---

## Core Architecture Principles

### 1. Functional Options Pattern

The project extensively uses the **functional options pattern** for configuration:

```go
// Creating a router with options
r := lit.NewRouter(ctx,
    lit.WithLivenessEndpoint("/alive"),
    lit.WithProfiling(),
)

// Creating a server with options
srv := lit.NewHttpServer(":8080", r.Handler(),
    lit.ServerShutdownGrace(10*time.Second),
    lit.ServerReadTimeout(time.Minute),
)
```

**When adding new components:**
- Define an options type (e.g., `RouterOption`, `ServerOption`)
- Create option functions that return the option type
- Apply options in the constructor

### 2. Context Propagation

All components are designed with **context-first** principles:

- Application context flows through all layers
- Request contexts carry monitoring, tracing, and user information
- Use `monitoring.FromContext(ctx)` to retrieve the monitor
- Use `iam.UserProfileFromContext(ctx)` for user information
- Never ignore contexts - always propagate them

### 3. Middleware and Handler Pattern

HTTP handlers follow a consistent pattern:

```go
type HandlerFunc func(Context) error

func MyHandler(c lit.Context) error {
    // Access monitor from context
    m := monitoring.FromContext(c.Request().Context())
    
    // Bind request data
    var req MyRequest
    if err := c.Bind(&req); err != nil {
        return err
    }
    
    // Business logic here
    
    // Return response
    return c.JSON(http.StatusOK, response)
}
```

### 4. Instrumentation Everywhere

All external interactions should be instrumented:

- **HTTP requests:** Use `instrumenthttp.StartOutgoingSegment`
- **gRPC calls:** Use `instrumentgrpc.StartUnaryCallSegment`
- **Kafka messages:** Use `instrumentkafka.StartSyncPublishSegment`
- **Database queries:** Use `instrumentpg` middleware
- Always call the `end()` function with defer

---

## Project Structure

```
lit/
├── *.go                    # Core HTTP router, server, context interfaces
├── broker/kafka/           # Kafka producer/consumer with instrumentation
├── caching/redis/          # Redis client wrapper
├── cors/                   # CORS middleware
├── docs/                   # Documentation (getting-started, auth, monitoring, etc.)
├── env/                    # Configuration management (viper-based)
├── grpcclient/             # gRPC client utilities
├── guard/                  # Authentication/authorization middleware
├── httpclient/             # Instrumented HTTP client
├── i18n/                   # Internationalization support
├── iam/                    # Identity & Access Management (JWT, Casbin)
├── ioutil/                 # I/O utilities
├── jwt/                    # JWT token generation and parsing
├── middleware/http/        # HTTP middleware (logging, recovery, etc.)
├── mocks/                  # Generated mocks (mockery)
├── monitoring/             # Logging, tracing, error reporting (Zap, OpenTelemetry, Sentry)
│   ├── instrumentgrpc/
│   ├── instrumenthttp/
│   ├── instrumentkafka/
│   └── instrumentpg/
├── postgres/               # PostgreSQL utilities
├── snowflake/              # Distributed ID generation
├── testutil/               # Testing helpers
└── testdata/               # Test fixtures
```

---

## Code Style Guidelines

### Formatting

- **Always run `go fmt ./...` before committing**
- Use **tabs for indentation** (configured in `.editorconfig`)
- Maximum line length: **120 characters**
- Insert final newline in all files

### Naming Conventions

- **Interfaces:** Descriptive names without "I" prefix (e.g., `Router`, `Context`, `Enforcer`)
- **Implementations:** Lowercase private structs (e.g., `router`, `context`)
- **Constructors:** Always use `New` prefix (e.g., `NewRouter`, `NewHttpServer`)
- **Options:** Use `With` prefix for option functions (e.g., `WithLivenessEndpoint`)
- **Test helpers:** Suffix with `ForTest` (e.g., `NewRouterForTest`)

### Package Organization

- Keep packages focused and cohesive
- Use `apis.go` for public APIs
- Use `spec.go` for interfaces/types
- Use `type.go` for type definitions
- Use `errors.go` for error definitions
- Group related functionality (e.g., `producer_async.go`, `producer_option.go`)

### Error Handling

- Use `github.com/pkg/errors` for error wrapping with context
- Always wrap errors with meaningful context:
  ```go
  return pkgerrors.Wrap(err, "failed to connect to database")
  ```
- Define sentinel errors as package-level variables:
  ```go
  var ErrEmptyTopic = errors.New("topic is empty")
  ```

---

## Testing Requirements

### Coverage Requirements

- **All logic MUST be covered by unit tests**
- CI/CD enforces test passing - **PRs with failing tests will not be merged**
- Use **table-driven tests** where applicable
- Mock external dependencies to keep tests isolated

### Test Structure

```go
func TestMyFunction(t *testing.T) {
    type arg struct {
        givenInput string
        expOutput  string
        expErr     error
    }
    tcs := map[string]arg{
        "success case": {
            givenInput: "test",
            expOutput:  "expected",
        },
        "error case": {
            givenInput: "",
            expErr:     errors.New("invalid input"),
        },
    }
    
    for scenario, tc := range tcs {
        tc := tc
        t.Run(scenario, func(t *testing.T) {
            t.Parallel()
            
            // Given
            // When
            got, err := MyFunction(tc.givenInput)
            
            // Then
            if tc.expErr != nil {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            testutil.Equal(t, tc.expOutput, got)
        })
    }
}
```

### Testing Tools & Conventions

**Tools:**
- **Assertions:** Use `github.com/stretchr/testify/require` (not `assert`)
- **Comparisons:** Use `testutil.Equal()` for deep comparisons with diffs
- **Mocking:** Use `mockery` for generating mocks (config in `.mockery.yaml`)
- **HTTP Testing:** Use `NewRouterForTest()` helper

**Test Pattern Conventions:**

1. **Use Map-Based Table-Driven Tests** (NOT slice-based)
   ```go
   // ✅ CORRECT: Use map[string]arg
   tcs := map[string]arg{
       "scenario name": { /* test case */ },
   }
   
   // ❌ WRONG: Don't use []struct
   tests := []struct {
       name string
       // ...
   }{}
   ```

2. **Naming Convention for Test Data:**
   - Use `arg` struct for test case definition
   - Prefix inputs with `given`: `givenInput`, `givenContext`, `givenOptions`
   - Prefix expected outputs with `exp`: `expOutput`, `expErr`, `expResult`

3. **Test Structure Pattern:**
   ```go
   for scenario, tc := range tcs {
       tc := tc  // Capture range variable for parallel tests
       t.Run(scenario, func(t *testing.T) {
           t.Parallel()  // Enable parallel execution
           
           // Given - Setup mocks and test data
           
           // When - Execute the function under test
           
           // Then - Assert results
       })
   }
   ```

4. **Always Include:**
   - `tc := tc` to capture range variable before parallel execution
   - `t.Parallel()` for independent test cases
   - `Given/When/Then` comments for test organization

5. **Complex Test Cases:**
   For tests with multiple mock setups:
   ```go
   type mockArg struct {
       givenContext context.Context
       expResult    *Result
   }
   type arg struct {
       givenMockSetup func() mockArg
       givenInput     string
       expErr         error
   }
   ```



### Running Tests

```bash
# Run all tests
make test

# Run specific package tests
go test -v ./guard/...

# Run with coverage
go test -coverprofile=coverage.out ./...

# Run benchmarks
make benchmark
```

---

## Common Patterns

### Creating HTTP Routers

```go
func NewServer(appCtx context.Context) *lit.Server {
    r := lit.NewRouter(appCtx,
        lit.WithLivenessEndpoint("/alive"),
        lit.WithProfiling(),
    )
    
    // Root middleware (authentication, logging, etc.)
    r.Use(customMiddleware())
    
    // Route groups
    r.Group("/api/v1", func(v1 lit.Router) {
        v1.Get("/users", listUsers)
        v1.Post("/users", createUser)
    }, authMiddleware())
    
    return lit.NewHttpServer(":8080", r.Handler(),
        lit.ServerShutdownGrace(10*time.Second),
    )
}
```

### Adding Middleware

```go
func MyMiddleware() lit.HandlerFunc {
    return func(c lit.Context) error {
        // Pre-processing
        ctx := c.Request().Context()
        m := monitoring.FromContext(ctx)
        m.Info("middleware executed")
        
        // Call next handler
        if err := c.Next(); err != nil {
            return err
        }
        
        // Post-processing
        return nil
    }
}
```

### Working with Monitoring

```go
func MyHandler(c lit.Context) error {
    ctx := c.Request().Context()
    m := monitoring.FromContext(ctx)
    
    // Add context to logs
    m = m.WithField(monitoring.String("user_id", "123"))
    
    // Log info
    m.Info("processing request")
    
    // Log errors (sends to Sentry automatically)
    if err != nil {
        m.Errorf(err, "failed to process request")
        return err
    }
    
    return nil
}
```

### Authentication & Authorization

```go
// Setup guard
guard := guard.New(validator, enforcer)

// User authentication
r.Group("/api", func(api lit.Router) {
    // User routes
    api.Get("/profile",
        guard.AuthenticateUserMiddleware(),
        getProfile,
    )
    
    // M2M routes with scope check
    api.Post("/admin",
        guard.AuthenticateM2MMiddleware(),
        guard.RequiredM2MScopeMiddleware("write:admin"),
        adminAction,
    )
    
    // Role-based permission check
    api.Post("/resources",
        guard.AuthenticateUserMiddleware(),
        guard.RolePermissionHandler(enforcer, "resources", "write", createResource),
    )
})
```

### Kafka Producer

```go
// Create producer
producer, err := kafka.NewAsyncProducer(cfg,
    kafka.WithSuccessHandler(func(msg *sarama.ProducerMessage) {
        log.Println("message sent successfully")
    }),
    kafka.WithErrorHandler(func(err *sarama.ProducerError) {
        log.Printf("failed to send: %v", err)
    }),
)

// Send message
ctx, end := instrumentkafka.StartAsyncEnqueueSegment(ctx, svcInfo, clientID, msg, key)
publisher := producer.SendMessage(msg)
end(publisher)
```

### Kafka Consumer

```go
consumer, err := kafka.NewConsumer(cfg,
    kafka.WithConsumerGroup("my-group"),
    kafka.WithTopics("my-topic"),
)

handler := func(ctx context.Context, msg *sarama.ConsumerMessage) error {
    m := monitoring.FromContext(ctx)
    m.Infof("received message", monitoring.String("key", string(msg.Key)))
    // Process message
    return nil
}

if err := consumer.Consume(ctx, handler); err != nil {
    log.Fatal(err)
}
```

### Configuration Management

```go
type AppConfig struct {
    ServerAddr string `mapstructure:"server_addr"`
    DBUrl      string `mapstructure:"db_url"`
    LogLevel   string `mapstructure:"log_level"`
}

func (c AppConfig) Validate() error {
    if c.ServerAddr == "" {
        return errors.New("server_addr is required")
    }
    return nil
}

// Read config from .env file
cfg, err := env.ReadAppConfig[AppConfig]()
```

---

## Dependencies

### Core Dependencies

- **HTTP Framework:** `github.com/gin-gonic/gin` (wrapped by lit.Context)
- **gRPC:** `google.golang.org/grpc`
- **Messaging:** `github.com/IBM/sarama` (Kafka)
- **Caching:** `github.com/redis/go-redis/v9`
- **Database:** `gorm.io/gorm`, `github.com/jackc/pgx/v5`
- **Logging:** `go.uber.org/zap`
- **Tracing:** `go.opentelemetry.io/otel`
- **Error Reporting:** `github.com/getsentry/sentry-go`
- **Authorization:** `github.com/casbin/casbin/v2`
- **Config:** `github.com/spf13/viper`
- **Testing:** `github.com/stretchr/testify`
- **Mocking:** Generated with `mockery`

### When Adding Dependencies

1. Check if the dependency already exists in `go.mod`
2. Prefer well-maintained, widely-used libraries
3. Consider security implications (see `.aikidoignore` for security scanning)
4. Update `go.mod` and `go.mod vendor` after adding

---

## Development Workflow

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/viebiz/lit.git
cd lit

# Start dependencies (Postgres, Redis, Kafka)
make init

# Run tests
make test

# Generate mocks
make gen-mocks

# Generate protobuf code
make gen-proto
```

### Docker Services

The project uses Docker Compose for local development:

- **Postgres:** localhost:54321 (user: `lit`, db: `master`)
- **Redis:** localhost:63791
- **Kafka:** localhost:9092
- **Jaeger:** localhost:16686 (UI), localhost:4317 (OTLP)

### Creating a Pull Request

1. **Create a feature branch:** `git checkout -b feature/my-feature`
2. **Make focused commits** with clear messages
3. **Add tests** for all new logic
4. **Run `make test`** and ensure all tests pass
5. **Format code:** `go fmt ./...`
6. **Update documentation** if adding new features
7. **Push and open a PR** with a clear description

### CI/CD

- **All tests must pass** before merge approval
- CI runs on every PR and checks:
  - Code formatting/linting
  - Unit test execution
  - Build validation
  - Coverage reporting

---

## Mock Generation

Mocks are generated using [mockery](https://github.com/vektra/mockery) with configuration in `.mockery.yaml`:

```bash
# Generate all mocks
make gen-mocks

# Or manually
mockery
```

**Mocking conventions:**
- Mocks are generated in the same package as the interface
- Naming: `Mock{{.InterfaceName}}` (e.g., `MockEnforcer`)
- External package mocks go in `mocks/` directory

---

## Common Gotchas

### 1. Context Usage

❌ **Don't:**
```go
func MyHandler(c lit.Context) error {
    // Using background context instead of request context
    ctx := context.Background()
    doSomething(ctx)
}
```

✅ **Do:**
```go
func MyHandler(c lit.Context) error {
    ctx := c.Request().Context()
    doSomething(ctx)
}
```

### 2. Error Wrapping

❌ **Don't:**
```go
if err != nil {
    return err  // Loses context
}
```

✅ **Do:**
```go
if err != nil {
    return pkgerrors.Wrap(err, "failed to process request")
}
```

### 3. Middleware Order

Remember that middleware wraps handlers like onions:

```go
r.Use(A, B, C)  // Execution: A → B → C → handler → C → B → A
```

### 4. Defer with Error Returns

❌ **Don't:**
```go
defer end()  // Doesn't capture error
return doSomething()
```

✅ **Do:**
```go
var err error
defer func() { end(err) }()
err = doSomething()
return err
```

### 5. Test Isolation

Always use mocks for external dependencies - don't make real HTTP/gRPC calls or connect to external services in unit tests.

### 6. Table-Driven Test Structure

Use map-based table-driven tests with descriptive scenario names:

```go
type arg struct {
    givenInput string
    expOutput  string
    expErr     error
}
tcs := map[string]arg{
    "scenario description": { /* test case */ },
}
for scenario, tc := range tcs {
    tc := tc  // Capture range variable
    t.Run(scenario, func(t *testing.T) {
        t.Parallel()
        // Test implementation
    })
}
```

### 6. Not Following the Map-Based Test Pattern

❌ **Problem:**
```go
tests := []struct {
    name string
    input string
}{}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) { /* ... */ })
}
```

✅ **Solution:**
```go
type arg struct {
    givenInput string
    expOutput  string
}
tcs := map[string]arg{
    "scenario": {givenInput: "test", expOutput: "result"},
}
for scenario, tc := range tcs {
    tc := tc
    t.Run(scenario, func(t *testing.T) {
        t.Parallel()
        // Given/When/Then
    })
}
```

---

## Documentation

### Where to Find Information

- **Getting Started:** `docs/getting-started.md`
- **HTTP Services:** `docs/http-services.md`
- **gRPC:** `docs/grpc.md`
- **Authentication:** `docs/auth-authorization.md`
- **Configuration:** `docs/configuration.md`
- **Monitoring:** `docs/monitoring.md`
- **Testing:** `docs/testing.md`
- **API Reference:** Package-level GoDoc comments

### Writing Documentation

- Add package-level comments explaining the purpose
- Document all exported types, functions, and interfaces
- Include usage examples in comments
- Update relevant docs in `docs/` when adding features

---

## Performance Considerations

1. **Use connection pools** for databases and HTTP clients
2. **Instrument long-running operations** with OpenTelemetry spans
3. **Set appropriate timeouts** on all external calls
4. **Use async Kafka producers** for high-throughput scenarios
5. **Cache expensive operations** with Redis when appropriate
6. **Monitor resource usage** with profiling endpoints (`WithProfiling()`)

---

## Security Best Practices

1. **Never hardcode secrets** - use environment variables
2. **Validate all inputs** using struct tags and custom validators
3. **Use HTTPS in production** - configure TLS with `ServerWithTLS()`
4. **Implement proper CORS** policies with the `cors` package
5. **Enforce authentication** on all protected routes
6. **Use Casbin** for fine-grained authorization
7. **Audit security with scorecard** (`.github/workflows/scorecard.yml`)

---

## When You're Stuck

1. **Check the documentation** in `docs/`
2. **Look at existing tests** - they demonstrate usage patterns
3. **Review similar implementations** in the codebase
4. **Check package-level comments** in `apis.go` and `spec.go` files
5. **Run tests** to verify your understanding: `make test`

---

## Key Takeaways for AI Assistants

✅ **Always prefer:**
- Functional options pattern for configuration
- Context propagation throughout the stack
- Table-driven tests with testify/require
- Error wrapping with meaningful context
- Instrumentation for observability
- Interface-based design for testability

❌ **Avoid:**
- Global state and singletons (except for tracer initialization)
- Ignoring contexts
- Missing test coverage
- Bare error returns without context
- Hardcoded configuration values
- Breaking changes to public APIs without migration path

---

## Additional Resources

- **Repository:** https://github.com/viebiz/lit
- **Go Reference:** https://pkg.go.dev/github.com/viebiz/lit
- **Contributing Guide:** `CONTRIBUTING.md`
- **License:** Apache-2.0 (see `LICENSE`)
- **Changelog:** `CHANGELOG.md`

---

**Remember:** Lit emphasizes simplicity, modularity, and observability. When extending the framework, maintain these principles and ensure backward compatibility whenever possible.