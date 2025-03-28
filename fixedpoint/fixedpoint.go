// Based on the General Decimal Arithmetic Specification 1.70 – 7 Apr 2009
// https://speleotrove.com/decimal/decarith.html
package fixedpoint

import (
	"strconv"
	"strings"
)

// FixedPoint represents a fixed-point arithmetic type with a wide dynamic range.
// Representable ranges:
//   - Positive values: from the smallest positive 1.0e-9999 to the largest positive 9.999999999999999e+9999
//   - Negative values: from the largest negative -9.999999999999999e+9999 to the smallest negative -1.0e-9999
type FixedPoint struct {
	coe  coefficient
	exp  exponent
	sign bool
	flg  flags
	ctx  context
}

type (
	Signal    uint8
	Precision uint8
	Rounding  uint8
)

type (
	coefficient uint64
	exponent    int16
	flags       struct {
		sig Signal
		inf bool
		nan bool
	}
	context struct {
		precision Precision
		rounding  Rounding
	}
)

var Zero = FixedPoint{
	coe: 0,
	exp: 0,
}

var NegZero = FixedPoint{
	sign: true,
	coe:  0,
	exp:  0,
}

var One = FixedPoint{
	coe: 1,
	exp: 0,
}

var NegOne = FixedPoint{
	sign: true,
	coe:  1,
	exp:  0,
}

var defaultContext = context{
	precision: fp_coe_max_len - fp_exp_limit_len,
	rounding:  RoundingNearestEven,
}

// Update constants for 64-bit coefficient
const (
	fp_coe_max_val   coefficient = 9_999_999_999_999_999_999 // 10^19 - 1
	fp_coe_max_len               = 19
	fp_exp_limit_val exponent    = 9_999
	fp_exp_limit_len             = 4
)

func (fp *FixedPoint) Init(sign bool, coe coefficient, exp exponent) *FixedPoint {
	fp.sign = sign
	fp.coe = coe
	fp.exp = exp
	fp.ctx = defaultContext
	return fp
}

func (fp *FixedPoint) Parse(s string) *FixedPoint {
	*fp = Parse(s)
	return fp
}

func (fp *FixedPoint) SetContext(precision Precision, rounding Rounding) *FixedPoint {
	fp.ctx.precision = precision
	fp.ctx.rounding = rounding
	return fp
}

// New returns a new FixedPoint value from a significand and an exponent.
func New(significand int64, exp int16) FixedPoint {
	sign := significand < 0
	abs := coefficient(significand)
	if sign {
		abs = coefficient(-significand)
	}

	if abs > fp_coe_max_val || abs == fp_coe_max_val && exp > 0 {
		return FixedPoint{
			sign: sign,
			coe:  0,
			exp:  exponent(exp),
			flg:  flags{inf: true, sig: SignalOverflow},
		}
	}
	return FixedPoint{
		sign: sign,
		coe:  abs,
		exp:  exponent(exp),
		ctx:  defaultContext,
	}
}

// Parse converts a string into a FixedPoint value.
func Parse(s string) FixedPoint {
	// Trim surrounding spaces.
	s = strings.TrimSpace(s)
	if s == "" {
		return FixedPoint{
			coe: 0,
			exp: 0,
			flg: flags{nan: true, sig: SignalConversionSyntax},
		}
	}

	// Handle special values (case-insensitive).
	lower := strings.ToLower(s)
	switch lower {
	case "nan", "+nan", "-nan":
		return FixedPoint{
			coe: 1,
			exp: 0,
			flg: flags{nan: true},
		}
	case "inf", "infinity", "+inf", "+infinity":
		return FixedPoint{
			coe: 1,
			exp: 0,
			flg: flags{inf: true},
		}
	case "-inf", "-infinity":
		return FixedPoint{
			sign: true,
			coe:  1,
			exp:  0,
			flg:  flags{inf: true},
		}
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
		return FixedPoint{
			coe: 0,
			exp: 0,
			flg: flags{nan: true, sig: SignalConversionSyntax},
		}
	}

	intPart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}
	if intPart == "" && fracPart == "" {
		return FixedPoint{
			coe: 0,
			exp: 0,
			flg: flags{nan: true, sig: SignalConversionSyntax},
		}
	}

	// Concatenate the integer and fractional digits.
	digits := intPart + fracPart

	// Attempt to parse the digits into a uint64.
	value, err := strconv.ParseUint(digits, 10, 64)
	if err != nil {
		return FixedPoint{
			coe: 0,
			exp: 0,
			flg: flags{nan: true, sig: SignalConversionSyntax},
		}
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
		return FixedPoint{
			coe: 0,
			exp: 0,
			flg: flags{inf: true, sig: SignalOverflow},
		}
	}

	return FixedPoint{
		sign: sign,
		coe:  coe,
		exp:  exp,
		ctx:  defaultContext,
	}
}

// Signal returns the current error signal of the FixedPoint.
func (a *FixedPoint) Signal() Signal {
	if a == nil {
		return SignalConversionSyntax
	}
	return a.flg.sig
}

// IsOk checks if the FixedPoint does not have an error signal.
func (a *FixedPoint) IsOk() bool {
	return a.Signal() == SignalClear
}

// Must returns a pointer to the FixedPoint or panics if it is not OK.
func (a *FixedPoint) Must() *FixedPoint {
	if !a.IsOk() {
		panic(a.Signal())
	}

	return a
}

// Must returns a pointer to the FixedPoint or panics if it is not OK.
func Must(a FixedPoint) *FixedPoint {
	return a.Must()
}

// Handle executes the provided handler if FixedPoint holds an error.
func (a *FixedPoint) Handle(h func(Signal) *FixedPoint) *FixedPoint {
	if a.IsOk() {
		return a
	}

	return h(a.Signal())
}

// Copy returns a shallow copy of the FixedPoint.
func (a *FixedPoint) Copy() FixedPoint {
	if a == nil {
		return FixedPoint{
			coe: 0,
			exp: 0,
			flg: flags{nan: true, sig: SignalInvalidOperation},
		}
	}

	return FixedPoint{
		sign: a.sign,
		coe:  a.coe,
		exp:  a.exp,
		flg:  a.flg,
	}
}

// Clone returns a new pointer that clones the FixedPoint.
func (a *FixedPoint) Clone() *FixedPoint {
	if a == nil {
		return nil
	}

	return &FixedPoint{
		sign: a.sign,
		coe:  a.coe,
		exp:  a.exp,
		flg:  a.flg,
	}
}

// IsSpecial returns true if the FixedPoint represents NaN or Infinity.
func (a *FixedPoint) IsSpecial() bool {
	return a != nil && (a.flg.nan || a.flg.inf)
}

// IsNaN returns true if the FixedPoint is Not-a-Number.
func (a *FixedPoint) IsNaN() bool {
	return a != nil && a.flg.nan
}

// IsInf returns true if the FixedPoint represents Infinity.
func (a *FixedPoint) IsInf() bool {
	return a != nil && a.flg.inf
}

// IsZero returns true if the FixedPoint is zero (ignoring sign).
func (a *FixedPoint) IsZero() bool {
	if a == nil {
		return false
	}
	if a.flg.nan || a.flg.inf {
		return false
	}
	return a.coe == 0
}

// IsNegative returns true if the FixedPoint is negative.
func (a *FixedPoint) IsNegative() bool {
	if a == nil {
		return false
	}
	if a.flg.nan || a.flg.inf {
		return false
	}
	return a.sign
}

// IsPositive returns true if the FixedPoint is positive.
func (a *FixedPoint) IsPositive() bool {
	if a == nil {
		return false
	}
	if a.flg.nan || a.flg.inf {
		return false
	}
	return !a.sign
}

// Add returns the sum of this FixedPoint and another.
func (a *FixedPoint) Add(b *FixedPoint) FixedPoint {
	// Replace nil-check with helper.
	if a == nil || b == nil {
		return invalid_operation()
	}

	// Align exponents: choose the smaller exponent for maximum precision.
	minExp := min(b.exp, a.exp)

	// Scale each coefficient to the common exponent.
	// For a, we need to multiply by 10^(a.exp - minExp)
	// For b, we need to multiply by 10^(b.exp - minExp)
	x_coe, xok := scale_coe(a.coe, a.exp-minExp)
	y_coe, yok := scale_coe(b.coe, b.exp-minExp)
	if !xok || !yok {
		return overflow_operation()
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
	if !ok {
		return overflow_operation()
	}

	// Check coefficient against maximum allowed value.
	if res_coe > fp_coe_max_val {
		return overflow_operation()
	}

	// The result uses the common (minimum) exponent.
	return FixedPoint{
		sign: res_sign,
		coe:  res_coe,
		exp:  minExp,
		ctx:  a.ctx,
	}
}

// Sub returns the difference between this FixedPoint and another.
func (a *FixedPoint) Sub(b *FixedPoint) FixedPoint {
	b_neg := b.Neg()
	return a.Add(&b_neg)
}

// Mul returns the product of this FixedPoint and another.
func (a *FixedPoint) Mul(b *FixedPoint) FixedPoint {
	// Return NaN if either operand is nil.
	if a == nil || b == nil {
		return FixedPoint{
			coe: 0,
			exp: 0,
			flg: flags{nan: true, sig: SignalInvalidOperation},
		}
	}

	// Safely multiply the absolute values.
	var res_abs coefficient
	if a.coe == 0 || b.coe == 0 {
		res_abs = 0
	} else {
		// Check for multiplication overflow.
		if a.coe > fp_coe_max_val/b.coe {
			return FixedPoint{
				coe: 0,
				exp: 0,
				flg: flags{inf: true, sig: SignalOverflow},
			}
		}
		res_abs = a.coe * b.coe
	}

	// Check if the result exceeds the maximum allowed coefficient.
	if res_abs > fp_coe_max_val {
		return FixedPoint{
			coe: 0,
			exp: 0,
			flg: flags{inf: true, sig: SignalOverflow},
		}
	}

	// Determine the sign of the result.
	resSign := a.sign != b.sign

	// Add the exponents.
	resExp := a.exp + b.exp

	return FixedPoint{
		sign: resSign,
		coe:  res_abs,
		exp:  resExp,
		ctx:  a.ctx,
	}
}

// Div returns the quotient of this FixedPoint divided by another.
func (a *FixedPoint) Div(b *FixedPoint) FixedPoint {
	// Replace nil-check with helper.
	if a == nil || b == nil {
		return invalid_operation()
	}

	// Handle divisor equals zero.
	if b.coe == 0 {
		switch a.coe {
		case 0:
			return FixedPoint{
				coe: 0,
				exp: 0,
				flg: flags{nan: true, sig: SignalDivisionImpossible},
			}
		default:
			return FixedPoint{
				sign: a.sign != b.sign,
				coe:  0,
				exp:  0,
				flg:  flags{inf: true, sig: SignalDivisionByZero},
			}
		}
	}

	// Long division algorithm.
	adjust := 0
	work_dividend := a.coe
	work_divisor := b.coe

	// If dividend is non-zero, adjust coefficients.
	if work_dividend != 0 {
		// Scale up dividend until it's >= divisor.
		for work_dividend < work_divisor {
			// Check for overflow is omitted for brevity.
			work_dividend *= 10
			adjust++
		}
		// Scale down divisor if dividend is too large.
		for work_dividend >= work_divisor*10 {
			work_divisor *= 10
			adjust--
		}
	}

	var res_coe coefficient = 0
	dig_count := 0
	// If dividend is zero, resultCoeff remains zero.
	if work_dividend != 0 {
		for {
			// Count how many times workingDivisor fits.
			var count coefficient = 0
			for work_dividend >= work_divisor {
				work_dividend -= work_divisor
				count++
			}
			res_coe = res_coe*10 + count
			dig_count++
			// Termination condition.
			if (work_dividend == 0 && adjust >= 0) || (dig_count >= fp_coe_max_len) {
				break
			}
			work_dividend *= 10
			adjust++
		}
	}

	// Compute the resulting exponent.
	res_exp := a.exp - b.exp - exponent(adjust)
	// Determine the result's sign (exclusive or).
	res_sign := a.sign != b.sign

	// Check for coefficient overflow.
	if res_coe > fp_coe_max_val {
		return FixedPoint{
			coe: 0,
			exp: 0,
			flg: flags{inf: true, sig: SignalOverflow},
		}
	}

	return FixedPoint{
		sign: res_sign,
		coe:  res_coe,
		exp:  res_exp,
		ctx:  a.ctx,
	}
}

// Neg returns the negation of the FixedPoint.
func (a *FixedPoint) Neg() FixedPoint {
	if a == nil {
		return invalid_operation()
	}

	return FixedPoint{
		sign: !a.sign,
		coe:  a.coe,
		exp:  a.exp,
		ctx:  a.ctx,
	}
}

// Abs returns the absolute value of the FixedPoint.
func (a *FixedPoint) Abs() FixedPoint {
	if a == nil {
		return FixedPoint{flg: flags{nan: true, sig: SignalInvalidOperation}}
	}

	if a.flg.nan || a.flg.inf {
		// Preserve special values but make sign positive
		result := *a
		result.sign = false
		return result
	}

	return FixedPoint{
		sign: false,
		coe:  a.coe,
		exp:  a.exp,
		ctx:  a.ctx,
	}
}

// Equal returns true if this FixedPoint is equal to the other.
func (a *FixedPoint) Equal(b *FixedPoint) bool {
	if a == nil || b == nil {
		return false
	}

	// If either is NaN or Infinity, check flags
	if a.flg.nan || b.flg.nan {
		return false // NaN is not equal to anything
	}
	if a.flg.inf && b.flg.inf {
		return a.sign == b.sign // Same infinity sign
	}
	if a.flg.inf || b.flg.inf {
		return false
	}

	// Handle zero specially
	if a.coe == 0 && b.coe == 0 {
		return true // Both are zero (sign doesn't matter)
	}

	// If signs differ, they're not equal
	if a.sign != b.sign {
		return false
	}

	// If exponents are equal, just compare coefficients
	if a.exp == b.exp {
		return a.coe == b.coe
	}

	// Need to align exponents to compare
	if a.exp > b.exp {
		scaledCoe, ok := scale_coe(a.coe, a.exp-b.exp)
		if !ok {
			return false
		}
		return scaledCoe == b.coe
	} else {
		scaledCoe, ok := scale_coe(b.coe, b.exp-a.exp)
		if !ok {
			return false
		}
		return a.coe == scaledCoe
	}
}

// LessThan returns true if this FixedPoint is less than the other.
func (a *FixedPoint) LessThan(b *FixedPoint) bool {
	if a == nil || b == nil || a.flg.nan || b.flg.nan {
		return false // NaN comparisons always return false
	}

	// Handle infinities
	if a.flg.inf {
		return a.sign // -Infinity is less than anything except itself
	}
	if b.flg.inf {
		return !b.sign // Anything is less than +Infinity
	}

	// If signs differ, negative < positive
	if a.sign != b.sign {
		return a.sign
	}

	// Same sign, align exponents and compare
	var cmpResult bool
	if a.exp == b.exp {
		cmpResult = a.coe < b.coe
	} else if a.exp > b.exp {
		scaledCoe, ok := scale_coe(a.coe, a.exp-b.exp)
		if !ok {
			return a.sign // Overflow means large number, sign determines result
		}
		cmpResult = scaledCoe < b.coe
	} else {
		scaledCoe, ok := scale_coe(b.coe, b.exp-a.exp)
		if !ok {
			return !a.sign // Overflow means large number, sign determines result
		}
		cmpResult = a.coe < scaledCoe
	}

	// If negative, reverse comparison result
	if a.sign {
		return !cmpResult
	}
	return cmpResult
}

// LessThanOrEqual returns true if this FixedPoint is less than or equal to the other.
func (a *FixedPoint) LessThanOrEqual(b *FixedPoint) bool {
	return a.LessThan(b) || a.Equal(b)
}

// GreaterThan returns true if this FixedPoint is greater than the other.
func (a *FixedPoint) GreaterThan(b *FixedPoint) bool {
	// Returns true if a > b.
	return b.LessThan(a)
}

// GreaterThanOrEqual returns true if this FixedPoint is greater than or equal to the other.
func (a *FixedPoint) GreaterThanOrEqual(b *FixedPoint) bool {
	// Returns true if a >= b.
	return !a.LessThan(b)
}

// Compare compares this FixedPoint with another returning -1, 0, or 1.
func (a *FixedPoint) Compare(b *FixedPoint) int {
	// Returns -1 if a < b, 0 if a == b, 1 if a > b.
	if a.Equal(b) {
		return 0
	} else if a.LessThan(b) {
		return -1
	}
	return 1
}

// Max returns a pointer to the maximum of two FixedPoint values.
func Max(a, b *FixedPoint) *FixedPoint {
	if a.Compare(b) < 0 {
		return b
	}
	return a
}

// Min returns a pointer to the minimum of two FixedPoint values.
func Min(a, b *FixedPoint) *FixedPoint {
	if a.Compare(b) > 0 {
		return b
	}
	return a
}

func invalid_operation() FixedPoint {
	return FixedPoint{
		coe: 0,
		exp: 0,
		flg: flags{nan: true, sig: SignalInvalidOperation},
	}
}

func overflow_operation() FixedPoint {
	return FixedPoint{
		coe: 0,
		exp: 0,
		flg: flags{inf: true, sig: SignalOverflow},
	}
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
