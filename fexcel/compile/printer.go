package compile

import (
	"fmt"
	"strings"
	"text/scanner"

	"github.com/unreal/fexcel/fexcel"
)

type Printer struct {
	Definitions map[string]map[string]int
	errors      ErrorList
	b           strings.Builder
}

func NewPrinter(fpath string, cfg fexcel.FileConfig) (*Printer, error) {
	var p Printer
	p.Definitions = make(map[string]map[string]int)

	spreadsheet, err := fexcel.OpenFile("test.xlsx", cfg)
	if err != nil {
		return nil, err
	}

	allDefs, err := spreadsheet.AllDefinitions()
	if err != nil {
		return nil, err
	}

	for t, _ := range spreadsheet.Locations {
		p.Definitions[t.String()] = make(map[string]int)
	}
	for t, defs := range allDefs {
		for _, def := range defs {
			p.Definitions[t.String()][def.Comment] = def.Id
		}
	}

	return &p, nil
}

func (p *Printer) error(pos scanner.Position, msg string) {
	p.errors.Add(pos, msg)
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
			if i, ok := p.Definitions[n.Type][n.Ident]; ok {
				fmt.Fprint(&p.b, fmt.Sprintf("%s[%d:%s]", n.Type, i, n.Ident))
			} else {
				p.error(n.Pos(), fmt.Sprintf("%s{%s} is undefined", n.Type, n.Ident))
			}
		}
	}

	return p.errors.Err()
}

func (p *Printer) Output() string {
	return p.b.String()
}
