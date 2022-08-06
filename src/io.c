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
