package fexcel

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Location struct {
	Type
	Axis   string // e.g. A2
	Sheet  string
	Offset int
}

var cellRegexp = regexp.MustCompile(`(\w+)(\{(\d+)\})?`)

// parses the Cell Axis and (optional) offset from the end
// of a cell spec (e.g. cell or cell{offset})
//
func parseCell(spec string) (cell string, offset int, err error) {
	if !strings.Contains(spec, "{") || !strings.Contains(spec, "}") {
		return spec, 0, nil
	}

	matches := cellRegexp.FindAllStringSubmatch(spec, -1)
	if len(matches) != 1 {
		return "", 0, fmt.Errorf("%q does not appear to be a valid cell spe", spec)
	}

	match := matches[0]
	if len(match) != 4 {
		panic("invalid cell match data")
	}

	// no offset specified
	if match[3] == "" {
		return spec, 0, nil
	}

	i, err := strconv.Atoi(match[3])
	if err != nil {
		return "", 0, err
	}

	return match[1], i, nil
}

// returns a Location based on a cell specification
// spec can be in the following forms:
//
//   Sheet:Cell
//   Sheet:Cell{offset}
//         Cell
//         Cell{offset}
//
// if the sheet is not provided in the spec, the default
// sheet is used.
//
func NewLocation(t Type, spec string, defaultSheet string) (*Location, error) {
	parts := strings.Split(spec, ":")

	switch len(parts) {
	case 2:
		axis, offset, err := parseCell(parts[1])
		if err != nil {
			return nil, err
		}

		return &Location{Type: t, Sheet: parts[0], Axis: axis, Offset: offset}, nil
	case 1:
		// e.g. A2 or A2{5}
		if defaultSheet == "" {
			return nil, fmt.Errorf("cell specification %q requires a default sheet, but none has been defined", spec)
		}

		axis, offset, err := parseCell(parts[0])
		if err != nil {
			return nil, err
		}

		return &Location{Type: t, Sheet: defaultSheet, Axis: axis, Offset: offset}, nil
	}

	return nil, fmt.Errorf("Cell specification %q is invalid. Should be in the form [Sheet:]Cell e.g. Sheet1:A2 or just A2.", spec)
}

type Definition struct {
	Type    Type
	Sheet   string
	Column  int
	Row     int
	Id      int
	Comment string
}

type File struct {
	path string
	xlsx *excelize.File

	Config    FileConfig
	Locations []*Location
	Warnings  []string
}

func (f *File) LocationsFor(t Type) []*Location {
	var locations []*Location
	for _, l := range f.Locations {
		if l.Type == t {
			locations = append(locations, l)
		}
	}

	return locations
}

func newFile(path string, cfg FileConfig) (*File, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	f := File{path: path, Config: cfg}

	locations, err := cfg.Locations()
	if err != nil {
		return nil, err
	}

	for _, locs := range locations {
		f.Locations = append(f.Locations, locs...)
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
	d.Sheet = sheet
	d.Column = col
	d.Row = row

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

	for _, l := range f.Locations {
		if l.Type == Constant {
			continue
		}

		d, err := f.Definitions(l.Type)
		if err != nil {
			return nil, err
		}

		defs[l.Type] = d
	}

	return defs, nil
}

func (f *File) Definitions(t Type) ([]Definition, error) {
	if f.Config.SpecFor(t) == "" {
		return nil, fmt.Errorf("Location for %s not defined", t)
	}

	ids := make(map[int]bool)
	var defs []Definition
	for _, loc := range f.Locations {
		if loc.Type != t {
			continue
		}

		col, row, err := excelize.CellNameToCoordinates(loc.Axis)
		if err != nil {
			return nil, fmt.Errorf("Invalid location for %s: %q", t, loc.Axis)
		}

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

			if ids[d.Id] {
				return nil, fmt.Errorf("%s[%d] already defined", t, d.Id)
			}
			ids[d.Id] = true

			defs = append(defs, d)

		}
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
	if f.Config.SpecFor(Constant) == "" {
		return nil, fmt.Errorf("Location for %s not defined", Constant)
	}

	constants := make(map[string]string)

	for _, loc := range f.Locations {
		if loc.Type != Constant {
			continue
		}

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

			value, err := f.readString(loc.Sheet, col+offset, row)
			if err != nil {
				return nil, err
			}
			if value == "" {
				f.Warnings = append(f.Warnings, fmt.Sprintf("Definition for constant %q is blank", id))
				continue
			}

			constants[id] = value
		}

	}

	return constants, nil
}
