package postgres

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
)

type config struct {
	pgxCfg       *pgx.ConnConfig
	pool         *sql.DB
	pingUponInit bool
}

// Option is an optional config used to modify the client's behaviour
type Option func(*config)

// PoolMaxConnLifetime sets the max duration a connection should be kept alive in the pool
func PoolMaxConnLifetime(v time.Duration) Option {
	return func(c *config) {
		c.pool.SetConnMaxLifetime(v)
	}
}

// AttemptPingUponStartup will ping the DB upon startup. If this fails, NewClient will return error
func AttemptPingUponStartup() Option {
	return func(c *config) {
		c.pingUponInit = true
	}
}
