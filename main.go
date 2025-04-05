package main

import (
	"fmt"

	fp "github.com/trippwill/go-currency/fixedpoint"
)

func main() {
	cx := fp.BasicContext()

	q := cx.Parse("-1234567.45")

	fmt.Printf("Packed: %#064b\n", q)
	fmt.Printf("Packed: %#016X\n", q)
	fmt.Printf("Packed: %020d\n", q)

	w := cx.Parse("-1.457845784578")

	fmt.Printf("Packed: %#064b\n", w)
	fmt.Printf("Packed: %#016X\n", w)
	fmt.Printf("Packed: %020d\n", w)

	e := cx.Parse("-Infinity")
	fmt.Printf("Packed: %#064b\n", e)
	fmt.Printf("Packed: %#016X\n", e)
	fmt.Printf("Packed: %020d\n", e)

	var a fp.X64
	err := a.Pack(3, -1, -1, 12345)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Packed: %#064b\n", a)
	fmt.Printf("Packed: %#016X\n", a)
	fmt.Printf("Packed: %020d\n", a)

	err = a.Pack(2, 1, 12, 345)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Packed: %#064b\n", a)
	fmt.Printf("Packed: %#016X\n", a)
	fmt.Printf("Packed: %020d\n", a)

	err = a.Pack(3, 1, -119, 98765432109875)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Packed: %#064b\n", a)
	fmt.Printf("Packed: %#016X\n", a)
	fmt.Printf("Packed: %020d\n", a)
}
