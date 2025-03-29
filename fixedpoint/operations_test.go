package fixedpoint

import (
	"testing"
)

var ctx = defaultContext

// For testing purposes, assume that FiniteNumber, Infinity, and NaN have the following minimal fields.
// (In real tests, these types are defined in the package and their Init methods initialize these fields.)

// Test FiniteNumber.Add
func TestFiniteNumberAdd(t *testing.T) {
	a := new(FiniteNumber).Init(false, 100, 0, ctx)
	b := new(FiniteNumber).Init(false, 23, 0, ctx)
	res := a.Add(b)

	fn, ok := res.(*FiniteNumber)
	if !ok {
		t.Fatalf("expected FiniteNumber, got %T", res)
	}
	if fn.coe != 123 {
		t.Errorf("expected coefficient 123, got %d", fn.coe)
	}
}

// Test FiniteNumber.Sub (by subtraction)
func TestFiniteNumberSub(t *testing.T) {
	a := new(FiniteNumber).Init(false, 150, 0, ctx)
	b := new(FiniteNumber).Init(false, 50, 0, ctx)
	res := a.Sub(b)

	fn, ok := res.(*FiniteNumber)
	if !ok {
		t.Fatalf("expected FiniteNumber, got %T", res)
	}
	if fn.coe != 100 {
		t.Errorf("expected coefficient 100, got %d", fn.coe)
	}
}

// Test FiniteNumber.Mul
func TestFiniteNumberMul(t *testing.T) {
	a := new(FiniteNumber).Init(false, 10, 0, ctx)
	b := new(FiniteNumber).Init(false, 20, 0, ctx)
	res := a.Mul(b)

	fn, ok := res.(*FiniteNumber)
	if !ok {
		t.Fatalf("expected FiniteNumber, got %T", res)
	}
	if fn.coe != 200 {
		t.Errorf("expected coefficient 200, got %d", fn.coe)
	}
}

// Test FiniteNumber.Div
func TestFiniteNumberDiv(t *testing.T) {
	a := new(FiniteNumber).Init(false, 200, 0, ctx)
	b := new(FiniteNumber).Init(false, 10, 0, ctx)
	res := a.Div(b)

	fn, ok := res.(*FiniteNumber)
	if !ok {
		t.Fatalf("expected FiniteNumber, got %T", res)
	}
	// Division may involve scaling – check for expected value (200/10 == 20).
	if fn.coe != 2 || fn.exp != 1 {
		t.Errorf("expected coefficient 2, exponent 1, got %d, %d", fn.coe, fn.exp)
	}
}

// Test FiniteNumber.Neg and Abs
func TestFiniteNumberNegAbs(t *testing.T) {
	a := new(FiniteNumber).Init(false, 250, 0, ctx)
	neg := a.Neg()

	fnNeg, ok := neg.(*FiniteNumber)
	if !ok {
		t.Fatalf("expected FiniteNumber from Neg(), got %T", neg)
	}
	// Check that negation flips the sign. (Assume 'true' means negative.)
	if fnNeg.sign != true {
		t.Errorf("expected negated sign true, got false")
	}

	abs := fnNeg.Abs()
	fnAbs, ok := abs.(*FiniteNumber)
	if !ok {
		t.Fatalf("expected FiniteNumber from Abs(), got %T", abs)
	}
	if fnAbs.sign != false {
		t.Errorf("expected absolute value sign false, got true")
	}
}

// Test FiniteNumber.Compare. For zero, comparison should be equal irrespective of sign.
func TestFiniteNumberCompare(t *testing.T) {
	a := new(FiniteNumber).Init(false, 0, 0, ctx)
	b := new(FiniteNumber).Init(true, 0, 0, ctx)
	if a.Compare(b) != 0 {
		t.Errorf("expected zero comparison, got non-zero")
	}

	// Compare two nonzero numbers.
	x := new(FiniteNumber).Init(false, 100, 0, ctx)
	y := new(FiniteNumber).Init(false, 200, 0, ctx)
	if x.Compare(y) != -1 {
		t.Errorf("expected x.Compare(y) == -1, got %d", x.Compare(y))
	}
	if y.Compare(x) != 1 {
		t.Errorf("expected y.Compare(x) == 1, got %d", y.Compare(x))
	}
}

// Test Infinity.Add with two Infinities of the same sign (invalid operation)
func TestInfinityAddSameSign(t *testing.T) {
	a := new(Infinity).Init(false, ctx)
	b := new(Infinity).Init(false, ctx)
	res := a.Add(b)

	_, ok := res.(*NaN)
	if !ok {
		t.Errorf("expected NaN as result of adding same-signed infinities")
	}
}

// Test Infinity.Add with opposite signs (should return Zero)
func TestInfinityAddOppositeSign(t *testing.T) {
	a := new(Infinity).Init(false, ctx)
	b := new(Infinity).Init(true, ctx)
	res := a.Add(b)

	// Assuming Zero is a predefined FiniteNumber representing 0.
	fn, ok := res.(*FiniteNumber)
	if !ok {
		t.Errorf("expected FiniteNumber (Zero), got %T", res)
	}
	if !fn.IsZero() {
		t.Errorf("expected Zero result")
	}
}
