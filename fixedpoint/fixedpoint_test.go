package fixedpoint

import (
	"strings"
	"testing"
)

func TestNewFixedPoint(t *testing.T) {
	tests := []struct {
		name         string
		significand  int64
		exp          int16
		want         string
		wantOverflow bool
	}{
		{"zero", 0, 0, "0", false},
		{"positive", 42, 0, "42", false},
		{"negative", -42, 0, "-42", false},
		{"with_exp_positive", 42, 2, "4200", false},
		{"with_exp_negative", 42, -2, "0.42", false},
		{"large_value", 999999999999999999, 0, "999999999999999999", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := New(tt.significand, tt.exp)
			if tt.wantOverflow {
				if (&fp).Signal()&SignalOverflow == 0 {
					t.Errorf("NewFixedPoint(%d, %d) expected overflow signal", tt.significand, tt.exp)
				}
			}
			if got := (&fp).String(); got != tt.want {
				t.Errorf("NewFixedPoint(%d, %d) = %s, want %s", tt.significand, tt.exp, got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       string
		wantSignal Signal
	}{
		{"zero", "0", "0", 0},
		{"positive_integer", "42", "42", 0},
		{"negative_integer", "-42", "-42", 0},
		{"positive_decimal", "42.5", "42.5", 0},
		{"negative_decimal", "-42.5", "-42.5", 0},
		{"leading_zero", "0.42", "0.42", 0},
		{"trailing_zero", "42.0", "42", 0},
		{"multiple_decimals", "42.5.3", "NaN", SignalInvalidOperation},
		{"empty_string", "", "NaN", SignalInvalidOperation},
		{"non_numeric", "abc", "NaN", SignalInvalidOperation},
		{"large_value", "999999999999999999", "999999999999999999", 0},
		{"overflow", "10000000000000000000", "Infinity", SignalOverflow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := Parse(tt.input)
			if got := (&fp).String(); got != tt.want {
				t.Errorf("Parse(%q) = %s, want %s", tt.input, got, tt.want)
			}
			if (&fp).Signal()&tt.wantSignal != tt.wantSignal {
				t.Errorf("Parse(%q) signal = %v, want %v", tt.input, (&fp).Signal(), tt.wantSignal)
			}
		})
	}
}

func TestFixedPoint_Add(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want string
	}{
		{"zero_add_zero", "0", "0", "0"},
		{"positive_add_positive", "42", "58", "100"},
		{"negative_add_negative", "-42", "-58", "-100"},
		{"positive_add_negative", "100", "-42", "58"},
		{"negative_add_positive", "-100", "42", "-58"},
		{"decimal_add_decimal", "42.5", "7.5", "50"},
		{"mixed_exponents", "42.5", "7.05", "49.55"},
		{"very_small_numbers", "0.000001", "0.000002", "0.000003"},
		{"large_valid_sum", "999999999999999990", "9", "999999999999999999"},
		{"overflow", "9999999999999999999", "1", "Infinity"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Parse(tt.a)
			b := Parse(tt.b)
			sum := a.Add(&b)
			if got := (&sum).String(); got != tt.want {
				t.Errorf("Add(%s, %s) = %s, want %s", tt.a, tt.b, got, tt.want)
			}
		})
	}

	t.Run("nil_handling", func(t *testing.T) {
		var a *FixedPoint = nil
		b := Parse("42")
		sum := a.Add(&b)
		if sum.String() != "NaN" {
			t.Errorf("Expected NaN when adding with nil, got %s", sum.String())
		}

		a = &b
		var c *FixedPoint = nil
		sum = a.Add(c)
		if sum.String() != "NaN" {
			t.Errorf("Expected NaN when adding with nil, got %s", sum.String())
		}
	})
}

func TestFixedPoint_Sub(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want string
	}{
		{"zero_sub_zero", "0", "0", "0"},
		{"positive_sub_positive", "100", "42", "58"},
		{"negative_sub_negative", "-100", "-42", "-58"},
		{"positive_sub_negative", "42", "-58", "100"},
		{"negative_sub_positive", "-42", "58", "-100"},
		{"decimal_sub_decimal", "50", "7.5", "42.5"},
		{"mixed_exponents", "49.55", "7.05", "42.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Parse(tt.a)
			b := Parse(tt.b)
			diff := a.Sub(&b)
			if got := (&diff).String(); got != tt.want {
				t.Errorf("Sub(%s, %s) = %s, want %s", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestFixedPoint_Mul(t *testing.T) {
	tests := []struct {
		name         string
		a            string
		b            string
		want         string
		wantOverflow bool
	}{
		{"zero_mul_zero", "0", "0", "0", false},
		{"zero_mul_nonzero", "0", "42", "0", false},
		{"nonzero_mul_zero", "42", "0", "0", false},
		{"positive_mul_positive", "42", "10", "420", false},
		{"negative_mul_positive", "-42", "10", "-420", false},
		{"positive_mul_negative", "42", "-10", "-420", false},
		{"negative_mul_negative", "-42", "-10", "420", false},
		{"decimal_mul_decimal", "4.2", "0.1", "0.42", false},
		{"large_valid_product", "999999999", "999999999", "999999998000000001", false},
		{"overflow", "9999999999999999999", "2", "Infinity", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Parse(tt.a)
			b := Parse(tt.b)
			product := a.Mul(&b)
			if got := (&product).String(); got != tt.want {
				t.Errorf("Mul(%s, %s) = %s, want %s", tt.a, tt.b, got, tt.want)
			}
			if tt.wantOverflow && (&product).Signal()&SignalOverflow == 0 {
				t.Errorf("Mul(%s, %s) expected overflow signal", tt.a, tt.b)
			}
		})
	}
}

func TestFixedPoint_Div(t *testing.T) {
	tests := []struct {
		name       string
		a          string
		b          string
		want       string
		wantSignal Signal
	}{
		{"nonzero_div_nonzero", "42", "2", "21", 0},
		{"zero_div_nonzero", "0", "42", "0", 0},
		{"nonzero_div_zero", "42", "0", "Infinity", SignalDivisionByZero},
		{"positive_div_positive", "420", "10", "42", 0},
		{"negative_div_positive", "-420", "10", "-42", 0},
		{"positive_div_negative", "420", "-10", "-42", 0},
		{"negative_div_negative", "-420", "-10", "42", 0},
		{"decimal_div_integer", "4.2", "2", "2.1", 0},
		{"repeating_decimal", "1", "3", "0.3333333333333333333", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Parse(tt.a)
			b := Parse(tt.b)
			quotient := a.Div(&b)
			if got := (&quotient).String(); got != tt.want {
				t.Errorf("Div(%s, %s) = %s, want %s", tt.a, tt.b, got, tt.want)
			}
			if (&quotient).Signal()&tt.wantSignal != tt.wantSignal {
				t.Errorf("Div(%s, %s) signal = %v, want %v", tt.a, tt.b, (&quotient).Signal(), tt.wantSignal)
			}
		})
	}
}

func TestFixedPoint_Neg(t *testing.T) {
	tests := []struct {
		name string
		a    string
		want string
	}{
		{"neg_zero", "0", "-0"},
		{"neg_positive", "42", "-42"},
		{"neg_negative", "-42", "42"},
		{"neg_decimal", "4.2", "-4.2"},
		{"neg_neg_decimal", "-4.2", "4.2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Parse(tt.a)
			negated := a.Neg()
			if got := (&negated).String(); got != tt.want {
				t.Errorf("Neg(%s) = %s, want %s", tt.a, got, tt.want)
			}
		})
	}

	t.Run("nil_handling", func(t *testing.T) {
		var a *FixedPoint = nil
		negated := a.Neg()
		if negated.String() != "NaN" {
			t.Errorf("Expected NaN when negating nil, got %s", negated.String())
		}
	})
}

func TestFixedPoint_String(t *testing.T) {
	tests := []struct {
		name string
		fp   FixedPoint
		want string
	}{
		{"zero", New(0, 0), "0"},
		{"positive", New(42, 0), "42"},
		{"negative", New(-42, 0), "-42"},
		{"positive_exp", New(42, 2), "4200"},
		{"negative_exp", New(42, -2), "0.42"},
		{"very_negative_exp", New(42, -4), "0.0042"},
		{"special_nan", FixedPoint{flg: flags{nan: true}}, "NaN"},
		{"special_pos_inf", FixedPoint{flg: flags{inf: true}}, "Infinity"},
		{"special_neg_inf", FixedPoint{sign: true, flg: flags{inf: true}}, "-Infinity"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (&tt.fp).String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("nil_handling", func(t *testing.T) {
		var fp *FixedPoint = nil
		if got := fp.String(); got != "NaN" {
			t.Errorf("String() on nil = %v, want NaN", got)
		}
	})
}

func TestFixedPoint_Equal(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"zero_equals_zero", "0", "0", true},
		{"zero_equals_negative_zero", "0", "-0", true},
		{"same_value", "42", "42", true},
		{"different_value", "42", "43", false},
		{"different_sign", "42", "-42", false},
		{"same_value_different_exp", "42", "42.0", true},
		{"same_value_different_exp2", "4.2", "4.20", true},
		{"different_value_close", "4.2", "4.21", false},
		{"nan_comparison", "NaN", "42", false},
		{"nan_nan_comparison", "NaN", "NaN", false},
		{"inf_comparison", "Infinity", "42", false},
		{"inf_inf_comparison", "Infinity", "Infinity", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Parse(tt.a)
			b := Parse(tt.b)
			if got := a.Equal(&b); got != tt.want {
				t.Errorf("Equal(%s, %s) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}

	t.Run("nil_handling", func(t *testing.T) {
		var a *FixedPoint = nil
		b := Parse("42")
		if a.Equal(&b) {
			t.Error("Expected nil not to equal a value")
		}

		a = &b
		var c *FixedPoint = nil
		if a.Equal(c) {
			t.Error("Expected a value not to equal nil")
		}
	})
}

func TestFixedPoint_LessThan(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"zero_less_than_positive", "0", "42", true},
		{"positive_less_than_larger", "42", "43", true},
		{"negative_less_than_zero", "-42", "0", true},
		{"negative_less_than_positive", "-42", "42", true},
		{"negative_less_than_less_negative", "-42", "-41", true},
		{"equal_values", "42", "42", false},
		{"larger_not_less_than_smaller", "43", "42", false},
		{"positive_not_less_than_negative", "42", "-42", false},
		{"decimal_comparison", "41.9", "42", true},
		{"decimal_comparison_equal", "42.0", "42", false},
		{"neg_infinity_less_than_anything", "-Infinity", "-999999999999999999", true},
		{"anything_less_than_pos_infinity", "999999999999999999", "Infinity", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Parse(tt.a)
			b := Parse(tt.b)
			if got := a.LessThan(&b); got != tt.want {
				t.Errorf("LessThan(%s, %s) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}

	t.Run("nan_comparisons", func(t *testing.T) {
		nan := Parse("NaN")
		num := Parse("42")

		if nan.LessThan(&num) {
			t.Error("NaN should not be less than any number")
		}

		if num.LessThan(&nan) {
			t.Error("No number should be less than NaN")
		}
	})
}

func TestFixedPoint_Abs(t *testing.T) {
	tests := []struct {
		name string
		a    string
		want string
	}{
		{"abs_zero", "0", "0"},
		{"abs_positive", "42", "42"},
		{"abs_negative", "-42", "42"},
		{"abs_decimal", "4.2", "4.2"},
		{"abs_neg_decimal", "-4.2", "4.2"},
		{"abs_nan", "NaN", "NaN"},
		{"abs_pos_infinity", "Infinity", "Infinity"},
		{"abs_neg_infinity", "-Infinity", "Infinity"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Parse(tt.a)
			abs := a.Abs()
			if got := (&abs).String(); got != tt.want {
				t.Errorf("Abs(%s) = %s, want %s", tt.a, got, tt.want)
			}
		})
	}
}

func TestZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		zero := Zero
		if got := (&zero).String(); got != "0" {
			t.Errorf("Zero() = %s, want \"0\"", got)
		}
	})
}

func TestOne(t *testing.T) {
	t.Run("one", func(t *testing.T) {
		one := One
		if got := (&one).String(); got != "1" {
			t.Errorf("One() = %s, want \"1\"", got)
		}
	})
}

func TestFixedPoint_Signal(t *testing.T) {
	fp := Parse("42")
	if s := fp.Signal(); s != SignalClear { // assuming SignalClear is 0
		t.Errorf("Expected SignalClear (0), got %v", s)
	}
}

func TestFixedPoint_IsOk(t *testing.T) {
	valid := Parse("42")
	if !valid.IsOk() {
		t.Error("Expected IsOk() to be true for valid FixedPoint")
	}
	invalid := Parse("")
	if invalid.IsOk() {
		t.Error("Expected IsOk() to be false for invalid FixedPoint")
	}
}

func TestFixedPoint_Must(t *testing.T) {
	valid := Parse("42")
	if valid.Must().String() != "42" {
		t.Error("Must() did not return correct FixedPoint for valid input")
	}
	invalid := Parse("")
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected Must() to panic for invalid FixedPoint")
		}
	}()
	_ = invalid.Must()
}

func TestFixedPoint_Handle(t *testing.T) {
	invalid := Parse("")
	handled := invalid.Handle(func(s Signal) *FixedPoint {
		// Return a valid FixedPoint (42)
		return &FixedPoint{coe: 42, exp: 0, flg: flags{sig: 0}}
	})
	if handled.String() != "42" {
		t.Error("Handle did not override invalid FixedPoint as expected")
	}
}

func TestFixedPoint_Copy(t *testing.T) {
	fp := Parse("-123.45")
	cp := fp.Copy()
	if cp.String() != fp.String() {
		t.Error("Copy() did not produce an identical FixedPoint")
	}
}

func TestFixedPoint_Clone(t *testing.T) {
	fp := Parse("123.45")
	clone := fp.Clone()
	if clone == nil || clone.String() != fp.String() {
		t.Error("Clone() did not produce an identical FixedPoint")
	}
}

func TestFixedPoint_Debug(t *testing.T) {
	fp := Parse("42.5")
	dbg := fp.Debug()
	if !strings.Contains(dbg, "sign:") || !strings.Contains(dbg, "coe:") || !strings.Contains(dbg, "exp:") {
		t.Errorf("Debug() output unexpected: %s", dbg)
	}
}

func TestFixedPoint_IsSpecial(t *testing.T) {
	nan := Parse("NaN")
	inf := Parse("inf")
	if !nan.IsSpecial() {
		t.Error("Expected NaN to be special")
	}
	if !inf.IsSpecial() {
		t.Error("Expected Infinity to be special")
	}
	pos := Parse("42")
	if pos.IsSpecial() {
		t.Error("Expected non-special number to not be special")
	}
}

func TestFixedPoint_IsNaN(t *testing.T) {
	nan := Parse("NaN")
	if !nan.IsNaN() {
		t.Error("Expected IsNaN() to be true for NaN")
	}
	num := Parse("42")
	if num.IsNaN() {
		t.Error("Expected IsNaN() to be false for a valid number")
	}
}

func TestFixedPoint_IsInf(t *testing.T) {
	inf := Parse("inf")
	if !inf.IsInf() {
		t.Error("Expected IsInf() to be true for Infinity")
	}
	num := Parse("42")
	if num.IsInf() {
		t.Error("Expected IsInf() to be false for a valid number")
	}
}

func TestFixedPoint_IsZero(t *testing.T) {
	zero := Parse("0")
	if !zero.IsZero() {
		t.Error("Expected IsZero() to be true for zero")
	}
	nonzero := Parse("1")
	if nonzero.IsZero() {
		t.Error("Expected IsZero() to be false for non-zero")
	}
}

func TestFixedPoint_IsNegative(t *testing.T) {
	neg := Parse("-42")
	if !neg.IsNegative() {
		t.Error("Expected IsNegative() to be true for negative number")
	}
	pos := Parse("42")
	if pos.IsNegative() {
		t.Error("Expected IsNegative() to be false for positive number")
	}
}

func TestFixedPoint_IsPositive(t *testing.T) {
	pos := Parse("42")
	if !pos.IsPositive() {
		t.Error("Expected IsPositive() to be true for positive number")
	}
	neg := Parse("-42")
	if neg.IsPositive() {
		t.Error("Expected IsPositive() to be false for negative number")
	}
}
