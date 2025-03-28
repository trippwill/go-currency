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
	val1 := fp.New(100, 0)
	val2 := fp.New(200, 0)
	expected := fp.New(300, 0)

	a := Amount[testCurrency]{Value: &val1, Currency: testCurrency{}}
	b := Amount[testCurrency]{Value: &val2, Currency: testCurrency{}}
	result := a.Add(b)

	// Fix: pass &expected to Equal.
	if !result.Value.Equal(&expected) {
		t.Errorf("Add failed: expected %v, got %v", expected, result.Value)
	}
}

// TestSub verifies the Sub method.
func TestSub(t *testing.T) {
	val1 := fp.New(300, 0)
	val2 := fp.New(100, 0)
	expected := fp.New(200, 0)

	a := Amount[testCurrency]{Value: &val1, Currency: testCurrency{}}
	b := Amount[testCurrency]{Value: &val2, Currency: testCurrency{}}
	result := a.Sub(b)

	// Fix: pass &expected to Equal.
	if !result.Value.Equal(&expected) {
		t.Errorf("Sub failed: expected %v, got %v", expected, result.Value)
	}
}

// BenchmarkAdd measures performance of Add.
func BenchmarkAdd(b *testing.B) {
	val1 := fp.New(1000, 0)
	val2 := fp.New(2000, 0)
	a := Amount[testCurrency]{Value: &val1, Currency: testCurrency{}}
	amt := Amount[testCurrency]{Value: &val2, Currency: testCurrency{}}

	for b.Loop() {
		_ = a.Add(amt)
	}
}

// Example benchmark for Mul.
func BenchmarkMul(b *testing.B) {
	val := fp.New(1000, 0)
	a := Amount[testCurrency]{Value: &val, Currency: testCurrency{}}
	factor := fp.New(2, 0)

	for b.Loop() {
		_ = a.Mul(factor)
	}
}
