package currency

import (
	"testing"

	fp "github.com/trippwill/go-currency/fixedpoint"
)

// Define a dummy currency type for testing.
type testCurrency struct{}

func (tc testCurrency) GetCode() Code              { return "TEST" }
func (tc testCurrency) GetSymbol() Symbol          { return "T" }
func (tc testCurrency) GetMinorUnitFactor() Factor { return 2 }

// TestAdd verifies the Add method.
func TestAdd(t *testing.T) {
	// Initialize dummy fixed point values.
	val1 := fp.Parse("100")
	val2 := fp.Parse("200")
	expected := fp.Parse("300")

	a := Amount[testCurrency]{Value: val1, Currency: testCurrency{}}
	b := Amount[testCurrency]{Value: val2, Currency: testCurrency{}}
	result := a.Add(b)

	// Fix: pass &expected to Equal.
	if !fp.Equals(result.Value, expected) {
		t.Errorf("Add failed: expected %v, got %v", expected, result.Value)
	}
}

// TestSub verifies the Sub method.
func TestSub(t *testing.T) {
	val1 := fp.Parse("300")
	val2 := fp.Parse("100")
	expected := fp.Parse("200")

	a := Amount[testCurrency]{Value: val1, Currency: testCurrency{}}
	b := Amount[testCurrency]{Value: val2, Currency: testCurrency{}}
	result := a.Sub(b)

	// Fix: use fp.Equals for comparison.
	if !fp.Equals(result.Value, expected) {
		t.Errorf("Sub failed: expected %v, got %v", expected, result.Value)
	}
}

// BenchmarkAdd measures performance of Add.
func BenchmarkAdd(b *testing.B) {
	val1 := fp.Parse("1000")
	val2 := fp.Parse("2000")
	a := Amount[testCurrency]{Value: val1, Currency: testCurrency{}}
	amt := Amount[testCurrency]{Value: val2, Currency: testCurrency{}}

	// Fix: replace b.Loop() with b.ResetTimer() and a proper loop.
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.Add(amt)
	}
}

// Example benchmark for Mul.
func BenchmarkMul(b *testing.B) {
	val := fp.Parse("1000")
	a := Amount[testCurrency]{Value: val, Currency: testCurrency{}}
	factor := fp.Parse("2.0")

	for b.Loop() {
		_ = a.Mul(factor)
	}
}
