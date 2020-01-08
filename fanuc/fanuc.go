package fanuc

import (
	"strconv"
)

type DataType int

const (
	Numreg DataType = iota
	Posreg
	Ualm
	Rin
	Rout
	Din
	Dout
	Gin
	Gout
	Ain
	Aout
	Sreg
	Flag
)

var dataTypes = [...]string{
	Numreg: "R",
	Posreg: "PR",
	Ualm:   "UALM",
	Rin:    "RI",
	Rout:   "RO",
	Din:    "DI",
	Dout:   "DO",
	Gin:    "GI",
	Gout:   "GO",
	Ain:    "AI",
	Aout:   "AO",
	Sreg:   "SR",
	Flag:   "F",
}

var verboseNames = [...]string{
	Numreg: "Numeric Register",
	Posreg: "Position Register",
	Ualm:   "User Alarm",
	Rin:    "Robot Input",
	Rout:   "Robot Output",
	Din:    "Digital Input",
	Dout:   "Digital Output",
	Gin:    "Group Input",
	Gout:   "Group Output",
	Ain:    "Analog Input",
	Aout:   "Analog Output",
	Sreg:   "String Register",
	Flag:   "Flag",
}

func (d DataType) String() string {
	s := ""
	if 0 <= d && d < DataType(len(dataTypes)) {
		s = dataTypes[d]
	}
	if s == "" {
		s = "DataType(" + strconv.Itoa(int(d)) + ")"
	}
	return s
}

func (d DataType) VerboseName() string {
	s := ""
	if 0 <= d && d < DataType(len(verboseNames)) {
		s = verboseNames[d]
	}
	if s == "" {
		s = "DataType(" + strconv.Itoa(int(d)) + ")"
	}
	return s
}
