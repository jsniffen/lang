#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include "file.h"
#include "lexer.h"


static struct lexeme lexemes[1000];
static int lexeme_index;

void push_lexeme(struct lexeme l)
{
	lexemes[lexeme_index++] = l;
}

int main()
{

	int i, line;
	char c;
	struct lexeme lex;

	struct file f;
	read_entire_file(&f, "main.language");
	line = 1;
	for (i = 0; i < f.length; ++i) {
		c = f.buffer[i];

		if (is_newline(c)) {
			++line;
			continue;
		}
		
		if (is_whitespace(c)) {
			continue;
		}

		if (is_comment(f.buffer + i, f.length)) {
			i += find_newline(f.buffer + i, f.length) - 1;
		}

		if (is_operator(f.buffer[i])) {
			do_operator(&lex, c, line);
			push_lexeme(lex);
		}

		if (is_delimiter(f.buffer[i])) {
			do_delimiter(&lex, c, line);
			push_lexeme(lex);
		}

		// if symbol

		// if number
	
		// if 
	}

	for (i = 0; i < lexeme_index; ++i) {
		print_lexeme(lexemes[i]);
	}

	free(f.buffer);
	return 0;
}
