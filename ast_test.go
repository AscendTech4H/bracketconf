package bracketconf

import (
	"encoding/json"
	"strings"
	"testing"
)

type testValue struct {
	Arr []interface{}
}

func (tv *testValue) add(v string) {
	tv.Arr = append(tv.Arr, v)
}

func (tv *testValue) sub() *testValue {
	s := &testValue{[]interface{}{}}
	tv.Arr = append(tv.Arr, s)
	return s
}

var testDirp DirectiveProcessor

func init() {
	testDirp = NewDirectiveProcessor(Directive{"hello", func(object interface{}, ans ...ASTNode) {
		tv := object.(*testValue)
		if len(ans) > 1 {
			tv = tv.sub()
		}
		for _, n := range ans {
			tv.hello(n)
		}
	}})
}

func (tv *testValue) hello(a ASTNode) {
	switch {
	case a.IsDir():
		fallthrough
	case a.IsBracket():
		a.Evaluate(tv.sub(), testDirp)
	case a.IsArr():
		sv := tv.sub()
		a.ForEach(func(i int, v ASTNode) {
			sv.hello(v)
		})
	default:
		tv.add(a.Text())
	}
}

func TestParseAST(t *testing.T) {
	a, err := ParseAST(strings.NewReader(`hello world {
		hello [(dank memes), pepe];
	};`), "testing.conf")
	if err != nil {
		t.Fatal(err.Error())
	}
	tv := &testValue{[]interface{}{}}
	tv.hello(*a)
	dat, err := json.Marshal(tv)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(dat))
}
