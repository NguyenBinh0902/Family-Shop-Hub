package kernel

import (
	"errors"
	"math"
	"testing"
)

func TestParseQuantity_BR34(t *testing.T) {
	tests := []struct {
		in      string
		want    Quantity
		wantErr bool
	}{
		{in: "2", want: 2000},
		{in: "0.5", want: 500},
		{in: "1.25", want: 1250},
		{in: "0.001", want: 1},
		{in: "-1.25", want: -1250},
		{in: "15000", want: 15_000_000},
		{in: "0.0005", wantErr: true}, // more than three decimals is refused, not rounded
		{in: "", wantErr: true},
		{in: "-", wantErr: true},
		{in: "1.", wantErr: true},
		{in: ".5", wantErr: true},
		{in: "1,5", wantErr: true},
		{in: "abc", wantErr: true},
		{in: "+1", wantErr: true},
		{in: "99999999999999999999", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseQuantity(tt.in)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidQuantity) {
					t.Fatalf("ParseQuantity(%q) err = %v, want ErrInvalidQuantity", tt.in, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseQuantity(%q) unexpected error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("ParseQuantity(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestQuantity_String(t *testing.T) {
	tests := []struct {
		q    Quantity
		want string
	}{
		{0, "0"},
		{500, "0.5"},
		{2000, "2"},
		{1250, "1.25"},
		{1, "0.001"},
		{-1250, "-1.25"},
		{-200_000, "-200"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.q.String(); got != tt.want {
				t.Errorf("Quantity(%d).String() = %q, want %q", tt.q, got, tt.want)
			}
			back, err := ParseQuantity(tt.want)
			if err != nil || back != tt.q {
				t.Errorf("round trip %q = %d, %v; want %d", tt.want, back, err, tt.q)
			}
		})
	}
}

func TestQuantity_MulMoney_BR33(t *testing.T) {
	tests := []struct {
		name    string
		q       Quantity
		price   Money
		want    Money
		wantErr error
	}{
		{"whole units", QuantityOf(2), 15_000, 30_000, nil},
		{"ml at đồng per ml", QuantityOf(5000), 36, 180_000, nil},     // UC-12
		{"half gram rounds once at the end", 500, 400, 200, nil},      // 0.5 g × 400 đ
		{"fraction rounds half up", 333, 10_001, 3330, nil},           // 3330.333
		{"exact half rounds up", 500, 3, 2, nil},                      // 1.5
		{"negative quantity rounds away from zero", -500, 3, -2, nil}, // reversal of the line above
		{"recipe line", QuantityOf(200), 30, 6000, nil},               // UC-11: 200 ml × 30 đ
		{"overflow is reported", Quantity(math.MaxInt64), math.MaxInt64, 0, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.MulMoney(tt.price)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("%s × %d = %d, want %d", tt.q, tt.price, got, tt.want)
			}
		})
	}
}
