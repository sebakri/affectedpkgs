# Usage

`affectedpkgs` is run from the command line inside your Go project (where `go.mod` is located).

## Basic Command

```bash
affectedpkgs [flags] <target-module-path>
```

*   **`<target-module-path>`**: The full import path of the module you are investigating (e.g., `github.com/gin-gonic/gin`).

## Flags

### `--roots`
**Description**: Filters the output to show only "root" packages. A root package is one that depends on the target but is not imported by any other package in the affected set.

**Use Case**: Useful when you want to find the entry points (like `main` packages or top-level libraries) that are responsible for bringing in the dependency.

```bash
affectedpkgs --roots github.com/lib/pq
```

### `--test`
**Description**: Includes test files (`*_test.go`) in the dependency analysis.

**Use Case**: Sometimes a library is only used in tests (e.g., a mocking library or assertion framework). By default, `affectedpkgs` looks at production code. Use this flag to broaden the search.

```bash
affectedpkgs --test github.com/stretchr/testify
```

### `--json`
**Description**: Outputs the result as a JSON array of strings.

**Use Case**: Ideal for scripting and CI/CD pipelines. For example, you could feed the output into a linter or a testing tool to only run tests for affected packages.

```bash
affectedpkgs --json github.com/sirupsen/logrus | jq length
```

## Examples

### Scenario 1: Upgrade Impact Analysis

You plan to upgrade `github.com/aws/aws-sdk-go` and want to know which packages in your monolith need attention.

```bash
affectedpkgs github.com/aws/aws-sdk-go
```
*Output:*
```text
mycompany.com/monolith/pkg/storage
mycompany.com/monolith/pkg/notifications
mycompany.com/monolith/cmd/worker
```

### Scenario 2: Clean up Dependencies

You think `github.com/pkg/errors` is no longer needed, or you want to see who is still using it.

```bash
affectedpkgs --roots github.com/pkg/errors
```

### Scenario 3: CI Integration

Run tests only for affected packages:

```bash
affectedpkgs --json github.com/changed/module > affected.json
go test $(cat affected.json | jq -r '.[]')
```
