typedef struct {
	char *data;
	char look;
	int index;
	int length;
} file_buffer;

char write_buffer[2048];
int write_index = 0;

void fb_init(file_buffer *fb, char *filename)
{
	read_file(filename, &fb->data, &fb->length);
	fb->index = 0;
	fb->look = fb->data[0];
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
		printf("%s\n", buf);
	}
}

void fb_eatwhitespace(file_buffer *fb)
{
	while (fb->look == ' ' || fb->look == '\t' || fb->look == '\n') {
		fb_next(fb);
	}
}

void fb_function(file_buffer *fb)
{
	fb_eatwhitespace(fb);
	while (isalpha(fb->look)) {
		printf("%c", fb->look);
		fb_next(fb);
	}
	fb_eatwhitespace(fb);
	fb_match(fb,'(');
	printf("\n");
}

void fb_parse(file_buffer *fb)
{
	int len = sprintf(write_buffer+write_index, "writing v2\n");
	write_index += len;
}

void fb_write(file_buffer *fb, char *filename)
{
	write_file(filename, write_buffer, write_index);
}
