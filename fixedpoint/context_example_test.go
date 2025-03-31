package fixedpoint_test

import (
	"fmt"

	"github.com/trippwill/go-currency/fixedpoint"
)

func ExampleBasicContext() {
	ctx := fixedpoint.BasicContext()
	fmt.Printf("%s", ctx)
	// Output:
	// context{precision: 9, rounding: RoundHalfUp,  traps: SignalInvalidOperation|SignalConversionSyntax}
}

func ExampleExtendedContext() {
	ctx := fixedpoint.ExtendedContext()
	fmt.Printf("%s", ctx)
	// Output:
	// context{precision: 15, rounding: RoundHalfEven,  traps: SignalOverflow|SignalInvalidOperation|SignalConversionSyntax}
}

func ExampleContext_Parse() {
	ctx := fixedpoint.BasicContext()
	fp := ctx.Parse("123.45")
	fmt.Println(fp)
	// Output:
	// 123.450000
}

func ExampleContext_Add() {
	ctx := fixedpoint.BasicContext()
	a := ctx.Parse("1.23")
	b := ctx.Parse("4.56")
	result := ctx.Add(a, b)
	fmt.Println(result)
	// Output:
	// 5.79000000
}

func ExampleContext_Sub() {
	ctx := fixedpoint.BasicContext()
	a := ctx.Parse("5.00")
	b := ctx.Parse("2.50")
	result := ctx.Sub(a, b)
	fmt.Println(result)
	// Output:
	// 2.50000000
}

func ExampleContext_Mul() {
	ctx := fixedpoint.BasicContext()
	a := ctx.Parse("2.00")
	b := ctx.Parse("3.50")
	result := ctx.Mul(a, b)
	fmt.Println(result)
	// Output:
	// 7.00000000
}

func ExampleContext_Div() {
	ctx := fixedpoint.BasicContext()
	a := ctx.Parse("7.00")
	b := ctx.Parse("2.00")
	result := ctx.Div(a, b)
	fmt.Println(result)
	// Output:
	// 3.50000000
}

func ExampleContext_Compare() {
	ctx := fixedpoint.BasicContext()
	a := ctx.Parse("1.23")
	b := ctx.Parse("4.56")
	fmt.Println(ctx.Compare(a, b))
	// Output:
	// -1
}
