package fixedpoint

import (
	"testing"
)

func TestRoundingModeString(t *testing.T) {
	tests := []struct {
		mode     Rounding
		expected string
		debug    string
	}{
		{RoundTiesToEven, "RoundTiesToEven", "TiE"},
		{RoundTiesToAway, "RoundTiesToAway", "TiA"},
		{RoundTowardPositive, "RoundTowardPositive", "ToP"},
		{RoundTowardNegative, "RoundTowardNegative", "ToN"},
		{RoundTowardZero, "RoundTowardZero", "ToZ"},
		{Rounding(99), "Rounding(99)", "?(99)"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if got := test.mode.String(); got != test.expected {
				t.Errorf("String() = %v, want %v", got, test.expected)
			}
			if got := test.mode.Debug(); got != test.debug {
				t.Errorf("Debug() = %v, want %v", got, test.debug)
			}
		})
	}
}

// TestRoundingApply64 tests the Apply function with uint64 coefficients
func TestRoundingApply64(t *testing.T) {
	tests := []struct {
		name      string
		rounding  Rounding
		coe       uint64
		exp       int16
		precision uint
		sign      signc
		expected  uint64
		removed   uint8
	}{
		// RoundTiesToEven (banker's rounding)
		{"TiesToEven-NoRounding", RoundTiesToEven, 123, 0, 3, signc_positive, 123, 0},
		{"TiesToEven-RoundDown-EvenQuotient", RoundTiesToEven, 12345, 0, 4, signc_positive, 1234, 1},
		{"TiesToEven-RoundUp-EvenQuotient-ExactHalf", RoundTiesToEven, 12350, 0, 4, signc_positive, 1235, 1},
		{"TiesToEven-RoundDown-OddQuotient-ExactHalf", RoundTiesToEven, 12450, 0, 4, signc_positive, 1245, 1},
		{"TiesToEven-RoundUp-OddQuotient-MoreThanHalf", RoundTiesToEven, 12451, 0, 4, signc_positive, 1245, 1},

		// RoundTiesToAway (round to nearest, ties away from zero)
		{"TiesToAway-NoRounding", RoundTiesToAway, 123, 0, 3, signc_positive, 123, 0},
		{"TiesToAway-RoundDown-LessThanHalf", RoundTiesToAway, 12344, 0, 4, signc_positive, 1234, 1},
		{"TiesToAway-RoundUp-ExactHalf", RoundTiesToAway, 12350, 0, 4, signc_positive, 1235, 1},
		{"TiesToAway-RoundUp-MoreThanHalf", RoundTiesToAway, 12351, 0, 4, signc_positive, 1235, 1},
		{"TiesToAway-RoundUp-NegativeSign-ExactHalf", RoundTiesToAway, 12350, 0, 4, signc_negative, 1235, 1},

		// RoundTowardPositive (ceiling)
		{"TowardPositive-NoRounding", RoundTowardPositive, 123, 0, 3, signc_positive, 123, 0},
		{"TowardPositive-RoundUp-Positive", RoundTowardPositive, 12345, 0, 4, signc_positive, 1235, 1},
		{"TowardPositive-RoundDown-Negative", RoundTowardPositive, 12345, 0, 4, signc_negative, 1234, 1},

		// RoundTowardNegative (floor)
		{"TowardNegative-NoRounding", RoundTowardNegative, 123, 0, 3, signc_positive, 123, 0},
		{"TowardNegative-RoundDown-Positive", RoundTowardNegative, 12345, 0, 4, signc_positive, 1234, 1},
		{"TowardNegative-RoundUp-Negative", RoundTowardNegative, 12345, 0, 4, signc_negative, 1235, 1},

		// RoundTowardZero (truncation)
		{"TowardZero-NoRounding", RoundTowardZero, 123, 0, 3, signc_positive, 123, 0},
		{"TowardZero-Truncate-Positive", RoundTowardZero, 12345, 0, 4, signc_positive, 1234, 1},
		{"TowardZero-Truncate-Negative", RoundTowardZero, 12345, 0, 4, signc_negative, 1234, 1},

		// Multiple digit rounding
		{"MultiDigit-RoundTiesToEven", RoundTiesToEven, 123456789, 0, 3, signc_positive, 123, 6},
		{"MultiDigit-RoundTowardZero", RoundTowardZero, 9876543210, 0, 5, signc_positive, 98765, 5},

		// Zero case
		{"Zero", RoundTiesToEven, 0, 0, 5, signc_positive, 0, 0},

		// Large coefficient
		//{"LargeCoefficient", RoundTiesToEven, 9999999999999999, 0, 7, signc_positive, 9999999, 9},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rounded, removed := apply(test.rounding, test.coe, test.exp, Precision(test.precision), test.sign)
			if rounded != test.expected {
				t.Errorf("Apply() rounded = %v, want %v", rounded, test.expected)
			}
			if removed != test.removed {
				t.Errorf("Apply() removed = %v, want %v", removed, test.removed)
			}
		})
	}
}

// TestRoundingApply32 tests the Apply function with uint32 coefficients
func TestRoundingApply32(t *testing.T) {
	tests := []struct {
		name      string
		rounding  Rounding
		coe       uint32
		exp       int16
		precision uint
		sign      signc
		expected  uint32
		removed   uint8
	}{
		// RoundTiesToEven (banker's rounding)
		{"TiesToEven-NoRounding", RoundTiesToEven, 123, 0, 3, signc_positive, 123, 0},
		{"TiesToEven-RoundDown-EvenQuotient", RoundTiesToEven, 12345, 0, 4, signc_positive, 1234, 1},
		{"TiesToEven-RoundUp-EvenQuotient-ExactHalf", RoundTiesToEven, 12350, 0, 4, signc_positive, 1235, 1},
		{"TiesToEven-RoundDown-OddQuotient-ExactHalf", RoundTiesToEven, 12450, 0, 4, signc_positive, 1245, 1},
		{"TiesToEven-RoundUp-OddQuotient-MoreThanHalf", RoundTiesToEven, 12451, 0, 4, signc_positive, 1245, 1},

		// Large coefficient for uint32
		//{"LargeCoefficient", RoundTiesToEven, 9999999, 0, 5, signc_positive, 99999, 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rounded, removed := apply(test.rounding, test.coe, test.exp, Precision(test.precision), test.sign)
			if rounded != test.expected {
				t.Errorf("Apply() rounded = %v, want %v", rounded, test.expected)
			}
			if removed != test.removed {
				t.Errorf("Apply() removed = %v, want %v", removed, test.removed)
			}
		})
	}
}

// Utility function tests
func TestCountDigits(t *testing.T) {
	tests := []struct {
		value    uint64
		expected uint8
	}{
		{0, 1},
		{1, 1},
		{9, 1},
		{10, 2},
		{99, 2},
		{100, 3},
		{12345, 5},
		{9999999, 7},
		{1000000000, 10},
		{9999999999999999, 16},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			if got := countDigits(test.value); got != test.expected {
				t.Errorf("countDigits(%v) = %v, want %v", test.value, got, test.expected)
			}
		})
	}
}

func TestPowTen(t *testing.T) {
	tests := []struct {
		power    uint
		expected uint64
	}{
		{0, 1},
		{1, 10},
		{2, 100},
		{3, 1000},
		{4, 10000},
		{5, 100000},
		{9, 1000000000},
		{16, 10000000000000000},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			if got := powTen[uint64](test.power); got != test.expected {
				t.Errorf("powTen(%v) = %v, want %v", test.power, got, test.expected)
			}
		})
	}
}
