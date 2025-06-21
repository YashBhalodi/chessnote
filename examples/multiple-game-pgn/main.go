package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/YashBhalodi/chessnote"
)

func main() {
	// Get the directory of the currently running file to build a reliable path.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Error: unable to get the current file path")
		os.Exit(1)
	}
	dir := filepath.Dir(filename)
	pgnPath := filepath.Join(dir, "Kasparov.pgn")

	// Read the multi-game PGN file.
	pgnData, err := os.ReadFile(pgnPath)
	if err != nil {
		fmt.Printf("Error reading PGN file: %v\n", err)
		os.Exit(1)
	}

	// Use the SplitMultiGame utility to get a slice of individual game strings.
	games := chessnote.SplitMultiGame(string(pgnData))

	fmt.Printf("Found %d games in the PGN file.\n\n", len(games))

	// Loop through and parse the first 5 games as a demonstration.
	fmt.Println("Parsing the first 5 games:")
	for i, gameStr := range games {
		if i >= 5 {
			break
		}

		game, err := chessnote.ParseString(gameStr)
		if err != nil {
			fmt.Printf("  Error parsing game %d: %v\n", i+1, err)
			continue
		}

		fmt.Printf("  Game %d: %s vs. %s, Result: %s\n",
			i+1,
			game.Tags["White"],
			game.Tags["Black"],
			game.Tags["Result"],
		)
	}
}
