package commenter

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
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
	Hosts            []string

	xlsx *excelize.File
}

func New(filename string, defaultSheet string, offset int, cfg Config, hosts []string) (*commenter, error) {
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
		Hosts:            hosts,
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

type Result struct {
	Flags   int
	Numregs int
	Posregs int
	Sregs   int
	Ualms   int

	Ains  int
	Aouts int
	Dins  int
	Douts int
	Gins  int
	Gouts int
	Rins  int
	Routs int
}

func (c *commenter) Update() (result Result, err error) {
	result.Numregs, err = c.processColumn(c.Config.Numregs, comtool.NUMREG, MaxDataLength)
	if err != nil {
		return
	}

	result.Posregs, err = c.processColumn(c.Config.Posregs, comtool.POSREG, MaxDataLength)
	if err != nil {
		return
	}

	result.Ualms, err = c.processColumn(c.Config.Ualms, comtool.UALM, MaxUalmLength)
	if err != nil {
		return
	}

	result.Rins, err = c.processColumn(c.Config.Rins, comtool.RIN, MaxIOLength)
	if err != nil {
		return
	}

	result.Routs, err = c.processColumn(c.Config.Routs, comtool.ROUT, MaxIOLength)
	if err != nil {
		return
	}

	result.Dins, err = c.processColumn(c.Config.Dins, comtool.DIN, MaxIOLength)
	if err != nil {
		return
	}

	result.Douts, err = c.processColumn(c.Config.Douts, comtool.DOUT, MaxIOLength)
	if err != nil {
		return
	}

	result.Gins, err = c.processColumn(c.Config.Gins, comtool.GIN, MaxIOLength)
	if err != nil {
		return
	}

	result.Gouts, err = c.processColumn(c.Config.Gouts, comtool.GOUT, MaxIOLength)
	if err != nil {
		return
	}

	result.Ains, err = c.processColumn(c.Config.Ains, comtool.AIN, MaxIOLength)
	if err != nil {
		return
	}

	result.Aouts, err = c.processColumn(c.Config.Aouts, comtool.AOUT, MaxIOLength)
	if err != nil {
		return
	}

	result.Sregs, err = c.processColumn(c.Config.Sregs, comtool.SREG, MaxDataLength)
	if err != nil {
		return
	}

	result.Flags, err = c.processColumn(c.Config.Flags, comtool.FLAG, MaxIOLength)
	if err != nil {
		return
	}

	return
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

func (c *commenter) truncateComment(comment string, maxLength int, row int, col int) (string, error) {
	if len(comment) > maxLength {
		cellName, err := excelize.CoordinatesToCellName(col+c.Offset, row)
		if err != nil {
			return "", err
		}

		oldComment := comment
		comment := comment[:maxLength]

		c.warn("Comment in cell %s truncated from '%s' to '%s'. (%d > %d)\n", cellName, oldComment, comment, len(oldComment), maxLength)
	}

	return comment, nil
}

func (c *commenter) processColumn(startCell string, fCode comtool.FunctionCode, maxLength int) (count int, err error) {
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

		comment, err = c.truncateComment(strings.TrimSpace(comment), maxLength, row, col)
		if err != nil {
			return count, err
		}

		var wg sync.WaitGroup
		errors := make(chan error, len(c.Hosts))
		for _, host := range c.Hosts {
			wg.Add(1)
			go func(host string, wg *sync.WaitGroup) {
				defer wg.Done()

				err := comtool.Set(fCode, id, comment, host, time.Duration(time.Duration(c.Config.Timeout)*time.Millisecond))
				if err != nil {
					errors <- fmt.Errorf("%s: %s", host, err)
				}
			}(host, &wg)
		}
		wg.Wait()
		close(errors)

		if len(errors) > 0 {
			s := ""
			for err := range errors {
				s = fmt.Sprintf("%s\n", err)
			}
			return count, fmt.Errorf("%s", s)
		}

		count++
		row++
	}

	return
}
