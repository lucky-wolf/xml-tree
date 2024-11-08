package etc

// Int is a type that represents any integer type
// todo: is there really no defined type for this in the standard library?
type Int interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

// multiply two numbers and divide by a third of the same type (often useful for integer math)
func MulDiv[T Int](value, numerator, denominator T) T {
	return value * numerator / denominator
}

// multiply two numbers and divide by a third of the same type, rounding up
func MulDivRoundUp[T Int](value, numerator, denominator T) T {
	return MulDiv(value+denominator-1, numerator, denominator)
}
