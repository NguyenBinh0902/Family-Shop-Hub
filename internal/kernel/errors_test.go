package kernel

import (
	"errors"
	"fmt"
	"testing"
)

var errTestImmutable = &Error{Code: "BR-15", Kind: ErrorKindConflict, Message: "Đơn đã hoàn tất, không thể sửa"}

func TestError_IsMatchesCopiesByCode(t *testing.T) {
	specific := errTestImmutable.WithMessage("Đơn 261201-001 đã hoàn tất")
	wrapped := fmt.Errorf("cancel order: %w", specific)

	if !errors.Is(wrapped, errTestImmutable) {
		t.Error("errors.Is should match a wrapped copy with the same code")
	}
	if errTestImmutable.Message != "Đơn đã hoàn tất, không thể sửa" {
		t.Error("WithMessage modified the shared package-level error")
	}

	var be *Error
	if !errors.As(wrapped, &be) || be.Message != "Đơn 261201-001 đã hoàn tất" {
		t.Errorf("errors.As = %v, want the specific copy", be)
	}

	other := &Error{Code: "BR-13", Kind: ErrorKindConflict}
	if errors.Is(wrapped, other) {
		t.Error("errors.Is matched a different code")
	}
}

func TestError_WithFieldsDoesNotShareBackingArray(t *testing.T) {
	base := NewValidationError(FieldError{Field: "name", Message: "Bắt buộc"})
	a := base.WithFields(FieldError{Field: "price", Message: "Phải ≥ 0"})
	b := base.WithFields(FieldError{Field: "unit", Message: "Bắt buộc"})

	if len(base.Fields) != 1 || len(a.Fields) != 2 || len(b.Fields) != 2 {
		t.Fatalf("unexpected lengths: base=%d a=%d b=%d", len(base.Fields), len(a.Fields), len(b.Fields))
	}
	if a.Fields[1].Field != "price" {
		t.Errorf("a.Fields[1] = %q, want price (overwritten by a later WithFields?)", a.Fields[1].Field)
	}
	if base.Kind != ErrorKindInvalid || base.Code != CodeValidation {
		t.Errorf("NewValidationError kind/code = %v/%s", base.Kind, base.Code)
	}
}

func TestNewID_IsVersion7AndOrdered(t *testing.T) {
	a, b := NewID(), NewID()
	if a.Version() != 7 {
		t.Errorf("version = %d, want 7", a.Version())
	}
	if a.String() >= b.String() {
		t.Errorf("IDs are not increasing: %s then %s", a, b)
	}
}
