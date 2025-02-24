package instrumentpg

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/viebiz/lit/postgres"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type instrumentedDB struct {
	postgres.ContextExecutor
}

func (db instrumentedDB) PrepareContext(ctx context.Context, query string) (stmt *sql.Stmt, err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		span.AddEvent("PrepareContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		db.recordError(span, err)
	}(time.Now())

	return db.ContextExecutor.PrepareContext(ctx, query)
}

func (db instrumentedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (rs sql.Result, err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		span.AddEvent("ExecContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		db.recordError(span, err)
	}(time.Now())

	return db.ContextExecutor.ExecContext(ctx, query, args...)
}

func (db instrumentedDB) QueryContext(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		span.AddEvent("QueryContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		db.recordError(span, err)
	}(time.Now())

	return db.ContextExecutor.QueryContext(ctx, query, args...)
}

func (db instrumentedDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) (row *sql.Row) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		span.AddEvent("QueryRowContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		db.recordError(span, row.Err())
	}(time.Now())

	return db.ContextExecutor.QueryRowContext(ctx, query, args...)
}

func (db instrumentedDB) recordError(span trace.Span, err error) {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			span.AddEvent("Database Error", trace.WithAttributes(
				attribute.String("Error", err.Error()),
				attribute.String("Code", pgErr.Code),
				attribute.String("Severity", pgErr.Severity),
				attribute.String("Message", pgErr.Message),
				attribute.String("Detail", pgErr.Detail),
			))
		} else {
			span.AddEvent("Database Error", trace.WithAttributes(
				attribute.String("Error", err.Error()),
			))
		}
	}
}
