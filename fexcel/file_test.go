package fexcel

import (
	"path/filepath"
	"testing"
)

const testDir = "testdata"

func TestOpenFile(t *testing.T) {
	fpath := filepath.Join(testDir, "test.xlsx")
	cfg := FileConfig{Offset: 1, Numregs: "Data:A2"}

	f, err := OpenFile(fpath, cfg)
	if err != nil {
		t.Fatal(err)
	}

	l := f.LocationsFor(Numreg)
	if len(l) != 1 {
		t.Errorf("expected 1 location for Numregs. got %d", len(l))
	}
	if l[0].Sheet != "Data" {
		t.Errorf("Bad sheet. Got %q, want %q", f.Locations[Numreg].Sheet, "Data")
	}
	if l[0].Axis != "A2" {
		t.Errorf("Bad axis. Got %q, want %q", f.Locations[Numreg].Axis, "A2")
	}
}

func TestMultipleLocations(t *testing.T) {
	fpath := filepath.Join(testDir, "test.xlsx")
	cfg := FileConfig{Offset: 1, Numregs: "Data:A2,Data:R2"}

	f, err := OpenFile(fpath, cfg)
	if err != nil {
		t.Fatal(err)
	}

	l := f.LocationsFor(Numreg)
	if len(l) != 2 {
		t.Errorf("expected 2 locations for Numregs. got %d", len(l))
	}
	if l[0].Sheet != "Data" {
		t.Errorf("Bad sheet. Got %q, want %q", f.Locations[Numreg].Sheet, "Data")
	}
	if l[0].Axis != "A2" {
		t.Errorf("Bad axis. Got %q, want %q", f.Locations[Numreg].Axis, "A2")
	}
	if l[1].Sheet != "Data" {
		t.Errorf("Bad sheet. Got %q, want %q", f.Locations[Numreg].Sheet, "Data")
	}
	if l[1].Axis != "R2" {
		t.Errorf("Bad axis. Got %q, want %q", f.Locations[Numreg].Axis, "R2")
	}

	numregs, err := f.Definitions(Numreg)
	if err != nil {
		t.Fatal(err)
	}

	if len(numregs) != 9 {
		t.Errorf("Bad # of numregs. Got %d, want %d", len(numregs), 9)
	}
}

func TestConstants(t *testing.T) {
	fpath := filepath.Join(testDir, "test.xlsx")

	f, err := OpenFile(fpath, FileConfig{
		Sheet:     "Data",
		Offset:    1,
		Constants: "M2",
	})
	if err != nil {
		t.Fatal(err)
	}

	constants, err := f.Constants()
	if err != nil {
		t.Fatal(err)
	}

	if constants["FOO"] != "bar" {
		t.Errorf("Expected %q, got %q", "bar", constants["FOO"])
	}

	if constants["BAZ"] != "3.14" {
		t.Errorf("Expected %q, got %q", "3.14", constants["BAZ"])
	}
}

func TestDefinitions(t *testing.T) {
	fpath := filepath.Join(testDir, "test.xlsx")

	cfg := FileConfig{
		Sheet:     "Data",
		Offset:    1,
		Constants: "M2",
		Numregs:   "A2",
		Posregs:   "D2",
		Sregs:     "G2",
		Flags:     "J2",
		Dins:      "IO:A2",
		Douts:     "IO:C2",
		Rins:      "IO:E2",
		Routs:     "IO:G2",
		Gins:      "IO:I2",
		Gouts:     "IO:K2",
		Ains:      "IO:M2",
		Aouts:     "IO:O2",
		Ualms:     "Alarms:A2",
	}

	f, err := OpenFile(fpath, cfg)
	if err != nil {
		t.Fatal(err)
	}

	expected := []struct {
		Type
		defs []Definition
	}{
		{
			Numreg,
			[]Definition{
				{Numreg, "Data", 1, 2, 1, "this is an extremely long comment"},
				{Numreg, "Data", 1, 3, 2, "two"},
				{Numreg, "Data", 1, 4, 3, "three"},
				{Numreg, "Data", 1, 5, 4, "four"},
				{Numreg, "Data", 1, 6, 5, "five"},
			},
		},
		{
			Posreg,
			[]Definition{
				{Posreg, "Data", 4, 2, 1, "pr1"},
				{Posreg, "Data", 4, 3, 2, "pr2"},
				{Posreg, "Data", 4, 4, 3, "pr3"},
				{Posreg, "Data", 4, 5, 4, "pr4"},
				{Posreg, "Data", 4, 6, 5, "pr5"},
			},
		},
		{
			Sreg,
			[]Definition{
				{Sreg, "Data", 7, 2, 1, "sreg1"},
				{Sreg, "Data", 7, 3, 2, "sreg2"},
			},
		},
		{
			Din,
			[]Definition{
				{Din, "IO", 1, 2, 1, "din1"},
				{Din, "IO", 1, 3, 2, "din2"},
				{Din, "IO", 1, 4, 3, "din3"},
			},
		},
		{
			Dout,
			[]Definition{
				{Dout, "IO", 3, 2, 1, "dout1"},
				{Dout, "IO", 3, 3, 2, "dout2"},
				{Dout, "IO", 3, 4, 3, "dout3"},
				{Dout, "IO", 3, 5, 4, "dout4"},
			},
		},
		{
			Rin,
			[]Definition{
				{Rin, "IO", 5, 2, 1, "rin1"},
				{Rin, "IO", 5, 3, 2, "rin2"},
			},
		},
		{
			Rout,
			[]Definition{
				{Rout, "IO", 7, 2, 1, "rout1"},
			},
		},
		{
			Gin,
			[]Definition{
				{Gin, "IO", 9, 2, 1, "gin1"},
			},
		},
		{
			Gout,
			[]Definition{
				{Gout, "IO", 11, 2, 1, "gout1"},
			},
		},
		{
			Ain,
			[]Definition{
				{Ain, "IO", 13, 2, 1, "ain1"},
			},
		},
		{
			Aout,
			[]Definition{
				{Aout, "IO", 15, 2, 1, "aout1"},
			},
		},
		{
			Ualm,
			[]Definition{
				{Ualm, "Alarms", 1, 2, 1, "test"},
				{Ualm, "Alarms", 1, 3, 2, "test two"},
				{Ualm, "Alarms", 1, 4, 3, "test three"},
				{Ualm, "Alarms", 1, 5, 4, "test four"},
			},
		},
	}

	for _, e := range expected {
		defs, err := f.Definitions(e.Type)
		if err != nil {
			t.Errorf("Failed to get defs for %s: %q", e.Type, err)
			continue
		}

		if len(defs) != len(e.defs) {
			t.Errorf("Bad # of defs for %s. Got %d, want %d", e.Type, len(defs), len(e.defs))
			continue
		}

		for id, def := range defs {
			if def.Type != e.Type {
				t.Errorf("Bad Type. Got %q, want %q", def.Type, e.Type)
			}

			if def.Sheet != e.defs[id].Sheet {
				t.Errorf("Bad sheet. Got %q want %q", def.Sheet, e.defs[id].Sheet)
			}

			if def.Column != e.defs[id].Column {
				t.Errorf("Bad column. Got %d want %d", def.Column, e.defs[id].Column)
			}

			if def.Row != e.defs[id].Row {
				t.Errorf("Bad row. Got %d want %d", def.Row, e.defs[id].Row)
			}

			if def.Id != e.defs[id].Id {
				t.Errorf("Bad id. Got %d, want %d", def.Id, e.defs[id].Id)
			}

			if def.Comment != e.defs[id].Comment {
				t.Errorf("Bad comment. Got %q, want %q", def.Comment, e.defs[id].Comment)
			}

		}
	}

	if len(f.Warnings) != 1 {
		t.Fatal("Expected 1 warning")
	}

	want := "comment in [Data]B2 for R[1] will be truncated to \"this is an extre\" (length 33 > max length 16 for Rs)"
	if f.Warnings[0] != want {
		t.Errorf("Bad warning. Got %q, want %q", f.Warnings[0], want)
	}
}

func TestNewLocation(t *testing.T) {
	tests := []struct {
		spec         string
		defaultSheet string
		expAxis      string
		expSheet     string
		expOffset    int
	}{
		{"A2", "Foo", "A2", "Foo", 0},
		{"Bar:A2", "Foo", "A2", "Bar", 0},
		{"D2", "Baz", "D2", "Baz", 0},
		{"Bar:A2{5}", "Foo", "A2", "Bar", 5},
	}

	for _, test := range tests {
		l, err := NewLocation(Numreg, test.spec, test.defaultSheet)
		if err != nil {
			t.Fatal(err)
		}

		if l.Axis != test.expAxis {
			t.Errorf("Bad axis. Got %q, want %q", l.Axis, test.expAxis)
		}
		if l.Sheet != test.expSheet {
			t.Errorf("Bad sheet. Got %q, want %q", l.Sheet, test.expSheet)
		}
	}
}

func TestNewFile(t *testing.T) {
	fpath := filepath.Join(testDir, "newfile.xlsx")
	cfg := FileConfig{Offset: 1, Numregs: "Data:A2"}

	_, err := NewFile(fpath, cfg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewFileAlreadyExists(t *testing.T) {
	fpath := filepath.Join(testDir, "test.xlsx")
	cfg := FileConfig{Offset: 1, Numregs: "Data:A2"}

	_, err := NewFile(fpath, cfg)
	if err == nil {
		t.Fatal("expected an error")
	}
}
