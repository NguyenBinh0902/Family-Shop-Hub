package kernel

import (
	"errors"
	"math"
	"math/bits"
)

// ErrOverflow reports that an arithmetic result does not fit in an int64.
var ErrOverflow = errors.New("kernel: arithmetic overflow")

// DivRound returns num/den rounded to the nearest integer, with halves rounded
// away from zero (BR-33): 5/2 = 3 and -5/2 = -3. Rounding away from zero keeps
// a reversal exactly equal to the negation of the original entry.
//
// Like the / operator, it panics if den is zero.
func DivRound(num, den int64) int64 {
	q, r := num/den, num%den
	if r == 0 {
		return q
	}
	absR, absDen := absUint(r), absUint(den)
	// Compare r with den-r instead of 2r with den so the check cannot overflow.
	if absR >= absDen-absR {
		if (num < 0) != (den < 0) {
			q--
		} else {
			q++
		}
	}
	return q
}

// MulDivRound returns a*b/den rounded half away from zero (BR-33). The product
// is computed with 128 bits, so it is exact even when a*b alone would overflow
// an int64; ErrOverflow is returned only when the final result does not fit.
//
// It panics if den is zero.
func MulDivRound(a, b, den int64) (int64, error) {
	if den == 0 {
		panic("kernel: MulDivRound by zero")
	}
	negative := (a < 0) != (b < 0) != (den < 0)

	hi, lo := bits.Mul64(absUint(a), absUint(b))
	d := absUint(den)
	if hi >= d {
		// The quotient would need more than 64 bits; bits.Div64 would panic.
		return 0, ErrOverflow
	}
	q, r := bits.Div64(hi, lo, d)
	if r >= d-r {
		q++
	}
	return signed(q, negative)
}

// signed converts a magnitude plus sign back to int64, rejecting values that
// do not fit.
func signed(mag uint64, negative bool) (int64, error) {
	if negative {
		if mag > uint64(math.MaxInt64)+1 {
			return 0, ErrOverflow
		}
		// -mag wraps around in uint64; converting that to int64 yields the
		// intended negative value (two's complement), range checked above.
		return int64(-mag), nil
	}
	if mag > math.MaxInt64 {
		return 0, ErrOverflow
	}
	return int64(mag), nil
}

// absUint returns |v| as uint64; it is exact even for math.MinInt64.
func absUint(v int64) uint64 {
	if v < 0 {
		return uint64(-(v + 1)) + 1
	}
	return uint64(v)
}
