package fexcel

import (
	"path/filepath"
	"testing"

	fanuc "github.com/onerobotics/go-fanuc"
)

const testDir = "testdata"

func validFile(t *testing.T) *File {
	fpath := filepath.Join(testDir, "test.xlsx")
	cfg := FileConfig{Offset: 1, Numregs: "Data:A2"}

	f, err := NewFile(fpath, cfg)
	if err != nil {
		t.Fatal(err)
	}

	return f
}

func TestNewFile(t *testing.T) {
	f := validFile(t)

	if f.Locations[fanuc.Numreg].Sheet != "Data" {
		t.Errorf("Bad sheet. Got %q, want %q", f.Locations[fanuc.Numreg].Sheet, "Data")
	}
	if f.Locations[fanuc.Numreg].Axis != "A2" {
		t.Errorf("Bad axis. Got %q, want %q", f.Locations[fanuc.Numreg].Axis, "A2")
	}
}

func TestDefinitions(t *testing.T) {
	fpath := filepath.Join(testDir, "test.xlsx")

	cfg := FileConfig{
		Sheet:   "Data",
		Offset:  1,
		Numregs: "A2",
		Posregs: "D2",
		Sregs:   "G2",
		Flags:   "J2",
		Dins:    "IO:A2",
		Douts:   "IO:C2",
		Rins:    "IO:E2",
		Routs:   "IO:G2",
		Gins:    "IO:I2",
		Gouts:   "IO:K2",
		Ains:    "IO:M2",
		Aouts:   "IO:O2",
		Ualms:   "Alarms:A2",
	}

	f, err := NewFile(fpath, cfg)
	if err != nil {
		t.Fatal(err)
	}
	err = f.Open()
	if err != nil {
		t.Fatal(err)
	}

	expected := []struct {
		fanuc.Type
		defs []Definition
	}{
		{
			fanuc.Numreg,
			[]Definition{
				{fanuc.Numreg, 1, "this is an extremely long comment"},
				{fanuc.Numreg, 2, "two"},
				{fanuc.Numreg, 3, "three"},
				{fanuc.Numreg, 4, "four"},
				{fanuc.Numreg, 5, "five"},
			},
		},
		{
			fanuc.Posreg,
			[]Definition{
				{fanuc.Posreg, 1, "pr1"},
				{fanuc.Posreg, 2, "pr2"},
				{fanuc.Posreg, 3, "pr3"},
				{fanuc.Posreg, 4, "pr4"},
				{fanuc.Posreg, 5, "pr5"},
			},
		},
		{
			fanuc.Sreg,
			[]Definition{
				{fanuc.Sreg, 1, "sreg1"},
				{fanuc.Sreg, 2, "sreg2"},
			},
		},
		{
			fanuc.Din,
			[]Definition{
				{fanuc.Din, 1, "din1"},
				{fanuc.Din, 2, "din2"},
				{fanuc.Din, 3, "din3"},
			},
		},
		{
			fanuc.Dout,
			[]Definition{
				{fanuc.Dout, 1, "dout1"},
				{fanuc.Dout, 2, "dout2"},
				{fanuc.Dout, 3, "dout3"},
				{fanuc.Dout, 4, "dout4"},
			},
		},
		{
			fanuc.Rin,
			[]Definition{
				{fanuc.Rin, 1, "rin1"},
				{fanuc.Rin, 2, "rin2"},
			},
		},
		{
			fanuc.Rout,
			[]Definition{
				{fanuc.Rout, 1, "rout1"},
			},
		},
		{
			fanuc.Gin,
			[]Definition{
				{fanuc.Gin, 1, "gin1"},
			},
		},
		{
			fanuc.Gout,
			[]Definition{
				{fanuc.Gout, 1, "gout1"},
			},
		},
		{
			fanuc.Ain,
			[]Definition{
				{fanuc.Ain, 1, "ain1"},
			},
		},
		{
			fanuc.Aout,
			[]Definition{
				{fanuc.Aout, 1, "aout1"},
			},
		},
		{
			fanuc.Ualm,
			[]Definition{
				{fanuc.Ualm, 1, "test"},
				{fanuc.Ualm, 2, "test two"},
				{fanuc.Ualm, 3, "test three"},
				{fanuc.Ualm, 4, "test four"},
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

			if def.Id != e.defs[id].Id {
				t.Errorf("Bad id. Got %d, want %d", def.Id, e.defs[id].Id)
			}

			if def.Comment != e.defs[id].Comment {
				t.Errorf("Bad comment. Got %q, want %q", def.Comment, e.defs[id].Comment)
			}

		}
	}
}

func TestNewLocation(t *testing.T) {
	tests := []struct {
		spec         string
		defaultSheet string
		expAxis      string
		expSheet     string
	}{
		{"A2", "Foo", "A2", "Foo"},
		{"Bar:A2", "Foo", "A2", "Bar"},
		{"D2", "Baz", "D2", "Baz"},
	}

	for _, test := range tests {
		l, err := NewLocation(test.spec, test.defaultSheet)
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

func TestOpen(t *testing.T) {
	f := validFile(t)
	if f.xlsx != nil {
		t.Fatal("xlsx should not be nil")
	}
	err := f.Open()
	if err != nil {
		t.Fatal(err)
	}
	if f.xlsx == nil {
		t.Fatal("xlsx is nil")
	}
}

func TestNew(t *testing.T) {
	var f File
	if f.xlsx != nil {
		t.Fatal("xlsx should not be nil")
	}
	f.New()
	if f.xlsx == nil {
		t.Fatal("xlsx is nil")
	}
}
