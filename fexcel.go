package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/onerobotics/fexcel/commenter"
)

var (
	sheetName string
	offset    int
	cfg commenter.Config
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
	fmt.Fprintf(os.Stderr, "Usage: fexcel [options] filename host(s)...\n\n")

	fmt.Fprintf(os.Stderr, "Example: fexcel -sheet Data -numregs A2 -posregs D2 spreadsheet.xlsx 127.0.0.101 127.0.0.102\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func init() {
	flag.StringVar(&sheetName, "sheet", "Sheet1", "the name of the sheet")
	flag.IntVar(&offset, "offset", 1, "column offset from ids to comments")
	flag.StringVar(&cfg.Numregs, "numregs", "", "start cell of numeric register ids")
	flag.StringVar(&cfg.Posregs, "posregs", "", "start cell of position register ids")
	flag.StringVar(&cfg.Ualms, "ualms", "", "start cell of user alarm ids")
	flag.StringVar(&cfg.Rins, "rins", "", "start cell of robot input ids")
	flag.StringVar(&cfg.Routs, "routs", "", "start cell of robot output ids")
	flag.StringVar(&cfg.Dins, "dins", "", "start cell of digital input ids")
	flag.StringVar(&cfg.Douts, "douts", "", "start cell of digital output ids")
	flag.StringVar(&cfg.Gins, "gins", "", "start cell of group input ids")
	flag.StringVar(&cfg.Gouts, "gouts", "", "start cell of group output ids")
	flag.StringVar(&cfg.Ains, "ains", "", "start cell of analog input ids")
	flag.StringVar(&cfg.Aouts, "aouts", "", "start cell of analog output ids")
	flag.StringVar(&cfg.Sregs, "sregs", "", "start cell of string register ids")
	flag.StringVar(&cfg.Flags, "flags", "", "start cell of flag ids")
}

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		usage()
	}

	filename := args[0]
	if filename == "" {
		usage()
	}

	hosts := args[1:]

	fmt.Printf(logo)

	c, err := commenter.New(filename, sheetName, offset, cfg)
	check(err)

	for _, host := range hosts {
		err = c.Update(host)
		check(err)

		fmt.Println()
	}
}
