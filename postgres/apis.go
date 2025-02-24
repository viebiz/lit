package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	pkgerrors "github.com/pkg/errors"
	"github.com/viebiz/lit/monitoring"
)

// NewPool opens a new DB connection pool, pings it and returns the pool
func NewPool(
	ctx context.Context,
	url string,
	poolMaxOpenConns int,
	poolMaxIdleConns int,
	opts ...Option,
) (BeginnerExecutor, error) {
	monitor := monitoring.FromContext(ctx)

	monitor.Infof("Initializing Postgres")

	connCfg, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, pkgerrors.WithStack(fmt.Errorf("parsing pgx config failed. err: %w", err))
	}

	pool, err := sql.Open("pgx", stdlib.RegisterConnConfig(connCfg))
	if err != nil {
		return nil, pkgerrors.WithStack(fmt.Errorf("opening DB failed. err: %w", err))
	}
	pool.SetConnMaxLifetime(29 * time.Minute) // Azure's default is 30 mins.
	pool.SetMaxOpenConns(poolMaxOpenConns)
	pool.SetMaxIdleConns(poolMaxIdleConns)
	cfg := config{
		pgxCfg: connCfg,
		pool:   pool,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	monitor.Infof("Postgres client created %s", connCfg.Database)

	if cfg.pingUponInit {
		monitor.Infof("Pinging DB...")
		if err = pool.PingContext(ctx); err != nil {
			return nil, pkgerrors.WithStack(fmt.Errorf("unable to ping DB. err: %w", err))
		}
		monitor.Infof("DB ping successful")
	}

	monitor.Infof("Postgres initialized")

	return pool, nil
}
