package imath

import "testing"

func TestAbs(t *testing.T) {
	if Abs(-5) != 5 {
		t.Errorf("Abs(-5) = %d; want 5", Abs(-5))
	}
	if Abs(5) != 5 {
		t.Errorf("Abs(5) = %d; want 5", Abs(5))
	}
}

func TestPow(t *testing.T) {
	if Pow[int, uint](2, 3) != 8 {
		t.Errorf("Pow(2, 3) = %d; want 8", Pow[int, uint](2, 3))
	}
	if Pow[int, uint](5, 0) != 1 {
		t.Errorf("Pow(5, 0) = %d; want 1", Pow[int, uint](5, 0))
	}
}

func TestGCD(t *testing.T) {
	if GCD(48, 18) != 6 {
		t.Errorf("GCD(48, 18) = %d; want 6", GCD(48, 18))
	}
	if GCD(7, 1) != 1 {
		t.Errorf("GCD(7, 1) = %d; want 1", GCD(7, 1))
	}
}

func TestLCM(t *testing.T) {
	if LCM(4, 6) != 12 {
		t.Errorf("LCM(4, 6) = %d; want 12", LCM(4, 6))
	}
	if LCM(0, 5) != 0 {
		t.Errorf("LCM(0, 5) = %d; want 0", LCM(0, 5))
	}
}
