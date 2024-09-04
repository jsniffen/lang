package lexer

import (
	"lang/token"
	"testing"
)

func TestLex(t *testing.T) {
	input := `
extern func puts(^u8) i32;

func main() {
	puts("hello world");
	return;
}
	`
	lexer := New(input)

	tests := []token.Token{
		token.Token{Type: token.EXTERN, Value: "extern"},
		token.Token{Type: token.FUNC, Value: "func"},
		token.Token{Type: token.IDENT, Value: "puts"},
		token.Token{Type: token.LPAREN, Value: "("},
		token.Token{Type: token.POINTER, Value: "^"},
		token.Token{Type: token.IDENT, Value: "u8"},
		token.Token{Type: token.RPAREN, Value: ")"},
		token.Token{Type: token.IDENT, Value: "i32"},
		token.Token{Type: token.SEMICOLON, Value: ";"},
		token.Token{Type: token.FUNC, Value: "func"},
		token.Token{Type: token.IDENT, Value: "main"},
		token.Token{Type: token.LPAREN, Value: "("},
		token.Token{Type: token.RPAREN, Value: ")"},
		token.Token{Type: token.LBRACE, Value: "{"},
		token.Token{Type: token.IDENT, Value: "puts"},
		token.Token{Type: token.LPAREN, Value: "("},
		token.Token{Type: token.STRING, Value: "\"hello world\""},
		token.Token{Type: token.RPAREN, Value: ")"},
		token.Token{Type: token.SEMICOLON, Value: ";"},
		token.Token{Type: token.RETURN, Value: "return"},
		token.Token{Type: token.SEMICOLON, Value: ";"},
		token.Token{Type: token.RBRACE, Value: "}"},
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
		token.Token{Type: token.IDENT, Value: "x"},
		token.Token{Type: token.ASSIGN, Value: "="},
		token.Token{Type: token.INT, Value: "10"},
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

func TestLexMath(t *testing.T) {
	input := "0 - 1 + 2 * 3 / 4"
	lexer := New(input)

	tests := []token.Token{
		token.Token{Type: token.INT, Value: "0"},
		token.Token{Type: token.MINUS, Value: "-"},
		token.Token{Type: token.INT, Value: "1"},
		token.Token{Type: token.PLUS, Value: "+"},
		token.Token{Type: token.INT, Value: "2"},
		token.Token{Type: token.ASTERISK, Value: "*"},
		token.Token{Type: token.INT, Value: "3"},
		token.Token{Type: token.SLASH, Value: "/"},
		token.Token{Type: token.INT, Value: "4"},
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

func TestLineColumn(t *testing.T) {
	input := `x = 1
y = 2`
	lexer := New(input)

	tests := []token.Token{
		token.New(token.IDENT, "x", 1, 0, ""),
		token.New(token.ASSIGN, "=", 1, 2, ""),
		token.New(token.INT, "1", 1, 4, ""),
		token.New(token.IDENT, "y", 2, 0, ""),
		token.New(token.ASSIGN, "=", 2, 2, ""),
		token.New(token.INT, "2", 2, 4, ""),
	}
	testIndex := 0
	for got := lexer.NextToken(); got.Type != token.EOF; got = lexer.NextToken() {
		want := tests[testIndex]

		if got.Type != want.Type || got.Value != want.Value {
			t.Fatalf("[%d] got: %v, want: %v", testIndex, got, want)
		}

		if got.Column != want.Column || got.Line != want.Line {
			t.Fatalf("[%d] got: %v, want: %v", testIndex, got, want)
		}
		testIndex += 1
	}

	if testIndex < len(tests) {
		t.Fatalf("Only produced %d token(s), wanted: %d", testIndex, len(tests))
	}
}
