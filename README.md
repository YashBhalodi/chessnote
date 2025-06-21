# ChessNote: High-Performance PGN Parsing in Go

[![Go Report Card](https://goreportcard.com/badge/github.com/YashBhalodi/chessnote)](https://goreportcard.com/report/github.com/YashBhalodi/chessnote)
[![Go Reference](https://pkg.go.dev/badge/github.com/YashBhalodi/chessnote.svg)](https://pkg.go.dev/github.com/YashBhalodi/chessnote)
[![CI](https://github.com/YashBhalodi/chessnote/actions/workflows/ci.yml/badge.svg)](https://github.com/YashBhalodi/chessnote/actions/workflows/ci.yml)

**ChessNote is a production-grade, high-performance Go library for parsing Portable Game Notation (PGN).** It's engineered from the ground up to be a foundational component for ambitious chess applications, from powerful analysis tools and database backends to beautiful game renderers.

We're not just parsing text; we're providing a reliable, rigorously-tested, and developer-friendly toolkit for bringing chess data to life.

## Why ChessNote?

In a world of hobbyist libraries, ChessNote is engineered for professional use. It's built on a philosophy of reliability and performance, making it the ideal choice for startups and developers who need a parser that just works.

*   **Rock-Solid Reliability:** Built with a strict **Test-Driven Development (TDD)** methodology, ChessNote boasts comprehensive test coverage. Every feature, from simple pawn moves to complex variations, is verified. The parser is designed to handle real-world, messy PGNs without panicking, returning structured errors that make debugging a breeze.

*   **Blazing-Fast Performance:** Written in pure, idiomatic Go, ChessNote is designed for speed. It minimizes allocations and uses an efficient scanning and parsing model to handle large PGN databases with ease. Its performance is validated by a comprehensive benchmark suite (see the `benchmarks/` directory).

*   **A Developer-First API:** The public API is clean, discoverable, and a joy to use. We believe developers should spend their time building great applications, not fighting with a clunky parser. Our data structures (`Game`, `Move`, `Square`) are intuitive and well-documented.

## Feature-Complete Parser

ChessNote implements the complete PGN standard, allowing you to parse virtually any PGN file you encounter.

- [x] **Full Tag Pair Support:** Parses the Seven Tag Roster and any other custom tags.
- [x] **Complete Movetext Parsing:**
    - [x] Standard pawn and piece moves (`e4`, `Nf3`)
    - [x] Captures (`exd5`, `Nxf3`)
    - [x] Checks (`+`) and Checkmates (`#`)
    - [x] Disambiguation (`Rdf8`, `N1c3`)
    - [x] Pawn Promotion (`e8=Q`)
    - [x] Castling (`O-O`, `O-O-O`)
- [x] **Advanced PGN Syntax:**
    - [x] Comments (`{...}` and `;...`)
    - [x] Recursive Annotation Variations (RAVs) `(...)`
    - [x] Numeric Annotation Glyphs (NAGs) `($1, $18)`
- [x] **Game Termination Markers:** Correctly identifies the game result (`1-0`, `0-1`, `1/2-1/2`, `*`).
- [x] **Robust Error Handling:** Returns detailed, structured errors for invalid syntax.

## Quick Start

Get the library:
```sh
go get github.com/YashBhalodi/chessnote
```

Parse a complete PGN game in just a few lines of code. Here's an example parsing Paul Morphy's famous "Opera Game":

```go
package main

import (
	"fmt"
	"log"

	"github.com/YashBhalodi/chessnote"
)

func main() {
	// PGN of the "Opera Game" (Paul Morphy vs. Duke Karl / Count Isouard, 1858)
	pgn := `
[Event "A Night at the Opera"]
[Site "Paris, France"]
[Date "1858.??.??"]
[White "Paul Morphy"]
[Black "Duke Karl / Count Isouard"]
[Result "1-0"]

1. e4 e5 2. Nf3 d6 3. d4 Bg4 4. dxe5 Bxf3 5. Qxf3 dxe5 6. Bc4 Nf6 7. Qb3 Qe7
8. Nc3 c6 9. Bg5 b5 10. Nxb5 cxb5 11. Bxb5+ Nbd7 12. O-O-O Rd8
13. Rxd7 Rxd7 14. Rd1 Qe6 15. Bxd7+ Nxd7 16. Qb8+ Nxb8 17. Rd8# 1-0
`
	game, err := chessnote.ParseString(pgn)
	if err != nil {
		log.Fatalf("Failed to parse PGN: %v", err)
	}

	fmt.Printf("Successfully parsed game: %s vs %s\n", game.Tags["White"], game.Tags["Black"])
	fmt.Printf("Result: %s\n", game.Result)
	fmt.Printf("Total moves: %d\n", len(game.Moves))

	// Let's inspect the brilliant final move!
	finalMove := game.Moves[len(game.Moves)-1]
	fmt.Printf("Final move: Rd8#\n")
	fmt.Printf("  - Piece: Rook\n")
	fmt.Printf("  - Destination: d8\n")
	fmt.Printf("  - Is Checkmate: %t\n", finalMove.IsMate)
}
```

## Running the Examples

The `examples/` directory contains standalone, runnable programs to showcase the library's features.

To run the basic parsing example:

```sh
# Navigate to the example directory
cd examples/basic_parser

# Run the program. It will parse the included 'opera_game.pgn' by default.
go run .

# You can also point it to your own PGN file
go run . --pgn_file=/path/to/your/game.pgn
```

To run the advanced iterator example, which demonstrates traversing a game's move tree including variations:

```sh
# Navigate to the example directory
cd examples/advanced_iterator

# Run the program.
go run .
```

## Project Philosophy

This project adheres to a strict set of engineering guidelines focused on creating professional, enterprise-grade software. We practice **Test-Driven Development (TDD)**, maintain a comprehensive test suite, and prioritize a clean, stable, and well-documented public API.

## Technical Documentation

For a deep dive into the parser's architecture, execution flow, and design decisions, please see our comprehensive **[Technical Documentation](./docs/README.md)**. This is the best resource for developers looking to contribute to the project or understand its inner workings.

## Roadmap & Contributing

ChessNote is currently feature-complete for parsing. Our future roadmap (Milestone 4) is focused on building tools on top of the library:

*   **`chessnote` CLI:** A command-line tool for validating and analyzing PGN files.
*   **GIF Renderer:** A tool to generate animated GIFs of chess games.

Contributions are welcome! Please open an issue to discuss any proposed changes.

## License

ChessNote is released under the [MIT License](./LICENSE).
