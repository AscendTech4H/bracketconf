package bracketconf

import (
	"fmt"
	"text/scanner"
)

//ConfErr is an error type emitted in case of a parsing error
type ConfErr struct {
	Pos scanner.Position
	Err error
}

func (ce ConfErr) Error() string {
	return fmt.Sprintf("%s: %s", ce.Pos.String(), ce.Err.Error())
}

func (a ASTNode) errWrap() {
	err := recover()
	if err != nil {
		if ce, ok := err.(ConfErr); ok {
			panic(ce)
		} else {
			panic(ConfErr{a.pos, err.(error)})
		}
	}
}
