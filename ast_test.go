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
	}},
		Directive{"int", func(object interface{}, ans ...ASTNode) {
			tv := object.(*testValue)
			tv.Arr = append(tv.Arr, ans[0].Int())
		}},
		Directive{"float", func(object interface{}, ans ...ASTNode) {
			tv := object.(*testValue)
			tv.Arr = append(tv.Arr, ans[0].Float())
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
	tv, err := Parse(strings.NewReader(`
	//Comment
	hello world {
		/* this is a comment */
		hello [(dank memes), pepe];
		int 65;
		float 52.4;
	};`), "testing.conf", testDirp, &testValue{[]interface{}{}})
	if err != nil {
		t.Fatal(err.Error())
	}
	dat, err := json.Marshal(tv.(*testValue))
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(dat) != `{"Arr":[{"Arr":["world",{"Arr":[{"Arr":[{"Arr":["dank","memes"]},"pepe"]},65,52.4]}]}]}` {
		t.Fatalf("Incorrect parse %s", string(dat))
	}
	_, err = Parse(strings.NewReader("4"), "65.conf", DirectiveProcessor{}, nil)
	if err.Error() != "65.conf:1:2: Directives must only contain lowercase letters" {
		t.Fatalf("Incorrect error %s", err.Error())
	}
}
