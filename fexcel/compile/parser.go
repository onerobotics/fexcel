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

func (p *parser) parseVar(typ string) Node {
	pos, lit := p.pos, p.lit
	p.expect(scanner.Ident)
	p.expectRbrace()

	return &VarNode{pos: pos, Type: typ, Ident: lit}
}

func (p *parser) parsePointer() Node {
	pos := p.pos
	p.next() // consume &
	typ := p.lit
	p.expect(scanner.Ident)
	p.expectLbrace()
	lit := p.lit
	p.expect(scanner.Ident)
	p.expectRbrace()

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
		msg := fmt.Sprintf("expected" + string(tok) + "but got" + got)
		p.error(p.scanner.Position, msg)
	}
}

func (p *parser) expectLbrace() {
	if p.lit == "{" {
		p.next()
	} else {
		got := scanner.TokenString(p.tok)
		msg := fmt.Sprintf("expected { but got " + got)
		p.error(p.scanner.Position, msg)
	}
}

func (p *parser) expectRbrace() {
	if p.lit == "}" {
		p.next()
	} else {
		got := scanner.TokenString(p.tok)
		msg := fmt.Sprintf("expected } but got " + got)
		p.error(p.scanner.Position, msg)
	}
}

func (p *parser) parseFile() *File {
	var f File

	for p.tok != scanner.EOF {
		switch p.tok {
		case scanner.Ident:
			typ := p.lit
			if p.scanner.Peek() == '{' {
				p.next()
				p.next()
				f.Nodes = append(f.Nodes, p.parseVar(typ))
			} else {
				f.Nodes = append(f.Nodes, p.parseText())
			}
		default:
			if p.lit == "&" {
				f.Nodes = append(f.Nodes, p.parsePointer())
			} else {
				f.Nodes = append(f.Nodes, p.parseText())
			}
		}
		//fmt.Printf("%s: %s %s\n", p.pos, scanner.TokenString(p.tok), p.lit)
		//p.next()
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
