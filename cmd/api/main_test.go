package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"familyshop-hub/internal/platform/db/dbtest"
	"familyshop-hub/internal/platform/httpserver"
	"familyshop-hub/internal/platform/problem"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	dbtest.Main(m)
}

func newTestEngine(pool *pgxpool.Pool) http.Handler {
	logger := slog.New(slog.DiscardHandler)
	return httpserver.NewEngine(logger, routes(newServer(pool, logger), logger))
}

func get(h http.Handler, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestHealthz_OK(t *testing.T) {
	d := dbtest.NewDB(t)

	rec := get(newTestEngine(d.App), "/api/v1/healthz")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body)
	}
	var body struct{ Status, Database string }
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Status != "ok" || body.Database != "ok" {
		t.Errorf("body = %+v", body)
	}
}

func TestHealthz_DBDown_Returns503Problem(t *testing.T) {
	// A pool pointing at a port nobody listens on. pgxpool connects lazily,
	// so creating it succeeds and the failure shows up on the first query.
	pool, err := pgxpool.New(context.Background(), "postgres://u:p@127.0.0.1:1/x?connect_timeout=1")
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	start := time.Now()
	rec := get(newTestEngine(pool), "/api/v1/healthz")

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body)
	}
	if ct := rec.Header().Get("Content-Type"); ct != problem.ContentType {
		t.Errorf("Content-Type = %q", ct)
	}
	var p problem.Problem
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil || p.Status != 503 || p.Detail == "" {
		t.Errorf("problem = %+v, err %v", p, err)
	}
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Errorf("health check took %s, want it bounded by the 2s timeout", elapsed)
	}
}

func TestUnknownAPIRoute_Returns404Problem(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), "postgres://u:p@127.0.0.1:1/x")
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	rec := get(newTestEngine(pool), "/api/v1/nope")
	if rec.Code != http.StatusNotFound || rec.Header().Get("Content-Type") != problem.ContentType {
		t.Errorf("status %d, content-type %q", rec.Code, rec.Header().Get("Content-Type"))
	}
}
