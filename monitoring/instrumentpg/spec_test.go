package instrumentpg

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
	"github.com/viebiz/lit/postgres"
	"go.opentelemetry.io/otel/trace"
)

func TestInstrumentedDB_PrepareContext(t *testing.T) {
	ctx := context.Background()
	be := postgres.NewMockBeginnerExecutor(t)
	be.EXPECT().PrepareContext(mock.Anything, mock.Anything).Return((*sql.Stmt)(nil), nil)

	db := instrumentedDB{BeginnerExecutor: be}
	_, _ = db.PrepareContext(ctx, "SELECT 1")
}

func TestInstrumentedDB_ExecContext(t *testing.T) {
	ctx := context.Background()
	be := postgres.NewMockBeginnerExecutor(t)
	be.EXPECT().ExecContext(mock.Anything, mock.Anything).Return((sql.Result)(nil), nil)

	db := instrumentedDB{BeginnerExecutor: be}
	_, _ = db.ExecContext(ctx, "UPDATE t SET a=1")
}

func TestInstrumentedDB_QueryContext(t *testing.T) {
	ctx := context.Background()
	be := postgres.NewMockBeginnerExecutor(t)
	be.EXPECT().QueryContext(mock.Anything, mock.Anything).Return((*sql.Rows)(nil), nil)

	db := instrumentedDB{BeginnerExecutor: be}
	_, _ = db.QueryContext(ctx, "SELECT 1")
}

func TestInstrumentedDB_QueryRowContext(t *testing.T) {
	ctx := context.Background()
	be := postgres.NewMockBeginnerExecutor(t)
	be.EXPECT().QueryRowContext(mock.Anything, mock.Anything).Return(&sql.Row{})

	db := instrumentedDB{BeginnerExecutor: be}
	_ = db.QueryRowContext(ctx, "SELECT 1")
}

func TestInstrumentedTx_PrepareContext(t *testing.T) {
	ctx := context.Background()
	ce := postgres.NewMockContextExecutor(t)
	ce.EXPECT().PrepareContext(mock.Anything, mock.Anything).Return((*sql.Stmt)(nil), nil)

	tx := instrumentedTx{ContextExecutor: ce}
	_, _ = tx.PrepareContext(ctx, "SELECT 1")
}

func TestInstrumentedTx_ExecContext(t *testing.T) {
	ctx := context.Background()
	ce := postgres.NewMockContextExecutor(t)
	ce.EXPECT().ExecContext(mock.Anything, mock.Anything).Return((sql.Result)(nil), nil)

	tx := instrumentedTx{ContextExecutor: ce}
	_, _ = tx.ExecContext(ctx, "DELETE FROM t")
}

func TestInstrumentedTx_QueryContext(t *testing.T) {
	ctx := context.Background()
	ce := postgres.NewMockContextExecutor(t)
	ce.EXPECT().QueryContext(mock.Anything, mock.Anything).Return((*sql.Rows)(nil), nil)

	tx := instrumentedTx{ContextExecutor: ce}
	_, _ = tx.QueryContext(ctx, "SELECT 1")
}

func TestInstrumentedTx_QueryRowContext(t *testing.T) {
	ctx := context.Background()
	ce := postgres.NewMockContextExecutor(t)
	ce.EXPECT().QueryRowContext(mock.Anything, mock.Anything).Return(&sql.Row{})

	tx := instrumentedTx{ContextExecutor: ce}
	_ = tx.QueryRowContext(ctx, "SELECT 1")
}

func TestInstrumentedDB_recordError(t *testing.T) {
	be := postgres.NewMockBeginnerExecutor(t)
	db := instrumentedDB{BeginnerExecutor: be}
	span := trace.SpanFromContext(context.Background())

	db.recordError(span, errors.New("generic error"))
	db.recordError(span, &pgconn.PgError{
		Code:     "XX000",
		Severity: "ERROR",
		Message:  "db error",
		Detail:   "detail",
	})
}

func TestInstrumentedTx_recordError(t *testing.T) {
	ce := postgres.NewMockContextExecutor(t)
	tx := instrumentedTx{ContextExecutor: ce}
	span := trace.SpanFromContext(context.Background())

	tx.recordError(span, errors.New("generic error"))
	tx.recordError(span, &pgconn.PgError{
		Code:     "XX001",
		Severity: "ERROR",
		Message:  "tx error",
		Detail:   "detail",
	})
}
