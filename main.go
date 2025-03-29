package main

import (
	"fmt"
	"unsafe"

	fp "github.com/trippwill/go-currency/fixedpoint"
)

func main() {
	println(unsafe.Sizeof(fp.FiniteNumber{}))
	println(unsafe.Sizeof(fp.Infinity{}))
	println(unsafe.Sizeof(fp.NaN{}))
	println("--------------------")

	a := fp.Parse("100.00")
	b := fp.Parse("200.00")
	c := fp.Must(a.Add(b))

	println(a.String(), "\t", a.Debug())
	println(b.String(), "\t", b.Debug())
	println(c.String(), "\t", c.Debug())
	println("--------------------")

	a = fp.Parse("-0.50")
	b = fp.Parse("37.50")
	c = fp.Must(a.Add(b))
	d := fp.Must(a.Sub(b))

	println(a.String(), "\t", a.Debug())
	println(b.String(), "\t", b.Debug())
	println(c.String(), "\t", c.Debug())
	println(d.String(), "\t", d.Debug())
	println("--------------------")

	a = fp.Parse("0.1")
	c = fp.Must(a.Mul(a))

	println(a.String(), "\t", a.Debug())
	println(c.String(), "\t", c.Debug())

	a = &fp.One
	b = fp.Parse("3")
	c = a.Div(b)

	println(a.String(), "\t", a.Debug())
	println(b.String(), "\t", b.Debug())
	println(c.String(), "\t", c.Debug())
	println("--------------------")

	a = fp.Parse("1")
	d = fp.Must(a.Add(c))

	println(a.String(), "\t", a.Debug())
	println(c.String(), "\t", c.Debug())
	println(d.String(), "\t", d.Debug())
	println("--------------------")

	b = fp.Parse("4")
	c = a.Div(b)

	println(a.String(), "\t", a.Debug())
	println(b.String(), "\t", b.Debug())
	println(c.String(), "\t", c.Debug())
	println("--------------------")

	x := fp.Parse("10000000000000000000")
	x.SetContext(4, fp.RoundHalfUp)
	y := fp.Parse("-0.00000012345")
	z := x.Add(y)

	println(x.String(), "\t", x.Debug())
	println(y.String(), "\t", y.Debug())
	println(z.String(), "\t", z.Debug())
	println("--------------------")

	// Demonstrate FixedPoint creation and parsing
	a = fp.Parse("123.45")
	b = fp.Parse("-67.89")
	fmt.Println("a:", a.String(), "b:", b.String())

	// Demonstrate addition
	sum := fp.Must(a.Add(b))
	fmt.Println("Sum:", sum.String())

	// Demonstrate subtraction
	diff := fp.Must(a.Sub(b))
	fmt.Println("Difference:", diff.String())

	// Demonstrate multiplication
	product := fp.Must(a.Mul(b))
	fmt.Println("Product:", product.String())

	// Demonstrate division
	quotient := a.Div(b)
	fmt.Println("Quotient:", quotient.String())

	// Demonstrate checks
	fmt.Println("a is finite:", a.IsFinite())
	fmt.Println("b is negative:", b.IsNegative())
	fmt.Println("Sum is positive:", sum.IsPositive())

	// Demonstrate special values
	inf := fp.Parse("Infinity")
	ninf := fp.Parse("-Infinity")
	nan := fp.Parse("NaN")
	fmt.Println("Infinity:", inf.String(), "NaN:", nan.String())
	fmt.Println("Infinity:", inf.String(), "Negative Infinity:", ninf.String())
	fmt.Println("Infinity is infinite:", inf.IsInf())
	fmt.Println("NaN is NaN:", nan.IsNaN())
}
