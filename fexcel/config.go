package fexcel

import (
	"errors"

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
