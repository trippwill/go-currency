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

	bc := fp.BasicContext()
	println(bc.Debug())
	println(sep)

	a := bc.Parse("100.00")
	b := bc.Parse("200.00")
	c := bc.Must(bc.Add(a, b))

	fmt.Printf(format, "a", a.String(), a.Debug())
	fmt.Printf(format, "b", b.String(), b.Debug())
	fmt.Printf(format, "a+b", c.String(), c.Debug())
	println(bc.Debug())
	println(sep)

	a = bc.Parse("9999997.255")
	b = bc.Parse("9999997.255")
	c = bc.Must(bc.Add(a, b))
	d := bc.Must(bc.Sub(b, a))

	fmt.Printf(format, "a", a.String(), a.Debug())
	fmt.Printf(format, "b", b.String(), b.Debug())
	fmt.Printf(format, "a+b", c.String(), c.Debug())
	fmt.Printf(format, "a-b", d.String(), d.Debug())
	println(bc.Debug())
	println(sep)

	a = bc.Parse("-0.50")
	b = bc.Parse("37.50")
	c = bc.Must(bc.Add(a, b))
	d = bc.Must(bc.Sub(b, c))

	fmt.Printf(format, "a", a.String(), a.Debug())
	fmt.Printf(format, "b", b.String(), b.Debug())
	fmt.Printf(format, "a+b", c.String(), c.Debug())
	fmt.Printf(format, "a-b", d.String(), d.Debug())
	println(bc.Debug())
	println(sep)

	a = bc.Parse("0.1")
	c = bc.Must(fp.All(bc.Mul, a, a, a, a))

	fmt.Printf(format, "a", a.String(), a.Debug())
	fmt.Printf(format, "a*a*a*a", c.String(), c.Debug())
	println(bc.Debug())
	println(sep)

	a = bc.Parse("1.0000")
	b = bc.Parse("3")
	c = bc.Div(a, b)

	fmt.Printf(format, "a", a.String(), a.Debug())
	fmt.Printf(format, "b", b.String(), b.Debug())
	fmt.Printf(format, "a/b", c.String(), c.Debug())
	println(bc.Debug())
	println(sep)

	x := bc.Parse("10000000000000000000")
	y := bc.Parse("-0.00000012345")
	z := bc.Add(x, y)

	fmt.Printf(format, "x", x.String(), x.Debug())
	fmt.Printf(format, "y", y.String(), y.Debug())
	fmt.Printf(format, "x+y", z.String(), z.Debug())
	println(bc.Debug())
	println(sep)

	// Demonstrate checks
	fmt.Println("a is finite:", a.IsFinite())
	fmt.Println("b is negative:", b.IsNegative())

	bcc := bc.Clone(true)

	// Demonstrate special values
	inf := bcc.Parse("Infinity")
	ninf := bcc.Parse("-Infinity")
	nan := bcc.Parse("NaN")
	fmt.Println("Infinity:", inf.String(), "NaN:", nan.String(), "-Infinity:", ninf.String())
	fmt.Println("Infinity is infinite:", inf.IsInf())
	fmt.Println("-Infinity is infinite:", ninf.IsInf())
	fmt.Println("-Infinity is negative:", ninf.IsNegative())
	fmt.Println("NaN is NaN:", nan.IsNaN())
}
