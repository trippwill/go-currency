package fixedpoint

const (
	RoundHalfUp Rounding = iota
	RoundHalfEven
	RoundDown
	RoundCeiling
	RoundFloor
)

func (r Rounding) Debug() string {
	switch r {
	case RoundHalfUp:
		return "HU"
	case RoundHalfEven:
		return "HE"
	case RoundDown:
		return "D"
	case RoundCeiling:
		return "C"
	case RoundFloor:
		return "F"
	default:
		return "?(0x" + r.String() + ")"
	}
}

func (r Rounding) String() string {
	switch r {
	case RoundHalfUp:
		return "RoundHalfUp"
	case RoundHalfEven:
		return "RoundHalfEven"
	case RoundDown:
		return "RoundDown"
	case RoundCeiling:
		return "RoundCeiling"
	case RoundFloor:
		return "RoundFloor"
	default:
		return "Rounding(0x" + r.String() + ")"
	}
}
