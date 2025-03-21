package currency

import (
	"fmt"
	"math"

	"golang.org/x/text/language"
)

type Amount[C Currency] struct {
	Value    FixedPoint
	Currency C
}

func NewAmount[C Currency](value FixedPoint) Amount[C] {
	var c C
	return Amount[C]{
		Value:    value,
		Currency: c,
	}
}

func NewAmountFromString[C Currency](value string, po *ParseOpts) Amount[C] {
	var c C
	return Amount[C]{
		Value:    NewFixedPoint(value, po),
		Currency: c,
	}
}

func AmountAdd[C Currency](a, b string, po *ParseOpts) Amount[C] {
	amountA := NewAmountFromString[C](a, po)
	amountB := NewAmountFromString[C](b, po)
	return amountA.Add(amountB)
}

func AmountSub[C Currency](a, b string, po *ParseOpts) Amount[C] {
	amountA := NewAmountFromString[C](a, po)
	amountB := NewAmountFromString[C](b, po)
	return amountA.Sub(amountB)
}

func AmountMul[C Currency](a, factor string, po *ParseOpts) Amount[C] {
	amountA := NewAmountFromString[C](a, po)
	f := NewFixedPoint(factor, po)
	return amountA.Mul(f)
}

func AmountDiv[C Currency](a, divisor string, po *ParseOpts) Amount[C] {
	amountA := NewAmountFromString[C](a, po)
	f := NewFixedPoint(divisor, po)
	return amountA.Div(f)
}

func (a Amount[C]) Format(tag language.Tag) string {
	return fmt.Sprintf("%s %s", a.Currency.GetSymbol(), a.Value.Format(tag))
}

func (a Amount[C]) String() string {
	return a.Format(language.Tag{})
}

func (a Amount[C]) Add(b Amount[C]) Amount[C] {
	// Normalize the precision
	maxPrecision := max(b.Value.Scale, a.Value.Scale)

	scaleA := int64(math.Pow10(int(maxPrecision - a.Value.Scale)))
	scaleB := int64(math.Pow10(int(maxPrecision - b.Value.Scale)))

	return Amount[C]{
		Value: FixedPoint{
			Base:  a.Value.Base*scaleA + b.Value.Base*scaleB,
			Scale: maxPrecision,
		},
		Currency: a.Currency,
	}
}

func (a Amount[C]) Sub(b Amount[C]) Amount[C] {
	// Normalize the precision
	maxPrecision := max(b.Value.Scale, a.Value.Scale)

	scaleA := int64(math.Pow10(int(maxPrecision - a.Value.Scale)))
	scaleB := int64(math.Pow10(int(maxPrecision - b.Value.Scale)))

	return Amount[C]{
		Value: FixedPoint{
			Base:  a.Value.Base*scaleA - b.Value.Base*scaleB,
			Scale: maxPrecision,
		},
		Currency: a.Currency,
	}
}

func (a Amount[C]) Mul(factor FixedPoint) Amount[C] {
	newBase := a.Value.Base * factor.Base
	newPrecision := a.Value.Scale + factor.Scale

	return Amount[C]{
		Value: FixedPoint{
			Base:  newBase,
			Scale: newPrecision,
		},
		Currency: a.Currency,
	}
}

func (a Amount[C]) Div(divisor FixedPoint) Amount[C] {
	if divisor.Base == 0 {
		panic("division by zero")
	}

	newBase := a.Value.Base * int64(math.Pow10(int(divisor.Scale))) / divisor.Base
	newPrecision := a.Value.Scale

	return Amount[C]{
		Value: FixedPoint{
			Base:  newBase,
			Scale: newPrecision,
		},
		Currency: a.Currency,
	}
}

// Neg returns the negation of the amount.
func (a Amount[C]) Neg() Amount[C] {
	return Amount[C]{
		Value: FixedPoint{
			Base:  -a.Value.Base,
			Scale: a.Value.Scale,
		},
		Currency: a.Currency,
	}
}

// Abs returns the absolute value of the amount.
func (a Amount[C]) Abs() Amount[C] {
	base := a.Value.Base
	if base < 0 {
		base = -base
	}
	return Amount[C]{
		Value: FixedPoint{
			Base:  base,
			Scale: a.Value.Scale,
		},
		Currency: a.Currency,
	}
}

// IsZero returns true if the amount is zero.
func (a Amount[C]) IsZero() bool {
	return a.Value.Base == 0
}

// Equal compares two amounts for equality after normalizing precision.
// It returns false if the currencies differ.
func (a Amount[C]) Equal(b Amount[C]) bool {
	maxPrecision := max(a.Value.Scale, b.Value.Scale)
	scaleA := int64(math.Pow10(int(maxPrecision - a.Value.Scale)))
	scaleB := int64(math.Pow10(int(maxPrecision - b.Value.Scale)))
	return a.Value.Base*scaleA == b.Value.Base*scaleB
}

func ConvIert[C, D Currency](a Amount[C], factor FixedPoint) Amount[D] {
	// Convert the amount using the conversion factor
	converted := a.Mul(factor)

	// Initialize the new currency type D
	var d D

	// Build the converted Amount with currency type D
	return Amount[D]{
		Value:    converted.Value,
		Currency: d,
	}
}
