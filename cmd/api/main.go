// Command api runs the FamilyShop Hub HTTP API.
//
// This file is the composition root: the only place that knows how config,
// database, contexts and HTTP fit together. Everything is wired by hand with
// constructors; there is no DI framework and no global state.
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"familyshop-hub/internal/platform/config"
	"familyshop-hub/internal/platform/db"
	"familyshop-hub/internal/platform/health"
	"familyshop-hub/internal/platform/httpapi"
	"familyshop-hub/internal/platform/httpserver"
)

func main() {
	// Cancelled on Ctrl+C or SIGTERM (docker stop), which starts the
	// graceful shutdown in httpserver.Run.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, os.LookupEnv, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "api:", err)
		os.Exit(1)
	}
}

// run holds the program so that deferred cleanups (closing the pool) run
// before main calls os.Exit, and so tests could drive it.
func run(ctx context.Context, lookup config.LookupFunc, logOut io.Writer) error {
	cfg, err := config.LoadAPI(lookup)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(logOut, &slog.HandlerOptions{Level: cfg.LogLevel}))
	if !cfg.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	srv := httpserver.New(httpserver.Options{
		Addr:     cfg.HTTPAddr,
		Logger:   logger,
		Register: routes(newServer(pool, logger), logger),
	})
	return srv.Run(ctx, cfg.ShutdownTimeout)
}

// server implements httpapi.StrictServerInterface by embedding one handler
// per context. Embedding promotes each handler's methods, so together they
// satisfy the single generated interface while each lives in its own package.
// When a milestone adds endpoints, its handler is added here.
type server struct {
	*health.Handler
}

// Compile-time check: the build fails here if an operation in openapi.yaml
// has no handler yet.
var _ httpapi.StrictServerInterface = server{}

func newServer(pool *pgxpool.Pool, logger *slog.Logger) server {
	return server{
		Handler: health.NewPostgresHandler(pool, logger),
	}
}

func routes(s server, logger *slog.Logger) func(*gin.Engine) {
	return func(r *gin.Engine) {
		httpapi.Register(r, s, logger)
	}
}
