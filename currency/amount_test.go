package currency

import (
	"testing"
)

// Define a dummy currency type for testing.
type testCurrency struct{}

func (tc testCurrency) GetCode() Code              { return "TEST" }
func (tc testCurrency) GetSymbol() Symbol          { return "T" }
func (tc testCurrency) GetMinorUnitFactor() Factor { return 2 }

// TestAdd verifies the Add method.
func TestAdd(t *testing.T) {
	// Initialize dummy fixed point values.
	val1 := DefaultContext.Parse("100")
	val2 := DefaultContext.Parse("200")
	expected := DefaultContext.Parse("300")

	a := Amount[testCurrency]{Value: val1, Currency: testCurrency{}}
	b := Amount[testCurrency]{Value: val2, Currency: testCurrency{}}
	result := a.Add(b)

	// Fix: pass &expected to Equal.
	if !DefaultContext.Equal(result.Value, expected) {
		t.Errorf("Add failed: expected %v, got %v", expected, result.Value)
	}
}

// TestSub verifies the Sub method.
func TestSub(t *testing.T) {
	val1 := DefaultContext.Parse("300")
	val2 := DefaultContext.Parse("100")
	expected := DefaultContext.Parse("200")

	a := Amount[testCurrency]{Value: val1, Currency: testCurrency{}}
	b := Amount[testCurrency]{Value: val2, Currency: testCurrency{}}
	result := a.Sub(b)

	// Fix: use fp.Equals for comparison.
	if !DefaultContext.Equal(result.Value, expected) {
		t.Errorf("Sub failed: expected %v, got %v", expected, result.Value)
	}
}

// BenchmarkAdd measures performance of Add.
func BenchmarkAdd(b *testing.B) {
	val1 := DefaultContext.Parse("1000")
	val2 := DefaultContext.Parse("2000")
	a := Amount[testCurrency]{Value: val1, Currency: testCurrency{}}
	amt := Amount[testCurrency]{Value: val2, Currency: testCurrency{}}

	// Fix: replace b.Loop() with b.ResetTimer() and a proper loop.

	for b.Loop() {
		_ = a.Add(amt)
	}
}

// Example benchmark for Mul.
func BenchmarkMul(b *testing.B) {
	val := DefaultContext.Parse("1000")
	a := Amount[testCurrency]{Value: val, Currency: testCurrency{}}
	factor := DefaultContext.Parse("2.0")

	for b.Loop() {
		_ = a.Mul(factor)
	}
}
