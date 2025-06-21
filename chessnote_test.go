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
