package util

// IsFile checks if a rune is a valid file character ('a'-'h').
func IsFile(r rune) bool {
	return r >= 'a' && r <= 'h'
}

// IsRank checks if a rune is a valid rank character ('1'-'8').
func IsRank(r rune) bool {
	return r >= '1' && r <= '8'
}

// IsWhitespace checks if a rune is a whitespace character.
func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

// IsLetter checks if a rune is a letter.
func IsLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// IsDigit checks if a rune is a digit.
func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
