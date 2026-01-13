# affectedpkgs

`affectedpkgs` is a Go CLI tool designed to identify which packages in your Go module depend on a specific target module. It is particularly useful for analyzing the impact of updating a dependency or understanding the reach of a specific library within your codebase.

## Features

- **Dependency Analysis**: precise detection of direct and transitive dependencies.
- **Root Cause Analysis**: `--roots` flag to identify the top-level packages that introduce the dependency.
- **Test Dependencies**: `--test` flag to include test-only dependencies in the analysis.
- **CI/CD Integration**: `--json` output for easy parsing and integration with automation pipelines.
- **Efficient**: Uses `go list` with JSON parsing to handle large repositories efficiently.

## Development

This project uses [Task](https://taskfile.dev/) for development automation.

### Build and Test

```bash
task build
task test
```

### Documentation

To preview the documentation locally:

```bash
task docs:serve
```

To build the static documentation site:

```bash
task docs:build
```

## Installation

### From Source

```bash
go install github.com/yourusername/affectedpkgs@latest
```

*(Note: Replace `github.com/yourusername/affectedpkgs` with the actual repository path if hosted)*

### Build Manually

```bash
git clone https://github.com/yourusername/affectedpkgs.git
cd affectedpkgs
go build -o affectedpkgs main.go
```

## Usage

The basic syntax is:

```bash
affectedpkgs [flags] <module-path>
```

### Examples

**Find all packages that depend on `logrus`:**

```bash
affectedpkgs github.com/sirupsen/logrus
```

**Find only root packages (top-level consumers):**

```bash
affectedpkgs --roots github.com/sirupsen/logrus
```

**Include test dependencies:**

```bash
affectedpkgs --test github.com/sirupsen/logrus
```

**Output in JSON format:**

```bash
affectedpkgs --json github.com/sirupsen/logrus
```

## Flags

- `--roots`: Only print root packages (packages not imported by any other affected package).
- `--json`: Output the list of affected packages in JSON format.
- `--test`: Include test imports in the dependency analysis.

## Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/amazing-feature`).
3. Commit your changes (`git commit -m 'Add some amazing feature'`).
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a Pull Request.

## License

[MIT License](LICENSE)
