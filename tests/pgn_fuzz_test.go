package chessnote_test

import (
	"testing"

	"github.com/YashBhalodi/chessnote"
)

func FuzzParse(f *testing.F) {
	// Add some valid PGN snippets as seed corpus.
	// This helps the fuzzer to start from a good baseline.
	f.Add("[Event \"F/S Return Match\"]")
	f.Add("1. e4 e5 2. Nf3 Nc6 *")
	f.Add("[White \"Kasparov, Garry\"] 1/2-1/2")

	f.Fuzz(func(t *testing.T, data string) {
		// The parser should handle any string input without panicking.
		// We don't need to check the error, as the goal of fuzzing here
		// is to find panic-inducing inputs.
		_, _ = chessnote.ParseString(data)
	})
}
