// Based on the General Decimal Arithmetic Specification 1.70 – 7 Apr 2009
// https://speleotrove.com/decimal/decarith.html
package fixedpoint

import (
	"fmt"
)

// FixedPoint represents a fixed-point arithmetic type with a wide dynamic range.
// It supports operations on finite numbers, infinities, and NaN (Not-a-Number) values.
type FixedPoint interface {
	fmt.Stringer
	Debug() string
	FixedPointChecks
	FixedPointOperations

	// Clone creates a deep copy of the FixedPoint value.
	Clone() FixedPoint
}

// FiniteNumber represents a finite fixed-point number with a coefficient, exponent, and sign.
// It also includes a context for precision and rounding mode.
type FiniteNumber struct {
	coe  coefficient // The coefficient (significand) of the number.
	exp  exponent    // The exponent of the number.
	sign bool        // The sign of the number (true for negative, false for positive).
}

// Infinity represents a positive or negative infinity value.
// It includes a sign and a context for precision and rounding.
type Infinity struct {
	sign bool // The sign of the infinity (true for negative, false for positive).
}

// NaN represents a quiet or signaling Not-a-Number (NaN) value.
// It includes diagnostic information, a sign, and a context.
type NaN struct {
	diag      diagnostic // Diagnostic information for the NaN value.
	sign      bool       // The sign of the NaN (true for negative, false for positive).
	signaling bool       // Indicates if the NaN is signaling (true) or quiet (false).
}

// coefficient represents the significand of a finite number.
type coefficient uint64

// exponent represents the exponent of a finite number.
type exponent int16

// diagnostic represents diagnostic information for NaN values.
type diagnostic uint64

// Predefined FixedPoint values for common constants.
var (
	Zero = FiniteNumber{
		coe: 0,
		exp: 0,
	}

	NegZero = FiniteNumber{
		sign: true,
		coe:  0,
		exp:  0,
	}
)

// Constants for coefficient and exponent limits.
const (
	fp_coe_max_val   coefficient = 9_999_999_999_999_999_999 // Maximum coefficient value (10^19 - 1).
	fp_coe_max_len               = 19                        // Maximum length of the coefficient.
	fp_exp_limit_val exponent    = 9_999                     // Maximum exponent value.
	fp_exp_limit_len             = 4                         // Maximum length of the exponent.
)

// Init initializes a FiniteNumber with the given sign, coefficient, exponent, and context.
func (fn *FiniteNumber) Init(sign bool, coe coefficient, exp exponent) *FiniteNumber {
	fn.sign = sign
	fn.coe = coe
	fn.exp = exp
	return fn
}

// Init initializes an Infinity with the given sign and context.
func (inf *Infinity) Init(sign bool) *Infinity {
	inf.sign = sign
	return inf
}

// Init initializes a NaN with the given signal and diagnostic information.
func (nan *NaN) Init(signaling bool, diag_skip int) *NaN {
	nan.signaling = signaling
	nan.diag = encodeDiagnosticInfo(getDiagnosticInfo(diag_skip + 1))
	return nan
}

func Parse(s string, ctx *Context) FixedPoint {
	if ctx == nil {
		panic("nil context")
	}

	return ctx.Parse(s)
}

// Clone creates a deep copy of a FiniteNumber.
func (a *FiniteNumber) Clone() FixedPoint {
	if a == nil {
		return nil
	}

	return &FiniteNumber{
		sign: a.sign,
		coe:  a.coe,
		exp:  a.exp,
	}
}

// Clone creates a deep copy of an Infinity.
func (a *Infinity) Clone() FixedPoint {
	if a == nil {
		return nil
	}

	return &Infinity{
		sign: a.sign,
	}
}

// Clone creates a deep copy of a NaN.
func (a *NaN) Clone() FixedPoint {
	if a == nil {
		return nil
	}

	return &NaN{
		sign:      a.sign,
		diag:      a.diag,
		signaling: a.signaling,
	}
}
