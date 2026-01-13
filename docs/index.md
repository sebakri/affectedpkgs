# affectedpkgs

**affectedpkgs** is a lightweight, efficient Command Line Interface (CLI) tool for Go developers. It helps you analyze your Go module's dependency graph to determine exactly which packages are affected by a specific external module.

## Why use affectedpkgs?

In large Go monorepos or multi-module projects, updating a core library (like a logger, a database driver, or a utility library) can have widespread effects. `affectedpkgs` helps you answer questions like:

*   "If I upgrade `github.com/sirupsen/logrus`, which parts of my application do I need to re-test?"
*   "Which of my services actually depend on this legacy library?"
*   "What are the top-level packages causing this dependency to be pulled in?"

## Key Features

*   **Fast Analysis**: Leverages `go list -deps -json` for native, efficient graph traversal.
*   **Transitive Awareness**: Detects dependencies deep in the graph, not just direct imports.
*   **Test Scope**: Optionally includes test files in the analysis to see if a dependency is only used for testing.
*   **Root Detection**: Can filter output to show only the "root" packages causing the dependency, reducing noise.
*   **Automation Ready**: JSON output mode makes it perfect for piping into other tools or CI scripts.

## Getting Started

Check out the [Installation](installation.md) guide to get set up, or jump straight to [Usage](usage.md) to see examples.
