package fixedpoint

import (
	"fmt"
	"strings"
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

	sign := ""
	if a.sign {
		sign = "-"
	}

	prec := int(a.context.precision)
	coe_str := fmt.Sprintf("%d", a.coe)

	if a.exp >= 0 {
		coe_str += strings.Repeat("0", int(a.exp))
		if prec > 0 {
			frac := strings.Repeat("0", prec)
			return sign + coe_str + "." + frac
		}
		return sign + coe_str
	}

	pos := len(coe_str) + int(a.exp)
	var int_part, frac_part string
	if pos <= 0 {
		int_part = "0"
		frac_part = strings.Repeat("0", -pos) + coe_str
	} else {
		int_part = coe_str[:pos]
		frac_part = coe_str[pos:]
	}

	if len(frac_part) > prec {
		frac_part = frac_part[:prec]
	} else {
		frac_part += strings.Repeat("0", prec-len(frac_part))
	}

	if prec > 0 {
		return sign + int_part + "." + frac_part
	}
	return sign + int_part
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
