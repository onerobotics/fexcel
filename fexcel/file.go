package fexcel

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Location struct {
	Axis   string // e.g. A2
	Sheet  string
	Offset int
}

// returns a Location based on a cell specification
// spec can be in the following forms:
//
//   Offset:Sheet:Cell
//          Sheet:Cell
//                Cell
//
// if the sheet is not provided in the spec, the default
// sheet is used.
//
func NewLocation(spec string, defaultSheet string) (*Location, error) {
	parts := strings.Split(spec, ":")

	switch len(parts) {
	case 3:
		offset, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		return &Location{Sheet: parts[1], Axis: parts[2], Offset: offset}, nil
	case 2:
		return &Location{Sheet: parts[0], Axis: parts[1]}, nil
	case 1:
		// e.g. A2
		if defaultSheet == "" {
			return nil, fmt.Errorf("cell specification %q requires a default sheet, but none has been defined", spec)
		}
		return &Location{Sheet: defaultSheet, Axis: spec}, nil
	}

	return nil, fmt.Errorf("Cell specification %q is invalid. Should be in the form [Sheet:]Cell e.g. Sheet1:A2 or just A2.", spec)
}

type Definition struct {
	Type    Type
	Id      int
	Comment string
}

type File struct {
	path string
	xlsx *excelize.File

	Config    FileConfig
	Locations map[Type]*Location
	Warnings  []string
}

func newFile(path string, cfg FileConfig) (*File, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	f := File{path: path, Config: cfg}

	// set locations based on config
	f.Locations = make(map[Type]*Location)
	types := []Type{Constant, Numreg, Posreg, Ualm, Ain, Aout, Din, Dout, Gin, Gout, Rin, Rout, Sreg, Flag}
	for _, t := range types {
		spec := cfg.SpecFor(t)
		if spec != "" {
			loc, err := NewLocation(spec, cfg.Sheet)
			if err != nil {
				return nil, err
			}

			f.Locations[t] = loc
		}
	}

	return &f, nil
}

func OpenFile(path string, cfg FileConfig) (*File, error) {
	f, err := newFile(path, cfg)
	if err != nil {
		return nil, err
	}

	xlsx, err := excelize.OpenFile(f.path)
	if err != nil {
		return nil, err
	}

	f.xlsx = xlsx

	return f, nil
}

func NewFile(path string, cfg FileConfig) (*File, error) {
	f, err := newFile(path, cfg)
	if err != nil {
		return nil, err
	}

	// file must not exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f.xlsx = excelize.NewFile()

		return f, nil
	}

	return nil, fmt.Errorf("File %q already exists", path)
}

func (f *File) Save() error {
	return f.xlsx.SaveAs(f.path)
}

func (f *File) readInt(sheet string, col, row int) (int, error) {
	axis, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return 0, err
	}

	value, err := f.xlsx.GetCellValue(sheet, axis)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (f *File) readString(sheet string, col, row int) (string, error) {
	axis, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return "", err
	}

	value, err := f.xlsx.GetCellValue(sheet, axis)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (f *File) readDefinition(t Type, sheet string, col, row, offset int) (d Definition, err error) {
	d.Type = t

	d.Id, err = f.readInt(sheet, col, row)
	if err != nil {
		return
	}

	d.Comment, err = f.readString(sheet, col+offset, row)
	if maxLength := MaxLengthFor(t); len(d.Comment) > maxLength {
		var axis string
		axis, err = excelize.CoordinatesToCellName(col+offset, row)
		if err != nil {
			return
		}

		f.Warnings = append(f.Warnings, fmt.Sprintf("comment in [%s]%s for %s[%d] will be truncated to %q (length %d > max length %d for %ss)", sheet, axis, d.Type, d.Id, d.Comment[:maxLength], len(d.Comment), maxLength, t))
	}

	return
}

func (f *File) AllDefinitions() (map[Type][]Definition, error) {
	defs := make(map[Type][]Definition)

	for t, _ := range f.Locations {
		if t == Constant {
			continue
		}

		d, err := f.Definitions(t)
		if err != nil {
			return nil, err
		}

		defs[t] = d
	}

	return defs, nil
}

func (f *File) Definitions(t Type) ([]Definition, error) {
	loc, defined := f.Locations[t]
	if !defined {
		return nil, fmt.Errorf("Location for %s not defined", t)
	}

	col, row, err := excelize.CellNameToCoordinates(loc.Axis)
	if err != nil {
		return nil, fmt.Errorf("Invalid location for %s: %q", t, loc.Axis)
	}

	var defs []Definition
	for ; ; row++ {
		// check for blank id
		s, err := f.readString(loc.Sheet, col, row)
		if err != nil {
			return nil, err
		}
		if s == "" {
			break
		}

		offset := loc.Offset
		if offset == 0 {
			offset = f.Config.Offset
		}

		d, err := f.readDefinition(t, loc.Sheet, col, row, offset)
		if err != nil {
			return nil, err
		}

		defs = append(defs, d)

	}

	return defs, nil
}

func (f *File) SetValue(sheet string, col int, row int, value interface{}) error {
	axis, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return err
	}

	return f.xlsx.SetCellValue(sheet, axis, value)
}

// excelize does not create a new sheet if it already exists
func (f *File) CreateSheet(name string) {
	f.xlsx.NewSheet(name)
}

func (f *File) Constants() (map[string]string, error) {
	loc, defined := f.Locations[Constant]
	if !defined {
		return nil, fmt.Errorf("Location for %s not defined", Constant)
	}

	constants := make(map[string]string)

	col, row, err := excelize.CellNameToCoordinates(loc.Axis)
	if err != nil {
		return nil, fmt.Errorf("Invalid location for %s: %q", Constant, loc.Axis)
	}

	for ; ; row++ {
		// check for blank identifier
		id, err := f.readString(loc.Sheet, col, row)
		if err != nil {
			return nil, err
		}
		if id == "" {
			break
		}

		offset := loc.Offset
		if offset == 0 {
			offset = f.Config.Offset
		}

		value, err := f.readString(loc.Sheet, col+1, row)
		if err != nil {
			return nil, err
		}
		if value == "" {
			f.Warnings = append(f.Warnings, fmt.Sprintf("Definition for constant %q is blank", id))
			continue
		}

		constants[id] = value
	}

	return constants, nil
}
