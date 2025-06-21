package chessnote

import (
	"bufio"
	"io"
	"strings"
)

// Game represents a single parsed PGN game.
type Game struct {
	Tags   map[string]string
	Moves  []Move
	Result string
}

// Move represents a single move by one player.
type Move struct {
	From      Square
	To        Square
	Piece     PieceType
	IsCapture bool
	IsCheck   bool
	IsMate    bool
}

// Square represents a single square on the board (e.g., e4).
type Square struct {
	File int // 0-7 for a-h
	Rank int // 0-7 for 1-8
}

// PieceType defines the type of chess piece.
type PieceType int

// Parser is a PGN parser.
type Parser struct {
	// TODO: Add fields for the parser state
}

// NewParser creates a new PGN parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads from an io.Reader and returns a parsed Game.
func (p *Parser) Parse(r io.Reader) (*Game, error) {
	game := &Game{
		Tags: make(map[string]string),
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			line = strings.TrimPrefix(line, "[")
			line = strings.TrimSuffix(line, "]")
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := strings.Trim(parts[1], "\"")
				game.Tags[key] = value
			}
		}
	}
	return game, nil
}

// ParseString is a helper to parse a PGN string.
func (p *Parser) ParseString(s string) (*Game, error) {
	return p.Parse(strings.NewReader(s))
} 