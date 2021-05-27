package compile

import (
	"fmt"
	"strings"
	"text/scanner"

	"github.com/onerobotics/fexcel/fexcel"
)

type Printer struct {
	Definitions map[string]map[string]int
	Constants   map[string]string
	errors      ErrorList
	b           strings.Builder
}

func NewPrinter(fpath string, cfg fexcel.FileConfig) (*Printer, error) {
	var p Printer
	p.Definitions = make(map[string]map[string]int)

	spreadsheet, err := fexcel.OpenFile(fpath, cfg)
	if err != nil {
		return nil, err
	}

	allDefs, err := spreadsheet.AllDefinitions()
	if err != nil {
		return nil, err
	}

	for t, _ := range spreadsheet.Locations {
		switch t {
		case fexcel.Constant:
			// noop
		default:
			p.Definitions[t.String()] = make(map[string]int)
		}
	}
	for t, defs := range allDefs {
		for _, def := range defs {
			p.Definitions[t.String()][def.Comment] = def.Id
		}
	}

	if spreadsheet.Locations[fexcel.Constant] != nil {
		p.Constants, err = spreadsheet.Constants()
		if err != nil {
			return nil, err
		}
	} else {
		p.Constants = make(map[string]string)
	}

	p.Definitions["UI"] = map[string]int{
		"IMSTP":      1,
		"Hold":       2,
		"SFSPD":      3,
		"CycleStop":  4,
		"FaultReset": 5,
		"Start":      6,
		"Home":       7,
		"Enable":     8,
		"ProdStart":  18,
	}
	p.Definitions["UO"] = map[string]int{
		"CmdEnabled":  1,
		"SystemReady": 2,
		"PrgRunning":  3,
		"PrgPaused":   4,
		"MotionHeld":  5,
		"Fault":       6,
		"AtPerch":     7,
		"TPEnabled":   8,
		"BattAlarm":   9,
		"Busy":        10,
	}
	p.Definitions["SI"] = map[string]int{
		"FaultReset": 1,
		"Remote":     2,
		"Hold":       3,
		"UserPB1":    4,
		"UserPB2":    5,
		"CycleStart": 6,
	}
	p.Definitions["SO"] = map[string]int{
		"RemoteLED":  0,
		"CycleStart": 1,
		"Hold":       2,
		"FaultLED":   3,
		"BattAlarm":  4,
		"UserLED1":   5,
		"UserLED2":   6,
		"TPEnabled":  7,
	}

	return &p, nil
}

func (p *Printer) error(pos scanner.Position, msg string) {
	p.errors.Add(pos, msg)
}

func (p *Printer) Reset() {
	p.errors.Reset()
	p.b.Reset()
}

func (p *Printer) Print(nodes ...Node) error {
	for _, node := range nodes {
		switch n := node.(type) {
		case *File:
			p.Print(n.Nodes...)
		case *PointerNode:
			if i, ok := p.Definitions[n.Type][n.Ident]; ok {
				fmt.Fprint(&p.b, fmt.Sprintf("%d", i))
			} else {
				p.error(n.Pos(), fmt.Sprintf("&%s{%s} is undefined", n.Type, n.Ident))
			}
		case *TextNode:
			fmt.Fprint(&p.b, n.Value)
		case *VarNode:
			if n.Type == "$" {
				if value, ok := p.Constants[n.Ident]; ok {
					fmt.Fprint(&p.b, value)
				} else {
					p.error(n.Pos(), fmt.Sprintf("${%s} is undefined", n.Ident))
				}
			} else {
				if i, ok := p.Definitions[n.Type][n.Ident]; ok {
					fmt.Fprint(&p.b, fmt.Sprintf("%s[%d:%s]", n.Type, i, n.Ident))
				} else {
					p.error(n.Pos(), fmt.Sprintf("%s{%s} is undefined", n.Type, n.Ident))
				}
			}
		}
	}

	return p.errors.Err()
}

func (p *Printer) Output() string {
	return p.b.String()
}
