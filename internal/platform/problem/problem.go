// Package problem turns errors into RFC 9457 Problem Details responses. It is
// the only place that decides how an error looks to the client: handlers just
// return errors.
package problem

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"familyshop-hub/internal/kernel"
	"familyshop-hub/internal/platform/requestid"
)

// ContentType is the media type of Problem Details responses.
const ContentType = "application/problem+json"

// Problem is the JSON body of an error response (RFC 9457) with two
// extensions: code (a business rule code such as "BR-15") and errors
// (per-field messages for forms).
type Problem struct {
	Type      string         `json:"type"`
	Title     string         `json:"title"`
	Status    int            `json:"status"`
	Detail    string         `json:"detail,omitempty"`
	Code      string         `json:"code,omitempty"`
	Errors    []FieldProblem `json:"errors,omitempty"`
	RequestID string         `json:"requestId,omitempty"`
}

// FieldProblem is a message about one input field.
type FieldProblem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Messages shown for errors that carry no business message of their own.
const (
	detailInternal = "Đã có lỗi xảy ra, vui lòng thử lại sau."
	detailNotFound = "Không tìm thấy dữ liệu."
	detailBadInput = "Yêu cầu không hợp lệ."
)

// New builds a Problem for a status code with a Vietnamese title.
func New(status int, detail string) Problem {
	return Problem{Type: "about:blank", Title: title(status), Status: status, Detail: detail}
}

// FromError maps err to a Problem:
//   - *kernel.Error → status from its Kind, with code, message and fields;
//   - pgx.ErrNoRows anywhere in the chain → 404;
//   - anything else → fallbackStatus if it is a 4xx (set by a request
//     decoder), otherwise 500 with a generic message so internals never leak.
func FromError(err error, fallbackStatus int) Problem {
	var be *kernel.Error
	switch {
	case errors.As(err, &be):
		p := New(statusOf(be.Kind), be.Message)
		p.Code = be.Code
		for _, f := range be.Fields {
			p.Errors = append(p.Errors, FieldProblem{Field: f.Field, Message: f.Message})
		}
		return p
	case errors.Is(err, pgx.ErrNoRows):
		return New(http.StatusNotFound, detailNotFound)
	case fallbackStatus >= 400 && fallbackStatus < 500:
		return New(fallbackStatus, detailBadInput)
	default:
		return New(http.StatusInternalServerError, detailInternal)
	}
}

// Write sends err as a Problem response. Server errors are logged with the
// full error; the client only sees the generic message.
func Write(c *gin.Context, logger *slog.Logger, err error, fallbackStatus int) {
	p := FromError(err, fallbackStatus)
	p.RequestID = requestid.From(c)
	if p.Status >= 500 {
		logger.ErrorContext(c.Request.Context(), "request failed",
			slog.String("request_id", p.RequestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Any("error", err),
		)
	}
	Render(c, p)
}

// Render writes p as the response body.
func Render(c *gin.Context, p Problem) {
	if p.RequestID == "" {
		p.RequestID = requestid.From(c)
	}
	c.Header("Content-Type", ContentType)
	c.AbortWithStatusJSON(p.Status, p)
}

// Middleware renders the last error a handler attached with c.Error, if the
// handler has not written a response itself. Generated OpenAPI handlers report
// errors this way, so they never have to format errors.
func Middleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}
		Write(c, logger, c.Errors.Last().Err, c.Writer.Status())
	}
}

func statusOf(k kernel.ErrorKind) int {
	switch k {
	case kernel.ErrorKindInvalid:
		return http.StatusUnprocessableEntity
	case kernel.ErrorKindConflict:
		return http.StatusConflict
	case kernel.ErrorKindForbidden:
		return http.StatusForbidden
	case kernel.ErrorKindNotFound:
		return http.StatusNotFound
	default:
		// A business error without a kind is a programming mistake; treat it
		// as invalid input rather than hiding the message behind a 500.
		return http.StatusUnprocessableEntity
	}
}

func title(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "Yêu cầu không hợp lệ"
	case http.StatusUnauthorized:
		return "Chưa đăng nhập"
	case http.StatusForbidden:
		return "Không có quyền"
	case http.StatusNotFound:
		return "Không tìm thấy"
	case http.StatusMethodNotAllowed:
		return "Phương thức không được hỗ trợ"
	case http.StatusConflict:
		return "Xung đột dữ liệu"
	case http.StatusUnprocessableEntity:
		return "Dữ liệu không hợp lệ"
	case http.StatusServiceUnavailable:
		return "Dịch vụ tạm thời không khả dụng"
	default:
		if status >= 500 {
			return "Lỗi máy chủ"
		}
		return http.StatusText(status)
	}
}
