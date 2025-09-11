# Testing Guide

This project includes helper packages for writing concise unit tests.

## Shared test setup

`[test_helper.go](../test_helper.go)` provides shared helpers:

- `NewRouterForTest` creates a Gin router and context pair for HTTP tests.
- `CreateTestContext` returns a standalone context for middleware or handler checks.

Example: `[bind_test.go](../bind_test.go#L280-L286)` uses `NewRouterForTest` to run requests against a temporary router.

## `testutil/` helpers

The [`testutil`](../testutil) package wraps `go-cmp` and exports utilities such as:

- `Equal` – compare expected and actual values and fail with a diff when they differ.
- Comparison options like `IgnoreUnexported`, `IgnoreSliceOrder`, and `EquateComparable` for fine‑grained comparisons.
- Trace helpers `NewTraceID`, `NewSpanID`, and `NewTraceState` for constructing OpenTelemetry identifiers.

Examples:

- `[grpc_unary_interceptor_test.go](../grpc_unary_interceptor_test.go#L131-L133)` compares responses while ignoring unexported fields.
- `[monitoring/instrumenthttp/server_test.go](../monitoring/instrumenthttp/server_test.go#L50-L52)` builds span contexts with `NewTraceID` and `NewSpanID`.

## `mocks/` usage

Reusable mocks for external dependencies live under [`mocks/`](../mocks). They are generated with [`mockery`](https://github.com/vektra/mockery) and work with `testify/mock`.

Typical pattern:

1. Instantiate the mock type.
2. Set expectations using `On(...).Return(...)`.
3. Inject the mock into the code under test and assert expectations.

Examples:

- `[iam/enforcer_test.go](../iam/enforcer_test.go#L45-L47)` stubs the Casbin `Enforce` method.
- `[caching/redis/client_test.go](../caching/redis/client_test.go#L53-L60)` verifies Redis interactions using `MockUniversalClient`.

