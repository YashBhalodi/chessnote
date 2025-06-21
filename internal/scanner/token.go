package scanner

// TokenType represents the type of a token.
type TokenType int

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = iota // An unknown token
	EOF                      // End of file

	// Literals
	IDENT  // e.g., "Event", "e4", "Nf3"
	NUMBER // e.g., "1", "29"
	STRING // e.g., "F/S Return Match"

	// Delimiters
	LBRACKET // [
	RBRACKET // ]
	LPAREN   // (
	RPAREN   // )
	DOT      // .
	ASTERISK // *

	// Keywords & Special
	COMMENT // A comment block or line
	NAG     // Numeric Annotation Glyph, e.g., $1
)
