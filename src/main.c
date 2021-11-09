#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#include "file.h"
#include "lexer.h"


static struct lexeme lexemes[1000];
static int lexeme_index;

int main()
{
	int i, line, character;
	char c;

	struct file f;
	read_entire_file(&f, "main.language");
	line = 1;
	character = 0;
	for (i = 0; i < f.length; ++i) {
		c = f.buffer[i];
		++character;

		if (is_newline(c)) {
			++line;
			character = 0;
			continue;
		}
		
		if (is_whitespace(c)) {
			continue;
		}

		if (is_comment(f.buffer + i, f.length)) {
			i += find_newline(f.buffer + i, f.length);
			++line;
			character = 0;
			continue;
		}

		struct lexeme *lex = lexemes + lexeme_index;
		lex->line = line;
		lex->character = character;

		if (is_operator(f.buffer[i])) {
			do_operator(lex, c);
			++lexeme_index;
			continue;
		}

		if (is_delimiter(f.buffer[i])) {
			do_delimiter(lex, c);
			++lexeme_index;
			continue;
		}

		if (is_identifier(f.buffer[i])) {
			int n;
			n = do_identifier(lex, f.buffer+i, f.length);
			i += n;
			character += n;
			++lexeme_index;
			continue;
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
