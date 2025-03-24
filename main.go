package main

import (
	"fmt"

	fp "github.com/trippwill/go-currency/fixedpoint"
)

func main() {
	// Example usage of FixedPoint
	fstr := "%-5s: %10s | %s\n"

	fp1, err := fp.NewFixedPoint("123.456")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fp2, err := fp.NewFixedPoint("-0.1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf(fstr, "FP 1", fp1, fp1.Debug())
	fmt.Printf(fstr, "FP 2", fp2, fp2.Debug())

	sum := fp1.Add(fp2)
	fmt.Printf(fstr, "Sum", sum, sum.Debug())

	fp3, err := fp.NewFixedPoint("NaN")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf(fstr, "NaN", fp3, fp3.Debug())
	fp4, err := fp.NewFixedPoint("inf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf(fstr, "Inf", fp4, fp4.Debug())
}
