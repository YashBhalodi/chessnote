# ChessNote

A high-performance, production-grade Go library for parsing Portable Game Notation (PGN).

ChessNote provides a robust, reliable, and easy-to-use tool to read, parse, and validate chess games. It is built from the ground up with a focus on excellent error handling, strong performance, and a clean, idiomatic Go API.

## Installation

To add ChessNote to your project, use `go get`:

```sh
go get github.com/HexaTech/chessnote
```

## Quick Start

Here's a simple example of how to parse the tags from a PGN string:

```go
package main

import (
	"fmt"

	"github.com/HexaTech/chessnote"
)

func main() {
	pgn := `
[Event "F/S Return Match"]
[Site "Belgrade, Serbia JUG"]
[Date "1992.11.04"]
[Round "29"]
[White "Fischer, Robert J."]
[Black "Spassky, Boris V."]
[Result "1/2-1/2"]

1. e4 e5 2. Nf3 Nc6 3. Bb5+ *
`

	game, err := chessnote.ParseString(pgn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Tags:")
	for key, value := range game.Tags {
		fmt.Printf("  %s: %s\n", key, value)
	}

	fmt.Println("\nMoves:")
	for i, move := range game.Moves {
		// This is a simplified representation of the move.
		// A real application would have a more sophisticated move formatter.
		fmt.Printf("  %d: %+v\n", i+1, move)
	}
}
```
