// Package health implements GET /healthz: it reports whether the API can
// reach and query the database.
package health

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"familyshop-hub/internal/platform/health/sqlcgen"
	"familyshop-hub/internal/platform/httpapi"
)

// timeout bounds the database check so a hung database makes the health check
// fail fast instead of hanging the caller (a load balancer, the web page).
const timeout = 2 * time.Second

// Pinger runs a trivial query against the database.
type Pinger interface {
	Ping(ctx context.Context) (int32, error)
}

// Handler serves the health endpoint.
type Handler struct {
	db     Pinger
	logger *slog.Logger
}

// NewHandler returns a Handler that checks the database through db. Pass
// sqlcgen.New(pool) in production.
func NewHandler(db Pinger, logger *slog.Logger) *Handler {
	return &Handler{db: db, logger: logger}
}

// NewPostgresHandler is a convenience for main.go: it builds the sqlc queries
// on top of a pool or connection.
func NewPostgresHandler(conn sqlcgen.DBTX, logger *slog.Logger) *Handler {
	return NewHandler(sqlcgen.New(conn), logger)
}

// GetHealthz implements httpapi.StrictServerInterface.
func (h *Handler) GetHealthz(ctx context.Context, _ httpapi.GetHealthzRequestObject) (httpapi.GetHealthzResponseObject, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if _, err := h.db.Ping(ctx); err != nil {
		// Unhealthy is an expected answer, not a failure of this handler, so
		// it is returned as a typed 503 response rather than as an error.
		h.logger.WarnContext(ctx, "health check: database unreachable", slog.Any("error", err))
		detail := "Không kết nối được cơ sở dữ liệu."
		return httpapi.GetHealthz503ApplicationProblemPlusJSONResponse{
			ServiceUnavailableApplicationProblemPlusJSONResponse: httpapi.ServiceUnavailableApplicationProblemPlusJSONResponse{
				Type:   "about:blank",
				Title:  "Dịch vụ tạm thời không khả dụng",
				Status: http.StatusServiceUnavailable,
				Detail: &detail,
			},
		}, nil
	}
	return httpapi.GetHealthz200JSONResponse{Status: httpapi.HealthStatusOk, Database: httpapi.HealthDatabaseOk}, nil
}
