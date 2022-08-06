#include <stdio.h>

char look;

void nextchar()
{
  look = getchar();
}

void error(char *s)
{
  printf("\nError: %s.\n", s);
}

void fail(char *s)
{
  error(s);
  exit(1);
}

void expected(char *s)
{
  char buf[256];
  sprintf(buf, "%s: Expected", s);
  fail(buf);
}

void match(char c)
{
  if (look == c) {
    nextchar();
  } else {
    char buf[256];
    sprintf(buf, "'%c'", c);
    expected(buf);
  }
}

char getname()
{
  char c;
  if (!isalpha(look)) {
    expected("Name");
  }
  c = look;
  nextchar();
  return c;
}

char getnum()
{
  char c;
  if (!isdigit(look)) {
    expected("Integer");
  }
  c = look;
  nextchar();
  return c;
}

void emit(char *s)
{
  printf("\t%s", s);
}

void emitln(char *s)
{
  emit(s);
  printf("\n");
}

void init()
{
  nextchar();
}

void term()
{
  char buf[256];
  sprintf(buf, "MOVE #%c,D0", getnum());
  emitln(buf);
}

void add()
{
  match('+');
  term();
  emitln("ADD D1,D0");
}

void subtract()
{
  match('-');
  term();
  emitln("SUB D1,D0");
  emitln("NEG D0");
}

void expression()
{
  term();
  while (look == '+' || look == '-') {
    emitln("Move D0,D1");
    switch (look) {
      case '+':
        add();
        break;
      case '-':
        subtract();
        break;
      default:
        expected("Addop");
        break;
    }
  }
}

int main()
{
  init();
  expression();
}
