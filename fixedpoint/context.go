package fixedpoint

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Context[X X64 | X32] interface {
	Parse(s string) X
	HandleSignals(original, fallback X) X
	ClearSignals()
	Signal() Signal
	Traps() Signal
	Precision() Precision
	Rounding() Rounding
}

var (
	_ Context[X64] = (*Context64)(nil)
	_ Context[X32] = (*Context32)(nil)
)

// Context64 represents the context for computing 64-bit decimal floating-point numbers.
type Context64 struct {
	context
}

// Context32 represents the context for computing 32-bit decimal floating-point numbers.
type Context32 struct {
	context
}

// context holds the width-independent elements of the context.
type context struct {
	traps     Signal    // The current signal traps.
	signals   Signal    // The current signal state.
	precision Precision // The precision (number of significant digits).
	rounding  Rounding  // The rounding mode.
	locale    Locale    // The locale settings.
}

type Locale struct {
	decimals  string // The decimal separators.
	thousands string // The thousands separators.
}

var DefaultLocale = Locale{
	decimals:  ".",
	thousands: ",_",
}

// Default Basic Context Values.
const (
	BasicRounding Rounding = DefaultRoundingMode
	BasicTraps    Signal   = SignalInvalidOperation | SignalOverflow | SignalUnderflow
)

// Default Extended Context values.
const ()

var (
	ErrUnsupportedPrecision = fmt.Errorf("unsupported precision")
	ErrUnknownRounding      = fmt.Errorf("unknown rounding mode")
)

func NewContext64(precision Precision, rounding Rounding, traps Signal, locale Locale) (*Context64, error) {
	context, err := newContext(precision, rounding, traps, locale, PrecisionMaximum64)
	if err != nil {
		return nil, err
	}

	return &Context64{
		context: context,
	}, nil
}

func NewContext32(precision Precision, rounding Rounding, traps Signal, locale Locale) (*Context32, error) {
	context, err := newContext(precision, rounding, traps, locale, PrecisionMaximum32)
	if err != nil {
		return nil, err
	}

	return &Context32{
		context: context,
	}, nil
}

// BasicContext32 returns a basic context with default values.
func BasicContext32() *Context32 {
	c, err := NewContext32(PrecisionDefault32, BasicRounding, BasicTraps, DefaultLocale)
	if err != nil {
		panic(err)
	}

	return c
}

// BasicContext64 returns a basic context with default values.
func BasicContext64() *Context64 {
	c, err := NewContext64(PrecisionDefault64, BasicRounding, BasicTraps, DefaultLocale)
	if err != nil {
		panic(err)
	}

	return c
}

// Parse converts a string into a FixedPoint value.
// It handles special values (e.g., "NaN", "Infinity") and parses finite numbers.
func (ctx *Context64) Parse(s string) X64 {
	if ctx == nil {
		ctx = BasicContext64()
	}

	sign, kind, coe, exp, signals := parseInput(&ctx.context, s, maxCoefficient64, eMax64)
	ctx.signals |= signals
	if kind != kind_finite {
		return newSpecial64(sign, kind)
	}

	var a X64
	err := a.pack(kind_finite, sign, exp, uint64(coe))
	if err != nil {
		ctx.signals |= SignalConversionSyntax
		return newSpecial64(signc_positive, kind_signaling)
	}

	err = a.Round(ctx.rounding, ctx.precision)
	if err != nil {
		log.Printf("Rounding error: %v", err)
		ctx.signals |= SignalInvalidOperation
		return newSpecial64(signc_positive, kind_signaling)
	}

	return a
}

func (ctx *Context32) Parse(s string) X32 {
	if ctx == nil {
		ctx = BasicContext32()
	}

	sign, kind, coe, exp, signals := parseInput(&ctx.context, s, maxCoefficient32, eMax32)
	ctx.signals |= signals
	if kind != kind_finite {
		return newSpecial32(sign, kind)
	}

	// Pack the parsed values into an X32 object.
	var a X32
	err := a.pack(kind_finite, sign, exp, uint32(coe))
	if err != nil {
		ctx.signals |= SignalConversionSyntax
		return newSpecial32(signc_positive, kind_signaling)
	}

	// Apply rounding to the result.
	err = a.Round(ctx.rounding, ctx.precision)
	if err != nil {
		// TODO: Add signal for rounding error
		log.Printf("Rounding error: %v", err)
		ctx.signals |= SignalInvalidOperation
		return newSpecial32(signc_positive, kind_signaling)
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

func (ctx *Context32) Clone(clear bool) *Context32 {
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

// HandleSignals checks the current signal state and returns the appropriate value.
func (ctx *Context64) HandleSignals(original, fallback X64) X64 {
	if ctx == nil {
		panic("Context64 is nil")
	}

	if ctx.signals&ctx.traps != 0 {
		return fallback
	}

	return original
}

// HandleSignals checks the current signal state and returns the appropriate value.
func (ctx *Context32) HandleSignals(original, fallback X32) X32 {
	if ctx == nil {
		panic("Context32 is nil")
	}

	if ctx.signals&ctx.traps != 0 {
		return fallback
	}

	return original
}

func (ctx *Context64) Add(a, b X64) X64 {
	if ctx == nil {
		panic("Context64 is nil")
	}

	akind, asign, aexp, acoe, err := a.unpack()
	if err != nil {
		ctx.signals |= SignalInvalidOperation
		return newSpecial64(signc_positive, kind_signaling)
	}

	if akind == kind_signaling || akind == kind_quiet {
		return a
	}

	bkind, bsign, bexp, bcoe, err := b.unpack()
	if err != nil {
		ctx.signals |= SignalInvalidOperation
		return newSpecial64(signc_positive, kind_signaling)
	}

	if bkind == kind_signaling || bkind == kind_quiet {
		return b
	}

	// Handle infinity
	if akind == kind_infinity || bkind == kind_infinity {
		if asign == bsign {
			return newSpecial64(asign, kind_infinity)
		}
		return newSpecial64(signc_positive, kind_signaling)
	}

	// adjust the coefficients so exp is the same
	if aexp > bexp {
		bexp += aexp - bexp
		bcoe >>= aexp - bexp
	}
	if bexp > aexp {
		aexp += bexp - aexp
		acoe >>= bexp - aexp
	}

	// add or subtract the coefficients according to the signs
	if asign == bsign {
		acoe += bcoe
	} else {
		if acoe > bcoe {
			acoe -= bcoe
		} else {
			acoe = bcoe - acoe
			asign = signc_negative
		}
	}

	var c X64
	err = c.pack(kind_finite, asign, aexp, acoe)
	if err != nil {
		ctx.signals |= SignalInvalidOperation
		return newSpecial64(signc_positive, kind_signaling)
	}

	err = c.Round(ctx.rounding, ctx.precision)
	if err != nil {
		ctx.signals |= SignalInvalidOperation
		return newSpecial64(signc_positive, kind_signaling)
	}

	return c
}

func (ctx *Context64) String() string {
	if ctx == nil {
		return "nil"
	}

	return fmt.Sprintf("Context64{precision: %d, rounding: %d, traps: %d, signals: %d}",
		ctx.precision, ctx.rounding, ctx.traps, ctx.signals)
}

func (ctx *Context32) String() string {
	if ctx == nil {
		return "nil"
	}

	return fmt.Sprintf("Context32{precision: %d, rounding: %d, traps: %d, signals: %d}",
		ctx.precision, ctx.rounding, ctx.traps, ctx.signals)
}

// ClearSignals clears the current signal state of the context.
func (ctx *context) ClearSignals() {
	if ctx != nil {
		ctx.signals = Signal(0)
	}
}

// Signal retrieves the current signal state of the context.
func (ctx *context) Signal() Signal {
	if ctx == nil {
		return SignalInvalidOperation
	}

	return ctx.signals
}

// Traps retrieves the current signal traps of the context.
func (ctx *context) Traps() Signal {
	if ctx == nil {
		return SignalInvalidOperation
	}

	return ctx.traps
}

// Precision retrieves the current precision of the context.
func (ctx *context) Precision() Precision {
	if ctx == nil {
		return Precision(0)
	}

	return ctx.precision
}

// Rounding retrieves the current rounding mode of the context.
func (ctx *context) Rounding() Rounding {
	if ctx == nil {
		return Rounding(0)
	}

	return ctx.rounding
}

func newContext(p Precision, r Rounding, traps Signal, l Locale, maxP Precision) (context, error) {
	if p < PrecisionMinimum || p > maxP {
		return context{}, ErrUnsupportedPrecision
	}
	if r < DefaultRoundingMode || r > MaxRoundingMode {
		return context{}, ErrUnknownRounding
	}

	return context{
		precision: p,
		rounding:  r,
		traps:     traps,
		signals:   Signal(0),
		locale:    l,
	}, nil
}

func parseInput[C uint64 | uint32, E int8 | int16](
	ctx *context,
	s string,
	maxCoefficient C,
	eMax E,
) (signc, kind, C, E, Signal) {
	if ctx == nil {
		return signc_positive, kind_signaling, 0, 0, SignalInvalidOperation
	}

	s = normalizeInput(s, ctx.locale)
	if s == "" {
		return signc_positive, kind_signaling, 0, 0, SignalConversionSyntax
	}

	sign, kind, isSpecial := isSpecial(s)
	if isSpecial {
		return sign, kind, 0, 0, Signal(0)
	}

	sign, digits, exp, ok := getDigitString[E](s)
	if !ok {
		return signc_positive, kind_signaling, 0, 0, SignalConversionSyntax
	}

	value, err := strconv.ParseUint(digits, 10, 64)
	if err != nil {
		return signc_positive, kind_signaling, 0, 0, SignalConversionSyntax
	}

	if value > uint64(maxCoefficient) || exp > eMax {
		return signc_positive, kind_signaling, 0, 0, SignalOverflow
	}

	return sign, kind, C(value), E(exp), Signal(0)
}

func normalizeInput(input string, locale Locale) string {
	// Trim surrounding spaces.
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	input = strings.ToLower(input)
	input = strings.ReplaceAll(input, " ", "")

	for _, sep := range locale.decimals {
		if sep != '.' {
			input = strings.ReplaceAll(input, string(sep), ".")
		}
	}

	for _, sep := range locale.thousands {
		input = strings.ReplaceAll(input, string(sep), "")
	}

	return input
}

func isSpecial(s string) (signc, kind, bool) {
	switch s {
	case "nan", "+nan":
		return signc_positive, kind_quiet, true
	case "-nan":
		return signc_negative, kind_quiet, true
	case "inf", "infinity", "+inf", "+infinity":
		return signc_positive, kind_infinity, true
	case "-inf", "-infinity":
		return signc_negative, kind_infinity, true
	default:
		return signc_error, kind_finite, false
	}
}

func getDigitString[E int8 | int16](s string) (signc, string, E, bool) {
	if s == "" {
		return signc_positive, "", 0, false
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

	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return signc_positive, "", 0, false
	}

	intPart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}
	if intPart == "" && fracPart == "" {
		return signc_positive, "", 0, false
	}

	// Determine the exponent.
	// For example, "123.45" becomes 12345 with an exponent of -2.
	return s_sign, intPart + fracPart, E(-len(fracPart)), true
}
