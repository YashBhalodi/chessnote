package chessnote

import (
	"bufio"
	"io"
	"strings"
)

// Game represents a single parsed PGN game, including its tag pairs,
// movetext, and final result.
type Game struct {
	// Tags is a map of PGN tag key-value pairs.
	Tags map[string]string
	// Moves is a slice of moves made in the game.
	Moves []Move
	// Result is the final result of the game (e.g., "1-0", "0-1").
	Result string
}

// Move represents a single move made by one player.
type Move struct {
	// From is the starting square of the move.
	From Square
	// To is the destination square of the move.
	To Square
	// Piece is the type of piece that was moved.
	Piece PieceType
	// IsCapture indicates whether the move was a capture.
	IsCapture bool
	// IsCheck indicates whether the move resulted in a check.
	IsCheck bool
	// IsMate indicates whether the move resulted in a checkmate.
	IsMate bool
}

// Square represents a single square on the board (e.g., e4).
type Square struct {
	// File is the file of the square, represented as 0-7 for files a-h.
	File int
	// Rank is the rank of the square, represented as 0-7 for ranks 1-8.
	Rank int
}

// PieceType defines the type of chess piece.
type PieceType int

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

// PieceSymbols maps a rune to a PieceType.
var PieceSymbols = map[rune]PieceType{
	'N': Knight,
	'B': Bishop,
	'R': Rook,
	'Q': Queen,
	'K': King,
}

// Parser is a PGN parser that reads from an io.Reader and parses it into a Game.
type Parser struct{}

// NewParser creates and returns a new PGN Parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads PGN data from an io.Reader, parses it, and returns a Game object.
// It processes the tag pairs and will be extended to parse the full movetext.
func (p *Parser) Parse(r io.Reader) (*Game, error) {
	game := &Game{
		Tags: make(map[string]string),
	}
	scanner := bufio.NewScanner(r)
	inMovetext := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Skip empty lines
		}

		// Try to parse as a tag pair first
		if !inMovetext && strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			line = strings.TrimPrefix(line, "[")
			line = strings.TrimSuffix(line, "]")
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := strings.Trim(parts[1], "\"")
				game.Tags[key] = value
				continue // Successfully parsed a tag, move to the next line
			}
		}

		// If it's not a tag pair, it must be movetext from here on out.
		inMovetext = true
		tokens := strings.Fields(line)
		for _, token := range tokens {
			// Is it a piece move? (e.g., Nf3)
			if len(token) == 3 {
				if piece, ok := PieceSymbols[rune(token[0])]; ok {
					dest := token[1:]
					if len(dest) == 2 && dest[0] >= 'a' && dest[0] <= 'h' && dest[1] >= '1' && dest[1] <= '8' {
						move := Move{
							Piece: piece,
							To: Square{
								File: int(dest[0] - 'a'),
								Rank: int(dest[1] - '1'),
							},
						}
						game.Moves = append(game.Moves, move)
						continue
					}
				}
			}

			// Is it a pawn move? (e.g., e4)
			if len(token) == 2 && token[0] >= 'a' && token[0] <= 'h' && token[1] >= '1' && token[1] <= '8' {
				move := Move{
					Piece: Pawn, // Assume Pawn for now
					To: Square{
						File: int(token[0] - 'a'),
						Rank: int(token[1] - '1'),
					},
				}
				game.Moves = append(game.Moves, move)
			}
		}
	}
	return game, nil
}

// ParseString is a convenience helper that parses a PGN string.
// It wraps the input string in a strings.Reader and calls Parse.
func (p *Parser) ParseString(s string) (*Game, error) {
	return p.Parse(strings.NewReader(s))
}
