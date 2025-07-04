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
- **Move Disambiguation**: The parser now handles movetext like `Rdf8` and `N1c3`.
- **Pawn Promotion**: The parser now handles movetext like `e8=Q` and `exd8=R+`.
- **Castling**: The parser now handles both kingside (`O-O`) and queenside (`O-O-O`) castling.
- **Code Health & Refactoring**: Completed a comprehensive cleanup phase which included:
  - Consolidation of all move-parsing tests into a single table-driven test.
  - Refactoring of the core move parser for robustness and clarity.
  - Hardening of the scanner tests.
  - A full review and update of all GoDoc comments and documentation.

---

### Next Steps (Milestone 2)

The primary goal of this milestone is to support the full PGN specification.

#### 1. Completed Features (Milestone 2)

- **Comments**: The parser now handles both single-line (`;`) and multi-line (`{...}`) comments by correctly ignoring them.
- **Recursive Annotation Variations (RAVs)**: The parser now correctly handles nested move lines, e.g., `(1. e5 d5)`. The parser was refactored to use a stateful, recursive descent approach.
- **Numeric Annotation Glyphs (NAGs)**: The parser now correctly handles NAGs, e.g., `$1`, `$18`, associating them with the preceding move.

### Milestone 2 Complete

With the implementation of NAGs, all major parsing features for the PGN standard are now complete.

---

### Next Steps (Milestone 3: Beta Release & Hardening)

The primary goal of this milestone is to prepare the library for a stable v1.0.0 beta release. The focus will shift from adding new features to improving robustness, documentation, and the public API.

#### 1. Immediate Next Step: Final Parser Review & Refinement

- **Task**: Conduct a thorough review of the entire parser and its data structures. Refine the public API, improve GoDoc comments, and ensure all code adheres to our engineering guidelines. This includes a final review of the `parseMove` and `parseCoreMove` functions for clarity and robustness.
- **Files to Update**: All `*.go` files.

#### 2. Upcoming Tasks

- **Enhanced Error Reporting**: Improve error messages to be more specific and user-friendly, including line and column numbers where possible.
- **Fuzz Testing**: Implement a comprehensive fuzz testing suite for the parser to ensure it can handle any malformed input without crashing.
- **Benchmarking**: Add a suite of benchmarks to measure and track the performance of the parser.
- **Examples**: Create a set of examples in the `examples/` directory to showcase how to use the library.
- **README Overhaul**: Update the `README.md` to reflect the complete feature set and provide a comprehensive guide for new users.

---

## Milestone 4: Polish and Release

- [x] Setup CI with GitHub Actions
- [x] Add project badges to README
- [x] Add LICENSE file
- [x] Refactor `parseCoreMove` to reduce cyclomatic complexity

## Next Step

- Publish the first version of the library (`v0.1.0`).

## Milestone 5: The First Release

- [x] Publish `v0.1.0`

---

### Milestone 6: Library Hardening

The primary goal of this milestone is to elevate the library to a true production-grade standard by adding benchmarks, fuzz testing, and contributor documentation. This will prepare us for a stable `v1.0.0` release.

- **Completed Features**:
    - **Benchmarking**: Added performance benchmarks, including a comprehensive test against a large, real-world PGN file (`Kasparov.pgn`), to validate performance claims and prevent regressions.
    - **Documentation**: Reorganized and enhanced project documentation.
    - **Fuzz Testing**: Implemented a fuzz test for the parser, ensuring it can handle a wide range of malformed inputs without crashing.
- **Next Step**: Create a `CONTRIBUTING.md` file to guide future contributors.

### Milestone 6 Complete

With the addition of comprehensive testing and documentation, the library is now significantly hardened and ready for a wider audience.

---

### Housekeeping & DX

- **Project Structure**: Reorganized the project to have dedicated `tests/` and `benchmarks/` directories, cleaning up the root.
- **Examples**: Added a new, runnable example in `examples/multiple-game-pgn/` to demonstrate how to use the `SplitMultiGame` utility.

### Future Milestones

- **`chessnote`** CLI Tool
- Semantic Validator

---

### Recently Completed

- **Parser Configuration**: Implemented a flexible configuration system for the parser using functional options. The parser now defaults to a strict mode that requires a game termination token, but can be switched to a more lenient "lax" mode. This aligns the parser's default behavior with the formal PGN specification while still supporting malformed PGNs.
- **`CONTRIBUTING.md`**: Created a comprehensive guide for new contributors.
- **Fuzz Testing**: Hardened the parser with a fuzz testing suite.
