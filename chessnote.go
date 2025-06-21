package chessnote

import (
	"fmt"
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
type Parser struct {
	s *Scanner
}

// NewParser creates and returns a new PGN Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse reads PGN data from an io.Reader, parses it, and returns a Game object.
// It processes the tag pairs and will be extended to parse the full movetext.
func (p *Parser) Parse() (*Game, error) {
	game := &Game{
		Tags: make(map[string]string),
	}

	for {
		tok := p.s.Scan()
		switch tok.Type {
		case EOF:
			return game, nil
		case LBRACKET:
			if err := p.parseTagPair(game); err != nil {
				return nil, err
			}
		case IDENT, NUMBER:
			// Once we see an ident or number outside a tag, we are in the movetext.
			if err := p.parseMovetext(tok, game); err != nil {
				return nil, err
			}
		}
	}
}

func (p *Parser) parseTagPair(g *Game) error {
	key := p.s.Scan()
	if key.Type != IDENT {
		return fmt.Errorf("expected ident for tag key, got %v", key)
	}

	value := p.s.Scan()
	if value.Type != STRING {
		return fmt.Errorf("expected string for tag value, got %v", value)
	}

	g.Tags[key.Literal] = value.Literal

	end := p.s.Scan()
	if end.Type != RBRACKET {
		return fmt.Errorf("expected ']' to close tag, got %v", end)
	}
	return nil
}

func (p *Parser) parseMovetext(firstToken Token, g *Game) error {
	// For now, put the first token back and re-scan in a loop.
	// This is not efficient, but it's a simple way to handle the transition.
	// We will build a more robust recursive descent parser later.
	p.s.r.UnreadRune() // This is a hack for now.

	// Re-create a reader from the rest of the stream. This is very inefficient.
	// A proper implementation would use a buffered scanner that can peek/unread.
	// But for this refactoring step, we will get it working first.
	// The rest of the implementation is left as before.

	// The old logic for parsing movetext from raw strings can be adapted here,
	// but it would be better to parse from the token stream.
	// Let's just re-implement the move parsing logic based on tokens.

	for tok := firstToken; tok.Type != EOF && tok.Type != ASTERISK; tok = p.s.Scan() {
		switch tok.Type {
		case IDENT:
			// Create a separate function to parse move text to keep this clean.
			move, ok := p.parseMove(tok.Literal)
			if ok {
				g.Moves = append(g.Moves, move)
			}
		}
	}

	return nil
}

func (p *Parser) parseMove(raw string) (Move, bool) {
	// Pawn captures like "exd5"
	if len(raw) == 4 && raw[1] == 'x' {
		// This could be a piece or a pawn capture. Check if the first char is a file.
		if raw[0] >= 'a' && raw[0] <= 'h' {
			dest := raw[2:]
			if len(dest) == 2 && dest[0] >= 'a' && dest[0] <= 'h' && dest[1] >= '1' && dest[1] <= '8' {
				return Move{
					Piece:     Pawn,
					IsCapture: true,
					From: Square{
						File: int(raw[0] - 'a'),
						// Rank is unknown from this notation
					},
					To: Square{
						File: int(dest[0] - 'a'),
						Rank: int(dest[1] - '1'),
					},
				}, true
			}
		}

		// Piece captures like "Nxf3"
		if piece, ok := PieceSymbols[rune(raw[0])]; ok {
			dest := raw[2:]
			if len(dest) == 2 && dest[0] >= 'a' && dest[0] <= 'h' && dest[1] >= '1' && dest[1] <= '8' {
				return Move{
					Piece:     piece,
					IsCapture: true,
					To: Square{
						File: int(dest[0] - 'a'),
						Rank: int(dest[1] - '1'),
					},
				}, true
			}
		}
	}

	// Piece moves like "Nf3"
	if len(raw) == 3 {
		if piece, ok := PieceSymbols[rune(raw[0])]; ok {
			dest := raw[1:]
			if len(dest) == 2 && dest[0] >= 'a' && dest[0] <= 'h' && dest[1] >= '1' && dest[1] <= '8' {
				return Move{
					Piece: piece,
					To: Square{
						File: int(dest[0] - 'a'),
						Rank: int(dest[1] - '1'),
					},
				}, true
			}
		}
	}

	// Pawn moves like "e4"
	if len(raw) == 2 && raw[0] >= 'a' && raw[0] <= 'h' && raw[1] >= '1' && raw[1] <= '8' {
		return Move{
			Piece: Pawn,
			To: Square{
				File: int(raw[0] - 'a'),
				Rank: int(raw[1] - '1'),
			},
		}, true
	}

	return Move{}, false
}

// ParseString is a convenience helper that parses a PGN string.
func ParseString(s string) (*Game, error) {
	p := NewParser(strings.NewReader(s))
	return p.Parse()
}
