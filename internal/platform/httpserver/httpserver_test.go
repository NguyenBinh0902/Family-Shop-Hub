package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"familyshop-hub/internal/platform/problem"
	"familyshop-hub/internal/platform/requestid"
)

func init() { gin.SetMode(gin.TestMode) }

func discardLogger() *slog.Logger { return slog.New(slog.DiscardHandler) }

func TestEngine_RequestID(t *testing.T) {
	r := NewEngine(discardLogger(), func(r *gin.Engine) {
		r.GET("/ok", func(c *gin.Context) { c.String(http.StatusOK, requestid.From(c)) })
	})

	t.Run("generated when missing", func(t *testing.T) {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ok", nil))
		id := rec.Header().Get(requestid.Header)
		if id == "" || rec.Body.String() != id {
			t.Errorf("header %q, body %q", id, rec.Body.String())
		}
	})
	t.Run("client value reused", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set(requestid.Header, "abc-123")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if got := rec.Header().Get(requestid.Header); got != "abc-123" {
			t.Errorf("request id = %q, want abc-123", got)
		}
	})
	t.Run("unsafe client value replaced", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set(requestid.Header, "bad id\nwith newline")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if got := rec.Header().Get(requestid.Header); strings.ContainsAny(got, " \n") || got == "" {
			t.Errorf("request id = %q, want a generated one", got)
		}
	})
}

func TestEngine_PanicBecomesProblem500(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	r := NewEngine(logger, func(r *gin.Engine) {
		r.GET("/panic", func(*gin.Context) { panic("bug") })
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/panic", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	var p problem.Problem
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil || p.Status != 500 || p.RequestID == "" {
		t.Errorf("body %q, err %v", rec.Body.String(), err)
	}
	if strings.Contains(rec.Body.String(), "bug") {
		t.Error("panic value leaked to the client")
	}
	if !strings.Contains(logs.String(), "panic recovered") {
		t.Error("panic was not logged")
	}
}

func TestEngine_UnknownRouteAndMethod(t *testing.T) {
	r := NewEngine(discardLogger(), func(r *gin.Engine) {
		r.GET("/only-get", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	})

	tests := []struct {
		method, path string
		want         int
	}{
		{http.MethodGet, "/nope", http.StatusNotFound},
		{http.MethodPost, "/only-get", http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))
		if rec.Code != tt.want || rec.Header().Get("Content-Type") != problem.ContentType {
			t.Errorf("%s %s: status %d content-type %q", tt.method, tt.path, rec.Code, rec.Header().Get("Content-Type"))
		}
	}
}

func TestServer_GracefulShutdownWaitsForInFlightRequest(t *testing.T) {
	started := make(chan struct{})
	srv := New(Options{
		Logger: discardLogger(),
		Register: func(r *gin.Engine) {
			r.GET("/slow", func(c *gin.Context) {
				close(started)
				time.Sleep(200 * time.Millisecond)
				c.String(http.StatusOK, "done")
			})
		},
	})
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	runErr := make(chan error, 1)
	go func() { runErr <- srv.Serve(ctx, ln, 2*time.Second) }()

	type result struct {
		body string
		err  error
	}
	got := make(chan result, 1)
	go func() {
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://"+ln.Addr().String()+"/slow", nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			got <- result{err: err}
			return
		}
		defer func() { _ = resp.Body.Close() }()
		b, err := io.ReadAll(resp.Body)
		got <- result{body: string(b), err: err}
	}()

	<-started
	cancel() // shutdown begins while /slow is still running

	res := <-got
	if res.err != nil || res.body != "done" {
		t.Errorf("in-flight request: body %q, err %v; want done", res.body, res.err)
	}
	if err := <-runErr; err != nil {
		t.Errorf("Serve returned %v", err)
	}
}
