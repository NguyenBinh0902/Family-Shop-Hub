// Command migrate applies the database migrations embedded in the binary.
//
//	migrate up            apply all pending migrations
//	migrate down          roll back the latest migration
//	migrate status        list migrations and whether they are applied
//	migrate create NAME   write db/migrations/NNNN_NAME.sql (no database needed)
//
// It connects with MIGRATION_DATABASE_URL, the owner role that may run DDL.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pressly/goose/v3"

	"familyshop-hub/db/migrations"
	"familyshop-hub/internal/platform/config"
)

const usage = "usage: migrate up | down | status | create NAME"

// migrationsDir is where `create` writes new files, relative to the repo root.
const migrationsDir = "db/migrations"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, os.Args[1:], os.LookupEnv, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "migrate:", err)
		os.Exit(1)
	}
}

// run holds the program logic; main only wires the real environment, which
// keeps run testable and lets deferred calls run before os.Exit.
func run(ctx context.Context, args []string, lookup config.LookupFunc, out io.Writer) error {
	if len(args) == 0 {
		return errors.New(usage)
	}

	if args[0] == "create" {
		if len(args) != 2 {
			return errors.New(usage)
		}
		path, err := migrations.Create(migrationsDir, args[1])
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, "created", path)
		return err
	}

	cfg, err := config.LoadMigration(lookup)
	if err != nil {
		return err
	}
	db, err := migrations.OpenDB(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	logger := slog.New(slog.NewTextHandler(out, nil))
	p, err := migrations.NewProvider(db, goose.WithSlog(logger))
	if err != nil {
		return err
	}

	switch args[0] {
	case "up":
		results, err := p.Up(ctx)
		if err != nil {
			return fmt.Errorf("up: %w", err)
		}
		_, err = fmt.Fprintf(out, "applied %d migration(s)\n", len(results))
		return err
	case "down":
		if _, err := p.Down(ctx); err != nil {
			return fmt.Errorf("down: %w", err)
		}
		return nil
	case "status":
		return printStatus(ctx, out, p)
	default:
		return errors.New(usage)
	}
}

func printStatus(ctx context.Context, out io.Writer, p *goose.Provider) error {
	statuses, err := p.Status(ctx)
	if err != nil {
		return fmt.Errorf("status: %w", err)
	}
	for _, s := range statuses {
		applied := "pending"
		if s.State == goose.StateApplied {
			applied = "applied " + s.AppliedAt.Format("2006-01-02 15:04:05")
		}
		if _, err := fmt.Fprintf(out, "%-40s %s\n", s.Source.Path, applied); err != nil {
			return err
		}
	}
	return nil
}
