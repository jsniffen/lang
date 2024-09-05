package types

import "lang/token"

type Type interface {
	IsNumeric() bool
	Name() string
}

var (
	TypeInt32 = &Int32{}
	TypeNil   = &Nil{}
)

type Pointer struct {
	To Type
}

func (p *Pointer) IsNumeric() bool { return false }
func (p *Pointer) Name() string    { return "ptr" }

type Int32 struct{}

func (i *Int32) IsNumeric() bool { return true }
func (i *Int32) Name() string    { return "i32" }

type Nil struct{}

func (n *Nil) IsNumeric() bool { return false }
func (n *Nil) Name() string    { return "" }

type Custom struct {
	name string
}

func (c *Custom) IsNumeric() bool { return false }
func (c *Custom) Name() string    { return c.name }

func FromToken(t token.Token) Type {
	switch t.Value {
	case "i32":
		return TypeInt32
	case "nil":
		return TypeNil
	default:
		return &Custom{name: t.Value}
	}
}
