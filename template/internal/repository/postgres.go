// Package repository is the persistence abstraction. Each resource gets
// an interface (e.g. user.go) and a concrete implementation named after
// its backing technology (e.g. user_postgres.go) — services only ever
// depend on the interface.
package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"{{ module_name }}/internal/apperror"
)

// NewPostgresPool connects to Postgres using the given DSN. Callers own
// the returned pool and must Close it on shutdown.
func NewPostgresPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

// mapPgError translates driver-level failures into apperror so the
// service layer never has to know it's talking to Postgres.
func mapPgError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.NotFound("resource not found")
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
		return apperror.Conflict("resource already exists")
	}

	return apperror.Internal("database error", err)
}
