package fixedpoint

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func (fp FixedPoint128) String() string {
	if fp.isNaN() {
		return "NaN"
	}
	if fp.isSNaN() {
		return "sNaN"
	}
	if fp.isInf() {
		if fp.sign() {
			return "-Infinity"
		}
		return "Infinity"
	}

	c := new(big.Int).Set(fp.coefficient())
	exp := fp.exponent()

	r := new(big.Rat).SetInt(c)
	if exp < 0 {
		denom := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-exp)), nil)
		r.Quo(r, new(big.Rat).SetInt(denom))
	} else if exp > 0 {
		num := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exp)), nil)
		r.Mul(r, new(big.Rat).SetInt(num))
	}

	// Use high precision to avoid rounding
	s := r.FloatString(34) // 34 digits is Decimal128 max
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")

	if fp.sign() {
		return "-" + s
	}
	return s
}

func (fp *FixedPoint64) String() string {
	if fp.isNaN() {
		return "NaN"
	}
	if fp.isSNaN() {
		return "sNaN"
	}
	if fp.isInf() {
		if fp.sign() {
			return "-Infinity"
		}
		return "Infinity"
	}
	c := new(big.Int).SetUint64(fp.coefficient())
	e := fp.exponent()

	r := new(big.Rat).SetInt(c)
	if e < 0 {
		r.Quo(r, new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-e)), nil)))
	} else if e > 0 {
		r.Mul(r, new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(e)), nil)))
	}

	s := r.FloatString(16)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if fp.sign() {
		return "-" + s
	}
	return s
}

func (fp *FixedPoint128) Scientific() string {
	if fp.isNaN() {
		return "NaN"
	}
	if fp.isSNaN() {
		return "sNaN"
	}
	if fp.isInf() {
		if fp.sign() {
			return "-Infinity"
		}
		return "Infinity"
	}

	coef := fp.coefficient()
	exp := fp.exponent()

	// Convert to scientific format: move decimal point after first digit
	coefStr := coef.String()
	digitCount := len(coefStr)
	if coef.Sign() == 0 {
		return "0e+0"
	}

	sciExp := exp + digitCount - 1
	var sb strings.Builder
	if fp.sign() {
		sb.WriteByte('-')
	}

	sb.WriteByte(coefStr[0])
	if digitCount > 1 {
		sb.WriteByte('.')
		sb.WriteString(coefStr[1:])
	}

	fmt.Fprintf(&sb, "e%+d", sciExp)
	return sb.String()
}

func (fp *FixedPoint64) Scientific() string {
	if fp.isNaN() {
		return "NaN"
	}
	if fp.isSNaN() {
		return "sNaN"
	}
	if fp.isInf() {
		if fp.sign() {
			return "-Infinity"
		}
		return "Infinity"
	}

	coef := fp.coefficient()
	if coef == 0 {
		return "0e+0"
	}

	str := strconv.FormatUint(coef, 10)
	sciExp := fp.exponent() + len(str) - 1

	var b strings.Builder
	if fp.sign() {
		b.WriteByte('-')
	}
	b.WriteByte(str[0])
	if len(str) > 1 {
		b.WriteByte('.')
		b.WriteString(str[1:])
	}
	fmt.Fprintf(&b, "e%+d", sciExp)
	return b.String()
}

func (fp *FixedPoint128) Debug() string {
	var kind string
	switch {
	case fp.isSNaN():
		kind = "Signaling NaN"
	case fp.isNaN():
		kind = "Quiet NaN"
	case fp.isInf():
		if fp.sign() {
			kind = "-Infinity"
		} else {
			kind = "+Infinity"
		}
	default:
		kind = "Finite"
	}

	return fmt.Sprintf("Kind: %s\nSign: %v\nExponent: %d\nCoefficient: %s\nRaw Hi: 0x%016X\nRaw Lo: 0x%016X",
		kind,
		fp.sign(),
		fp.exponent(),
		fp.coefficient().String(),
		fp.hi,
		fp.lo,
	)
}

func (fp *FixedPoint64) Debug() string {
	sign := fp.sign()
	coef := fp.coefficient()
	exp := fp.exponent()

	var kind string
	switch {
	case fp.isNaN():
		kind = "Quiet NaN"
	case fp.isSNaN():
		kind = "Signaling NaN"
	case fp.isInf():
		if sign {
			kind = "-Infinity"
		} else {
			kind = "+Infinity"
		}
	default:
		kind = "Finite"
	}

	return fmt.Sprintf(
		"Kind: %s\nSign: %v\nExponent: %d\nCoefficient: %d\nRaw Bits: 0x%016X",
		kind,
		sign,
		exp,
		coef,
		fp.bits,
	)
}
