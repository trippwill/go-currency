package fixedpoint

import (
	"testing"
)

func BenchmarkPack(b *testing.B) {
	var x X64
	kind := kind_finite
	sign := signc_positive
	exp := int16(102)
	coe := uint64(0x7FFFFFFFFFFFF)

	for b.Loop() {
		err := x.pack(kind, sign, exp, coe)
		if err != nil {
			b.Fatalf("pack failed: %v", err)
		}
	}
}

func BenchmarkUnpack(b *testing.B) {
	x := X64{0x7FFFFFFFFFFFF} // Example packed value

	for b.Loop() {
		_, _, _, _, err := x.unpack()
		if err != nil {
			b.Fatalf("unpack failed: %v", err)
		}
	}
}

func BenchmarkQuantize(b *testing.B) {
	var x X64
	if err := x.pack(kind_finite, signc_positive, 0, 123456789012345); err != nil {
		b.Fatalf("pack failed: %v", err)
	}

	for b.Loop() {
		_, err := quantize64(x, 0, RoundTiesToEven)
		if err != Signal(0) {
			b.Fatalf("quantize failed: %v", err)
		}
	}
}

func BenchmarkString(b *testing.B) {
	var x X64
	if err := x.pack(kind_finite, signc_positive, 0, 0); err != nil {
		b.Fatalf("pack failed: %v", err)
	}

	for b.Loop() {
		_ = x.String()
	}
}
