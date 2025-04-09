package fixedpoint

import (
	"testing"
)

// IEEE 754-2008 Conformance Tests
// These tests verify compliance with the IEEE 754-2008 standard for decimal floating-point arithmetic

// TestQuantizationRounding tests the quantization operation with different rounding modes.
// Quantization is an operation defined in IEEE 754-2008 that adjusts the exponent of a number
// while preserving its value as much as possible through rounding.
func TestQuantizationRounding(t *testing.T) {
	tests := []struct {
		name      string
		value     X64
		expTarget int16
		mode      Rounding
		expected  string
	}{
		{
			name: "Quantize-NoRoundingNeeded",
			value: func() X64 {
				var x X64
				_ = x.pack(kind_finite, signc_positive, 0, 12345)
				return x
			}(),
			expTarget: 0,
			mode:      RoundTiesToEven,
			expected:  "12345",
		},
		{
			name: "Quantize-RoundTiesToEven",
			value: func() X64 {
				var x X64
				_ = x.pack(kind_finite, signc_positive, -2, 12345)
				return x
			}(),
			expTarget: 0,
			mode:      RoundTiesToEven,
			expected:  "123",
		},
		{
			name: "Quantize-RoundTowardPositive",
			value: func() X64 {
				var x X64
				_ = x.pack(kind_finite, signc_positive, -1, 12345)
				return x
			}(),
			expTarget: 0,
			mode:      RoundTowardPositive,
			expected:  "1235",
		},
		{
			name: "Quantize-RoundTowardNegative",
			value: func() X64 {
				var x X64
				_ = x.pack(kind_finite, signc_negative, -1, 12345)
				return x
			}(),
			expTarget: 0,
			mode:      RoundTowardNegative,
			expected:  "-1235",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For this test, we need to implement a quantize function
			// that adjusts the exponent while preserving the value
			result, _ := quantize64(tt.value, tt.expTarget, tt.mode)
			got := result.String()
			if got != tt.expected {
				t.Errorf("quantize() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// Add tests for X32.pack and X32.unpack
func TestX32PackUnpackConformance(t *testing.T) {
	tests := []struct {
		name      string
		kind      kind
		sign      signc
		exp       int8
		coe       uint32
		expectErr bool
	}{
		{
			name: "ValidFinitePositive",
			kind: kind_finite,
			sign: signc_positive,
			exp:  0,
			coe:  12345,
		},
		{
			name: "ValidFiniteNegative",
			kind: kind_finite,
			sign: signc_negative,
			exp:  -5,
			coe:  67890,
		},
		{
			name: "Infinity",
			kind: kind_infinity,
			sign: signc_positive,
		},
		{
			name: "QuietNaN",
			kind: kind_quiet,
			sign: signc_positive,
		},
		{
			name: "SignalingNaN",
			kind: kind_signaling,
			sign: signc_negative,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var x X32
			err := x.pack(tt.kind, tt.sign, tt.exp, tt.coe)
			if (err != nil) != tt.expectErr {
				t.Errorf("pack() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr {
				k, s, e, c, err := x.unpack()
				if err != nil {
					t.Errorf("unpack() error = %v", err)
					return
				}
				if k != tt.kind || s != tt.sign || e != tt.exp || c != tt.coe {
					t.Errorf("unpack() = (%v, %v, %v, %v), want (%v, %v, %v, %v)", k, s, e, c, tt.kind, tt.sign, tt.exp, tt.coe)
				}
			}
		})
	}
}

// Add tests for X64.pack and X64.unpack
func TestX64PackUnpackConformance(t *testing.T) {
	tests := []struct {
		name      string
		kind      kind
		sign      signc
		exp       int16
		coe       uint64
		expectErr bool
	}{
		{
			name: "ValidFinitePositive",
			kind: kind_finite,
			sign: signc_positive,
			exp:  0,
			coe:  123456789012345,
		},
		{
			name: "ValidFiniteNegative",
			kind: kind_finite,
			sign: signc_negative,
			exp:  -10,
			coe:  987654321098765,
		},
		{
			name: "Infinity",
			kind: kind_infinity,
			sign: signc_positive,
		},
		{
			name: "QuietNaN",
			kind: kind_quiet,
			sign: signc_positive,
		},
		{
			name: "SignalingNaN",
			kind: kind_signaling,
			sign: signc_negative,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var x X64
			err := x.pack(tt.kind, tt.sign, tt.exp, tt.coe)
			if (err != nil) != tt.expectErr {
				t.Errorf("pack() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr {
				k, s, e, c, err := x.unpack()
				if err != nil {
					t.Errorf("unpack() error = %v", err)
					return
				}
				if k != tt.kind || s != tt.sign || e != tt.exp || c != tt.coe {
					t.Errorf("unpack() = (%v, %v, %v, %v), want (%v, %v, %v, %v)", k, s, e, c, tt.kind, tt.sign, tt.exp, tt.coe)
				}
			}
		})
	}
}
