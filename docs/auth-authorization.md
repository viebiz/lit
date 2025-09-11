# Authentication and Authorization

This document describes how `lit` handles authentication and authorization.
It focuses on the `jwt`, `iam`, and `guard` packages.

## JWT: token generation and verification

The `jwt` package provides helpers for signing and parsing JSON Web Tokens.
It supports multiple signing algorithms such as HMAC and RSA. To generate a
signed token:

```go
// choose a signing method
sm := jwt.NewRS256()
// create a signer
priv, _ := rsa.GenerateKey(rand.Reader, 4096)
// fill claims
claims := jwt.RegisteredClaims{Issuer: "https://issuer.example", Subject: "user|123"}
// build and sign the token
signed, _ := jwt.NewToken(sm, claims).SignedString(priv)
```

Verification uses a parser with a key lookup function:

```go
parser := jwt.NewParser[jwt.RegisteredClaims]()
tk, err := parser.Parse(signed, func(_ string) (crypto.PublicKey, error) {
    return &priv.PublicKey, nil
})
```

The returned token contains both the header and the claims for further use.

## IAM: Casbin enforcer and profiles

The `iam` package wraps the [Casbin](https://casbin.org) enforcer. `NewEnforcer`
creates an instance backed by a PostgreSQL adapter and loads a model with a
custom `hasPermission` matcher for fine‑grained permission checks. The enforcer
expects policies of the form `(subject, object, action)` and returns an error
when the action is not permitted.

`iam` also defines helpers to extract machine‑to‑machine and user profiles from
JWT claims, storing identifiers, roles and scopes that can later be retrieved
from the request context.

## Guard middlewares

`guard` combines token validation and authorization checks for HTTP handlers.
Create an `AuthGuard` with a token `Validator` and an `Enforcer`:

```go
guard := guard.New(validator, enforcer)
```

### Token authentication

- `AuthenticateUserMiddleware` validates a user access token, extracts the user
  profile and injects it into the request context.
- `AuthenticateM2MMiddleware` performs the same for machine‑to‑machine tokens.

### Authorization helpers

- `RequiredM2MScopeMiddleware` ensures an authenticated M2M profile contains at
  least one of the required scopes.
- `RolePermissionHandler` wraps a handler and checks a user's roles against the
  Casbin enforcer before allowing access.

### Usage examples

```go
r := lit.New()

// Require M2M scope to create resources
r.POST("/resources",
    guard.AuthenticateM2MMiddleware(),
    guard.RequiredM2MScopeMiddleware("write:resources"),
    createResource)

// Allow users with READ permission on the "items" resource
r.GET("/items",
    guard.AuthenticateUserMiddleware(),
    guard.RolePermissionHandler(listItems, "items", guard.ActionRead))
```

These middlewares provide a consistent way to enforce both token validity and
fine‑grained authorization across HTTP routes.

