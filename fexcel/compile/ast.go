package compile

import (
	"text/scanner"
)

type File struct {
	pos   scanner.Position
	Nodes []Node
}

type Node interface {
	Pos() scanner.Position
}

type PointerNode struct {
	pos   scanner.Position
	Type  string
	Ident string
}

type TextNode struct {
	pos   scanner.Position
	Value string
}

type VarNode struct {
	pos   scanner.Position
	Type  string // e.g. R, PR
	Ident string // e.g. foo
}

func (f *File) Pos() scanner.Position        { return f.pos }
func (n *PointerNode) Pos() scanner.Position { return n.pos }
func (n *TextNode) Pos() scanner.Position    { return n.pos }
func (n *VarNode) Pos() scanner.Position     { return n.pos }
