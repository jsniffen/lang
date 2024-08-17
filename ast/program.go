package ast

import "bytes"

type Program struct {
	Statements []Statement
}

func (p *Program) DebugString(i int) string {
	var out bytes.Buffer
	out.WriteString("Program")
	for _, s := range p.Statements {
		out.WriteString(s.DebugString(i + 1))
	}
	return out.String()
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	return out.String()
}
