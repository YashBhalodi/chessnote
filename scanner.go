package chessnote

import (
	"bufio"
	"io"
	"strconv"

	"github.com/YashBhalodi/chessnote/internal/util"
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

	if util.IsWhitespace(r) {
		s.unread()
		return s.scanWhitespace()
	} else if util.IsLetter(r) || util.IsDigit(r) {
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
	case '(':
		return Token{Type: LPAREN, Literal: string(r)}
	case ')':
		return Token{Type: RPAREN, Literal: string(r)}
	case '"':
		return s.scanString()
	case '.':
		return Token{Type: DOT, Literal: string(r)}
	case '*':
		return Token{Type: ASTERISK, Literal: string(r)}
	case '{':
		return s.scanCommentBlock()
	case ';':
		return s.scanCommentLine()
	case '$':
		return s.scanNAG()
	}

	return Token{Type: ILLEGAL, Literal: string(r)}
}

func (s *Scanner) scanWhitespace() Token {
	var lit string
	for {
		r := s.read()
		if r == eof {
			break
		} else if !util.IsWhitespace(r) {
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
		} else if !util.IsLetter(r) && !util.IsDigit(r) && r != '_' && r != '+' && r != '#' && r != 'x' && r != '=' && r != '-' {
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

func (s *Scanner) scanCommentBlock() Token {
	var lit string
	for {
		r := s.read()
		if r == '}' || r == eof {
			break
		}
		lit += string(r)
	}
	return Token{Type: COMMENT, Literal: lit}
}

func (s *Scanner) scanCommentLine() Token {
	var lit string
	for {
		r := s.read()
		if r == '\n' || r == eof {
			break
		}
		lit += string(r)
	}
	return Token{Type: COMMENT, Literal: lit}
}

func (s *Scanner) scanNAG() Token {
	var lit string
	for {
		r := s.read()
		if !util.IsDigit(r) {
			s.unread()
			break
		}
		lit += string(r)
	}
	return Token{Type: NAG, Literal: lit}
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
