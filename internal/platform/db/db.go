// Package db owns the PostgreSQL connection pool and transactions.
//
// Every tenant transaction sets app.household_id so Row-Level Security acts
// as a second barrier behind the household filter in each query (NFR-05).
// Repositories never open transactions themselves: a use case opens one
// through TxRunner and receives the repositories bound to it (ADR-0004).
package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool connects to url and checks the connection, so a wrong URL or a
// database that is down is reported at startup.
func NewPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}

// ErrNoHousehold reports a tenant transaction started without a household.
// It is a programming error: tenant data must always be scoped.
var ErrNoHousehold = errors.New("db: tenant transaction without household id")

// TxRunner runs a function inside a transaction and hands it repositories of
// type R built on that transaction. R is a struct of repository interfaces
// declared by the calling context's app package, and newRepos (wired in
// main.go) builds it from the transaction. A use case that touches several
// contexts therefore gets all of them in one atomic transaction, with nothing
// hidden in context.Context.
//
// A context's app package declares the interface it needs, for example
//
//	type TxRunner interface {
//		InTx(ctx context.Context, householdID uuid.UUID, fn func(context.Context, Repos) error) error
//	}
//
// and *db.TxRunner[Repos] satisfies it without naming it.
type TxRunner[R any] struct {
	pool     *pgxpool.Pool
	newRepos func(pgx.Tx) R
}

// NewTxRunner returns a TxRunner that builds repositories with newRepos.
func NewTxRunner[R any](pool *pgxpool.Pool, newRepos func(pgx.Tx) R) *TxRunner[R] {
	return &TxRunner[R]{pool: pool, newRepos: newRepos}
}

// InTx runs fn in a read-write transaction scoped to householdID. The
// transaction commits if fn returns nil and rolls back otherwise.
func (r *TxRunner[R]) InTx(ctx context.Context, householdID uuid.UUID, fn func(context.Context, R) error) error {
	if householdID == uuid.Nil {
		return ErrNoHousehold
	}
	return r.run(ctx, pgx.TxOptions{}, householdID, fn)
}

// InReadTx runs fn in a read-only transaction scoped to householdID. Reads go
// through a transaction too, because the RLS setting only lives inside one.
func (r *TxRunner[R]) InReadTx(ctx context.Context, householdID uuid.UUID, fn func(context.Context, R) error) error {
	if householdID == uuid.Nil {
		return ErrNoHousehold
	}
	return r.run(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly}, householdID, fn)
}

// InSystemTx runs fn without a household, for tables outside any tenant such
// as users and sessions. RLS-protected tables return no rows here.
func (r *TxRunner[R]) InSystemTx(ctx context.Context, fn func(context.Context, R) error) error {
	return r.run(ctx, pgx.TxOptions{}, uuid.Nil, fn)
}

// run opens the transaction; householdID may be uuid.Nil only for InSystemTx.
func (r *TxRunner[R]) run(ctx context.Context, opts pgx.TxOptions, householdID uuid.UUID, fn func(context.Context, R) error) error {
	return pgx.BeginTxFunc(ctx, r.pool, opts, func(tx pgx.Tx) error {
		if householdID != uuid.Nil {
			if err := SetHousehold(ctx, tx, householdID); err != nil {
				return err
			}
		}
		return fn(ctx, r.newRepos(tx))
	})
}

// SetHousehold sets app.household_id for the rest of tx only (the third
// argument true makes it transaction-local, like SET LOCAL). Unlike SET,
// set_config accepts a bind parameter, so no SQL is built from strings.
func SetHousehold(ctx context.Context, tx pgx.Tx, householdID uuid.UUID) error {
	if _, err := tx.Exec(ctx, "SELECT set_config('app.household_id', $1, true)", householdID.String()); err != nil {
		return fmt.Errorf("set app.household_id: %w", err)
	}
	return nil
}
