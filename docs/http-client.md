# HTTP Client

The `httpclient` package offers a composable HTTP client with connection pooling, retry support and convenient request helpers. It is intended for use by applications that need robust outbound HTTP interactions.

## Connection Pool

Two constructors are provided for managing HTTP connection pools:

- `httpclient.NewSharedPool()` returns a standard `*http.Client` configured with sensible defaults and connection pooling.
- `httpclient.NewSharedCustomPool()` returns a `*httpclient.SharedCustomPool` where the client timeout is disabled so that request timeouts are controlled by the `Client` itself.

Both accept `PoolOption` functions to tweak the underlying `http.Client` and its `http.Transport`:

```go
pool := httpclient.NewSharedCustomPool(
    httpclient.OverridePoolTimeoutDuration(30 * time.Second),
    httpclient.OverridePoolMaxIdleConns(100),
    httpclient.OverridePoolMaxConnsPerHost(50),
)
```

A `PoolOption` receives both the `*http.Client` and `*http.Transport`, allowing advanced customisation such as supplying a custom transport:

```go
pool := httpclient.NewSharedCustomPool(func(_ *http.Client, t *http.Transport) {
    // Example: trust all certificates (only for testing!)
    t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
})
```

## Retry and Timeout Logic

Retries are configured through `OverrideTimeoutAndRetryOption`. You can control the retry count, per‑try timeout, overall deadline and which HTTP status codes or timeouts trigger another attempt:

```go
client, err := httpclient.NewUnauthenticated(
    httpclient.Config{
        ServiceName: "github",
        URL:         "https://api.github.com/users/:user",
        Method:      http.MethodGet,
    },
    pool,
    httpclient.OverrideTimeoutAndRetryOption(
        3,                // max retries
        5*time.Second,    // timeout per try
        20*time.Second,   // overall deadline
        true,             // retry on timeout
        []int{http.StatusTooManyRequests},
    ),
)
```

Internally the client uses exponential backoff starting at a 2 s interval and stops when either the retry count is exhausted or the overall deadline is reached.

## Request Helpers

`Client.Send` issues the request based on a `Payload`:

```go
resp, err := client.Send(ctx, httpclient.Payload{
    Body:        []byte(`{"message":"hello"}`),
    QueryParams: url.Values{"page": []string{"1"}},
    PathVars:    map[string]string{"user": "octocat"},
    Header:      map[string]string{"X-Trace-ID": traceID},
})
```

The `Payload` struct lets you specify the request body, query parameters, path variables and per‑request headers. Authentication helpers are also provided:

- `httpclient.NewWithAPIKey` attaches a static API key header.
- `httpclient.NewWithOAuth` configures OAuth2 client‑credentials flow.

## Examples

### Configure Timeouts and Custom Transport

```go
pool := httpclient.NewSharedCustomPool(
    httpclient.OverridePoolTimeoutDuration(10*time.Second),
    func(_ *http.Client, t *http.Transport) {
        t.MaxIdleConnsPerHost = 500
    },
)

client, _ := httpclient.NewUnauthenticated(
    httpclient.Config{
        ServiceName: "httpbin",
        URL:         "https://httpbin.org/get",
        Method:      http.MethodGet,
    },
    pool,
)
```

### Interacting with an External API

```go
ctx := context.Background()

resp, err := client.Send(ctx, httpclient.Payload{
    QueryParams: url.Values{"show_env": []string{"1"}},
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("status:", resp.Status)
fmt.Println(string(resp.Body))
```

The above snippet performs a GET request to [httpbin](https://httpbin.org) and prints the response body.
