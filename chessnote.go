// Package chessnote provides a high-performance, production-grade Go library
// for parsing Portable Game Notation (PGN), the universal standard for chess game data.
package chessnote

import (
	"fmt"
	"io"
	"strings"

	"github.com/HexaTech/chessnote/internal/util"
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
// It processes tag pairs and the core movetext, including moves, captures, and checks.
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
	move := Move{}

	// 1. Handle check/mate suffix first
	if strings.HasSuffix(raw, "+") {
		move.IsCheck = true
		raw = strings.TrimSuffix(raw, "+")
	} else if strings.HasSuffix(raw, "#") {
		move.IsMate = true
		raw = strings.TrimSuffix(raw, "#")
	}

	// 2. The last two characters are almost always the destination square.
	if len(raw) < 2 {
		return Move{}, false
	}
	destStr := raw[len(raw)-2:]
	dest, ok := newSquare(destStr)
	if !ok {
		return Move{}, false
	}
	move.To = dest

	pre := raw[:len(raw)-2]

	// 3. If pre is empty, it's a pawn move (e.g., "e4")
	if len(pre) == 0 {
		move.Piece = Pawn
		return move, true
	}

	// 4. Handle piece moves and pawn captures
	firstChar := rune(pre[0])

	if piece, ok := PieceSymbols[firstChar]; ok {
		// It's a piece move
		move.Piece = piece
		disambStr := pre[1:]

		if strings.Contains(disambStr, "x") {
			move.IsCapture = true
			disambStr = strings.Replace(disambStr, "x", "", 1)
		}

		if len(disambStr) == 1 {
			char := rune(disambStr[0])
			if util.IsFile(char) {
				move.From.File = int(char - 'a')
			} else if util.IsRank(char) {
				move.From.Rank = int(char - '1')
			}
		}

		return move, true
	} else if util.IsFile(firstChar) {
		// It's a pawn capture (e.g., "exd5")
		if pre == string(firstChar)+"x" {
			move.Piece = Pawn
			move.IsCapture = true
			move.From.File = int(firstChar - 'a')
			return move, true
		}
	}

	return Move{}, false
}

func newSquare(s string) (Square, bool) {
	if len(s) != 2 || !util.IsFile(rune(s[0])) || !util.IsRank(rune(s[1])) {
		return Square{}, false
	}
	return Square{
		File: int(s[0] - 'a'),
		Rank: int(s[1] - '1'),
	}, true
}

// ParseString is a convenience helper that parses a PGN string.
func ParseString(s string) (*Game, error) {
	p := NewParser(strings.NewReader(s))
	return p.Parse()
}
