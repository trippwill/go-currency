package fixedpoint

const (
	RoundingNearestEven Rounding = iota
	RoundingNearestAway
	RoundingZero
	RoundingPositiveInf
	RoundingNegativeInf
	RoundingPositiveZero
	RoundingNegativeZero
)

func (r Rounding) Debug() string {
	switch r {
	case RoundingNearestEven:
		return "NE"
	case RoundingNearestAway:
		return "NA"
	case RoundingZero:
		return "Z"
	case RoundingPositiveInf:
		return "+Inf"
	case RoundingNegativeInf:
		return "-Inf"
	case RoundingPositiveZero:
		return "+Z"
	case RoundingNegativeZero:
		return "-Z"
	default:
		return "?(0x" + r.String() + ")"
	}
}

func (r Rounding) String() string {
	switch r {
	case RoundingNearestEven:
		return "RoundingNearestEven"
	case RoundingNearestAway:
		return "RoundingNearestAway"
	case RoundingZero:
		return "RoundingZero"
	case RoundingPositiveInf:
		return "RoundingPositiveInf"
	case RoundingNegativeInf:
		return "RoundingNegativeInf"
	case RoundingPositiveZero:
		return "RoundingPositiveZero"
	case RoundingNegativeZero:
		return "RoundingNegativeZero"
	default:
		return "Rounding(0x" + r.String() + ")"
	}
}
