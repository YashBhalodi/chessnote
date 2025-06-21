package chessnote_test

import (
	"reflect"
	"strings"
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
	p := chessnote.NewParser()
	game, err := p.Parse(strings.NewReader(pgn))
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
	p := chessnote.NewParser()
	game, err := p.Parse(strings.NewReader(pgn))
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
