package fexcel

import (
	"testing"

	fanuc "github.com/onerobotics/go-fanuc"
)

func TestNewTarget(t *testing.T) {
	target, err := NewTarget("127.0.0.1", 5)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := target.client.(*fanuc.HTTPClient); !ok {
		t.Errorf("Expected an HTTPClient. Got %T", target.client)
	}

	target, err = NewTarget("testdata", 0)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := target.client.(*fanuc.FileClient); !ok {
		t.Errorf("Expected a FileClient. Got %T", target.client)
	}
}

func TestGetComments(t *testing.T) {
	target, err := NewTarget("testdata", 0)
	if err != nil {
		t.Fatal(err)
	}

	err = target.GetComments(Numreg)
	if err != nil {
		t.Fatal(err)
	}

	if numregs, ok := target.Comments[Numreg]; ok {
		if len(numregs) != 200 {
			t.Fatalf("Only got %d numregs. Want 200", len(numregs))
		}
		tests := []struct {
			id      int
			comment string
		}{
			{1, "this is an extre"},
			{2, "two"},
			{10, "UngripDelay"},
		}

		for _, test := range tests {
			if r, ok := numregs[test.id]; ok {
				if r != test.comment {
					t.Errorf("Bad comment for R[%d]. Got %q, want %q", test.id, r, test.comment)
				}
			} else {
				t.Errorf("R[%d] undefined", test.id)
			}
		}
	} else {
		t.Errorf("numregs not found")
	}
}
