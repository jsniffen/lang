struct file {
	char *buffer;
	long length;
};

void read_entire_file(struct file *f, char *fn) {
	FILE *fi = fopen(fn, "rb");
	fseek(fi, 0, SEEK_END);
	f->length = ftell(fi);
	fseek(fi, 0, SEEK_SET);

	f->buffer = (char *)malloc(f->length + 1);
	fread(f->buffer, 1, f->length, fi);
	fclose(fi);
}
