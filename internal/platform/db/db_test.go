package db_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"familyshop-hub/internal/kernel"
	"familyshop-hub/internal/platform/db"
	"familyshop-hub/internal/platform/db/dbtest"
)

func TestMain(m *testing.M) { dbtest.Main(m) }

// repos is the simplest possible Repos type: the transaction itself.
func newRunner(d *dbtest.DB) *db.TxRunner[pgx.Tx] {
	return db.NewTxRunner(d.App, func(tx pgx.Tx) pgx.Tx { return tx })
}

func currentHousehold(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}) (*uuid.UUID, error) {
	var id *uuid.UUID
	err := q.QueryRow(ctx, "SELECT app_current_household()").Scan(&id)
	return id, err
}

func TestTxRunner_SetsHouseholdOnlyInsideTransaction(t *testing.T) {
	d := dbtest.NewDB(t)
	ctx := context.Background()
	runner := newRunner(d)
	household := kernel.NewID()

	err := runner.InTx(ctx, household, func(ctx context.Context, tx pgx.Tx) error {
		got, err := currentHousehold(ctx, tx)
		if err != nil {
			return err
		}
		if got == nil || *got != household {
			t.Errorf("inside tx household = %v, want %s", got, household)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("InTx: %v", err)
	}

	// The setting is transaction-local: a pooled connection reused afterwards
	// must not carry the previous request's household.
	for range 5 {
		got, err := currentHousehold(ctx, d.App)
		if err != nil {
			t.Fatal(err)
		}
		if got != nil {
			t.Fatalf("household leaked outside the transaction: %s", got)
		}
	}
}

func TestTxRunner_RequiresHousehold(t *testing.T) {
	d := dbtest.NewDB(t)
	runner := newRunner(d)
	noop := func(context.Context, pgx.Tx) error { return nil }

	if err := runner.InTx(context.Background(), uuid.Nil, noop); !errors.Is(err, db.ErrNoHousehold) {
		t.Errorf("InTx(uuid.Nil) = %v, want ErrNoHousehold", err)
	}
	if err := runner.InReadTx(context.Background(), uuid.Nil, noop); !errors.Is(err, db.ErrNoHousehold) {
		t.Errorf("InReadTx(uuid.Nil) = %v, want ErrNoHousehold", err)
	}
}

func TestTxRunner_RollsBackOnError(t *testing.T) {
	d := dbtest.NewDB(t)
	ctx := context.Background()
	mustExec(t, d, "CREATE TABLE items (n int)")
	runner := newRunner(d)

	errBoom := errors.New("boom")
	err := runner.InTx(ctx, kernel.NewID(), func(ctx context.Context, tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "INSERT INTO items VALUES (1)"); err != nil {
			return err
		}
		return errBoom
	})
	if !errors.Is(err, errBoom) {
		t.Fatalf("InTx error = %v, want the callback's error", err)
	}
	if n := count(t, d, "items"); n != 0 {
		t.Errorf("rows after rollback = %d, want 0", n)
	}
}

func TestTxRunner_InReadTxRejectsWrites(t *testing.T) {
	d := dbtest.NewDB(t)
	mustExec(t, d, "CREATE TABLE items (n int)")

	err := newRunner(d).InReadTx(context.Background(), kernel.NewID(), func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "INSERT INTO items VALUES (1)")
		return err
	})
	if err == nil {
		t.Error("write inside a read-only transaction succeeded")
	}
}

// TestRLS_BR01_AppRoleCannotSeeOtherHousehold proves the second barrier of
// BR-01/NFR-05: even a query with no household filter at all only returns the
// current household's rows when it runs as the application role.
func TestRLS_BR01_AppRoleCannotSeeOtherHousehold(t *testing.T) {
	d := dbtest.NewDB(t)
	ctx := context.Background()
	mustExec(t, d, `
		CREATE TABLE notes (household_id uuid NOT NULL, body text NOT NULL);
		ALTER TABLE notes ENABLE ROW LEVEL SECURITY;
		CREATE POLICY tenant ON notes
			USING (household_id = app_current_household())
			WITH CHECK (household_id = app_current_household());`)

	householdA, householdB := kernel.NewID(), kernel.NewID()
	mustExec(t, d, "INSERT INTO notes VALUES ($1, 'a1'), ($1, 'a2'), ($2, 'b1')", householdA, householdB)
	runner := newRunner(d)

	t.Run("app role is not superuser and cannot bypass RLS", func(t *testing.T) {
		var super, bypass bool
		err := d.App.QueryRow(ctx,
			"SELECT rolsuper, rolbypassrls FROM pg_roles WHERE rolname = current_user").Scan(&super, &bypass)
		if err != nil {
			t.Fatal(err)
		}
		if super || bypass {
			t.Fatalf("application role has superuser=%v bypassrls=%v", super, bypass)
		}
	})

	t.Run("unfiltered select sees only the current household", func(t *testing.T) {
		err := runner.InReadTx(ctx, householdA, func(ctx context.Context, tx pgx.Tx) error {
			rows, err := tx.Query(ctx, "SELECT body FROM notes ORDER BY body")
			if err != nil {
				return err
			}
			bodies, err := pgx.CollectRows(rows, pgx.RowTo[string])
			if err != nil {
				return err
			}
			if len(bodies) != 2 || bodies[0] != "a1" || bodies[1] != "a2" {
				t.Errorf("household A sees %v, want [a1 a2]", bodies)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("no household means no rows", func(t *testing.T) {
		err := runner.InSystemTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
			var n int
			if err := tx.QueryRow(ctx, "SELECT count(*) FROM notes").Scan(&n); err != nil {
				return err
			}
			if n != 0 {
				t.Errorf("system transaction sees %d tenant rows, want 0", n)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("cannot write into another household", func(t *testing.T) {
		err := runner.InTx(ctx, householdA, func(ctx context.Context, tx pgx.Tx) error {
			_, err := tx.Exec(ctx, "INSERT INTO notes VALUES ($1, 'sneaky')", householdB)
			return err
		})
		if err == nil {
			t.Error("inserting a row for household B while scoped to A succeeded")
		}
	})
}

func mustExec(t *testing.T, d *dbtest.DB, sql string, args ...any) {
	t.Helper()
	if _, err := d.Owner.Exec(context.Background(), sql, args...); err != nil {
		t.Fatalf("exec %q: %v", sql, err)
	}
}

func count(t *testing.T, d *dbtest.DB, table string) int {
	t.Helper()
	var n int
	if err := d.Owner.QueryRow(context.Background(), "SELECT count(*) FROM "+table).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}
