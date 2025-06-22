// Package chessnote provides a high-performance, production-grade Go library
// for parsing Portable Game Notation (PGN), the universal standard for chess game data.
package chessnote

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/YashBhalodi/chessnote/internal/scanner"
	"github.com/YashBhalodi/chessnote/internal/util"
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
	// Variations lists any alternative move sequences that could have been
	// played. This is used for representing Recursive Annotation Variations (RAVs).
	Variations [][]Move
	// NAGs is a slice of Numeric Annotation Glyphs (e.g., $1, $2)
	// associated with the move.
	NAGs []int
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

// ParserConfig holds configuration settings for the parser.
type ParserConfig struct {
	// Strict mode requires that a PGN game must end with a valid result token
	// (*, 1-0, 0-1, or 1/2-1/2). When disabled (lax mode), a game can end
	// at the end of the file without a result token.
	// It is enabled by default.
	Strict bool
}

// A ParserOption configures a Parser.
type ParserOption func(*ParserConfig)

// WithLaxParsing returns a ParserOption that disables strict parsing mode.
// In lax mode, the parser will not require a final game result token and will
// successfully parse a game that ends abruptly at the end of the file.
func WithLaxParsing() ParserOption {
	return func(c *ParserConfig) {
		c.Strict = false
	}
}

// Parser is a PGN parser that reads from an io.Reader and parses it into a Game.
// It implements a standard recursive descent parser.
type Parser struct {
	s      *scanner.Scanner
	tok    scanner.Token // The current token
	config ParserConfig
}

// NewParser creates and returns a new PGN Parser for the given reader.
// By default, it operates in strict mode. Behavior can be customized with
// ParserOptions, such as WithLaxParsing().
func NewParser(r io.Reader, opts ...ParserOption) *Parser {
	// Default configuration
	config := ParserConfig{
		Strict: true,
	}

	// Apply all options
	for _, opt := range opts {
		opt(&config)
	}

	p := &Parser{
		s:      scanner.NewScanner(r),
		config: config,
	}
	p.scan() // Initialize the first token
	return p
}

// scan moves to the next token and sets it as the parser's current token.
func (p *Parser) scan() {
	p.tok = p.s.Scan()
}

// Parse reads and parses the entire PGN data from the reader, returning a
// single Game object. It expects the PGN data to contain exactly one game.
// The parser stops at the first game-terminating symbol (*, 1-0, etc.).
func (p *Parser) Parse() (*Game, error) {
	game := &Game{
		Tags: make(map[string]string),
	}

	for {
		switch p.tok.Type {
		case scanner.EOF:
			// In strict mode, a game must end with a result token.
			// Reaching EOF without one is an error.
			if p.config.Strict && len(game.Moves) > 0 {
				return nil, fmt.Errorf("unexpected EOF: game must end with a result token")
			}
			return game, nil
		case scanner.LBRACKET:
			// If we are already parsing moves and see a new tag, the game has ended
			// without a result marker.
			if len(game.Moves) > 0 {
				return game, nil
			}
			if err := p.parseTagPair(game); err != nil {
				return nil, err
			}
		case scanner.COMMENT:
			p.scan() // Ignore comments
		case scanner.IDENT, scanner.NUMBER:
			// Once we see an ident or number outside a tag, we are in the movetext.
			if err := p.parseMovetext(&game.Moves); err != nil {
				return nil, err
			}
			// After parsing movetext, we might have a result token.
			if isResult(p.tok) {
				game.Result = p.tok.Literal
			} else if p.config.Strict {
				// If we finish parsing moves and don't have a result, it's an error in strict mode.
				return nil, fmt.Errorf("game must end with a result token, got %v", p.tok)
			}
			return game, nil
		default:
			return nil, fmt.Errorf("unexpected token at start of game: %v", p.tok)
		}
	}
}

func (p *Parser) parseTagPair(g *Game) error {
	p.scan() // Consume '['
	key := p.tok
	if key.Type != scanner.IDENT {
		return fmt.Errorf("expected ident for tag key, got %v", key)
	}

	p.scan() // Consume key
	value := p.tok
	if value.Type != scanner.STRING {
		return fmt.Errorf("expected string for tag value, got %v", value)
	}
	g.Tags[key.Literal] = value.Literal

	p.scan() // Consume value
	if p.tok.Type != scanner.RBRACKET {
		return fmt.Errorf("expected ']' to close tag, got %v", p.tok)
	}
	p.scan() // Consume ']'
	return nil
}

func (p *Parser) parseMovetext(moves *[]Move) error {
	for {
		switch p.tok.Type {
		case scanner.EOF, scanner.ASTERISK, scanner.RPAREN, scanner.LBRACKET:
			return nil // Let caller handle termination
		case scanner.IDENT:
			if isResult(p.tok) {
				return nil // Let caller handle result
			}
			move, err := p.parseMove()
			if err != nil {
				return err
			}
			*moves = append(*moves, move)
		case scanner.NAG:
			if len(*moves) == 0 {
				return fmt.Errorf("found NAG before any moves")
			}
			lastMove := &(*moves)[len(*moves)-1]
			nag, err := strconv.Atoi(p.tok.Literal)
			if err != nil {
				// This should not happen if the scanner is correct.
				return fmt.Errorf("invalid NAG value: %v", p.tok.Literal)
			}
			if lastMove.NAGs == nil {
				lastMove.NAGs = make([]int, 0)
			}
			lastMove.NAGs = append(lastMove.NAGs, nag)
			p.scan()
		case scanner.NUMBER, scanner.DOT, scanner.COMMENT:
			p.scan() // Ignore
		case scanner.LPAREN:
			if len(*moves) == 0 {
				return fmt.Errorf("found variation before any moves")
			}
			lastMove := &(*moves)[len(*moves)-1]
			if err := p.parseRAV(lastMove); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected token in movetext: %v", p.tok)
		}
	}
}

func (p *Parser) parseRAV(parentMove *Move) error {
	p.scan() // Consume '('
	var variationMoves []Move
	if err := p.parseMovetext(&variationMoves); err != nil {
		return err
	}

	if p.tok.Type != scanner.RPAREN {
		return fmt.Errorf("expected ')' to close variation, got %v", p.tok)
	}
	p.scan() // Consume ')'

	if parentMove.Variations == nil {
		parentMove.Variations = make([][]Move, 0)
	}
	parentMove.Variations = append(parentMove.Variations, variationMoves)
	return nil
}

func isResult(tok scanner.Token) bool {
	if tok.Type == scanner.ASTERISK {
		return true
	}
	if tok.Type == scanner.IDENT && (tok.Literal == "1-0" || tok.Literal == "0-1" || tok.Literal == "1/2-1/2") {
		return true
	}
	return false
}

func (p *Parser) parseMove() (Move, error) {
	raw := p.tok.Literal
	move, ok := p.parseMoveFromRaw(raw)
	if !ok {
		return Move{}, fmt.Errorf("invalid move: %s", raw)
	}

	p.scan() // Consume the move token.
	return move, nil
}

// parseMoveFromRaw is the old implementation that works on a string.
// We will phase this out.
func (p *Parser) parseMoveFromRaw(raw string) (Move, bool) {
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

	// Identify and parse the rest of the move components from the prefix.
	movetext, move.Piece = parsePiece(movetext)
	movetext, fromSquare := parseDisambiguation(movetext)
	move.From = fromSquare

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

func parsePiece(movetext string) (string, PieceType) {
	if len(movetext) > 0 {
		if piece, ok := PieceSymbols[rune(movetext[0])]; ok {
			return movetext[1:], piece
		}
	}
	return movetext, Pawn
}

func parseDisambiguation(movetext string) (string, Square) {
	from := Square{}
	if len(movetext) == 0 {
		return movetext, from
	}

	// It can't be a capture 'x' at this stage. If it is, it's part of
	// the next parsing step.
	if movetext[0] == 'x' {
		return movetext, from
	}

	// Disambiguation can be one char (file or rank) or two chars (file and rank).
	// But we don't handle the two-char case yet (e.g. "R1a2").
	char := rune(movetext[0])
	if util.IsFile(char) {
		from.File = int(char - 'a')
		return movetext[1:], from
	} else if util.IsRank(char) {
		from.Rank = int(char - '1')
		return movetext[1:], from
	}

	return movetext, from
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
func ParseString(s string, opts ...ParserOption) (*Game, error) {
	// Trim the UTF-8 Byte Order Mark (BOM) if it exists. Some PGN files
	// may contain this, which can cause parsing to fail.
	s = strings.TrimPrefix(s, "\uFEFF")
	p := NewParser(strings.NewReader(s), opts...)
	return p.Parse()
}
