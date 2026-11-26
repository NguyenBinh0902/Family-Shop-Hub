// Package requestid gives every HTTP request an ID that appears in logs and in
// error responses, so a user's error report can be matched to the server log.
package requestid

import (
	"github.com/gin-gonic/gin"

	"familyshop-hub/internal/kernel"
)

// Header is the HTTP header that carries the request ID.
const Header = "X-Request-ID"

// ginKey is unexported so no other package can overwrite the value.
const ginKey = "familyshop.request_id"

// maxLen bounds IDs supplied by clients so they cannot flood the logs.
const maxLen = 64

// Middleware reuses a well-formed X-Request-ID from the client (useful when a
// proxy already assigned one) or generates a new one, and echoes it back.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(Header)
		if !valid(id) {
			id = kernel.NewID().String()
		}
		c.Set(ginKey, id)
		c.Header(Header, id)
		c.Next()
	}
}

// From returns the request ID of c, or "" if the middleware did not run.
func From(c *gin.Context) string {
	return c.GetString(ginKey)
}

// valid accepts short IDs made of characters that are safe to log.
func valid(id string) bool {
	if id == "" || len(id) > maxLen {
		return false
	}
	for _, r := range id {
		ok := r == '-' || r == '_' || r == '.' ||
			(r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
		if !ok {
			return false
		}
	}
	return true
}
