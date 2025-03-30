package fixedpoint

import (
	"testing"
)

func newTestContext(prec int, rounding Rounding) context {
	return context{
		precision: Precision(prec),
		rounding:  rounding,
	}
}

// Fuzz test for FiniteNumber exported operations.
func FuzzFiniteNumberOperations(f *testing.F) {
	// Seed corpus.
	f.Add(true, uint64(12345), 2, false, uint64(67890), 3, uint64(10))
	f.Fuzz(func(t *testing.T, sign1 bool, coeff1 uint64, exp1 int, sign2 bool, coeff2 uint64, exp2 int, precision uint64) {
		// Limit values to avoid overflow.
		prec := int(precision%10 + 1)
		ctx := newTestContext(prec, RoundHalfUp)
		// Create FiniteNumbers (limit coefficient and exponent ranges).
		a := new(FiniteNumber).Init(sign1, coefficient(coeff1%100000), exponent(exp1%10), ctx)
		b := new(FiniteNumber).Init(sign2, coefficient(coeff2%100000), exponent(exp2%10), ctx)
		// Call exported operations.
		_ = a.Add(b)
		_ = a.Sub(b)
		_ = a.Mul(b)
		_ = a.Div(b)
		_ = a.Neg()
		_ = a.Abs()
		_ = a.Compare(b)
	})
}

// Fuzz test for Infinity exported operations.
func FuzzInfinityOperations(f *testing.F) {
	// Seed corpus.
	f.Add(true, uint64(12345), 2, false, uint64(67890), 3, uint64(10))
	f.Fuzz(func(t *testing.T, sign bool, coeff uint64, exp int, finSign bool, finCoeff uint64, finExp int, precision uint64) {
		prec := int(precision%10 + 1)
		ctx := newTestContext(prec, RoundHalfUp)
		a := new(Infinity).Init(sign, ctx)
		b := new(FiniteNumber).Init(finSign, coefficient(finCoeff%100000), exponent(finExp%10), ctx)
		// Ensure finite number coefficient is nonzero to avoid division-by-zero.
		if b.coe == 0 {
			b.coe = 1
		}
		_ = a.Add(b)
		_ = a.Sub(b)
		_ = a.Mul(b)
		_ = a.Div(b)
		_ = a.Neg()
		_ = a.Abs()
		_ = a.Compare(b)
	})
}
