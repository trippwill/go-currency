package currency

import (
	"fmt"
	"testing"
)

// DummyCurrency is a simple implementation of the Currency interface for testing.
type DummyCurrency struct{}

func (d DummyCurrency) GetCode() Code              { return "XXX" }
func (d DummyCurrency) GetSymbol() Symbol          { return "$" }
func (d DummyCurrency) GetMinorUnitFactor() Factor { return 1 }

// Assume ParseOpts is defined in the main package.
// For testing purposes, we define a default instance.
var defaultParseOpts = &ParseOpts{
	thousands: ',',
	decimal:   '.',
}

func TestAmountAdd(t *testing.T) {
	tests := []struct {
		name     string
		a        Amount[USD]
		b        Amount[USD]
		expected Amount[USD]
	}{
		{
			name: "Same precision",
			a: Amount[USD]{
				Value:    FixedPoint{Base: 1000, Scale: 2},
				Currency: USD{},
			},
			b: Amount[USD]{
				Value:    FixedPoint{Base: 2000, Scale: 2},
				Currency: USD{},
			},
			expected: Amount[USD]{
				Value:    FixedPoint{Base: 3000, Scale: 2},
				Currency: USD{},
			},
		},
		{
			name: "Different precision",
			a: Amount[USD]{
				Value:    FixedPoint{Base: 1000, Scale: 2},
				Currency: USD{},
			},
			b: Amount[USD]{
				Value:    FixedPoint{Base: 200, Scale: 1},
				Currency: USD{},
			},
			expected: Amount[USD]{
				Value:    FixedPoint{Base: 3000, Scale: 2},
				Currency: USD{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Add(tt.b)
			if result.Value.Base != tt.expected.Value.Base || result.Value.Scale != tt.expected.Value.Scale {
				t.Errorf("Add() = %v, want %v", result.Value, tt.expected.Value)
			}
		})
	}
}

func TestAmountSub(t *testing.T) {
	tests := []struct {
		name     string
		a        Amount[USD]
		b        Amount[USD]
		expected Amount[USD]
	}{
		{
			name: "Same precision",
			a: Amount[USD]{
				Value:    FixedPoint{Base: 3000, Scale: 2},
				Currency: USD{},
			},
			b: Amount[USD]{
				Value:    FixedPoint{Base: 1000, Scale: 2},
				Currency: USD{},
			},
			expected: Amount[USD]{
				Value:    FixedPoint{Base: 2000, Scale: 2},
				Currency: USD{},
			},
		},
		{
			name: "Different precision",
			a: Amount[USD]{
				Value:    FixedPoint{Base: 3000, Scale: 2},
				Currency: USD{},
			},
			b: Amount[USD]{
				Value:    FixedPoint{Base: 200, Scale: 1},
				Currency: USD{},
			},
			expected: Amount[USD]{
				Value:    FixedPoint{Base: 1000, Scale: 2},
				Currency: USD{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Sub(tt.b)
			if result.Value.Base != tt.expected.Value.Base || result.Value.Scale != tt.expected.Value.Scale {
				t.Errorf("Sub() = %v, want %v", result.Value, tt.expected.Value)
			}
		})
	}
}

func TestAmountMul(t *testing.T) {
	tests := []struct {
		name     string
		a        Amount[USD]
		f        FixedPoint
		expected Amount[USD]
	}{
		{
			name: "Simple multiplication",
			a: Amount[USD]{
				Value:    FixedPoint{Base: 1000, Scale: 2},
				Currency: USD{},
			},
			f: FixedPoint{Base: 2, Scale: 0},
			expected: Amount[USD]{
				Value:    FixedPoint{Base: 2000, Scale: 2},
				Currency: USD{},
			},
		},
		{
			name: "Multiplication with precision",
			a: Amount[USD]{
				Value:    FixedPoint{Base: 1000, Scale: 2},
				Currency: USD{},
			},
			f: FixedPoint{Base: 25, Scale: 1},
			expected: Amount[USD]{
				Value:    FixedPoint{Base: 25000, Scale: 3},
				Currency: USD{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Mul(tt.f)
			if result.Value.Base != tt.expected.Value.Base || result.Value.Scale != tt.expected.Value.Scale {
				t.Errorf("Mul() = %v, want %v", result.Value, tt.expected.Value)
			}
		})
	}
}

func TestAmountDiv(t *testing.T) {
	tests := []struct {
		name     string
		a        Amount[USD]
		f        FixedPoint
		expected Amount[USD]
	}{
		{
			name: "Simple division",
			a: Amount[USD]{
				Value:    FixedPoint{Base: 1000, Scale: 2},
				Currency: USD{},
			},
			f: FixedPoint{Base: 2, Scale: 0},
			expected: Amount[USD]{
				Value:    FixedPoint{Base: 500, Scale: 2},
				Currency: USD{},
			},
		},
		{
			name: "Division with scale adjustment",
			a: Amount[USD]{
				Value:    FixedPoint{Base: 1000, Scale: 2},
				Currency: USD{},
			},
			f: FixedPoint{Base: 25, Scale: 1},
			expected: Amount[USD]{
				Value:    FixedPoint{Base: 400, Scale: 2},
				Currency: USD{},
			},
		},
		{
			name: "Division by larger scale",
			a: Amount[USD]{
				Value:    FixedPoint{Base: 1000, Scale: 2},
				Currency: USD{},
			},
			f: FixedPoint{Base: 1, Scale: 3},
			expected: Amount[USD]{
				Value:    FixedPoint{Base: 1000000, Scale: 2},
				Currency: USD{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Div(tt.f)
			if result.Value.Base != tt.expected.Value.Base || result.Value.Scale != tt.expected.Value.Scale {
				t.Errorf("Div() = %v, want %v", result.Value, tt.expected.Value)
			}
		})
	}
}

// Define a second dummy currency for conversion tests.
type OtherCurrency struct{}

func (o OtherCurrency) GetCode() Code              { return "YYY" }
func (o OtherCurrency) GetSymbol() Symbol          { return "€" }
func (o OtherCurrency) GetMinorUnitFactor() Factor { return 1 }

func TestConvIert(t *testing.T) {
	tests := []struct {
		name     string
		input    Amount[DummyCurrency]
		factor   FixedPoint
		expected FixedPoint
	}{
		{
			name: "Identity conversion",
			input: Amount[DummyCurrency]{
				Value:    FixedPoint{Base: 10000, Scale: 2}, // 100.00
				Currency: DummyCurrency{},
			},
			factor:   FixedPoint{Base: 1, Scale: 0}, // Multiply by 1
			expected: FixedPoint{Base: 10000, Scale: 2},
		},
		{
			name: "Double conversion",
			input: Amount[DummyCurrency]{
				Value:    FixedPoint{Base: 5000, Scale: 2}, // 50.00
				Currency: DummyCurrency{},
			},
			factor:   FixedPoint{Base: 2, Scale: 0}, // Multiply by 2
			expected: FixedPoint{Base: 10000, Scale: 2},
		},
		{
			name: "Conversion with precision change",
			input: Amount[DummyCurrency]{
				Value:    FixedPoint{Base: 2000, Scale: 2}, // 20.00
				Currency: DummyCurrency{},
			},
			factor:   FixedPoint{Base: 3, Scale: 1}, // Expected: 2000*3=6000, scale = 2+1 = 3
			expected: FixedPoint{Base: 6000, Scale: 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvIert[DummyCurrency, OtherCurrency](tc.input, tc.factor)
			if result.Value.Base != tc.expected.Base || result.Value.Scale != tc.expected.Scale {
				t.Errorf("ConvIert() Value = %v, want %v", result.Value, tc.expected)
			}
			// Verify the returned currency type by checking its symbol.
			if result.Currency.GetSymbol() != "€" {
				t.Errorf("ConvIert() Currency symbol = %v, want %v", result.Currency.GetSymbol(), "€")
			}
		})
	}
}
func ExampleAmountAdd() {
	// Example: 100.00 + 50.50 = 150.50
	result := AmountAdd[DummyCurrency]("100.00", "50.50", defaultParseOpts)
	fmt.Println(result.String())
	// Output: $ 150.50
}

func ExampleAmountSub() {
	// Example: 100.00 - 50.50 = 49.50
	result := AmountSub[DummyCurrency]("100.00", "50.50", defaultParseOpts)
	fmt.Println(result.String())
	// Output: $ 49.50
}

func ExampleAmountMul() {
	// Example: 100.00 * 2.50 = 250.00
	result := AmountMul[DummyCurrency]("100.00", "2.50", defaultParseOpts)
	fmt.Println(result.String())
	// Output: $ 250.0000
}

func ExampleAmountDiv() {
	// Example: 100.00 / 2.50 = 40.00
	result := AmountDiv[DummyCurrency]("100.00", "2.50", defaultParseOpts)
	fmt.Println(result.String())
	// Output: $ 40.00
}
