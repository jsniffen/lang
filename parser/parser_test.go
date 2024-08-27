package parser

import (
	"fmt"
	"lang/ir"
	"lang/lexer"
	"strings"
	"testing"
)

func TestParseFuncDecl(t *testing.T) {
	input := `
	func a(a i32, b i32)
	func b() i32 {
		return 12345
	}
	`

	want := `
declare void @a(i32 %1, i32 %2)
define i32 @b() {
	ret i32 12345
}
`

	test(t, input, want)
}

func test(t *testing.T, input string, expected string) {
	l := lexer.New(input)
	p := New(l)

	prog, _ := p.ParseProgram()

	for _, err := range p.Errors {
		t.Fatalf(err)
	}

	var w ir.Writer
	prog.Codegen(&w)

	var want, got string

	got = strings.TrimSpace(w.String())
	got = strings.ReplaceAll(got, "\t", "")
	got = strings.ReplaceAll(got, " ", "")
	got = strings.ReplaceAll(got, "\n", "")
	got = strings.ReplaceAll(got, "\r", "")

	want = strings.TrimSpace(expected)
	want = strings.ReplaceAll(want, "\t", "")
	want = strings.ReplaceAll(want, " ", "")
	want = strings.ReplaceAll(want, "\n", "")
	want = strings.ReplaceAll(want, "\r", "")

	if got != want {
		for i := range got {
			if got[i] != want[i] {
				fmt.Printf("[%d], '%s' != '%s'\n", i, string(got[i]), string(want[i]))
				break
			}
		}
		t.Fatalf("got: \n'%s', want: \n'%s'", w.String(), expected)
	}
}
