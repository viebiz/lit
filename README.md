# Lightning

Lightning is a modular, production-ready collection of Go libraries that helps backend developers build reliable applications faster. It reduces boilerplate by offering reusable components for key tasks such as:

- HTTP routing with authentication guards and middleware
- gRPC server and client setup with service registration
- Internationalization (i18n) support
- Kafka integration 
- Redis integration
- Structured, context-aware logging
- Configuration management with environment overrides

Designed to be simple, extensible, and easy to adopt, Lit lets teams focus on business logic instead of infrastructure scaffolding.

## Project status

| Name            | Status                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
|-----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Pipeline        | ![GitHub branch status](https://img.shields.io/github/checks-status/viebiz/lit/main) [![CircleCI](https://dl.circleci.com/status-badge/img/circleci/Nur6mXEFG9qEiztTeZh7R9/AKZQkEe9aCbcR1kLJk4amp/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/circleci/Nur6mXEFG9qEiztTeZh7R9/AKZQkEe9aCbcR1kLJk4amp/tree/main)                                                                                                                                                                                                                                                                                                            |
| Coverage        | [![Coverage Status](https://coveralls.io/repos/github/viebiz/lit/badge.svg?branch=main)](https://coveralls.io/github/viebiz/lit?branch=main) [![Codacy Badge](https://app.codacy.com/project/badge/Coverage/c6d7a11459994e3984fd2ae2008839d1)](https://app.codacy.com/gh/viebiz/lit/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage) [![codecov](https://codecov.io/github/viebiz/lit/graph/badge.svg?token=MIJM38CDIP)](https://codecov.io/github/viebiz/lit)                                                                                                                                                       |
| Code Quality    | [![Codacy Badge](https://app.codacy.com/project/badge/Grade/c6d7a11459994e3984fd2ae2008839d1)](https://app.codacy.com/gh/viebiz/lit/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade) [![Open Source Helpers](https://www.codetriage.com/viebiz/lit/badges/users.svg)](https://www.codetriage.com/viebiz/lit) [![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/viebiz/lit/badge?style=flat)](https://securityscorecards.dev/viewer/?uri=github.com/viebiz/lit) [![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10175/badge)](https://www.bestpractices.dev/projects/10175) |
| Go Report       | [![Go Report Card](https://goreportcard.com/badge/github.com/viebiz/lit)](https://goreportcard.com/report/github.com/viebiz/lit)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| Go Reference    | [![Go Reference](https://pkg.go.dev/badge/github.com/viebiz/lit?status.svg)](https://pkg.go.dev/github.com/viebiz/lit?tab=doc)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| Release Version | ![GitHub Release](https://img.shields.io/github/v/release/viebiz/lit)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| Tag Version     | ![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/viebiz/lit)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| License         | ![GitHub License](https://img.shields.io/github/license/viebiz/lit)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |


## 🤖 AI Assistant Support

Lightning (Lit) includes comprehensive instructions for AI coding assistants to help you develop faster:

- **[GitHub Copilot & Cursor](.github/copilot-instructions.md)** - Complete framework reference for inline suggestions
- **[Claude AI](CLAUDE.md)** - Detailed architectural guide with common pitfalls and best practices
- **[Generic AI Assistants](AGENTS.md)** - Quick-start guide for all AI coding tools

These files provide guidance on:
- Core architectural patterns (functional options, context propagation)
- Testing requirements and best practices
- Common development workflows
- Security and performance considerations

See [.github/AI-INSTRUCTIONS-README.md](.github/AI-INSTRUCTIONS-README.md) for details on which file to use.

## 📚 Documentation

Detailed documentation is available in the [`docs/`](docs/) directory:

- [Getting Started](docs/getting-started.md) - Installation and hello world
- [HTTP Services](docs/http-services.md) - Building HTTP APIs
- [gRPC Services](docs/grpc.md) - gRPC server and client setup
- [Authentication & Authorization](docs/auth-authorization.md) - JWT and Casbin integration
- [Monitoring](docs/monitoring.md) - Logging, tracing, and error reporting
- [Kafka Messaging](docs/kafka-messaging.md) - Message broker integration
- [Redis Caching](docs/redis-caching.md) - Caching layer
- [Testing](docs/testing.md) - Testing utilities and patterns

## 🚀 Getting Started

Install Lightning:

```bash
go get github.com/viebiz/lit@latest
```

Create a simple HTTP server:

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

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:

- Code style and formatting
- Testing requirements
- Development workflow
- Pull request process

## 📝 License

Lightning is licensed under the Apache-2.0 License. See [LICENSE](LICENSE) for details.
