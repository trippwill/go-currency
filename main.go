package main

import (
	"fmt"

	"github.com/trippwill/go-currency/currency"
)

func main() {
	a := currency.NewAmountFromString[currency.USD]("0.100", currency.DefaultParseOpts)
	fmt.Println("0.100", a)

	a = currency.Amount[currency.USD]{Value: currency.FixedPoint{
		Base:  1005000,
		Scale: 5,
	}}

	fmt.Println("10.05000", a)

	b := currency.NewAmount[currency.GBP](
		currency.NewFixedPoint("258.0214", nil),
	)

	fmt.Println("258.0214", b)

	// Adding two amounts of the same currency that were created with different precisions
	usdAmount1 := currency.NewAmountFromString[currency.USD]("1.23", currency.DefaultParseOpts)
	usdAmount2 := currency.NewAmountFromString[currency.USD]("4.56789", currency.DefaultParseOpts)
	sumUSD := usdAmount1.Add(usdAmount2)
	fmt.Println("USD Sum:", sumUSD) // 5.79789

	// gbpAmount := currency.NewAmountFromString(currency.GBP{}, "2.50", currency.DefaultParseOpts)
	// mixedSum := usdAmount1.Add(gbpAmount)
	// fmt.Println("Mixed Sum (USD + GBP):", mixedSum)

	mulUSD := sumUSD.Mul(currency.NewFixedPoint("2.12", nil))
	fmt.Println("USD Mul:", mulUSD) // 12.2915268
}
