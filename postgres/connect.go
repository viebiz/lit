package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
	pkgerrors "github.com/pkg/errors"
)

const (
	defaultMaxIdleConn = 2
	defaultMaxOpenConn = 20
)

// Connect returns the singleton instance of the database
func Connect(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "opening postgres connection")
	}

	db.SetMaxIdleConns(defaultMaxIdleConn)
	db.SetMaxOpenConns(defaultMaxOpenConn)

	return db, nil
}
