package fixedpoint

import (
	"fmt"
	"strings"
)

// Debug returns a formatted string with debug information about the FixedPoint.
func (a *FixedPoint) Debug() string {
	sign := '+'
	if a.sign {
		sign = '-'
	}
	return fmt.Sprintf(
		"fp{sign: %c, coe: %d, exp: %d, flg: %s, ctx: %s}",
		sign,
		a.coe,
		a.exp,
		a.flg.Debug(),
		a.ctx.Debug())
}

// String returns the human-readable string representation of the FixedPoint.
func (a *FixedPoint) String() string {
	if a == nil {
		return "NaN"
	}
	if a.flg.nan {
		return "NaN"
	}
	if a.flg.inf {
		if a.sign {
			return "-Infinity"
		}
		return "Infinity"
	}

	// Convert the absolute value to a string
	str := fmt.Sprintf("%d", a.coe)

	// Apply the exponent
	exp := int(a.exp)
	if exp >= 0 {
		// Add trailing zeros
		for range exp {
			str += "0"
		}
	} else {
		// Insert decimal point
		expAbs := -exp
		if len(str) <= expAbs {
			// Need to pad with leading zeros
			padding := expAbs - len(str)
			str = "0." + strings.Repeat("0", padding) + str
		} else {
			// Insert decimal point at the correct position
			pos := len(str) - expAbs
			str = str[:pos] + "." + str[pos:]
		}
		// Trim trailing zeros after decimal point
		str = strings.TrimRight(str, "0")
		str = strings.TrimRight(str, ".")
	}

	// Add sign
	if a.sign {
		str = "-" + str
	}

	return str
}

func (f flags) String() string {
	return fmt.Sprintf("f{s: %s, i: %t, n: %t}", f.sig, f.inf, f.nan)
}

func (f flags) Debug() string {
	fs := f.sig.Debug()
	if f.inf {
		fs += ":INF"
	}
	if f.nan {
		fs += ":NAN"
	}
	return fs
}

func (c context) Debug() string {
	return fmt.Sprintf("%s:%v", c.rounding.Debug(), c.precision)
}

func (c context) String() string {
	return fmt.Sprintf("context{rounding: %s, precision: %v}", c.rounding, c.precision)
}
