package main

import (
	"fmt"
	"strings"

	"github.com/lucky-wolf/xml-tree/format"
)

func main() {
	// Table-driven demo for format.Natural
	testCases := []struct {
		desc  string
		value float64
	}{
		{"Fractional", 0.75},
		{"Small float", 0.0080526},
		{"Large int", 1234567},
		{"Zero", 0},
		{"Small int", 7},
		{"Small float", 0.000123},
		{"Large float", 987654321.123},
		{"Negative int", -42},
		{"Negative float", -0.000987},
		{"Pi", 3.1415926535},
		{"Tiny float", 1e-9},
		{"Huge float", 1.23e15},
		{"Recurring decimal", 1.0 / 3.0},
	}

	fmt.Printf("%-20s %-20s %-20s %-20s\n", "Description", "Value", "format.Natural (3)", "format.Natural (5)")
	fmt.Println(strings.Repeat("-", 80))
	for _, tc := range testCases {
		result3 := format.Natural(tc.value, 3)
		result5 := format.Natural(tc.value, 5)
		fmt.Printf("%-20s %-20g %-20s %-20s\n", tc.desc, tc.value, result3, result5)
	}
}
