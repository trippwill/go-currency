package fixedpoint

type Precision uint8 // Precision represents the number of significant digits in a FixedPoint value.

const (
	PrecisionMinimum   Precision = 3
	PrecisionDefault32 Precision = 5
	PrecisionMaximum32 Precision = 7
	PrecisionDefault64 Precision = 9
	PrecisionMaximum64 Precision = 16
)
