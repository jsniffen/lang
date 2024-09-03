package token

import "fmt"

type TokenType string

const (
	ASSIGN    = "="
	ASTERISK  = "*"
	COMMA     = ","
	EOF       = "EOF"
	EXTERN    = "EXTERN"
	FUNC      = "FUNC"
	IDENT     = "IDENT"
	INT       = "INT"
	LBRACE    = "{"
	LPAREN    = "("
	MINUS     = "-"
	PLUS      = "+"
	RBRACE    = "}"
	RETURN    = "RETURN"
	RPAREN    = ")"
	SEMICOLON = ";"
	SLASH     = "/"
	STRING    = "STRING"
	VAR       = "VAR"
)

var KeywordsMap = map[string]TokenType{
	"extern": EXTERN,
	"func":   FUNC,
	"return": RETURN,
	"var":    VAR,
}

type Token struct {
	Column   int
	Filename string
	Line     int
	Type     TokenType
	Value    string
}

func New(t TokenType, v string, l, c int, f string) Token {
	return Token{
		Column:   c,
		Filename: f,
		Line:     l,
		Type:     t,
		Value:    v,
	}
}

func (t Token) Path() string {
	return fmt.Sprintf("%s:%d:%d", t.Filename, t.Line, t.Column)
}
