char *read_file(char *filename, char **data, int *length)
{
	FILE *file;
	file = fopen(filename, "rb");
	if (file == 0) {
		return "error opening file";
	}
	if (fseek(file, 0, SEEK_END) != 0) {
		fclose(file);
		return "error seeking file";
	}
	*length = ftell(file);
	if (*length == -1) {
		fclose(file);
		return "error ftell";
	}
	if (fseek(file, 0, SEEK_SET)) {
		fclose(file);
		return "error seeking file";
	}

	*data = malloc(*length);
	if (*length != fread(*data, 1, *length, file)) {
		fclose(file);
		free(*data);
		return "error reading file";
	}

	return 0;
}

char *write_file(char *filename, char *data, int length)
{
	FILE *file;
	int lenwrite;
	file = fopen(filename, "wb");
	if (file == 0) {
		return "error opening file";
	}
	lenwrite = fwrite(data, 1, length, file);
	fclose(file);
	return lenwrite == length ? 0 : "error writing file";
}

typedef struct {
	char *data;
	int cursor;
	int length;
} write_buffer;

void wb_init(write_buffer *wb)
{
	wb->cursor = 0;
	wb->length = 256;
	wb->data = malloc(256);
	memset(wb->data, 0, 256);
}

void wb_write(write_buffer *wb, char *s, int length)
{
	if (wb->cursor + length >= wb->length-1) {
		int newlen = 2*(wb->length+length);
		char *newdata = malloc(newlen);
		memset(newdata, 0, newlen);
		memcpy(newdata, wb->data, wb->length);
		free(wb->data);
		wb->data = newdata;
		wb->length = newlen;
	}
	memcpy(wb->data+wb->cursor, s, length);
	wb->cursor += length;
}
