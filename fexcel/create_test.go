package fexcel

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewCreator(t *testing.T) {
	cfg := Config{
		FileConfig: FileConfig{Offset: 1, Sheet: "Sheet1", Numregs: "A1"},
	}

	// must be an xlsx file
	_, err := NewCreator("foo.bar", cfg, false, "127.0.0.1")
	if err == nil {
		t.Fatal("Expected an error")
	}
	want := "File path must end in .xlsx"
	if err.Error() != want {
		t.Errorf("Bad error msg. Got %q, want %q", err.Error(), want)
	}

	// file can't already exist
	_, err = NewCreator("./testdata/test.xlsx", cfg, false, "testdata")
	if err == nil {
		t.Fatal("expected an error")
	}
	want = "File \"./testdata/test.xlsx\" already exists"
	if err.Error() != want {
		t.Errorf("Bad error msg. Got %q, want %q", err.Error(), want)
	}

	// header option fail
	_, err = NewCreator("./testdata/test.xlsx", cfg, true, "testdata")
	if err == nil {
		t.Fatal("expected an error")
	}
	want = "Cell spec for Rs (A1) must be in row 2 or lower for headers option"
	if err.Error() != want {
		t.Errorf("Bad error msg. Got %q, want %q", err.Error(), want)
	}

	// this one should be good
	_, err = NewCreator("./testdata/test2.xlsx", cfg, false, "testdata")
	if err != nil {
		t.Fatal(err)
	}

	// overlaps
	cfg.FileConfig = FileConfig{Offset: 1, Sheet: "Sheet1", Numregs: "A2", Posregs: "B2"}
	_, err = NewCreator("./testdata/test2.xlsx", cfg, false, "testdata")
	if err == nil {
		t.Error("Expected an overlap error. Got none.")
	} else {
		if err.Error() != "configuration has overlapping columns" {
			t.Errorf("Bad overlap error msg. Got %q, want %q", err.Error(), "configuration has overlapping columns")
		}
	}
}

func TestCreatorCreate(t *testing.T) {
	dir, err := ioutil.TempDir("testdata", "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fpath := filepath.Join(dir, "test.xlsx")
	cfg := Config{
		FileConfig: FileConfig{Offset: 1, Sheet: "Sheet1", Numregs: "A2"},
		Timeout:    500,
	}

	c, err := NewCreator(fpath, cfg, true, "testdata")
	if err != nil {
		t.Fatal(err)
	}

	err = c.Create(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}

	// verify with a DiffCommand
	cmd, err := NewDiffCommand(fpath, cfg, "testdata")
	if err != nil {
		t.Fatal(err)
	}

	comparisons, err := cmd.Compare(Numreg)
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
