enum lexeme_id {
	ADD,
	SUBTRACT,
	DIVIDE,
	MULTIPLY,
	MODULUS,
	IDENTIFIER,
	LEFT_BRACE,
	RIGHT_BRACE,
	LEFT_PAREN,
	RIGHT_PAREN,
};

struct lexeme {
	int line;
	enum token_id id;
};

void print_lexeme(struct lexeme lex)
{
	char *id;
	switch (lex.id) {
		case ADD:
			id = "ADD";
			break;
		case SUBTRACT:
			id = "SUBTRACT";
			break;
		case DIVIDE:
			id = "DIVIDE";
			break;
		case MULTIPLY:
			id = "MULTIPLY";
			break;
		case MODULUS:
			id = "MODULUS";
			break;
		case IDENTIFIER:
			id = "IDENTIFIER";
			break;
		case LEFT_BRACE:
			id = "LEFT_BRACE";
			break;
		case RIGHT_BRACE:
			id = "RIGHT_BRACE";
			break;
		case LEFT_PAREN:
			id = "LEFT_PAREN";
			break;
		case RIGHT_PAREN:
			id = "RIGHT_PAREN";
			break;
	}
	printf("line: %d, id: %s\n", lex.line, id);
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

void do_operator(struct lexeme *lex, char c, int line)
{
	lex->line = line;
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
	return c == '{' || c == '}' || c == '(' || c == ')';
}

void do_delimiter(struct lexeme *lex, char c, int line)
{
	lex->line = line;
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
	};
}

int find_newline(char *buf, int len)
{
	int count = 0;
	while (!is_newline(*buf++) && --len > 0) ++count;
	return count;
}
