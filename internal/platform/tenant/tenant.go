// Package tenant defines the scope of a request: which household and store it
// acts on, and who is acting with which role. M1 adds the middleware that
// resolves the scope from the URL and the user's memberships.
package tenant

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Role is a member's role (roles-permissions.md). The zero value is invalid.
type Role int

// Roles.
const (
	RoleUnknown Role = iota
	RoleOwner        // Chủ hộ: household level, every store
	RoleManager      // Quản lý: one store
	RoleStaff        // Nhân viên: one store
)

var roleNames = map[Role]string{
	RoleOwner:   "OWNER",
	RoleManager: "MANAGER",
	RoleStaff:   "STAFF",
}

// String returns the code used in the database and API ("OWNER").
func (r Role) String() string {
	if s, ok := roleNames[r]; ok {
		return s
	}
	return "UNKNOWN"
}

// ParseRole converts a role code back to a Role; ok is false if s is unknown.
func ParseRole(s string) (r Role, ok bool) {
	for role, name := range roleNames {
		if name == s {
			return role, true
		}
	}
	return RoleUnknown, false
}

// Scope is the tenant context of a request. Use cases take it as an explicit
// parameter rather than reading it from context.Context, so every function
// that touches tenant data shows that in its signature.
type Scope struct {
	HouseholdID uuid.UUID
	StoreID     uuid.UUID // uuid.Nil for household-level requests
	UserID      uuid.UUID
	Role        Role
}

// ginKey is unexported so only this package can set the scope.
const ginKey = "familyshop.tenant_scope"

// SetScope stores s on the request. Only the tenant middleware calls it.
func SetScope(c *gin.Context, s Scope) { c.Set(ginKey, s) }

// ScopeFrom returns the scope of the request; ok is false if the tenant
// middleware did not run, which handlers must treat as an internal error.
func ScopeFrom(c *gin.Context) (s Scope, ok bool) {
	v, exists := c.Get(ginKey)
	if !exists {
		return Scope{}, false
	}
	s, ok = v.(Scope)
	return s, ok
}
