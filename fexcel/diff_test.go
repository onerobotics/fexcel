package fexcel

import (
	"testing"

	fanuc "github.com/onerobotics/go-fanuc"
)

func TestDiffCompare(t *testing.T) {
	cfg := Config{
		FileConfig: FileConfig{
			Numregs: "Data:A2",
			Offset:  1,
		},
	}

	cmd, err := NewDiffCommand("testdata/test.xlsx", cfg, "testdata")
	if err != nil {
		t.Fatal(err)
	}

	results, err := cmd.Compare(fanuc.Numreg)
	if err != nil {
		t.Fatal(err)
	}

	wants := []struct {
		id   int
		want string
		got  string
		eql  bool
	}{
		{1, "this is an extremely long comment", "this is an extre", false},
		{2, "two", "two", true},
		{3, "three", "three", true},
		{4, "four", "four", true},
		{5, "five", "five", true},
	}

	for id, want := range wants {
		result := results[id]
		if result.Id != want.id {
			t.Errorf("Bad id. Got %d, want %d", result.Id, want.id)
		}
		if result.Want != want.want {
			t.Errorf("Bad want. Got %q, want %q", result.Want, want.want)
		}
		if result.Got["testdata"] != want.got {
			t.Errorf("Bad got. Got %q, want %q", result.Got["testdata"], want.got)
		}
		if result.Equal() != want.eql {
			t.Errorf("Bad eql. Got %t, want %t", result.Equal(), want.eql)
		}
	}
}
