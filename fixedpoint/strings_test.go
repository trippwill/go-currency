package fixedpoint

import (
	"strings"
	"testing"
)

func TestFixedPoint128_String(t *testing.T) {
	tests := []struct {
		name     string
		fp       FixedPoint128
		expected string
	}{
		{
			name:     "zero",
			fp:       mustParse128("0"),
			expected: "0",
		},
		{
			name:     "positive integer",
			fp:       mustParse128("42"),
			expected: "42",
		},
		{
			name:     "negative integer",
			fp:       mustParse128("-42"),
			expected: "-42",
		},
		{
			name:     "decimal with positive exponent",
			fp:       mustParse128("12300"),
			expected: "12300",
		},
		{
			name:     "decimal with negative exponent",
			fp:       mustParse128("1.23"),
			expected: "1.23",
		},
		{
			name:     "decimal with negative exponent ending with zero",
			fp:       mustParse128("1.230"),
			expected: "1.23",
		},
		{
			name:     "negative decimal with negative exponent",
			fp:       mustParse128("-1.23"),
			expected: "-1.23",
		},
		{
			name:     "NaN",
			fp:       mustParse128("NaN"),
			expected: "NaN",
		},
		{
			name:     "sNaN",
			fp:       mustParse128("sNaN"),
			expected: "sNaN",
		},
		{
			name:     "positive infinity",
			fp:       mustParse128("Infinity"),
			expected: "Infinity",
		},
		{
			name:     "negative infinity",
			fp:       mustParse128("-Infinity"),
			expected: "-Infinity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fp.String()
			if result != tt.expected {
				t.Errorf("FixedPoint128.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFixedPoint128_Scientific(t *testing.T) {
	tests := []struct {
		name     string
		fp       FixedPoint128
		expected string
	}{
		{
			name:     "zero",
			fp:       mustParse128("0"),
			expected: "0e+0",
		},
		{
			name:     "positive integer",
			fp:       mustParse128("42"),
			expected: "4.2e+1",
		},
		{
			name:     "negative integer",
			fp:       mustParse128("-42"),
			expected: "-4.2e+1",
		},
		{
			name:     "large number",
			fp:       mustParse128("1234567"),
			expected: "1.234567e+6",
		},
		{
			name:     "small decimal",
			fp:       mustParse128("0.01234"),
			expected: "1.234e-2",
		},
		{
			name:     "negative small decimal",
			fp:       mustParse128("-0.01234"),
			expected: "-1.234e-2",
		},
		{
			name:     "NaN",
			fp:       mustParse128("NaN"),
			expected: "NaN",
		},
		{
			name:     "sNaN",
			fp:       mustParse128("sNaN"),
			expected: "sNaN",
		},
		{
			name:     "positive infinity",
			fp:       mustParse128("Infinity"),
			expected: "Infinity",
		},
		{
			name:     "negative infinity",
			fp:       mustParse128("-Infinity"),
			expected: "-Infinity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fp.Scientific()
			if result != tt.expected {
				t.Errorf("FixedPoint128.Scientific() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFixedPoint128_Debug(t *testing.T) {
	tests := []struct {
		name         string
		fp           FixedPoint128
		expectedKind string
	}{
		{
			name:         "regular number",
			fp:           mustParse128("42"),
			expectedKind: "Kind: Finite",
		},
		{
			name:         "NaN",
			fp:           mustParse128("NaN"),
			expectedKind: "Kind: Quiet NaN",
		},
		{
			name:         "sNaN",
			fp:           mustParse128("sNaN"),
			expectedKind: "Kind: Signaling NaN",
		},
		{
			name:         "positive infinity",
			fp:           mustParse128("Infinity"),
			expectedKind: "Kind: +Infinity",
		},
		{
			name:         "negative infinity",
			fp:           mustParse128("-Infinity"),
			expectedKind: "Kind: -Infinity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fp.Debug()
			if !strings.HasPrefix(result, tt.expectedKind) {
				t.Errorf("FixedPoint128.Debug() = %v, expected to start with %v", result, tt.expectedKind)
			}
		})
	}
}

// Helper function to handle errors from Parse128
func mustParse128(s string) FixedPoint128 {
	fp, err := Parse128(s)
	if err != nil {
		panic("Failed to parse: " + s + ", error: " + err.Error())
	}
	return fp
}
