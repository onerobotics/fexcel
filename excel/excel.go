package excel

import (
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/onerobotics/fexcel/fanuc"
)

type Location struct {
	Axis  string // e.g. A2
	Sheet string
}

type File struct {
	offset int
	xlsx   *excelize.File

	Locations map[fanuc.DataType]Location
}

func NewFile(path string, offset int) (*File, error) {
	xlsx, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	if offset == 0 {
		return nil, fmt.Errorf("offset must be nonzero")
	}

	f := File{offset: offset, xlsx: xlsx}
	f.Locations = make(map[fanuc.DataType]Location)

	return &f, nil
}

func (f *File) SetLocation(dataType fanuc.DataType, axis, sheet string) {
	f.Locations[dataType] = Location{axis, sheet}
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

func (f *File) readDefinition(dataType fanuc.DataType, sheet string, col, row int) (d fanuc.Definition, err error) {
	d.DataType = dataType

	d.Id, err = f.readInt(sheet, col, row)
	if err != nil {
		return
	}

	d.Comment, err = f.readString(sheet, col+f.offset, row)
	return
}

func (f *File) Definitions(dataType fanuc.DataType) ([]fanuc.Definition, error) {
	loc, defined := f.Locations[dataType]
	if !defined {
		return nil, fmt.Errorf("Location for %s not defined", dataType)
	}

	col, row, err := excelize.CellNameToCoordinates(loc.Axis)
	if err != nil {
		return nil, fmt.Errorf("Invalid location for %s: %q", dataType, loc.Axis)
	}

	var defs []fanuc.Definition
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
