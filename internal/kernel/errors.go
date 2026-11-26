package kernel

// ErrorKind classifies a business error so the HTTP layer can choose a status
// code without knowing every error. The zero value is invalid on purpose.
type ErrorKind int

// Error kinds and the HTTP status platform/problem maps them to.
const (
	ErrorKindUnknown   ErrorKind = iota
	ErrorKindInvalid             // 422: input breaks a rule
	ErrorKindConflict            // 409: current state forbids the action
	ErrorKindForbidden           // 403: the user may not do this
	ErrorKindNotFound            // 404: the thing does not exist for this user
)

// FieldError attaches a message to one input field, such as "name".
type FieldError struct {
	Field   string
	Message string
}

// Error is a business rule violation that is shown to the user. Code is the
// rule code ("BR-15") or a generic code such as "VALIDATION"; Message is in
// Vietnamese because it is displayed as is.
//
// Each context declares its errors as package-level values in its domain
// package:
//
//	var ErrOrderImmutable = &kernel.Error{
//		Code: "BR-15", Kind: kernel.ErrorKindConflict,
//		Message: "Đơn đã hoàn tất, không thể sửa",
//	}
type Error struct {
	Code    string
	Message string
	Kind    ErrorKind
	Fields  []FieldError
}

// CodeValidation is the code of input validation errors that are not tied to
// a numbered business rule.
const CodeValidation = "VALIDATION"

// NewValidationError returns an ErrorKindInvalid error listing the bad fields.
func NewValidationError(fields ...FieldError) *Error {
	return &Error{
		Code:    CodeValidation,
		Message: "Dữ liệu không hợp lệ",
		Kind:    ErrorKindInvalid,
		Fields:  fields,
	}
}

// Error implements the error interface.
func (e *Error) Error() string { return e.Code + ": " + e.Message }

// Is makes errors.Is match on Code, so a copy made by WithMessage or
// WithFields still matches the package-level value it came from.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && t.Code == e.Code
}

// WithMessage returns a copy of e with a more specific message, for example
// one that names the product that is out of stock. e itself is not modified,
// because it is usually a shared package-level value.
func (e *Error) WithMessage(msg string) *Error {
	c := *e
	c.Message = msg
	return &c
}

// WithFields returns a copy of e carrying the given field errors.
func (e *Error) WithFields(fields ...FieldError) *Error {
	c := *e
	c.Fields = append(append([]FieldError(nil), e.Fields...), fields...)
	return &c
}
