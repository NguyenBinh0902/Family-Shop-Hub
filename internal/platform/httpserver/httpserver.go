// Package httpserver builds the Gin engine with the middleware every request
// goes through (request ID, access log, panic recovery, problem rendering) and
// runs it with graceful shutdown.
package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"

	"familyshop-hub/internal/platform/problem"
	"familyshop-hub/internal/platform/requestid"
)

// Options configures a Server.
type Options struct {
	Addr   string
	Logger *slog.Logger
	// Register adds the application's routes to the engine.
	Register func(r *gin.Engine)
}

// Server is the HTTP server of the API.
type Server struct {
	srv    *http.Server
	logger *slog.Logger
}

// New builds the engine and wraps it in an http.Server with conservative
// timeouts, so slow clients cannot hold connections forever.
func New(opts Options) *Server {
	return &Server{
		srv: &http.Server{
			Addr:              opts.Addr,
			Handler:           NewEngine(opts.Logger, opts.Register),
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       120 * time.Second,
		},
		logger: opts.Logger,
	}
}

// NewEngine returns the Gin engine with the shared middleware and routes.
// It is exported so tests can exercise it with httptest without a socket.
func NewEngine(logger *slog.Logger, register func(r *gin.Engine)) *gin.Engine {
	r := gin.New()
	r.HandleMethodNotAllowed = true
	// Middleware order matters: the request ID must exist before anything
	// logs, and recovery sits inside the logger so a panic is still logged
	// with its final 500 status.
	r.Use(requestid.Middleware(), accessLog(logger), recovery(logger), problem.Middleware(logger))

	r.NoRoute(func(c *gin.Context) {
		problem.Render(c, problem.New(http.StatusNotFound, "Không tìm thấy đường dẫn."))
	})
	r.NoMethod(func(c *gin.Context) {
		problem.Render(c, problem.New(http.StatusMethodNotAllowed, "Phương thức không được hỗ trợ."))
	})

	if register != nil {
		register(r)
	}
	return r
}

// Run listens on the configured address and serves until ctx is cancelled.
// See Serve for the shutdown behaviour.
func (s *Server) Run(ctx context.Context, shutdownTimeout time.Duration) error {
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", s.srv.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", s.srv.Addr, err)
	}
	return s.Serve(ctx, ln, shutdownTimeout)
}

// Serve serves on ln until ctx is cancelled, then stops accepting connections
// and waits up to shutdownTimeout for in-flight requests to finish. Taking a
// listener lets tests use a random free port.
func (s *Server) Serve(ctx context.Context, ln net.Listener, shutdownTimeout time.Duration) error {
	s.logger.InfoContext(ctx, "http server listening", slog.String("addr", ln.Addr().String()))

	serveErr := make(chan error, 1)
	go func() { serveErr <- s.srv.Serve(ln) }()

	select {
	case err := <-serveErr:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}

	s.logger.InfoContext(ctx, "http server shutting down", slog.Duration("timeout", shutdownTimeout))
	// The parent ctx is already cancelled, so shutdown needs a fresh one.
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	if err := <-serveErr; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve: %w", err)
	}
	s.logger.InfoContext(shutdownCtx, "http server stopped")
	return nil
}

// accessLog writes one structured line per request.
func accessLog(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		}
		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		logger.LogAttrs(c.Request.Context(), level, "http request",
			slog.String("request_id", requestid.From(c)),
			slog.String("method", c.Request.Method),
			slog.String("route", route),
			slog.Int("status", status),
			slog.Duration("duration", time.Since(start)),
		)
	}
}

// recovery turns a panic into a logged 500 problem response instead of
// crashing the process. Handlers must not panic on purpose; this is a safety
// net for bugs.
func recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if rec == http.ErrAbortHandler { //nolint:errorlint // sentinel compared by identity, as net/http does
				panic(rec) // let net/http abort the connection as intended
			}
			logger.ErrorContext(c.Request.Context(), "panic recovered",
				slog.String("request_id", requestid.From(c)),
				slog.Any("panic", rec),
				slog.String("stack", string(debug.Stack())),
			)
			if !c.Writer.Written() {
				problem.Render(c, problem.New(http.StatusInternalServerError, "Đã có lỗi xảy ra, vui lòng thử lại sau."))
			}
			c.Abort()
		}()
		c.Next()
	}
}
