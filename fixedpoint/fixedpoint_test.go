package fixedpoint

import (
	"testing"
)

func TestNewFixedPoint(t *testing.T) {
	tests := []struct {
		input    string
		expected FixedPoint
		hasError bool
	}{
		{"123.45", FixedPoint{s: false, c: 12345, e: -2}, false},
		{"-0.001", FixedPoint{s: true, c: 1, e: -3}, false},
		{"invalid", FixedPoint{}, true},
	}

	for _, test := range tests {
		fp, err := NewFixedPoint(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("NewFixedPoint(%q) error = %v, wantErr %v", test.input, err, test.hasError)
			continue
		}
		if !test.hasError && !fp.Equals(test.expected) {
			t.Errorf("NewFixedPoint(%q) = %+v, want %+v", test.input, fp, test.expected)
		}
	}
}

func TestFixedPointAdd(t *testing.T) {
	tests := []struct {
		a, b     string
		expected string
	}{
		{"1.23", "4.56", "5.79"},
		{"-1.23", "1.23", "0"},
		{"0", "0", "0"},
	}

	for _, test := range tests {
		fp1, err1 := NewFixedPoint(test.a)
		fp2, err2 := NewFixedPoint(test.b)
		expected, err3 := NewFixedPoint(test.expected)

		if err1 != nil || err2 != nil || err3 != nil {
			t.Fatalf("Error creating FixedPoint: %v, %v, %v", err1, err2, err3)
		}

		result := fp1.Add(fp2)
		if !result.Equals(expected) {
			t.Errorf("Add(%q, %q) = %+v, want %+v", test.a, test.b, result, expected)
		}
	}
}

func TestFixedPointSpecialValues(t *testing.T) {
	if !NaN.IsNaN() {
		t.Error("NaN should be NaN")
	}
	if !Inf.IsInf() || Inf.IsNegative() {
		t.Error("Inf should be positive infinity")
	}
	if !NegInf.IsInf() || !NegInf.IsNegative() {
		t.Error("NegInf should be negative infinity")
	}
	if !Zero.IsZero() || Zero.IsNegative() {
		t.Error("Zero should be positive zero")
	}
	if !NegZero.IsZero() || !NegZero.IsNegative() {
		t.Error("NegZero should be negative zero")
	}
}

func TestAddUsingContext(t *testing.T) {
	tests := []struct {
		name     string
		x        FixedPoint
		y        FixedPoint
		expected FixedPoint
	}{
		{
			name: "Infinite values",
			x:    Inf,
			y:    NegInf,
			expected: FixedPoint{
				x: x_sNan,
			},
		},
		{
			name: "NaN values",
			x:    NaN,
			y:    FixedPoint{c: 123, e: 0},
			expected: FixedPoint{
				x: x_Nan,
			},
		},
		{
			name: "Zero values",
			x:    FixedPoint{c: 0, e: 0},
			y:    FixedPoint{c: 0, e: 0},
			expected: FixedPoint{
				c: 0,
				e: 0,
			},
		},
		{
			name: "Opposite signs",
			x:    FixedPoint{c: 123, e: 0, s: false},
			y:    FixedPoint{c: 123, e: 0, s: true},
			expected: FixedPoint{
				c: 0,
				e: 0,
			},
		},
		{
			name: "Large exponent differences",
			x:    FixedPoint{c: 1, e: 10},
			y:    FixedPoint{c: 1, e: 0},
			expected: FixedPoint{
				c: 10000000001,
				e: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := DefaultContext
			result := tt.x.AddUsingContext(tt.y, ctx)
			if !result.Equals(tt.expected) {
				t.Errorf("AddUsingContext(%v, %v) = %v, want %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

func TestFixedPointString(t *testing.T) {
	tests := []struct {
		name     string
		fp       FixedPoint
		expected string
	}{
		{
			name:     "NaN",
			fp:       NaN,
			expected: "NaN",
		},
		{
			name:     "Positive Infinity",
			fp:       Inf,
			expected: "Inf",
		},
		{
			name:     "Negative Infinity",
			fp:       NegInf,
			expected: "-Inf",
		},
		{
			name:     "Zero",
			fp:       Zero,
			expected: "0",
		},
		{
			name:     "Negative Zero",
			fp:       NegZero,
			expected: "0",
		},
		{
			name:     "Positive FixedPoint with no fractional part",
			fp:       FixedPoint{s: false, c: 12345, e: 0},
			expected: "12345",
		},
		{
			name:     "Negative FixedPoint with no fractional part",
			fp:       FixedPoint{s: true, c: 12345, e: 0},
			expected: "-12345",
		},
		{
			name:     "Positive FixedPoint with fractional part",
			fp:       FixedPoint{s: false, c: 12345, e: -2},
			expected: "123.45",
		},
		{
			name:     "Negative FixedPoint with fractional part",
			fp:       FixedPoint{s: true, c: 12345, e: -2},
			expected: "-123.45",
		},
		{
			name:     "Positive FixedPoint with large exponent",
			fp:       FixedPoint{s: false, c: 12345, e: 2},
			expected: "1234500",
		},
		{
			name:     "Negative FixedPoint with large exponent",
			fp:       FixedPoint{s: true, c: 12345, e: 2},
			expected: "-1234500",
		},
		{
			name:     "Positive FixedPoint with small fractional part",
			fp:       FixedPoint{s: false, c: 1, e: -5},
			expected: "0.00001",
		},
		{
			name:     "Negative FixedPoint with small fractional part",
			fp:       FixedPoint{s: true, c: 1, e: -5},
			expected: "-0.00001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fp.String()
			if result != tt.expected {
				t.Errorf("FixedPoint.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFixedPointEquals(t *testing.T) {
	tests := []struct {
		name     string
		fp1      FixedPoint
		fp2      FixedPoint
		expected bool
	}{
		{
			name:     "Equal finite values",
			fp1:      FixedPoint{s: false, c: 12345, e: -2},
			fp2:      FixedPoint{s: false, c: 12345, e: -2},
			expected: true,
		},
		{
			name:     "Different signs",
			fp1:      FixedPoint{s: false, c: 12345, e: -2},
			fp2:      FixedPoint{s: true, c: 12345, e: -2},
			expected: false,
		},
		{
			name:     "Different coefficients",
			fp1:      FixedPoint{s: false, c: 12345, e: -2},
			fp2:      FixedPoint{s: false, c: 54321, e: -2},
			expected: false,
		},
		{
			name:     "Different exponents",
			fp1:      FixedPoint{s: false, c: 12345, e: -2},
			fp2:      FixedPoint{s: false, c: 12345, e: -3},
			expected: false,
		},
		{
			name:     "Both NaN",
			fp1:      NaN,
			fp2:      NaN,
			expected: true,
		},
		{
			name:     "One NaN",
			fp1:      NaN,
			fp2:      FixedPoint{s: false, c: 12345, e: -2},
			expected: false,
		},
		{
			name:     "Both Inf with same sign",
			fp1:      Inf,
			fp2:      Inf,
			expected: true,
		},
		{
			name:     "Both Inf with different signs",
			fp1:      Inf,
			fp2:      NegInf,
			expected: false,
		},
		{
			name:     "One Inf",
			fp1:      Inf,
			fp2:      FixedPoint{s: false, c: 12345, e: -2},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fp1.Equals(tt.fp2)
			if result != tt.expected {
				t.Errorf("Equals(%v, %v) = %v, want %v", tt.fp1, tt.fp2, result, tt.expected)
			}
		})
	}
}

func FuzzFixedPointAdd(f *testing.F) {
	// Seed fuzz inputs with varied ranges, including negatives, very large and very small numbers.
	f.Add("1.23", "4.56")
	f.Add("-1.23", "4.56")
	f.Add("123.456", "-78.90")
	f.Add("0", "9876.54321")
	f.Add("-0.001", "-999.99")
	f.Add("9999999999999999999", "0")
	f.Add("0", "0.000000000000000001")
	f.Add("-9999999999999999999", "0")
	f.Add("0", "-0.000000000000000001")

	f.Fuzz(func(t *testing.T, aStr, bStr string) {
		fp1, err1 := NewFixedPoint(aStr)
		fp2, err2 := NewFixedPoint(bStr)
		if err1 != nil || err2 != nil {
			t.Skip() // skip invalid inputs
		}

		result1 := fp1.Add(fp2)
		result2 := fp2.Add(fp1)
		// Test commutativity only for finite values.
		if fp1.IsFinite() && fp2.IsFinite() {
			if !result1.Equals(result2) {
				t.Errorf("Addition is not commutative: %v + %v = %v, but %v + %v = %v", fp1, fp2, result1, fp2, fp1, result2)
			}
		}
		// Ensure the result's String() is not empty.
		if result1.String() == "" {
			t.Errorf("Result string should not be empty for %v + %v", fp1, fp2)
		}
	})
}
