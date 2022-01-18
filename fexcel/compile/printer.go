package compile

import (
	"fmt"
	"strings"
	"text/scanner"

	"github.com/onerobotics/fexcel/fexcel"
)

type Printer struct {
	*cache
	errors ErrorList
	b      strings.Builder
}

func NewPrinter(fpath string, cfg fexcel.FileConfig) (*Printer, error) {
	cache, err := newCache(fpath, cfg)
	if err != nil {
		return nil, err
	}

	var p Printer
	p.cache = cache

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
