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
	Type  TokenType
	Value string
}
