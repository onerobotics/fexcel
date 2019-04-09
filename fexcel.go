package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/onerobotics/comtool"
)

var (
	sheetName string
	offset    int
	numregs   string
	posregs   string
	ualms     string
	rins      string
	routs     string
	dins      string
	douts     string
	gins      string
	gouts     string
	ains      string
	aouts     string
	sregs     string
	flags     string
)

const logo = `  __                  _
 / _|                | |
| |_ _____  _____ ___| |
|  _/ _ \ \/ / __/ _ \ |
| ||  __/>  < (_|  __/ |
|_| \___/_/\_\___\___|_|

by ONE Robotics Company
www.onerobotics.com

`

func usage() {
	fmt.Fprintf(os.Stderr, logo)
	fmt.Fprintf(os.Stderr, "Usage: fexcel [options] filename host\n\n")
	fmt.Fprintf(os.Stderr, "Example: fexcel -sheet Data -numregs A2 -posregs D2 spreadsheet.xlsx 127.0.0.101\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func init() {
	flag.StringVar(&sheetName, "sheet", "Sheet1", "the name of the sheet")
	flag.IntVar(&offset, "offset", 1, "column offset from ids to comments")
	flag.StringVar(&numregs, "numregs", "", "start cell of numeric register ids")
	flag.StringVar(&posregs, "posregs", "", "start cell of position register ids")
	flag.StringVar(&ualms, "ualms", "", "start cell of user alarm ids")
	flag.StringVar(&rins, "rins", "", "start cell of robot input ids")
	flag.StringVar(&routs, "routs", "", "start cell of robot output ids")
	flag.StringVar(&dins, "dins", "", "start cell of digital input ids")
	flag.StringVar(&douts, "douts", "", "start cell of digital output ids")
	flag.StringVar(&gins, "gins", "", "start cell of group input ids")
	flag.StringVar(&gouts, "gouts", "", "start cell of group output ids")
	flag.StringVar(&ains, "ains", "", "start cell of analog input ids")
	flag.StringVar(&aouts, "aouts", "", "start cell of analog output ids")
	flag.StringVar(&sregs, "sregs", "", "start cell of string register ids")
	flag.StringVar(&flags, "flags", "", "start cell of flag ids")
}

func getIdAndComment(col, row int, xlsx *excelize.File) (id int, comment string, err error) {
	axis, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return
	}

	value, err := xlsx.GetCellValue(sheetName, axis)
	if err != nil {
		return
	}

	id, err = strconv.Atoi(value)
	if err != nil {
		return
	}

	axis, err = excelize.CoordinatesToCellName(col+offset, row)
	if err != nil {
		return
	}

	comment, err = xlsx.GetCellValue(sheetName, axis)
	if err != nil {
		return
	}

	return
}

func setComments(startCell string, xlsx *excelize.File, f comtool.FunctionCode, host string) (count int, err error) {
	if startCell == "" {
		return
	}

	col, row, err := excelize.CellNameToCoordinates(startCell)
	if err != nil {
		return
	}

	for {
		id, comment, err := getIdAndComment(col, row, xlsx)
		if err != nil {
			break
		}

		err = comtool.Set(f, id, comment, host)
		if err != nil {
			return count, err
		}

		count++
		row++
	}

	return
}

func main() {
	flag.Parse()

	filename := flag.Arg(0)
	if filename == "" {
		usage()
	}

	host := flag.Arg(1)
	if host == "" {
		usage()
	}

	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf(logo)

	count, err := setComments(numregs, xlsx, comtool.NUMREG, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d numeric registers\n", count)

	count, err = setComments(posregs, xlsx, comtool.POSREG, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d position registers\n", count)

	count, err = setComments(ualms, xlsx, comtool.UALM, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d user alarms\n", count)

	count, err = setComments(rins, xlsx, comtool.RIN, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d robot inputs\n", count)

	count, err = setComments(routs, xlsx, comtool.ROUT, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d robot outputs\n", count)

	count, err = setComments(dins, xlsx, comtool.DIN, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d digital inputs\n", count)

	count, err = setComments(douts, xlsx, comtool.DOUT, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d digital outputs\n", count)

	count, err = setComments(gins, xlsx, comtool.GIN, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d group inputs\n", count)

	count, err = setComments(gouts, xlsx, comtool.GOUT, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d group outputs\n", count)

	count, err = setComments(ains, xlsx, comtool.AIN, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d analog inputs\n", count)

	count, err = setComments(aouts, xlsx, comtool.AOUT, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d analog outputs\n", count)

	count, err = setComments(sregs, xlsx, comtool.SREG, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d string registers\n", count)

	count, err = setComments(flags, xlsx, comtool.FLAG, host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Updated %d flags\n", count)
}
