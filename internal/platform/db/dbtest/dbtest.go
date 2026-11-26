// Package dbtest runs integration tests against a real PostgreSQL started with
// testcontainers-go.
//
// One container is started per test package, on first use. The migrations
// are applied once to a template database; every NewDB call then clones the
// template (CREATE DATABASE ... TEMPLATE), which takes milliseconds and gives
// each test an isolated, fully migrated database.
//
// Usage in a package's tests:
//
//	func TestMain(m *testing.M) { dbtest.Main(m) }
//
//	func TestSomething(t *testing.T) {
//		d := dbtest.NewDB(t)
//		// d.App connects as familyshop_app, so RLS applies.
//	}
package dbtest

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"familyshop-hub/db/migrations"
)

// Image is the PostgreSQL image used by tests; keep it in step with
// docker-compose.yml.
const Image = "postgres:19-alpine"

const (
	ownerUser    = "owner"
	ownerPass    = "owner"
	appUser      = "familyshop_app"
	appPass      = "app"
	templateName = "familyshop_template"
)

// DB is an isolated, migrated database for one test.
type DB struct {
	// App connects as the application role: not superuser, no BYPASSRLS.
	// Use it for everything the application itself would do.
	App *pgxpool.Pool
	// Owner connects as the owner (superuser in the test container). Use it
	// only for test setup that the application may not do, such as DDL or
	// seeding rows of another household.
	Owner *pgxpool.Pool
}

var (
	startOnce sync.Once
	startErr  error
	container *tcpostgres.PostgresContainer
	baseURL   *url.URL // owner connection to the "postgres" database
	dbCounter atomic.Int64
)

// Main runs the package's tests and then stops the container, if one was
// started. Call it from TestMain.
func Main(m *testing.M) {
	code := m.Run()
	if container != nil {
		if err := testcontainers.TerminateContainer(container); err != nil {
			fmt.Fprintln(os.Stderr, "dbtest: terminate container:", err)
		}
	}
	os.Exit(code)
}

// NewDB returns a fresh database cloned from the migrated template. It is
// dropped when the test ends. Tests that call it are skipped with -short.
func NewDB(t testing.TB) *DB {
	t.Helper()
	if testing.Short() {
		t.Skip("integration test: needs Docker (skipped with -short)")
	}
	startOnce.Do(func() { startErr = start() })
	if startErr != nil {
		t.Fatalf("dbtest: start postgres: %v", startErr)
	}

	ctx := context.Background()
	name := fmt.Sprintf("test_%d_%d", os.Getpid(), dbCounter.Add(1))

	admin := connect(t, ownerURL("postgres"))
	if _, err := admin.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", name, templateName)); err != nil {
		t.Fatalf("dbtest: create database: %v", err)
	}

	d := &DB{
		App:   connect(t, appURL(name)),
		Owner: connect(t, ownerURL(name)),
	}
	t.Cleanup(func() {
		d.App.Close()
		d.Owner.Close()
		// WITH (FORCE) disconnects anything a test leaked.
		if _, err := admin.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %s WITH (FORCE)", name)); err != nil {
			t.Errorf("dbtest: drop database: %v", err)
		}
		admin.Close()
	})
	return d
}

// start launches the container, creates the application role (done by
// deploy/postgres/init in docker-compose) and builds the template database.
func start() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	c, err := tcpostgres.Run(ctx, Image,
		tcpostgres.WithDatabase(templateName),
		tcpostgres.WithUsername(ownerUser),
		tcpostgres.WithPassword(ownerPass),
		tcpostgres.BasicWaitStrategies(),
	)
	if c != nil {
		container = c // set even on error so Main can clean up
	}
	if err != nil {
		return fmt.Errorf("run container: %w", err)
	}

	connStr, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return fmt.Errorf("connection string: %w", err)
	}
	if baseURL, err = url.Parse(connStr); err != nil {
		return fmt.Errorf("parse connection string: %w", err)
	}

	conn, err := pgx.Connect(ctx, ownerURL(templateName))
	if err != nil {
		return fmt.Errorf("connect to template: %w", err)
	}
	_, err = conn.Exec(ctx, fmt.Sprintf(
		"CREATE ROLE %s LOGIN PASSWORD '%s' NOSUPERUSER NOBYPASSRLS", appUser, appPass))
	_ = conn.Close(ctx)
	if err != nil {
		return fmt.Errorf("create app role: %w", err)
	}

	if err := migrate(ctx, ownerURL(templateName)); err != nil {
		return err
	}
	return nil
}

func migrate(ctx context.Context, dsn string) error {
	sqlDB, err := migrations.OpenDB(dsn)
	if err != nil {
		return err
	}
	// Closing matters: CREATE DATABASE ... TEMPLATE fails while anyone is
	// connected to the template.
	defer func() { _ = sqlDB.Close() }()

	p, err := migrations.NewProvider(sqlDB, goose.WithSlog(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		return err
	}
	if _, err := p.Up(ctx); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

func connect(t testing.TB, dsn string) *pgxpool.Pool {
	t.Helper()
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("dbtest: connect: %v", err)
	}
	return pool
}

func ownerURL(dbName string) string { return withDB(dbName, url.UserPassword(ownerUser, ownerPass)) }

func appURL(dbName string) string { return withDB(dbName, url.UserPassword(appUser, appPass)) }

func withDB(dbName string, user *url.Userinfo) string {
	u := *baseURL
	u.User = user
	u.Path = "/" + dbName
	return u.String()
}
