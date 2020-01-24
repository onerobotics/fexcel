package fexcel

import (
	"errors"
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	fanuc "github.com/onerobotics/go-fanuc"
)

type FileConfig struct {
	Numregs string // e.g. A2 or Sheet1:A2
	Posregs string
	Ualms   string
	Rins    string
	Routs   string
	Dins    string
	Douts   string
	Gins    string
	Gouts   string
	Ains    string
	Aouts   string
	Sregs   string
	Flags   string
	Sheet   string
	Offset  int
}

type Config struct {
	FileConfig
	NoUpdate bool
	Timeout  int
}

func (c *FileConfig) Specs() []string {
	return []string{c.Numregs, c.Posregs, c.Ualms, c.Rins, c.Routs, c.Dins, c.Douts, c.Gins, c.Gouts, c.Ains, c.Aouts, c.Sregs, c.Flags}
}

func (c *FileConfig) Count() (i int) {
	for _, spec := range c.Specs() {
		if spec != "" {
			i++
		}
	}

	return i
}

func (c *FileConfig) Locations() (map[fanuc.Type][]*Location, error) {
	types := []fanuc.Type{fanuc.Numreg, fanuc.Posreg, fanuc.Ualm, fanuc.Rin, fanuc.Rout, fanuc.Din, fanuc.Dout, fanuc.Gin, fanuc.Gout, fanuc.Ain, fanuc.Aout, fanuc.Sreg, fanuc.Flag}

	locations := make(map[fanuc.Type][]*Location)
	for _, t := range types {
		spec := c.SpecFor(t)
		if spec != "" {
			l, err := NewLocation(spec, c.Sheet)
			if err != nil {
				return nil, err
			}

			locations[t] = append(locations[t], l)
		}
	}

	return locations, nil
}

func (c *FileConfig) CheckHeaders() error {
	locations, err := c.Locations()
	if err != nil {
		return nil
	}

	for t, locs := range locations {
		for _, loc := range locs {
			_, row, err := excelize.CellNameToCoordinates(loc.Axis)
			if err != nil {
				return err
			}

			if row < 2 {
				return fmt.Errorf("Cell spec for %ss (%s) must be in row 2 or lower for headers option", t.VerboseName(), c.SpecFor(t))
			}
		}
	}

	return nil
}

func (c *FileConfig) HasOverlaps() (bool, error) {
	locations, err := c.Locations()
	if err != nil {
		return false, err
	}

	sheets := make(map[string]map[int]bool)
	for _, locs := range locations {
		for _, loc := range locs {
			if _, defined := sheets[loc.Sheet]; !defined {
				sheets[loc.Sheet] = make(map[int]bool)
			}

			col, _, err := excelize.CellNameToCoordinates(loc.Axis)
			if err != nil {
				return false, err
			}

			// we consider the start Axis all the way through the offset to be an overlap
			// e.g. numregs starting in column A with an offset of 5 will
			// prevent other items from using columns A, B, C, D and E
			for i := col; i <= col+c.Offset; i++ {
				if sheets[loc.Sheet][i] {
					return true, nil
				} else {
					sheets[loc.Sheet][i] = true
				}
			}
		}
	}

	return false, nil
}

func (c *FileConfig) Validate() error {
	if c.Count() < 1 {
		return errors.New("no cell locations defined")
	}

	if c.Offset == 0 {
		return errors.New("offset must be nonzero")
	}

	return nil
}

func (c *FileConfig) SpecFor(t fanuc.Type) string {
	switch t {
	case fanuc.Numreg:
		return c.Numregs
	case fanuc.Posreg:
		return c.Posregs
	case fanuc.Ualm:
		return c.Ualms
	case fanuc.Rin:
		return c.Rins
	case fanuc.Rout:
		return c.Routs
	case fanuc.Din:
		return c.Dins
	case fanuc.Dout:
		return c.Douts
	case fanuc.Gin:
		return c.Gins
	case fanuc.Gout:
		return c.Gouts
	case fanuc.Ain:
		return c.Ains
	case fanuc.Aout:
		return c.Aouts
	case fanuc.Sreg:
		return c.Sregs
	case fanuc.Flag:
		return c.Flags
	}

	return ""
}
