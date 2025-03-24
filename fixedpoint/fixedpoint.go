// This is an implementation of decimal floating-point arithmetic
// Based on the General Decimal Arithmetic Specification: http://speleotrove.com/decimal/decarith.html
package fixedpoint

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var g_context = DefaultContext

type Context struct {
	condition SIG
	precision int16
	rounding  ROUND
}

type FixedPoint struct {
	s bool   // sign
	c uint64 // coefficient
	e int16  // exponent
	x rune   // special
}

var (
	NaN     = FixedPoint{x: x_Nan}
	SNaN    = FixedPoint{x: x_sNan}
	Inf     = FixedPoint{x: x_Inf}
	NegInf  = FixedPoint{x: x_Inf, s: true}
	Zero    = FixedPoint{}
	NegZero = FixedPoint{s: true}
)

const (
	c_max_digits        = 19
	c_max_value  uint64 = 9_999_999_999_999_999_999
)

const (
	x_Nan  = 'n'
	x_sNan = 'N'
	x_Inf  = 'F'
)

var (
	DefaultContext  = new(Context).Init(c_max_digits, ROUND_HALF_UP)
	BasicContext    = new(Context).Init(c_max_digits, ROUND_HALF_UP)
	ExtendedContext = new(Context).Init(c_max_digits, ROUND_HALF_UP)
)

var (
	ErrFixedPoint         = errors.New("fixed point error")
	ErrClamped            = errors.New("clamped error")
	ErrDivisionByZero     = errors.New("division by zero")
	ErrInexact            = errors.New("inexact error")
	ErrRounded            = errors.New("rounded error")
	ErrSubnormal          = errors.New("subnormal error")
	ErrOverflow           = errors.New("overflow error")
	ErrUnderflow          = errors.New("underflow error")
	ErrFloatOperation     = errors.New("float operation error")
	ErrDivisionImpossible = errors.New("division impossible error")
	ErrInvalidContext     = errors.New("invalid context error")
	ErrConversionSyntax   = errors.New("conversion syntax error")
	ErrDivisionUndefined  = errors.New("division undefined error")
)

type (
	ROUND int
	SIG   int
)

const (
	ROUND_DOWN ROUND = iota
	ROUND_HALF_UP
	ROUND_HALF_EVEN
	ROUND_CEILING
	ROUND_FLOOR
	ROUND_UP
	ROUND_HALF_DOWN
	ROUND_05UP
)

const (
	SIG_NONE SIG = iota
	SIG_CLAMPED
	SIG_DIVISION_BY_ZERO
	SIG_INEXACT
	SIG_OVERFLOW
	SIG_ROUNDED
	SIG_UNDERFLOW
	SIG_INVALID_OPERATION
	SIG_SUBNORMAL
	SIG_FLOAT_OPERATION
)

var condition_map = map[error]SIG{
	ErrConversionSyntax:   SIG_INVALID_OPERATION,
	ErrDivisionImpossible: SIG_INVALID_OPERATION,
	ErrDivisionUndefined:  SIG_INVALID_OPERATION,
	ErrInvalidContext:     SIG_INVALID_OPERATION,
}

func GetGlobalContext() *Context    { return g_context }
func SetGlobalContext(ctx *Context) { g_context = ctx }

func (ctx *Context) Init(precision int16, rounding ROUND) *Context {
	ctx.precision = precision
	ctx.rounding = rounding
	ctx.condition = SIG_NONE
	return ctx
}

func (ctx *Context) setCondition(cond SIG) {
	if cond == 0 || cond < SIG_CLAMPED || cond > SIG_FLOAT_OPERATION {
		panic("invalid condition")
	}

	ctx.condition |= cond
}

func (ctx *Context) clearCondition(cond SIG) {
	ctx.condition &= ^cond
}

func NewFixedPoint(value string) (fp FixedPoint, err error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return FixedPoint{}, ErrConversionSyntax
	}

	switch strings.ToLower(value) {
	case "nan":
		return NaN, nil
	case "snan":
		return SNaN, nil
	case "inf", "+inf":
		return Inf, nil
	case "-inf":
		return NegInf, nil
	}

	var sign bool
	if value[0] == '-' {
		sign = true
		value = value[1:]
	} else if value[0] == '+' {
		value = value[1:]
	}

	parts := strings.SplitN(value, ".", 2)
	intPart := parts[0]
	var fracPart string
	if len(parts) == 2 {
		fracPart = parts[1]
	}

	if len(intPart) == 0 && len(fracPart) == 0 {
		return FixedPoint{}, ErrConversionSyntax
	}

	combined := intPart + fracPart
	if len(combined) > c_max_digits {
		return FixedPoint{}, ErrOverflow
	}

	coeff, err := strconv.ParseUint(combined, 10, 64)
	if err != nil {
		return FixedPoint{}, ErrConversionSyntax
	}

	exp := int16(0)
	if len(fracPart) > 0 {
		exp = -int16(len(fracPart))
	}

	return FixedPoint{
		s: sign,
		c: coeff,
		e: exp,
	}, nil
}

func NewFromFixedPoint(fp FixedPoint) FixedPoint {
	return FixedPoint{
		s: fp.s,
		c: fp.c,
		e: fp.e,
		x: fp.x,
	}
}

func (x FixedPoint) Add(y FixedPoint) FixedPoint {
	return x.AddUsingContext(y, g_context)
}

func (x FixedPoint) AddUsingContext(y FixedPoint, ctx *Context) FixedPoint {
	switch {
	case x.IsNaN() || y.IsNaN():
		ctx.setCondition(SIG_INVALID_OPERATION)
		return SNaN
	case x.IsInf() || y.IsInf():
		if x.IsInf() && y.IsInf() && x.s != y.s {
			ctx.setCondition(SIG_INVALID_OPERATION)
			return SNaN
		}
		if x.IsInf() {
			return x
		}
		return y
	}

	// Normalize exponents
	exp := min(x.e, y.e)
	scaleX := int64(math.Pow10(int(x.e - exp)))
	scaleY := int64(math.Pow10(int(y.e - exp)))

	// Scale coefficients
	coeffX := int64(x.c) * scaleX
	coeffY := int64(y.c) * scaleY

	// Apply signs
	if x.s {
		coeffX = -coeffX
	}
	if y.s {
		coeffY = -coeffY
	}

	// Add coefficients
	resultCoeff := coeffX + coeffY
	resultSign := resultCoeff < 0

	// Take absolute value for the result coefficient
	if resultCoeff < 0 {
		resultCoeff = -resultCoeff
	}

	// Create result FixedPoint
	result := FixedPoint{
		s: resultSign,
		c: uint64(resultCoeff),
		e: exp,
	}

	// Rescale result based on context rounding
	return result.rescale(result.e, ctx.rounding)
}

func (fp *FixedPoint) IsNormal() bool {
	return fp.x == 0 && fp.e >= -c_max_digits && fp.e <= c_max_digits
}

func (fp *FixedPoint) IsZero() bool      { return fp.c == 0 && fp.x == 0 }
func (fp *FixedPoint) IsInteger() bool   { return fp.e == 0 && fp.c != 0 && fp.x == 0 }
func (fp *FixedPoint) IsNegative() bool  { return fp.s }
func (fp *FixedPoint) IsPositive() bool  { return !fp.s && fp.x == 0 }
func (fp *FixedPoint) IsSubnormal() bool { return fp.e < -c_max_digits && fp.x == 0 }
func (fp *FixedPoint) IsSpecial() bool   { return fp.x != 0 }
func (fp *FixedPoint) IsInf() bool       { return fp.x == x_Inf }
func (fp *FixedPoint) IsNaN() bool       { return fp.x == x_Nan || fp.x == x_sNan }
func (fp *FixedPoint) IsFinite() bool    { return fp.x == 0 && fp.c != 0 }

func min(a, b int16) int16 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}

func (fp FixedPoint) rescale(exp int16, round ROUND) FixedPoint {
	if fp.e == exp {
		return fp
	}

	if fp.IsNaN() || fp.IsInf() {
		return fp
	}

	if fp.c == 0 {
		return FixedPoint{}
	}

	if exp < -c_max_digits || exp > c_max_digits {
		fp.x = x_Inf
		fp.s = false
		fp.c = 0
		fp.e = 0
		return fp
	}

	if round == ROUND_DOWN && (fp.e > exp) {
		scaleFactor := int(fp.e - exp)
		fp.c /= uint64(math.Pow10(scaleFactor))
		fp.e = exp
		return fp
	}

	if round == ROUND_UP && (fp.e < exp) {
		scaleFactor := int(exp - fp.e)
		fp.c *= uint64(math.Pow10(scaleFactor))
		fp.e = exp
		return fp
	}

	return fp
}

// normalize a and b to have the same exp and length of coefficient
func normalize(a FixedPoint, b FixedPoint, prec int16) (FixedPoint, FixedPoint) {
	var tmp, other FixedPoint
	if a.e < b.e {
		tmp = a
		other = b
	} else {
		tmp = b
		other = a
	}

	tmp_len := dlen(tmp.c)
	other_len := dlen(other.c)
	exp := tmp.e + min(-1, int16(tmp_len)-prec-2)
	if other_len+uint(other.e)-1 < uint(exp) {
		other.c = 1
		other.e = exp
	}

	tmp.c = tmp.c * uint64(math.Pow10(int(exp-tmp.e)))
	tmp.e = other.e
	return tmp, other
}

// dlen returns the number of digits in the coefficient
// for example, dlen(1234) = 4
// dlen(0) = 1
// dlen(1234567890123456789) = 19
func dlen(i uint64) uint {
	switch {
	case i == 0:
		return 1
	case i > 1 && i < 10:
		return 1
	case i > 9 && i < 100:
		return 2
	case i > 99 && i < 1000:
		return 3
	case i > 999 && i < 10000:
		return 4
	case i > 9999 && i < 100000:
		return 5
	case i > 99999 && i < 1000000:
		return 6
	case i > 999999 && i < 10000000:
		return 7
	case i > 9999999 && i < 100000000:
		return 8
	case i > 99999999 && i < 1000000000:
		return 9
	case i > 999999999 && i < 10000000000:
		return 10
	}

	count := uint(0)

	for i > 0 {
		i /= 10
		count++
	}

	return count
}

// String implements the Stringer interface for FixedPoint.
func (fp FixedPoint) String() string {
	if fp.IsNaN() {
		if fp.x == x_sNan {
			return "sNaN"
		}
		return "NaN"
	}
	if fp.IsInf() {
		if fp.s {
			return "-Inf"
		}
		return "Inf"
	}

	// Build sign
	sign := ""
	if fp.s && fp.c != 0 {
		sign = "-"
	}

	coeffStr := strconv.FormatUint(fp.c, 10)

	if fp.e >= 0 {
		// For non-negative exponent, append zeros
		return sign + coeffStr + strings.Repeat("0", int(fp.e))
	}

	// Negative exponent: insert a decimal point.
	// The number of fraction digits equals -fp.e.
	dotPos := len(coeffStr) + int(fp.e) // calculate insertion position for the decimal point
	if dotPos <= 0 {
		// Not enough digits: pad with zeros on the left.
		intPart := "0"
		fracPart := strings.Repeat("0", -dotPos) + coeffStr
		return sign + intPart + "." + fracPart
	}

	intPart := coeffStr[:dotPos]
	fracPart := coeffStr[dotPos:]
	return sign + intPart + "." + fracPart
}

func (fp *FixedPoint) Debug() string {
	return fmt.Sprintf("FixedPoint{x: %c, sign: %t, coeff: %d, exp: %d}", fp.x, fp.s, fp.c, fp.e)
}

func (fp FixedPoint) Equals(other FixedPoint) bool {
	// Check if both values represent NaN
	if ((fp.x == x_Nan) || (fp.x == x_sNan)) && ((other.x == x_Nan) || (other.x == x_sNan)) {
		return true
	}

	// Check if one is NaN and the other is not
	if ((fp.x == x_Nan) || (fp.x == x_sNan)) || ((other.x == x_Nan) || (other.x == x_sNan)) {
		return false
	}

	// Check if both are Inf with the same sign
	if (fp.x == x_Inf) && (other.x == x_Inf) && fp.s == other.s {
		return true
	}

	// Check if one is Inf and the other is not
	if (fp.x == x_Inf) || (other.x == x_Inf) {
		return false
	}

	// Handle zero and negative zero
	if fp.c == 0 && other.c == 0 {
		return true
	}

	// Compare sign, coefficient, and exponent
	return fp.s == other.s && fp.c == other.c && fp.e == other.e
}
