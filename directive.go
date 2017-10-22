package bracketconf

import (
	"errors"
	"fmt"
)

//Directive is an object used to specify a directive to the AST processor
type Directive struct {
	Name     string
	Callback func(object interface{}, ast ...ASTNode)
}

//DirectiveProcessor is a utility object used to execute directives
type DirectiveProcessor struct {
	tbl map[string]func(interface{}, ...ASTNode)
}

//NewDirectiveProcessor creates a DirectiveProcessor using the specified set of directives
func NewDirectiveProcessor(ds ...Directive) (dp DirectiveProcessor) {
	dp.tbl = make(map[string]func(interface{}, ...ASTNode))
	for _, d := range ds {
		dp.tbl[d.Name] = d.Callback
	}
	return
}

//Evaluate evaluates the directive or bracket into the provided object using the specified DirectiveProcessor
func (a ASTNode) Evaluate(object interface{}, processor DirectiveProcessor) {
	if a.IsBracket() {
		for _, v := range a.val.(astBrack) {
			v.Evaluate(object, processor)
		}
	} else {
		if !a.IsDir() {
			panic(ConfErr{a.pos, errors.New("Not a directive")})
		}
		d := a.val.(astDir)
		f := processor.tbl[d.name]
		if f == nil {
			panic(ConfErr{a.pos, fmt.Errorf("Undefined directive %q", d.name)})
		}
		defer a.errWrap()
		f(object, d.args...)
	}
}
