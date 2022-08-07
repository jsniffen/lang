#include <stdio.h>
#include <stdlib.h>
#include <ctype.h>
#include <string.h>

#include "src/io.c"
#include "src/file_buffer.c"

int main()
{
	file_buffer fb;
	fb_init(&fb, "main.lang");
	fb_parse(&fb);
	fb_writefile(&fb, "main.asm");
	return 0;
}
