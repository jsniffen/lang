package types

import "lang/token"

type Type interface {
	IsNumeric() bool
	Name() string
}

var (
	Int32 = &int32{}
	Nil   = &nil{}
)

type Pointer struct {
	To Type
}

func (p *Pointer) IsNumeric() bool { return false }
func (p *Pointer) Name() string    { return "ptr" }

type int32 struct{}

func (i *int32) IsNumeric() bool { return true }
func (i *int32) Name() string    { return "i32" }

type nil struct{}

func (n *nil) IsNumeric() bool { return false }
func (n *nil) Name() string    { return "" }

type Custom struct {
	name string
}

func (c *Custom) IsNumeric() bool { return false }
func (c *Custom) Name() string    { return c.name }

func FromToken(t token.Token) Type {
	switch t.Value {
	case "i32":
		return Int32
	case "nil":
		return Nil
	default:
		return &Custom{name: t.Value}
	}
}
