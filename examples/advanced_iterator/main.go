// This example demonstrates a more advanced use case: iterating through a game's
// entire move tree, including all variations (Recursive Annotation Variations).
//
// Usage:
//
//	go run .
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/YashBhalodi/chessnote"
)

func main() {
	// For this example, we will always parse the same file.
	const pgnFilePath = "fischer_petrosian_1959.pgn"

	log.Printf("Parsing PGN file: %s", pgnFilePath)

	file, err := os.Open(pgnFilePath)
	if err != nil {
		log.Fatalf("Error opening PGN file: %v", err)
	}
	defer file.Close()

	parser := chessnote.NewParser(file)
	game, err := parser.Parse()
	if err != nil {
		log.Fatalf("Error parsing PGN: %v", err)
	}

	fmt.Printf("\n--- Game Tree for %s vs. %s ---\n", game.Tags["White"], game.Tags["Black"])

	// Begin the recursive traversal of the move tree, starting with the main line.
	printMoveTree(game.Moves, 0)

	fmt.Println("----------------------------------")
}

// printMoveTree recursively prints the moves of a game, indenting variations.
func printMoveTree(moves []chessnote.Move, depth int) {
	// The `depth` parameter tracks how deep we are in the variation tree.
	// We use it to calculate the indentation for printing.
	indent := strings.Repeat("    ", depth)

	for i, move := range moves {
		// We need a simple way to represent the move in SAN (Standard Algebraic Notation).
		// Since the library currently only provides the structured `Move` object,
		// we'll create a simple (and incomplete) formatter for this example.
		// A full SAN generator would be a great addition to the library itself!
		san := formatMoveToSAN(move)

		// Print the move number only for the main line (depth 0).
		if depth == 0 {
			// For Black's move, we print "..."
			if i%2 != 0 {
				fmt.Printf("%s%d... %s\n", indent, (i/2)+1, san)
			} else {
				fmt.Printf("%s%d. %s", indent, (i/2)+1, san)
			}
		} else {
			// In variations, we don't print move numbers to avoid confusion.
			fmt.Printf("%s- %s\n", indent, san)
		}

		// If the current move has variations, recursively call this function for each one.
		if len(move.Variations) > 0 {
			for _, variation := range move.Variations {
				// Start a new line for the variation and increase the depth.
				fmt.Println()
				printMoveTree(variation, depth+1)
			}
		}
	}
}

// formatMoveToSAN creates a simplified string representation of a move.
// Note: This is not a complete or fully compliant SAN generator.
func formatMoveToSAN(move chessnote.Move) string {
	if move.IsKingsideCastle {
		return "O-O"
	}
	if move.IsQueensideCastle {
		return "O-O-O"
	}

	var sb strings.Builder

	piece := pieceToChar(move.Piece)
	if move.Piece != chessnote.Pawn {
		sb.WriteString(piece)
	}

	// This is a simplification; real SAN disambiguation is more complex.
	if move.From.File > 0 {
		sb.WriteRune(rune('a' + move.From.File))
	}
	if move.From.Rank > 0 {
		sb.WriteRune(rune('1' + move.From.Rank))
	}

	if move.IsCapture {
		if move.Piece == chessnote.Pawn && sb.Len() == 0 {
			sb.WriteRune(rune('a' + move.From.File))
		}
		sb.WriteString("x")
	}

	sb.WriteString(squareToString(move.To))

	if move.Promotion != chessnote.Pawn {
		sb.WriteString("=")
		sb.WriteString(pieceToChar(move.Promotion))
	}

	if move.IsMate {
		sb.WriteString("#")
	} else if move.IsCheck {
		sb.WriteString("+")
	}

	return sb.String()
}

func squareToString(s chessnote.Square) string {
	return string(rune('a'+s.File)) + string(rune('1'+s.Rank))
}

func pieceToChar(p chessnote.PieceType) string {
	switch p {
	case chessnote.Knight:
		return "N"
	case chessnote.Bishop:
		return "B"
	case chessnote.Rook:
		return "R"
	case chessnote.Queen:
		return "Q"
	case chessnote.King:
		return "K"
	default:
		return ""
	}
}
