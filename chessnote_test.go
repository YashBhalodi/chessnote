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

func TestParsePawnMove(t *testing.T) {
	t.Parallel()
	pgn := `1. e4 *`
	game, err := chessnote.ParseString(pgn)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	want := []chessnote.Move{
		{
			To: chessnote.Square{File: 4, Rank: 3}, // e4
		},
	}

	if !reflect.DeepEqual(game.Moves, want) {
		t.Errorf("Parse() got = %v, want %v", game.Moves, want)
	}
}

func TestParsePieceMove(t *testing.T) {
	t.Parallel()
	pgn := `1. Nf3 *`
	game, err := chessnote.ParseString(pgn)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	want := []chessnote.Move{
		{
			Piece: chessnote.Knight,
			To:    chessnote.Square{File: 5, Rank: 2}, // f3
		},
	}

	if !reflect.DeepEqual(game.Moves, want) {
		t.Errorf("Parse() got = %+v, want %+v", game.Moves, want)
	}
}

func TestParseCapture(t *testing.T) {
	t.Parallel()
	pgn := `1. Nxf3 *`
	game, err := chessnote.ParseString(pgn)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	want := []chessnote.Move{
		{
			Piece:     chessnote.Knight,
			To:        chessnote.Square{File: 5, Rank: 2}, // f3
			IsCapture: true,
		},
	}

	if !reflect.DeepEqual(game.Moves, want) {
		t.Errorf("Parse() got = %+v, want %+v", game.Moves, want)
	}
}

func TestParsePawnCapture(t *testing.T) {
	t.Parallel()
	pgn := `1. exd5 *`
	game, err := chessnote.ParseString(pgn)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	want := []chessnote.Move{
		{
			Piece:     chessnote.Pawn,
			From:      chessnote.Square{File: 4},          // From 'e' file
			To:        chessnote.Square{File: 3, Rank: 4}, // To 'd5'
			IsCapture: true,
		},
	}

	// Custom comparison to ignore the From.Rank, which is ambiguous for this move.
	if len(game.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(game.Moves))
	}
	got := game.Moves[0]
	if got.Piece != want[0].Piece ||
		got.From.File != want[0].From.File ||
		got.To != want[0].To ||
		got.IsCapture != want[0].IsCapture {
		t.Errorf("Parse() got = %+v, want %+v", got, want[0])
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

func TestParseCheck(t *testing.T) {
	t.Parallel()
	pgn := `1. e4+ *`
	game, err := chessnote.ParseString(pgn)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	want := []chessnote.Move{
		{
			Piece:   chessnote.Pawn,
			To:      chessnote.Square{File: 4, Rank: 3}, // e4
			IsCheck: true,
		},
	}

	if !reflect.DeepEqual(game.Moves, want) {
		t.Errorf("Parse() got = %+v, want %+v", game.Moves, want)
	}
}

func TestParseCheckmate(t *testing.T) {
	t.Parallel()
	pgn := `1. e4# *`
	game, err := chessnote.ParseString(pgn)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	want := []chessnote.Move{
		{
			Piece:  chessnote.Pawn,
			To:     chessnote.Square{File: 4, Rank: 3}, // e4
			IsMate: true,
		},
	}

	if !reflect.DeepEqual(game.Moves, want) {
		t.Errorf("Parse() got = %+v, want %+v", game.Moves, want)
	}
}

func TestParseDisambiguation(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		pgn  string
		want chessnote.Move
	}{
		{
			name: "file disambiguation",
			pgn:  "1. Rdf8 *",
			want: chessnote.Move{
				Piece: chessnote.Rook,
				From:  chessnote.Square{File: 3},          // 'd' file
				To:    chessnote.Square{File: 5, Rank: 7}, // f8
			},
		},
		{
			name: "rank disambiguation",
			pgn:  "1. N1c3 *",
			want: chessnote.Move{
				Piece: chessnote.Knight,
				From:  chessnote.Square{Rank: 0},          // '1' rank
				To:    chessnote.Square{File: 2, Rank: 2}, // c3
			},
		},
		{
			name: "file disambiguation with capture",
			pgn:  "1. Rdxf8 *",
			want: chessnote.Move{
				Piece:     chessnote.Rook,
				From:      chessnote.Square{File: 3},          // 'd' file
				To:        chessnote.Square{File: 5, Rank: 7}, // f8
				IsCapture: true,
			},
		},
		{
			name: "rank disambiguation with capture",
			pgn:  "1. N1xc3 *",
			want: chessnote.Move{
				Piece:     chessnote.Knight,
				From:      chessnote.Square{Rank: 0},          // '1' rank
				To:        chessnote.Square{File: 2, Rank: 2}, // c3
				IsCapture: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game, err := chessnote.ParseString(tc.pgn)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if len(game.Moves) != 1 {
				t.Fatalf("expected 1 move, got %d", len(game.Moves))
			}

			got := game.Moves[0]

			// Custom comparison to ignore ambiguous parts of From square
			if got.Piece != tc.want.Piece ||
				got.To != tc.want.To ||
				got.IsCapture != tc.want.IsCapture {
				t.Errorf("Parse() got basic fields = %+v, want %+v", got, tc.want)
			}

			// a bit verbose, but clear.
			switch tc.name {
			case "file disambiguation", "file disambiguation with capture":
				if got.From.File != tc.want.From.File {
					t.Errorf("Parse() got From.File = %d, want %d", got.From.File, tc.want.From.File)
				}
			case "rank disambiguation", "rank disambiguation with capture":
				if got.From.Rank != tc.want.From.Rank {
					t.Errorf("Parse() got From.Rank = %d, want %d", got.From.Rank, tc.want.From.Rank)
				}
			}
		})
	}
}

func TestParsePromotion(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		pgn  string
		want chessnote.Move
	}{
		{
			name: "simple promotion",
			pgn:  "1. e8=Q *",
			want: chessnote.Move{
				Piece:     chessnote.Pawn,
				To:        chessnote.Square{File: 4, Rank: 7}, // e8
				Promotion: chessnote.Queen,
			},
		},
		{
			name: "promotion with capture",
			pgn:  "1. exd8=R *",
			want: chessnote.Move{
				Piece:     chessnote.Pawn,
				From:      chessnote.Square{File: 4},          // 'e' file
				To:        chessnote.Square{File: 3, Rank: 7}, // d8
				IsCapture: true,
				Promotion: chessnote.Rook,
			},
		},
		{
			name: "promotion with check",
			pgn:  "1. e8=Q+ *",
			want: chessnote.Move{
				Piece:     chessnote.Pawn,
				To:        chessnote.Square{File: 4, Rank: 7}, // e8
				Promotion: chessnote.Queen,
				IsCheck:   true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game, err := chessnote.ParseString(tc.pgn)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if len(game.Moves) != 1 {
				t.Fatalf("expected 1 move, got %d", len(game.Moves))
			}
			got := game.Moves[0]

			// Custom comparison for pawn capture promotion
			if tc.name == "promotion with capture" {
				if got.From.File != tc.want.From.File {
					t.Errorf("Parse() got From.File = %d, want %d", got.From.File, tc.want.From.File)
				}
				// create a copy and clear the From field for the DeepEqual check
				got.From = chessnote.Square{}
				tc.want.From = chessnote.Square{}
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Parse() got = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestParseCastling(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		pgn  string
		want chessnote.Move
	}{
		{
			name: "kingside castle",
			pgn:  "1. O-O *",
			want: chessnote.Move{
				Piece:            chessnote.King,
				IsKingsideCastle: true,
			},
		},
		{
			name: "queenside castle",
			pgn:  "1. O-O-O *",
			want: chessnote.Move{
				Piece:             chessnote.King,
				IsQueensideCastle: true,
			},
		},
		{
			name: "kingside castle with check",
			pgn:  "1. O-O+ *",
			want: chessnote.Move{
				Piece:            chessnote.King,
				IsKingsideCastle: true,
				IsCheck:          true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game, err := chessnote.ParseString(tc.pgn)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if len(game.Moves) != 1 {
				t.Fatalf("expected 1 move, got %d", len(game.Moves))
			}
			got := game.Moves[0]

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Parse() got = %+v, want %+v", got, tc.want)
			}
		})
	}
}
