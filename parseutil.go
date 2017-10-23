package bracketconf

import (
	"io"
	"os"
	"path/filepath"
)

//ParseFileAST parses AST from a config at the given file path
func ParseFileAST(filename string) (*ASTNode, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	a, err := ParseAST(f, filepath.Base(filename))
	if err != nil {
		return nil, err
	}
	return a, nil
}

//Parse parses a config from a reader using the given DirectiveProcessor for the root node into out
//Out is returned for convenience
func Parse(r io.Reader, filename string, dp DirectiveProcessor, out interface{}) (o interface{}, e error) {
	a, err := ParseAST(r, filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		ev := recover()
		if ev != nil {
			o = nil
			e = ev.(error)
		}
	}()
	a.Evaluate(out, dp)
	return out, nil
}

//ParseFile parses a config file - like Parse but opens and reads from a file
func ParseFile(filename string, dp DirectiveProcessor, out interface{}) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f, filename, dp, out)
}
