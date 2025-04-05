package fixedpoint

const (
	mask_trailing_sig32 uint32 = 0x1F_FFFF // Trailing significand (lowest 20 bits)
	shift_spec32        uint32 = 24
	exp_bias32                 = 101       // Bias for 32-bit exponent
	exp_limit32                = 191       // Maximum encoded exponent for 32-bit
	exp_min32                  = -95       // Minimum exponent value for 32-bit
	exp_max32                  = 96        // Maximum exponent value for 32-bit
	coe_dlen32                 = 7         // Number of digits in significand
	coe_maxv32                 = 9_999_999 // Maximum value for significand
)

func (x *X32) unpack() (kind kind, sign sign, exp int8, coe uint32, err error) {
	if x == nil {
		err = newInternalError(x, "unpack nil")
		return
	}

	// Create a copy for safety
	data := x.uint32

	msb := spec(data >> shift_spec32) // Most significant byte
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
		coe = uint32(data & mask_trailing_sig32)
		return kind, sign, 0, coe, nil
	case kind_finite:
		if form == form_none {
			err = newInternalError(x, "invalid form for finite number")
			return
		}
	}

	exp, coe, err = decode32_exp_coe(form == form_large, data)
	return
}

func decode32_exp_coe(large bool, data uint32) (int8, uint32, error) {
	if large {
		// Large form: extract exponent and significand
		exp := int8((data >> 25) & 0x3FF)            // 10 bits for exponent
		coe := (data & 0x1FFFFF)                     // 21 bits for significand
		return exp - exp_bias32, coe | 0x200000, nil // Add implicit leading "10"
	} else {
		// Small form: extract exponent and significand
		exp := int8((data >> 23) & 0x3FF) // 10 bits for exponent
		coe := (data & 0x7FFFFF)          // 23 bits for significand
		return exp - exp_bias32, coe, nil
	}
}

func (x *X32) pack(kind kind, sign sign, exp int8, coe uint32) error {
	if x == nil {
		return newInternalError(x, "pack nil")
	}

	var data uint32

	// Set the sign bit [31:31]
	switch sign {
	case sign_negative:
		data |= (0b1 << 31)
	case sign_positive:
		data = 0 // Clear the sign bit
	default:
		return newInternalError(x, "invalid sign")
	}

	// Set the kind bits [30:24]
	switch kind {
	case kind_infinity:
		data |= (0b0_1111_000 << 24) // Set the kind bits for infinity
		*x = X32{data}
		return nil
	case kind_quiet:
		data |= (0b0_1111_100 << 24)
		data |= (coe & uint32(mask_trailing_sig32)) // Set the kind bits for quiet NaN
		*x = X32{data}
		return nil
	case kind_signaling:
		data |= (0b0_1111_110 << 24)
		data |= (coe & uint32(mask_trailing_sig32)) // Set the kind bits for signaling NaN
		*x = X32{data}
		return nil
	}

	// Encode finite numbers
	msb4 := (coe >> 21) & 0xF // Extract the most significant 4 bits of the significand
	large := msb4 >= 0x8      // Check if MSB4 is in the range 1000 (binary) or higher
	if large {
		data |= (uint32(exp+exp_bias32) & 0x3FF) << 21 // Encode 10-bit exponent shifted 2 bits right
		data |= (coe & 0x1FFFFF)                       // Encode 21-bit significand
		data |= 0b011 << 28                            // Set the large form bits
	} else {
		data |= (uint32(exp+exp_bias32) & 0x3FF) << 23 // Encode 10-bit exponent
		data |= (coe & 0x7FFFFF)                       // Encode 23-bit significand
	}

	*x = X32{data}
	return nil
}
