package main

import (
	"fmt"

	fp "github.com/trippwill/go-currency/fixedpoint"
)

func main() {
	cx := fp.BasicContext[fp.Context64]()

	q := cx.Parse("-1234567.45")
	fmt.Printf("%s\t%s\n", q.String(), q.Debug())
	fmt.Printf("Packed: %#064b\n", q)
	fmt.Printf("Packed: %#016X\n", q)
	fmt.Printf("Packed: %020d\n", q)
	println("--------------------")

	w := cx.Parse("-1.457845784578")
	fmt.Printf("%s\t%s\n", w.String(), w.Debug())
	fmt.Printf("Packed: %#064b\n", w)
	fmt.Printf("Packed: %#016X\n", w)
	fmt.Printf("Packed: %020d\n", w)
	println("--------------------")

	e := cx.Parse("-Infinity")
	fmt.Printf("%s\t%s\n", e.String(), e.Debug())
	fmt.Printf("Packed: %#064b\n", e)
	fmt.Printf("Packed: %#016X\n", e)
	fmt.Printf("Packed: %020d\n", e)
	println("--------------------")

	r := cx.Parse("NaN")
	fmt.Printf("%s\t%s\n", r.String(), r.Debug())
	fmt.Printf("Packed: %#064b\n", r)
	fmt.Printf("Packed: %#016X\n", r)
	fmt.Printf("Packed: %020d\n", r)
	println("--------------------")

	var a fp.X64
	err := a.Pack(3, -1, -1, 12345)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\t%s\n", a.String(), a.Debug())
	fmt.Printf("Packed: %#064b\n", a)
	fmt.Printf("Packed: %#016X\n", a)
	fmt.Printf("Packed: %020d\n", a)
	println("--------------------")

	err = a.Pack(2, 1, 12, 345)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\t%s\n", a.String(), a.Debug())
	fmt.Printf("Packed: %#064b\n", a)
	fmt.Printf("Packed: %#016X\n", a)
	fmt.Printf("Packed: %020d\n", a)
	println("--------------------")

	err = a.Pack(3, 1, -119, 98765432109875)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\t%s\n", a.String(), a.Debug())
	fmt.Printf("Packed: %#064b\n", a)
	fmt.Printf("Packed: %#016X\n", a)
	fmt.Printf("Packed: %020d\n", a)
	println("--------------------\n")
}
