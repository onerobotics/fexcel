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
