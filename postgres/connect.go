package postgres

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	pkgerrors "github.com/pkg/errors"
)

const (
	defaultMaxIdleConn     = 10
	defaultMaxOpenConn     = 200
	defaultMaxIdleTime     = 10 * time.Minute
	defaultMaxConnLifetime = 30 * time.Minute
)

// Connect returns a connection pool for postgres conn
func Connect(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "open pgx connection")
	}

	db.SetMaxIdleConns(defaultMaxIdleConn)
	db.SetMaxOpenConns(defaultMaxOpenConn)
	db.SetConnMaxLifetime(defaultMaxConnLifetime)
	db.SetConnMaxIdleTime(defaultMaxIdleTime)

	return db, nil
}
