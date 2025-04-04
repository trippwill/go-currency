package currency

var _ Currency = &USD{}
var _ Currency = &GBP{}
var _ Currency = &JPY{}

type USD struct{}

// GetCode implements Currency.
func (u USD) GetCode() Code {
	return "USD"
}

// GetMinorUnitFactor implements Currency.
func (u USD) GetMinorUnitFactor() Factor {
	return 100
}

// GetSymbol implements Currency.
func (u USD) GetSymbol() Symbol {
	return "$"
}

type GBP struct{}

// GetCode implements Currency.
func (g GBP) GetCode() Code {
	return "GBP"
}

// GetMinorUnitFactor implements Currency.
func (g GBP) GetMinorUnitFactor() Factor {
	return 100
}

// GetSymbol implements Currency.
func (g GBP) GetSymbol() Symbol {
	return "£"
}

type JPY struct{}

// GetCode implements Currency.
func (j JPY) GetCode() Code {
	return "JPY"
}

// GetMinorUnitFactor implements Currency.
func (j JPY) GetMinorUnitFactor() Factor {
	return 1
}

// GetSymbol implements Currency.
func (j JPY) GetSymbol() Symbol {
	return "¥"
}
