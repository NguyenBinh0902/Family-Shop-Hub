package kernel

import (
	"errors"
	"strconv"
	"strings"
)

// QuantityScale is the number of Quantity units in one display unit:
// quantities are stored in thousandths (BR-34).
const QuantityScale = 1000

// Quantity is an amount of a product in thousandths of its unit, so 1 cup is
// 1000 and 0.5 kg is 500 (BR-34). Using an integer keeps stock arithmetic exact.
type Quantity int64

// ErrInvalidQuantity reports a string that is not a decimal number with at
// most three fractional digits.
var ErrInvalidQuantity = errors.New("kernel: invalid quantity")

// QuantityOf returns n whole units as a Quantity.
func QuantityOf(n int64) Quantity { return Quantity(n * QuantityScale) }

// ParseQuantity parses a decimal string such as "2", "0.5" or "-1.25".
// It rejects more than three fractional digits instead of silently rounding,
// because a lost digit would be a lost amount of stock.
func ParseQuantity(s string) (Quantity, error) {
	negative := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")

	whole, frac, hasDot := strings.Cut(s, ".")
	if whole == "" || (hasDot && frac == "") || len(frac) > 3 || !isDigits(whole) || !isDigits(frac) {
		return 0, ErrInvalidQuantity
	}

	frac += strings.Repeat("0", 3-len(frac))
	n, err := strconv.ParseInt(whole+frac, 10, 64)
	if err != nil {
		return 0, ErrInvalidQuantity
	}
	if negative {
		n = -n
	}
	return Quantity(n), nil
}

// String formats q as a decimal without trailing zeros: 500 → "0.5",
// 2000 → "2", -1250 → "-1.25".
func (q Quantity) String() string {
	sign := ""
	mag := absUint(int64(q))
	if q < 0 {
		sign = "-"
	}
	whole := mag / QuantityScale
	frac := mag % QuantityScale
	if frac == 0 {
		return sign + strconv.FormatUint(whole, 10)
	}
	fracStr := strings.TrimRight(strconv.FormatUint(frac+QuantityScale, 10)[1:], "0")
	return sign + strconv.FormatUint(whole, 10) + "." + fracStr
}

// MulMoney returns q × unit price, i.e. the value of q at price m per unit.
// The division by QuantityScale happens once at the end and is rounded half
// away from zero (BR-33, BR-34), so 0.333 × 10.000 đ = 3.330 đ exactly.
func (q Quantity) MulMoney(m Money) (Money, error) {
	v, err := MulDivRound(int64(q), int64(m), QuantityScale)
	return Money(v), err
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
