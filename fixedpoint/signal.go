package fixedpoint

type Signal uint8

const SignalClear Signal = 0

const (
	SignalOverflow Signal = 1 << iota
	SignalUnderflow
	SignalDivisionByZero
	SignalDivisionImpossible
	SignalInexact
	SignalInvalidOperation
	s_conversionSyntax
)

const (
	SignalConversionSyntax = s_conversionSyntax | SignalInvalidOperation
)

func (s Signal) String() string {
	switch s {
	case SignalClear:
		return "SignalClear"
	case SignalOverflow:
		return "SignalOverflow"
	case SignalUnderflow:
		return "SignalUnderflow"
	case SignalDivisionByZero:
		return "SignalDivisionByZero"
	case SignalDivisionImpossible:
		return "SignalDivisionImpossible"
	case SignalInexact:
		return "SignalInexact"
	case SignalInvalidOperation:
		return "SignalInvalidOperation"
	case SignalConversionSyntax:
		return "SignalConversionSyntax"
	default:
		return "Signal(0x" + s.String() + ")"
	}
}
