package checker

import (
	"lang/lexer"
	"lang/parser"
	"testing"
)

func TestChecker(t *testing.T) {
	input := `
func main() {
	var x i32 = 1
}
	`
	l := lexer.New(input)
	p := parser.New(l)
	prog, ok := p.ParseProgram()
	if !ok {
		for _, err := range p.Errors {
			t.Fatalf(err)
		}
	}

	checker := New(prog)
	if !checker.Check() {
		t.Fail()
		return
	}
}
