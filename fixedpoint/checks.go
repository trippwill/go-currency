package fixedpoint

type FixedPointChecks interface {
	IsFinite() bool
	IsSpecial() bool
	IsNaN() bool
	IsInf() bool
	IsZero() bool
	IsNegative() bool
	IsPositive() bool
}

var (
	_ FixedPointChecks = (*FiniteNumber)(nil)
	_ FixedPointChecks = (*Infinity)(nil)
	_ FixedPointChecks = (*NaN)(nil)
)

func (a *FiniteNumber) IsFinite() bool { return a != nil }
func (a *Infinity) IsFinite() bool     { return false }
func (a *NaN) IsFinite() bool          { return false }

func (a *FiniteNumber) IsSpecial() bool { return a != nil }
func (a *Infinity) IsSpecial() bool     { return true }
func (a *NaN) IsSpecial() bool          { return true }

func (a *FiniteNumber) IsNaN() bool { return a == nil }
func (a *Infinity) IsNaN() bool     { return false }
func (a *NaN) IsNaN() bool          { return true }

func (a *FiniteNumber) IsInf() bool { return false }
func (a *Infinity) IsInf() bool     { return a != nil }
func (a *NaN) IsInf() bool          { return false }

func (a *FiniteNumber) IsZero() bool { return a != nil && a.coe == 0 && a.exp == 0 }
func (a *Infinity) IsZero() bool     { return false }
func (a *NaN) IsZero() bool          { return false }

func (a *FiniteNumber) IsNegative() bool { return a != nil && a.sign }
func (a *Infinity) IsNegative() bool     { return a != nil && a.sign }
func (a *NaN) IsNegative() bool          { return false }

func (a *FiniteNumber) IsPositive() bool { return a != nil && !a.sign }
func (a *Infinity) IsPositive() bool     { return a != nil && !a.sign }
func (a *NaN) IsPositive() bool          { return false }
