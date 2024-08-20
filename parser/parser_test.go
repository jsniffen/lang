package parser

import (
	"lang/lexer"
	"testing"
)

func TestParser(t *testing.T) {
	input := "func main() i32 {}"
	l := lexer.New(input)
	p := New(l)

	prog, ok := p.ParseProgram()
	if !ok {
		p.PrintErrors()
		t.Fail()
		return
	}

	statements := []string{
		"func main() i32",
	}

	if len(prog.Statements) != len(statements) {
		t.Fatalf("got: %d program statements, want: %d", len(prog.Statements), len(statements))
	}

	for i, s := range prog.Statements {
		got := s.String()
		want := statements[i]

		if got != want {
			t.Fatalf("got: %v, want: %v", got, want)
		}
	}
}

func TestParseTypes(t *testing.T) {
	input := `func main() {
		return 1 + 2
	}`
	l := lexer.New(input)
	p := New(l)

	_, ok := p.ParseProgram()
	if !ok {
		p.PrintErrors()
		t.Fail()
		return
	}
}
