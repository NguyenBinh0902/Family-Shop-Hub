package kernel

import (
	"errors"
	"strings"
)

// Phone is a normalized Vietnamese phone number: exactly 10 digits starting
// with 0. Construct it with ParsePhone; the zero value means "no phone".
type Phone string

// ErrInvalidPhone reports input that does not normalize to a valid phone.
var ErrInvalidPhone = errors.New("kernel: invalid phone number")

// ParsePhone normalizes s following domain/identity.md: it drops spaces, dots
// and dashes, turns a leading +84 or 84 into 0, and requires 10 digits
// starting with 0. "+84 912.345.678" becomes "0912345678".
func ParsePhone(s string) (Phone, error) {
	s = strings.NewReplacer(" ", "", ".", "", "-", "").Replace(s)
	switch {
	case strings.HasPrefix(s, "+84"):
		s = "0" + s[3:]
	case strings.HasPrefix(s, "84"):
		s = "0" + s[2:]
	}
	if len(s) != 10 || s[0] != '0' || !isDigits(s) {
		return "", ErrInvalidPhone
	}
	return Phone(s), nil
}

// String returns the normalized number.
func (p Phone) String() string { return string(p) }

// Masked hides the middle four digits ("091****678") so phone numbers can be
// logged without exposing them in full (NFR-13).
func (p Phone) Masked() string {
	if len(p) != 10 {
		return "***"
	}
	return string(p[:3]) + "****" + string(p[7:])
}
