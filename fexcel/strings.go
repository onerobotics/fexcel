package fexcel

import (
	fanuc "github.com/onerobotics/go-fanuc"
)

// No, this does not account for special cases (like words ending in s or
// or the word "fish"), but it's good enough for fexcel's purposes
func Pluralize(word string, i int) string {
	if i == 1 {
		return word
	} else {
		return word + "s"
	}
}

func MaxLengthFor(t fanuc.Type) int {
	switch t {
	case fanuc.Numreg, fanuc.Posreg, fanuc.Sreg:
		return 16
	case fanuc.Ualm:
		return 29
	default:
		return 24
	}
}

func Truncated(s string, t fanuc.Type) string {
	if max := MaxLengthFor(t); len(s) > max {
		return s[:max]
	}
	return s
}
