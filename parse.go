package bracketconf

import (
	"errors"
	"fmt"
	"io"
	"text/scanner"
)

type tcls uint8

const (
	tnilc tcls = iota
	teof
	tmeh
	tobrack
	tcbrack
	toparen
	tcparen
	tosqbrack
	tcsqbrack
	tcomma
	tsemicolon
)

func classify(s *scanner.Scanner) (tcls, string, scanner.Position) {
	n := s.Scan()
	tt := s.TokenText()
	pos := s.Pos()
	class := tmeh
	if n == scanner.EOF {
		class = teof
		goto ret
	}
	switch tt {
	case "":
		class = tnilc
	case "{":
		class = tobrack
	case "}":
		class = tcbrack
	case "(":
		class = toparen
	case ")":
		class = tcparen
	case "[":
		class = tosqbrack
	case "]":
		class = tcsqbrack
	case ",":
		class = tcomma
	case ";":
		class = tsemicolon
	}
ret:
	return class, tt, pos
}

func mehToAST(s string, p scanner.Position) *ASTNode {
	return &ASTNode{p, tval, s}
}

func parseParenArr(p scanner.Position, s *scanner.Scanner) *ASTNode {
	vals := astArr{}
	for class, val, pos := classify(s); class != tcparen; {
		switch class {
		case tnilc:
		case tmeh:
			vals = append(vals, *mehToAST(val, pos))
		case teof:
			panic(ConfErr{pos, io.ErrUnexpectedEOF})
		case tcbrack:
			panic(ConfErr{pos, errors.New("Closing bracket does not correspond to open bracket")})
		case tcomma:
			panic(ConfErr{pos, errors.New("Commas not allowed in shell-style arrays (consider using a JS-style array)")})
		case tcsqbrack:
			panic(ConfErr{pos, errors.New("Closing square bracket does not correspond to open square bracket")})
		default:
			panic(ConfErr{pos, fmt.Errorf("Token %q not allowed in shell-style array", val)})
		}
		class, val, pos = classify(s)
	}
	return &ASTNode{p, tarr, vals}
}

func parseSqBrackArr(p scanner.Position, s *scanner.Scanner) *ASTNode {
	vals := astArr{}
	var v *ASTNode
	for class, val, pos := classify(s); class != tcsqbrack; {
		switch class {
		case tnilc:
		case tmeh:
			if v != nil {
				panic(ConfErr{pos, errors.New("Multiple values not seperated by commas")})
			}
			v = mehToAST(val, pos)
		case toparen:
			if v != nil {
				panic(ConfErr{pos, errors.New("Multiple values not seperated by commas")})
			}
			v = parseParenArr(pos, s)
		case tosqbrack:
			if v != nil {
				panic(ConfErr{pos, errors.New("Multiple values not seperated by commas")})
			}
			v = parseSqBrackArr(pos, s)
		case tcomma:
			if v == nil {
				panic(ConfErr{pos, errors.New("No value preceding comma")})
			}
			vals = append(vals, *v)
			v = nil
		case tcparen:
			panic(ConfErr{pos, errors.New("Closing parenthese does not correspond to open parenthese")})
		case teof:
			panic(ConfErr{pos, io.ErrUnexpectedEOF})
		default:
			panic(ConfErr{pos, fmt.Errorf("Token %q not allowed in JS-style array", val)})
		}
		class, val, pos = classify(s)
	}
	switch {
	case v != nil:
		vals = append(vals, *v)
	case len(vals) != 0:
		panic(ConfErr{s.Pos(), errors.New("No value after comma")})
	}
	return &ASTNode{p, tarr, vals}
}

func parseDirective(name string, p scanner.Position, s *scanner.Scanner) *ASTNode {
	if name == "" {
		panic(ConfErr{p, errors.New("Empty directive")})
	}
	for _, c := range []rune(name) {
		if c < 'a' || c > 'z' {
			panic(ConfErr{p, errors.New("Directives must only contain lowercase letters")})
		}
	}
	args := []ASTNode{}
	for class, val, pos := classify(s); class != tsemicolon; {
		switch class {
		case tnilc:
		case tmeh:
			args = append(args, *mehToAST(val, pos))
		case teof:
			panic(ConfErr{pos, io.ErrUnexpectedEOF})
		case tobrack:
			args = append(args, *parseBracket(pos, s))
		case tosqbrack:
			args = append(args, *parseSqBrackArr(pos, s))
		case toparen:
			args = append(args, *parseParenArr(pos, s))
		default:
			panic(ConfErr{pos, fmt.Errorf("Token %q not allowed in directive", val)})
		}
		class, val, pos = classify(s)
	}
	return &ASTNode{p, tdir, astDir{name: name, args: args}}
}

func parseBracket(p scanner.Position, s *scanner.Scanner) *ASTNode {
	dirs := astBrack{}
	for class, val, pos := classify(s); class != tcbrack; {
		switch class {
		case tnilc:
		case tmeh:
			dirs = append(dirs, *parseDirective(val, pos, s))
		default:
			panic(ConfErr{pos, fmt.Errorf("Token %q not allowed in current context", val)})
		}
		class, val, pos = classify(s)
	}
	return &ASTNode{p, tbrack, dirs}
}

//ParseAST parses the data from an io.Reader into an AST and returns the root node
func ParseAST(r io.Reader, filnename string) (a *ASTNode, err error) {
	defer func() {
		e := recover()
		if e != nil {
			a = nil
			err = e.(error)
		}
	}()
	s := new(scanner.Scanner)
	s.Filename = filnename
	s.Init(r)
	posi := s.Position
	dirs := astBrack{}
	for class, val, pos := classify(s); class != teof; class, val, pos = classify(s) {
		switch class {
		case tnilc:
		case tmeh:
			dirs = append(dirs, *parseDirective(val, pos, s))
		default:
			panic(ConfErr{pos, fmt.Errorf("Token %q not allowed in current context", val)})
		}
	}
	return &ASTNode{posi, tbrack, dirs}, nil
}
