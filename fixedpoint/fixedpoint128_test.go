package fixedpoint

import (
	"math/big"
	"testing"
)

func TestFixedPoint128_Sign(t *testing.T) {
	tests := []struct {
		name     string
		fp       FixedPoint128
		expected bool
	}{
		{
			name:     "Positive",
			fp:       FixedPoint128{hi: 0, lo: 123},
			expected: false,
		},
		{
			name:     "Negative",
			fp:       FixedPoint128{hi: 1 << 63, lo: 123},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fp.sign()
			if got != tt.expected {
				t.Errorf("sign() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFixedPoint128_SetSign(t *testing.T) {
	tests := []struct {
		name     string
		fp       FixedPoint128
		sign     bool
		expected FixedPoint128
	}{
		{
			name:     "Set Positive",
			fp:       FixedPoint128{hi: 1 << 63, lo: 123},
			sign:     false,
			expected: FixedPoint128{hi: 0, lo: 123},
		},
		{
			name:     "Set Negative",
			fp:       FixedPoint128{hi: 0, lo: 123},
			sign:     true,
			expected: FixedPoint128{hi: 1 << 63, lo: 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := tt.fp
			fp.setSign(tt.sign)
			if fp.hi != tt.expected.hi || fp.lo != tt.expected.lo {
				t.Errorf("setSign() got = {%v, %v}, want {%v, %v}",
					fp.hi, fp.lo, tt.expected.hi, tt.expected.lo)
			}
		})
	}
}

func TestFixedPoint128_Coefficient(t *testing.T) {
	tests := []struct {
		name     string
		fp       FixedPoint128
		expected *big.Int
	}{
		{
			name:     "Zero",
			fp:       FixedPoint128{hi: 0, lo: 0},
			expected: big.NewInt(0),
		},
		{
			name:     "Low bits only",
			fp:       FixedPoint128{hi: 0, lo: 0xFFFFFFFFFFFFFFFF},
			expected: new(big.Int).SetUint64(0xFFFFFFFFFFFFFFFF),
		},
		{
			name:     "High bits only",
			fp:       FixedPoint128{hi: fp128_coe_mask_hi, lo: 0},
			expected: new(big.Int).Lsh(new(big.Int).SetUint64(fp128_coe_mask_hi), 64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fp.coefficient()
			if got.Cmp(tt.expected) != 0 {
				t.Errorf("coefficient() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFixedPoint128_SetCoefficient(t *testing.T) {
	tests := []struct {
		name     string
		fp       FixedPoint128
		coef     *big.Int
		expected FixedPoint128
		sig      SIG
	}{
		{
			name:     "Zero",
			fp:       FixedPoint128{hi: 0, lo: 0},
			coef:     big.NewInt(0),
			expected: FixedPoint128{hi: 0, lo: 0},
			sig:      SIG_NONE,
		},
		{
			name:     "Small value",
			fp:       FixedPoint128{hi: 0, lo: 0},
			coef:     big.NewInt(123456),
			expected: FixedPoint128{hi: 0, lo: 123456},
			sig:      SIG_NONE,
		},
		{
			name:     "Negative",
			fp:       FixedPoint128{hi: 0, lo: 0},
			coef:     big.NewInt(-1),
			expected: FixedPoint128{hi: 0, lo: 0},
			sig:      SIG_INVALID_OPERATION,
		},
		{
			name:     "Overflow",
			fp:       FixedPoint128{hi: 0, lo: 0},
			coef:     new(big.Int).Lsh(big.NewInt(1), 114),
			expected: FixedPoint128{hi: 0, lo: 0},
			sig:      SIG_OVERFLOW,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := tt.fp
			sig := fp.setCoefficient(tt.coef)
			if sig != tt.sig {
				t.Errorf("setCoefficient() signal = %v, want %v", sig, tt.sig)
			}
			if sig == SIG_NONE && (fp.hi != tt.expected.hi || fp.lo != tt.expected.lo) {
				t.Errorf("setCoefficient() result = {%v, %v}, want {%v, %v}",
					fp.hi, fp.lo, tt.expected.hi, tt.expected.lo)
			}
		})
	}
}

func TestFixedPoint128_Exponent(t *testing.T) {
	tests := []struct {
		name     string
		fp       FixedPoint128
		expected int
	}{
		{
			name:     "Zero",
			fp:       FixedPoint128{hi: 0, lo: 0},
			expected: -fp128_exp_bias,
		},
		{
			name:     "One",
			fp:       FixedPoint128{hi: uint64(fp128_exp_bias+1) << 49, lo: 0},
			expected: 1,
		},
		{
			name:     "Negative",
			fp:       FixedPoint128{hi: uint64(fp128_exp_bias-1) << 49, lo: 0},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fp.exponent()
			if got != tt.expected {
				t.Errorf("exponent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFixedPoint128_SetExponent(t *testing.T) {
	tests := []struct {
		name     string
		fp       FixedPoint128
		exp      int
		expected FixedPoint128
		sig      SIG
	}{
		{
			name:     "Zero",
			fp:       FixedPoint128{hi: 0, lo: 0},
			exp:      0,
			expected: FixedPoint128{hi: uint64(fp128_exp_bias) << 49, lo: 0},
			sig:      SIG_NONE,
		},
		{
			name:     "Invalid too large",
			fp:       FixedPoint128{hi: 0, lo: 0},
			exp:      9000,
			expected: FixedPoint128{hi: 0, lo: 0},
			sig:      SIG_INVALID_OPERATION,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := tt.fp
			sig := fp.setExponent(tt.exp)
			if sig != tt.sig {
				t.Errorf("setExponent() signal = %v, want %v", sig, tt.sig)
			}
			if sig == SIG_NONE && (fp.hi != tt.expected.hi || fp.lo != tt.expected.lo) {
				t.Errorf("setExponent() result = {%v, %v}, want {%v, %v}",
					fp.hi, fp.lo, tt.expected.hi, tt.expected.lo)
			}
		})
	}
}

func TestFixedPoint128_SpecialValues(t *testing.T) {
	t.Run("NaN", func(t *testing.T) {
		var fp FixedPoint128
		fp.setNaN(false)
		if !fp.isNaN() {
			t.Errorf("Expected isNaN() to be true after setNaN()")
		}
		if fp.isInf() {
			t.Errorf("Expected isInf() to be false for NaN")
		}
	})

	t.Run("SNaN", func(t *testing.T) {
		var fp FixedPoint128
		fp.setSNaN(false)
		if !fp.isSNaN() {
			t.Errorf("Expected isSNaN() to be true after setSNaN()")
		}
		if !fp.isNaN() {
			t.Errorf("Expected isNaN() to also be true for SNaN")
		}
	})

	t.Run("Infinity", func(t *testing.T) {
		var fp FixedPoint128
		fp.setInf(false)
		if !fp.isInf() {
			t.Errorf("Expected isInf() to be true after setInf()")
		}
		if fp.isNaN() {
			t.Errorf("Expected isNaN() to be false for Infinity")
		}
	})

	t.Run("Sign preservation", func(t *testing.T) {
		// Test negative NaN
		var fp FixedPoint128
		fp.setNaN(true)
		if !fp.sign() {
			t.Errorf("Expected sign to be preserved as negative for NaN")
		}

		// Test negative Infinity
		fp.setInf(true)
		if !fp.sign() {
			t.Errorf("Expected sign to be preserved as negative for Infinity")
		}
	})
}

func TestFixedPoint128_IsFinite(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*FixedPoint128)
		expected bool
	}{
		{
			name:     "Zero is finite",
			setup:    func(fp *FixedPoint128) {},
			expected: true,
		},
		{
			name: "Normal number is finite",
			setup: func(fp *FixedPoint128) {
				fp.setExponent(10)
				fp.setCoefficient(big.NewInt(123))
			},
			expected: true,
		},
		{
			name:     "Infinity is not finite",
			setup:    func(fp *FixedPoint128) { fp.setInf(false) },
			expected: false,
		},
		{
			name:     "NaN is not finite",
			setup:    func(fp *FixedPoint128) { fp.setNaN(false) },
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fp FixedPoint128
			tt.setup(&fp)
			got := fp.isFinite()
			if got != tt.expected {
				t.Errorf("isFinite() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParse128(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantFP  func(fp FixedPoint128) bool // predicate to validate fp
		wantErr error
	}{
		{
			name:  "NaN",
			input: "NaN",
			wantFP: func(fp FixedPoint128) bool {
				return fp.isNaN() && !fp.sign() // NaN with positive sign
			},
			wantErr: nil,
		},
		{
			name:  "-NaN",
			input: "-NaN",
			wantFP: func(fp FixedPoint128) bool {
				return fp.isSNaN() && fp.sign() // SNaN with negative sign
			},
			wantErr: nil,
		},
		{
			name:  "sNaN",
			input: "sNaN",
			wantFP: func(fp FixedPoint128) bool {
				return fp.isSNaN() && fp.sign() // SNaN with negative sign (per implementation)
			},
			wantErr: nil,
		},
		{
			name:  "Infinity",
			input: "Infinity",
			wantFP: func(fp FixedPoint128) bool {
				return fp.isInf() && !fp.isNaN() && !fp.sign()
			},
			wantErr: nil,
		},
		{
			name:  "+Infinity",
			input: "+Infinity",
			wantFP: func(fp FixedPoint128) bool {
				return fp.isInf() && !fp.isNaN() && !fp.sign()
			},
			wantErr: nil,
		},
		{
			name:  "-Infinity",
			input: "-Infinity",
			wantFP: func(fp FixedPoint128) bool {
				return fp.isInf() && !fp.isNaN() && fp.sign()
			},
			wantErr: nil,
		},
		{
			name:  "Positive number 1",
			input: "1",
			wantFP: func(fp FixedPoint128) bool {
				// For 1, we expect exponent = 0 and coefficient = 1, positive sign.
				return !fp.sign() &&
					fp.exponent() == 0 &&
					fp.coefficient().Cmp(big.NewInt(1)) == 0
			},
			wantErr: nil,
		},
		{
			name:  "Negative number -1",
			input: "-1",
			wantFP: func(fp FixedPoint128) bool {
				// For -1, coefficient = 1, exponent = 0, negative sign preserved.
				return fp.sign() &&
					fp.exponent() == 0 &&
					fp.coefficient().Cmp(big.NewInt(1)) == 0
			},
			wantErr: nil,
		},
		{
			name:  "Decimal number 123.45",
			input: "123.45",
			wantFP: func(fp FixedPoint128) bool {
				// Conversion of 123.45:
				// rat = 12345/100 => loop multiplies numerator twice: num becomes 1234500 then trailing zero removal divides twice => coefficient = 12345, exponent = 0.
				return !fp.sign() &&
					fp.exponent() == 0 &&
					fp.coefficient().Cmp(big.NewInt(12345)) == 0
			},
			wantErr: nil,
		},
		{
			name:    "Invalid syntax",
			input:   "abc",
			wantFP:  nil,
			wantErr: ErrConversionSyntax,
		},
		{
			name:    "Non-integral decimal expansion",
			input:   "1/2",
			wantFP:  nil,
			wantErr: ErrConversionSyntax,
		},
		{
			name: "Overflow number",
			// A number that results in a coefficient with bit length > 113.
			// For example, a 36-digit number is roughly > 2^113.
			input:   "100000000000000000000000000000000000",
			wantFP:  nil,
			wantErr: ErrOverflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp, err := Parse128(tt.input)
			if tt.wantErr != nil {
				if err == nil || err != tt.wantErr {
					t.Errorf("Parse128(%q) error = %v, want error %v", tt.input, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Parse128(%q) unexpected error: %v", tt.input, err)
				return
			}
			if tt.wantFP != nil && !tt.wantFP(fp) {
				t.Errorf("Parse128(%q) produced unexpected FixedPoint128 representation", tt.input)
			}
		})
	}
}

// Fuzz test for coefficient setting and retrieval
func FuzzFixedPoint128_Coefficient(f *testing.F) {
	seeds := []uint64{0, 1, 123456, 0xFFFFFFFFFFFFFFFF}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, val uint64) {
		var fp FixedPoint128
		coef := new(big.Int).SetUint64(val)
		sig := fp.setCoefficient(coef)

		if sig != SIG_NONE {
			return // Skip values that would cause overflow
		}

		got := fp.coefficient()
		if got.Cmp(coef) != 0 {
			t.Errorf("Coefficient mismatch: set %v but got %v", coef, got)
		}
	})
}

// Fuzz test for exponent setting and retrieval
func FuzzFixedPoint128_Exponent(f *testing.F) {
	// Add some seed values within the valid exponent range
	seeds := []int{-6176, -1000, 0, 1000, 8000}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, exp int) {
		var fp FixedPoint128
		sig := fp.setExponent(exp)

		if sig != SIG_NONE {
			return // Skip values outside the valid range
		}

		got := fp.exponent()
		if got != exp {
			t.Errorf("Exponent mismatch: set %v but got %v", exp, got)
		}
	})
}
