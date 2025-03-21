package currency

type Code string
type Symbol string
type Factor uint16

type Currency interface {
	GetCode() Code
	GetSymbol() Symbol
	GetMinorUnitFactor() Factor
}
