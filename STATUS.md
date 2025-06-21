# Project Status & Next Steps

This document tracks the current development status of the ChessNote PGN parser. All engineers should refer to this file to understand the current state of the project and what the immediate next steps are.

**Please update this file before every commit to reflect your progress.**

---

### Current Milestone

We have successfully **completed Milestone 1** and are now beginning **Milestone 2: Advanced Syntax & Semantic Validation**.

### Completed Features (Milestone 1)

- **Tag Pair Parsing**: Correctly parses the Seven Tag Roster and other arbitrary tags.
- **Core Movetext Parsing**:
  - Standard Pawn Moves (`e4`)
  - Standard Piece Moves (`Nf3`)
  - Piece Captures (`Nxf3`)
  - Pawn Captures (`exd5`)
  - Checks (`+`) and Checkmates (`#`)
- **Core Data Models**: Established `Game`, `Move`, `Square`, and `PieceType` structs.
- **Project Structure**:
  - A comprehensive `README.md` with installation and usage examples.
  - A `docs/` directory with EBNF grammar.
  - Full GoDoc comments for the public API.
  - A robust, table-driven test suite with high coverage.
- **Architectural Refinements**:
  - The parser has been refactored to use a `Scanner` for efficient tokenization.
  - Common logic has been centralized into an `internal/util` package with its own tests.

---

### Next Steps (Milestone 2)

The primary goal of this milestone is to support the full PGN specification.

#### 1. Immediate Next Step: Move Disambiguation

The next feature to be implemented is handling moves that require disambiguation.

- **Task**: Update the parser to handle movetext like `Rdf8` (specifying the 'd' file for the Rook) or `N1c3` (specifying the '1' rank for the Knight).
- **Files to Update**: `chessnote.go`, `chessnote_test.go`.
- **EBNF to Add**: Update `move` in `grammar.ebnf` to include optional rank/file specifiers.

#### 2. Upcoming Features

Once disambiguation is complete, we will proceed with the following features in order:

- **Pawn Promotion**: Parsing moves like `e8=Q`.
- **Castling**: Parsing both kingside (`O-O`) and queenside (`O-O-O`) castling.
- **Comments**: Handling both single-line (`;`) and multi-line (`{...}`) comments.
- **Recursive Annotation Variations (RAVs)**: Parsing nested move lines, e.g., `(1. e5 d5)`.
- **Numeric Annotation Glyphs (NAGs)**: Parsing annotations like `$1`, `$2`, etc. 
