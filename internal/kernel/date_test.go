package kernel

import (
	"testing"
	"time"
)

func TestBusinessDate_BR25_UTCvsVietnam(t *testing.T) {
	tests := []struct {
		name string
		at   time.Time
		want Date
	}{
		// UC-16: completed at 23:30 on 01/12 Vietnam time = 16:30 UTC.
		{"late evening in Vietnam is still the same day", time.Date(2026, 12, 1, 16, 30, 0, 0, time.UTC), Date{2026, 12, 1}},
		{"midnight in Vietnam is the next day", time.Date(2026, 12, 1, 17, 0, 0, 0, time.UTC), Date{2026, 12, 2}},
		{"early UTC morning", time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC), Date{2026, 12, 1}},
		{"new year in Vietnam", time.Date(2026, 12, 31, 17, 0, 0, 0, time.UTC), Date{2027, 1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BusinessDate(tt.at); got != tt.want {
				t.Errorf("BusinessDate(%s) = %s, want %s", tt.at, got, tt.want)
			}
		})
	}
}

func TestDate_AddMonthsClamped(t *testing.T) {
	tests := []struct {
		from   Date
		months int
		want   Date
	}{
		{Date{2026, 1, 31}, 1, Date{2026, 2, 28}},
		{Date{2028, 1, 31}, 1, Date{2028, 2, 29}}, // leap year
		{Date{2026, 3, 31}, 1, Date{2026, 4, 30}},
		{Date{2026, 12, 15}, 12, Date{2027, 12, 15}},
		{Date{2026, 11, 30}, 3, Date{2027, 2, 28}},
		{Date{2026, 5, 10}, 0, Date{2026, 5, 10}},
		{Date{2026, 3, 31}, -1, Date{2026, 2, 28}},
	}
	for _, tt := range tests {
		t.Run(tt.from.String(), func(t *testing.T) {
			if got := tt.from.AddMonthsClamped(tt.months); got != tt.want {
				t.Errorf("%s + %d months = %s, want %s", tt.from, tt.months, got, tt.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	d, err := ParseDate("2026-12-01")
	if err != nil {
		t.Fatalf("ParseDate: %v", err)
	}
	if want := (Date{2026, 12, 1}); d != want {
		t.Errorf("ParseDate = %v, want %v", d, want)
	}
	if d.String() != "2026-12-01" {
		t.Errorf("String() = %q", d.String())
	}
	if _, err := ParseDate("01/12/2026"); err == nil {
		t.Error("ParseDate accepted a non-ISO date")
	}
}

func TestDate_Compare(t *testing.T) {
	a, b := Date{2026, 12, 1}, Date{2026, 12, 2}
	if a.Compare(b) != -1 || b.Compare(a) != 1 || a.Compare(a) != 0 {
		t.Error("Compare returned the wrong order")
	}
	if !(Date{}).IsZero() || a.IsZero() {
		t.Error("IsZero is wrong")
	}
}
