# AI Coding Assistant Guide for Lightning (Lit)

> **For:** GitHub Copilot, Claude, Cursor, Codeium, and other AI coding assistants  
> **Project:** Lightning (Lit) - Modular Go framework for backend services  
> **Repository:** https://github.com/viebiz/lit

---

## 🚀 Project Overview

Lightning (Lit) is a production-ready Go framework providing modular components for:

- **HTTP/gRPC servers** with middleware and routing
- **Authentication & Authorization** (JWT, Casbin RBAC)
- **Observability** (structured logging, distributed tracing, error reporting)
- **Message brokers** (Kafka producer/consumer)
- **Caching** (Redis integration)
- **Configuration** (environment-based with Viper)
- **Testing utilities** (mocks, helpers, fixtures)

**Key Stats:**
- Go Version: 1.24.6
- License: Apache-2.0
- Test Coverage: High (CI-enforced)
- Dependencies: Gin, gRPC, Sarama, go-redis, Zap, OpenTelemetry, Casbin

---

## 🎯 Core Principles

When working with this codebase, always:

### 1. Use Functional Options Pattern

```go
// ✅ Correct way to configure components
r := lit.NewRouter(ctx,
    lit.WithLivenessEndpoint("/alive"),
    lit.WithProfiling(),
)

srv := lit.NewHttpServer(":8080", r.Handler(),
    lit.ServerShutdownGrace(10*time.Second),
    lit.ServerReadTimeout(time.Minute),
)
```

### 2. Propagate Contexts Always

```go
// ✅ Correct: Propagate request context
func MyHandler(c lit.Context) error {
    ctx := c.Request().Context()
    m := monitoring.FromContext(ctx)
    result, err := service.Process(ctx, data)
    return c.JSON(http.StatusOK, result)
}

// ❌ Wrong: Creates new context, loses tracing
func MyHandler(c lit.Context) error {
    ctx := context.Background()  // Don't do this!
    service.Process(ctx, data)
}
```

### 3. Wrap Errors with Context

```go
// ✅ Correct: Wrap with meaningful context
import pkgerrors "github.com/pkg/errors"

if err != nil {
    return pkgerrors.Wrap(err, "failed to fetch user from database")
}

// ❌ Wrong: Loses context
if err != nil {
    return err
}
```

### 4. Instrument External Calls

```go
// ✅ Correct: Instrument HTTP calls
ctx, end := instrumenthttp.StartOutgoingSegment(ctx, req)
defer end(resp.StatusCode, err)
resp, err := client.Do(req)

// ✅ Correct: Instrument gRPC calls
ctx, end := instrumentgrpc.StartUnaryCallSegment(ctx, svcInfo, method)
defer end(err)
resp, err := grpcClient.Call(ctx, req)
```

### 5. Write Tests for Everything

```go
// ✅ Use map-based table-driven tests
func TestMyFunction(t *testing.T) {
    type arg struct {
        givenInput string
        expOutput  string
        expErr     error
    }
    tcs := map[string]arg{
        "success": {givenInput: "valid", expOutput: "result"},
        "error":   {givenInput: "", expErr: errors.New("invalid input")},
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

---

## 📦 Package Structure

```
lit/
├── Core HTTP
│   ├── lit.go              # Router, Context, HandlerFunc interfaces
│   ├── router.go           # Router implementation (wraps Gin)
│   ├── context.go          # Context implementation
│   └── server.go           # HTTP server with graceful shutdown
│
├── Authentication & Authorization
│   ├── jwt/                # JWT signing and verification
│   ├── iam/                # Casbin enforcer, user/M2M profiles
│   └── guard/              # Auth middleware (user, M2M, scopes)
│
├── Observability
│   ├── monitoring/         # Monitor (Zap + Sentry)
│   ├── monitoring/instrumenthttp/     # HTTP instrumentation
│   ├── monitoring/instrumentgrpc/     # gRPC instrumentation
│   ├── monitoring/instrumentkafka/    # Kafka instrumentation
│   └── monitoring/instrumentpg/       # PostgreSQL instrumentation
│
├── Infrastructure
│   ├── broker/kafka/       # Kafka producer/consumer
│   ├── caching/redis/      # Redis client wrapper
│   ├── postgres/           # PostgreSQL utilities
│   ├── grpcclient/         # gRPC client helpers
│   └── httpclient/         # HTTP client with retry/timeout
│
├── Utilities
│   ├── env/                # Configuration (Viper)
│   ├── i18n/               # Internationalization
│   ├── middleware/http/    # HTTP middleware
│   ├── cors/               # CORS support
│   ├── ioutil/             # I/O helpers
│   └── snowflake/          # Distributed ID generation
│
└── Testing
    ├── testutil/           # Test helpers (Equal, trace IDs)
    ├── mocks/              # Generated mocks
    └── test_helper.go      # Test router/context creation
```

---

## 🛠️ Common Patterns

### Creating an HTTP Server

```go
func main() {
    ctx := context.Background()
    
    // Create router
    r := lit.NewRouter(ctx,
        lit.WithLivenessEndpoint("/alive"),
        lit.WithProfiling(),
    )
    
    // Add middleware
    r.Use(corsMiddleware(), authMiddleware())
    
    // Register routes
    r.Get("/ping", func(c lit.Context) error {
        return c.String(http.StatusOK, "pong")
    })
    
    // Group routes
    r.Group("/api/v1", func(api lit.Router) {
        api.Get("/users", listUsers)
        api.Post("/users", createUser)
    }, authMiddleware())
    
    // Start server
    srv := lit.NewHttpServer(":8080", r.Handler(),
        lit.ServerShutdownGrace(10*time.Second),
    )
    
    if err := srv.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### Writing Handlers

```go
func CreateUser(c lit.Context) error {
    // Get monitor from context
    ctx := c.Request().Context()
    m := monitoring.FromContext(ctx)
    
    // Bind request
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return err  // Automatically handled by error middleware
    }
    
    // Get user profile if authenticated
    profile, _ := iam.UserProfileFromContext(ctx)
    m = m.WithField(monitoring.String("user_id", profile.Subject))
    
    // Business logic
    user, err := userService.Create(ctx, req)
    if err != nil {
        m.Errorf(err, "failed to create user")
        return pkgerrors.Wrap(err, "create user failed")
    }
    
    m.Info("user created successfully")
    return c.JSON(http.StatusCreated, user)
}
```

### Writing Middleware

```go
func MyMiddleware() lit.HandlerFunc {
    return func(c lit.Context) error {
        ctx := c.Request().Context()
        m := monitoring.FromContext(ctx)
        
        // Pre-processing
        start := time.Now()
        m.Info("middleware start")
        
        // Call next handler
        err := c.Next()
        
        // Post-processing
        duration := time.Since(start)
        m.Infof("middleware end", 
            monitoring.Duration("duration", duration),
        )
        
        return err
    }
}
```

### Authentication Setup

```go
// Setup JWT validator and Casbin enforcer
validator := jwt.NewValidator(publicKey)
enforcer, _ := iam.NewEnforcer(db, modelPath)

// Create guard
guard := guard.New(validator, enforcer)

// Protect routes
api := r.Group("/api", guard.AuthenticateUserMiddleware())

// Require specific scopes for M2M
admin := r.Group("/admin",
    guard.AuthenticateM2MMiddleware(),
    guard.RequiredM2MScopeMiddleware("write:admin"),
)

// Check permissions with Casbin
r.Post("/resources",
    guard.AuthenticateUserMiddleware(),
    guard.RolePermissionHandler(enforcer, "resources", "write", createResource),
)
```

### Kafka Producer

```go
// Create async producer
producer, err := kafka.NewAsyncProducer(cfg,
    kafka.WithSuccessHandler(func(msg *sarama.ProducerMessage) {
        log.Println("sent:", msg.Offset)
    }),
    kafka.WithErrorHandler(func(err *sarama.ProducerError) {
        log.Printf("error: %v", err)
    }),
)

// Send message with instrumentation
msg := &sarama.ProducerMessage{
    Topic: "events",
    Key:   sarama.StringEncoder(key),
    Value: sarama.ByteEncoder(payload),
}

ctx, end := instrumentkafka.StartAsyncEnqueueSegment(ctx, svcInfo, clientID, msg, key)
publisher := producer.SendMessage(msg)
end(publisher)
```

### Kafka Consumer

```go
consumer, err := kafka.NewConsumer(cfg,
    kafka.WithConsumerGroup("my-service"),
    kafka.WithTopics("events"),
)

handler := func(ctx context.Context, msg *sarama.ConsumerMessage) error {
    m := monitoring.FromContext(ctx)
    m.Info("processing message",
        monitoring.String("topic", msg.Topic),
        monitoring.Int64("offset", msg.Offset),
    )
    
    // Process message
    return processEvent(ctx, msg.Value)
}

if err := consumer.Consume(ctx, handler); err != nil {
    log.Fatal(err)
}
```

### Configuration

```go
type AppConfig struct {
    ServerAddr string `mapstructure:"server_addr"`
    DBUrl      string `mapstructure:"db_url"`
    KafkaURL   string `mapstructure:"kafka_url"`
}

func (c AppConfig) Validate() error {
    if c.ServerAddr == "" {
        return errors.New("server_addr required")
    }
    return nil
}

// Reads .env file + APP_* environment variables
cfg, err := env.ReadAppConfig[AppConfig]()
if err != nil {
    log.Fatal(err)
}
```

---

## ✅ Code Style Checklist

Before committing code, ensure:

- [ ] Run `go fmt ./...`
- [ ] All tests pass: `make test`
- [ ] Added tests for new code
- [ ] Errors wrapped with `pkgerrors.Wrap`
- [ ] Contexts propagated correctly
- [ ] External calls instrumented
- [ ] No hardcoded configuration
- [ ] Documentation updated if adding public APIs
- [ ] Mocks generated if interfaces changed: `make gen-mocks`

---

## 🧪 Testing Guidelines

### Test File Naming
- Test files: `*_test.go`
- Same package as code under test
- Use `require` (not `assert`) from testify
- Use map-based table-driven tests: `tcs := map[string]arg{...}`
- Always use `t.Parallel()` for independent tests
- Structure tests with `Given/When/Then` comments

### Test Structure

```go
func TestMyFunction(t *testing.T) {
    type arg struct {
        givenInput string
        expOutput  string
        expErr     error
    }
    tcs := map[string]arg{
        "success": {
            givenInput: "test",
            expOutput:  "expected",
        },
        "error": {
            givenInput: "",
            expErr:     errors.New("invalid input"),
        },
    }
    
    for scenario, tc := range tcs {
        tc := tc
        t.Run(scenario, func(t *testing.T) {
            t.Parallel()
            
            // Given
            ctx := context.Background()
            mock := new(MockDependency)
            mock.On("Method", mock.Anything).Return(tc.expOutput, tc.expErr)
            
            // When
            got, err := MyFunction(ctx, mock, tc.givenInput)
            
            // Then
            if tc.expErr != nil {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            testutil.Equal(t, tc.expOutput, got)
            mock.AssertExpectations(t)
        })
    }
}
```

### HTTP Handler Tests

```go
func TestMyHandler(t *testing.T) {
    r, _ := lit.NewRouterForTest()
    r.Post("/users", MyHandler)
    
    body := `{"name":"John"}`
    req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    r.Handler().ServeHTTP(w, req)
    
    require.Equal(t, http.StatusCreated, w.Code)
}
```

### Running Tests

```bash
make test                    # All tests
go test -v ./guard/...       # Specific package
go test -run TestMyFunc      # Specific test
make benchmark               # Benchmarks
```

---

## 🚨 Common Mistakes to Avoid

| ❌ Don't | ✅ Do |
|---------|------|
| `context.Background()` in handlers | `c.Request().Context()` |
| `return err` without wrapping | `return pkgerrors.Wrap(err, "context")` |
| Skip instrumentation | Always instrument external calls |
| Forget `defer end()` | Always defer instrumentation end functions |
| Hardcode config | Use `env.ReadAppConfig` |
| Real connections in tests | Use mocks |
| `assert.Equal` | `require.Equal` (fails fast) |
| Global state | Pass dependencies via constructors |
| `tests := []struct{...}` | `tcs := map[string]arg{...}` |
| Skip `t.Parallel()` | Always use `t.Parallel()` |

---

## 📚 Quick Command Reference

```bash
# Development
make init                    # Start Docker services
make test                    # Run all tests
make gen-mocks              # Generate mocks
make gen-proto              # Generate protobuf
make benchmark              # Run benchmarks
make tear-down              # Stop all services

# Testing
go test ./...               # All tests
go test -v ./pkg/...        # Verbose
go test -cover ./...        # With coverage
go test -run TestName       # Specific test

# Code Quality
go fmt ./...                # Format code
go vet ./...                # Vet code
go mod tidy                 # Clean dependencies
go mod vendor               # Vendor dependencies
```

---

## 🔗 Important Files

- **Contributing:** `CONTRIBUTING.md`
- **Documentation:** `docs/` directory
- **Examples:** Test files (`*_test.go`)
- **Configuration:** `.mockery.yaml`, `.editorconfig`
- **CI/CD:** `.github/workflows/ci.yml`

---

## 💡 Tips for AI Assistants

When generating code:

1. **Look at existing patterns** in similar files first
2. **Follow the functional options pattern** for constructors
3. **Always add tests** with your code
4. **Use map-based table-driven tests** (`tcs := map[string]arg{...}`)
5. **Mock external dependencies** in tests
6. **Instrument external calls** (HTTP, gRPC, Kafka, DB)
7. **Propagate contexts** - never create `Background()` mid-chain
8. **Wrap errors** with meaningful messages
9. **Add logging** with structured fields
10. **Document public APIs** with GoDoc comments
11. **Use `t.Parallel()`** for independent test cases
12. **Structure tests** with Given/When/Then comments

---

## 📖 Documentation

For detailed guides, see the `docs/` directory:

- `getting-started.md` - Project overview and hello world
- `http-services.md` - HTTP server and routing
- `grpc.md` - gRPC services
- `auth-authorization.md` - JWT and Casbin
- `monitoring.md` - Logging, tracing, error reporting
- `kafka-messaging.md` - Kafka producer/consumer
- `redis-caching.md` - Redis integration
- `testing.md` - Testing utilities and patterns

---

**Project Philosophy:** Lightning (Lit) prioritizes developer experience, production readiness, and observability. Code should be clear, testable, and instrumented. When in doubt, follow existing patterns in the codebase.