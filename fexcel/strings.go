package fexcel

// No, this does not account for special cases (like words ending in s or
// or the word "fish"), but it's good enough for fexcel's purposes
func Pluralize(word string, i int) string {
	if i == 1 {
		return word
	} else {
		return word + "s"
	}
}

func MaxLengthFor(t Type) int {
	switch t {
	case Numreg, Posreg, Sreg:
		return 16
	case Ualm:
		return 29
	default:
		return 24
	}
}

func Truncated(s string, t Type) string {
	if max := MaxLengthFor(t); len(s) > max {
		return s[:max]
	}
	return s
}
