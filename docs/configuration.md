# Configuration Management

The `env` package provides a simple way to load configuration from files and override values with environment variables. It uses [`viper`](https://github.com/spf13/viper) under the hood.

## Default loading

`ReadAppConfig` reads `config.env` in the current directory using the `env` format:

```go
cfg, err := env.ReadAppConfig[Config]()
```

The default behavior:

- Config path: `.`
- Config type: `env`
- Config file: `config.env`
- Environment variable prefix: `APP_`

## Environment overrides

Environment variables override file values. Keys are prefixed with `APP_` and nested fields use underscores:

```bash
APP_WEB_HOST=192.168.0.1
APP_LANG=vi,fr
```

## Defining structs

Define a struct that mirrors your configuration. Use `mapstructure` tags to map fields:

```go
type Database struct {
    URL string `mapstructure:"URL"`
}

type Web struct {
    Host string `mapstructure:"HOST"`
    Port int    `mapstructure:"PORT"`
}

type Config struct {
    AppName string   `mapstructure:"APP_NAME"`
    Lang    []string `mapstructure:"LANG"`
    DB      Database `mapstructure:"DB"`
    Web     Web      `mapstructure:"WEB"`
}
```

## Using the config in a service

```go
package main

import (
    "log"

    "github.com/viebiz/lit/env"
)

func main() {
    cfg, err := env.ReadAppConfig[Config]()
    if err != nil {
        log.Fatal(err)
    }

    // use cfg.DB.URL, cfg.Web.Host, ...
}
```

`ReadAppConfigWithOptions` allows loading from other formats or locations:

```go
cfg, err := env.ReadAppConfigWithOptions[Config]("yaml", "config", "./configs")
```

## Sample `config.yaml`

```yaml
APP_NAME: lightning
LANG:
  - en
  - vi
DB:
  URL: postgres:thisisurl
WEB:
  HOST: 0.0.0.0
  PORT: 8080
```

Running with environment overrides:

```bash
APP_WEB_HOST=192.168.0.1 APP_LANG=vi,fr go run ./cmd/app
```

The environment variables replace the corresponding values from the file.
