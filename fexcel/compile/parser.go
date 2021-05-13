package compile

import (
	"fmt"
	"strings"
	"text/scanner"
)

type parser struct {
	scanner scanner.Scanner
	errors  ErrorList

	pos scanner.Position
	tok rune
	lit string
}

func (p *parser) init(filename string, src string) {
	p.scanner.Init(strings.NewReader(src))
	p.scanner.Mode = scanner.ScanIdents | scanner.ScanInts
	p.scanner.Filename = filename
	p.scanner.Whitespace = 0
	p.scanner.Error = func(s *scanner.Scanner, msg string) {
		p.errors.Add(s.Position, msg)
	}

	p.next()
}

func (p *parser) error(pos scanner.Position, msg string) {
	p.errors.Add(pos, msg)
}

func (p *parser) next() {
	p.tok = p.scanner.Scan()
	p.pos = p.scanner.Position
	p.lit = p.scanner.TokenText()
}

func (p *parser) parseVar() Node {
	pos, typ := p.pos, p.lit
	p.next() // typ
	p.next() // {
	lit := p.lit
	p.expect(scanner.Ident)
	p.expectLit("}")

	return &VarNode{pos: pos, Type: typ, Ident: lit}
}

func (p *parser) parsePointer() Node {
	pos := p.pos
	p.next() // consume &
	typ := p.lit
	p.expect(scanner.Ident)
	p.expectLit("{")
	lit := p.lit
	p.expect(scanner.Ident)
	p.expectLit("}")

	return &PointerNode{pos: pos, Type: typ, Ident: lit}
}

func (p *parser) parseText() Node {
	pos, lit := p.pos, p.lit
	p.next()
	return &TextNode{pos: pos, Value: lit}
}

func (p *parser) expect(tok rune) {
	if p.tok == tok {
		p.next()
	} else {
		got := scanner.TokenString(p.tok)
		msg := fmt.Sprintf("expcted %q but got %q", string(tok), got)
		p.error(p.scanner.Position, msg)
	}
}

func (p *parser) expectLit(lit string) {
	if p.lit == lit {
		p.next()
	} else {
		got := scanner.TokenString(p.tok)
		msg := fmt.Sprintf("expected %q but got %q", lit, got)
		p.error(p.scanner.Position, msg)
	}
}

func (p *parser) parseFile() *File {
	var f File

	for p.tok != scanner.EOF {
		switch p.tok {
		case scanner.Ident:
			if p.scanner.Peek() == '{' {
				f.Nodes = append(f.Nodes, p.parseVar())
			} else {
				f.Nodes = append(f.Nodes, p.parseText())
			}
		default:
			switch p.lit {
			case "&":
				f.Nodes = append(f.Nodes, p.parsePointer())
			case "$":
				if p.scanner.Peek() == '{' {
					f.Nodes = append(f.Nodes, p.parseVar())
				} else {
					f.Nodes = append(f.Nodes, p.parseText())
				}
			default:
				f.Nodes = append(f.Nodes, p.parseText())
			}
		}
	}
	return &f
}

func Parse(filename string, src string) (*File, error) {
	var p parser
	p.init(filename, src)
	f := p.parseFile()

	err := p.errors.Err()
	return f, err
}
