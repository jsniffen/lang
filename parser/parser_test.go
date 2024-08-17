package parser

import (
	"lang/lexer"
	"testing"
)

func TestParser(t *testing.T) {
	input := "x = 1"
	l := lexer.New(input)
	p := New(l)

	prog := p.ParseProgram()
	if p.HasErrors() {
		p.PrintErrors()
		t.Fail()
		return
	}

	statements := []string{
		"x = 1",
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
