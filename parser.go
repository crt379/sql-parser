package sqlparser

import (
	"bufio"
	"io"
	"strings"
)

type Parser struct {
	reader *bufio.Reader
}

func New(r io.Reader) Parser {
	return Parser{
		reader: bufio.NewReader(r),
	}
}

type FixedString struct {
	rs  []rune
	len int
	si  int
	ei  int
	b   strings.Builder
}

func NewFixedString(len int) FixedString {
	return FixedString{
		rs:  make([]rune, len),
		len: len,
		si:  0,
		ei:  0,
	}
}

func (s *FixedString) Append(char rune) {
	// 将字符写入当前索引位置
	s.rs[s.ei%s.len] = char
	s.ei++

	// 如果 ei 达到缓冲区末尾，重置 ei 并更新 si
	if s.ei > len(s.rs) {
		s.si++
	}
}

func (s *FixedString) String() string {
	s.b.Reset()
	for i := s.si; i < s.ei; i++ {
		s.b.WriteRune(s.rs[i%s.len])
	}

	return s.b.String()
}

func (p *Parser) Parse() (string, bool) {
	var build strings.Builder
	state := 0
	symbol := NewFixedString(20)
	endsymbol := ";"
	for {
		var c rune
		var err error
		if c, _, err = p.reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			}
		}

		symbol.Append(c)
		char := string(c)
		switch state {
		case 0:
			if char == ";" || char == "\n" || char == " " {
				continue
			}

			if char == "-" {
				state = 30
				continue
			}

			if char == "/" {
				state = 20
				continue
			}

			build.WriteString(char)
			state = 1
			continue
		case 1:
			build.WriteString(char)
			if char == "'" {
				state = 10
				continue
			}

			symbolstr := symbol.String()
			if endsymbol == ";" && (strings.HasSuffix(symbolstr, "DELIMITER") || strings.HasSuffix(symbolstr, "RETURNS TRIGGER AS")) {
				state = 40
				continue
			}

			if strings.HasSuffix(symbol.String(), endsymbol) {
				return build.String(), false
			}
		case 10:
			build.WriteString(char)
			if char == "'" {
				state = 1
				continue
			}
		case 20:
			if char == "*" {
				state = 21
				continue
			}
		case 21:
			if char == "*" {
				state = 22
				continue
			}
		case 22:
			if char == "/" {
				state = 0
				continue
			}
		case 30:
			if char == "-" {
				state = 31
				continue
			}
		case 31:
			if char == "\n" {
				state = 0
				continue
			}
		case 40:
			build.WriteString(char)
			if char != " " {
				endsymbol = char
				state = 41
				continue
			}
		case 41:
			build.WriteString(char)
			if char == "\n" || char == " " {
				state = 1
				continue
			}

			endsymbol += char
		}

	}

	return build.String(), true
}
