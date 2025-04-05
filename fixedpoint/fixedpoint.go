// Base on the General Decimal Arithmetic Specification 1.70 – 7 Apr 2009
// https://speleotrove.com/decimal/decarith.html
package fixedpoint

// type FixedPoint[X X64 | X32] interface {
// 	fmt.Stringer
// 	Add(X) X
// 	Sub(X) X
// 	Mul(X) X
// 	Div(X) X
// }

// var (
//
//	_ FixedPoint[X64] = (*X64)(nil)
//	_ FixedPoint[X32] = (*X32)(nil)
//
// )
type X64 struct {
	uint64
}

type X32 struct {
	uint32
}

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

var (
	_ packed[int16, uint64] = (*X64)(nil)
	_ packed[int8, uint32]  = (*X32)(nil)
)

type (
	spec uint8
)

const (
	mask_spec spec = 0b1_1111_110
)

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
