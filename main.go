package main

import (
	"fmt"

	fp "github.com/trippwill/go-currency/fixedpoint"
)

func main() {
	a, _ := fp.Parse128("1.23")
	b, _ := fp.Parse128("4.56")

	fmt.Println("a: ", a)
	fmt.Println("b: ", b)

	fmt.Println("a: ", a.Scientific())
	fmt.Println("b: ", b.Scientific())

	fmt.Println("a: ", a.Debug())
	fmt.Println("b: ", b.Debug())
}
