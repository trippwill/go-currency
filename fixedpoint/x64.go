package fixedpoint

const (
	mask_trailing_sig64 uint64 = 0x3_FFFF_FFFF_FFFF // Trailing significand (lowest 50 bits)
	shift_spec64        uint64 = 56
	exp_bias64                 = 398                   // Bias for 64-bit exponent
	exp_limit64                = 767                   // Maximum encoded exponent for 64-bit
	exp_min64                  = -383                  // Minimum exponent value for 64-bit
	exp_max64                  = 384                   // Maximum exponent value for 64-bit
	coe_dlen64                 = 16                    // Number of digits in significand
	coe_maxv64                 = 9_999_999_999_999_999 // Maximum value for significand
)

func (x *X64) Pack(kind kind, sign sign, exp int16, coe uint64) error {
	if x == nil {
		return newInternalError(x, "pack nil")
	}
	return x.pack(kind, sign, exp, coe)
}

func (x *X64) unpack() (kind kind, sign sign, exp int16, coe uint64, err error) {
	if x == nil {
		err = newInternalError(x, "unpack nil")
		return
	}

	// Create a copy for safety
	data := x.uint64

	msb := spec(data >> shift_spec64) // Most significant byte
	sign, kind, form, err := decode_spec(msb)
	if err != nil || sign == sign_error {
		return
	}

	switch kind {
	case kind_infinity:
		// Only the sign is relevant for infinity
		return kind, sign, 0, 0, nil
	case kind_quiet, kind_signaling:
		// NaNs may encode diagnostic information in the trailing significand
		coe = (data & mask_trailing_sig64)
		return kind, sign, 0, coe, nil
	case kind_finite:
		if form == form_none {
			err = newInternalError(x, "invalid form for finite number")
			return
		}
	}

	exp, coe, err = decode64_exp_coe(form == form_large, data)
	return
}

// decode64_exp_coe decodes the significand and exponent from the given data.
// For small form, then the exponent field consists of the 10 bits following the sign bit,
// and the significand is the remaining 53 bits, with an implicit leading 0 bit.
// For large form, the 10-bit exponent field is shifted 2 bits to the right (after both the sign bit and the "11" bits thereafter),
// and the represented significand is in the remaining 51 bits. In this case there is an implicit (that is, not stored)
// leading 3-bit sequence "100" for the MSB bits of the true significand.
func decode64_exp_coe(large bool, data uint64) (int16, uint64, error) {
	if large {
		// Large form: extract exponent and significand
		exp := int16((data >> 51) & 0x3FF)                  // 10 bits for exponent
		coe := (data & 0x7FFFFFFFFFFFF)                     // 51 bits for significand
		return exp - exp_bias64, coe | 0x4000000000000, nil // Add implicit leading "100"
	} else {
		// Small form: extract exponent and significand
		exp := int16((data >> 53) & 0x3FF) // 10 bits for exponent
		coe := (data & 0x1FFFFFFFFFFFFF)   // 53 bits for significand
		return exp - exp_bias64, coe, nil
	}
}

func (x *X64) pack(kind kind, sign sign, exp int16, coe uint64) error {
	if x == nil {
		return newInternalError(x, "pack nil")
	}

	var data uint64

	// Set the sign bit [63:63]
	switch sign {
	case sign_negative:
		data |= (0b1 << 63)
	case sign_positive:
		data = 0 // Clear the sign bit
	default:
		return newInternalError(x, "invalid sign")
	}

	// Set the kind bits [62:56]
	switch kind {
	case kind_infinity:
		data |= (0b0_1111_000 << 56) // Set the kind bits for infinity
		*x = X64{data}
		return nil
	case kind_quiet:
		data |= (0b0_1111_100 << 56)
		data |= (coe & mask_trailing_sig64) // Set the kind bits for quiet NaN
		*x = X64{data}
		return nil
	case kind_signaling:
		data |= (0b0_1111_110 << 56)
		data |= (coe & mask_trailing_sig64) // Set the kind bits for signaling NaN
		*x = X64{data}
		return nil
	}

	// Encode finite numbers
	msb4 := (coe >> 49) & 0xF // Extract the most significant 4 bits of the significand
	large := msb4 >= 0x8      // Check if MSB4 is in the range 1000 (binary) or higher
	if large {                // Large form
		data |= (uint64(exp+exp_bias64) & 0x3FF) << 51 // Encode 10-bit exponent shifted 2 bits right
		data |= (coe & 0x7FFFFFFFFFFFF)                // Encode 51-bit significand
		data |= 0b011 << 60                            // Set the large form bits
	} else {
		data |= (uint64(exp+exp_bias64) & 0x3FF) << 53 // Encode 10-bit exponent
		data |= (coe & 0x1FFFFFFFFFFFFF)               // Encode 53-bit significand
	}

	*x = X64{data}
	return nil
}
