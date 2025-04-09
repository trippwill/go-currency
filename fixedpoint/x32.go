package fixedpoint

// X32 implements the IEEE 754-2008 decimal32 format
// using Binary Integer Decimal (BID) encoding
type X32 struct {
	uint32
}

var _ packed[int8, uint32] = (*X32)(nil)

// Constants for decimal32 format according to IEEE 754-2008
const (
	// eLimit32 is the maximum encoded exponent value (3 * 2^ecbits - 1)
	eLimit32 int16 = 191 // 3 * 2^6 - 1
	// eMax32 is the maximum decoded exponent value ((Elimit/2) + 1)
	eMax32 int8 = 96 // (191/2) + 1
	// eMin32 is the minimum decoded exponent value (-Elimit/2)
	eMin32 int8 = -95 // -191/2
	// eTiny32 is the exponent of the smallest possible subnormal (Emin - (precision-1))
	eTiny32 int16 = -101 // -95 - (7-1)
	// bias32 is the value to add to decoded exponent to get encoded exponent (-Emin + precision - 1)
	bias32 int16 = 101 // -(-95) + 7 - 1
	// maxCoefficient32 is the maximum coefficient value (10^precision - 1)
	maxCoefficient32 uint32 = 9999999 // 10^7 - 1
)

// pack implements the packed interface by encoding components into BID format.
// According to IEEE 754-2008, decimal32 has:
// - 1 bit for sign
// - 11 bits for combination field (including bits from exponent and coefficient)
// - 20 bits for remaining coefficient bits
func (x *X32) pack(k kind, sign signc, exp int8, coe uint32) error {
	// Validate inputs
	if sign != signc_negative && sign != signc_positive {
		return newInternalError(sign, "invalid sign")
	}

	if coe > maxCoefficient32 && k == kind_finite {
		return newInternalError(coe, "coefficient overflow")
	}

	if (exp > eMax32 || exp < eMin32) && k == kind_finite {
		return newInternalError(exp, "exponent out of range")
	}

	// Check for subnormal values (non-zero coefficient with minimum exponent)
	// and return a signaling NaN immediately
	if k == kind_finite && exp == eMin32 && coe > 0 {
		// Set as signaling NaN
		x.uint32 = 0x7E000000
		if sign == signc_negative {
			x.uint32 |= 1 << 31
		}
		return nil
	}

	// Start with zero
	var result uint32 = 0

	// Set sign bit (bit 31)
	if sign == signc_negative {
		result |= 1 << 31
	}

	// Process based on kind
	switch k {
	case kind_finite:
		// Add bias to get encoded exponent
		biasedExp := uint32(int16(exp) + bias32)

		// Check if coefficient fits in 20 bits (2^20 = 1048576)
		if coe < (1 << 20) {
			// Normal format: G0..G10=eeeeeeeeeee, remaining bits are coefficient
			// Exponent bits first (8 bits)
			result |= (biasedExp & 0xFF) << 23
			// Then coefficient bits
			result |= coe & 0xFFFFF
		} else {
			// Large coefficient - need to use alternative encoding
			// Set first 2 bits of exp in combination field
			result |= ((biasedExp >> 6) & 0x3) << 29
			// Set special pattern 11 to indicate this format
			result |= 3 << 27
			// Set remaining 6 bits of exponent
			result |= (biasedExp & 0x3F) << 21
			// Set coefficient bits
			result |= coe & 0x1FFFFF
		}

	case kind_infinity:
		// Infinity: G0..G4=11110, G5..G10=0
		result |= 0x78000000

	case kind_quiet:
		// Quiet NaN: G0..G4=11111, G5=0
		result |= 0x7C000000

	case kind_signaling:
		// Signaling NaN: G0..G4=11111, G5=1
		result |= 0x7E000000

	default:
		return newInternalError(k, "invalid kind")
	}

	// Store the result
	x.uint32 = result
	return nil
}

// unpack implements the packed interface by decoding BID format into components.
func (x *X32) unpack() (kind, signc, int8, uint32, error) {
	if x == nil {
		return kind_signaling, signc_error, 0, 0, newInternalError(nil, "nil receiver")
	}

	// Get the bits
	bits := x.uint32

	// Extract sign (bit 31)
	sign := signc_positive
	if bits&(1<<31) != 0 {
		sign = signc_negative
	}

	// Extract combination field to identify special values
	g0g4 := (bits >> 26) & 0x1F // First 5 bits of combination field

	// Check for special values based on the first 5 bits
	switch g0g4 {
	case 0x1E: // 11110
		// Positive or negative infinity
		return kind_infinity, sign, 0, 0, nil
	case 0x1F: // 11111
		// NaN - determine if quiet or signaling using G5 bit
		if (bits>>25)&0x1 == 1 {
			return kind_signaling, sign, 0, 0, nil
		}
		return kind_quiet, sign, 0, 0, nil
	}

	// Handle normal values
	g0g1 := (bits >> 29) & 0x3 // First 2 bits of combination field

	// Extract exponent and coefficient for finite numbers
	var exp int8
	var coe uint32

	if g0g1 == 0x3 { // Large coefficient format
		// Extract encoded exponent: 2 bits in combination field + 6 bits in exponent continuation field
		encodedExp := int16(((bits >> 29) & 0x3) << 6)
		encodedExp |= int16((bits >> 21) & 0x3F)
		exp = int8(encodedExp - bias32) // Remove bias to get decoded exponent

		// Extract coefficient
		coe = bits & 0x1FFFFF
	} else {
		// Normal format
		// Extract encoded exponent: 8 bits after sign
		encodedExp := int16((bits >> 23) & 0xFF)
		exp = int8(encodedExp - bias32) // Remove bias to get decoded exponent

		// Extract coefficient
		coe = bits & 0xFFFFF
	}

	return kind_finite, sign, exp, coe, nil
}

// isZero returns true if the X32 value is zero (positive or negative).
func (x *X32) isZero() bool {
	k, _, _, coe, err := x.unpack()
	if err != nil || k != kind_finite {
		return false
	}
	return coe == 0
}

// isNaN returns true if the X32 value is Not-a-Number (quiet or signaling).
func (x *X32) isNaN() bool {
	k, _, _, _, err := x.unpack()
	if err != nil {
		return false
	}
	return k == kind_quiet || k == kind_signaling
}

// isInf returns true if the X32 value is infinity (positive or negative).
func (x *X32) isInf() bool {
	k, _, _, _, err := x.unpack()
	if err != nil {
		return false
	}
	return k == kind_infinity
}

// Round applies the specified rounding mode to an X32 value to achieve the target precision.
// It implements the rounding behavior defined in IEEE 754-2008.
func (x *X32) Round(mode Rounding, precision uint) error {
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
	if digits <= precision {
		return nil
	}

	// Apply rounding to the coefficient
	newCoe, digitsRemoved := Apply(mode, coe, exp, precision, sign)

	// If digits were removed, adjust the exponent
	if digitsRemoved > 0 {
		exp += int8(digitsRemoved)
	}

	// For special cases of subnormal or extreme values
	if exp < eMin32 || exp > eMax32 {
		if exp < eMin32 {
			// If exponent is too small, try to adjust by reducing precision
			// This is a simplification - full subnormal handling would be more complex
			if newCoe == 0 {
				// Zero can be represented with any exponent
				exp = 0
			} else if (newCoe % 10) == 0 {
				// Can shift right to increase exponent
				for exp < eMin32 && (newCoe%10) == 0 {
					newCoe /= 10
					exp++
				}
			}

			// If still too small, return error or set to zero
			if exp < eMin32 {
				if newCoe == 0 {
					return x.pack(kind_finite, sign, 0, 0) // Return zero
				}
				return newInternalError(exp, "exponent out of range")
			}
		} else if exp > eMax32 {
			// If exponent is too large, return infinity
			return x.pack(kind_infinity, sign, 0, 0)
		}
	}

	// Pack the result back
	return x.pack(k, sign, exp, newCoe)
}
