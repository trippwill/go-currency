package fixedpoint

import (
	"fmt"
	"strconv"
	"strings"
)

// Context represents the context in which FixedPoint values are computed.
type Context struct {
	traps     Signal    // The current signal traps.
	signals   Signal    // The current signal state.
	precision Precision // The precision (number of significant digits).
	rounding  Rounding  // The rounding mode.
}

type (
	Signal    uint8 // Signal represents the signal state of the context.
	Precision uint8 // Precision represents the number of significant digits in a FixedPoint value.
	Rounding  uint8 // Rounding represents the rounding mode used in the context.
)

const (
	PrecisionDefault Precision = 9
	PrecisionMinimum Precision = 7
	PrecisionMaximum Precision = fp_coe_max_len - fp_exp_limit_len
)

// Default Basic Context Values.
const (
	BasicPrecision Precision = PrecisionDefault
	BasicTraps     Signal    = SignalConversionSyntax | SignalInvalidOperation
	BasicRounding  Rounding  = RoundHalfUp
)

// Default Extended Context values.
const (
	ExtendedPrecision Precision = PrecisionMaximum
	ExtendedTraps     Signal    = SignalConversionSyntax | SignalInvalidOperation | SignalOverflow
	ExtendedRounding  Rounding  = RoundHalfEven
)

// ErrUnsupportedPrecision is an error that indicates an unsupported precision value.
var ErrUnsupportedPrecision = fmt.Errorf("precision must be between %v and %v", PrecisionMinimum, PrecisionMaximum)

// NewContext creates a new context with the specified precision, rounding mode, and enabled traps.
func NewContext(precision Precision, rounding Rounding, traps Signal) (*Context, error) {
	if precision < PrecisionMinimum || precision > PrecisionMaximum {
		return nil, ErrUnsupportedPrecision
	}

	return &Context{
		precision: precision,
		rounding:  rounding,
		traps:     traps,
		signals:   Signal(0),
	}, nil
}

// BasicContext returns a basic context with default values.
func BasicContext() *Context {
	return &Context{
		precision: BasicPrecision,
		rounding:  BasicRounding,
		traps:     BasicTraps,
		signals:   Signal(0),
	}
}

// ExtendedContext returns an extended context with default values.
func ExtendedContext() *Context {
	return &Context{
		precision: ExtendedPrecision,
		rounding:  ExtendedRounding,
		traps:     ExtendedTraps,
		signals:   Signal(0),
	}
}

// Clone creates a copy of the context, optionally clearing the signal state.
func (ctx *Context) Clone(clear bool) *Context {
	if ctx == nil {
		return nil
	}

	signals := ctx.signals
	if clear {
		signals = Signal(0)
	}

	return &Context{
		precision: ctx.precision,
		rounding:  ctx.rounding,
		traps:     ctx.traps,
		signals:   signals,
	}
}

// ClearSignals clears the current signal state of the context.
func (ctx *Context) ClearSignals() {
	ctx.signals = Signal(0)
}

// Signal retrieves the current signal state of the context.
func (ctx *Context) Signal() Signal {
	if ctx == nil {
		return SignalInvalidOperation
	}

	return ctx.signals
}

// Parse converts a string into a FixedPoint value.
// It handles special values (e.g., "NaN", "Infinity") and parses finite numbers.
func (ctx *Context) Parse(s string) FixedPoint {
	// Trim surrounding spaces.
	s = strings.TrimSpace(s)
	if s == "" {
		ctx.signals |= SignalConversionSyntax
		return new_snan()
	}

	// Handle special values (case-insensitive).
	lower := strings.ToLower(s)
	switch lower {
	case "nan", "+nan", "-nan":
		return new_qnan()
	case "inf", "infinity", "+inf", "+infinity":
		return new(Infinity).Init(false)
	case "-inf", "-infinity":
		return new(Infinity).Init(true)
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
		ctx.signals |= SignalConversionSyntax
		return new_snan()
	}

	intPart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}
	if intPart == "" && fracPart == "" {
		ctx.signals |= SignalConversionSyntax
		return new_snan()
	}

	// Concatenate the integer and fractional digits.
	digits := intPart + fracPart

	// Attempt to parse the digits into a uint64.
	value, err := strconv.ParseUint(digits, 10, 64)
	if err != nil {
		ctx.signals |= SignalConversionSyntax
		return new_snan()
	}

	// Determine the exponent.
	// For example, "123.45" becomes 12345 with an exponent of -2.
	exp := int16(-len(fracPart))
	coe := coefficient(value)

	// Check for coefficient overflow.
	if coe > fp_coe_max_val || exp > fp_exp_limit_val {
		ctx.signals |= SignalOverflow
		return new_snan()
	}

	return apply_precision(
		new(FiniteNumber).Init(sign, coe, exp),
		ctx)
}

func (ctx *Context) Add(a, b FixedPoint) FixedPoint {
	switch v := a.(type) {
	case *FiniteNumber:
		return v.Add(b, ctx)
	case *Infinity:
		return v.Add(b, ctx)
	case *NaN:
		return v.Add(b, ctx)
	default:
		return a.Add(b, ctx)
	}
}

// Sub performs subtraction of two FixedPoint values.
func (ctx *Context) Sub(a, b FixedPoint) FixedPoint {
	switch v := a.(type) {
	case *FiniteNumber:
		return v.Sub(b, ctx)
	case *Infinity:
		return v.Sub(b, ctx)
	case *NaN:
		return v.Sub(b, ctx)
	default:
		return a.Sub(b, ctx)
	}
}

// Mul performs multiplication of two FixedPoint values.
func (ctx *Context) Mul(a, b FixedPoint) FixedPoint {
	switch v := a.(type) {
	case *FiniteNumber:
		return v.Mul(b, ctx)
	case *Infinity:
		return v.Mul(b, ctx)
	case *NaN:
		return v.Mul(b, ctx)
	default:
		return a.Mul(b, ctx)
	}
}

// Div performs division of two FixedPoint values.
func (ctx *Context) Div(a, b FixedPoint) FixedPoint {
	switch v := a.(type) {
	case *FiniteNumber:
		return v.Div(b, ctx)
	case *Infinity:
		return v.Div(b, ctx)
	case *NaN:
		return v.Div(b, ctx)
	default:
		return a.Div(b, ctx)
	}
}

// Neg performs negation of a FixedPoint value.
func (ctx *Context) Neg(a FixedPoint) FixedPoint {
	switch v := a.(type) {
	case *FiniteNumber:
		return v.Neg(ctx)
	case *Infinity:
		return v.Neg(ctx)
	case *NaN:
		return v.Neg(ctx)
	default:
		return a.Neg(ctx)
	}
}

// Abs performs absolute value of a FixedPoint value.
func (ctx *Context) Abs(a FixedPoint) FixedPoint {
	switch v := a.(type) {
	case *FiniteNumber:
		return v.Abs(ctx)
	case *Infinity:
		return v.Abs(ctx)
	case *NaN:
		return v.Abs(ctx)
	default:
		return a.Abs(ctx)
	}
}

// Compare compares two FixedPoint values and returns an integer indicating their relative order.
// It returns -1 if a < b, 0 if a == b, and 1 if a > b.
// panics if ctx is nil.
func (ctx *Context) Compare(a, b FixedPoint) int {
	if ctx == nil {
		panic("context is nil")
	}

	switch v := a.(type) {
	case *FiniteNumber:
		return v.Compare(b, ctx)
	case *Infinity:
		return v.Compare(b, ctx)
	case *NaN:
		return v.Compare(b, ctx)
	default:
		return a.Compare(b, ctx)
	}
}

// Comparison functions for FixedPoint values.
func (ctx *Context) Equal(a, b FixedPoint) bool              { return a.Compare(b, ctx) == 0 }
func (ctx *Context) LessThan(a, b FixedPoint) bool           { return a.Compare(b, ctx) < 0 }
func (ctx *Context) GreaterThan(a, b FixedPoint) bool        { return a.Compare(b, ctx) > 0 }
func (ctx *Context) LessThanOrEqual(a, b FixedPoint) bool    { return a.Compare(b, ctx) <= 0 }
func (ctx *Context) GreaterThanOrEqual(a, b FixedPoint) bool { return a.Compare(b, ctx) >= 0 }

// Must panics if the traps set in the context are triggered.
func (ctx *Context) Must(FixedPoint FixedPoint) FixedPoint {
	if ctx.traps&ctx.signals != 0 {
		panic(fmt.Sprintf("signals: %v", ctx.signals))
	}
	return FixedPoint
}

// Handle traps and return the FixedPoint value.
func (ctx *Context) Trap(handler func(*Context, FixedPoint) FixedPoint, a FixedPoint) FixedPoint {
	if ctx.traps&ctx.signals != 0 {
		return handler(ctx, a)
	}
	return a
}

// All applies the specified operation to all FixedPoint values in the args slice.
func All(op func(a, b FixedPoint) FixedPoint, args ...FixedPoint) FixedPoint {
	if len(args) == 0 {
		return nil
	}

	result := args[0]
	for _, arg := range args[1:] {
		result = op(result, arg)
	}

	return result
}

func TestMustAll(ctx *Context, a, b, c, d FixedPoint) FixedPoint {
	return ctx.Must(All(ctx.Add, a, b, c, d))
}
