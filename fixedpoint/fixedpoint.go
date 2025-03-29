// Based on the General Decimal Arithmetic Specification 1.70 – 7 Apr 2009
// https://speleotrove.com/decimal/decarith.html
package fixedpoint

import (
	"fmt"
	"strconv"
	"strings"
)

// FixedPoint represents a fixed-point arithmetic type with a wide dynamic range.
// It supports operations on finite numbers, infinities, and NaN (Not-a-Number) values.
type FixedPoint interface {
	fmt.Stringer
	Debug() string
	// FixedPointChecks defines methods for validation and context management.
	FixedPointChecks
	// FixedPointOperations defines arithmetic operations for FixedPoint types.
	FixedPointOperations
	// SetContext sets the precision and rounding mode for the FixedPoint value.
	SetContext(Precision, Rounding) error
	// Signal retrieves the current signal state of the FixedPoint value.
	Signal() Signal
	// Clone creates a deep copy of the FixedPoint value.
	Clone() FixedPoint
}

// FiniteNumber represents a finite fixed-point number with a coefficient, exponent, and sign.
// It also includes a context for precision and rounding mode.
type FiniteNumber struct {
	coe     coefficient // The coefficient (significand) of the number.
	exp     exponent    // The exponent of the number.
	sign    bool        // The sign of the number (true for negative, false for positive).
	context             // The context for precision and rounding.
}

// Infinity represents a positive or negative infinity value.
// It includes a sign and a context for precision and rounding.
type Infinity struct {
	context
	sign bool // The sign of the infinity (true for negative, false for positive).
}

// NaN represents a quiet or signaling Not-a-Number (NaN) value.
// It includes diagnostic information, a sign, and a context.
type NaN struct {
	diag    diagnostic // Diagnostic information for the NaN value.
	sign    bool       // The sign of the NaN (true for negative, false for positive).
	context            // The context for precision and rounding.
}

// Signal represents the signal state of a FixedPoint value.
type Signal uint8

// Precision represents the number of significant digits for a FixedPoint value.
type Precision uint8

// Rounding represents the rounding mode for a FixedPoint value.
type Rounding uint8

// coefficient represents the significand of a finite number.
type coefficient uint64

// exponent represents the exponent of a finite number.
type exponent int16

// diagnostic represents diagnostic information for NaN values.
type diagnostic uint64

// context represents the precision, rounding mode, and signal state of a FixedPoint value.
type context struct {
	signal    Signal    // The current signal state.
	precision Precision // The precision (number of significant digits).
	rounding  Rounding  // The rounding mode.
}

// Predefined FixedPoint values for common constants.
var (
	Zero = FiniteNumber{
		coe:     0,
		exp:     0,
		context: defaultContext,
	}

	NegZero = FiniteNumber{
		sign:    true,
		coe:     0,
		exp:     0,
		context: defaultContext,
	}

	One = FiniteNumber{
		coe:     1,
		exp:     0,
		context: defaultContext,
	}

	NegOne = FiniteNumber{
		sign:    true,
		coe:     1,
		exp:     0,
		context: defaultContext,
	}
)

// defaultContext defines the default precision and rounding mode.
var defaultContext = context{
	precision: 9,           // Default precision is 9 significant digits.
	rounding:  RoundHalfUp, // Default rounding mode is "Round Half Up".
}

// Constants for coefficient and exponent limits.
const (
	fp_coe_max_val   coefficient = 9_999_999_999_999_999_999 // Maximum coefficient value (10^19 - 1).
	fp_coe_max_len               = 19                        // Maximum length of the coefficient.
	fp_exp_limit_val exponent    = 9_999                     // Maximum exponent value.
	fp_exp_limit_len             = 4                         // Maximum length of the exponent.
)

// Init initializes a FiniteNumber with the given sign, coefficient, exponent, and context.
func (fn *FiniteNumber) Init(sign bool, coe coefficient, exp exponent, ctx context) *FiniteNumber {
	fn.sign = sign
	fn.coe = coe
	fn.exp = exp
	fn.context = ctx
	return fn
}

// Init initializes an Infinity with the given sign and context.
func (inf *Infinity) Init(sign bool, ctx context) *Infinity {
	inf.sign = sign
	inf.context = ctx
	return inf
}

// Init initializes a NaN with the given signal and diagnostic information.
func (nan *NaN) Init(signal Signal, diag_skip int) *NaN {
	nan.context.signal = signal
	nan.diag = encodeDiagnosticInfo(getDiagnosticInfo(diag_skip + 1))
	return nan
}

// SetContext sets the precision and rounding mode for the context.
// Returns an error if the precision is out of range (1 to 19).
func (c *context) SetContext(precision Precision, rounding Rounding) error {
	if c == nil {
		return fmt.Errorf("context is nil")
	}

	if precision < 1 || precision > 19 {
		return fmt.Errorf("precision must be between 1 and 19")
	}
	c.precision = precision
	c.rounding = rounding
	return nil
}

// Parse converts a string into a FixedPoint value.
// It handles special values (e.g., "NaN", "Infinity") and parses finite numbers.
func Parse(s string) FixedPoint {
	// Trim surrounding spaces.
	s = strings.TrimSpace(s)
	if s == "" {
		return new(NaN).Init(SignalConversionSyntax, 2)
	}

	// Handle special values (case-insensitive).
	lower := strings.ToLower(s)
	switch lower {
	case "nan", "+nan", "-nan":
		return new(NaN).Init(SignalClear, 2)
	case "inf", "infinity", "+inf", "+infinity":
		return new(Infinity).Init(false, defaultContext)
	case "-inf", "-infinity":
		return new(Infinity).Init(true, defaultContext)
	}

	// Determine sign.
	sign := false
	switch s[0] {
	case '-':
		sign = true
		s = s[1:]
	case '+':
		s = s[1:]
	}

	// Split the input on the decimal point.
	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return new(NaN).Init(SignalConversionSyntax, 2)
	}

	intPart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}
	if intPart == "" && fracPart == "" {
		return new(NaN).Init(SignalConversionSyntax, 2)
	}

	// Concatenate the integer and fractional digits.
	digits := intPart + fracPart

	// Attempt to parse the digits into a uint64.
	value, err := strconv.ParseUint(digits, 10, 64)
	if err != nil {
		return new(NaN).Init(SignalConversionSyntax, 2)
	}

	// Determine the exponent.
	// For example, "123.45" becomes 12345 with an exponent of -2.
	var exp exponent
	if fracPart != "" {
		exp = exponent(-len(fracPart))
	}

	coe := coefficient(value)

	// Check for coefficient overflow.
	if coe > fp_coe_max_val {
		return new(NaN).Init(SignalOverflow, 2)
	}

	return apply_rounding(new(FiniteNumber).Init(sign, coe, exp, defaultContext))
}

// Signal retrieves the current signal state of the context.
func (c *context) Signal() Signal {
	if c == nil {
		return SignalInvalidOperation
	}

	return c.signal
}

// Clone creates a deep copy of a FiniteNumber.
func (a *FiniteNumber) Clone() FixedPoint {
	if a == nil {
		return nil
	}

	return &FiniteNumber{
		sign:    a.sign,
		coe:     a.coe,
		exp:     a.exp,
		context: a.context,
	}
}

// Clone creates a deep copy of an Infinity.
func (a *Infinity) Clone() FixedPoint {
	if a == nil {
		return nil
	}

	return &Infinity{
		sign:    a.sign,
		context: a.context,
	}
}

// Clone creates a deep copy of a NaN.
func (a *NaN) Clone() FixedPoint {
	if a == nil {
		return nil
	}

	return &NaN{
		sign:    a.sign,
		diag:    a.diag,
		context: a.context,
	}
}

func Equals(a, b FixedPoint) bool             { return a.Compare(b) == 0 }
func LessThan(a, b FixedPoint) bool           { return a.Compare(b) < 0 }
func GreaterThan(a, b FixedPoint) bool        { return a.Compare(b) > 0 }
func LessThanOrEqual(a, b FixedPoint) bool    { return a.Compare(b) <= 0 }
func GreaterThanOrEqual(a, b FixedPoint) bool { return a.Compare(b) >= 0 }

func Must(a FixedPoint) FixedPoint {
	if a == nil {
		panic("nil fixed point value")
	}
	return a
}
