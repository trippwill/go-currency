package fixedpoint

import "fmt"

// X64 implements the IEEE 754-2008 decimal64 format
// using Binary Integer Decimal (BID) encoding
type X64 struct {
	uint64
}

var _ packed[int16, uint64] = (*X64)(nil)

// Constants for decimal64 format according to IEEE 754-2008
const (
	// eLimit64 is the maximum encoded exponent value (3 * 2^ecbits - 1)
	eLimit64 int16 = 767 // 3 * 2^8 - 1
	// eMax64 is the maximum decoded exponent value ((Elimit/2) + 1)
	eMax64 int16 = 384 // (767/2) + 1
	// eMin64 is the minimum decoded exponent value (-Elimit/2)
	eMin64 int16 = -383 // -767/2
	// eTiny64 is the exponent of the smallest possible subnormal (Emin - (precision-1))
	eTiny64 int16 = -398 // -383 - (16-1)
	// bias64 is the value to add to decoded exponent to get encoded exponent (-Emin + precision - 1)
	bias64 int16 = 398 // -(-383) + 16 - 1
	// maxCoefficient64 is the maximum coefficient value (10^precision - 1)
	maxCoefficient64 uint64 = 9999999999999999 // 10^16 - 1
)

func (x *X64) Pack(k kind, sign signc, exp int16, coe uint64) error {
	if x == nil {
		return fmt.Errorf("nil receiver")
	}

	return x.pack(k, sign, exp, coe)
}

// pack implements the packed interface by encoding components into BID format.
// According to IEEE 754-2008, decimal64 has:
// - 1 bit for sign
// - 10 bits for combination field (including bits from exponent and coefficient)
// - 53 bits for remaining coefficient bits
func (x *X64) pack(k kind, sign signc, exp int16, coe uint64) error {
	// Validate inputs
	if sign != signc_negative && sign != signc_positive {
		return newInternalError(sign, "invalid sign")
	}

	if coe > maxCoefficient64 && k == kind_finite {
		return newInternalError(coe, "coefficient overflow")
	}

	if (exp > eMax64 || exp < eMin64) && k == kind_finite {
		return newInternalError(exp, "exponent out of range")
	}

	// Check for subnormal values (non-zero coefficient with minimum exponent)
	// and return a signaling NaN immediately
	if k == kind_finite && exp == eMin64 && coe > 0 {
		// Set as signaling NaN
		x.uint64 = 0x7E00000000000000
		if sign == signc_negative {
			x.uint64 |= 1 << 63
		}
		return nil
	}

	// Start with zero
	var result uint64 = 0

	// Set sign bit (bit 63)
	if sign == signc_negative {
		result |= 1 << 63
	}

	// Process based on kind
	switch k {
	case kind_finite:
		// Add bias to get encoded exponent
		biasedExp := uint64(exp + bias64)

		// Check if coefficient fits in 53 bits (2^53 = 9007199254740992)
		if coe < (1 << 53) {
			// Normal format: G0..G9=eeeeeeeeee, remaining bits are coefficient
			// Exponent bits first (10 bits)
			result |= (biasedExp & 0x3FF) << 53
			// Then coefficient bits
			result |= coe & 0x1FFFFFFFFFFFFF
		} else {
			// Large coefficient - need to use alternative encoding
			// Set first 2 bits of exp in combination field
			result |= ((biasedExp >> 8) & 0x3) << 61
			// Set special pattern 11 to indicate this format
			result |= 3 << 59
			// Set remaining 8 bits of exponent
			result |= (biasedExp & 0xFF) << 51
			// Set coefficient bits
			result |= coe & 0x7FFFFFFFFFFFF
		}

	case kind_infinity:
		// Infinity: G0..G4=11110, G5..G9=0
		result |= 0x7800000000000000

	case kind_quiet:
		// Quiet NaN: G0..G4=11111, G5=0
		result |= 0x7C00000000000000

	case kind_signaling:
		// Signaling NaN: G0..G4=11111, G5=1
		result |= 0x7E00000000000000

	default:
		return newInternalError(k, "invalid kind")
	}

	// Store the result
	x.uint64 = result
	return nil
}

// unpack implements the packed interface by decoding BID format into components.
func (x *X64) unpack() (kind, signc, int16, uint64, error) {
	if x == nil {
		return kind_signaling, signc_error, 0, 0, newInternalError(nil, "nil receiver")
	}

	// Get the bits
	bits := x.uint64

	// Extract sign (bit 63)
	sign := signc_positive
	if bits&(1<<63) != 0 {
		sign = signc_negative
	}

	// Extract combination field to identify special values
	g0g4 := (bits >> 58) & 0x1F // First 5 bits of combination field

	// Check for special values based on the first 5 bits
	switch g0g4 {
	case 0x1E: // 11110
		// Positive or negative infinity
		return kind_infinity, sign, 0, 0, nil
	case 0x1F: // 11111
		// NaN - determine if quiet or signaling using G5 bit
		if (bits>>57)&0x1 == 1 {
			return kind_signaling, sign, 0, 0, nil
		}
		return kind_quiet, sign, 0, 0, nil
	}

	// Handle normal values
	g0g1 := (bits >> 61) & 0x3 // First 2 bits of combination field

	// Extract exponent and coefficient for finite numbers
	var exp int16
	var coe uint64

	if g0g1 == 0x3 { // Large coefficient format
		// Extract encoded exponent: 2 bits in combination field + 8 bits in exponent continuation field
		encodedExp := int16(((bits >> 61) & 0x3) << 8)
		encodedExp |= int16((bits >> 51) & 0xFF)
		exp = encodedExp - bias64 // Remove bias to get decoded exponent

		// Extract coefficient
		coe = bits & 0x7FFFFFFFFFFFF
	} else {
		// Normal format
		// Extract encoded exponent: 10 bits after sign
		encodedExp := int16((bits >> 53) & 0x3FF)
		exp = encodedExp - bias64 // Remove bias to get decoded exponent

		// Extract coefficient
		coe = bits & 0x1FFFFFFFFFFFFF
	}

	return kind_finite, sign, exp, coe, nil
}

// isZero returns true if the X64 value is zero (positive or negative).
func (x *X64) isZero() bool {
	k, _, _, coe, err := x.unpack()
	if err != nil || k != kind_finite {
		return false
	}
	return coe == 0
}

// isNaN returns true if the X64 value is Not-a-Number (quiet or signaling).
func (x *X64) isNaN() bool {
	k, _, _, _, err := x.unpack()
	if err != nil {
		return false
	}
	return k == kind_quiet || k == kind_signaling
}

// isInf returns true if the X64 value is infinity (positive or negative).
func (x *X64) isInf() bool {
	k, _, _, _, err := x.unpack()
	if err != nil {
		return false
	}
	return k == kind_infinity
}

// Round applies the specified rounding mode to an X64 value to achieve the target precision.
// It implements the rounding behavior defined in IEEE 754-2008.
func (x *X64) Round(mode Rounding, prec Precision) error {
	k, sign, exp, coe, err := x.unpack()
	if err != nil {
		return err
	}

	// Only finite numbers can be rounded
	if k != kind_finite {
		return nil
	}

	// Count digits in coefficient
	digits := countDigits(coe)

	// If we're already at or below the target precision, no rounding needed
	if digits <= uint8(prec) {
		return nil
	}

	// Apply rounding to the coefficient
	newCoe, digitsRemoved := apply(mode, coe, exp, prec, sign)

	// If digits were removed, adjust the exponent
	if digitsRemoved > 0 {
		exp += int16(digitsRemoved)
	}

	// For special cases of subnormal or extreme values
	if exp < eMin64 || exp > eMax64 {
		if exp < eMin64 {
			// If exponent is too small, try to adjust by reducing precision
			// This is a simplification - full subnormal handling would be more complex
			if newCoe == 0 {
				// Zero can be represented with any exponent
				exp = 0
			} else if (newCoe % 10) == 0 {
				// Can shift right to increase exponent
				for exp < eMin64 && (newCoe%10) == 0 {
					newCoe /= 10
					exp++
				}
			}

			// If still too small, return error or set to zero
			if exp < eMin64 {
				if newCoe == 0 {
					return x.pack(kind_finite, sign, 0, 0) // Return zero
				}
				return newInternalError(exp, "exponent out of range")
			}
		} else if exp > eMax64 {
			// If exponent is too large, return infinity
			return x.pack(kind_infinity, sign, 0, 0)
		}
	}

	// Pack the result back
	return x.pack(k, sign, exp, newCoe)
}
