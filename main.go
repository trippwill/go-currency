package main

import (
	"fmt"
	"unsafe"

	fp "github.com/trippwill/go-currency/fixedpoint"
)

func main() {
	println("FiniteNumber:", unsafe.Sizeof(fp.FiniteNumber{}))
	println("Infinity:", unsafe.Sizeof(fp.Infinity{}))
	println("NaN:", unsafe.Sizeof(fp.NaN{}))
	println("--------------------")

	format := "%-5s\t%12s\t%s\n"
	sep := "-------------------------------------"

	a := fp.Parse("100.00")
	a.SetContext(6, fp.RoundHalfUp)
	b := fp.Parse("200.00")
	c := fp.Must(a.Add(b))

	fmt.Printf(format, "a", a.String(), a.Debug())
	fmt.Printf(format, "b", b.String(), b.Debug())
	fmt.Printf(format, "a+b", c.String(), c.Debug())
	println(sep)

	a = fp.Parse("-0.50")
	b = fp.Parse("37.50")
	c = fp.Must(a.Add(b))
	d := fp.Must(a.Sub(b))

	fmt.Printf(format, "a", a.String(), a.Debug())
	fmt.Printf(format, "b", b.String(), b.Debug())
	fmt.Printf(format, "a+b", c.String(), c.Debug())
	fmt.Printf(format, "a-b", d.String(), d.Debug())
	println(sep)

	a = fp.Parse("0.1")
	c = fp.Must(a.Mul(a))

	fmt.Printf(format, "a", a.String(), a.Debug())
	fmt.Printf(format, "a*a", c.String(), c.Debug())
	println(sep)

	a = &fp.One
	b = fp.Parse("3")
	c = a.Div(b)

	fmt.Printf(format, "a", a.String(), a.Debug())
	fmt.Printf(format, "b", b.String(), b.Debug())
	fmt.Printf(format, "a/b", c.String(), c.Debug())
	println(sep)

	a = fp.Parse("1")
	d = fp.Must(a.Add(c))

	fmt.Printf(format, "a", a.String(), a.Debug())
	fmt.Printf(format, "c", c.String(), c.Debug())
	fmt.Printf(format, "a+c", d.String(), d.Debug())
	println(sep)

	b = fp.Parse("4")
	c = a.Div(b)

	fmt.Printf(format, "a", a.String(), a.Debug())
	fmt.Printf(format, "b", b.String(), b.Debug())
	fmt.Printf(format, "a/b", c.String(), c.Debug())
	println(sep)

	x := fp.Parse("10000000000000000000")
	x.SetContext(4, fp.RoundHalfUp)
	y := fp.Parse("-0.00000012345")
	z := x.Add(y)

	fmt.Printf(format, "x", x.String(), x.Debug())
	fmt.Printf(format, "y", y.String(), y.Debug())
	fmt.Printf(format, "x+y", z.String(), z.Debug())
	println(sep)

	// Demonstrate FixedPoint creation and parsing
	a = fp.Parse("123.45")
	b = fp.Parse("-67.89")
	fmt.Println("a:", a.String(), "b:", b.String())

	// Demonstrate checks
	fmt.Println("a is finite:", a.IsFinite())
	fmt.Println("b is negative:", b.IsNegative())

	// Demonstrate special values
	inf := fp.Parse("Infinity")
	ninf := fp.Parse("-Infinity")
	nan := fp.Parse("NaN")
	fmt.Println("Infinity:", inf.String(), "NaN:", nan.String(), "-Infinity:", ninf.String())
	fmt.Println("Infinity is infinite:", inf.IsInf())
	fmt.Println("-Infinity is infinite:", ninf.IsInf())
	fmt.Println("-Infinity is negative:", ninf.IsNegative())
	fmt.Println("NaN is NaN:", nan.IsNaN())
}
