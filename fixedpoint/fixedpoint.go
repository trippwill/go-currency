// Package fixedpoint implements the IEEE 754-2008 Decimal Floating-Point Arithmetic standard
// using the Binary Integer Decimal (BID) encoding.
package fixedpoint

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
