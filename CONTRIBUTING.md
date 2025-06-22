# Contributing to ChessNote

First off, thank you for considering contributing to ChessNote! It's people like you that make open source such a great community. We welcome any and all contributions.

This document provides guidelines for contributing to the project. Please read it carefully to ensure a smooth and effective collaboration process.

## How Can I Contribute?

There are many ways to contribute, from writing code and documentation to reporting bugs and suggesting features.

- **Reporting Bugs:** If you find a bug, please open an issue on our GitHub repository. Describe the bug in detail, including steps to reproduce it, the expected outcome, and the actual outcome.
- **Suggesting Enhancements:** If you have an idea for a new feature or an improvement to an existing one, please open an issue to discuss it.
- **Pull Requests:** If you're ready to contribute code or documentation, please submit a pull request.

## Development Workflow

We follow a strict Test-Driven Development (TDD) and Pull Request-based workflow. All contributions must adhere to these guidelines.

1.  **Fork the Repository:** Start by forking the official `chessnote` repository on GitHub.
2.  **Create a Branch:** Create a new branch for your feature or bug fix. Name it descriptively (e.g., `feat/add-new-parser-feature`, `fix/issue-with-castling`).
    ```bash
    git checkout -b your-branch-name
    ```
3.  **Follow TDD:**
    - **Write a Failing Test:** Before you write any implementation code, write a test that describes the feature or reproduces the bug. This test should fail.
    - **Write Code to Pass the Test:** Write the simplest possible code to make the test pass.
    - **Refactor:** Clean up your code, ensuring all tests still pass.
4.  **Ensure All Checks Pass:** Before submitting, make sure your code is formatted with `go fmt` and passes all linter checks and existing tests.
    ```bash
    go fmt ./...
    go test ./...
    ```
5.  **Update Documentation:** All documentation related to a code change **must** be included in the same pull request. This includes GoDoc comments, the `README.md`, and any files in the `docs/` directory. Accurate documentation is a requirement for merging.
6.  **Update `STATUS.md`:** Before you commit, you **must** update the `STATUS.md` file to reflect the work you have completed and what the new "Next Step" is. This is crucial for keeping the team aligned.
7.  **Commit Your Changes:** Commit your changes with a descriptive, well-formatted commit message. We follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.
    ```bash
    git commit -m "feat: Briefly describe the new feature"
    ```
8.  **Push and Open a Pull Request:** Push your branch to your fork and open a pull request against the `main` branch of the official repository.
9.  **Code Review:** Your pull request will be reviewed by the maintainers. Be prepared to discuss your changes and make any necessary adjustments.

## Coding Standards

Please adhere to the standards outlined in our `README.md` and the general best practices of the Go community.

- Write idiomatic Go.
- Keep the public API clean and well-documented.
- Minimize dependencies.

Thank you for your contribution! 
