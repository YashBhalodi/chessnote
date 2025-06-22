package main

import (
	"fmt"

	"github.com/YashBhalodi/chessnote"
)

func main() {
	// A PGN string that is missing the required game termination token.
	// According to the PGN standard, every game must end with a result
	// like "1-0", "0-1", "1/2-1/2", or "*".
	pgnWithoutResult := `
[Event "A Game Without a Result"]
[White "Player A"]
[Black "Player B"]

1. e4 e5 2. Nf3
`

	fmt.Println("--- Attempting to parse with default (strict) settings ---")
	// By default, the parser is in strict mode and will return an error
	// if the game termination token is missing.
	_, err := chessnote.ParseString(pgnWithoutResult)
	if err != nil {
		fmt.Println("Successfully caught expected error in strict mode:")
		fmt.Println(err)
	}

	fmt.Println("\n--- Attempting to parse with WithLaxParsing() option ---")
	// We can use the WithLaxParsing() option to create a more lenient parser
	// that does not require the final result token.
	game, err := chessnote.ParseString(pgnWithoutResult, chessnote.WithLaxParsing())
	if err != nil {
		fmt.Println("Unexpected error in lax mode:", err)
		return
	}

	fmt.Println("Successfully parsed in lax mode!")
	fmt.Printf("Found %d moves.\n", len(game.Moves))
	fmt.Printf("White: %s\n", game.Tags["White"])
	fmt.Printf("Black: %s\n", game.Tags["Black"])
}
