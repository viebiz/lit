package postgres

import (
	"context"
	"database/sql"
)

// Executor can perform SQL queries.
type Executor interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// ContextExecutor can perform SQL queries with context
type ContextExecutor interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	Executor
}

// Transactor can commit and rollback, on top of being able to execute queries.
type Transactor interface {
	Commit() error
	Rollback() error

	Executor
}

// Beginner begins transactions.
type Beginner interface {
	Begin() (*sql.Tx, error)
}

// ContextTransactor can commit and rollback, on top of being able to execute
// context-aware queries.
type ContextTransactor interface {
	Commit() error
	Rollback() error

	ContextExecutor
}

// ContextBeginner allows creation of context aware transactions with options.
type ContextBeginner interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)

	ContextExecutor
}
