package fixedpoint

type FixedPointChecks interface {
	IsOk() bool
	IsFinite() bool
	IsSpecial() bool
	IsNaN() bool
	IsInf() bool
	IsZero() bool
	IsNegative() bool
	IsPositive() bool
}

func (c *context) IsOk() bool {
	switch {
	case c == nil:
		return false
	case c.signal != SignalClear:
		return false
	case c.precision < 7 || c.precision > 19:
		return false
	}

	return true
}

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
