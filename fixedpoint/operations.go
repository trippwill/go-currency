package fixedpoint

import (
	"math"
	"strconv"
)

// FixedPointOperations defines a set of arithmetic operations for fixed point numbers.
type FixedPointOperations interface {
	// Neg returns the negation of the FixedPoint.
	Neg(*Context) FixedPoint
	// Add returns the sum of this FixedPoint and another.
	Add(FixedPoint, *Context) FixedPoint
	// Sub returns the difference between this FixedPoint and another.
	Sub(FixedPoint, *Context) FixedPoint
	// Mul returns the product of this FixedPoint and another.
	Mul(FixedPoint, *Context) FixedPoint
	// Div returns the quotient of this FixedPoint divided by another.
	Div(FixedPoint, *Context) FixedPoint
	// Abs returns the absolute value of this FixedPoint.
	Abs(*Context) FixedPoint
	// Compare compares this FixedPoint with another.
	// It returns -1 if this FixedPoint is less than the other, 0 if they are equal,
	// and 1 if this FixedPoint is greater than the other.
	Compare(FixedPoint) int
}

var (
	_ FixedPointOperations = (*FiniteNumber)(nil)
	_ FixedPointOperations = (*Infinity)(nil)
	_ FixedPointOperations = (*NaN)(nil)
)

func (a *FiniteNumber) Add(right FixedPoint, ctx *Context) FixedPoint {
	if res, ok := val_operands(a, right, ctx); !ok {
		return res
	}

	switch b := right.(type) {
	case *Infinity:
		return b.Add(a, ctx)

	case *FiniteNumber:
		// Align exponents: choose the smaller exponent for maximum precision.
		min_exp := min(b.exp, a.exp)

		// Scale each coefficient to the common exponent.
		// For a, we need to multiply by 10^(a.exp - minExp)
		// For b, we need to multiply by 10^(b.exp - minExp)
		x_coe, xok := scale_coe(a.coe, a.exp-min_exp)
		y_coe, yok := scale_coe(b.coe, b.exp-min_exp)
		if !xok || !yok {
			ctx.signals |= SignalOverflow
			return new_qnan()
		}

		var res_coe coefficient
		var res_sign bool
		var ok bool

		// When signs are identical, do simple addition.
		if a.sign == b.sign {
			res_coe, ok = safe_add(x_coe, y_coe)
			res_sign = a.sign
		} else {
			// When signs differ, subtract the smaller magnitude from the larger.
			if x_coe >= y_coe {
				res_coe, ok = safe_sub(x_coe, y_coe)
				res_sign = a.sign
			} else {
				res_coe, ok = safe_sub(y_coe, x_coe)
				res_sign = b.sign
			}
			// Zero result should always be positive.
			if res_coe == 0 {
				res_sign = false
			}
		}

		// Overflow during addition or subtraction.
		if !ok || res_coe_overflow(res_coe) {
			ctx.signals |= SignalOverflow
		}

		return apply_rounding(
			new(FiniteNumber).Init(
				res_sign,
				res_coe,
				min_exp),
			ctx)
	}

	panic(a)
}

func (a *Infinity) Add(right FixedPoint, ctx *Context) FixedPoint {
	if res, ok := val_operands(a, right, ctx); !ok {
		return res
	}

	switch b := right.(type) {
	case *Infinity:
		// Infinity + Infinity is invalid operation if signs match
		if a.sign == b.sign {
			ctx.signals |= SignalInvalidOperation
			return new_snan()
		}
		// Infinity plus opposite infinity is Zero
		return Zero.Clone()
	case *FiniteNumber:
		// Infinity plus finite number is infinity if signs match
		if a.sign == b.sign {
			return a.Clone()
		}
		ctx.signals |= SignalInvalidOperation
		return new_snan()
	}

	panic(a)
}

func (a *NaN) Add(right FixedPoint, ctx *Context) FixedPoint {
	if res, ok := val_operands(a, right, ctx); !ok {
		return res
	}

	panic(a)
}

func (a *FiniteNumber) Sub(b FixedPoint, ctx *Context) FixedPoint {
	b_neg := b.Neg(ctx)
	return a.Add(b_neg, ctx)
}

func (a *Infinity) Sub(right FixedPoint, ctx *Context) FixedPoint {
	if res, ok := val_operands(a, right, ctx); !ok {
		return res
	}

	switch b := right.(type) {
	case *Infinity:
		// Infinity - Infinity is invalid operation if signs match
		if a.sign == b.sign {
			return new(NaN).Init(false, 2)
		}
		// Infinity minus opposite infinity is infinity
		return a
	case *FiniteNumber:
		// Infinity minus finite number is infinity
		return a
	}

	panic(a)
}

func (a *NaN) Sub(right FixedPoint, ctx *Context) FixedPoint {
	if res, ok := val_operands(a, right, ctx); !ok {
		return res
	}
	panic(a)
}

func (a *FiniteNumber) Mul(right FixedPoint, ctx *Context) FixedPoint {
	if res, ok := val_operands(a, right, ctx); !ok {
		return res
	}

	switch b := right.(type) {
	case *FiniteNumber:
		resCoe, ok := safe_mul(a.coe, b.coe)
		if !ok || res_coe_overflow(resCoe) {
			ctx.signals |= SignalOverflow
			return new_qnan()
		}

		result := new(FiniteNumber).Init(mul_sign(a.sign, b.sign), resCoe, mul_exp(a.exp, b.exp))
		return apply_rounding(result, ctx)

	case *Infinity:
		if a.coe == 0 {
			ctx.signals |= SignalInvalidOperation
			return new_snan()
		}
		return new(Infinity).Init(mul_sign(a.sign, b.sign))
	}

	panic(a)
}

func (a *Infinity) Mul(right FixedPoint, ctx *Context) FixedPoint {
	if res, ok := val_operands(a, right, ctx); !ok {
		return res
	}

	switch b := right.(type) {
	case *FiniteNumber:
		if b.coe == 0 {
			ctx.signals |= SignalInvalidOperation
			return new_snan()
		}
		return new(Infinity).Init(mul_sign(a.sign, b.sign))

	case *Infinity:
		return new(Infinity).Init(mul_sign(a.sign, b.sign))
	}

	panic(a)
}

func (a *NaN) Mul(right FixedPoint, ctx *Context) FixedPoint {
	if res, ok := val_operands(a, right, ctx); !ok {
		return res
	}
	return a
}

func res_coe_overflow(coe coefficient) bool { return coe > fp_coe_max_val }
func mul_sign(a_sign, b_sign bool) bool     { return a_sign != b_sign }

func mul_exp(a_exp, b_exp exponent) exponent {
	if a_exp > b_exp {
		return a_exp + b_exp
	}
	return b_exp + a_exp
}

func (a *FiniteNumber) Div(right FixedPoint, ctx *Context) FixedPoint {
	if res, ok := val_operands(a, right, ctx); !ok {
		return res
	}

	switch b := right.(type) {
	case *Infinity:
		return Zero.Clone()

	case *FiniteNumber:
		if b.coe == 0 {
			if a.coe == 0 {
				ctx.signals |= SignalDivisionImpossible
				return new_snan()
			}
			return new(Infinity).Init(a.sign != b.sign)
		}

		// Perform division with scaling to maintain precision
		dividend := a.coe
		divisor := b.coe
		adjust := int(a.exp - b.exp)

		// Scale dividend to maintain maximum precision
		scale_factor := coefficient(1)
		max_scale := fp_coe_max_val / divisor

		// Scale up dividend as much as possible without overflow
		for scale_factor < coefficient(math.Pow10(int(ctx.precision))) && scale_factor < max_scale {
			scale_factor *= 10
			adjust--
		}

		dividend, ok := safe_mul(dividend, scale_factor)
		if !ok {
			ctx.signals |= SignalOverflow
		}

		quotient := dividend / divisor
		remainder := dividend % divisor

		// Apply proper rounding according to context
		if remainder != 0 {
			switch ctx.rounding {
			case RoundHalfUp:
				if remainder*2 >= divisor {
					quotient++
				}
			case RoundHalfEven:
				if remainder*2 > divisor {
					quotient++
				}
				if remainder*2 == divisor && quotient%2 == 1 {
					quotient++
				}
			case RoundDown:
				// Truncate
			case RoundCeiling:
				if !a.sign && remainder != 0 {
					quotient++
				}
			case RoundFloor:
				if a.sign && remainder != 0 {
					quotient++
				}
			default:
				if remainder*2 >= divisor {
					quotient++
				}
			}
		}

		if quotient > fp_coe_max_val {
			ctx.signals |= SignalOverflow
		}

		return apply_rounding(
			new(FiniteNumber).Init(
				a.sign != b.sign,
				coefficient(quotient),
				exponent(adjust)),
			ctx)
	}

	panic(a)
}

func (a *Infinity) Div(right FixedPoint, ctx *Context) FixedPoint {
	if res, ok := val_operands(a, right, ctx); !ok {
		return res
	}

	switch b := right.(type) {
	case *FiniteNumber:
		if b.coe == 0 {
			ctx.signals |= SignalInvalidOperation
			return new_snan()
		}
		return new(Infinity).Init(a.sign != b.sign)

	case *Infinity:
		ctx.signals |= SignalInvalidOperation
		return new_snan()
	}

	panic(a)
}

func (a *NaN) Div(right FixedPoint, ctx *Context) FixedPoint {
	if res, ok := val_operands(a, right, ctx); !ok {
		return res
	}
	panic(a)
}

// Neg returns the negation of the FixedPoint.
func (a *FiniteNumber) Neg(ctx *Context) FixedPoint {
	if a == nil {
		ctx.signals |= SignalInvalidOperation
		return new_snan()
	}

	return new(FiniteNumber).Init(!a.sign, a.coe, a.exp)
}

func (a *Infinity) Neg(ctx *Context) FixedPoint {
	if a == nil {
		ctx.signals |= SignalInvalidOperation
		return new_snan()
	}

	return new(Infinity).Init(!a.sign)
}

func (a *NaN) Neg(ctx *Context) FixedPoint {
	if a == nil {
		ctx.signals |= SignalInvalidOperation
		return new_snan()
	}

	return new(NaN).Init(false, 2)
}

func (a *FiniteNumber) Abs(ctx *Context) FixedPoint {
	if a == nil {
		ctx.signals |= SignalInvalidOperation
		return new_snan()
	}

	return new(FiniteNumber).Init(false, a.coe, a.exp)
}

func (a *Infinity) Abs(ctx *Context) FixedPoint {
	if a == nil {
		ctx.signals |= SignalInvalidOperation
		return new_snan()
	}

	return new(Infinity).Init(false)
}

func (a *NaN) Abs(ctx *Context) FixedPoint {
	if a == nil {
		ctx.signals |= SignalInvalidOperation
		return new_snan()
	}

	return a.Clone()
}

func (a *FiniteNumber) Compare(b FixedPoint) int {
	switch b := b.(type) {
	case *NaN:
		return 1
	case *Infinity:
		if a.sign {
			if b.sign {
				return 1
			}
			return -1
		} else {
			if b.sign {
				return 1
			}
			return -1
		}
	case *FiniteNumber:
		// Zeroes compare equal irrespective of sign.
		if a.IsZero() && b.IsZero() {
			return 0
		}
		// Different signs: negative < positive.
		if a.sign != b.sign {
			if a.sign {
				return -1
			}
			return 1
		}
		// Align exponents for comparison.
		min_exp := min(a.exp, b.exp)
		a_oe, a_ok := scale_coe(a.coe, a.exp-min_exp)
		b_coe, b_ok := scale_coe(b.coe, b.exp-min_exp)
		switch {
		case !a_ok || !b_ok:
			if a_oe < b_coe {
				if a.sign {
					return 1
				}
				return -1
			} else if a_oe > b_coe {
				if a.sign {
					return -1
				}
				return 1
			}
			return 0

		case a_oe < b_coe:
			if a.sign {
				return 1
			}
			return -1
		case a_oe > b_coe:
			if a.sign {
				return -1
			}
			return 1
		default:
			return 0
		}
	}

	panic(a)
}

func (a *Infinity) Compare(b FixedPoint) int {
	switch b := b.(type) {
	case *NaN:
		return 1
	case *Infinity:
		if a.sign == b.sign {
			return 0
		}
		if a.sign {
			return -1
		}
		return 1
	case *FiniteNumber:
		if a.sign {
			return 1
		}
		return -1
	}

	panic(a)
}

func (a *NaN) Compare(b FixedPoint) int {
	// NaN compares equal to another NaN; otherwise, it is considered less.
	if _, ok := b.(*NaN); ok {
		return 0
	}
	return -1
}

func new_snan() FixedPoint {
	return new(NaN).Init(true, 3)
}

func new_qnan() FixedPoint {
	return new(NaN).Init(false, 3)
}

// scale_coe always adjusts the coefficient to the desired exponent without losing precision.
func scale_coe(c coefficient, diff exponent) (coefficient, bool) {
	if diff > 0 {
		// Multiply absolute value by 10 for each increment in diff.
		for i := exponent(0); i < diff; i++ {
			// Check for multiplication overflow.
			if c > fp_coe_max_val/10 {
				return 0, false
			}
			c *= 10
		}
	} else if diff < 0 {
		// Divide absolute value by 10 for each decrement in diff,
		// ensuring no remainder is lost.
		for i := diff; i < 0; i++ {
			if c%10 != 0 {
				return 0, false
			}
			c /= 10
		}
	}
	return c, true
}

func safe_add[C ~uint64](x, y C) (C, bool) {
	if x > C(fp_coe_max_val)-y {
		return 0, false
	}

	return x + y, true
}

func safe_sub[C ~uint64](x, y C) (C, bool) {
	if x < y {
		return 0, false
	}

	return x - y, true
}

func safe_mul[C ~uint64](x, y C) (C, bool) {
	if x == 0 || y == 0 {
		return 0, true
	}
	if x > C(fp_coe_max_val)/y {
		return 0, false
	}
	return x * y, true
}

// val_operands checks if the operands are valid for arithmetic operations.
// if operands are valid for computation, it returns nil, true.
// if operands are invalid, it returns one of the original NaN values and false.
// if either operand is nil, it returns a new NaN and false.
func val_operands(a, b FixedPoint, ctx *Context) (FixedPoint, bool) {
	if a == nil || b == nil {
		ctx.signals |= SignalInvalidOperation
		return new(NaN).Init(true, 2), false
	}

	// The result of any arithmetic operation which has an operand which is a NaN (a quiet NaN or a signaling NaN)
	// is [s,qNaN] or [s,qNaN,d]. The sign and any diagnostic information is copied from the first operand which
	// is a signaling NaN, or if neither is signaling then from the first operand which is a NaN.
	// Whenever a result is a NaN, the sign of the result depends only on the copied operand (the following rules do not apply).

	a_nan, a_ok := a.(*NaN)
	b_nan, b_ok := b.(*NaN)

	switch {
	case a_ok && b_ok && a_nan.signaling:
		ctx.signals |= SignalInvalidOperation
		return a.Clone(), false
	case a_ok && b_ok && b_nan.signaling:
		ctx.signals |= SignalInvalidOperation
		return b.Clone(), false
	case a_ok && b_ok:
		return a.Clone(), false
	case a_ok && !b_ok:
		return a.Clone(), false
	case !a_ok && b_ok:
		return b.Clone(), false
	}
	return nil, true
}

// apply_rounding rounds a FiniteNumber based on context.precision and context.rounding.
// the resulting coefficient is adjusted to the desired precision.
func apply_rounding(fn *FiniteNumber, ctx *Context) FixedPoint {
	digits := dlen(fn.coe)
	prec := int(ctx.precision)
	if digits == prec {
		return fn
	}

	switch {
	case digits == 0:
		// If the coefficient is zero, no rounding is needed.
		return fn
	case digits == prec:
		// If the number of digits is equal to precision, no rounding is needed.
		return fn
	case digits < prec:
		// If the number of digits is less than precision, scale up.
		required_mult := coefficient(1)
		diff := prec - digits
		// Scale up the coefficient by 10^(prec - digits)
		for range diff {
			if fn.coe > fp_coe_max_val/10 {
				return fn
			}
			required_mult *= 10
		}

		if res_coe, ok := safe_mul(fn.coe, required_mult); ok {
			fn.coe = res_coe
			fn.exp -= exponent(diff)
		}

		return fn
	case digits > prec:
		// If the number of digits is greater than precision, scale down.
		drop := digits - prec
		divisor := uint64(1)
		for range drop {
			divisor *= 10
		}

		quotient := fn.coe / coefficient(divisor)
		remainder := fn.coe % coefficient(divisor)

		switch ctx.rounding {
		case RoundHalfUp:
			if uint64(remainder)*2 >= divisor {
				quotient++
			}
		case RoundHalfEven:
			if uint64(remainder)*2 > divisor {
				quotient++
			} else if uint64(remainder)*2 == divisor {
				if quotient%2 == 1 {
					quotient++
				}
			}
		case RoundDown:
			// truncate
		case RoundCeiling:
			if !fn.sign && remainder != 0 {
				quotient++
			}
		case RoundFloor:
			if fn.sign && remainder != 0 {
				quotient++
			}
		default:
			if uint64(remainder)*2 >= divisor {
				quotient++
			}
		}

		new_exp := fn.exp + exponent(drop)
		fn.coe = quotient
		fn.exp = new_exp

		return fn
	}

	panic(fn)
}

func dlen[C ~uint64](c C) int {
	switch {
	case c == 0:
		return 1
	case c < 10:
		return 1
	case c < 100:
		return 2
	case c < 1000:
		return 3
	case c < 10000:
		return 4
	case c < 100000:
		return 5
	case c < 1000000:
		return 6
	case c < 10000000:
		return 7
	case c < 100000000:
		return 8
	case c < 1000000000:
		return 9
	case c < 10000000000:
		return 10
	case c < 100000000000:
		return 11
	case c < 1000000000000:
		return 12
	case c < 10000000000000:
		return 13
	case c < 100000000000000:
		return 14
	case c < 1000000000000000:
		return 15
	case c < 10000000000000000:
		return 16
	case c < 100000000000000000:
		return 17
	case c < 1000000000000000000:
		return 18
	case c < 10000000000000000000:
		return 19
	}

	return len(strconv.FormatUint(uint64(c), 10))
}
