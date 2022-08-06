#include <stdio.h>
#include <stdlib.h>
#include <ctype.h>

int is_whitespace(char c)
{
	return c == '\n'
		|| c == ' '
		|| c == '\t';

}

void error(char *msg)
{
	printf(msg);
	printf("\n");
	exit(1);
}

typedef struct {
	char *data;
	char look;
	int index;
	int length;
} file_buffer;

void fb_init(file_buffer *fb, char *filename)
{
	FILE *file;
	int length;
	char *buffer;

	file = fopen(filename, "rb");
	if (file == 0) {
		error("error opening file");
	}
	if (fseek(file, 0, SEEK_END) != 0) {
		fclose(file);
		error("error seeking file");
	}
	length = ftell(file);
	if (length == -1) {
		fclose(file);
		error("error ftell");
	}
	if (fseek(file, 0, SEEK_SET)) {
		fclose(file);
		error("error seeking file");
	}

	buffer = malloc(length);
	if (length != fread(buffer, 1, length, file)) {
		fclose(file);
		free(buffer);
		error("error reading file");
	}

	fclose(file);
	fb->length = length;
	fb->data = buffer;
	fb->index = 0;
	fb->look = buffer[0];
}

void fb_next(file_buffer *fb)
{
	if (fb->index >= fb->length-1) {
		printf("reached end of file\n");
		return;
	}
	fb->look = fb->data[++fb->index];
}

void fb_match(file_buffer *fb, char c)
{
	if (fb->look != c) {
		char buf[256];
		sprintf(buf, "error: expected %c\n", c);
		error(buf);
	}
}

void fb_eatwhitespace(file_buffer *fb)
{
	while (is_whitespace(fb->look)) {
		fb_next(fb);
	}
}

void fb_function(file_buffer *fb) {
	fb_eatwhitespace(fb);
	while (isalpha(fb->look)) {
		printf("%c", fb->look);
		fb_next(fb);
	}
	fb_eatwhitespace(fb);
	fb_match(fb,')');
	printf("\n");
}


int main()
{
	file_buffer fb;
	fb_init(&fb, "main.lang");
	fb_function(&fb);
	return 0;
}
