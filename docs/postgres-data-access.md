# Postgres Data Access

This guide shows how to connect to Postgres, run context-aware queries, manage transactions, and monitor the database connection in services built with this repository.

## Connection pool setup

Use [`postgres.NewPool`](../postgres/apis.go) to establish a connection pool. A `monitoring.Monitor` in the context enables structured logs during initialization.

```go
m, _ := monitoring.New(monitoring.Config{})
ctx := monitoring.SetInContext(context.Background(), m)

pool, err := postgres.NewPool(
    ctx,
    os.Getenv("DATABASE_URL"), // e.g. "postgresql://user:pass@localhost:5432/db?sslmode=disable"
    10, // max open connections
    5,  // max idle connections
    postgres.AttemptPingUponStartup(),
)
if err != nil {
    log.Fatalf("cannot create pool: %v", err)
}

// add OpenTelemetry span events for queries
pool = instrumentpg.WithInstrumentation(pool)
```

The pool uses `database/sql` under the hood, sets sane defaults, and can be tuned via options like `PoolMaxConnLifetime`. Monitoring messages during startup help trace connection issues.

## Context-aware queries

Every method accepts a `context.Context` to propagate deadlines and tracing information. For example:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

var count int
er := pool.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE active = $1`, true).Scan(&count)
if err != nil {
    return err
}
```

## Transactions

Create transactions with `BeginTx` to gain full context support. Wrap the returned transaction with instrumentation if you want span events for each statement.

```go
ctx := r.Context()
tx, err := pool.BeginTx(ctx, &sql.TxOptions{})
if err != nil {
    return err
}
txExec := instrumentpg.WithInstrumentationTx(tx)

if _, err := txExec.ExecContext(ctx, `INSERT INTO audit(msg) VALUES($1)`, "created"); err != nil {
    tx.Rollback()
    return err
}
return tx.Commit()
```

## Migration guidance

Run database migrations before starting the service. Tools such as [`tern`](https://github.com/jackc/tern) work well with `pgx`:

```bash
tern migrate --migrations ./migrations --url "$DATABASE_URL"
```

Ensure migrations are idempotent so application restarts are safe.

## Ping and health checks

`postgres.AttemptPingUponStartup` verifies connectivity during pool creation. For runtime health checks, call `PingContext` with a short timeout:

```go
func dbHealthHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), time.Second)
    defer cancel()

    if err := pool.PingContext(ctx); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
}
```

This endpoint can be wired into your router to expose database status to load balancers or orchestration systems.
