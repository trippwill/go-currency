// package imath provides integer math functions for signed and unsigned integers.
package imath

// integer is a type constraint that matches all integer types, both signed and unsigned.
type integer interface {
	signed | unsigned
}

// signed is a type constraint that matches all signed integer types.
type signed interface {
	int | int8 | int16 | int32 | int64
}

// unsigned is a type constraint that matches all unsigned integer types.
type unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64
}

// Abs returns the absolute value of x.
func Abs[I integer](x I) I {
	if x < 0 {
		return -x
	}
	return x
}

// Neg returns the negation of x.
func Neg[I integer, S signed](x I) S {
	return S(-x)
}

// Pow returns x raised to the power of y.
func Pow[I integer, U unsigned](x I, y U) I {
	if y == 0 {
		return I(1)
	}
	result := I(1)
	for ; y > 0; y-- {
		result *= x
	}
	return result
}

// Clamp restricts the value `x` to the range [min, max].
// If `x` is less than `min`, it returns `min`.
// If `x` is greater than `max`, it returns `max`.
// Otherwise, it returns `x`.
func Clamp[I integer](x, min, max I) I {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

// GCD computes the greatest common divisor of `a` and `b` using the Euclidean algorithm.
// Both `a` and `b` must be non-negative integers.
func GCD[I integer](a, b I) I {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM computes the least common multiple of `a` and `b`.
// It uses the formula: LCM(a, b) = abs(a * b) / GCD(a, b).
func LCM[I integer](a, b I) I {
	if a == 0 || b == 0 {
		return 0
	}
	return Abs(a*b) / GCD(a, b)
}

// Sign returns:
// - 1 if `x` is positive,
// - -1 if `x` is negative,
// - 0 if `x` is zero.
func Sign[S signed](x S) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}
