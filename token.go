package chessnote

// TokenType represents a lexical token type.
type TokenType int

// Token represents a lexical token returned by the scanner.
type Token struct {
	// Type is the type of the token.
	Type TokenType
	// Literal is the literal value of the token.
	Literal string
}

const (
	// Special token types
	ILLEGAL TokenType = iota // An illegal token
	EOF                      // End of file

	// Literals
	IDENT   // e.g., Event, White, Nf3, e4, Nxf3+
	COMMENT // e.g., { A comment }
	STRING  // e.g., "F/S Return Match"
	NUMBER  // e.g., 1, 29

	// Punctuation
	LBRACKET // [
	RBRACKET // ]
	LPAREN   // (
	RPAREN   // )
	ASTERISK // *
	DOT      // .
)
