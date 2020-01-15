package fexcel

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	fanuc "github.com/onerobotics/go-fanuc"
)

func TestNewCreator(t *testing.T) {
	cfg := Config{Offset: 1, Sheet: "Sheet1", Numregs: "A2"}

	// must be an xlsx file
	_, err := NewCreator("foo.bar", cfg, "127.0.0.1")
	if err == nil {
		t.Fatal("Expected an error")
	}
	want := "File path must end in .xlsx"
	if err.Error() != want {
		t.Errorf("Bad error msg. Got %q, want %q", err.Error(), want)
	}

	// file can't already exist
	_, err = NewCreator("./testdata/test.xlsx", cfg, "testdata")
	if err == nil {
		t.Fatal("expected an error")
	}
	want = "File already exists"
	if err.Error() != want {
		t.Errorf("Bad error msg. Got %q, want %q", err.Error(), want)
	}

	_, err = NewCreator("./testdata/test2.xlsx", cfg, "testdata")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreatorCreate(t *testing.T) {
	dir, err := ioutil.TempDir("testdata", "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fpath := filepath.Join(dir, "test.xlsx")
	cfg := Config{Offset: 1, Sheet: "Sheet1", Numregs: "A2"}

	c, err := NewCreator(fpath, cfg, "testdata")
	if err != nil {
		t.Fatal(err)
	}

	err = c.Create(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}

	// verify with a DiffCommand
	cmd, err := NewDiffCommand(fpath, cfg, 500, "testdata")
	if err != nil {
		t.Fatal(err)
	}

	comparisons, err := cmd.Compare(fanuc.Numreg)
	if err != nil {
		t.Fatal(err)
	}

	if len(comparisons) != 200 {
		t.Errorf("Got %d comparisons. Want 200", len(comparisons))
	}

	for _, c := range comparisons {
		if !c.Equal() {
			t.Error("Comparison not equal")
		}
	}
}
