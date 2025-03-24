// This is an implementation of decimal floating-point arithmetic
// Based on the General Decimal Arithmetic Specification: http://speleotrove.com/decimal/decarith.html
package fixedpoint

import (
	"errors"
)

var g_context = DefaultContext

type Context struct {
	condition SIG
	precision uint64
	rounding  ROUND
}

var (
	DefaultContext  = new(Context).Init(9999999999999999999, ROUND_HALF_UP)
	BasicContext    = new(Context).Init(9999999999999999999, ROUND_HALF_UP)
	ExtendedContext = new(Context).Init(9999999999999999999, ROUND_HALF_UP)
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

type ROUND int

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

type SIG int

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

func GetContext() *Context    { return g_context }
func SetContext(ctx *Context) { g_context = ctx }

func (ctx *Context) Init(precision uint64, rounding ROUND) *Context {
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
