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
