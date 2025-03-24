package fixedpoint

type FixedPoint64 struct {
	bits uint64
}

const (
	fp64_exp_bits = 10
	fp64_exp_bias = 398
	fp64_exp_mask = uint64(0x3FF)            // 10 bits
	fp64_coe_mask = uint64(0x1FFFFFFFFFFFFF) // 53 bits
)

func (fp *FixedPoint64) sign() bool {
	return fp.bits>>63 != 0
}

func (fp *FixedPoint64) setSign(neg bool) {
	if neg {
		fp.bits |= 1 << 63
	} else {
		fp.bits &^= 1 << 63
	}
}

func (fp *FixedPoint64) exponent() int {
	biased := (fp.bits >> 53) & fp64_exp_mask
	return int(biased) - fp64_exp_bias
}

func (fp *FixedPoint64) setExponent(exp int) {
	biased := uint64(exp + fp64_exp_bias)
	if biased >= (1 << fp64_exp_bits) {
		panic("exponent out of range")
	}
	fp.bits &^= fp64_exp_mask << 53
	fp.bits |= biased << 53
}

func (fp *FixedPoint64) coefficient() uint64 {
	return fp.bits & fp64_coe_mask
}

func (fp *FixedPoint64) setCoefficient(c uint64) {
	if c >= (1 << 53) {
		panic("coefficient too large")
	}
	fp.bits &^= fp64_coe_mask
	fp.bits |= c
}

func (fp *FixedPoint64) isNaN() bool {
	class := (fp.bits >> 58) & 0b11111
	return class == 0b11110 || class == 0b11111
}

func (fp *FixedPoint64) setNaN(sign bool) {
	fp.bits = 0
	if sign {
		fp.bits |= 1 << 63
	}
	fp.bits |= uint64(0b11110) << 58
	fp.bits |= 1 // ensure coefficient is non-zero
}

func (fp *FixedPoint64) isSNaN() bool {
	return (fp.bits>>58)&0b11111 == 0b11111
}

func (fp *FixedPoint64) setSNaN(sign bool) {
	fp.bits = 0
	if sign {
		fp.bits |= 1 << 63
	}
	fp.bits |= uint64(0b11111) << 58
	fp.bits |= 1 // ensure coefficient is non-zero
}

func (fp *FixedPoint64) isInf() bool {
	return (fp.bits>>58)&0b11111 == 0b11100
}

func (fp *FixedPoint64) setInf(sign bool) {
	fp.bits = 0
	if sign {
		fp.bits |= 1 << 63
	}
	fp.bits |= uint64(0b11100) << 58
}

func (fp *FixedPoint64) isFinite() bool {
	return (fp.bits>>58)&0b11111 < 0b11100
}
