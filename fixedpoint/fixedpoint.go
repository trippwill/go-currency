// Base on the General Decimal Arithmetic Specification 1.70 – 7 Apr 2009
// https://speleotrove.com/decimal/decarith.html
package fixedpoint

import "fmt"

type FixedPoint[X X64 | X32] interface {
	// fmt.Stringer
	// Add(X) X
	// Sub(X) X
	// Mul(X) X
	// Div(X) X
}

var (
	_ FixedPoint[X64] = (*X64)(nil)
	_ FixedPoint[X32] = (*X32)(nil)
)

type (
	X64 uint64
	X32 uint32
)

type sign int8

const (
	sign_negative sign = -1 // Negative
	sign_positive sign = 1  // Positive
	sign_error    sign = 0  // Error
)

type kind uint8

const (
	kind_signaling kind = iota // Signaling NaN
	kind_quiet                 // Quiet NaN
	kind_infinity              // flag_infinity
	kind_finite                // flag_finite
)

type form uint8

const (
	form_none  form = iota // Not finite form
	form_small             // Small form
	form_large             // Large form
)

// FixedPoint implements the Decimal Arithmetic Encoding for 128-bit decimal numbers.
// See https://speleotrove.com/decimal/decbits.html

type packed[E int8 | int16, C uint32 | uint64] interface {
	pack(kind kind, sign sign, exp E, coe C) error
	unpack() (kind, sign, E, C, error)
}

var _ packed[int16, uint64] = (*X64)(nil)

// _ packed[int8, uint32]  = (*X32)(nil)

type internalError struct {
	data any
	msg  string
}

func (e *internalError) Error() string {
	return fmt.Sprintf("internal error: %s: %v", e.msg, e.data)
}

func newInternalError(data any, msg string) error {
	return &internalError{
		data: data,
		msg:  msg,
	}
}

type (
	spec uint8
)

const (
	mask_spec           spec   = 0b1_1111_110
	mask_trailing_sig64 uint64 = 0x3_FFFF_FFFF_FFFF // Trailing significand (lowest 50 bits)
	mask_trailing_sig32 uint32 = 0x1F_FFFF          // Trailing significand (lowest 20 bits)
	shift_spec64               = 56
	shift_spec32               = 24
	exp_bias64                 = 398  // Bias for 64-bit exponent
	exp_bias32                 = 101  // Bias for 32-bit exponent
	exp_limit64                = 767  // Maximum encoded exponent for 64-bit
	exp_limit32                = 191  // Maximum encoded exponent for 32-bit
	exp_min64                  = -383 // Minimum exponent value for 64-bit
	exp_min32                  = -95  // Minimum exponent value for 32-bit
	exp_max64                  = 384  // Maximum exponent value for 64-bit
	exp_max32                  = 96   // Maximum exponent value for 32-bit
)

func (x *X64) unpack() (kind kind, sign sign, exp int16, coe uint64, err error) {
	if x == nil {
		err = newInternalError(x, "unpack nil")
		return
	}

	// Create a copy for safety
	data := uint64(*x)

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

func decode_spec(msb spec) (sign, kind, form, error) {
	if msb == 0 {
		return sign_positive, kind_finite, form_small, nil
	}

	// Check for special values

	switch msb & mask_spec {
	case 0b1_1111_000:
		return sign_negative, kind_infinity, form_none, nil
	case 0b0_1111_000:
		return sign_positive, kind_infinity, form_none, nil
	case 0b1_1111_100:
		return sign_negative, kind_quiet, form_none, nil
	case 0b0_1111_100:
		return sign_positive, kind_quiet, form_none, nil
	case 0b1_1111_110:
		return sign_negative, kind_signaling, form_none, nil
	case 0b0_1111_110:
		return sign_positive, kind_signaling, form_none, nil
	}

	// The following patterns are for finite numbers

	form_mask := spec(0b1_11_11_000)
	switch msb & form_mask {
	case 0b1_11_00_000:
		fallthrough
	case 0b1_11_01_000:
		fallthrough
	case 0b1_11_10_000:
		return sign_negative, kind_finite, form_large, nil
	case 0b0_11_00_000:
		fallthrough
	case 0b0_11_01_000:
		fallthrough
	case 0b0_11_10_000:
		return sign_positive, kind_finite, form_large, nil
	}

	form_mask = spec(0b1_11_00_000)
	switch msb & form_mask {
	case 0b1_00_00_000:
		fallthrough
	case 0b1_01_00_000:
		fallthrough
	case 0b1_10_00_000:
		return sign_negative, kind_finite, form_small, nil
	case 0b0_00_00_000:
		fallthrough
	case 0b0_01_00_000:
		fallthrough
	case 0b0_10_00_000:
		return sign_positive, kind_finite, form_small, nil
	}
	// If we reach here, the spec is not recognized
	return sign_error, kind_signaling, form_none, newInternalError(msb, "invalid spec")
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

func (x *X64) Pack(kind kind, sign sign, exp int16, coe uint64) error {
	if x == nil {
		return newInternalError(x, "pack nil")
	}
	return x.pack(kind, sign, exp, coe)
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
		*x = X64(data)
		return nil
	case kind_quiet:
		data |= (0b0_1111_100 << 56)
		data |= (coe & mask_trailing_sig64) // Set the kind bits for quiet NaN
		*x = X64(data)
		return nil
	case kind_signaling:
		data |= (0b0_1111_110 << 56)
		data |= (coe & mask_trailing_sig64) // Set the kind bits for signaling NaN
		*x = X64(data)
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

	*x = X64(data)
	return nil
}

func (x *X32) unpack() (kind kind, sign sign, exp int8, coe uint32, err error) {
	if x == nil {
		err = newInternalError(x, "unpack nil")
		return
	}

	// Create a copy for safety
	data := uint32(*x)

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
		*x = X32(data)
		return nil
	case kind_quiet:
		data |= (0b0_1111_100 << 24)
		data |= (coe & uint32(mask_trailing_sig32)) // Set the kind bits for quiet NaN
		*x = X32(data)
		return nil
	case kind_signaling:
		data |= (0b0_1111_110 << 24)
		data |= (coe & uint32(mask_trailing_sig32)) // Set the kind bits for signaling NaN
		*x = X32(data)
		return nil
	}

	// Encode finite numbers
	msb4 := (coe >> 21) & 0xF // Extract the most significant 4 bits of the significand
	large := msb4 >= 0x8      // Check if MSB4 is in the range 1000 (binary) or higher
	if large {                // Large form
		data |= (uint32(exp+exp_bias32) & 0x3FF) << 21 // Encode 10-bit exponent shifted 2 bits right
		data |= (coe & 0x1FFFFF)                       // Encode 21-bit significand
		data |= 0b011 << 28                            // Set the large form bits
	} else {
		data |= (uint32(exp+exp_bias32) & 0x3FF) << 23 // Encode 10-bit exponent
		data |= (coe & 0x7FFFFF)                       // Encode 23-bit significand
	}

	*x = X32(data)
	return nil
}
