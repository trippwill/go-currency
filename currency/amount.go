package currency

import (
	fp "github.com/trippwill/go-currency/fixedpoint"
	"golang.org/x/text/language"
)

var DefaultContext = fp.BasicContext()

type Amount[C Currency] struct {
	Value    fp.FixedPoint
	Currency C
}

func (a Amount[C]) Format(tag language.Tag) string {
	panic("not implemented")
}

func (a Amount[C]) String() string {
	return a.Format(language.Tag{})
}

func (a Amount[C]) Add(b Amount[C]) Amount[C] {
	sum := a.Value.Add(b.Value, DefaultContext)
	return Amount[C]{
		Value:    sum,
		Currency: a.Currency,
	}
}

func (a Amount[C]) Sub(b Amount[C]) Amount[C] {
	diff := a.Value.Sub(b.Value, DefaultContext)
	return Amount[C]{
		Value:    diff,
		Currency: a.Currency,
	}
}

func (a Amount[C]) Mul(factor fp.FixedPoint) Amount[C] {
	product := a.Value.Mul(factor, DefaultContext)
	return Amount[C]{
		Value:    product,
		Currency: a.Currency,
	}
}

func (a Amount[C]) Div(divisor fp.FixedPoint) (Amount[C], error) {
	quotient := a.Value.Div(divisor, DefaultContext)
	return Amount[C]{
		Value:    quotient,
		Currency: a.Currency,
	}, nil
}

func (a Amount[C]) Neg() Amount[C] {
	negated := a.Value.Neg(DefaultContext)
	return Amount[C]{
		Value:    negated,
		Currency: a.Currency,
	}
}

func (a Amount[C]) Abs() Amount[C] {
	absVal := a.Value.Abs(DefaultContext)
	return Amount[C]{
		Value:    absVal,
		Currency: a.Currency,
	}
}

func (a Amount[C]) IsZero() bool {
	return a.Value.IsZero()
}

func (a Amount[C]) Equal(b Amount[C]) bool {
	return DefaultContext.Equal(a.Value, b.Value)
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
