package bracketconf

import (
	"errors"
	"strconv"
	"strings"
)

//Int returns the value as an int
func (a ASTNode) Int() int {
	if a.t != tval {
		panic(ConfErr{a.pos, errors.New("Not a basic value")})
	}
	v, err := strconv.Atoi(a.val.(string))
	if err != nil {
		panic(ConfErr{a.pos, err})
	}
	return v
}

//Float returns the value as a float64
func (a ASTNode) Float() float64 {
	if a.t != tval {
		panic(ConfErr{a.pos, errors.New("Not a basic value")})
	}
	v, err := strconv.ParseFloat(a.val.(string), 64)
	if err != nil {
		panic(ConfErr{a.pos, err})
	}
	return v
}

//Text returns the value as a string
func (a ASTNode) Text() string {
	if a.t != tval {
		panic(ConfErr{a.pos, errors.New("Not a basic value")})
	}
	v := a.val.(string)
	if strings.HasPrefix(v, "\"") {
		u, err := strconv.Unquote(v)
		if err != nil {
			panic(ConfErr{a.pos, err})
		}
		v = u
	}
	return v
}

//Len returns the number of elements in the array
//If used on a bracketed value, it will return the number of directives
func (a ASTNode) Len() int {
	if !a.IsArr() {
		panic(ConfErr{a.pos, errors.New("Not an array")})
	}
	return len(a.val.(astArr))
}

//Index returns the value in the array at index n
//If used on a bracketed value, it will return the n'th directive
func (a ASTNode) Index(n int) ASTNode {
	if n >= a.Len() {
		panic(ConfErr{a.pos, errors.New("Index out of bounds")})
	}
	return a.val.(astArr)[n]
}

//ForEach runs a function with each index-value pair in the array
func (a ASTNode) ForEach(f func(i int, v ASTNode)) {
	if !a.IsArr() {
		panic(ConfErr{a.pos, errors.New("Not an array")})
	}
	arr := a.val.(astArr)
	for i, v := range arr {
		f(i, v)
	}
}
