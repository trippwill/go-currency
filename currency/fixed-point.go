package currency

import (
	"math"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

// FixedPoint is a fixed-point number that combines an integral unit (Base)
// with a specified number of fractional digits (Precision)
type FixedPoint struct {
	Base  int64
	Scale uint8
}

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

func NewFixedPoint(value string, o *ParseOpts) FixedPoint {
	if o == nil {
		o = DefaultParseOpts
	}

	// Remove thousands separators if present.
	if o.thousands != 0 {
		value = strings.ReplaceAll(value, string(o.thousands), "")
	}

	sign := int64(1)
	if strings.HasPrefix(value, "-") {
		sign = -1
		value = value[1:]
	}

	parts := strings.Split(value, string(o.decimal))
	if len(parts) == 0 || len(parts) > 2 {
		panic("invalid number format")
	}

	integerPart, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		panic(err)
	}

	precision := uint8(0)
	if len(parts) == 2 {
		precision = uint8(len(parts[1]))
		fractionalPart, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			panic(err)
		}
		if fractionalPart == 0 && integerPart == 0 {
			precision = 0
		} else {
			scale := int64(math.Pow10(int(precision)))
			integerPart = sign * (integerPart*scale + fractionalPart)
		}
	}

	return FixedPoint{
		Base:  sign * integerPart,
		Scale: precision,
	}
}

// Format formats the FixedPoint value based on the given language tag.
func (fp *FixedPoint) Format(tag language.Tag) string {
	scaled := float64(fp.Base) / math.Pow10(int(fp.Scale))

	// Create a Printer for the desired locale
	p := message.NewPrinter(tag)

	// Format with exactly fp.Precision digits after the decimal point
	// (so 0.01000 is displayed as "0.01000" in e.g. English, or "0,01000" in e.g. French).
	return p.Sprintf("%v",
		number.Decimal(
			scaled,
			number.Scale(int(fp.Scale)),
		),
	)
}

func (fp *FixedPoint) String() string {
	return fp.Format(language.Tag{})
}
