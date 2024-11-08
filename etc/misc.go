package etc

import "golang.org/x/exp/constraints"

// multiply two numbers and divide by a third of the same type (often useful for constraints.Integereger math)
func MulDiv[T constraints.Integer](value, numerator, denominator T) T {
	return value * numerator / denominator
}

// multiply two numbers and divide by a third of the same type, rounding up
func MulDivRoundUp[T constraints.Integer](value, numerator, denominator T) T {
	return MulDiv(value+denominator-1, numerator, denominator)
}
