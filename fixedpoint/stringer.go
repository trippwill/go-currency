package fixedpoint

import (
	"fmt"
)

// Debug returns a formatted string with debug information about the FixedPoint.
func (a *FiniteNumber) Debug() string {
	sign := '+'
	if a.sign {
		sign = '-'
	}
	return fmt.Sprintf(
		"fn{sign: %c, coe: %d, exp: %d, ctx: %s}",
		sign,
		a.coe,
		a.exp,
		a.context.Debug())
}

// Debug returns a formatted string with debug information about the Infinity.
func (a *Infinity) Debug() string {
	sign := '+'
	if a.sign {
		sign = '-'
	}
	return fmt.Sprintf(
		"inf{sign: %c, ctx: %s}",
		sign,
		a.context.Debug())
}

func (a *NaN) Debug() string {
	sign := '+'
	if a.sign {
		sign = '-'
	}

	if diagnostic, ok := DecodePayload(a.diag); ok {
		return fmt.Sprintf(
			"nan{sign: %c, signal: %s, diag: %v}",
			sign,
			a.context.signal,
			diagnostic)
	}

	return fmt.Sprintf(
		"nan{sign: %c, signal: %s, diag: %d}",
		sign,
		a.context.signal,
		a.diag)
}

// String returns a string representation as a decimal number.
// It formats the number according to the context's precision and rounding mode.
func (a *FiniteNumber) String() string {
	if a.IsZero() {
		return "0"
	}

	// Use a separate sign string.
	sign := ""
	if a.sign {
		sign = "-"
	}

	// Convert coefficient to string.
	coe_str := fmt.Sprintf("%d", a.coe)

	// If the exponent is zero or positive, pad with trailing zeros.
	if a.exp >= 0 {
		// Append a.exp zeros to the coefficient string.
		for range int(a.exp) {
			coe_str += "0"
		}
		return sign + coe_str
	}

	// For a negative exponent, insert a decimal point.
	// Let d be the number of digits after the decimal point.
	d := int(-a.exp)

	// If the coefficient string has fewer digits than required, pad with leading zeros.
	if len(coe_str) <= d {
		zeros := ""
		for range d - len(coe_str) {
			zeros += "0"
		}
		return sign + "0." + zeros + coe_str
	}

	// Otherwise, insert the decimal point at the correct location.
	int_part := coe_str[:len(coe_str)-d]
	frac_part := coe_str[len(coe_str)-d:]
	return sign + int_part + "." + frac_part
}

// String returns a string representation of the Infinity.
func (a *Infinity) String() string {
	if a.sign {
		return "-Infinity"
	}
	return "Infinity"
}

// String returns a string representation of the NaN.
func (a *NaN) String() string {
	sign := '+'
	if a.sign {
		sign = '-'
	}
	return fmt.Sprintf("%cNaN{%s}:%v", sign, a.context.signal, a.diag)
}

func (c context) Debug() string {
	return fmt.Sprintf("%s:%v", c.rounding.Debug(), c.precision)
}

func (c context) String() string {
	return fmt.Sprintf("context{rounding: %s, precision: %v}", c.rounding, c.precision)
}
