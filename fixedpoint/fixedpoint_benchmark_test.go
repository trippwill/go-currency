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
