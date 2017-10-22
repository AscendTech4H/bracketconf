package bracketconf

import "text/scanner"

//ASTNode is an AST node
type ASTNode struct {
	pos scanner.Position
	t   uint8
	val interface{}
}

const (
	tnil = iota
	tval
	tbrack
	tarr
	tdir
)

type astBrack []ASTNode

type astArr []ASTNode

type astDir struct {
	name string
	args []ASTNode
}

//Position returns the location of the value in the config file
func (a ASTNode) Position() scanner.Position {
	return a.pos
}

//IsBracket returns whether the ASTNode is a bracketed value
func (a ASTNode) IsBracket() bool {
	return a.t == tbrack
}

//IsArr returns whether the ASTNode is an array value
func (a ASTNode) IsArr() bool {
	return a.t == tarr
}

//IsDir returns whether the ASTNode is a directive
func (a ASTNode) IsDir() bool {
	return a.t == tdir
}

//IsValue returns whether the ASTNode is a value (array or otherwise)
func (a ASTNode) IsValue() bool {
	return a.IsArr() || a.t == tval
}
