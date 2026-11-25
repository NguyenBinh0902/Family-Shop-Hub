package kernel

import (
	"fmt"
	"time"
)

// Date is a calendar date with no time of day and no time zone, such as a
// business date or a warranty end date. A separate type avoids the classic bug
// of a time.Time at midnight UTC showing up as the previous day in Vietnam.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

const dateLayout = "2006-01-02"

// DateOf returns the calendar date of t in t's own location.
func DateOf(t time.Time) Date {
	y, m, d := t.Date()
	return Date{Year: y, Month: m, Day: d}
}

// ParseDate parses an ISO date such as "2026-12-01".
func ParseDate(s string) (Date, error) {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return Date{}, fmt.Errorf("parse date %q: %w", s, err)
	}
	return DateOf(t), nil
}

// String formats d as "2006-01-02".
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, int(d.Month), d.Day)
}

// IsZero reports whether d is the zero Date.
func (d Date) IsZero() bool { return d == Date{} }

// Compare returns -1, 0 or +1 when d is before, equal to or after o.
func (d Date) Compare(o Date) int {
	return d.midnightUTC().Compare(o.midnightUTC())
}

// AddMonthsClamped adds n months, clamping the day to the end of the target
// month: 31/01 + 1 month is 28/02 (or 29/02 in a leap year). time.AddDate
// would instead overflow into March. Used for warranty end dates (BR-12).
func (d Date) AddMonthsClamped(n int) Date {
	first := time.Date(d.Year, d.Month+time.Month(n), 1, 0, 0, 0, 0, time.UTC)
	y, m := first.Year(), first.Month()
	return Date{Year: y, Month: m, Day: min(d.Day, daysIn(y, m))}
}

func (d Date) midnightUTC() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
}

// daysIn returns the number of days in month m of year y. Day 0 of the next
// month normalizes to the last day of m.
func daysIn(y int, m time.Month) int {
	return time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
