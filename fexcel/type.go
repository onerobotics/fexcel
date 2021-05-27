package fexcel

import "strconv"

type Type int

const (
	Constant Type = iota
	Numreg
	Posreg
	Ualm
	Ain
	Aout
	Din
	Dout
	Gin
	Gout
	Rin
	Rout
	Sreg
	Flag
	Uin
	Uout
	Sin
	Sout
)

var types = [...]string{
	Constant: "Constant",
	Numreg:   "R",
	Posreg:   "PR",
	Ualm:     "UALM",
	Ain:      "AI",
	Aout:     "AO",
	Din:      "DI",
	Dout:     "DO",
	Gin:      "GI",
	Gout:     "GO",
	Rin:      "RI",
	Rout:     "RO",
	Sreg:     "SR",
	Flag:     "F",
	Uin:      "UI",
	Uout:     "UOUT",
	Sin:      "SI",
	Sout:     "SO",
}

func (t Type) String() string {
	s := ""
	if 0 <= t && t < Type(len(types)) {
		s = types[t]
	}
	if s == "" {
		s = "type(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}
