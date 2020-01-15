package fexcel

import (
	"testing"
)

func TestPluralize(t *testing.T) {
	tests := []struct {
		word  string
		count int
		want  string
	}{
		{"Numeric Register", 0, "Numeric Registers"},
		{"Numeric Register", 1, "Numeric Register"},
		{"Numeric Register", 2, "Numeric Registers"},
		{"Ualm", 5, "Ualms"},
		{"Digital Input", 1, "Digital Input"},
	}

	for _, test := range tests {
		got := Pluralize(test.word, test.count)
		if got != test.want {
			t.Errorf("Got %q, want %q", got, test.want)
		}
	}
}
