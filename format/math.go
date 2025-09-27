package format

import (
	"math"
	"strconv"
	"strings"
)

var siPrefixes = map[int]string{
	-18: "a", // atto
	-15: "f", // femto
	-12: "p", // pico
	-9:  "n", // nano
	-6:  "µ", // micro
	-3:  "m", // milli
	0:   "",  // (none)
	3:   "k", // kilo
	6:   "M", // mega
	9:   "G", // giga
	12:  "T", // tera
	15:  "P", // peta
	18:  "E", // exa
}

// formats a float64 into a human natural representation.
// significant controls significant digits.
// uses engineering notation for large/small numbers, otherwise natural decimal.
func Natural(value float64, significant int) (s string) {
	if value == 0 {
		s = "0"
		return
	}
	mantissa, exp := Naturalize(value)
	s = Mantissa(mantissa, significant)
	if exp != 0 {
		s += Exponent(exp)
	}
	return
}

// formats a float64 into a human natural representation.
// significant controls significant digits.
// uses SI notation for large/small numbers, otherwise natural decimal.
func NaturalSI(value float64, significant int) (s string) {
	if value == 0 {
		s = "0"
		return
	}
	mantissa, exp := Naturalize(value)
	s = Mantissa(mantissa, significant)
	if exp != 0 {
		if prefix, ok := siPrefixes[exp]; ok {
			s += prefix
		} else {
			s += Exponent(exp)
		}
	}
	return
}

// Naturalize takes a value and returns mantissa + exponent (multiple of 3)
// such that value = mantissa * 10^exp, with any "human" adjustments applied.
func Naturalize(value float64) (mantissa float64, exp int) {
	if value == 0 {
		return
	}
	abs := math.Abs(value)

	// Find nearest engineering exponent
	exp = int(math.Floor(math.Log10(abs)/3.0) * 3)

	// Adjust down if the previous exponent gives a more natural mantissa
	prevExp := exp - 3
	prevMant := value / math.Pow(10, float64(prevExp))
	if math.Abs(prevMant) < 1e6 {
		exp = prevExp
	}

	// Adjust up if the next exponent gives a more natural mantissa (small numbers)
	nextExp := exp + 3
	nextMant := value / math.Pow(10, float64(nextExp))
	if math.Abs(nextMant) >= 1 {
		exp = nextExp
	}

	mantissa = value / math.Pow(10, float64(exp))
	return
}

// Exponent formats an Exponent as "eN"
func Exponent(exp int) string {
	if exp == 0 {
		return ""
	}
	return "e" + strconv.Itoa(exp)
}

// ExponentToSI formats an exponent as an SI prefix, or falls back to "eN"
func ExponentToSI(exp int) string {
	if exp == 0 {
		return ""
	}
	if prefix, ok := siPrefixes[exp]; ok {
		return prefix
	}
	return Exponent(exp)
}

// formats the mantissa portion of a float with the given number of significant digits
func Mantissa(value float64, significant int) string {
	if value == 0 {
		return "0"
	}

	abs := math.Abs(value)
	exp10 := int(math.Floor(math.Log10(abs)))
	// digits after decimal = significant - (integer digits)
	// example: value=468750 -> exp10=5, significant=5 -> digits=significant-exp10-1 = -1 => 0 decimals
	digits := max(significant-exp10-1, 0)

	return SimplifyNumber(strconv.FormatFloat(value, 'f', digits, 64))
}

// trims trailing zeros and a trailing decimal point if needed
func SimplifyNumber(s string) string {
	if strings.Contains(s, ".") {
		s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	}
	return s
}
