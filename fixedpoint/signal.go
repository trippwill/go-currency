package fixedpoint

import "strings"

const SignalClear Signal = 0

const (
	SignalOverflow Signal = 1 << iota
	SignalUnderflow
	SignalDivisionByZero
	SignalDivisionImpossible
	SignalInexact
	SignalRounding
	SignalInvalidOperation
	sig_conversionSyntax
)

const (
	SignalConversionSyntax = sig_conversionSyntax | SignalInvalidOperation
)

var debugFlags = []struct {
	symbol string
	flag   Signal
}{
	{"o", SignalOverflow},
	{"u", SignalUnderflow},
	{"0", SignalDivisionByZero},
	{"P", SignalDivisionImpossible},
	{"i", SignalInexact},
	{"r", SignalRounding},
	{"X", SignalInvalidOperation},
	{"c", SignalConversionSyntax},
}

var stringFlags = []struct {
	name string
	flag Signal
}{
	{"SignalOverflow", SignalOverflow},
	{"SignalUnderflow", SignalUnderflow},
	{"SignalDivisionByZero", SignalDivisionByZero},
	{"SignalDivisionImpossible", SignalDivisionImpossible},
	{"SignalInexact", SignalInexact},
	{"SignalRounding", SignalRounding},
	{"SignalInvalidOperation", SignalInvalidOperation},
	{"SignalConversionSyntax", SignalConversionSyntax},
}

func (s Signal) Debug() string {
	var signals []string
	for _, f := range debugFlags {
		if s&f.flag != 0 {
			signals = append(signals, f.symbol)
		}
	}

	switch len(signals) {
	case 0:
		return "*"
	default:
		return strings.Join(signals, "")
	}
}

func (s Signal) String() string {
	var signals []string
	for _, f := range stringFlags {
		if s&f.flag != 0 {
			signals = append(signals, f.name)
		}
	}
	return strings.Join(signals, "|")
}
