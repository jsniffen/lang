file = open("main.lang").read()

def eatline(buffer, i):
    while i < len(buffer) and buffer[i] != "\n":
        i += 1
    return i+1 if i < len(buffer)-1 else i

i = 0
line, column = 1, 1

while i < len(file):
    char = file[i]

    if char == "/" and i < len(file)-1 and file[i+1] == "/":
        print(f"{line}:{column} COMMENT")
        i = eatline(file, i)
        line += 1
        column = 1
        continue

    if char == "\n":
        line += 1
        column = 1
        i += 1
        continue


    column += 1
    i += 1

