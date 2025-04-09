// Package fixedpoint implements the IEEE 754-2008 Decimal Floating-Point Arithmetic standard
// using the Binary Integer Decimal (BID) encoding.
package fixedpoint

import (
	"log"
	"unsafe"
)

// Sign represents the sign of a decimal floating-point number
type signc int8

const (
	signc_negative signc = -1 // Negative
	signc_positive signc = 1  // Positive
	signc_error    signc = 0  // Error
)

// Kind represents the type of a decimal floating-point number
type kind uint8

const (
	kind_signaling kind = iota // Signaling NaN
	kind_quiet                 // Quiet NaN
	kind_infinity              // Infinity
	kind_finite                // Finite number
)

// packed is an internal interface for decimal floating-point types
// following the IEEE 754-2008 standard with BID encoding.
type packed[E int8 | int16, C uint32 | uint64] interface {
	// pack converts the components of a decimal floating-point number
	// into its BID encoding representation
	pack(kind kind, sign signc, exp E, coe C) error

	// unpack extracts the components of a decimal floating-point number
	// from its BID encoding representation
	unpack() (kind, signc, E, C, error)

	// isZero returns true if the number is zero (positive or negative)
	isZero() bool

	// isNaN returns true if the number is Not-a-Number (quiet or signaling)
	isNaN() bool

	// isInf returns true if the number is infinity (positive or negative)
	isInf() bool
}

// Create special values for decimal64 and decimal32

// newSpecial64 creates a special value (NaN, Infinity) for decimal64
func newSpecial64(sign signc, kind kind) X64 {
	var res X64
	switch kind {
	case kind_signaling, kind_quiet, kind_infinity:
		if err := res.pack(kind, sign, 0, 0); err != nil {
			panic(err)
		}
	default:
		panic(newInternalError(res, "invalid kind"))
	}
	return res
}

// newSpecial32 creates a special value (NaN, Infinity) for decimal32
func newSpecial32(sign signc, kind kind) X32 {
	var res X32
	switch kind {
	case kind_signaling, kind_quiet, kind_infinity:
		if err := res.pack(kind, sign, 0, 0); err != nil {
			panic(err)
		}
	default:
		panic(newInternalError(res, "invalid kind"))
	}
	return res
}

// quantize64 adjusts the decimal64 value to the target exponent using the specified rounding mode.
// quantize64 implements the IEEE 754-2008 quantize operation.
func quantize64(x X64, expTarget int16, mode Rounding) (X64, Signal) {
	k, sign, exp, coe, err := x.unpack()
	if err != nil || k != kind_finite {
		return x, SignalInvalidOperation // Return the original for special values
	}

	// Calculate the shift amount
	shift := expTarget - exp
	if shift == 0 {
		return x, Signal(0)
	}

	var result X64
	if shift < 0 {
		// Reduce precision (move decimal point right)
		multiplier := pow10[uint64](uint(-shift))
		if coe > maxCoefficient64/multiplier {
			// Overdlow
			return X64{}, SignalOverflow
		}
		coe *= multiplier
	} else {
		// Increase precision (move decimal point left)
		divisor := pow10[uint64](uint(shift))
		quotient, remainder := coe/divisor, coe%divisor
		halfDivisor := divisor / 2

		switch mode {
		case RoundTiesToEven:
			if remainder > halfDivisor || (remainder == halfDivisor && (quotient&1) == 1) {
				quotient++
			}
		case RoundTiesToAway:
			if remainder >= halfDivisor {
				quotient++
			}
		case RoundTowardPositive:
			if remainder > 0 && sign == signc_positive {
				quotient++
			}
		case RoundTowardNegative:
			if remainder > 0 && sign == signc_negative {
				quotient++
			}
		}
		coe = quotient
	}

	err = result.pack(k, sign, expTarget, coe)
	if err != nil {
		log.Println("Error packing result:", err)
		return X64{}, SignalInvalidOperation
	}

	return result, Signal(0)
}

// quantize32 adjusts the decimal32 value to the target exponent using the specified rounding mode.
// quantize32 implements the IEEE 754-2008 quantize operation.
func quantize32(x X32, expTarget int8, mode Rounding) (X32, Signal) {
	k, sign, exp, coe, err := x.unpack()
	if err != nil || k != kind_finite {
		return x, SignalInvalidOperation // Return the original for special values
	}

	// Calculate the shift amount
	shift := expTarget - exp
	if shift == 0 {
		return x, Signal(0)
	}

	var result X32
	if shift < 0 {
		// Reduce precision (move decimal point right)
		multiplier := pow10[uint32](uint(-shift))
		if coe > maxCoefficient32/multiplier {
			// Overflow
			return X32{}, SignalOverflow
		}
		coe *= multiplier
	} else {
		// Increase precision (move decimal point left)
		divisor := pow10[uint32](uint(shift))
		quotient, remainder := coe/divisor, coe%divisor
		halfDivisor := divisor / 2

		switch mode {
		case RoundTiesToEven:
			if remainder > halfDivisor || (remainder == halfDivisor && (quotient&1) == 1) {
				quotient++
			}
		case RoundTiesToAway:
			if remainder >= halfDivisor {
				quotient++
			}
		case RoundTowardPositive:
			if remainder > 0 && sign == signc_positive {
				quotient++
			}
		case RoundTowardNegative:
			if remainder > 0 && sign == signc_negative {
				quotient++
			}
		}
		coe = quotient
	}

	err = result.pack(k, sign, expTarget, coe)
	if err != nil {
		return X32{}, SignalInvalidOperation
	}

	return result, Signal(0)
}

var pow10Lookup = []uint64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
	1000000000000000000,
	10000000000000000000,
}

// pow10 computes 10^n
func pow10[U uint32 | uint64](n uint) (res U) {
	const maxIndex4 = 9
	const maxIndex8 = 19

	switch unsafe.Sizeof(res) {
	case 4:
		if uintptr(n) <= maxIndex4 {
			return U(pow10Lookup[n])
		}
	case 8:
		if uintptr(n) <= maxIndex8 {
			return U(pow10Lookup[n])
		}
	}

	return 0
}
