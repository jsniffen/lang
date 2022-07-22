from enum import Enum

class Token(Enum):
    COMMENT = 1
    IDENTIFIER = 2
    FUNCTION = 3
    LBRACE = 4
    RBRACE = 5
    LPAREN = 6
    RPAREN = 7
    STRING = 8

file = open("lang/hello.lang").read()

def eatline(buffer, i):
    j = i
    while j < len(buffer) and buffer[j] != "\n":
        j += 1
    return j-i

def eatstring(buffer, i):
    j = i+1
    while j < len(buffer) and buffer[j] != "\"":
        j += 1
    return j-i

def getword(buffer, i):
    j = i
    chars = []
    while j < len(buffer) and buffer[j].isalpha():
        chars.append(buffer[j])
        j += 1
    return "".join(chars), j-i-1

tokens = []
i = 0
line, column = 1, 1

while i < len(file):
    eat = 0
    newline = False

    char = file[i]

    if char == "\n":
        newline = True
    elif char == "/" and i < len(file)-1 and file[i+1] == "/":
        eat = eatline(file, i)
        tokens.append(Token.COMMENT)
    if char == "{":
        tokens.append(Token.LBRACE)
    elif char == "}":
        tokens.append(Token.RBRACE)
    elif char == "(":
        tokens.append(Token.LPAREN)
    elif char == ")":
        tokens.append(Token.RPAREN)
    elif char == "\"":
        eat = eatstring(file, i)
        tokens.append(Token.STRING)
    elif char.isalpha():
        word, eat = getword(file, i)
        tokens.append(Token.IDENTIFIER)
    if newline:
        column = 1
        line += 1
    else:
        column += 1

    i += 1 + eat

for t in tokens:
    print(t)
