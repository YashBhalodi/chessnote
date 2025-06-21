package chessnote

import (
	"bufio"
	"io"
	"strconv"
)

// Scanner is responsible for lexical analysis of a PGN input stream.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next PGN token and its literal value.
func (s *Scanner) Scan() Token {
	r := s.read()

	if isWhitespace(r) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(r) || isDigit(r) {
		s.unread()
		return s.scanIdent()
	}

	switch r {
	case eof:
		return Token{Type: EOF}
	case '[':
		return Token{Type: LBRACKET, Literal: string(r)}
	case ']':
		return Token{Type: RBRACKET, Literal: string(r)}
	case '"':
		return s.scanString()
	case '.':
		return Token{Type: DOT, Literal: string(r)}
	case '*':
		return Token{Type: ASTERISK, Literal: string(r)}
	}

	return Token{Type: ILLEGAL, Literal: string(r)}
}

func (s *Scanner) scanWhitespace() Token {
	var lit string
	for {
		r := s.read()
		if r == eof {
			break
		} else if !isWhitespace(r) {
			s.unread()
			break
		}
		lit += string(r)
	}
	// Whitespace is not a token, so we recursively call Scan to get the next one.
	return s.Scan()
}

func (s *Scanner) scanIdent() Token {
	var lit string
	for {
		r := s.read()
		if r == eof {
			break
		} else if !isLetter(r) && !isDigit(r) && r != '_' {
			s.unread()
			break
		}
		lit += string(r)
	}

	if _, err := strconv.Atoi(lit); err == nil {
		return Token{Type: NUMBER, Literal: lit}
	}
	return Token{Type: IDENT, Literal: lit}
}

func (s *Scanner) scanString() Token {
	var lit string
	for {
		r := s.read()
		if r == '"' || r == eof {
			break
		}
		lit += string(r)
	}
	return Token{Type: STRING, Literal: lit}
}

func (s *Scanner) read() rune {
	r, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return r
}

func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

var eof = rune(0)

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
