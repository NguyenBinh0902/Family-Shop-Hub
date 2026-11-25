package problem

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"familyshop-hub/internal/kernel"
	"familyshop-hub/internal/platform/requestid"
)

func init() { gin.SetMode(gin.TestMode) }

var errBR13 = &kernel.Error{Code: "BR-13", Kind: kernel.ErrorKindConflict, Message: "Không đủ tồn kho"}

func TestFromError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		fallback   int
		wantStatus int
		wantCode   string
		wantDetail string
	}{
		{"business conflict", fmt.Errorf("sell: %w", errBR13), 500, 409, "BR-13", "Không đủ tồn kho"},
		{"validation", kernel.NewValidationError(kernel.FieldError{Field: "name", Message: "Bắt buộc"}), 500, 422, "VALIDATION", "Dữ liệu không hợp lệ"},
		{"forbidden", &kernel.Error{Code: "FORBIDDEN", Kind: kernel.ErrorKindForbidden, Message: "x"}, 500, 403, "FORBIDDEN", "x"},
		{"not found kind", &kernel.Error{Code: "NOT_FOUND", Kind: kernel.ErrorKindNotFound, Message: "x"}, 500, 404, "NOT_FOUND", "x"},
		{"no rows", fmt.Errorf("load order: %w", pgx.ErrNoRows), 500, 404, "", detailNotFound},
		{"decoder 400 keeps status", errors.New("invalid character"), 400, 400, "", detailBadInput},
		{"internal error is hidden", errors.New("dial tcp: connection refused"), 500, 500, "", detailInternal},
		{"internal error with 200 fallback", errors.New("boom"), 200, 500, "", detailInternal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := FromError(tt.err, tt.fallback)
			if p.Status != tt.wantStatus || p.Code != tt.wantCode || p.Detail != tt.wantDetail {
				t.Errorf("got status=%d code=%q detail=%q, want %d %q %q",
					p.Status, p.Code, p.Detail, tt.wantStatus, tt.wantCode, tt.wantDetail)
			}
			if p.Type != "about:blank" || p.Title == "" {
				t.Errorf("type/title not set: %+v", p)
			}
		})
	}
}

func TestFromError_CarriesFieldErrors(t *testing.T) {
	err := kernel.NewValidationError(
		kernel.FieldError{Field: "name", Message: "Bắt buộc"},
		kernel.FieldError{Field: "sellPrice", Message: "Phải lớn hơn hoặc bằng 0"},
	)
	p := FromError(err, 500)
	if len(p.Errors) != 2 || p.Errors[1].Field != "sellPrice" {
		t.Errorf("errors = %+v", p.Errors)
	}
}

func TestMiddleware_RendersHandlerError(t *testing.T) {
	r := gin.New()
	r.Use(requestid.Middleware(), Middleware(slog.New(slog.NewTextHandler(io.Discard, nil))))
	r.GET("/business", func(c *gin.Context) { _ = c.Error(errBR13) })
	r.GET("/internal", func(c *gin.Context) {
		_ = c.Error(errors.New("secret detail"))
		c.Status(http.StatusInternalServerError)
	})
	r.GET("/written", func(c *gin.Context) {
		_ = c.Error(errBR13)
		c.String(http.StatusTeapot, "already written")
	})

	t.Run("business error", func(t *testing.T) {
		rec := serve(r, "/business")
		var p Problem
		decode(t, rec, &p)
		if rec.Code != http.StatusConflict || p.Code != "BR-13" || p.RequestID == "" {
			t.Errorf("status=%d problem=%+v", rec.Code, p)
		}
		if ct := rec.Header().Get("Content-Type"); ct != ContentType {
			t.Errorf("Content-Type = %q", ct)
		}
	})
	t.Run("internal error does not leak", func(t *testing.T) {
		rec := serve(r, "/internal")
		var p Problem
		decode(t, rec, &p)
		if rec.Code != http.StatusInternalServerError || p.Detail != detailInternal {
			t.Errorf("status=%d problem=%+v", rec.Code, p)
		}
	})
	t.Run("handler response is kept", func(t *testing.T) {
		rec := serve(r, "/written")
		if rec.Code != http.StatusTeapot || rec.Body.String() != "already written" {
			t.Errorf("status=%d body=%q", rec.Code, rec.Body.String())
		}
	})
}

func serve(h http.Handler, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), v); err != nil {
		t.Fatalf("decode %q: %v", rec.Body.String(), err)
	}
}
