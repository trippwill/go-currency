package fixedpoint

import (
	"math/big"
	"strconv"
	"strings"
)

// FixedPoint128 represents a fixed-point number in 128-bit format.
type FixedPoint128 struct {
	hi uint64
	lo uint64
}

func Parse128(s string) (FixedPoint128, error) {
	var d FixedPoint128

	s = strings.TrimSpace(s)
	lower := strings.ToLower(s)
	switch lower {
	case "nan", "+nan":
		d.setNaN(false)
		return d, nil
	case "-nan":
		d.setNaN(true)
		return d, nil
	case "inf", "infinity", "+inf", "+infinity":
		d.setInf(false)
		return d, nil
	case "-inf", "-infinity":
		d.setInf(true)
		return d, nil
	}

	// Determine the sign.
	neg := false
	if s != "" {
		if s[0] == '-' {
			neg = true
			s = s[1:]
		} else if s[0] == '+' {
			s = s[1:]
		}
	}

	// Split into base and exponent parts.
	var basePart, expPart string
	if i := strings.IndexAny(s, "eE"); i != -1 {
		basePart = s[:i]
		expPart = s[i+1:]
	} else {
		basePart = s
	}

	// Parse the explicit exponent (if any); default is 0.
	expVal := 0
	if expPart != "" {
		var err error
		expVal, err = strconv.Atoi(expPart)
		if err != nil {
			return FixedPoint128{}, ErrConversionSyntax
		}
	}

	// If there is a decimal point, remove it and note how many digits were after it.
	decDigits := 0
	if i := strings.Index(basePart, "."); i != -1 {
		decDigits = len(basePart) - i - 1
		basePart = strings.Replace(basePart, ".", "", 1)
	}

	// Remove any leading zeros.
	basePart = strings.TrimLeft(basePart, "0")
	if basePart == "" {
		// The value is zero.
		d.setSign(neg)
		d.setExponent(0)
		d.setCoefficient(big.NewInt(0))
		return d, nil
	}

	// Adjust the total exponent: exponent from the scientific notation minus the number of digits after the decimal point.
	totalExp := expVal - decDigits

	// Enforce a maximum of 34 digits.
	if len(basePart) > 34 {
		return FixedPoint128{}, ErrOverflow
	}

	coef := new(big.Int)
	if _, ok := coef.SetString(basePart, 10); !ok {
		return FixedPoint128{}, ErrConversionSyntax
	}

	d.setSign(neg)
	d.setExponent(totalExp)
	d.setCoefficient(coef)
	return d, nil
}

const (
	fp128_exp_bias    = 6176
	fp128_exp_mask    = uint64(0x3FFF)             // 14 bits
	fp128_coe_mask_hi = uint64(0x0001FFFFFFFFFFFF) // 49 bits
)

func (fp *FixedPoint128) sign() bool {
	return fp.hi>>63 != 0
}

func (fp *FixedPoint128) setSign(s bool) {
	if s {
		fp.hi |= 1 << 63
	} else {
		fp.hi &^= 1 << 63
	}
}

func (fp *FixedPoint128) coefficient() *big.Int {
	_h := fp.hi & fp128_coe_mask_hi
	_l := fp.lo

	hi := new(big.Int).Lsh(new(big.Int).SetUint64(_h), 64)
	lo := new(big.Int).SetUint64(_l)
	return hi.Or(hi, lo)
}

func (fp *FixedPoint128) setCoefficient(c *big.Int) SIG {
	if c.Sign() < 0 {
		return SIG_INVALID_OPERATION
	}
	if c.BitLen() > 113 {
		return SIG_OVERFLOW
	}

	lo := new(big.Int).And(c, big.NewInt(0).SetUint64(0xFFFFFFFFFFFFFFFF))
	hi := new(big.Int).Rsh(c, 64)

	fp.lo = lo.Uint64()
	fp.hi &^= fp128_coe_mask_hi              // clear top 49 bits
	fp.hi |= hi.Uint64() & fp128_coe_mask_hi // set top 49 bits

	return SIG_NONE
}

func (fp *FixedPoint128) exponent() int {
	biased := (fp.hi >> 49) & fp128_exp_mask
	return int(biased) - fp128_exp_bias
}

func (fp *FixedPoint128) setExponent(exp int) SIG {
	biased := uint64(exp + fp128_exp_bias)
	if biased >= (1 << 14) {
		return SIG_INVALID_OPERATION
	}
	// Clear old exponent bits
	fp.hi &^= fp128_exp_mask << 49
	// Set new exponent
	fp.hi |= biased << 49

	return SIG_NONE
}

func (fp *FixedPoint128) combinationField() uint8 {
	return uint8((fp.hi >> 58) & 0x1F)
}

func (fp *FixedPoint128) isNaN() bool {
	cf := fp.combinationField()
	return cf == 0b11110 || cf == 0b11111
}

func (fp *FixedPoint128) setNaN(sign bool) {
	fp.hi = 0
	fp.lo = 0
	if sign {
		fp.hi |= 1 << 63
	}
	fp.hi |= uint64(0b11110) << 58
	fp.hi |= 1 // set some coefficient bit to distinguish from signaling NaN
}

func (fp *FixedPoint128) isSNaN() bool {
	return fp.combinationField() == 0b11111
}

func (fp *FixedPoint128) setSNaN(sign bool) {
	fp.hi = 0
	fp.lo = 0
	if sign {
		fp.hi |= 1 << 63
	}
	fp.hi |= uint64(0b11111) << 58
	fp.hi |= 1 // set some coefficient bit
}

func (fp *FixedPoint128) isInf() bool {
	return fp.combinationField() == 0b11100
}

func (fp *FixedPoint128) setInf(sign bool) {
	fp.hi = 0
	fp.lo = 0
	if sign {
		fp.hi |= 1 << 63 // set sign bit
	}
	fp.hi |= uint64(0b11100) << 58
}

func (fp *FixedPoint128) isFinite() bool {
	cf := fp.combinationField()
	return cf < 0b11100
}
