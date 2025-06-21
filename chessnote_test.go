package chessnote_test

import (
	"reflect"
	"testing"

	"github.com/HexaTech/chessnote"
)

func TestParseTagPairs(t *testing.T) {
	t.Parallel()
	pgn := `
[Event "F/S Return Match"]
[Site "Belgrade, Serbia JUG"]
[Date "1992.11.04"]
[Round "29"]
[White "Fischer, Robert J."]
[Black "Spassky, Boris V."]
[Result "1/2-1/2"]
`
	game, err := chessnote.ParseString(pgn)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	want := map[string]string{
		"Event":  "F/S Return Match",
		"Site":   "Belgrade, Serbia JUG",
		"Date":   "1992.11.04",
		"Round":  "29",
		"White":  "Fischer, Robert J.",
		"Black":  "Spassky, Boris V.",
		"Result": "1/2-1/2",
	}

	if !reflect.DeepEqual(game.Tags, want) {
		t.Errorf("Parse() got = %v, want %v", game.Tags, want)
	}
}

func TestParseMoves(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		pgn  string
		want chessnote.Move
	}{
		// Pawn Moves
		{
			name: "pawn move",
			pgn:  "1. e4 *",
			want: chessnote.Move{Piece: chessnote.Pawn, To: chessnote.Square{File: 4, Rank: 3}},
		},
		{
			name: "pawn move with check",
			pgn:  "1. e4+ *",
			want: chessnote.Move{Piece: chessnote.Pawn, To: chessnote.Square{File: 4, Rank: 3}, IsCheck: true},
		},
		{
			name: "pawn move with mate",
			pgn:  "1. e4# *",
			want: chessnote.Move{Piece: chessnote.Pawn, To: chessnote.Square{File: 4, Rank: 3}, IsMate: true},
		},
		// Piece Moves
		{
			name: "piece move",
			pgn:  "1. Nf3 *",
			want: chessnote.Move{Piece: chessnote.Knight, To: chessnote.Square{File: 5, Rank: 2}},
		},
		// Captures
		{
			name: "piece capture",
			pgn:  "1. Nxf3 *",
			want: chessnote.Move{Piece: chessnote.Knight, To: chessnote.Square{File: 5, Rank: 2}, IsCapture: true},
		},
		{
			name: "pawn capture",
			pgn:  "1. exd5 *",
			want: chessnote.Move{Piece: chessnote.Pawn, From: chessnote.Square{File: 4}, To: chessnote.Square{File: 3, Rank: 4}, IsCapture: true},
		},
		// Disambiguation
		{
			name: "file disambiguation",
			pgn:  "1. Rdf8 *",
			want: chessnote.Move{Piece: chessnote.Rook, From: chessnote.Square{File: 3}, To: chessnote.Square{File: 5, Rank: 7}},
		},
		{
			name: "rank disambiguation",
			pgn:  "1. N1c3 *",
			want: chessnote.Move{Piece: chessnote.Knight, From: chessnote.Square{Rank: 0}, To: chessnote.Square{File: 2, Rank: 2}},
		},
		{
			name: "file disambiguation with capture",
			pgn:  "1. Rdxf8 *",
			want: chessnote.Move{Piece: chessnote.Rook, From: chessnote.Square{File: 3}, To: chessnote.Square{File: 5, Rank: 7}, IsCapture: true},
		},
		{
			name: "rank disambiguation with capture",
			pgn:  "1. N1xc3 *",
			want: chessnote.Move{Piece: chessnote.Knight, From: chessnote.Square{Rank: 0}, To: chessnote.Square{File: 2, Rank: 2}, IsCapture: true},
		},
		// Promotion
		{
			name: "simple promotion",
			pgn:  "1. e8=Q *",
			want: chessnote.Move{Piece: chessnote.Pawn, To: chessnote.Square{File: 4, Rank: 7}, Promotion: chessnote.Queen},
		},
		{
			name: "promotion with capture",
			pgn:  "1. exd8=R *",
			want: chessnote.Move{Piece: chessnote.Pawn, From: chessnote.Square{File: 4}, To: chessnote.Square{File: 3, Rank: 7}, IsCapture: true, Promotion: chessnote.Rook},
		},
		{
			name: "promotion with check",
			pgn:  "1. e8=Q+ *",
			want: chessnote.Move{Piece: chessnote.Pawn, To: chessnote.Square{File: 4, Rank: 7}, Promotion: chessnote.Queen, IsCheck: true},
		},
		// Castling
		{
			name: "kingside castle",
			pgn:  "1. O-O *",
			want: chessnote.Move{Piece: chessnote.King, IsKingsideCastle: true},
		},
		{
			name: "queenside castle",
			pgn:  "1. O-O-O *",
			want: chessnote.Move{Piece: chessnote.King, IsQueensideCastle: true},
		},
		{
			name: "kingside castle with check",
			pgn:  "1. O-O+ *",
			want: chessnote.Move{Piece: chessnote.King, IsKingsideCastle: true, IsCheck: true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game, err := chessnote.ParseString(tc.pgn)
			if err != nil {
				t.Fatalf("Parse() error = %v, for PGN: %s", err, tc.pgn)
			}
			if len(game.Moves) != 1 {
				t.Fatalf("expected 1 move, got %d", len(game.Moves))
			}
			got := game.Moves[0]

			// For ambiguous moves, we only care about the part that was specified.
			// The reflect.DeepEqual check below will handle the rest of the fields.
			if tc.name == "pawn capture" || tc.name == "file disambiguation" || tc.name == "file disambiguation with capture" || tc.name == "promotion with capture" {
				if got.From.File != tc.want.From.File {
					t.Errorf("Parse() got From.File = %d, want %d", got.From.File, tc.want.From.File)
				}
				// Clear the rank for the final DeepEqual check since it's ambiguous.
				got.From.Rank = 0
				tc.want.From.Rank = 0
			}
			if tc.name == "rank disambiguation" || tc.name == "rank disambiguation with capture" {
				if got.From.Rank != tc.want.From.Rank {
					t.Errorf("Parse() got From.Rank = %d, want %d", got.From.Rank, tc.want.From.Rank)
				}
				// Clear the file for the final DeepEqual check since it's ambiguous.
				got.From.File = 0
				tc.want.From.File = 0
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Parse() got = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestParseWithComments(t *testing.T) {
	t.Parallel()
	pgn := `
[Event "Test Game"]
[Result "1-0"]

{This is a comment at the start.}
1. e4 ; This is a move comment
{This is a comment between moves.} e5
2. Nf3 Nc6 1-0
`
	game, err := chessnote.ParseString(pgn)
	if err != nil {
		t.Fatalf("ParseString() failed: %v", err)
	}

	if len(game.Moves) != 4 {
		t.Fatalf("expected 4 moves, got %d", len(game.Moves))
	}

	expectedTags := map[string]string{
		"Event":  "Test Game",
		"Result": "1-0",
	}
	if !reflect.DeepEqual(game.Tags, expectedTags) {
		t.Errorf("got tags %v, want %v", game.Tags, expectedTags)
	}

	expectedMoves := []chessnote.Move{
		{Piece: chessnote.Pawn, To: chessnote.Square{File: 4, Rank: 3}},
		{Piece: chessnote.Pawn, To: chessnote.Square{File: 4, Rank: 4}},
		{Piece: chessnote.Knight, To: chessnote.Square{File: 5, Rank: 2}},
		{Piece: chessnote.Knight, To: chessnote.Square{File: 2, Rank: 5}},
	}

	if !reflect.DeepEqual(game.Moves, expectedMoves) {
		t.Errorf("got moves\n%v\nwant\n%v", game.Moves, expectedMoves)
	}

	if game.Result != "1-0" {
		t.Errorf("got result %q, want %q", game.Result, "1-0")
	}
}

func TestParseWithRAV(t *testing.T) {
	t.Parallel()
	pgn := `1. e4 (1. d4) 1... e5 *`
	game, err := chessnote.ParseString(pgn)
	if err != nil {
		t.Fatalf("ParseString() failed: %v", err)
	}

	if len(game.Moves) != 2 {
		t.Fatalf("expected 2 moves, got %d", len(game.Moves))
	}

	mainMove := game.Moves[0]
	if mainMove.Piece != chessnote.Pawn || mainMove.To != (chessnote.Square{File: 4, Rank: 3}) {
		t.Errorf("unexpected main move: got %+v", mainMove)
	}

	if len(mainMove.Variations) != 1 {
		t.Fatalf("expected 1 variation, got %d", len(mainMove.Variations))
	}

	variation := mainMove.Variations[0]
	if len(variation) != 1 {
		t.Fatalf("expected 1 move in variation, got %d", len(variation))
	}

	variationMove := variation[0]
	expectedVariationMove := chessnote.Move{Piece: chessnote.Pawn, To: chessnote.Square{File: 3, Rank: 3}}
	if !reflect.DeepEqual(variationMove, expectedVariationMove) {
		t.Errorf("got variation move %+v, want %+v", variationMove, expectedVariationMove)
	}
}

func TestParseWithNAGs(t *testing.T) {
	t.Parallel()
	pgn := `1. e4 $1 1... e5 $2 $18 *`
	game, err := chessnote.ParseString(pgn)
	if err != nil {
		t.Fatalf("ParseString() failed: %v", err)
	}

	if len(game.Moves) != 2 {
		t.Fatalf("expected 2 moves, got %d", len(game.Moves))
	}

	move1 := game.Moves[0]
	want1 := []int{1}
	if !reflect.DeepEqual(move1.NAGs, want1) {
		t.Errorf("move 1: got NAGs %v, want %v", move1.NAGs, want1)
	}

	move2 := game.Moves[1]
	want2 := []int{2, 18}
	if !reflect.DeepEqual(move2.NAGs, want2) {
		t.Errorf("move 2: got NAGs %v, want %v", move2.NAGs, want2)
	}
}

// newSquare is an unexported function, so we test it here in the same package.
func TestNewSquare(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		s      string
		want   chessnote.Square
		wantOk bool
	}{
		{"valid square", "e4", chessnote.Square{File: 4, Rank: 3}, true},
		{"invalid file", "z4", chessnote.Square{}, false},
		{"invalid rank", "e9", chessnote.Square{}, false},
		{"too short", "e", chessnote.Square{}, false},
		{"too long", "e4e5", chessnote.Square{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a bit of a hack to test an unexported function.
			// In a real project, you might expose this as a public utility
			// if it were needed, or keep it tested implicitly via the parser.
			// For this case, direct testing is clearer.
			// We can't call chessnote.newSquare directly, so we parse a move.
			game, err := chessnote.ParseString(tt.s)
			if err != nil && tt.wantOk {
				t.Fatalf("ParseString() error = %v", err)
			}

			var got chessnote.Square
			var gotOk bool
			if len(game.Moves) == 1 {
				got = game.Moves[0].To
				gotOk = true
			}

			if gotOk != tt.wantOk {
				t.Fatalf("newSquare() ok = %v, want %v", gotOk, tt.wantOk)
			}
			if gotOk && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSquare() got = %v, want %v", got, tt.want)
			}
		})
	}
}
