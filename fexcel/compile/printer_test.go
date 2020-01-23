package compile

import (
	"testing"

	"github.com/unreal/fexcel/fexcel"
)

func TestPrinter(t *testing.T) {
	p, err := NewPrinter("testdata/test.xlsx", fexcel.FileConfig{
		Numregs: "A2",
		Posregs: "D2",
		Sheet:   "Data",
		Offset:  1,
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

func TestPrinterErrors(t *testing.T) {
	p, err := NewPrinter("testdata/test.xlsx", fexcel.FileConfig{
		Numregs: "A2",
		Posregs: "D2",
		Sheet:   "Data",
		Offset:  1,
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
