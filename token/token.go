package token

type TokenType string

const (
	ASSIGN    = "="
	ASTERISK  = "ASTERISK"
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
)

var KeywordsMap = map[string]TokenType{
	"extern": EXTERN,
	"func":   FUNC,
	"return": RETURN,
}

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

func New(t TokenType, v string, l, c int) Token {
	return Token{
		Type:   t,
		Value:  v,
		Line:   l,
		Column: c,
	}
}
