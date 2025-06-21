// Package chessnote provides a high-performance, production-grade Go library
// for parsing Portable Game Notation (PGN), the universal standard for chess game data.
package chessnote

import "strings"

// SplitMultiGame takes a string containing multiple PGN games and splits them
// into a slice of individual game strings. It normalizes line endings to
// handle different file formats (e.g., Windows-style \r\n).
//
// This utility is useful for pre-processing PGN files that contain an entire
// database of games before passing each individual game to the parser.
func SplitMultiGame(pgn string) []string {
	// Normalize line endings to \n to handle \r\n from Windows files.
	pgn = strings.ReplaceAll(pgn, "\r\n", "\n")
	var games []string
	var currentGame strings.Builder
	lines := strings.Split(pgn, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "[Event ") && currentGame.Len() > 0 {
			// Found the start of a new game, so save the previous one.
			gameStr := strings.TrimSpace(currentGame.String())
			if gameStr != "" {
				games = append(games, gameStr)
			}
			currentGame.Reset()
		}
		currentGame.WriteString(line)
		currentGame.WriteString("\n")
	}

	// Add the last game in the file.
	gameStr := strings.TrimSpace(currentGame.String())
	if gameStr != "" {
		games = append(games, gameStr)
	}

	return games
}
