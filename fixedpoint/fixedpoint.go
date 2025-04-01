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
	coe      coefficient // The coefficient (significand) of the number.
	sign_exp t_sign_exp  // The packed sign and exponent of the number.
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

// t_sign_exp represents the t_sign_exp of a finite number.
type t_sign_exp uint16

// diagnostic represents diagnostic information for NaN values.
type diagnostic uint64

// Predefined FixedPoint values for common constants.
var (
	Zero = FiniteNumber{
		coe:      0,
		sign_exp: pack_sign_exp(false, 0),
	}

	NegZero = FiniteNumber{
		sign_exp: pack_sign_exp(true, 0),
		coe:      0,
	}
)

var (
	_ FixedPoint = (*FiniteNumber)(nil)
	_ FixedPoint = (*Infinity)(nil)
	_ FixedPoint = (*NaN)(nil)
)

// Constants for coefficient and exponent limits.
const (
	fp_coe_max_val   coefficient = 9_999_999_999_999_999_999 // Maximum coefficient value (10^19 - 1).
	fp_coe_max_len               = 19                        // Maximum length of the coefficient.
	fp_exp_limit_val int16       = maxExponent               // Maximum exponent value.
	fp_exp_limit_len             = 4                         // Maximum length of the exponent.
)

// Init initializes a FiniteNumber with the given sign, coefficient, exponent, and context.
func (fn *FiniteNumber) Init(sign bool, coe coefficient, exp int16) *FiniteNumber {
	fn.coe = coe
	fn.sign_exp = pack_sign_exp(sign, exp)
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

const (
	minExponent = -1343 // Minimum exponent value.
	maxExponent = 1344  // Maximum exponent value.
	bias        = 1344  // Bias for the exponent.
)

// pack_sign_exp packs a sign and exponent into a 16-bit unsigned integer.
// - sign: overall sign; false means positive, true means negative.
// - exponent: must be in the range [-1343, 1344].
// The layout is as follows:
//
// [15] | [14:12] | [11:0]
// sign | unused  | biased exponent
func pack_sign_exp(sign bool, exponent int16) t_sign_exp {
	if exponent < minExponent || exponent > maxExponent {
		panic(fmt.Sprintf("exponent %d out of range [%d, %d]", exponent, minExponent, maxExponent))
	}
	// Compute the biased exponent.
	e := uint16(exponent + bias)

	// Determine the sign bit.
	s := uint16(0)
	if sign {
		s = 1
	}

	// Pack the sign in bit 15 and the exponent in bits 0-11.
	return t_sign_exp((s << 15) | e)
}

// unpack_sign_exp unpacks a 16-bit field into its sign and exponent components.
// Returns:
// - sign: false for positive, true for negative.
// - exponent: the original exponent.
func unpack_sign_exp(packed t_sign_exp) (bool, int16) {
	// Extract the sign bit (bit 15).
	s := (packed >> 15) & 0x1

	// Extract the biased exponent (bits 0-11).
	e := int16(packed & 0x0FFF)

	// Reverse the bias to get the original exponent.
	exponent := e - bias
	sign := s == 1

	return sign, exponent
}

// Clone creates a deep copy of a FiniteNumber.
func (a *FiniteNumber) Clone() FixedPoint {
	if a == nil {
		return nil
	}

	return &FiniteNumber{
		coe:      a.coe,
		sign_exp: a.sign_exp,
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
