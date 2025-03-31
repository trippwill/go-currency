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
	return fmt.Sprintf("fn{%c, %d, %d}", sign, a.coe, a.exp)
}

// Debug returns a formatted string with debug information about the Infinity.
func (a *Infinity) Debug() string {
	sign := '+'
	if a.sign {
		sign = '-'
	}
	return fmt.Sprintf("inf{%c}", sign)
}

func (a *NaN) Debug() string {
	sign := '+'
	if a.sign {
		sign = '-'
	}

	if diagnostic, ok := DecodePayload(a.diag); ok {
		return fmt.Sprintf("nan{%c, %v}", sign, diagnostic)
	}

	return fmt.Sprintf("nan{%c, %d}", sign, a.diag)
}

// String returns a string representation as a decimal number.
func (fn *FiniteNumber) String() string {
	if fn == nil {
		return "nil"
	}

	a := &FiniteNumber{
		sign: fn.sign,
		coe:  fn.coe,
		exp:  fn.exp,
	}

	prec := dlen(a.coe)
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
	return fmt.Sprintf("%cNaN{%v}", sign, a.diag)
}

func (c Context) Debug() string {
	return fmt.Sprintf("%s:%v:%s[%s]", c.rounding.Debug(), c.precision, c.signals.Debug(), c.traps.Debug())
}

func (c Context) String() string {
	return fmt.Sprintf("context{precision: %v, rounding: %s,  traps: %s}", c.precision, c.rounding, c.traps)
}
