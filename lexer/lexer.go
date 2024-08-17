package lexer

import (
	"lang/token"
	"os"
)

type Lexer struct {
	pos      int
	next     int
	ch       byte
	data     string
	line     int
	col      int
	Filename string
}

func New(s string) *Lexer {
	l := &Lexer{data: s}
	l.advance()
	l.line = 1
	l.col = 0
	return l
}

func FromFile(f string) (*Lexer, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	l := New(string(b))
	l.Filename = f
	return l, nil
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.eatWhitespace()

	line, col := l.line, l.col
	switch ch := l.ch; ch {
	case '(':
		tok = token.New(token.LPAREN, string(l.ch), line, col)
	case ')':
		tok = token.New(token.RPAREN, string(l.ch), line, col)
	case '{':
		tok = token.New(token.LBRACE, string(l.ch), line, col)
	case '}':
		tok = token.New(token.RBRACE, string(l.ch), line, col)
	case '*':
		tok = token.New(token.ASTERISK, string(l.ch), line, col)
	case '+':
		tok = token.New(token.PLUS, string(l.ch), line, col)
	case ';':
		tok = token.New(token.SEMICOLON, string(l.ch), line, col)
	case '=':
		tok = token.New(token.ASSIGN, string(l.ch), line, col)
	case '/':
		tok = token.New(token.SLASH, string(l.ch), line, col)
	case '-':
		tok = token.New(token.MINUS, string(l.ch), line, col)
	case '"':
		value := l.eatString()
		return token.New(token.STRING, value, line, col)
	default:
		if isAlpha(ch) {
			value := l.eatIdent()
			tokenType, ok := token.KeywordsMap[value]
			if ok {
				return token.New(tokenType, value, line, col)
			} else {
				return token.New(token.IDENT, value, line, col)
			}
		} else if isDigit(ch) {
			value := l.eatInt()
			return token.New(token.INT, value, line, col)
		} else {
			return token.New(token.EOF, "", line, col)
		}
	}

	l.advance()

	return tok
}

func (l *Lexer) advance() {
	if l.ch == '\n' {
		l.line += 1
		l.col = 0
	} else {
		l.col += 1
	}

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
