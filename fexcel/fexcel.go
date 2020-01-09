package fexcel

import (
	"fmt"
	"strings"

	"github.com/onerobotics/fexcel/excel"
	"github.com/onerobotics/fexcel/fanuc"
)

const Version = "2.0.0"

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

func Logo() string {
	return fmt.Sprintf(`  __                  _
 / _|                | |
| |_ _____  _____ ___| |
|  _/ _ \ \/ / __/ _ \ |
| ||  __/>  < (_|  __/ |
|_| \___/_/\_\___\___|_|
                  v%s

by ONE Robotics Company
www.onerobotics.com

`, Version)
}

func Pluralize(word string, i int) string {
	if i == 1 {
		return word
	} else {
		return word + "s"
	}
}

func parseLocationSpec(spec string, defaultSheet string) (sheet string, axis string, err error) {
	if spec == "" {
		return
	}

	parts := strings.Split(spec, ":")

	switch len(parts) {
	case 2:
		sheet, axis = parts[0], parts[1]
		return
	case 1:
		// e.g. A2
		sheet = defaultSheet
		axis = spec
		return
	}

	err = fmt.Errorf("Cell specification %q is invalid. Should be in the form [Sheet:]Cell e.g. Sheet1:A2 or just A2.", spec)

	return
}

func setLocation(f *excel.File, dataType fanuc.DataType, spec string, defaultSheet string) error {
	sheet, axis, err := parseLocationSpec(spec, defaultSheet)
	if err != nil {
		return err
	}

	f.SetLocation(dataType, axis, sheet)

	return nil
}

func PrepareFile(fpath string, cfg Config) (*excel.File, error) {
	f, err := excel.NewFile(fpath, cfg.Offset)
	if err != nil {
		return nil, err
	}

	typeSpecs := []struct {
		dataType fanuc.DataType
		spec     string
	}{
		{fanuc.Numreg, cfg.Numregs},
		{fanuc.Posreg, cfg.Posregs},
		{fanuc.Ualm, cfg.Ualms},
		{fanuc.Rin, cfg.Rins},
		{fanuc.Rout, cfg.Routs},
		{fanuc.Din, cfg.Dins},
		{fanuc.Dout, cfg.Douts},
		{fanuc.Gin, cfg.Gins},
		{fanuc.Gout, cfg.Gouts},
		{fanuc.Ain, cfg.Ains},
		{fanuc.Aout, cfg.Aouts},
		{fanuc.Sreg, cfg.Sregs},
		{fanuc.Flag, cfg.Flags},
	}

	for _, ts := range typeSpecs {
		if ts.spec != "" {
			err = setLocation(f, ts.dataType, ts.spec, cfg.Sheet)
			if err != nil {
				return nil, err
			}
		}
	}

	return f, nil
}
