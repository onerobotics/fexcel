package commenter

import "testing"

func TestParseStartCell(t *testing.T) {
	c := &commenter{DefaultSheetName: "Default"}

	sheetName, cellName, err := c.parseStartCell("A2")
	if err != nil {
		t.Fatal(err)
	}
	if sheetName != "Default" {
		t.Errorf("got %s, want %s", sheetName, "Default")
	}
	if cellName != "A2" {
		t.Errorf("got %s, want %s", cellName, "A2")
	}

	sheetName, cellName, err = c.parseStartCell("Sheet2:D2")
	if err != nil {
		t.Fatal(err)
	}
	if sheetName != "Sheet2" {
		t.Errorf("got %s, want %s", sheetName, "Sheet2")
	}
	if cellName != "D2" {
		t.Errorf("got %s, want %s", cellName, "D2")
	}

	sheetName, cellName, err = c.parseStartCell("invalid:cell:spec")
	if err == nil {
		t.Error("Expected an error")
	}
	if err.Error() != "Invalid cell string: `invalid:cell:spec`" {
		t.Errorf("got %s, want %s", err.Error(), "Invalid cell string: `invalid:cell:spec`")
	}
}

func TestWithSpreadsheet(t *testing.T) {
	var cfg Config
	cfg.Numregs = "A2"

	c, err := New("testdata/test.xlsx", "Sheet1", 1, cfg, nil)
	if err != nil {
		t.Error(err)
	}

	result, err := c.Update()
	if err != nil {
		t.Error(err)
	}

	if result.Numregs != 10 {
		t.Errorf("got %d, want %d", result.Numregs, 10)
	}
}

func TestWithSpreadsheetFormatting(t *testing.T) {
	var cfg Config
	cfg.Numregs = "A2"

	c, err := New("testdata/test.xlsx", "Formatting", 1, cfg, nil)
	if err != nil {
		t.Error(err)
	}

	result, err := c.Update()
	if err != nil {
		t.Error(err)
	}

	if result.Numregs != 10 {
		t.Errorf("got %d, want %d", result.Numregs, 10)
	}
}
