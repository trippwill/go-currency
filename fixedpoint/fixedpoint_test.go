package fixedpoint

import (
	"testing"
)

func TestParseFinite(t *testing.T) {
	fp := Parse("123.45")
	// Expect a FiniteNumber for valid finite input.
	_, ok := fp.(*FiniteNumber)
	if !ok {
		t.Fatalf("expected FiniteNumber, got %T", fp)
	}
	// Check the value and exponent.
	if fp.String() != "123.45" {
		t.Errorf("expected 123.45, got %s", fp.String())
	}
	if fp.(*FiniteNumber).exp != -2 {
		t.Errorf("expected exponent -2, got %d", fp.(*FiniteNumber).exp)
	}
	// Check the sign.
	if fp.(*FiniteNumber).sign {
		t.Error("expected positive sign, got negative")
	}
	// Check the coefficient.
	if fp.(*FiniteNumber).coe != 12345 {
		t.Errorf("expected coefficient 12345, got %d", fp.(*FiniteNumber).coe)
	}
	// Check the context.
	if fp.(*FiniteNumber).context.signal != defaultContext.signal {
		t.Errorf("expected signal valid, got %s", fp.(*FiniteNumber).context.signal)
	}
	// Check the context's precision.
	if fp.(*FiniteNumber).context.precision != defaultContext.precision {
		t.Errorf("expected precision 2, got %d", fp.(*FiniteNumber).context.precision)
	}
	// Check the context's rounding mode.
	if fp.(*FiniteNumber).context.rounding != defaultContext.rounding {
		t.Errorf("expected rounding HalfEven, got %s", fp.(*FiniteNumber).context.rounding)
	}
}

func TestParseInvalid(t *testing.T) {
	fp := Parse("invalid")
	// Expect a NaN for non-numeric input.
	_, ok := fp.(*NaN)
	if !ok {
		t.Fatalf("expected NaN, got %T", fp)
	}
	// Check the signal.
	if fp.(*NaN).context.signal != SignalConversionSyntax {
		t.Errorf("expected signal SignalConversionSyntax, got %s", fp.(*NaN).context.signal)
	}
	// Check the diagnostic.
	if fp.(*NaN).diag == 0 {
		t.Error("expected non-zero diagnostic, got 0")
	}
	// Check the sign.
	if fp.(*NaN).sign {
		t.Error("expected positive sign, got negative")
	}
}

func TestEqualityAndComparison(t *testing.T) {
	a := Parse("123.45")
	b := Parse("123.45")
	if !Equals(a, b) {
		t.Error("expected a == b")
	}
	c := Parse("100.00")
	if !LessThan(c, a) {
		t.Error("expected 100.00 < 123.45")
	}
	if !GreaterThan(a, c) {
		t.Error("expected 123.45 > 100.00")
	}
	if !LessThanOrEqual(a, b) {
		t.Error("expected a <= b")
	}
	if !GreaterThanOrEqual(a, b) {
		t.Error("expected a >= b")
	}
}

func TestMust(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected Must(nil) to panic")
		}
	}()
	// Must should panic when provided a nil FixedPoint value.
	Must(nil)
}

func TestClone(t *testing.T) {
	a := Parse("123.45")
	clone := a.Clone()
	if a == clone {
		t.Error("Clone should return a different instance")
	}
	// Additional tests can compare Debug or String outputs if available.
}
