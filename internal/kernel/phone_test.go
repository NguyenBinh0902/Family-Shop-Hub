package kernel

import (
	"errors"
	"testing"
)

func TestParsePhone(t *testing.T) {
	tests := []struct {
		in      string
		want    Phone
		wantErr bool
	}{
		{in: "0912345678", want: "0912345678"},
		{in: "0912 345 678", want: "0912345678"},
		{in: "0912.345.678", want: "0912345678"},
		{in: "0912-345-678", want: "0912345678"},
		{in: "+84912345678", want: "0912345678"},
		{in: "+84 912.345.678", want: "0912345678"},
		{in: "84912345678", want: "0912345678"},
		{in: "", wantErr: true},
		{in: "091234567", wantErr: true},   // 9 digits
		{in: "09123456789", wantErr: true}, // 11 digits
		{in: "1912345678", wantErr: true},  // does not start with 0
		{in: "09123a5678", wantErr: true},
		{in: "(091)2345678", wantErr: true}, // parentheses are not in the normalization rule
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParsePhone(tt.in)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidPhone) {
					t.Fatalf("ParsePhone(%q) err = %v, want ErrInvalidPhone", tt.in, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParsePhone(%q) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("ParsePhone(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestPhone_Masked_NFR13(t *testing.T) {
	if got, want := Phone("0912345678").Masked(), "091****678"; got != want {
		t.Errorf("Masked() = %q, want %q", got, want)
	}
	if got := Phone("").Masked(); got != "***" {
		t.Errorf("Masked() of empty phone = %q, want ***", got)
	}
}
