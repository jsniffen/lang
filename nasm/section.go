package nasm

import "bytes"

type Section struct {
	ind int
	out bytes.Buffer
}

func newSection(name string) *Section {
	s := &Section{}

	if name != "" {
		s.write("section .")
		s.write(name)
		s.newLine()
	}

	return s
}

func (s *Section) String() string {
	return s.out.String()
}

func (s *Section) write(str string) {
	s.out.WriteString(str)
}

func (s *Section) indent() {
	s.ind += 1
}

func (s *Section) unindent() {
	s.ind -= 1
}

func (s *Section) newLine() {
	s.out.WriteString("\n")
	for i := 0; i < s.ind; i++ {
		s.out.WriteString("\t")
	}
}
