package currency

import (
	"fmt"
	"testing"

	"golang.org/x/text/language"
)

func ExampleFixedPoint_Format() {
	fp := FixedPoint{
		Base:  12345,
		Scale: 2,
	}
	fmt.Println(fp.Format(language.English))
	fmt.Println(fp.Format(language.German))
	// Output:
	// 123.45
	// 123,45
}

func ExampleFixedPoint_String() {
	fp := FixedPoint{
		Base:  12345,
		Scale: 2,
	}
	fmt.Println(fp.String())
	// Output:
	// 123.45
}
func TestNewFixedPoint(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		wantBase      int64
		wantPrecision uint8
	}{
		{"Zero", "0.0", 0, 0},
		{"Simple", "12.34", 1234, 2},
		{"SingleDigitFrac", "1.1", 11, 1},
		{"TrailingZero", "12.50", 1250, 2},
		{"LargerFraction", "100.234", 100234, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := NewFixedPoint(tt.value, nil)
			if fp.Base != tt.wantBase {
				t.Errorf("NewFixedPoint(%v): Base = %v, want %v", tt.value, fp.Base, tt.wantBase)
			}
			if fp.Scale != tt.wantPrecision {
				t.Errorf("NewFixedPoint(%v): Precision = %v, want %v", tt.value, fp.Scale, tt.wantPrecision)
			}
		})
	}
}

func ExampleNewFixedPoint() {
	fp := NewFixedPoint("123.45", DefaultParseOpts)
	// fp.Base is 12345 and fp.Precision is 2, so String() prints "123.45"
	fmt.Println(fp.String())
	// Output:
	// 123.45
}
