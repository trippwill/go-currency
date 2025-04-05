package fixedpoint

import (
	"testing"
)

func TestParse(t *testing.T) {
	ctx := BasicContext()

	tests := []struct {
		input      string
		expectSign sign
		expectErr  bool
		expectKind kind
		expectExp  int16
		expectCoe  uint64
	}{
		{"123.45", sign_positive, false, kind_finite, -2, 12345},
		{"-123.45", sign_negative, false, kind_finite, -2, 12345},
		{"", sign_positive, true, kind_signaling, 0, 0},
		{"abc", sign_positive, true, kind_signaling, 0, 0},
		{"123..45", sign_positive, true, kind_signaling, 0, 0},
		{"NaN", sign_positive, false, kind_quiet, 0, 0},
		{"Infinity", sign_positive, false, kind_infinity, 0, 0},
		{"-Infinity", sign_negative, false, kind_infinity, 0, 0},
	}

	for _, tt := range tests {
		ctx.signals = 0 // Clear signals before each test
		result := ctx.Parse(tt.input)

		if tt.expectErr && ctx.signals&SignalConversionSyntax == 0 {
			t.Errorf("expected conversion syntax error for input %q", tt.input)
		}
		if !tt.expectErr && ctx.signals&SignalConversionSyntax != 0 {
			t.Errorf("unexpected conversion syntax error for input %q", tt.input)
		}
		kind, sign, exp, coe, err := result.unpack()
		if err != nil {
			t.Errorf("unexpected error unpacking result for input %q: %v", tt.input, err)
		}
		if !tt.expectErr && kind != tt.expectKind && sign != tt.expectSign && exp != tt.expectExp && coe != tt.expectCoe {
			t.Errorf("unexpected unpack for input %q: got (%v, %v, %v, %v), want (%v, %v, %v, %v)",
				tt.input, kind, sign, exp, coe, tt.expectKind, tt.expectSign, tt.expectExp, tt.expectCoe)
		}
	}
}
