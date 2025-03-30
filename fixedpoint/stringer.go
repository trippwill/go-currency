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
		"fn{%c, %d, %d, %s}",
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
		"inf{%c, ctx: %s}",
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
			"nan{%c, signal: %s, diag: %v}",
			sign,
			a.context.signal,
			diagnostic)
	}

	return fmt.Sprintf(
		"nan{%c, signal: %s, diag: %d}",
		sign,
		a.context.signal,
		a.diag)
}

// String returns a string representation as a decimal number.
func (fn *FiniteNumber) String() string {
	if fn == nil {
		return "nil"
	}

	_fn := &FiniteNumber{
		sign:    fn.sign,
		coe:     fn.coe,
		exp:     fn.exp,
		context: fn.context,
	}

	var a *FiniteNumber
	_fp := apply_rounding(_fn)
	switch v := _fp.(type) {
	case *FiniteNumber:
		a = v
	default:
		return v.String()
	}

	prec := int(a.context.precision)
	// Fast path for zero coefficient.
	if a.coe == 0 {
		if prec > 1 {
			return "0." + strings.Repeat("0", prec-1)
		}
		return "0"
	}

	coe_str := fmt.Sprintf("%d", a.coe)
	pos := len(coe_str) + int(a.exp)
	var int_part, frac_part string

	if a.exp < 0 {
		if pos <= 0 {
			int_part = "0"
			frac_part = strings.Repeat("0", -pos) + coe_str
		} else {
			int_part = coe_str[:pos]
			frac_part = coe_str[pos:]
		}
	} else if a.exp > 0 {
		int_part = coe_str + strings.Repeat("0", int(a.exp))
		frac_part = ""
	} else {
		int_part = coe_str
		frac_part = ""
	}

	frac_part = strings.TrimRight(frac_part, "0")
	if frac_part == "" {
		frac_part = "0"
	}

	// Build result and trim unnecessary zeros.
	result := int_part + "." + frac_part
	result = strings.TrimRight(result, "0")
	// Ensure the numeric part (excluding the decimal point) meets the precision.
	if len(result)-1 < prec {
		result += strings.Repeat("0", prec-(len(result)-1))
	}

	// Apply sign: note that a.sign being true implies a negative number.
	if a.sign {
		result = "-" + result
	}

	return result
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
