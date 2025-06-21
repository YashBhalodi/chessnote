package benchmarks

import (
	"os"
	"testing"

	"github.com/YashBhalodi/chessnote"
)

func BenchmarkParseOperaGame(b *testing.B) {
	pgn, err := os.ReadFile("../examples/basic_parser/opera_game.pgn")
	if err != nil {
		b.Fatalf("failed to read PGN file: %v", err)
	}
	pgnStr := string(pgn)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := chessnote.ParseString(pgnStr)
		if err != nil {
			b.Fatalf("ParseString() failed: %v", err)
		}
	}
}

func BenchmarkParseFischerPetrosian(b *testing.B) {
	pgn, err := os.ReadFile("../examples/advanced_iterator/fischer_petrosian_1959.pgn")
	if err != nil {
		b.Fatalf("failed to read PGN file: %v", err)
	}
	pgnStr := string(pgn)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := chessnote.ParseString(pgnStr)
		if err != nil {
			b.Fatalf("ParseString() failed: %v", err)
		}
	}
}

func BenchmarkParseKasparovGames(b *testing.B) {
	pgn, err := os.ReadFile("../examples/Kasparov.pgn")
	if err != nil {
		b.Fatalf("failed to read PGN file: %v", err)
	}
	games := chessnote.SplitMultiGame(string(pgn))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, game := range games {
			_, err := chessnote.ParseString(game)
			if err != nil {
				b.Fatalf("ParseString() failed: %v\nPGN:\n%s", err, game)
			}
		}
	}
}
