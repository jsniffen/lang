typedef struct {
	char *data;
	char look;
	int index;
	int length;
} file_buffer;

char write_buffer[2048];
int write_index = 0;


void fb_putchar(file_buffer *fb, char c)
{
	write_buffer[write_index++] = c;
}

void fb_putstr(file_buffer *fb, char *s)
{
	int i;
	for (i = 0; i < strlen(s); ++i) {
		fb_putchar(fb, s[i]);
	}
}

void fb_putnum(file_buffer *fb, int n)
{
	char buf[256];
	sprintf(buf, "%d", n);
	fb_putstr(fb, buf);
}

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
		exit(1);
	}
	fb_next(fb);
}

int fb_matchstr(file_buffer *fb, char *str)
{
	int i;
	for (i = 0; i < strlen(str); ++i) {
		if (fb->look != str[i]) {
			return 0;
		}
		fb_next(fb);
	}
	return 1;
}

void fb_eatwhitespace(file_buffer *fb)
{
	while (fb->look == ' ' || fb->look == '\t' || fb->look == '\n') {
		fb_next(fb);
	}
}

int fb_number(file_buffer *fb)
{
	char c;
	if (!isdigit(fb->look)) {
		printf("error: expected number\n");
	}
	c = fb->look;
	fb_next(fb);
	return c;
}

void fb_type(file_buffer *fb)
{
	fb_eatwhitespace(fb);
	if (fb_matchstr(fb, "i32")) {
	} else {
		printf("error: expected type\n");
	}
}

void fb_identifier(file_buffer *fb)
{
	while (isalpha(fb->look)) {
		fb_putchar(fb, fb->look);
		fb_next(fb);
	}
	fb_putchar(fb, ':');
}

void fb_expression(file_buffer *fb)
{
	char c;
	c = fb_number(fb);
	fb_putstr(fb, "\tmov rax, ");
	fb_putchar(fb, c);
	fb_putchar(fb, '\n');
	fb_match(fb, '+');
	c = fb_number(fb);
	fb_putstr(fb, "\tmov rdi, ");
	fb_putchar(fb, c);
	fb_putchar(fb, '\n');
	fb_putstr(fb, "\tadd rax, rdi");
	fb_putchar(fb, '\n');
}

void fb_function(file_buffer *fb)
{
	fb_eatwhitespace(fb);
	fb_identifier(fb);
	fb_eatwhitespace(fb);
	fb_match(fb,'(');
	fb_match(fb,')');
	fb_eatwhitespace(fb);
	fb_type(fb);
	fb_eatwhitespace(fb);
	fb_match(fb,'{');
	fb_eatwhitespace(fb);
	fb_expression(fb);
	fb_eatwhitespace(fb);
	fb_match(fb,'}');
}

void fb_parse(file_buffer *fb)
{
	fb_function(fb);
}

void fb_writefile(file_buffer *fb, char *filename)
{
	write_file(filename, write_buffer, write_index);
}
