package kernel

import "strconv"

// Money is an amount in Vietnamese đồng. It is always an integer: VND has no
// minor unit in practice and floating point is never used for money (NFR-09).
type Money int64

// Add returns m + o.
func (m Money) Add(o Money) Money { return m + o }

// Sub returns m - o.
func (m Money) Sub(o Money) Money { return m - o }

// Neg returns -m, used when writing reversal entries.
func (m Money) Neg() Money { return -m }

// String returns the plain integer amount, for logs and debugging only.
// Formatting for people (thousand separators, "₫") is the frontend's job.
func (m Money) String() string { return strconv.FormatInt(int64(m), 10) }
