package fixedpoint_test

import (
	"testing"

	"github.com/trippwill/go-currency/fixedpoint"
)

func TestFiniteNumber_Add(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected string
	}{
		{"Add positive numbers", "1.23", "4.56", "5.79000000"},
		{"Add negative numbers", "-1.23", "-4.56", "-5.79000000"},
		{"Add positive and negative", "1.23", "-4.56", "-3.33000000"},
		{"Add zero", "1.23", "0", "1.23000000"},
	}

	ctx := fixedpoint.BasicContext()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ctx.Parse(tt.a)
			b := ctx.Parse(tt.b)
			result := a.Add(b, ctx)
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestFiniteNumber_Sub(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected string
	}{
		{"Subtract positive numbers", "5.00", "2.50", "2.50000000"},
		{"Subtract negative numbers", "-5.00", "-2.50", "-2.50000000"},
		{"Subtract positive and negative", "5.00", "-2.50", "7.50000000"},
		{"Subtract zero", "5.00", "0", "5.00000000"},
	}

	ctx := fixedpoint.BasicContext()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ctx.Parse(tt.a)
			b := ctx.Parse(tt.b)
			result := a.Sub(b, ctx)
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestFiniteNumber_Mul(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected string
	}{
		{"Multiply positive numbers", "2.00", "3.50", "7.00000000"},
		{"Multiply negative numbers", "-2.00", "-3.50", "7.00000000"},
		{"Multiply positive and negative", "2.00", "-3.50", "-7.00000000"},
		{"Multiply by zero", "2.00", "0", "0."},
	}

	ctx := fixedpoint.BasicContext()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ctx.Parse(tt.a)
			b := ctx.Parse(tt.b)
			result := a.Mul(b, ctx)
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestFiniteNumber_Div(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected string
	}{
		{"Divide positive numbers", "7.00", "2.00", "3.50000000"},
		{"Divide negative numbers", "-7.00", "-2.00", "3.50000000"},
		{"Divide positive and negative", "7.00", "-2.00", "-3.50000000"},
		{"Divide by one", "7.00", "1.00", "7.00000000"},
	}

	ctx := fixedpoint.BasicContext()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ctx.Parse(tt.a)
			b := ctx.Parse(tt.b)
			result := a.Div(b, ctx)
			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}
