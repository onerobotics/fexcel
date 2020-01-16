package fexcel

import (
	"errors"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	fanuc "github.com/onerobotics/go-fanuc"
)

type Config struct {
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

	Sheet    string
	Offset   int
	NoUpdate bool
}

func (c *Config) Count() (i int) {
	items := []string{c.Numregs, c.Posregs, c.Ualms, c.Rins, c.Routs, c.Dins, c.Douts, c.Gins, c.Gouts, c.Ains, c.Aouts, c.Sregs, c.Flags}

	for _, item := range items {
		if item != "" {
			i++
		}
	}

	return i
}

func (c *Config) HasOverlaps() (bool, error) {
	specs := []string{c.Numregs, c.Posregs, c.Ualms, c.Rins, c.Routs, c.Dins, c.Douts, c.Gins, c.Gouts, c.Ains, c.Aouts, c.Sregs, c.Flags}

	var locations []*Location
	for _, spec := range specs {
		if spec != "" {
			l, err := NewLocation(spec, c.Sheet)
			if err != nil {
				return false, err
			}

			locations = append(locations, l)
		}
	}

	sheets := make(map[string]map[int]bool)
	for _, loc := range locations {
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

	return false, nil
}

func (c *Config) Validate() error {
	if c.Count() < 1 {
		return errors.New("no cell locations defined")
	}

	if c.Offset == 0 {
		return errors.New("offset must be nonzero")
	}

	return nil
}

func (c *Config) SpecFor(t fanuc.Type) string {
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
