package lexer

import (
	"lang/token"
	"testing"
)

func TestLex(t *testing.T) {
	input := `
extern func puts(*u8) i32;

func main() {
	puts("hello world");
	return;
}
	`
	lexer := New(input)

	tests := []token.Token{
		token.Token{token.EXTERN, "extern"},
		token.Token{token.FUNC, "func"},
		token.Token{token.IDENT, "puts"},
		token.Token{token.LPAREN, "("},
		token.Token{token.ASTERISK, "*"},
		token.Token{token.IDENT, "u8"},
		token.Token{token.RPAREN, ")"},
		token.Token{token.IDENT, "i32"},
		token.Token{token.SEMICOLON, ";"},
		token.Token{token.FUNC, "func"},
		token.Token{token.IDENT, "main"},
		token.Token{token.LPAREN, "("},
		token.Token{token.RPAREN, ")"},
		token.Token{token.LBRACE, "{"},
		token.Token{token.IDENT, "puts"},
		token.Token{token.LPAREN, "("},
		token.Token{token.STRING, "\"hello world\""},
		token.Token{token.RPAREN, ")"},
		token.Token{token.SEMICOLON, ";"},
		token.Token{token.RETURN, "return"},
		token.Token{token.SEMICOLON, ";"},
		token.Token{token.RBRACE, "}"},
	}
	testIndex := 0
	for got := lexer.NextToken(); got.Type != token.EOF; got = lexer.NextToken() {
		want := tests[testIndex]

		if got.Type != want.Type || got.Value != want.Value {
			t.Fatalf("[%d] got: %v, want: %v", testIndex, got, want)
		}
		testIndex += 1
	}

	if testIndex < len(tests) {
		t.Fatalf("Only produced %d token(s), wanted: %d", testIndex, len(tests))
	}
}

func TestLexInt(t *testing.T) {
	input := "x = 10"
	lexer := New(input)

	tests := []token.Token{
		token.Token{token.IDENT, "x"},
		token.Token{token.ASSIGN, "="},
		token.Token{token.INT, "10"},
	}
	testIndex := 0
	for got := lexer.NextToken(); got.Type != token.EOF; got = lexer.NextToken() {
		want := tests[testIndex]

		if got.Type != want.Type || got.Value != want.Value {
			t.Fatalf("[%d] got: %v, want: %v", testIndex, got, want)
		}
		testIndex += 1
	}

	if testIndex < len(tests) {
		t.Fatalf("Only produced %d token(s), wanted: %d", testIndex, len(tests))
	}
}
