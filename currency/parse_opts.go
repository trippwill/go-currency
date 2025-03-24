package currency

import "fmt"

type ParseOpts struct {
	thousands rune
	decimal   rune
}

var DefaultParseOpts *ParseOpts = new(ParseOpts).Init(',', '.')

func (po *ParseOpts) Init(thousands, decimal rune) *ParseOpts {
	po.thousands = thousands
	po.decimal = decimal
	return po
}

type ParseError struct {
	Input string
	Inner error
}

func (pe ParseError) Error() string {
	return fmt.Sprintf("failed to parse %q: %s", pe.Input, pe.Inner.Error())
}
