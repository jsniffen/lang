look = None

def get_char():
    global look
    look = input()[0]

def error(s):
    print(f"Error: {s}.")

def abort(s):
    error(s)
    exit(1)

def expected(s):
    abort(f"{s} Expected")

def match(c):
    if c == look:
        pass
    else:
        expected(f"'{c}'")

def get_name():
    if not look.isalpha():
        expected("Name")
    return look

def get_num():
    if not look.isdigit():
        expected("Integer")
    return look

def emit(s):
    print(f"\t{s}", end="")

def emit_ln(s):
    emit(s)
    print()


def init():
    get_char()

def expression():
    emit_ln(f"Move #{get_num()},D0")

init()
expression()



