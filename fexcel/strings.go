package fexcel

func Pluralize(word string, i int) string {
	if i == 1 {
		return word
	} else {
		return word + "s"
	}
}
