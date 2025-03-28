package main

import (
	"unsafe"

	fp "github.com/trippwill/go-currency/fixedpoint"
)

func main() {
	println(unsafe.Sizeof(fp.FixedPoint{}))
	println("--------------------")

	a := fp.New(10000, -2)
	b := fp.New(20000000, -4)
	c := fp.Must(a.Add(&b))

	println(a.String(), "\t", a.Debug())
	println(b.String(), "\t", b.Debug())
	println(c.String(), "\t", c.Debug())
	println("--------------------")

	a = fp.New(-50, -2)
	b = fp.New(3750, -2)
	c = fp.Must(a.Add(&b))
	d := a.Sub(&b)

	println(a.String(), "\t", a.Debug())
	println(b.String(), "\t", b.Debug())
	println(c.String(), "\t", c.Debug())
	println(d.String(), "\t", d.Debug())
	println("--------------------")

	a = fp.Parse("0.1")
	c = fp.Must(a.Mul(&a))

	println(a.String(), "\t", a.Debug())
	println(c.String(), "\t", c.Debug())

	a = fp.One
	b = fp.Parse("3")
	c = fp.Must(a.Div(&b))

	println(a.String(), "\t", a.Debug())
	println(b.String(), "\t", b.Debug())
	println(c.String(), "\t", c.Debug())
	println("--------------------")

	a = fp.New(1, 0)
	d = a.Add(c)

	println(a.String(), "\t", a.Debug())
	println(c.String(), "\t", c.Debug())
	println(d.String(), "\t", d.Debug())
	println("--------------------")

	b = fp.Parse("4")
	c = fp.Must(a.Div(&b))

	println(a.String(), "\t", a.Debug())
	println(b.String(), "\t", b.Debug())
	println(c.String(), "\t", c.Debug())
	println("--------------------")
}
