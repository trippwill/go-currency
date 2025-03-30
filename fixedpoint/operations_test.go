package fixedpoint

import (
	"testing"
)

// Use a default context for testing.
// It must match the context struct defined in the package.
var defaultCtx = defaultContext

func TestFiniteAdd(t *testing.T) {
	// Test adding two finite numbers with same exponent.
	// Represent 100 as coefficient 100 and exponent 0, and 200 similarly.
	a := new(FiniteNumber).Init(false, 100, 0, defaultCtx)
	b := new(FiniteNumber).Init(false, 200, 0, defaultCtx)
	expected := new(FiniteNumber).Init(false, 300, 0, defaultCtx)

	res := a.Add(b)
	result, ok := res.(*FiniteNumber)
	if !ok {
		t.Errorf("Expected FiniteNumber result, got %T", res)
		return
	}
	if !Equals(result, expected) {
		t.Errorf("Finite add failed: expected %+v, got %+v", expected, result)
	}
}

func TestFiniteAddDiffExponent(t *testing.T) {
	// Test adding two finite numbers with different exponents.
	// Represent 12.3 as (123, exp=1) and 0.77 as (77, exp=2).
	a := new(FiniteNumber).Init(false, 123, -1, defaultCtx)
	b := new(FiniteNumber).Init(false, 77, -2, defaultCtx)
	// Expected: Align exponents to the lower precision.
	expected := new(FiniteNumber).Init(false, 1307, -2, defaultCtx)

	res := a.Add(b)
	result, ok := res.(*FiniteNumber)
	if !ok {
		t.Errorf("Expected FiniteNumber result, got %T", res)
		return
	}

	if !Equals(result, expected) {
		t.Errorf("Finite add with different exponents failed: expected %+v, got %+v", expected, result)
	}
}

func TestInfinityAddFinite(t *testing.T) {
	// Test adding Infinity and a finite number.
	inf := new(Infinity).Init(false, defaultCtx) // positive infinity
	finite := new(FiniteNumber).Init(false, 500, 0, defaultCtx)
	// According to rules, Infinity + finite = Infinity when signs match.
	res := inf.Add(finite)
	resultInf, ok := res.(*Infinity)
	if !ok {
		t.Errorf("Expected Infinity result, got %T", res)
		return
	}
	if resultInf.sign != inf.sign {
		t.Errorf("Infinity add failed: expected sign %v, got %v", inf.sign, resultInf.sign)
	}

	// Test reverse: finite + Infinity should yield same result.
	res2 := finite.Add(inf)
	result2, ok := res2.(*Infinity)
	if !ok {
		t.Errorf("Expected Infinity result, got %T", res2)
		return
	}
	if result2.sign != inf.sign {
		t.Errorf("Finite + Infinity add failed: expected sign %v, got %v", inf.sign, result2.sign)
	}
}

func TestInfinityAddInfinity(t *testing.T) {
	// Test adding two infinities with same sign (should yield NaN).
	inf1 := new(Infinity).Init(false, defaultCtx)
	inf2 := new(Infinity).Init(false, defaultCtx)

	res := inf1.Add(inf2)
	if !res.IsNaN() {
		t.Errorf("Expected NaN for Infinity + Infinity with same sign, got %T", res)
	}

	// Test adding infinities with opposite signs.
	infPos := new(Infinity).Init(false, defaultCtx)
	infNeg := new(Infinity).Init(true, defaultCtx)

	res2 := infPos.Add(infNeg)
	if res2Finite, ok := res2.(*FiniteNumber); ok {
		// Infinity + (-Infinity) should result in Zero (represented as a finite zero).
		if res2Finite.coe != 0 {
			t.Errorf("Expected Zero for Infinity + -Infinity, got %+v", res2Finite)
		}
	} else {
		t.Errorf("Expected Zero (FiniteNumber) for Infinity + -Infinity, got %T", res2)
	}
}

func TestFiniteAddInfinityDifferentSign(t *testing.T) {
	// Finite + Infinity when signs differ should yield NaN.
	finite := new(FiniteNumber).Init(false, 250, 0, defaultCtx)
	inf := new(Infinity).Init(true, defaultCtx) // negative infinity

	res := finite.Add(inf)
	if !res.IsNaN() {
		t.Errorf("Expected NaN for finite + Infinity with different signs, got %T", res)
	}
}

func TestAddWithNaN(t *testing.T) {
	// When one of the operands is NaN, the result should be a clone of that NaN.
	nan := new(NaN).Init(SignalInvalidOperation, 1)
	finite := new(FiniteNumber).Init(false, 100, 0, defaultCtx)

	res := nan.Add(finite)
	if !res.IsNaN() {
		t.Errorf("Expected NaN result when first operand is NaN, got %T", res)
	}

	res2 := finite.Add(nan)
	if !res2.IsNaN() {
		t.Errorf("Expected NaN result when second operand is NaN, got %T", res2)
	}
}
