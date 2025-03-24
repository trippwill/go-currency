package currency

import (
	"golang.org/x/text/language"

	fp "github.com/trippwill/go-currency/fixedpoint"
)

type Amount[C Currency] struct {
	Value    fp.FixedPoint
	Currency C
}

func NewAmount[C Currency](value fp.FixedPoint) Amount[C] {
	var c C
	return Amount[C]{
		Value:    value,
		Currency: c,
	}
}

func NewAmountFromString[C Currency](value string, po *ParseOpts) (amt Amount[C], err error) {
	var c C
	fp, err := fp.NewFixedPoint(value)
	if err != nil {
		return
	}

	return Amount[C]{
		Value:    fp,
		Currency: c,
	}, nil
}

func (a Amount[C]) Format(tag language.Tag) string {
	panic("not implemented")
}

func (a Amount[C]) String() string {
	return a.Format(language.Tag{})
}

func (a Amount[C]) Add(b Amount[C]) Amount[C] {
	return Amount[C]{
		Value:    a.Value.Add(b.Value),
		Currency: a.Currency,
	}
}

func (a Amount[C]) Sub(b Amount[C]) Amount[C] {
	panic("not implemented")
}

func (a Amount[C]) Mul(factor fp.FixedPoint) Amount[C] {
	panic("not implemented")
}

func (a Amount[C]) Div(divisor fp.FixedPoint) (res Amount[C], err error) {
	panic("not implemented")
}

// Neg returns the negation of the amount.
func (a Amount[C]) Neg() Amount[C] {
	panic("not implemented")
}

// Abs returns the absolute value of the amount.
func (a Amount[C]) Abs() Amount[C] {
	panic("not implemented")
}

// IsZero returns true if the amount is zero.
func (a Amount[C]) IsZero() bool {
	panic("not implemented")
}

// Equal compares two amounts for equality after normalizing precision.
// It returns false if the currencies differ.
func (a Amount[C]) Equal(b Amount[C]) bool {
	return false
}

func Convert[C, D Currency](a Amount[C], factor fp.FixedPoint) Amount[D] {
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
