package fexcel

import (
	"testing"
)

func TestHasOverlaps(t *testing.T) {
	tests := []struct {
		numregs     string
		posregs     string
		dins        string
		offset      int
		hasOverlaps bool
	}{
		{"A1", "C1", "E1", 1, false},
		{"A1", "C1", "E1", 2, true},
		{"A1", "Foo:C1", "Bar:E1", 2, false}, // no overlap because on different sheets
		{"A1", "B1", "C1", 3, true},          // not _really_, but this is a stupid spreadsheet design
		{"A1", "C1", "E1,B1", 1, true},       // DINs B1 overlaps with numregs A1
	}

	for id, test := range tests {
		cfg := FileConfig{
			Numregs: test.numregs,
			Posregs: test.posregs,
			Dins:    test.dins,
			Offset:  test.offset,
			Sheet:   "Default",
		}

		result, err := cfg.HasOverlaps()
		if err != nil {
			t.Fatal(err)
		}

		if result != test.hasOverlaps {
			t.Errorf("HasOverlaps(%d): Got %t, want %t", id, result, test.hasOverlaps)
		}
	}
}

func TestSpecWithComma(t *testing.T) {
	cfg := FileConfig{
		Numregs: "Sheet1:A2,Sheet2:A2",
	}

	locs, err := cfg.Locations()
	if err != nil {
		t.Fatal(err)
	}

	if len(locs[Numreg]) != 2 {
		t.Errorf("expected 2 location for numregs, got %d", len(locs[Numreg]))
	}
}
