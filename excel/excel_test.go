package excel

import (
	"path/filepath"
	"testing"

	"github.com/onerobotics/fexcel/fanuc"
)

const testDir = "testdata"

func TestNewFile(t *testing.T) {
	fpath := filepath.Join(testDir, "test.xlsx")

	f, err := NewFile(fpath, 5)
	if err != nil {
		t.Fatal(err)
	}

	if f.offset != 5 {
		t.Errorf("Bad offset. Got %d, want %d", f.offset, 5)
	}

	if f.xlsx.Path != fpath {
		t.Errorf("Bad path. Got %q, want %q", f.xlsx.Path, fpath)
	}
}

func TestSetLocation(t *testing.T) {
	fpath := filepath.Join(testDir, "test.xlsx")

	f, err := NewFile(fpath, 5)
	if err != nil {
		t.Fatal(err)
	}

	if _, defined := f.Locations[fanuc.Numreg]; defined {
		t.Error("Expected locations to be undefined by default")
	}

	f.SetLocation(fanuc.Numreg, "A2", "Sheet1")
	if f.Locations[fanuc.Numreg].Axis != "A2" {
		t.Errorf("Bad axis. Got %q, want %q", f.Locations[fanuc.Numreg].Axis, "A2")
	}
	if f.Locations[fanuc.Numreg].Sheet != "Sheet1" {
		t.Errorf("Bad sheet. Got %q, want %q", f.Locations[fanuc.Numreg].Sheet, "Sheet1")
	}
}

func TestDefinitions(t *testing.T) {
	fpath := filepath.Join(testDir, "test.xlsx")

	f, err := NewFile(fpath, 1)
	if err != nil {
		t.Fatal(err)
	}

	f.SetLocation(fanuc.Numreg, "A2", "Data")
	f.SetLocation(fanuc.Posreg, "D2", "Data")
	f.SetLocation(fanuc.Sreg, "G2", "Data")
	f.SetLocation(fanuc.Flag, "J2", "Data")
	f.SetLocation(fanuc.Din, "A2", "IO")
	f.SetLocation(fanuc.Dout, "C2", "IO")
	f.SetLocation(fanuc.Rin, "E2", "IO")
	f.SetLocation(fanuc.Rout, "G2", "IO")
	f.SetLocation(fanuc.Gin, "I2", "IO")
	f.SetLocation(fanuc.Gout, "K2", "IO")
	f.SetLocation(fanuc.Ain, "M2", "IO")
	f.SetLocation(fanuc.Aout, "O2", "IO")
	f.SetLocation(fanuc.Ualm, "A2", "Alarms")

	expected := []struct {
		fanuc.DataType
		defs []fanuc.Definition
	}{
		{
			fanuc.Numreg,
			[]fanuc.Definition{
				{fanuc.Numreg, 1, "this is an extremely long comment"},
				{fanuc.Numreg, 2, "two"},
				{fanuc.Numreg, 3, "three"},
				{fanuc.Numreg, 4, "four"},
				{fanuc.Numreg, 5, "five"},
			},
		},
		{
			fanuc.Posreg,
			[]fanuc.Definition{
				{fanuc.Posreg, 1, "pr1"},
				{fanuc.Posreg, 2, "pr2"},
				{fanuc.Posreg, 3, "pr3"},
				{fanuc.Posreg, 4, "pr4"},
				{fanuc.Posreg, 5, "pr5"},
			},
		},
		{
			fanuc.Sreg,
			[]fanuc.Definition{
				{fanuc.Sreg, 1, "sreg1"},
				{fanuc.Sreg, 2, "sreg2"},
			},
		},
		{
			fanuc.Din,
			[]fanuc.Definition{
				{fanuc.Din, 1, "din1"},
				{fanuc.Din, 2, "din2"},
				{fanuc.Din, 3, "din3"},
			},
		},
		{
			fanuc.Dout,
			[]fanuc.Definition{
				{fanuc.Dout, 1, "dout1"},
				{fanuc.Dout, 2, "dout2"},
				{fanuc.Dout, 3, "dout3"},
				{fanuc.Dout, 4, "dout4"},
			},
		},
		{
			fanuc.Rin,
			[]fanuc.Definition{
				{fanuc.Rin, 1, "rin1"},
				{fanuc.Rin, 2, "rin2"},
			},
		},
		{
			fanuc.Rout,
			[]fanuc.Definition{
				{fanuc.Rout, 1, "rout1"},
			},
		},
		{
			fanuc.Gin,
			[]fanuc.Definition{
				{fanuc.Gin, 1, "gin1"},
			},
		},
		{
			fanuc.Gout,
			[]fanuc.Definition{
				{fanuc.Gout, 1, "gout1"},
			},
		},
		{
			fanuc.Ain,
			[]fanuc.Definition{
				{fanuc.Ain, 1, "ain1"},
			},
		},
		{
			fanuc.Aout,
			[]fanuc.Definition{
				{fanuc.Aout, 1, "aout1"},
			},
		},
		{
			fanuc.Ualm,
			[]fanuc.Definition{
				{fanuc.Ualm, 1, "test"},
				{fanuc.Ualm, 2, "test two"},
				{fanuc.Ualm, 3, "test three"},
				{fanuc.Ualm, 4, "test four"},
			},
		},
	}

	for _, e := range expected {
		defs, err := f.Definitions(e.DataType)
		if err != nil {
			t.Errorf("Failed to get defs for %s: %q", e.DataType, err)
			continue
		}

		if len(defs) != len(e.defs) {
			t.Errorf("Bad # of defs for %s. Got %d, want %d", e.DataType, len(defs), len(e.defs))
			continue
		}

		for id, def := range defs {
			if def.DataType != e.DataType {
				t.Errorf("Bad DataType. Got %q, want %q", def.DataType, e.DataType)
			}

			if def.Id != e.defs[id].Id {
				t.Errorf("Bad id. Got %d, want %d", def.Id, e.defs[id].Id)
			}

			if def.Comment != e.defs[id].Comment {
				t.Errorf("Bad comment. Got %q, want %q", def.Comment, e.defs[id].Comment)
			}

		}
	}
}
