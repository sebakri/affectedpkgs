# Installation

`affectedpkgs` is written in Go, so you can install it easily using standard Go tools or by building from source.

## Prerequisites

*   **Go**: You need Go installed on your machine (version 1.18 or later is recommended). Verify with `go version`.

## Method 1: Using Shell Script (Fastest)

You can install the latest binary with a single command:

```bash
curl -sSL https://raw.githubusercontent.com/sebakri/affectedpkgs/main/install.sh | bash
```

This script automatically detects your OS and architecture, downloads the latest release, and installs it to `/usr/local/bin`.

## Method 2: Using `go install`

To install the latest version directly:

```bash
go install github.com/yourusername/affectedpkgs@latest
```

Ensure that your `$GOPATH/bin` is in your system's `PATH`.

## Method 2: Building from Source

If you want to modify the code or contribute:

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/yourusername/affectedpkgs.git
    cd affectedpkgs
    ```

2.  **Build the binary:**
    ```bash
    go build -o affectedpkgs main.go
    ```

3.  **Move to PATH (Optional):**
    ```bash
    mv affectedpkgs /usr/local/bin/
    ```

## Verification

Run the help command to verify the installation:

```bash
affectedpkgs -h
```
