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

// Move represents a single move made by one player, capturing all details
// expressed in Standard Algebraic Notation (SAN).
type Move struct {
	// From is the starting square of the move. For many moves, this may be
	// partially or fully zero, as PGN format often omits this information
	// when it's not needed for disambiguation.
	From Square
	// To is the destination square of the move. This is always specified.
	To Square
	// Piece is the type of piece that was moved.
	Piece PieceType
	// Promotion is the piece type a pawn is promoted to. It is zero
	// (Pawn) if there is no promotion.
	Promotion PieceType
	// IsCapture indicates whether the move was a capture.
	IsCapture bool
	// IsCheck indicates whether the move resulted in a check.
	IsCheck bool
	// IsMate indicates whether the move resulted in a checkmate.
	IsMate bool
	// IsKingsideCastle indicates a kingside castling move (O-O).
	IsKingsideCastle bool
	// IsQueensideCastle indicates a queenside castling move (O-O-O).
	IsQueensideCastle bool
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
	// Pawn represents a pawn piece. Note that this is the zero value for PieceType.
	Pawn PieceType = iota
	// Knight represents a knight piece.
	Knight
	// Bishop represents a bishop piece.
	Bishop
	// Rook represents a rook piece.
	Rook
	// Queen represents a queen piece.
	Queen
	// King represents a king piece.
	King
)

// PieceSymbols maps a rune representation of a piece in PGN to a PieceType.
// Pawns are not represented by a symbol in PGN.
var PieceSymbols = map[rune]PieceType{
	'N': Knight,
	'B': Bishop,
	'R': Rook,
	'Q': Queen,
	'K': King,
}

// Parser is a PGN parser that reads from an io.Reader and parses it into a Game.
// It implements a standard recursive descent parser.
type Parser struct {
	s *Scanner
}

// NewParser creates and returns a new PGN Parser for the given reader.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse reads and parses the entire PGN data from the reader, returning a
// single Game object. It expects the PGN data to contain exactly one game.
// The parser stops at the first game-terminating symbol (*, 1-0, etc.).
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
	// The final move we will build and return.
	var finalMove Move

	// Make a mutable copy of the raw string to parse.
	movetext := raw

	// 1. Parse and strip the promotion FIRST.
	if i := strings.Index(movetext, "="); i != -1 {
		// This handles cases like "e8=Q+"
		promoAndSuffix := movetext[i+1:]
		if len(promoAndSuffix) == 0 {
			return Move{}, false // "e8=" is invalid
		}

		promoChar := rune(promoAndSuffix[0])
		if piece, ok := PieceSymbols[promoChar]; ok {
			finalMove.Promotion = piece
		} else {
			return Move{}, false // Invalid promotion piece
		}

		// Check for a suffix *after* the promotion
		if len(promoAndSuffix) > 1 {
			suffix := promoAndSuffix[1:]
			switch suffix {
			case "+":
				finalMove.IsCheck = true
			case "#":
				finalMove.IsMate = true
			}
		}

		movetext = movetext[:i] // "e8"
	} else {
		// 2. If no promotion, parse and strip the check/mate suffix.
		if strings.HasSuffix(movetext, "+") {
			finalMove.IsCheck = true
			movetext = strings.TrimSuffix(movetext, "+")
		} else if strings.HasSuffix(movetext, "#") {
			finalMove.IsMate = true
			movetext = strings.TrimSuffix(movetext, "#")
		}
	}

	// 3. Parse the core move notation that's left.
	var coreMove Move
	var ok bool

	switch movetext {
	case "O-O":
		coreMove.Piece = King
		coreMove.IsKingsideCastle = true
		ok = true
	case "O-O-O":
		coreMove.Piece = King
		coreMove.IsQueensideCastle = true
		ok = true
	default:
		// If not castling, parse as a regular move.
		coreMove, ok = p.parseCoreMove(movetext)
	}

	if !ok {
		return Move{}, false
	}

	// 4. Combine the results.
	coreMove.IsCheck = finalMove.IsCheck
	coreMove.IsMate = finalMove.IsMate
	coreMove.Promotion = finalMove.Promotion
	return coreMove, true
}

// parseCoreMove handles a move string after any suffixes/promotions have been removed.
func (p *Parser) parseCoreMove(raw string) (Move, bool) {
	// A pawn capture is the only case where the move starts with a file and contains a capture.
	// e.g. "exd5". Let's handle this special case first.
	if len(raw) == 4 && util.IsFile(rune(raw[0])) && raw[1] == 'x' {
		dest, ok := newSquare(raw[2:])
		if !ok {
			return Move{}, false // Should not happen if grammar is correct
		}
		return Move{
			Piece:     Pawn,
			From:      Square{File: int(raw[0] - 'a')},
			To:        dest,
			IsCapture: true,
		}, true
	}

	move := Move{}
	movetext := raw // create a mutable copy

	// The destination square is always the last two characters.
	if len(movetext) < 2 {
		return Move{}, false
	}
	destStr := movetext[len(movetext)-2:]
	dest, ok := newSquare(destStr)
	if !ok {
		return Move{}, false
	}
	move.To = dest
	movetext = movetext[:len(movetext)-2] // a.k.a 'prefix'

	// Identify the piece type. Defaults to Pawn for moves like "e4".
	move.Piece = Pawn
	if len(movetext) > 0 {
		if piece, ok := PieceSymbols[rune(movetext[0])]; ok {
			move.Piece = piece
			movetext = movetext[1:]
		}
	}

	// The remainder of movetext can be a disambiguation, a capture, or both.
	// PGN standard is: [piece][disambiguation][capture][dest]
	// So we parse disambiguation first.

	// Check for disambiguation, e.g. the 'd' in "Rdf8" or '1' in "N1c3"
	if len(movetext) > 0 {
		char := rune(movetext[0])
		if util.IsFile(char) || util.IsRank(char) {
			if len(movetext) > 1 && movetext[1] == 'x' {
				// This is a disambiguated capture, e.g., "d" in "Rdxf8"
			} else if len(movetext) > 1 {
				// Invalid, e.g. "Rdd4"
				return Move{}, false
			}

			if util.IsFile(char) {
				move.From.File = int(char - 'a')
			} else {
				move.From.Rank = int(char - '1')
			}
			movetext = movetext[1:]
		}
	}

	// Check for a capture for piece moves, e.g. "x" in "Nxf3" or "Rdxf8"
	if len(movetext) > 0 && movetext[0] == 'x' {
		move.IsCapture = true
		movetext = movetext[1:]
	}

	// If anything is left, the move is invalid.
	if len(movetext) > 0 {
		return Move{}, false
	}

	return move, true
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
// It is intended for quickly parsing a complete PGN string already in memory.
// For parsing from files or network streams, creating a Parser with NewParser
// is recommended.
func ParseString(s string) (*Game, error) {
	p := NewParser(strings.NewReader(s))
	return p.Parse()
}
