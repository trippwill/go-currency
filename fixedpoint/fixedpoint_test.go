package fixedpoint

import (
	"testing"
)

func TestX64PackUnpackRoundtrip(t *testing.T) {
	tests := []struct {
		kind kind
		sign sign
		exp  int16
		coe  uint64
	}{
		{kind_finite, sign_positive, 0, 1},
		{kind_finite, sign_negative, -10, 12345},
		{kind_infinity, sign_positive, 0, 0},
		{kind_quiet, sign_negative, 0, 42},
		{kind_signaling, sign_positive, 0, 99},
	}

	for _, tt := range tests {
		var x X64
		err := x.pack(tt.kind, tt.sign, tt.exp, tt.coe)
		if err != nil {
			t.Fatalf("pack failed: %v", err)
		}

		unpackedKind, unpackedSign, unpackedExp, unpackedCoe, err := x.unpack()
		if err != nil {
			t.Fatalf("unpack failed: %v", err)
		}

		if unpackedKind != tt.kind || unpackedSign != tt.sign || unpackedExp != tt.exp || unpackedCoe != tt.coe {
			t.Errorf("roundtrip mismatch: got (%v, %v, %v, %v), want (%v, %v, %v, %v)",
				unpackedKind, unpackedSign, unpackedExp, unpackedCoe, tt.kind, tt.sign, tt.exp, tt.coe)
		}
	}
}

func TestX32PackUnpackRoundtrip(t *testing.T) {
	var x X32
	tests := []struct {
		kind kind
		sign sign
		exp  int8
		coe  uint32
	}{
		{kind_finite, sign_positive, 0, 1},
		{kind_finite, sign_negative, -5, 12345},
		{kind_infinity, sign_positive, 0, 0},
		{kind_quiet, sign_negative, 0, 42},
		{kind_finite, sign_negative, -95, 12345},
	}

	for _, test := range tests {
		err := x.pack(test.kind, test.sign, test.exp, test.coe)
		if err != nil {
			t.Fatalf("pack failed: %v", err)
		}

		kind, sign, exp, coe, err := x.unpack()
		if err != nil {
			t.Fatalf("unpack failed: %v", err)
		}

		if kind != test.kind || sign != test.sign || exp != test.exp || coe != test.coe {
			t.Errorf("roundtrip mismatch: got (%v, %v, %v, %v), want (%v, %v, %v, %v)",
				kind, sign, exp, coe, test.kind, test.sign, test.exp, test.coe)
		}
	}
}

func FuzzX64PackUnpackRoundtrip(f *testing.F) {
	f.Add(uint8(kind_finite), int8(sign_positive), int16(0), uint64(1))
	f.Add(uint8(kind_finite), int8(sign_negative), int16(-10), uint64(12345))
	f.Add(uint8(kind_infinity), int8(sign_positive), int16(0), uint64(0))
	f.Add(uint8(kind_quiet), int8(sign_negative), int16(0), uint64(42))
	f.Add(uint8(kind_signaling), int8(sign_positive), int16(0), uint64(99))

	f.Fuzz(func(t *testing.T, _kind uint8, _sign int8, exp int16, coe uint64) {
		switch _kind {
		case uint8(kind_finite), uint8(kind_infinity), uint8(kind_quiet), uint8(kind_signaling):
			// valid kinds
		default:
			t.Skipf("invalid kind: %v", _kind)
		}

		if _sign != int8(sign_positive) && _sign != int8(sign_negative) {
			t.Skipf("invalid sign: %v", _sign)
		}
		var x X64
		err := x.pack(kind(_kind), sign(_sign), exp, coe)
		if err != nil {
			t.Skipf("pack failed: %v", err)
		}

		unpackedKind, unpackedSign, unpackedExp, unpackedCoe, err := x.unpack()
		if err != nil {
			t.Fatalf("unpack failed: %v", err)
		}

		if unpackedKind == kind_signaling || unpackedKind == kind_quiet {
			exp = 0
		}

		if unpackedKind == kind_infinity {
			coe = 0
			exp = 0
		}

		if unpackedKind != kind(_kind) || unpackedSign != sign(_sign) || unpackedExp != exp || unpackedCoe != coe {
			t.Errorf("roundtrip mismatch: got (%v, %v, %v, %v), want (%v, %v, %v, %v)",
				unpackedKind, unpackedSign, unpackedExp, unpackedCoe, kind(_kind), sign(_sign), exp, coe)
		}
	})
}

func FuzzX32PackUnpackRoundtrip(f *testing.F) {
	f.Add(uint8(kind_finite), int8(sign_positive), int8(0), uint32(1))
	f.Add(uint8(kind_finite), int8(sign_negative), int8(-10), uint32(12345))
	f.Add(uint8(kind_infinity), int8(sign_positive), int8(0), uint32(0))
	f.Add(uint8(kind_quiet), int8(sign_negative), int8(0), uint32(42))
	f.Add(uint8(kind_signaling), int8(sign_positive), int8(0), uint32(99))

	f.Fuzz(func(t *testing.T, _kind uint8, _sign int8, exp int8, coe uint32) {
		switch _kind {
		case uint8(kind_finite), uint8(kind_infinity), uint8(kind_quiet), uint8(kind_signaling):
			// valid kinds
		default:
			t.Skipf("invalid kind: %v", _kind)
		}
		if _sign != int8(sign_positive) && _sign != int8(sign_negative) {
			t.Skipf("invalid sign: %v", _sign)
		}
		var x X32
		err := x.pack(kind(_kind), sign(_sign), exp, coe)
		if err != nil {
			t.Skipf("pack failed: %v", err)
		}

		unpackedKind, unpackedSign, unpackedExp, unpackedCoe, err := x.unpack()
		if err != nil {
			t.Fatalf("unpack failed: %v", err)
		}

		if unpackedKind == kind_signaling || unpackedKind == kind_quiet {
			exp = 0
		}

		if unpackedKind == kind_infinity {
			coe = 0
			exp = 0
		}

		if unpackedKind != kind(_kind) || unpackedSign != sign(_sign) || unpackedExp != exp || unpackedCoe != coe {
			t.Errorf("roundtrip mismatch: got (%v, %v, %v, %v), want (%v, %v, %v, %v)",
				unpackedKind, unpackedSign, unpackedExp, unpackedCoe, kind(_kind), sign(_sign), exp, coe)
		}
	})
}
