package fixedpoint

import (
	"fmt"
	"strconv"
	"strings"
)

// String implements the fmt.Stringer interface for X64.
// It returns a human-readable representation of the X64 decimal floating-point number.
func (x X64) String() string {
	k, sign, exp, coe, err := x.unpack()
	if err != nil {
		return fmt.Sprintf("X64{ERROR: %v}", err)
	}

	switch k {
	case kind_quiet:
		if sign == signc_negative {
			return "-qNaN"
		}
		return "qNaN"
	case kind_signaling:
		if sign == signc_negative {
			return "-sNaN"
		}
		return "sNaN"
	case kind_infinity:
		if sign == signc_negative {
			return "-Infinity"
		}
		return "Infinity"
	}

	// For finite numbers, handle the sign, coefficient, and exponent
	signStr := ""
	if sign == signc_negative {
		signStr = "-"
	}

	// If coefficient is zero, return "0"
	if coe == 0 {
		return signStr + "0"
	}

	coeStr := strconv.FormatUint(coe, 10)

	// Apply scientific notation if the exponent is out of a reasonable range
	absExp := exp
	if absExp < 0 {
		absExp = -absExp
	}

	if absExp > 6 {
		// Scientific notation: c.ccc...e±exp
		digits := len(coeStr)
		adjExp := exp + int16(digits-1)

		// Format the coefficient with decimal point
		var formatted string
		if digits > 1 {
			formatted = coeStr[:1] + "." + coeStr[1:]
		} else {
			formatted = coeStr + ".0"
		}

		// Trim trailing zeros after decimal point, but keep at least one digit
		parts := strings.Split(formatted, ".")
		parts[1] = strings.TrimRight(parts[1], "0")
		if parts[1] == "" {
			parts[1] = "0"
		}
		formatted = parts[0] + "." + parts[1]

		return fmt.Sprintf("%s%se%+d", signStr, formatted, adjExp)
	}

	// Regular decimal notation
	if exp >= 0 {
		// Positive exponent - append zeros
		return signStr + coeStr + strings.Repeat("0", int(exp))
	} else {
		// Negative exponent - insert decimal point
		absExp := int(-exp)
		if absExp >= len(coeStr) {
			// Need to prepend zeros: 0.000ccc
			zeros := strings.Repeat("0", absExp-len(coeStr))
			return signStr + "0." + zeros + coeStr
		} else {
			// Insert decimal point: cc.ccc
			pos := len(coeStr) - absExp
			return signStr + coeStr[:pos] + "." + coeStr[pos:]
		}
	}
}

// Debug returns a debug representation of the X64 value showing the internal components.
func (x X64) Debug() string {
	k, sign, exp, coe, err := x.unpack()
	if err != nil {
		return fmt.Sprintf("X64{ERROR: %v}", err)
	}

	signChar := '+'
	if sign == signc_negative {
		signChar = '-'
	}

	switch k {
	case kind_quiet:
		return fmt.Sprintf("X64{qNaN, %c}", signChar)
	case kind_signaling:
		return fmt.Sprintf("X64{sNaN, %c}", signChar)
	case kind_infinity:
		return fmt.Sprintf("X64{Inf, %c}", signChar)
	default:
		return fmt.Sprintf("X64{%c, %d, %d}", signChar, coe, exp)
	}
}

// String implements the fmt.Stringer interface for X32.
// It returns a human-readable representation of the X32 decimal floating-point number.
func (x X32) String() string {
	k, sign, exp, coe, err := x.unpack()
	if err != nil {
		return fmt.Sprintf("X32{ERROR: %v}", err)
	}

	switch k {
	case kind_quiet:
		if sign == signc_negative {
			return "-qNaN"
		}
		return "qNaN"
	case kind_signaling:
		if sign == signc_negative {
			return "-sNaN"
		}
		return "sNaN"
	case kind_infinity:
		if sign == signc_negative {
			return "-Infinity"
		}
		return "Infinity"
	}

	// For finite numbers, handle the sign, coefficient, and exponent
	signStr := ""
	if sign == signc_negative {
		signStr = "-"
	}

	// If coefficient is zero, return "0"
	if coe == 0 {
		return signStr + "0"
	}

	coeStr := strconv.FormatUint(uint64(coe), 10)

	// Apply scientific notation if the exponent is out of a reasonable range
	absExp := exp
	if absExp < 0 {
		absExp = -absExp
	}

	if absExp > 6 {
		// Scientific notation: c.ccc...e±exp
		digits := len(coeStr)
		adjExp := exp + int8(digits-1)

		// Format the coefficient with decimal point
		var formatted string
		if digits > 1 {
			formatted = coeStr[:1] + "." + coeStr[1:]
		} else {
			formatted = coeStr + ".0"
		}

		// Trim trailing zeros after decimal point, but keep at least one digit
		parts := strings.Split(formatted, ".")
		parts[1] = strings.TrimRight(parts[1], "0")
		if parts[1] == "" {
			parts[1] = "0"
		}
		formatted = parts[0] + "." + parts[1]

		return fmt.Sprintf("%s%se%+d", signStr, formatted, adjExp)
	}

	// Regular decimal notation
	if exp >= 0 {
		// Positive exponent - append zeros
		return signStr + coeStr + strings.Repeat("0", int(exp))
	} else {
		// Negative exponent - insert decimal point
		absExp := int(-exp)
		if absExp >= len(coeStr) {
			// Need to prepend zeros: 0.000ccc
			zeros := strings.Repeat("0", absExp-len(coeStr))
			return signStr + "0." + zeros + coeStr
		} else {
			// Insert decimal point: cc.ccc
			pos := len(coeStr) - absExp
			return signStr + coeStr[:pos] + "." + coeStr[pos:]
		}
	}
}

// Debug returns a debug representation of the X32 value showing the internal components.
func (x X32) Debug() string {
	k, sign, exp, coe, err := x.unpack()
	if err != nil {
		return fmt.Sprintf("X32{ERROR: %v}", err)
	}

	signChar := '+'
	if sign == signc_negative {
		signChar = '-'
	}

	switch k {
	case kind_quiet:
		return fmt.Sprintf("X32{qNaN, %c}", signChar)
	case kind_signaling:
		return fmt.Sprintf("X32{sNaN, %c}", signChar)
	case kind_infinity:
		return fmt.Sprintf("X32{Inf, %c}", signChar)
	default:
		return fmt.Sprintf("X32{%c, %d, %d}", signChar, coe, exp)
	}
}
