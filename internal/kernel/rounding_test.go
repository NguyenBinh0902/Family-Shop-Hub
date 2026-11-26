package kernel

import (
	"errors"
	"math"
	"testing"
)

func TestDivRound_BR33_HalfAwayFromZero(t *testing.T) {
	tests := []struct {
		name     string
		num, den int64
		want     int64
	}{
		{"exact", 10, 2, 5},
		{"below half rounds down", 74_000, 7, 10_571}, // 10.571,43 (inventory.md example)
		{"half rounds up", 5, 2, 3},
		{"above half rounds up", 7, 4, 2},
		{"negative half rounds away from zero", -5, 2, -3},
		{"negative divisor half", 5, -2, -3},
		{"both negative half", -5, -2, 3},
		{"negative below half", -74_000, 7, -10_571},
		{"zero numerator", 0, 7, 0},
		{"large divisor no overflow", math.MaxInt64, math.MaxInt64 - 1, 1},
		{"min int64 numerator", math.MinInt64, 2, math.MinInt64 / 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DivRound(tt.num, tt.den); got != tt.want {
				t.Errorf("DivRound(%d, %d) = %d, want %d", tt.num, tt.den, got, tt.want)
			}
		})
	}
}

func TestDivRound_ReversalIsExactNegation(t *testing.T) {
	// The reason for rounding away from zero: reversing an entry must give
	// exactly the negated amount.
	for num := int64(-50); num <= 50; num++ {
		if DivRound(-num, 4) != -DivRound(num, 4) {
			t.Fatalf("DivRound(%d, 4) is not symmetric", num)
		}
	}
}

func TestMulDivRound_BR33(t *testing.T) {
	tests := []struct {
		name      string
		a, b, den int64
		want      int64
		wantErr   error
	}{
		{"simple", 3, 10_000, 1000, 30, nil},
		{"half up", 1, 1500, 1000, 2, nil},
		{"negative half away from zero", -1, 1500, 1000, -2, nil},
		{"intermediate exceeds int64 but result fits", math.MaxInt64, 1000, 1000, math.MaxInt64, nil},
		{"result overflows", math.MaxInt64, 2, 1, 0, ErrOverflow},
		{"negative result overflows", math.MinInt64, 2, 1, 0, ErrOverflow},
		{"min int64 fits", math.MinInt64, 1, 1, math.MinInt64, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MulDivRound(tt.a, tt.b, tt.den)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("MulDivRound(%d, %d, %d) = %d, want %d", tt.a, tt.b, tt.den, got, tt.want)
			}
		})
	}
}
