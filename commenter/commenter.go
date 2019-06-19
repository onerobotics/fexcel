package commenter

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/onerobotics/comtool"
)

const (
	MaxDataLength = 16
	MaxIOLength   = 24
	MaxUalmLength = 29
)

type Config struct {
	Numregs string // e.g. A2 or SheetName:A2
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
	Timeout int // milliseconds
}

type commenter struct {
	Config
	DefaultSheetName string
	Offset           int

	xlsx *excelize.File
}

func New(filename string, defaultSheet string, offset int, cfg Config) (*commenter, error) {
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	index := xlsx.GetSheetIndex(defaultSheet)
	if index == 0 {
		return nil, fmt.Errorf("Could not find sheet with name '%s' in file: %s\n", defaultSheet, filename)
	}

	c := &commenter{
		Config:           cfg,
		DefaultSheetName: defaultSheet,
		Offset:           offset,
	}
	c.xlsx = xlsx

	return c, nil
}

func (c *commenter) getIdAndComment(sheetName string, col, row int) (id int, comment string, err error) {
	axis, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return
	}

	value, err := c.xlsx.GetCellValue(sheetName, axis)
	if err != nil {
		return
	}

	id, err = strconv.Atoi(value)
	if err != nil {
		return
	}

	axis, err = excelize.CoordinatesToCellName(col+c.Offset, row)
	if err != nil {
		return
	}

	comment, err = c.xlsx.GetCellValue(sheetName, axis)
	if err != nil {
		return
	}

	return
}

func (c *commenter) Update(host string) error {
	fmt.Printf("Updating %s...\n", host)
	fmt.Println("----------------------------")

	count, err := c.processColumn(c.Config.Numregs, comtool.NUMREG, host, MaxDataLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d numeric registers\n", count)

	count, err = c.processColumn(c.Config.Posregs, comtool.POSREG, host, MaxDataLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d position registers\n", count)

	count, err = c.processColumn(c.Config.Ualms, comtool.UALM, host, MaxUalmLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d user alarms\n", count)

	count, err = c.processColumn(c.Config.Rins, comtool.RIN, host, MaxIOLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d robot inputs\n", count)

	count, err = c.processColumn(c.Config.Routs, comtool.ROUT, host, MaxIOLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d robot outputs\n", count)

	count, err = c.processColumn(c.Config.Dins, comtool.DIN, host, MaxIOLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d digital inputs\n", count)

	count, err = c.processColumn(c.Config.Douts, comtool.DOUT, host, MaxIOLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d digital outputs\n", count)

	count, err = c.processColumn(c.Config.Gins, comtool.GIN, host, MaxIOLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d group inputs\n", count)

	count, err = c.processColumn(c.Config.Gouts, comtool.GOUT, host, MaxIOLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d group outputs\n", count)

	count, err = c.processColumn(c.Config.Ains, comtool.AIN, host, MaxIOLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d analog inputs\n", count)

	count, err = c.processColumn(c.Config.Aouts, comtool.AOUT, host, MaxIOLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d analog outputs\n", count)

	count, err = c.processColumn(c.Config.Sregs, comtool.SREG, host, MaxDataLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d string registers\n", count)

	count, err = c.processColumn(c.Config.Flags, comtool.FLAG, host, MaxIOLength)
	if err != nil {
		return err
	}
	fmt.Printf("Updated %d flags\n", count)

	return nil
}

func (c *commenter) warn(format string, args ...interface{}) {
	fmt.Printf("WARNING: "+format, args...)
}

// cellString may be A2 or SheetName:A2
func (c *commenter) parseStartCell(cellString string) (sheetName string, cellName string, err error) {
	parts := strings.Split(cellString, ":")
	if len(parts) > 2 {
		err = fmt.Errorf("Invalid cell string: `%s`", cellString)
		return
	}

	if len(parts) == 1 {
		sheetName = c.DefaultSheetName
		cellName = cellString
		return
	}

	sheetName = parts[0]
	cellName = parts[1]

	return
}

func (c *commenter) processColumn(startCell string, fCode comtool.FunctionCode, host string, maxLength int) (count int, err error) {
	if startCell == "" {
		return
	}

	sheetName, cellName, err := c.parseStartCell(startCell)
	if err != nil {
		return
	}

	col, row, err := excelize.CellNameToCoordinates(cellName)
	if err != nil {
		return
	}

	for {
		id, comment, err := c.getIdAndComment(sheetName, col, row)
		if err != nil {
			break
		}

		comment = strings.TrimSpace(comment)
		if len(comment) > maxLength {
			cellName, err := excelize.CoordinatesToCellName(col+c.Offset, row)
			if err != nil {
				return count, err
			}

			oldComment := comment
			comment := comment[:maxLength]

			c.warn("Comment in cell %s truncated from '%s' to '%s'. (%d > %d)\n", cellName, oldComment, comment, len(oldComment), maxLength)
		}

		err = comtool.Set(fCode, id, comment, host, time.Duration(time.Duration(c.Config.Timeout)*time.Millisecond))
		if err != nil {
			return count, err
		}

		count++
		row++
	}

	return
}
