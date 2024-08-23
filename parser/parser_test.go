package parser

import (
	"bytes"
	"lang/ast"
	"lang/lexer"
	"lang/types"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestParseFuncDecl(t *testing.T) {
	input := `
	func a(a i32, b i32)
	func b() i32 {}
	`
	stmts := &ast.Program{
		Statements: []ast.Statement{
			&ast.FuncDecl{
				Args: []*ast.FuncArg{
					&ast.FuncArg{
						Name:     "a",
						Type:     &ast.Type{types.Int32},
						Location: "%a",
					},
					&ast.FuncArg{
						Name:     "b",
						Type:     &ast.Type{types.Int32},
						Location: "%b",
					},
				},
				Name:       "a",
				ReturnType: &ast.Type{types.Void},
				Extern:     true,
			},
			&ast.FuncDecl{
				Body:       []ast.Statement{},
				Name:       "b",
				ReturnType: &ast.Type{types.Int32},
				Extern:     false,
			},
		},
	}

	test(t, input, stmts)
}

func test(t *testing.T, input string, want *ast.Program) {
	l := lexer.New(input)
	p := New(l)

	got, _ := p.ParseProgram()

	for _, err := range p.Errors {
		t.Fatalf(err)
	}

	if !reflect.DeepEqual(got, want) {
		var gotBuf bytes.Buffer
		var wantBuf bytes.Buffer

		got.Codegen(&gotBuf)

		want.Codegen(&wantBuf)

		t.Fatalf("got: \n%s, want: \n%s, got: %s, want: %s", gotBuf.String(), wantBuf.String(), spew.Sdump(got), spew.Sdump(want))
	}
}
