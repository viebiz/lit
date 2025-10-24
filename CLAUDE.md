# Claude AI Assistant Instructions for Lightning (Lit)

This document provides comprehensive guidance for Claude AI when working with the **Lightning (Lit)** project - a modular, production-ready Go framework for building reliable backend services.

---

## 🎯 Project Mission

Lightning (Lit) is a comprehensive Go framework that eliminates boilerplate and accelerates backend development by providing battle-tested, modular components for common infrastructure tasks. The framework emphasizes:

- **Developer Experience:** Intuitive APIs with functional options pattern
- **Production Ready:** Built-in observability, error handling, and graceful degradation
- **Modularity:** Use only what you need - no vendor lock-in
- **Testability:** Comprehensive testing utilities and mock generation

**Repository:** `github.com/viebiz/lit`  
**Go Version:** 1.24.6  
**License:** Apache-2.0

---

## 📐 Architectural Philosophy

### Core Design Principles

1. **Context-First Design**
   - Application context flows from `main()` through all layers
   - Request contexts carry monitoring, tracing, user identity, and request metadata
   - Never create `context.Background()` in the middle of a call chain
   - Always propagate contexts to maintain observability chain

2. **Functional Options Pattern**
   - Constructors accept variadic option functions for flexibility
   - Options provide sensible defaults with explicit overrides
   - Type-safe configuration prevents runtime errors
   - Example: `NewRouter(ctx, WithLivenessEndpoint("/alive"), WithProfiling())`

3. **Interface Segregation**
   - Small, focused interfaces (e.g., `Route`, `Router`, `Context`)
   - Easy to mock and test in isolation
   - Implementations are private; APIs are public

4. **Observability by Default**
   - Every external interaction is instrumented (HTTP, gRPC, Kafka, PostgreSQL)
   - Structured logging with Zap
   - Distributed tracing with OpenTelemetry
   - Error reporting with Sentry integration

5. **Fail-Fast with Graceful Degradation**
   - Validate inputs early with clear error messages
   - Wrap errors with context using `pkg/errors`
   - Support graceful shutdown with configurable grace periods
   - Panic recovery in all HTTP/gRPC handlers

---

## 🗂️ Package Architecture

### Core HTTP Components

**`lit.go`, `router.go`, `context.go`, `server.go`**

The foundation of HTTP services:

```go
// Router wraps gin.Engine with lit's context and middleware
type Router interface {
    Route                  // HTTP method registration
    Group(prefix string, routerFunc func(Router), middleware ...HandlerFunc) Router
    Route(prefix string, middleware ...HandlerFunc) Router
    Routes() RoutesInfo
    Handler() http.Handler
}

// Context wraps gin.Context with additional helpers
type Context interface {
    context.Context        // Standard context methods
    Request() *http.Request
    Writer() ResponseWriter
    Bind(obj interface{}) error
    JSON(code int, obj interface{}) error
    // ... 30+ helper methods
}

// HandlerFunc is the signature for all handlers and middleware
type HandlerFunc func(Context) error
```

**Key Pattern:** The router uses Gin internally but exposes a cleaner API with error-returning handlers.

### Authentication & Authorization Stack

**`jwt/`, `iam/`, `guard/`**

Three-layer security model:

1. **JWT Package:** Token generation, parsing, and verification
   - Supports HMAC and RSA signing
   - RFC 9068 (JWT Profile for OAuth 2.0 Access Tokens) compliant
   - JWK (JSON Web Key) support for key rotation

2. **IAM Package:** Identity and access management
   - Wraps Casbin enforcer for RBAC/ABAC
   - Custom `hasPermission` function for fine-grained checks
   - PostgreSQL adapter for policy storage
   - User and M2M (machine-to-machine) profile extraction from JWT

3. **Guard Package:** HTTP middleware integration
   - `AuthenticateUserMiddleware()` - validates user JWT, extracts profile
   - `AuthenticateM2MMiddleware()` - validates M2M JWT, extracts client info
   - `RequiredM2MScopeMiddleware(scopes...)` - enforces OAuth2 scopes
   - `RolePermissionHandler(enforcer, resource, action, handler)` - Casbin integration

**Usage Pattern:**
```go
guard := guard.New(jwtValidator, casbinEnforcer)

api := r.Group("/api/v1", guard.AuthenticateUserMiddleware())
api.Get("/profile", getProfile)
api.Post("/admin", 
    guard.RequiredM2MScopeMiddleware("write:admin"),
    adminAction,
)
```

### Monitoring & Observability Stack

**`monitoring/`, `monitoring/instrument*/`**

Unified observability layer:

```go
// Monitor combines structured logging and error reporting
type Monitor struct {
    logger *zap.Logger
    sentry *sentry.Hub
}

// Methods automatically add context and propagate to Sentry
m.Info("message")
m.Infof("template %s", value)
m.Error(err)
m.Errorf(err, "context: %s", detail)
```

**Instrumentation Packages:**

- **`instrumenthttp/`** - HTTP server and client spans
  - `StartIncomingRequest(m, r, routePattern)` - creates span for incoming HTTP
  - `StartOutgoingGroupSegment(ctx, svc, name, method, url)` - groups retry attempts
  - `StartOutgoingSegment(ctx, req)` - instruments individual HTTP call

- **`instrumentgrpc/`** - gRPC unary interceptors
  - `StartUnaryIncomingCall(ctx, m, fullMethod, req)` - server-side
  - `StartUnaryCallSegment(ctx, svcInfo, fullMethod)` - client-side

- **`instrumentkafka/`** - Kafka producer tracing
  - `StartSyncPublishSegment(ctx, svc, clientID, msg, key)` - sync producer
  - `StartAsyncEnqueueSegment(ctx, svc, clientID, msg, key)` - async producer

- **`instrumentpg/`** - PostgreSQL query tracing (via GORM callbacks)

**Pattern:** All instrumentations return an `end()` function to be called with `defer`:

```go
ctx, reqMeta, end := instrumenthttp.StartIncomingRequest(m, r, "/users/:id")
defer end(statusCode, err)

// Business logic here
```

### Message Broker Integration

**`broker/kafka/`**

Production-ready Kafka integration with:

- Async and sync producer implementations
- Consumer group management with automatic rebalancing
- Instrumentation for distributed tracing
- Configurable retry with exponential backoff
- Error and success callback handlers
- TLS support

**Producer Pattern:**
```go
producer, err := kafka.NewAsyncProducer(cfg,
    kafka.WithSuccessHandler(onSuccess),
    kafka.WithErrorHandler(onError),
)

ctx, end := instrumentkafka.StartAsyncEnqueueSegment(ctx, svcInfo, clientID, msg, key)
publisher := producer.SendMessage(msg)
end(publisher)
```

**Consumer Pattern:**
```go
consumer, err := kafka.NewConsumer(cfg,
    kafka.WithConsumerGroup("my-group"),
    kafka.WithTopics("topic1", "topic2"),
)

handler := func(ctx context.Context, msg *sarama.ConsumerMessage) error {
    // Process message with automatic instrumentation
    return nil
}

err = consumer.Consume(ctx, handler)
```

### Data Access Layer

**`postgres/`, `caching/redis/`**

PostgreSQL utilities:
- Connection management with pgx
- GORM integration with automatic instrumentation
- Migration helpers

Redis caching:
- Wrapper around go-redis/v9
- Context-aware operations
- Instrumented commands

### Configuration Management

**`env/`**

Viper-based configuration with environment override:

```go
type AppConfig struct {
    ServerAddr string `mapstructure:"server_addr"`
    DBUrl      string `mapstructure:"db_url"`
}

func (c AppConfig) Validate() error {
    // Custom validation logic
    return nil
}

// Reads from .env file + APP_* environment variables
cfg, err := env.ReadAppConfig[AppConfig]()
```

**Environment Variable Mapping:**
- File: `server.addr: ":8080"`
- Env: `APP_SERVER_ADDR=":9090"` (takes precedence)
- Delimiter `.` becomes `_` in environment variables

### Internationalization

**`i18n/`**

Built on `go-i18n/v2`:

- Load translations from JSON/YAML files
- Context-aware localization
- Pluralization support
- Template interpolation

### HTTP Client Utilities

**`httpclient/`**

Instrumented HTTP client builder:
- Timeout configuration
- Retry with exponential backoff
- Automatic tracing propagation
- Circuit breaker pattern support

### Testing Infrastructure

**`testutil/`, `mocks/`, `test_helper.go`**

Comprehensive testing support:

1. **testutil Package:**
   - `Equal(t, expected, actual)` - deep comparison with diff output
   - `IgnoreUnexported(types...)` - compare options
   - `NewTraceID()`, `NewSpanID()` - OpenTelemetry test helpers

2. **Mock Generation:**
   - Configuration in `.mockery.yaml`
   - Mocks generated with `make gen-mocks`
   - Testify-compatible mock objects

3. **Test Helpers:**
   - `NewRouterForTest()` - creates router + Gin context for HTTP tests
   - `CreateTestContext()` - standalone context for unit tests

---

## 🎨 Code Style & Conventions

### Formatting Rules (from `.editorconfig`)

```
- Charset: UTF-8
- Line ending: LF (Unix-style)
- Indentation: Tabs for Go files, spaces elsewhere
- Max line length: 120 characters
- Insert final newline: true
```

**Always run before committing:**
```bash
go fmt ./...
```

### Naming Conventions

| Type | Pattern | Examples |
|------|---------|----------|
| Interfaces | Descriptive noun, no "I" prefix | `Router`, `Context`, `Enforcer` |
| Implementations | Lowercase private struct | `router`, `context`, `enforcer` |
| Constructors | `New` prefix | `NewRouter`, `NewHttpServer`, `NewMonitor` |
| Option Functions | `With` prefix | `WithLivenessEndpoint`, `WithProfiling` |
| Test Helpers | `ForTest` suffix | `NewRouterForTest`, `CreateTestContext` |
| Error Variables | `Err` prefix | `ErrEmptyTopic`, `ErrInvalidToken` |

### Package Organization

Standard file naming pattern:

- `apis.go` - Public API functions and constructors
- `spec.go` - Interface definitions and type contracts
- `type.go` - Struct type definitions
- `errors.go` - Error type and constant definitions
- `*_option.go` - Option function definitions
- `*_test.go` - Test files (same package)
- `mock_*.go` - Generated mocks (by mockery)

### Error Handling Pattern

```go
import pkgerrors "github.com/pkg/errors"

// Wrap errors with context
if err != nil {
    return pkgerrors.Wrap(err, "failed to connect to database")
}

// Wrap with formatted message
if err != nil {
    return pkgerrors.Wrapf(err, "failed to process user %s", userID)
}

// Define sentinel errors
var (
    ErrEmptyTopic = errors.New("topic is empty")
    ErrNotFound   = errors.New("resource not found")
)
```

### Context Handling Pattern

```go
// ✅ CORRECT: Always propagate request context
func MyHandler(c lit.Context) error {
    ctx := c.Request().Context()
    m := monitoring.FromContext(ctx)
    
    // Pass ctx to all downstream calls
    result, err := service.Process(ctx, data)
    if err != nil {
        m.Errorf(err, "failed to process")
        return err
    }
    
    return c.JSON(http.StatusOK, result)
}

// ❌ WRONG: Creating background context loses tracing
func MyHandler(c lit.Context) error {
    ctx := context.Background()  // Loses tracing, user info, etc.
    service.Process(ctx, data)
}
```

---

## 🧪 Testing Philosophy

### Test Coverage Requirements

**Mandatory:**
- All public functions must have tests
- All error paths must be tested
- Critical business logic must have >90% coverage
- CI/CD blocks merge if tests fail

### Table-Driven Test Pattern

```go
func TestMyFunction(t *testing.T) {
    type arg struct {
        givenInput InputType
        expOutput  OutputType
        expErr     error
    }
    tcs := map[string]arg{
        "successful case": {
            givenInput: InputType{Field: "value"},
            expOutput:  OutputType{Result: "expected"},
        },
        "validation error": {
            givenInput: InputType{},
            expErr:     ErrInvalidInput,
        },
        "processing error": {
            givenInput: InputType{Field: "bad"},
            expErr:     ErrProcessingFailed,
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
                require.ErrorIs(t, err, tc.expErr)
                return
            }
            
            require.NoError(t, err)
            testutil.Equal(t, tc.expOutput, got)
        })
    }
}
```

### Mocking External Dependencies

```go
func TestServiceWithMock(t *testing.T) {
    // Create mock
    mockClient := new(MockHTTPClient)
    
    // Set expectations
    mockClient.On("Get", mock.Anything, "https://api.example.com/data").
        Return(&Response{Data: "test"}, nil).
        Once()
    
    // Inject mock
    service := NewService(mockClient)
    
    // Test
    result, err := service.FetchData(context.Background())
    require.NoError(t, err)
    require.Equal(t, "test", result.Data)
    
    // Verify expectations
    mockClient.AssertExpectations(t)
}
```

### HTTP Handler Testing

```go
func TestMyHandler(t *testing.T) {
    r, c := NewRouterForTest()
    r.Post("/users", MyHandler)
    
    body := `{"name":"John","email":"john@example.com"}`
    req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    r.Handler().ServeHTTP(w, req)
    
    require.Equal(t, http.StatusCreated, w.Code)
    
    var resp UserResponse
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    require.NoError(t, err)
    require.Equal(t, "John", resp.Name)
}
```

---

## 🔧 Development Workflow

### Local Development Setup

```bash
# 1. Clone repository
git clone https://github.com/viebiz/lit.git
cd lit

# 2. Start infrastructure services
make init  # Starts Postgres, Redis, Kafka via Docker Compose

# 3. Run tests
make test

# 4. Generate mocks (after interface changes)
make gen-mocks

# 5. Generate protobuf code (if proto files changed)
make gen-proto

# 6. Run benchmarks
make benchmark

# 7. Cleanup
make tear-down
```

### Docker Services (docker-compose.yml)

| Service | Port | Credentials |
|---------|------|-------------|
| PostgreSQL | 54321 | user: `lit`, db: `master`, no password |
| Redis | 63791 | No auth |
| Kafka | 9092 | KRaft mode (no Zookeeper) |
| Jaeger UI | 16686 | Tracing dashboard |
| Jaeger OTLP | 4317 | OpenTelemetry collector |

### Git Workflow

```bash
# 1. Create feature branch
git checkout -b feature/add-awesome-feature

# 2. Make changes with focused commits
git add .
git commit -m "feat: add awesome feature"

# 3. Ensure tests pass
make test

# 4. Format code
go fmt ./...

# 5. Push and create PR
git push origin feature/add-awesome-feature
```

### Commit Message Convention

Follow conventional commits:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test additions or changes
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `chore:` - Maintenance tasks

---

## 🚨 Common Pitfalls & Solutions

### Pitfall 1: Losing Context Chain

❌ **Problem:**
```go
func Handler(c lit.Context) error {
    ctx := context.Background()  // Breaks tracing!
    service.Call(ctx)
}
```

✅ **Solution:**
```go
func Handler(c lit.Context) error {
    ctx := c.Request().Context()  // Preserves tracing
    service.Call(ctx)
}
```

### Pitfall 2: Not Calling Instrumentation End Function

❌ **Problem:**
```go
ctx, end := instrumenthttp.StartOutgoingSegment(ctx, req)
// Forgot to call end() - span never closes!
resp, err := client.Do(req)
```

✅ **Solution:**
```go
ctx, end := instrumenthttp.StartOutgoingSegment(ctx, req)
defer end(resp.StatusCode, err)  // Always defer!
resp, err := client.Do(req)
```

### Pitfall 3: Ignoring Error Context

❌ **Problem:**
```go
if err != nil {
    return err  // What failed? Where?
}
```

✅ **Solution:**
```go
if err != nil {
    return pkgerrors.Wrap(err, "failed to fetch user profile from database")
}
```

### Pitfall 4: Middleware Execution Order Confusion

Middleware wraps like an onion - first registered executes first (outer) and last (after handler):

```go
r.Use(LoggingMiddleware, AuthMiddleware, ValidationMiddleware)

// Execution order:
// 1. LoggingMiddleware (pre)
// 2. AuthMiddleware (pre)
// 3. ValidationMiddleware (pre)
// 4. Handler
// 5. ValidationMiddleware (post)
// 6. AuthMiddleware (post)
// 7. LoggingMiddleware (post)
```

### Pitfall 5: Not Mocking in Unit Tests

❌ **Problem:**
```go
func TestMyService(t *testing.T) {
    svc := NewService()
    // Makes real HTTP call - test is slow, flaky, requires network
    result, err := svc.FetchData(context.Background())
}
```

✅ **Solution:**
```go
func TestMyService(t *testing.T) {
    mockClient := new(MockHTTPClient)
    mockClient.On("Get", mock.Anything, mock.Anything).
        Return(&Response{Data: "test"}, nil)
    
    svc := NewService(WithHTTPClient(mockClient))
    result, err := svc.FetchData(context.Background())
}
```

---

## 📚 Quick Reference Guide

### Most Common Tasks

**Creating an HTTP Server:**
```go
r := lit.NewRouter(ctx, lit.WithLivenessEndpoint("/alive"))
r.Get("/hello", func(c lit.Context) error {
    return c.String(http.StatusOK, "Hello!")
})
srv := lit.NewHttpServer(":8080", r.Handler())
srv.Run()
```

**Writing Table-Driven Tests:**
```go
type arg struct {
    givenInput string
    expOutput  string
    expErr     error
}
tcs := map[string]arg{
    "success case": {givenInput: "test", expOutput: "result"},
    "error case": {givenInput: "", expErr: ErrInvalid},
}
for scenario, tc := range tcs {
    tc := tc
    t.Run(scenario, func(t *testing.T) {
        t.Parallel()
        // Given/When/Then
        got, err := MyFunc(tc.givenInput)
        if tc.expErr != nil {
            require.Error(t, err)
            return
        }
        require.NoError(t, err)
        testutil.Equal(t, tc.expOutput, got)
    })
}
```

**Adding Authentication:**
```go
guard := guard.New(validator, enforcer)
api := r.Group("/api", guard.AuthenticateUserMiddleware())
api.Get("/profile", getProfile)
```

**Logging with Context:**
```go
m := monitoring.FromContext(ctx)
m = m.WithField(monitoring.String("user_id", userID))
m.Info("processing request")
```

**Creating Spans:**
```go
ctx, reqMeta, end := instrumenthttp.StartIncomingRequest(m, r, "/users/:id")
defer end(http.StatusOK, nil)
```

**Publishing to Kafka:**
```go
producer, _ := kafka.NewAsyncProducer(cfg)
msg := &sarama.ProducerMessage{Topic: "events", Value: sarama.StringEncoder(data)}
producer.SendMessage(msg)
```

**Reading Configuration:**
```go
cfg, err := env.ReadAppConfig[MyConfig]()
```

---

## 🎯 When Assisting with Code

### As Claude, I should:

1. **Always check existing patterns** before suggesting new approaches
2. **Prefer composition over reinvention** - use existing utilities
3. **Add comprehensive tests** for all new code
4. **Include error wrapping** with meaningful context
5. **Add instrumentation** for external calls
6. **Follow functional options pattern** for new constructors
7. **Update documentation** when adding public APIs
8. **Consider backward compatibility** - avoid breaking changes
9. **Use map-based table-driven tests** with `Given/When/Then` structure
10. **Mock external dependencies** in unit tests
11. **Always use `t.Parallel()`** for independent test cases
12. **Capture range variables** with `tc := tc` before goroutines

### I should NOT:

1. Create `context.Background()` in the middle of request processing
2. Return bare errors without wrapping
3. Add code without corresponding tests
4. Hardcode configuration values
5. Skip instrumentation on external calls
6. Use global variables for state
7. Break existing public APIs
8. Ignore the functional options pattern
9. Create real connections in unit tests
10. Forget to call `defer end()` on instrumentation
11. Use slice-based table-driven tests (`[]struct`) instead of map-based
12. Skip `t.Parallel()` in test cases

---

## 📖 Documentation References

For more detailed information:

- **Getting Started:** `docs/getting-started.md`
- **HTTP Services:** `docs/http-services.md`
- **gRPC Services:** `docs/grpc.md`
- **HTTP Middleware:** `docs/http-middleware.md`
- **Authentication:** `docs/auth-authorization.md`
- **Configuration:** `docs/configuration.md`
- **Localization:** `docs/localization.md`
- **Redis Caching:** `docs/redis-caching.md`
- **Kafka Messaging:** `docs/kafka-messaging.md`
- **Postgres:** `docs/postgres-data-access.md`
- **HTTP Client:** `docs/http-client.md`
- **Monitoring:** `docs/monitoring.md`
- **Utilities:** `docs/utilities.md`
- **Testing:** `docs/testing.md`

---

## 🤝 Contributing Philosophy

Lightning (Lit) welcomes contributions that:

- Maintain the existing design philosophy
- Add value without adding complexity
- Come with comprehensive tests and documentation
- Consider backward compatibility
- Follow the established code style

See `CONTRIBUTING.md` for full guidelines.

---

**Remember:** When working with this codebase, clarity and observability are paramount. Every function should be testable, every error should be traceable, and every external interaction should be instrumented. Write code that the next developer (including future you) will thank you for.