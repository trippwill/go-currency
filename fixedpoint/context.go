package fixedpoint

import (
	"fmt"
	"strconv"
	"strings"
)

// context represents the context in which FixedPoint values are computed.
type context struct {
	traps     Signal    // The current signal traps.
	signals   Signal    // The current signal state.
	precision Precision // The precision (number of significant digits).
	rounding  Rounding  // The rounding mode.
}

type (
	Signal    uint8 // Signal represents the signal state of the context.
	Precision uint8 // Precision represents the number of significant digits in a FixedPoint value.
)

const (
	SignalInvalidOperation Signal = 1 << iota
	SignalOverflow
	SignalUnderflow
	SignalConversionSyntax
)

const (
	PrecisionDefault Precision = 9
	PrecisionMinimum Precision = 7
	PrecisionMaximum Precision = 11
)

// Default Basic Context Values.
const (
	BasicPrecision Precision = PrecisionDefault
	BasicRounding  Rounding  = DefaultRoundingMode
	BasicTraps     Signal    = SignalInvalidOperation | SignalOverflow | SignalUnderflow
)

// Default Extended Context values.
const ()

var ErrUnsupportedPrecision = fmt.Errorf("unsupported precision")

// NewContext creates a new context with the specified precision, rounding mode, and enabled traps.
func NewContext[C Context64 | Context32](precision Precision, rounding Rounding, traps Signal) (*C, error) {
	if precision < PrecisionMinimum || precision > PrecisionMaximum {
		return nil, ErrUnsupportedPrecision
	}

	return &C{
		context: context{
			precision: precision,
			rounding:  rounding,
			traps:     traps,
			signals:   Signal(0),
		},
	}, nil
}

type Context64 struct {
	context
}

type Context32 struct {
	context
}

// BasicContext returns a basic context with default values.
func BasicContext[C Context32 | Context64]() *C {
	return &C{
		context: context{
			precision: BasicPrecision,
			rounding:  BasicRounding,
			traps:     BasicTraps,
			signals:   Signal(0),
		},
	}
}

// Parse converts a string into a FixedPoint value.
// It handles special values (e.g., "NaN", "Infinity") and parses finite numbers.
func (ctx *Context64) Parse(s string) X64 {
	// Trim surrounding spaces.
	s = strings.TrimSpace(s)
	if s == "" {
		ctx.signals |= SignalConversionSyntax
		return new_special64(signc_positive, kind_signaling)
	}

	// Handle special values (case-insensitive).
	lower := strings.ToLower(s)
	switch lower {
	case "nan", "+nan":
		return new_special64(signc_positive, kind_quiet)
	case "-nan":
		return new_special64(signc_negative, kind_quiet)
	case "inf", "infinity", "+inf", "+infinity":
		return new_special64(signc_positive, kind_infinity)
	case "-inf", "-infinity":
		return new_special64(signc_negative, kind_infinity)
	}

	// Determine s_sign.
	s_sign := signc_positive
	switch s[0] {
	case '-':
		s_sign = signc_negative
		s = s[1:]
	case '+':
		s = s[1:]
	}

	// Split the input on the decimal point.
	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		ctx.signals |= SignalConversionSyntax
		return new_special64(signc_positive, kind_signaling)
	}

	intPart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}
	if intPart == "" && fracPart == "" {
		ctx.signals |= SignalConversionSyntax
		return new_special64(signc_positive, kind_signaling)
	}

	// Concatenate the integer and fractional digits.
	digits := intPart + fracPart

	// Attempt to parse the digits into a uint64.
	value, err := strconv.ParseUint(digits, 10, 64)
	if err != nil {
		ctx.signals |= SignalConversionSyntax
		return new_special64(signc_positive, kind_signaling)
	}

	// Determine the exponent.
	// For example, "123.45" becomes 12345 with an exponent of -2.
	exp := int16(-len(fracPart))
	coe := value

	// Check for coefficient overflow.
	if coe > MaxCoefficient64 || exp > Emax64 {
		ctx.signals |= SignalOverflow
		return new_special64(signc_positive, kind_signaling)
	}

	var a X64
	err = a.pack(kind_finite, s_sign, exp, coe)
	if err != nil {
		ctx.signals |= SignalConversionSyntax
		return new_special64(signc_positive, kind_signaling)
	}

	return a

	// return apply_precision(
	// 	new(FiniteNumber).Init(sign, coe, exp),
	// 	ctx)
}

func (ctx *Context32) Parse(s string) X32 {
	// Trim surrounding spaces.
	s = strings.TrimSpace(s)
	if s == "" {
		ctx.signals |= SignalConversionSyntax
		return new_special32(signc_positive, kind_signaling)
	}

	// Handle special values (case-insensitive).
	lower := strings.ToLower(s)
	switch lower {
	case "nan", "+nan":
		return new_special32(signc_positive, kind_quiet)
	case "-nan":
		return new_special32(signc_negative, kind_quiet)
	case "inf", "infinity", "+inf", "+infinity":
		return new_special32(signc_positive, kind_infinity)
	case "-inf", "-infinity":
		return new_special32(signc_negative, kind_infinity)
	}

	// Determine s_sign.
	s_sign := signc_positive
	switch s[0] {
	case '-':
		s_sign = signc_negative
		s = s[1:]
	case '+':
		s = s[1:]
	}

	// Split the input on the decimal point.
	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		ctx.signals |= SignalConversionSyntax
		return new_special32(signc_positive, kind_signaling)
	}

	intPart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}
	if intPart == "" && fracPart == "" {
		ctx.signals |= SignalConversionSyntax
		return new_special32(signc_positive, kind_signaling)
	}

	// Concatenate the integer and fractional digits.
	digits := intPart + fracPart

	// Attempt to parse the digits into a uint32.
	value, err := strconv.ParseUint(digits, 10, 32)
	if err != nil {
		ctx.signals |= SignalConversionSyntax
		return new_special32(signc_positive, kind_signaling)
	}

	// Determine the exponent.
	// For example, "123.45" becomes 12345 with an exponent of -2.
	exp := int8(-len(fracPart))
	coe := uint32(value)

	// Check for coefficient overflow.
	if coe > MaxCoefficient32 || exp > Emax32 {
		ctx.signals |= SignalOverflow
		return new_special32(signc_positive, kind_signaling)
	}

	var a X32
	err = a.pack(kind_finite, s_sign, exp, coe)
	if err != nil {
		ctx.signals |= SignalConversionSyntax
		return new_special32(signc_positive, kind_signaling)
	}

	return a
}

// Clone creates a copy of the context, optionally clearing the signal state.
func (ctx *Context64) Clone(clear bool) *Context64 {
	if ctx == nil {
		return nil
	}

	signals := ctx.signals
	if clear {
		signals = Signal(0)
	}

	return &Context64{
		context: context{
			precision: ctx.precision,
			rounding:  ctx.rounding,
			traps:     ctx.traps,
			signals:   signals,
		},
	}
}

func (ctx *Context32) CLone(clear bool) *Context32 {
	if ctx == nil {
		return nil
	}

	signals := ctx.signals
	if clear {
		signals = Signal(0)
	}

	return &Context32{
		context: context{
			precision: ctx.precision,
			rounding:  ctx.rounding,
			traps:     ctx.traps,
			signals:   signals,
		},
	}
}

// ClearSignals clears the current signal state of the context.
func (ctx *Context64) ClearSignals() {
	ctx.signals = Signal(0)
}

// Signal retrieves the current signal state of the context.
func (ctx *Context64) Signal() Signal {
	if ctx == nil {
		return SignalInvalidOperation
	}

	return ctx.signals
}
