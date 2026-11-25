package tenant

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestRole_StringAndParse(t *testing.T) {
	for _, r := range []Role{RoleOwner, RoleManager, RoleStaff} {
		got, ok := ParseRole(r.String())
		if !ok || got != r {
			t.Errorf("ParseRole(%q) = %v, %v", r.String(), got, ok)
		}
	}
	if _, ok := ParseRole("ADMIN"); ok {
		t.Error("ParseRole accepted an unknown role")
	}
	if RoleUnknown.String() != "UNKNOWN" {
		t.Error("zero Role should be UNKNOWN")
	}
}

func TestScope_SetAndGet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	if _, ok := ScopeFrom(c); ok {
		t.Fatal("ScopeFrom reported a scope before SetScope")
	}
	want := Scope{HouseholdID: uuid.New(), StoreID: uuid.New(), UserID: uuid.New(), Role: RoleStaff}
	SetScope(c, want)
	if got, ok := ScopeFrom(c); !ok || got != want {
		t.Errorf("ScopeFrom = %+v, %v; want %+v", got, ok, want)
	}
}
