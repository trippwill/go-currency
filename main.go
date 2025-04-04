package main

import (
	"fmt"

	fp "github.com/trippwill/go-currency/fixedpoint"
)

func main() {
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
