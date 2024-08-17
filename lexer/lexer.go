package lexer

import (
	"lang/token"
)

type Lexer struct {
	pos  int
	next int
	ch   byte
	data string
}

func New(s string) *Lexer {
	l := &Lexer{data: s}
	l.advance()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.eatWhitespace()

	switch ch := l.ch; ch {
	case '(':
		tok = token.Token{token.LPAREN, string(l.ch)}
	case ')':
		tok = token.Token{token.RPAREN, string(l.ch)}
	case '{':
		tok = token.Token{token.LBRACE, string(l.ch)}
	case '}':
		tok = token.Token{token.RBRACE, string(l.ch)}
	case '*':
		tok = token.Token{token.ASTERISK, string(l.ch)}
	case '+':
		tok = token.Token{token.PLUS, string(l.ch)}
	case ';':
		tok = token.Token{token.SEMICOLON, string(l.ch)}
	case '=':
		tok = token.Token{token.ASSIGN, string(l.ch)}
	case '"':
		value := l.eatString()
		return token.Token{token.STRING, value}
	default:
		if isAlpha(ch) {
			value := l.eatIdent()
			tokenType, ok := token.KeywordsMap[value]
			if ok {
				return token.Token{tokenType, value}
			} else {
				return token.Token{token.IDENT, value}
			}
		} else if isDigit(ch) {
			value := l.eatInt()
			return token.Token{token.INT, value}
		} else {
			return token.Token{token.EOF, ""}
		}
	}

	l.advance()

	return tok
}

func (l *Lexer) advance() {
	if l.next >= len(l.data) {
		l.pos = l.next
		l.ch = 0
	} else {
		l.pos = l.next
		l.ch = l.data[l.pos]
		l.next += 1
	}
}

func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.advance()
	}
}

func (l *Lexer) eatIdent() string {
	pos := l.pos
	for isAlpha(l.ch) || isDigit(l.ch) {
		l.advance()
	}
	return l.data[pos:l.pos]
}

func (l *Lexer) eatInt() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.advance()
	}
	return l.data[pos:l.pos]
}

func (l *Lexer) eatString() string {
	pos := l.pos
	l.advance()
	for l.ch != '"' {
		l.advance()
	}
	l.advance()
	return l.data[pos:l.pos]
}

func isAlpha(b byte) bool {
	return 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_'
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}
