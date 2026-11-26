package kernel

import "time"

// Vietnam is the business time zone (Asia/Ho_Chi_Minh). Vietnam has used
// UTC+7 without daylight saving since 1975, so a fixed zone is exact for every
// date this application handles, and it needs no tzdata in the container.
var Vietnam = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60)

// Clock supplies the current time. Domain code receives a Clock (or a time
// value) instead of calling time.Now, so tests can control time.
type Clock interface {
	Now() time.Time
}

// SystemClock reads the real time, in UTC.
type SystemClock struct{}

// Now returns the current time in UTC.
func (SystemClock) Now() time.Time { return time.Now().UTC() }

// FixedClock always returns T. It is meant for tests.
type FixedClock struct{ T time.Time }

// Now returns c.T.
func (c FixedClock) Now() time.Time { return c.T }

// BusinessDate returns the business date of instant t: its calendar date in
// Vietnam. An order completed at 16:30 UTC on 1 December belongs to
// 1 December, while 17:00 UTC already belongs to 2 December (BR-25).
func BusinessDate(t time.Time) Date {
	return DateOf(t.In(Vietnam))
}
