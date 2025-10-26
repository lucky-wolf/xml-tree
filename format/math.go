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

// returns the SI prefix for a given exponent, if it exists.
func SIPrefixForExponent(exp int) (prefix string, ok bool) {
	prefix, ok = siPrefixes[exp]
	return
}

// formats a float64 into a human natural representation.
// significant controls significant digits.
// uses engineering notation for large/small numbers, otherwise natural decimal.
func Natural(value float64, significant int) (s string) {
	if value == 0 {
		s = "0"
		return
	}
	mantissa, exp := Naturalize(value, 1e-3, 1e6)
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
	mantissa, exp := Naturalize(value, 1e-3, 1e3)
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
func Naturalize(value, floor, ceiling float64) (mantissa float64, exp int) {
	if value == 0 {
		return
	}
	abs := math.Abs(value)

	// natural cutoff comes first
	if abs >= floor && abs < ceiling {
		mantissa = value
		return
	}

	// find nearest engineering exponent
	exp = int(math.Floor(math.Log10(abs)/3.0) * 3)

	// if the mantissa would be >= 1000, bump the exponent up by 3	// otherwise, just use the original exponent
	mantissa = value / math.Pow(10, float64(exp))
	return
}

// Exponent formats an Exponent as "e+-N"
func Exponent(exp int) string {
	switch {
	case exp > 9:
		return "e+" + strconv.Itoa(exp)
	case exp < -9:
		return "e-" + strconv.Itoa(-exp)
	case exp > 0:
		return "e+0" + strconv.Itoa(exp)
	case exp < 0:
		return "e-0" + strconv.Itoa(-exp)
	}
	return ""
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

func Mantissa(value float64, significant int) (s string) {
	if value == 0 {
		return "0"
	}

	abs := math.Abs(value)
	exp10 := int(math.Floor(math.Log10(abs)))

	// how many digits after the decimal to give `significant` total
	digits := max(significant-exp10-1, 0)

	// round to that many places
	scale := math.Pow(10, float64(digits))
	rounded := math.Round(value*scale) / scale

	s = strconv.FormatFloat(rounded, 'f', digits, 64)
	s = SimplifyNumber(s)
	return
}

// trims trailing zeros and a trailing decimal point if needed
func SimplifyNumber(s string) string {
	if strings.Contains(s, ".") {
		s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	}
	return s
}
