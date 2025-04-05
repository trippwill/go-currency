package fixedpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext64Parse(t *testing.T) {
	ctx := BasicContext[Context64]()

	tests := []struct {
		input      string
		expectSign signc
		expectErr  bool
		expectKind kind
		expectExp  int16
		expectCoe  uint64
	}{
		{"123.45", signc_positive, false, kind_finite, -2, 12345},
		{"-123.45", signc_negative, false, kind_finite, -2, 12345},
		{"", signc_positive, true, kind_signaling, 0, 0},
		{"abc", signc_positive, true, kind_signaling, 0, 0},
		{"123..45", signc_positive, true, kind_signaling, 0, 0},
		{"NaN", signc_positive, false, kind_quiet, 0, 0},
		{"Infinity", signc_positive, false, kind_infinity, 0, 0},
		{"-Infinity", signc_negative, false, kind_infinity, 0, 0},
	}

	for _, tt := range tests {
		ctx.signals = 0 // Clear signals before each test
		result := ctx.Parse(tt.input)

		if tt.expectErr {
			assert.NotZero(t, ctx.signals&SignalConversionSyntax, "expected conversion syntax error for input %q", tt.input)
		} else {
			assert.Zero(t, ctx.signals&SignalConversionSyntax, "unexpected conversion syntax error for input %q", tt.input)
		}

		kind, sign, exp, coe, err := result.unpack()
		assert.NoError(t, err, "unexpected error unpacking result for input %q", tt.input)

		if !tt.expectErr {
			assert.Equal(t, tt.expectKind, kind, "unexpected kind for input %q", tt.input)
			assert.Equal(t, tt.expectSign, sign, "unexpected sign for input %q", tt.input)
			assert.Equal(t, tt.expectExp, exp, "unexpected exponent for input %q", tt.input)
			assert.Equal(t, tt.expectCoe, coe, "unexpected coefficient for input %q", tt.input)
		}
	}
}

func TestContext32Parse(t *testing.T) {
	ctx := BasicContext[Context32]()

	tests := []struct {
		input      string
		expectSign signc
		expectErr  bool
		expectKind kind
		expectExp  int8
		expectCoe  uint32
	}{
		{"123.45", signc_positive, false, kind_finite, -2, 12345},
		{"-123.45", signc_negative, false, kind_finite, -2, 12345},
		{"", signc_positive, true, kind_signaling, 0, 0},
		{"abc", signc_positive, true, kind_signaling, 0, 0},
		{"123..45", signc_positive, true, kind_signaling, 0, 0},
		{"NaN", signc_positive, false, kind_quiet, 0, 0},
		{"Infinity", signc_positive, false, kind_infinity, 0, 0},
		{"-Infinity", signc_negative, false, kind_infinity, 0, 0},
	}

	for _, tt := range tests {
		ctx.signals = 0 // Clear signals before each test
		result := ctx.Parse(tt.input)

		if tt.expectErr {
			assert.NotZero(t, ctx.signals&SignalConversionSyntax, "expected conversion syntax error for input %q", tt.input)
		} else {
			assert.Zero(t, ctx.signals&SignalConversionSyntax, "unexpected conversion syntax error for input %q", tt.input)
		}

		kind, sign, exp, coe, err := result.unpack()
		assert.NoError(t, err, "unexpected error unpacking result for input %q", tt.input)

		if !tt.expectErr {
			assert.Equal(t, tt.expectKind, kind, "unexpected kind for input %q", tt.input)
			assert.Equal(t, tt.expectSign, sign, "unexpected sign for input %q", tt.input)
			assert.Equal(t, tt.expectExp, exp, "unexpected exponent for input %q", tt.input)
			assert.Equal(t, tt.expectCoe, coe, "unexpected coefficient for input %q", tt.input)
		}
	}
}
