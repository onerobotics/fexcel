package compile

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/onerobotics/fexcel/fexcel"
)

func TestPrinter(t *testing.T) {
	p, err := NewPrinter("testdata/test.xlsx", fexcel.FileConfig{
		Constants: "G2",
		Numregs:   "A2",
		Posregs:   "D2",
		Sheet:     "Data",
		Offset:    1,
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		src string
		exp string
	}{
		{"R{one}", "R[1:one]"},
		{"R{two}", "R[2:two]"},
		{"R{three}", "R[3:three]"},
		{"&R{one}", "1"},
		{"&R{two}", "2"},
		{"&R{three}", "3"},
		{"PR{home}", "PR[4:home]"},
		{"PR{lpos}", "PR[5:lpos]"},
		{"PR{jpos}", "PR[6:jpos]"},
		{"R{one}=&PR{home}", "R[1:one]=4"},
		{"R[1:foobar]", "R[1:foobar]"},
		{"! testing {} ;", "! testing {} ;"},
		{"${HOME_SPEED}", "100"},
		{"${HOME_CNT}", "0"},
	}

	for _, test := range tests {
		p.Reset()

		f, err := Parse("", test.src)
		if err != nil {
			t.Errorf("Parse(%s): %s", test.src, err)
			continue
		}

		err = p.Print(f)
		if err != nil {
			t.Errorf("Print(%s): %s", test.src, err)
			continue
		}

		got := p.Output()
		if got != test.exp {
			t.Errorf("Output(%s). Got %q, want %q", test.src, got, test.exp)
		}
	}

}

func TestGolden(t *testing.T) {
	p, err := NewPrinter("testdata/test.xlsx", fexcel.FileConfig{
		Constants: "G2",
		Numregs:   "A2",
		Posregs:   "D2",
		Sheet:     "Data",
		Offset:    1,
	})
	if err != nil {
		t.Fatal(err)
	}

	filenames := []string{"test"}
	for _, fname := range filenames {
		src, err := ioutil.ReadFile(filepath.Join("testdata", fname+".ls"))
		if err != nil {
			t.Fatal(err)
		}

		p.Reset()
		f, err := Parse(fname+".ls", string(src))
		if err != nil {
			t.Errorf("Parse(%s): %s", fname+".ls", err)
			continue
		}

		err = p.Print(f)
		if err != nil {
			t.Errorf("Print(%s): %s", fname+".ls", err)
			continue
		}

		// compare against golden file
		golden, err := ioutil.ReadFile(filepath.Join("testdata", fname+".golden"))
		if err != nil {
			t.Fatal(err)
		}

		sLines := strings.Split(p.Output(), "\n")
		gLines := strings.Split(string(golden), "\n")

		if len(sLines) != len(gLines) {
			t.Errorf("line count mismatch, src: %d, golden: %d", len(sLines), len(gLines))
		}

		for i, _ := range sLines {
			if sLines[i] != gLines[i] {
				t.Errorf("Compare(%s) line %d: %q vs %q", fname+".ls", i+1, sLines[i], gLines[i])
			}
		}
	}
}

func TestPrinterErrors(t *testing.T) {
	p, err := NewPrinter("testdata/test.xlsx", fexcel.FileConfig{
		Constants: "G2",
		Numregs:   "A2",
		Posregs:   "D2",
		Sheet:     "Data",
		Offset:    1,
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		src string
		exp string
	}{
		{"R{lpos}", "test.ls:1:1: R{lpos} is undefined"},
		{"&R{undefined}", "test.ls:1:1: &R{undefined} is undefined"},
		{"${asdfasdf}", "test.ls:1:1: ${asdfasdf} is undefined"},
	}

	for _, test := range tests {
		p.Reset()

		f, err := Parse("test.ls", test.src)
		if err != nil {
			t.Errorf("Parse(%s): %s", test.src, err)
			continue
		}

		err = p.Print(f)
		if err == nil {
			t.Errorf("Print(%s): didn't get an error", test.src)
			continue
		}

		if err.Error() != test.exp {
			t.Errorf("Output(%s). Got %q, want %q", test.src, err.Error(), test.exp)
		}
	}
}
