package llvm

import (
	"fmt"
	"lang/ast"
	"lang/token"
	"lang/types"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Generator struct {
	module   *ir.Module
	function *ir.Func
	block    *ir.Block
	add      *ir.Func
	funcs    map[string]*ir.Func
	vars     map[string]value.Value
}

func NewGenerator() *Generator {
	return &Generator{
		module: ir.NewModule(),
		funcs:  make(map[string]*ir.Func),
		vars:   make(map[string]value.Value),
	}
}

func (g *Generator) Generate(program *ast.Program) string {
	for _, s := range program.Statements {
		g.genNode(s)
	}

	return g.module.String()
}

func (g *Generator) genNode(n ast.Node) value.Value {
	switch v := n.(type) {
	case *ast.FuncCall:
		return g.genFuncCall(v)
	case *ast.FuncDecl:
		return g.genFuncDecl(v)
	case *ast.InfixExpression:
		return g.genInfixExpression(v)
	case *ast.IntLiteral:
		return g.genIntLiteral(v)
	case *ast.Return:
		return g.genReturn(v)
	case *ast.Var:
		return g.genVar(v)
	case *ast.VarDecl:
		return g.genVarDecl(v)
	default:
		panic(fmt.Sprintf("cannot generate %T", v))
	}
}

func (g *Generator) genFuncCall(fc *ast.FuncCall) value.Value {
	if g.block == nil {
		panic("block is nil")
	}

	args := make([]value.Value, 0)
	for _, n := range fc.Args {
		args = append(args, g.genNode(n))
	}

	f, ok := g.funcs[fc.Token.Value]
	if !ok {
		panic(fmt.Sprintf("Cannot find func %s", fc.Token.Value))
	}
	return g.block.NewCall(f, args...)
}

func (g *Generator) genFuncDecl(fd *ast.FuncDecl) value.Value {
	ip := make([]*ir.Param, 0)
	for _, p := range fd.Params {
		ip = append(ip, ir.NewParam(p.Token.Value, irType(p.Type.Type)))
	}

	g.function = g.module.NewFunc(fd.Token.Value, irType(fd.ReturnType.Type), ip...)

	if !fd.Extern {
		g.block = g.function.NewBlock("")
		for _, n := range fd.Body {
			g.genNode(n)
		}
		g.block = nil
	}

	g.funcs[fd.Token.Value] = g.function
	g.function = nil

	return nil
}

func (g *Generator) genInfixExpression(ie *ast.InfixExpression) value.Value {
	l := g.genNode(ie.Left)
	r := g.genNode(ie.Right)

	switch ie.Token.Type {
	case token.PLUS:
		return g.block.NewAdd(l, r)
	case token.ASTERISK:
		return g.block.NewMul(l, r)
	default:
		panic(fmt.Sprintf("cannot generate %s", ie.Token.Type))
	}

}

func (g *Generator) genIntLiteral(il *ast.IntLiteral) value.Value {
	t := irType(il.Type()).(*irtypes.IntType)
	v := int64(il.Value)
	return constant.NewInt(t, v)
}

func (g *Generator) genReturn(r *ast.Return) value.Value {
	if g.block == nil {
		panic("g.block is nil")
	}
	v := g.genNode(r.Value)
	g.block.NewRet(v)

	return nil
}

func (g *Generator) genVar(v *ast.Var) value.Value {
	if g.function != nil {
		for _, p := range g.function.Params {
			if p.LocalIdent.LocalName == v.Token.Value {
				return p
			}
		}

		vd, ok := g.vars[v.Token.Value]
		if ok {
			aa := vd.(*ir.InstAlloca)
			return g.block.NewLoad(aa.Typ.ElemType, vd)
		}
	}
	panic(fmt.Sprintf("Could not find var %s", v.Token.Value))
}

func (g *Generator) genVarDecl(vd *ast.VarDecl) value.Value {
	if g.block == nil {
		panic("block is nil")
	}

	src := g.genNode(vd.Value)
	dst := g.block.NewAlloca(irType(vd.Type.Type))
	g.block.NewStore(src, dst)

	g.vars[vd.Token.Value] = dst

	return nil
}

func irType(t types.Type) irtypes.Type {
	switch v := t.(type) {
	case *types.Int32:
		return irtypes.NewInt(32)
	case *types.Pointer:
		return irtypes.NewPointer(irType(v.To))
	default:
		panic(fmt.Sprintf("cannot convert %T", v))
	}
}
