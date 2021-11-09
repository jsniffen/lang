enum lexeme_id {
	ADD,
	SUBTRACT,
	DIVIDE,
	MULTIPLY,
	MODULUS,
	IDENT,
	LEFT_BRACE,
	RIGHT_BRACE,
	LEFT_PAREN,
	RIGHT_PAREN,
	COMMA,
};

struct lexeme {
	int line;
	int character;
	enum lexeme_id id;
	char name[256];
};

void print_lexeme(struct lexeme lex)
{
	char *id;
	switch (lex.id) {
		case ADD:
			id = "+";
			break;
		case SUBTRACT:
			id = "-";
			break;
		case DIVIDE:
			id = "/";
			break;
		case MULTIPLY:
			id = "*";
			break;
		case MODULUS:
			id = "%";
			break;
		case IDENT:
			id = "IDENT";
			break;
		case LEFT_BRACE:
			id = "{";
			break;
		case RIGHT_BRACE:
			id = "}";
			break;
		case LEFT_PAREN:
			id = "(";
			break;
		case RIGHT_PAREN:
			id = ")";
			break;
		case COMMA:
			id = ",";
			break;
	}
	printf("%d:%d\t%s\t\"%s\"\n", lex.line, lex.character, id, lex.name);
}

bool is_newline(char c)
{
	return c == '\n';
}

bool is_whitespace(char c)
{
	return c == ' ' || c == '\t';
}

bool is_comment(char *buf, int len)
{
	if (len < 2) return false;
	return buf[0] == '/' && buf[1] == '/';
}

bool is_operator(char c)
{
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '%';
}

bool is_number(char c)
{
	return c >= '0' && c <= '9';
}

bool is_character(char c)
{
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z');
}

bool is_identifier(char c)
{
	return is_number(c) || is_character(c);
}

int do_identifier(struct lexeme *lex, char *buffer, int len)
{
	lex->id = IDENT;
	int i;
	for (i = 0; i < len; ++i) {
		if (!is_character(buffer[i]) && !is_number(buffer[i])) break;

		lex->name[i] = buffer[i];
	}
	lex->name[i] = '\0';
	return i-1;
}

void do_operator(struct lexeme *lex, char c)
{
	switch (c) {
		case '+':
			lex->id = ADD;
			break;
		case '-':
			lex->id = SUBTRACT;
			break;
		case '*':
			lex->id = MULTIPLY;
			break;
		case '/':
			lex->id = DIVIDE;
			break;
		case '%':
			lex->id = MODULUS;
			break;
	}
}

bool is_delimiter(char c) {
	return  c == '{' ||
		c == '}' ||
		c == '(' ||
		c == ')' ||
		c == ',';
}

void do_delimiter(struct lexeme *lex, char c)
{
	switch (c) {
		case '{':
			lex->id = LEFT_BRACE;
			break;
		case '}':
			lex->id = RIGHT_BRACE;
			break;
		case '(':
			lex->id = LEFT_PAREN;
			break;
		case ')':
			lex->id = RIGHT_PAREN;
			break;
		case ',':
			lex->id = COMMA;
			break;
	};
}

int find_newline(char *buf, int len)
{
	int count = 0;
	while (!is_newline(*buf++) && --len > 0) ++count;
	return count;
}
