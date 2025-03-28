// Based on the General Decimal Arithmetic Specification 1.70 – 7 Apr 2009
// https://speleotrove.com/decimal/decarith.html
package fixedpoint

import (
	"fmt"
	"strconv"
	"strings"
)

// FixedPoint represents a fixed-point arithmetic type with a wide dynamic range.
// Representable ranges:
//   - Positive values: from the smallest positive 1.0e-9999 to the largest positive 9.999999999999999e+9999
//   - Negative values: from the largest negative -9.999999999999999e+9999 to the smallest negative -1.0e-9999
type FixedPoint struct {
	coe coefficient // sign[63] + coefficient[62:0]
	exp exponent    // exponent
	flg flags       // flags
	_   uint8       // padding
}

type (
	coefficient uint64
	exponent    int16
	flags       struct {
		sig Signal
		inf bool
		nan bool
		_   bool
		_   bool
	}
)

var Zero = FixedPoint{
	coe: 0,
	exp: 0,
}

var NegZero = FixedPoint{
	coe: pack(-1, 0),
	exp: 0,
}

var One = FixedPoint{
	coe: pack(1, 1),
	exp: 0,
}

var NegOne = FixedPoint{
	coe: pack(-1, 1),
	exp: 0,
}

const (
	fp_coe_max_val   uint64   = 999_999_999_999_999_999
	fp_coe_max_len            = 18
	fp_exp_limit_val exponent = 9_999
	fp_exp_limit_len          = 4
)

// New returns a new FixedPoint value from a significand and an exponent.
func New(significand int64, exp int16) FixedPoint {
	coe, ok := convert_signed(significand)
	if !ok {
		return FixedPoint{
			coe: coe,
			exp: exponent(exp),
			flg: flags{nan: true, sig: SignalOverflow},
		}
	}

	return FixedPoint{
		coe: coe,
		exp: exponent(exp),
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
			coe: pack(1, 0),
			exp: 0,
			flg: flags{nan: true},
		}
	case "inf", "infinity", "+inf", "+infinity":
		return FixedPoint{
			coe: pack(1, 0),
			exp: 0,
			flg: flags{inf: true},
		}
	case "-inf", "-infinity":
		return FixedPoint{
			coe: pack(-1, 0),
			exp: 0,
			flg: flags{inf: true},
		}
	}

	// Determine sign.
	sign := int8(1)
	if s[0] == '-' {
		sign = -1
		s = s[1:]
	} else if s[0] == '+' {
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

	// Check for coefficient overflow.
	if value > fp_coe_max_val {
		return FixedPoint{
			coe: 0,
			exp: 0,
			flg: flags{inf: true, sig: SignalOverflow},
		}
	}

	coe := pack(sign, value)
	return FixedPoint{
		coe: coe,
		exp: exp,
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
		coe: a.coe,
		exp: a.exp,
		flg: a.flg,
	}
}

// Clone returns a new pointer that clones the FixedPoint.
func (a *FixedPoint) Clone() *FixedPoint {
	if a == nil {
		return nil
	}

	return &FixedPoint{
		coe: a.coe,
		exp: a.exp,
		flg: a.flg,
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
	_, abs := unpack(a.coe)
	return abs == 0
}

// IsNegative returns true if the FixedPoint is negative.
func (a *FixedPoint) IsNegative() bool {
	if a == nil {
		return false
	}
	if a.flg.nan || a.flg.inf {
		return false
	}
	sign, _ := unpack(a.coe)
	return sign < 0
}

// IsPositive returns true if the FixedPoint is positive.
func (a *FixedPoint) IsPositive() bool {
	if a == nil {
		return false
	}
	if a.flg.nan || a.flg.inf {
		return false
	}
	sign, _ := unpack(a.coe)
	return sign > 0
}

// Debug returns a formatted string with debug information about the FixedPoint.
func (a *FixedPoint) Debug() string {
	s, abs := unpack(a.coe)
	return fmt.Sprintf("sign: %d, abs: %d, exp: %d, flg: %v", s, abs, a.exp, a.flg)
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

	x_s, x_abs := unpack(x_coe)
	y_s, y_abs := unpack(y_coe)

	var res_abs uint64
	var res_sign int8
	var ok bool

	// When signs are identical, do simple addition.
	if x_s == y_s {
		res_abs, ok = safe_add(x_abs, y_abs)
		res_sign = x_s
	} else {
		// When signs differ, subtract the smaller magnitude from the larger.
		if x_abs >= y_abs {
			res_abs, ok = safe_sub(x_abs, y_abs)
			res_sign = x_s
		} else {
			res_abs, ok = safe_sub(y_abs, x_abs)
			res_sign = y_s
		}
		// Zero result should always be positive.
		if res_abs == 0 {
			res_sign = 1
		}
	}

	// Overflow during addition or subtraction.
	if !ok {
		return overflow_operation()
	}

	// Check coefficient against maximum allowed value.
	if res_abs > fp_coe_max_val {
		return overflow_operation()
	}

	// Pack the coefficient with its sign.
	res_coe := pack(res_sign, res_abs)

	// The result uses the common (minimum) exponent.
	return FixedPoint{
		coe: res_coe,
		exp: minExp,
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

	// Unpack the coefficients.
	aSign, aAbs := unpack(a.coe)
	bSign, bAbs := unpack(b.coe)

	// Safely multiply the absolute values.
	var resAbs uint64
	if aAbs == 0 || bAbs == 0 {
		resAbs = 0
	} else {
		// Check for multiplication overflow.
		if aAbs > fp_coe_max_val/bAbs {
			return FixedPoint{
				coe: 0,
				exp: 0,
				flg: flags{inf: true, sig: SignalOverflow},
			}
		}
		resAbs = aAbs * bAbs
	}

	// Check if the result exceeds the maximum allowed coefficient.
	if resAbs > fp_coe_max_val {
		return FixedPoint{
			coe: 0,
			exp: 0,
			flg: flags{inf: true, sig: SignalOverflow},
		}
	}

	// Determine the sign of the result.
	resSign := aSign * bSign

	// Add the exponents.
	resExp := a.exp + b.exp

	// Pack the result coefficient.
	resCoe := pack(resSign, resAbs)

	return FixedPoint{
		coe: resCoe,
		exp: resExp,
	}
}

// Div returns the quotient of this FixedPoint divided by another.
func (a *FixedPoint) Div(b *FixedPoint) FixedPoint {
	// Replace nil-check with helper.
	if a == nil || b == nil {
		return invalid_operation()
	}

	// Unpack the coefficients.
	aSign, a_abs := unpack(a.coe)
	bSign, b_abs := unpack(b.coe)

	// Handle divisor equals zero.
	if b_abs == 0 {
		switch a_abs {
		case 0:
			return FixedPoint{
				coe: 0,
				exp: 0,
				flg: flags{nan: true, sig: SignalDivisionImpossible},
			}
		default:
			return FixedPoint{
				coe: pack(aSign*bSign, 0),
				exp: 0,
				flg: flags{inf: true, sig: SignalDivisionByZero},
			}
		}
	}

	// Long division algorithm.
	var adjust int = 0
	var workingDividend uint64 = a_abs
	var workingDivisor uint64 = b_abs

	// If dividend is non-zero, adjust coefficients.
	if workingDividend != 0 {
		// Scale up dividend until it's >= divisor.
		for workingDividend < workingDivisor {
			// Check for overflow is omitted for brevity.
			workingDividend *= 10
			adjust++
		}
		// Scale down divisor if dividend is too large.
		for workingDividend >= workingDivisor*10 {
			workingDivisor *= 10
			adjust--
		}
	}

	var resultCoeff uint64 = 0
	var digitCount int = 0
	// If dividend is zero, resultCoeff remains zero.
	if workingDividend != 0 {
		for {
			// Count how many times workingDivisor fits.
			var count uint64 = 0
			for workingDividend >= workingDivisor {
				workingDividend -= workingDivisor
				count++
			}
			resultCoeff = resultCoeff*10 + count
			digitCount++
			// Termination condition.
			if (workingDividend == 0 && adjust >= 0) || (digitCount >= fp_coe_max_len) {
				break
			}
			workingDividend *= 10
			adjust++
		}
	}

	// Compute the resulting exponent.
	resExp := a.exp - b.exp - exponent(adjust)
	// Determine the result's sign (exclusive or).
	resSign := aSign * bSign

	// Check for coefficient overflow.
	if resultCoeff > fp_coe_max_val {
		return FixedPoint{
			coe: 0,
			exp: 0,
			flg: flags{inf: true, sig: SignalOverflow},
		}
	}

	// Pack the result.
	resCoe := pack(resSign, resultCoeff)
	return FixedPoint{
		coe: resCoe,
		exp: resExp,
	}
}

// Neg returns the negation of the FixedPoint.
func (a *FixedPoint) Neg() FixedPoint {
	if a == nil {
		return invalid_operation()
	}

	s, abs := unpack(a.coe)
	return FixedPoint{
		coe: pack(-s, abs),
		exp: a.exp,
	}
}

// String returns the human-readable string representation of the FixedPoint.
func (a *FixedPoint) String() string {
	if a == nil {
		return "NaN"
	}
	if a.flg.nan {
		return "NaN"
	}
	if a.flg.inf {
		sign, _ := unpack(a.coe)
		if sign < 0 {
			return "-Infinity"
		}
		return "Infinity"
	}

	sign, abs := unpack(a.coe)

	// Convert the absolute value to a string
	str := fmt.Sprintf("%d", abs)

	// Apply the exponent
	exp := int(a.exp)
	if exp >= 0 {
		// Add trailing zeros
		for range exp {
			str += "0"
		}
	} else {
		// Insert decimal point
		expAbs := -exp
		if len(str) <= expAbs {
			// Need to pad with leading zeros
			padding := expAbs - len(str)
			str = "0." + strings.Repeat("0", padding) + str
		} else {
			// Insert decimal point at the correct position
			pos := len(str) - expAbs
			str = str[:pos] + "." + str[pos:]
		}
		// Trim trailing zeros after decimal point
		str = strings.TrimRight(str, "0")
		str = strings.TrimRight(str, ".")
	}

	// Add sign
	if sign < 0 {
		str = "-" + str
	}

	return str
}

func (f flags) String() string {
	return fmt.Sprintf("flags{sig: %s, inf: %t, nan: %t}", f.sig, f.inf, f.nan)
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
		signA, _ := unpack(a.coe)
		signB, _ := unpack(b.coe)
		return signA == signB // Same infinity sign
	}
	if a.flg.inf || b.flg.inf {
		return false
	}

	// Normalize both values to compare
	signA, absA := unpack(a.coe)
	signB, absB := unpack(b.coe)

	// Handle zero specially
	if absA == 0 && absB == 0 {
		return true // Both are zero (sign doesn't matter)
	}

	// If signs differ, they're not equal
	if signA != signB {
		return false
	}

	// If exponents are equal, just compare coefficients
	if a.exp == b.exp {
		return absA == absB
	}

	// Need to align exponents to compare
	if a.exp > b.exp {
		scaledCoe, ok := scale_coe(a.coe, a.exp-b.exp)
		if !ok {
			return false
		}
		_, scaledAbs := unpack(scaledCoe)
		return scaledAbs == absB
	} else {
		scaledCoe, ok := scale_coe(b.coe, b.exp-a.exp)
		if !ok {
			return false
		}
		_, scaledAbs := unpack(scaledCoe)
		return absA == scaledAbs
	}
}

// LessThan returns true if this FixedPoint is less than the other.
func (a *FixedPoint) LessThan(b *FixedPoint) bool {
	if a == nil || b == nil || a.flg.nan || b.flg.nan {
		return false // NaN comparisons always return false
	}

	// Handle infinities
	if a.flg.inf {
		signA, _ := unpack(a.coe)
		return signA < 0 // -Infinity is less than anything except itself
	}
	if b.flg.inf {
		signB, _ := unpack(b.coe)
		return signB > 0 // Anything is less than +Infinity
	}

	// Extract signs and absolute values
	signA, absA := unpack(a.coe)
	signB, absB := unpack(b.coe)

	// If signs differ, negative < positive
	if signA != signB {
		return signA < signB
	}

	// Same sign, align exponents and compare
	var cmpResult bool
	if a.exp == b.exp {
		cmpResult = absA < absB
	} else if a.exp > b.exp {
		scaledCoe, ok := scale_coe(a.coe, a.exp-b.exp)
		if !ok {
			return signA < 0 // Overflow means large number, sign determines result
		}
		_, scaledAbs := unpack(scaledCoe)
		cmpResult = scaledAbs < absB
	} else {
		scaledCoe, ok := scale_coe(b.coe, b.exp-a.exp)
		if !ok {
			return signA > 0 // Overflow means large number, sign determines result
		}
		_, scaledAbs := unpack(scaledCoe)
		cmpResult = absA < scaledAbs
	}

	// If negative, reverse comparison result
	if signA < 0 {
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

// Abs returns the absolute value of the FixedPoint.
func (a *FixedPoint) Abs() FixedPoint {
	if a == nil {
		return FixedPoint{flg: flags{nan: true, sig: SignalInvalidOperation}}
	}

	if a.flg.nan || a.flg.inf {
		// Preserve special values but make sign positive
		result := *a
		_, abs := unpack(a.coe)
		result.coe = pack(1, abs)
		return result
	}

	_, abs := unpack(a.coe)
	return FixedPoint{
		coe: pack(1, abs),
		exp: a.exp,
	}
}

type signed interface {
	int | int8 | int16 | int32 | int64 | exponent
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

func convert_signed[I signed](x I) (coefficient, bool) {
	var sign uint64
	var abs uint64
	v := int64(x)

	if v < 0 {
		sign = 1
		// Compute absolute value using bitwise complement plus one,
		// which correctly handles v == math.MinInt64.
		abs = uint64(^v) + 1
	} else {
		abs = uint64(v)
	}

	// Check for overflow
	ok := abs <= fp_coe_max_val

	// Ensure the sign bit (bit 63) is set appropriately.
	return coefficient((sign << 63) | (abs & 0x7FFFFFFFFFFFFFFF)), ok
}

func pack(s int8, abs uint64) coefficient {
	var signBit uint64
	if s < 0 {
		signBit = 1
	}

	return coefficient((signBit << 63) | (abs & 0x7FFFFFFFFFFFFFFF))
}

func unpack(c coefficient) (int8, uint64) {
	var sign int8
	var abs uint64

	// Extract the sign bit (bit 63).
	if c&(1<<63) != 0 {
		sign = -1
	} else {
		sign = 1
	}

	// Extract the absolute value (lower 63 bits).
	abs = uint64(c & 0x7FFFFFFFFFFFFFFF)

	// Return the signed value.
	return sign, uint64(abs)
}

// scale_coe always adjusts the coefficient to the desired exponent without losing precision.
func scale_coe(c coefficient, diff exponent) (coefficient, bool) {
	s, abs := unpack(c)
	if diff > 0 {
		// Multiply absolute value by 10 for each increment in diff.
		for i := exponent(0); i < diff; i++ {
			// Check for multiplication overflow.
			if abs > fp_coe_max_val/10 {
				return 0, false
			}
			abs *= 10
		}
	} else if diff < 0 {
		// Divide absolute value by 10 for each decrement in diff,
		// ensuring no remainder is lost.
		for i := diff; i < 0; i++ {
			if abs%10 != 0 {
				return 0, false
			}
			abs /= 10
		}
	}
	return pack(s, abs), true
}

func safe_add(x, y uint64) (uint64, bool) {
	if x > fp_coe_max_val-y {
		return 0, false
	}

	return x + y, true
}

func safe_sub(x, y uint64) (uint64, bool) {
	if x < y {
		return 0, false
	}

	return x - y, true
}
