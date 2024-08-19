

declare i32 @puts(i8*)

define void @main() {
	%s [13 x i8] = c"hello world\0"
	call i32 @puts([13 x i8] %s)
}
