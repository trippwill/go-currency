// Package fixedpoint implements IEEE 754-2008 decimal floating-point arithmetic.
package fixedpoint

import (
	"fmt"
)

// Rounding defines the rounding modes according to IEEE 754-2008
type Rounding int

const (
	// RoundTiesToEven rounds to the nearest value; if the number falls midway,
	// it is rounded to the nearest value with an even least significant digit.
	// This is the default rounding mode defined in IEEE 754-2008.
	RoundTiesToEven Rounding = iota

	// RoundTiesToAway rounds to the nearest value; if the number falls midway,
	// it is rounded to the nearest value with magnitude away from zero.
	RoundTiesToAway

	// RoundTowardPositive rounds toward positive infinity.
	// Also known as "ceiling" or "up" rounding.
	RoundTowardPositive

	// RoundTowardNegative rounds toward negative infinity.
	// Also known as "floor" or "down" rounding.
	RoundTowardNegative

	// RoundTowardZero rounds toward zero.
	// Also known as "truncation" rounding.
	RoundTowardZero
)

// DefaultRoundingMode is the default rounding mode (RoundTiesToEven)
const DefaultRoundingMode = RoundTiesToEven

// String returns the string representation of the rounding mode.
func (r Rounding) String() string {
	switch r {
	case RoundTiesToEven:
		return "RoundTiesToEven"
	case RoundTiesToAway:
		return "RoundTiesToAway"
	case RoundTowardPositive:
		return "RoundTowardPositive"
	case RoundTowardNegative:
		return "RoundTowardNegative"
	case RoundTowardZero:
		return "RoundTowardZero"
	default:
		return fmt.Sprintf("Rounding(%d)", r)
	}
}

// Debug returns a short string representation of the rounding mode.
func (r Rounding) Debug() string {
	switch r {
	case RoundTiesToEven:
		return "TiE"
	case RoundTiesToAway:
		return "TiA"
	case RoundTowardPositive:
		return "ToP"
	case RoundTowardNegative:
		return "ToN"
	case RoundTowardZero:
		return "ToZ"
	default:
		return fmt.Sprintf("?(%d)", uint8(r))
	}
}

// Apply applies the specified rounding mode to a coefficient to reduce it to the target precision.
// It returns the rounded coefficient and the number of digits removed.
func Apply[E int8 | int16, C uint32 | uint64](mode Rounding, coef C, exp E, precision uint, sign signc) (C, uint) {
	if coef == 0 {
		return 0, 0 // Zero doesn't need rounding
	}

	digits := countDigits(coef)

	// If we're already at or below the target precision, no rounding needed
	if digits <= precision {
		return coef, 0
	}

	// Calculate how many digits need to be removed
	digitsToRemove := digits - precision

	if digitsToRemove == 0 {
		return coef, 0
	}

	// Calculate divisor (10^digitsToRemove)
	var divisor, powerOfTen C = 1, 10
	for i := uint(1); i <= digitsToRemove; i++ {
		divisor *= powerOfTen
	}

	// Calculate half of the divisor for tie-breaking
	halfDivisor := divisor / 2

	// Quotient and remainder
	quotient := coef / divisor
	remainder := coef % divisor

	// Apply the rounding mode
	switch mode {
	case RoundTiesToEven:
		// If remainder is exactly half, round to even
		if remainder == halfDivisor {
			// If quotient is odd, round up to make it even
			if quotient%2 == 1 {
				quotient++
			}
		} else if remainder > halfDivisor {
			// If remainder is more than half, round up
			quotient++
		}
	case RoundTiesToAway:
		// If remainder is half or more, round away from zero
		if remainder >= halfDivisor {
			quotient++
		}
	case RoundTowardPositive:
		// If positive and any remainder, round up
		if sign == signc_positive && remainder > 0 {
			quotient++
		}
	case RoundTowardNegative:
		// If negative and any remainder, round down (more negative)
		if sign == signc_negative && remainder > 0 {
			quotient++
		}
	case RoundTowardZero:
		// Truncate (do nothing, quotient is already truncated)
	}

	return quotient, digitsToRemove
}

// countDigits returns the number of decimal digits in a number.
func countDigits[T uint32 | uint64](n T) uint {
	if n == 0 {
		return 1
	}

	var count uint = 0
	for n > 0 {
		n /= 10
		count++
	}
	return count
}

// powTen returns 10^n for the given n.
func powTen[T uint32 | uint64](n uint) T {
	if n == 0 {
		return 1
	}

	result := T(1)
	for i := uint(0); i < n; i++ {
		result *= 10
	}
	return result
}
