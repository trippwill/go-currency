package fixedpoint

import (
	"testing"
)

func TestX64String(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() X64
		expected string
	}{
		{
			name: "positive zero",
			setup: func() X64 {
				var x X64
				x.pack(kind_finite, signc_positive, 0, 0)
				return x
			},
			expected: "0",
		},
		{
			name: "negative zero",
			setup: func() X64 {
				var x X64
				x.pack(kind_finite, signc_negative, 0, 0)
				return x
			},
			expected: "-0",
		},
		{
			name: "positive infinity",
			setup: func() X64 {
				var x X64
				x.pack(kind_infinity, signc_positive, 0, 0)
				return x
			},
			expected: "Infinity",
		},
		{
			name: "negative infinity",
			setup: func() X64 {
				var x X64
				x.pack(kind_infinity, signc_negative, 0, 0)
				return x
			},
			expected: "-Infinity",
		},
		{
			name: "quiet NaN",
			setup: func() X64 {
				var x X64
				x.pack(kind_quiet, signc_positive, 0, 0)
				return x
			},
			expected: "qNaN",
		},
		{
			name: "negative quiet NaN",
			setup: func() X64 {
				var x X64
				x.pack(kind_quiet, signc_negative, 0, 0)
				return x
			},
			expected: "-qNaN",
		},
		{
			name: "signaling NaN",
			setup: func() X64 {
				var x X64
				x.pack(kind_signaling, signc_positive, 0, 0)
				return x
			},
			expected: "sNaN",
		},
		{
			name: "negative signaling NaN",
			setup: func() X64 {
				var x X64
				x.pack(kind_signaling, signc_negative, 0, 0)
				return x
			},
			expected: "-sNaN",
		},
		{
			name: "simple integer",
			setup: func() X64 {
				var x X64
				x.pack(kind_finite, signc_positive, 0, 123)
				return x
			},
			expected: "123",
		},
		{
			name: "negative integer",
			setup: func() X64 {
				var x X64
				x.pack(kind_finite, signc_negative, 0, 456)
				return x
			},
			expected: "-456",
		},
		{
			name: "positive exponent",
			setup: func() X64 {
				var x X64
				x.pack(kind_finite, signc_positive, 2, 789)
				return x
			},
			expected: "78900",
		},
		{
			name: "negative exponent with decimal",
			setup: func() X64 {
				var x X64
				x.pack(kind_finite, signc_positive, -2, 12345)
				return x
			},
			expected: "123.45",
		},
		{
			name: "large negative exponent",
			setup: func() X64 {
				var x X64
				x.pack(kind_finite, signc_positive, -5, 67)
				return x
			},
			expected: "0.00067",
		},
		{
			name: "scientific notation for large exponent",
			setup: func() X64 {
				var x X64
				x.pack(kind_finite, signc_positive, 10, 1234)
				return x
			},
			expected: "1.234e+13",
		},
		{
			name: "scientific notation for small exponent",
			setup: func() X64 {
				var x X64
				x.pack(kind_finite, signc_negative, -10, 5678)
				return x
			},
			expected: "-5.678e-7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := tt.setup()
			got := x.String()
			if got != tt.expected {
				t.Errorf("X64.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestX32String(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() X32
		expected string
	}{
		{
			name: "positive zero",
			setup: func() X32 {
				var x X32
				x.pack(kind_finite, signc_positive, 0, 0)
				return x
			},
			expected: "0",
		},
		{
			name: "negative zero",
			setup: func() X32 {
				var x X32
				x.pack(kind_finite, signc_negative, 0, 0)
				return x
			},
			expected: "-0",
		},
		{
			name: "positive infinity",
			setup: func() X32 {
				var x X32
				x.pack(kind_infinity, signc_positive, 0, 0)
				return x
			},
			expected: "Infinity",
		},
		{
			name: "negative infinity",
			setup: func() X32 {
				var x X32
				x.pack(kind_infinity, signc_negative, 0, 0)
				return x
			},
			expected: "-Infinity",
		},
		{
			name: "quiet NaN",
			setup: func() X32 {
				var x X32
				x.pack(kind_quiet, signc_positive, 0, 0)
				return x
			},
			expected: "qNaN",
		},
		{
			name: "decimal with zero exponent",
			setup: func() X32 {
				var x X32
				x.pack(kind_finite, signc_positive, 0, 42)
				return x
			},
			expected: "42",
		},
		{
			name: "decimal with positive exponent",
			setup: func() X32 {
				var x X32
				x.pack(kind_finite, signc_positive, 3, 123)
				return x
			},
			expected: "123000",
		},
		{
			name: "decimal with negative exponent",
			setup: func() X32 {
				var x X32
				x.pack(kind_finite, signc_negative, -2, 456)
				return x
			},
			expected: "-4.56",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := tt.setup()
			got := x.String()
			if got != tt.expected {
				t.Errorf("X32.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}
