package parser

import (
	"fmt"
	"lang/ast"
	"lang/lexer"
	"lang/token"
	"lang/types"
	"testing"
)

func TestVar(t *testing.T) {
	input := `
var x i32 = y
	`
	want := []ast.Statement{
		&ast.VarDecl{
			Token: token.Token{
				Type:  token.IDENT,
				Value: "x",
			},
			Type: &ast.Type{types.Int32},
			Value: &ast.Var{
				Token: token.Token{
					Type:  token.IDENT,
					Value: "y",
				},
			},
		},
	}
	test(t, input, want)
}

func TestFuncDecl(t *testing.T) {
	input := `
func main() i32 {
	return 1
}
func test()
	`
	want := []ast.Statement{
		&ast.FuncDecl{
			Token: token.Token{
				Type:  token.IDENT,
				Value: "main",
			},
			Body: []ast.Statement{
				&ast.Return{
					HasValue: true,
					Value:    &ast.IntLiteral{Value: 1},
				},
			},
			ReturnType: &ast.Type{types.Int32},
		},
		&ast.FuncDecl{
			Token: token.Token{
				Type:  token.IDENT,
				Value: "test",
			},
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
			Token: token.Token{
				Type:  token.IDENT,
				Value: "x",
			},
			Type:  &ast.Type{types.Int32},
			Value: &ast.IntLiteral{Value: 1},
		},
		&ast.VarDecl{
			Token: token.Token{
				Type:  token.IDENT,
				Value: "y",
			},
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
		if err := checkNode(got[i], want[i]); err != nil {
			t.Fatalf("[%d] %v", i, err)
		}
	}
}

func checkNode(gotNode, wantNode ast.Node) error {
	switch got := gotNode.(type) {
	case *ast.FuncDecl:
		want, ok := wantNode.(*ast.FuncDecl)
		if !ok {
			return fmt.Errorf("got *ast.FuncDecl, wanted %v", wantNode)
		}
		if err := checkFuncDecl(got, want); err != nil {
			return fmt.Errorf("*ast.FuncDecl: %v", err)
		}
	case *ast.Return:
		want, ok := wantNode.(*ast.Return)
		if !ok {
			return fmt.Errorf("got *ast.Return, wanted %v", wantNode)
		}
		if err := checkReturn(got, want); err != nil {
			return fmt.Errorf("*ast.Return: %v", err)
		}
	case *ast.Type:
		want, ok := wantNode.(*ast.Type)
		if !ok {
			return fmt.Errorf("got *ast.Type, wanted %v", wantNode)
		}
		if err := checkType(got, want); err != nil {
			return fmt.Errorf("*ast.Type: %v", err)
		}
	case *ast.IntLiteral:
		want, ok := wantNode.(*ast.IntLiteral)
		if !ok {
			return fmt.Errorf("got *ast.IntLiteral, wanted %v", wantNode)
		}
		if err := checkIntLiteral(got, want); err != nil {
			return fmt.Errorf("*ast.IntLiteral: %v", err)
		}
	case *ast.Var:
		want, ok := wantNode.(*ast.Var)
		if !ok {
			return fmt.Errorf("got *ast.Var, wanted %v", wantNode)
		}
		if err := checkVar(got, want); err != nil {
			return fmt.Errorf("*ast.Var: %v", err)
		}
	case *ast.VarDecl:
		want, ok := wantNode.(*ast.VarDecl)
		if !ok {
			return fmt.Errorf("got *ast.VarDecl, wanted %v", wantNode)
		}
		if err := checkVarDecl(got, want); err != nil {
			return fmt.Errorf("*ast.VarDecl: %v", err)
		}
	default:
		return fmt.Errorf("unsupported type, %v", gotNode)
	}
	return nil
}

func checkVar(got, want *ast.Var) error {
	if got.VarDecl != nil {
		return fmt.Errorf("VarDecl: expected nil, got %v", got.VarDecl)
	}

	if err := checkToken(got.Token, want.Token); err != nil {
		return fmt.Errorf("Name: %v", err)
	}

	return nil
}

func checkToken(got, want token.Token) error {
	if err := checkString(string(got.Type), string(want.Type)); err != nil {
		return fmt.Errorf("Type: %v", err)
	}

	if err := checkString(got.Value, want.Value); err != nil {
		return fmt.Errorf("Value: %v", err)
	}

	return nil
}

func checkVarDecl(got, want *ast.VarDecl) error {
	if err := checkToken(got.Token, want.Token); err != nil {
		return fmt.Errorf("Token: %v", err)
	}

	if err := checkType(got.Type, want.Type); err != nil {
		return fmt.Errorf("Type: %v", err)
	}

	if err := checkNode(got.Value, want.Value); err != nil {
		return fmt.Errorf("Value: %v", err)
	}

	return nil
}

func checkFuncDecl(got, want *ast.FuncDecl) error {
	if len(got.Args) != len(want.Args) {
		return fmt.Errorf("got %d args, want %d", len(got.Args), len(want.Args))
	}
	for i := range got.Args {
		if err := checkNode(got.Args[i], want.Args[i]); err != nil {
			return fmt.Errorf("args [%d]: %v", i, err)
		}
	}

	if len(got.Body) != len(want.Body) {
		return fmt.Errorf("got %d body statements, want %d", len(got.Body), len(want.Body))
	}
	for i := range got.Body {
		if err := checkNode(got.Body[i], want.Body[i]); err != nil {
			return fmt.Errorf("body [%d]: %v", i, err)
		}
	}

	if err := checkBool(got.Extern, want.Extern); err != nil {
		return fmt.Errorf("extern: %v", err)
	}

	if err := checkToken(got.Token, want.Token); err != nil {
		return fmt.Errorf("Token: %v", err)
	}

	if err := checkBool(got.Extern, want.Extern); err != nil {
		return err
	}

	if err := checkNode(got.ReturnType, want.ReturnType); err != nil {
		return err
	}

	return nil
}

func checkReturn(got, want *ast.Return) error {
	if err := checkBool(got.HasValue, want.HasValue); err != nil {
		return fmt.Errorf("HasValue: %v", err)
	}

	if err := checkNode(got.Value, want.Value); err != nil {
		return fmt.Errorf("Value: %v", err)
	}

	return nil
}

func checkType(got, want *ast.Type) error {
	return checkString(got.Name, want.Name)
}

func checkIntLiteral(got, want *ast.IntLiteral) error {
	return checkInt(got.Value, want.Value)
}

func checkInt(got, want int) error {
	if got != want {
		return fmt.Errorf("got %d, want %d", got, want)
	}
	return nil
}

func checkString(got, want string) error {
	if got != want {
		return fmt.Errorf("got %s, want %s", got, want)
	}
	return nil
}

func checkBool(got, want bool) error {
	if got != want {
		return fmt.Errorf("got %v, want %v", got, want)
	}
	return nil
}
