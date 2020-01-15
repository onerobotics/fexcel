package fexcel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	fanuc "github.com/onerobotics/go-fanuc"
)

type Location struct {
	Axis  string // e.g. A2
	Sheet string
}

// returns a Location based on a cell specification
// spec can be in the form of Sheet:Cell or just Cell
// if the sheet is not provided in the spec, the default
// sheet is used.
func NewLocation(spec string, defaultSheet string) (*Location, error) {
	parts := strings.Split(spec, ":")

	switch len(parts) {
	case 2:
		return &Location{Sheet: parts[0], Axis: parts[1]}, nil
	case 1:
		// e.g. A2
		return &Location{Sheet: defaultSheet, Axis: spec}, nil
	}

	return nil, fmt.Errorf("Cell specification %q is invalid. Should be in the form [Sheet:]Cell e.g. Sheet1:A2 or just A2.", spec)
}

type File struct {
	path string
	xlsx *excelize.File

	Config
	Locations map[fanuc.Type]*Location
}

func NewFile(path string, cfg Config) (*File, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	f := File{path: path, Config: cfg}

	// set locations based on config
	f.Locations = make(map[fanuc.Type]*Location)
	types := []fanuc.Type{fanuc.Numreg, fanuc.Posreg, fanuc.Ualm, fanuc.Ain, fanuc.Aout, fanuc.Din, fanuc.Dout, fanuc.Gin, fanuc.Gout, fanuc.Rin, fanuc.Rout, fanuc.Sreg, fanuc.Flag}
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

func (f *File) Open() error {
	xlsx, err := excelize.OpenFile(f.path)
	if err != nil {
		return err
	}

	f.xlsx = xlsx

	return nil
}

func (f *File) New() {
	f.xlsx = excelize.NewFile()
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

func (f *File) readDefinition(dataType fanuc.Type, sheet string, col, row int) (d Definition, err error) {
	d.Type = dataType

	d.Id, err = f.readInt(sheet, col, row)
	if err != nil {
		return
	}

	d.Comment, err = f.readString(sheet, col+f.Config.Offset, row)
	return
}

func (f *File) Definitions(dataType fanuc.Type) ([]Definition, error) {
	loc, defined := f.Locations[dataType]
	if !defined {
		return nil, fmt.Errorf("Location for %s not defined", dataType)
	}

	col, row, err := excelize.CellNameToCoordinates(loc.Axis)
	if err != nil {
		return nil, fmt.Errorf("Invalid location for %s: %q", dataType, loc.Axis)
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

		d, err := f.readDefinition(dataType, loc.Sheet, col, row)
		if err != nil {
			return nil, err
		}

		defs = append(defs, d)

	}

	return defs, nil
}
