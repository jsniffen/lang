package ir

import "bytes"

type Writer struct {
	w      bytes.Buffer
	indent int
}

func (w *Writer) String() string {
	return w.w.String()
}

func (w *Writer) Bytes() []byte {
	return w.w.Bytes()
}

func (w *Writer) Write(s string) {
	w.w.WriteString(s)
}

func (w *Writer) NewLine() {
	w.w.WriteString("\n")
	for i := 0; i < w.indent; i += 1 {
		w.w.WriteString("\t")
	}
}

func (w *Writer) Indent() {
	w.indent += 1
}

func (w *Writer) DeIndent() {
	w.indent -= 1
}
