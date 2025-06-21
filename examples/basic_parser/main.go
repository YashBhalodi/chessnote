// This example demonstrates how to parse a PGN file from your local disk.
//
// Usage:
//
//  1. Run with the default PGN file:
//     go run .
//
//  2. Run with your own PGN file:
//     go run . --pgn_file=/path/to/your/game.pgn
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/YashBhalodi/chessnote"
)

func main() {
	// Use Go's flag package to create a command-line flag for the PGN file path.
	// This allows a user to easily point this example to their own files.
	pgnFilePath := flag.String("pgn_file", "opera_game.pgn", "Path to the PGN file to parse.")
	flag.Parse()

	log.Printf("Parsing PGN file: %s", *pgnFilePath)

	// Open the PGN file. The os.Open function returns an *os.File object,
	// which satisfies the io.Reader interface that our parser needs.
	file, err := os.Open(*pgnFilePath)
	if err != nil {
		log.Fatalf("Error opening PGN file: %v", err)
	}
	defer file.Close()

	// Create a new parser for the file.
	// The NewParser function is the most flexible way to create a parser,
	// as it works with any io.Reader (files, network connections, in-memory buffers, etc.).
	parser := chessnote.NewParser(file)

	// Call the Parse method to process the PGN data.
	game, err := parser.Parse()
	if err != nil {
		log.Fatalf("Error parsing PGN: %v", err)
	}

	// Now that we have the parsed game, we can inspect its data.
	fmt.Println("\n--- Game Information ---")
	fmt.Printf("Event: %s\n", game.Tags["Event"])
	fmt.Printf("Site: %s\n", game.Tags["Site"])
	fmt.Printf("Date: %s\n", game.Tags["Date"])
	fmt.Printf("White: %s\n", game.Tags["White"])
	fmt.Printf("Black: %s\n", game.Tags["Black"])
	fmt.Printf("Result: %s\n", game.Result)
	fmt.Printf("Total Moves: %d\n", len(game.Moves))
	fmt.Println("----------------------")
}
