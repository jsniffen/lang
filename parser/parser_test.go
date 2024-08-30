package parser

import (
	"fmt"
	"lang/ast"
	"lang/lexer"
	"lang/types"
	"testing"
)

func TestFuncDecl(t *testing.T) {
	input := `
func main() i32 {}
func test()
	`
	want := []ast.Statement{
		&ast.FuncDecl{
			Name:       "main",
			Body:       []ast.Statement{},
			ReturnType: &ast.Type{types.Int32},
		},
		&ast.FuncDecl{
			Name:       "test",
			Body:       []ast.Statement{},
			ReturnType: &ast.Type{types.Void},
			Extern:     true,
		},
	}
	test(t, input, want)
}

func TestVarDecl(t *testing.T) {
	input := `
var x i32 = 1
var y i32 = 2
	`
	want := []ast.Statement{
		&ast.VarDecl{
			Name:  "x",
			Type:  &ast.Type{types.Int32},
			Value: &ast.IntLiteral{Value: 1},
		},
		&ast.VarDecl{
			Name:  "y",
			Type:  &ast.Type{types.Int32},
			Value: &ast.IntLiteral{Value: 2},
		},
	}
	test(t, input, want)
}

func test(t *testing.T, input string, want []ast.Statement) {
	l := lexer.New(input)
	p := New(l)

	prog, ok := p.ParseProgram()
	if !ok {
		for _, err := range p.Errors {
			t.Fatalf(err)
		}
	}

	got := prog.Statements
	if len(got) != len(want) {
		t.Fatalf("got %d statements, wanted %d", len(got), len(want))
	}

	for i := range got {
		if err := sameNode(got[i], want[i]); err != nil {
			t.Fatalf("[%d] %v", i, err)
		}
	}
}

func sameNode(gotNode, wantNode ast.Node) error {
	switch got := gotNode.(type) {
	case *ast.VarDecl:
		want, ok := wantNode.(*ast.VarDecl)
		if !ok {
			return fmt.Errorf("got *ast.VarDecl, wanted %v", wantNode)
		}
		if err := sameVarDecl(got, want); err != nil {
			return fmt.Errorf("*ast.VarDecl: %v", err)
		}
	case *ast.FuncDecl:
		want, ok := wantNode.(*ast.FuncDecl)
		if !ok {
			return fmt.Errorf("got *ast.FuncDecl, wanted %v", wantNode)
		}
		if err := sameFuncDecl(got, want); err != nil {
			return fmt.Errorf("*ast.FuncDecl: %v", err)
		}
	case *ast.Type:
		want, ok := wantNode.(*ast.Type)
		if !ok {
			return fmt.Errorf("got *ast.Type, wanted %v", wantNode)
		}
		if err := sameType(got, want); err != nil {
			return fmt.Errorf("*ast.Type: %v", err)
		}
	case *ast.IntLiteral:
		want, ok := wantNode.(*ast.IntLiteral)
		if !ok {
			return fmt.Errorf("got *ast.IntLiteral, wanted %v", wantNode)
		}
		if err := sameIntLiteral(got, want); err != nil {
			return fmt.Errorf("*ast.IntLiteral: %v", err)
		}
	default:
		return fmt.Errorf("unsupported type, %v", gotNode)
	}
	return nil
}

func sameVarDecl(got, want *ast.VarDecl) error {
	if err := sameString(got.Name, want.Name); err != nil {
		return fmt.Errorf("name: %v", err)
	}

	if err := sameType(got.Type, want.Type); err != nil {
		return fmt.Errorf("type: %v", err)
	}

	if err := sameNode(got.Value, want.Value); err != nil {
		return fmt.Errorf("value: %v", err)
	}

	return nil
}

func sameFuncDecl(got, want *ast.FuncDecl) error {
	if len(got.Args) != len(want.Args) {
		return fmt.Errorf("got %d args, want %d", len(got.Args), len(want.Args))
	}
	for i := range got.Args {
		if err := sameNode(got.Args[i], want.Args[i]); err != nil {
			return fmt.Errorf("args [%d]: %v", i, err)
		}
	}

	if len(got.Body) != len(want.Body) {
		return fmt.Errorf("got %d body statements, want %d", len(got.Body), len(want.Body))
	}
	for i := range got.Body {
		if err := sameNode(got.Body[i], want.Body[i]); err != nil {
			return fmt.Errorf("body [%d]: %v", i, err)
		}
	}

	if err := sameBool(got.Extern, want.Extern); err != nil {
		return fmt.Errorf("extern: %v", err)
	}

	if err := sameString(got.Name, want.Name); err != nil {
		return fmt.Errorf("name: %v", err)
	}

	if err := sameBool(got.Extern, want.Extern); err != nil {
		return err
	}

	if err := sameNode(got.ReturnType, want.ReturnType); err != nil {
		return err
	}

	return nil
}

func sameType(got, want *ast.Type) error {
	return sameString(got.Name, want.Name)
}

func sameIntLiteral(got, want *ast.IntLiteral) error {
	return sameInt(got.Value, want.Value)
}

func sameInt(got, want int) error {
	if got != want {
		return fmt.Errorf("got %d, want %d", got, want)
	}
	return nil
}

func sameString(got, want string) error {
	if got != want {
		return fmt.Errorf("got %s, want %s", got, want)
	}
	return nil
}

func sameBool(got, want bool) error {
	if got != want {
		return fmt.Errorf("got %v, want %v", got, want)
	}
	return nil
}
