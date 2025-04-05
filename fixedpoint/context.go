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

type Context64 struct {
	Context
}

type Context32 struct {
	Context
}

type (
	Signal    uint8 // Signal represents the signal state of the context.
	Precision uint8 // Precision represents the number of significant digits in a FixedPoint value.
	Rounding  uint8 // Rounding represents the rounding mode used in the context.
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

const (
	RoundingDefault Rounding = 0
	RoundingUp      Rounding = 1
	RoundingDown    Rounding = 2
	RoundingHalfUp  Rounding = 3
)

// Default Basic Context Values.
const (
	BasicPrecision Precision = PrecisionDefault
	BasicRounding  Rounding  = RoundingDefault
	BasicTraps     Signal    = SignalInvalidOperation | SignalOverflow | SignalUnderflow
)

// Default Extended Context values.
const ()

var ErrUnsupportedPrecision = fmt.Errorf("unsupported precision")

// NewContext creates a new context with the specified precision, rounding mode, and enabled traps.
func NewContext64(precision Precision, rounding Rounding, traps Signal) (*Context64, error) {
	if precision < PrecisionMinimum || precision > PrecisionMaximum {
		return nil, ErrUnsupportedPrecision
	}

	return &Context64{
		Context: Context{
			precision: precision,
			rounding:  rounding,
			traps:     traps,
			signals:   Signal(0),
		},
	}, nil
}

// BasicContext returns a basic context with default values.
func BasicContext() *Context64 {
	return &Context64{
		Context: Context{
			precision: BasicPrecision,
			rounding:  BasicRounding,
			traps:     BasicTraps,
			signals:   Signal(0),
		},
	}
}

// Clone creates a copy of the context, optionally clearing the signal state.
// func (ctx *Context) Clone(clear bool) *Context {
// 	if ctx == nil {
// 		return nil
// 	}
//
// 	signals := ctx.signals
// 	if clear {
// 		signals = Signal(0)
// 	}
//
// 	return &Context{
// 		precision: ctx.precision,
// 		rounding:  ctx.rounding,
// 		traps:     ctx.traps,
// 		signals:   signals,
// 	}
// }
//
// // ClearSignals clears the current signal state of the context.
// func (ctx *Context) ClearSignals() {
// 	ctx.signals = Signal(0)
// }
//
// // Signal retrieves the current signal state of the context.
// func (ctx *Context) Signal() Signal {
// 	if ctx == nil {
// 		return SignalInvalidOperation
// 	}
//
// 	return ctx.signals
// }

// Parse converts a string into a FixedPoint value.
// It handles special values (e.g., "NaN", "Infinity") and parses finite numbers.
func (ctx *Context64) Parse(s string) X64 {
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
		return new_infinity(sign_positive)
	case "-inf", "-infinity":
		return new_infinity(sign_negative)
	}

	// Determine sg.
	sg := sign_positive
	switch s[0] {
	case '-':
		sg = sign_negative
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
	coe := value

	// Check for coefficient overflow.
	if coe > coe_maxv64 || exp > exp_max64 {
		ctx.signals |= SignalOverflow
		return new_snan()
	}

	var a X64
	// Pack the data into the X64 type.
	err = a.pack(kind_finite, sg, exp, coe)
	if err != nil {
		ctx.signals |= SignalConversionSyntax
		return new_snan()
	}

	return a

	// return apply_precision(
	// 	new(FiniteNumber).Init(sign, coe, exp),
	// 	ctx)
}

func new_snan() X64 {
	var res X64
	if err := res.pack(kind_signaling, sign_positive, 0, 0); err != nil {
		panic(err)
	}
	return res
}

func new_qnan() X64 {
	var res X64
	if err := res.pack(kind_quiet, sign_positive, 0, 0); err != nil {
		panic(err)
	}
	return res
}

func new_infinity(sign sign) X64 {
	var res X64
	if err := res.pack(kind_infinity, sign, 0, 0); err != nil {
		panic(err)
	}
	return res
}
