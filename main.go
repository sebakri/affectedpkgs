package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Package struct {
	ImportPath   string
	Module       *Module
	Imports      []string
	TestImports  []string
	XTestImports []string
	Standard     bool
}

type Module struct {
	Path string
}

func main() {
	rootsFlag := flag.Bool("roots", false, "only print root packages")
	jsonFlag := flag.Bool("json", false, "output in JSON format")
	testFlag := flag.Bool("test", false, "include test imports")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: affectedpkgs [flags] <module-path>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	targetModule := flag.Arg(0)

	// Step 1: Identify the packages in the current module (the "targets").
	// We run go list ./... to get the packages we are analyzing.
	targetImportPaths, err := getPackageList(false, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing packages: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Build the full dependency graph.
	// We run go list -deps -json ./... to get everything.
	graph, err := getPackageGraph(*testFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting dependency graph: %v\n", err)
		os.Exit(1)
	}

	// Step 3: Find affected packages.
	affected := make([]string, 0)
	memo := make(map[string]bool)

	// Helper to check dependency recursively
	var check func(string) bool
	check = func(importPath string) bool {
		if res, ok := memo[importPath]; ok {
			return res
		}

		pkg, ok := graph[importPath]
		if !ok {
			// If package is missing from graph (shouldn't happen with -deps), assume false
			// or it might be a standard lib package not fully detailed?
			// Standard packages usually don't depend on 3rd party modules.
			memo[importPath] = false
			return false
		}

		// Direct dependency check
		if pkg.Module != nil && pkg.Module.Path == targetModule {
			memo[importPath] = true
			return true
		}

		// Recurse
		// Combine Imports, TestImports, XTestImports based on flags/context
		// Note: go list -deps usually lists dependencies of imports.
		// For the *target* packages, we might care about their TestImports.
		// But for transitive dependencies, we usually only care about Imports.
		// However, if A imports B, and B's test imports Target, A is NOT affected by Target (unless A runs B's tests, which is not how it works).
		// So TestImports are only relevant for the *starting* packages?
		// Actually, standard `go list -deps -test` logic is complex.
		// Simplification:
		// deps = pkg.Imports
		// if *testFlag && (isTarget(pkg)) { deps += TestImports + XTestImports }
		// Wait, `go list -deps -test` produces a graph where test dependencies are included in the stream.
		// But the `Imports` list in the JSON might still be separated.
		// The `affectedpkgs` tool usually asks: "Does this package depend on module M?"
		// If I run `go test ./pkg`, will it download module M?
		// If `pkg` imports `M`, yes. If `pkg` test-imports `M`, yes.
		// So for the *roots*, we must check TestImports/XTestImports if --test is on.
		// For transitive deps, we only follow `Imports` because `pkg` importing `dep` doesn't pull in `dep`'s test deps.

		// Optimization: Construct the list of dependencies to check.
		depsToCheck := pkg.Imports
		
		// We only check TestImports if this package is one of the roots (targets) AND --test is specified.
		// Because `go list -deps` for a dependency doesn't normally include its test deps unless we are running tests for THAT dependency.
		// But here we are running tests for `./...`.
		// So only the packages in `./...` have their test dependencies active.
		isRoot := false
		for _, t := range targetImportPaths {
			if t == importPath {
				isRoot = true
				break
			}
		}

		if *testFlag && isRoot {
			depsToCheck = append(depsToCheck, pkg.TestImports...)
			depsToCheck = append(depsToCheck, pkg.XTestImports...)
		}

		for _, depPath := range depsToCheck {
			if check(depPath) {
				memo[importPath] = true
				return true
			}
		}

		memo[importPath] = false
		return false
	}

	for _, p := range targetImportPaths {
		if check(p) {
			affected = append(affected, p)
		}
	}

	if len(affected) == 0 {
		os.Exit(1)
	}

	// Step 4: Handle --roots flag
	// "Only print root packages".
	// Interpretation: Filter out packages that are imported by other packages in the affected set.
	// Note: We only filter based on *Import* relationships, not TestImports.
	// If A imports B, and both are affected. 
	// A is a "root" (nothing imports it). B is not.
	// If A test-imports B. A depends on B.
	if *rootsFlag {
		affectedSet := make(map[string]bool)
		for _, p := range affected {
			affectedSet[p] = true
		}

		isImportedByOther := make(map[string]bool)
		for _, p := range affected {
			pkg := graph[p]
			// Check who this package imports.
			// If p imports q, and q is in affectedSet, then q is imported by p.
			// So q is NOT a root.
		
			// We must include TestImports here if --test is on, because if A test-imports B, A is the "parent" causing B to be included.
			deps := pkg.Imports
			if *testFlag {
				deps = append(deps, pkg.TestImports...)
				deps = append(deps, pkg.XTestImports...)
			}

			for _, dep := range deps {
				if affectedSet[dep] {
					isImportedByOther[dep] = true
				}
			}
		}

		newAffected := make([]string, 0)
		for _, p := range affected {
			if !isImportedByOther[p] {
				newAffected = append(newAffected, p)
			}
		}
		affected = newAffected
	}

	// Step 5: Output
	if *jsonFlag {
		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(affected); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
			os.Exit(1)
		}
	} else {
		for _, p := range affected {
			fmt.Println(p)
		}
	}
}

// getPackageList returns the list of ImportPaths matching ./...
func getPackageList(test bool, deps bool) ([]string, error) {
	args := []string{"list", "-json"}
	if test {
		args = append(args, "-test")
	}
	if deps {
		args = append(args, "-deps")
	}
	args = append(args, "./...")

	cmd := exec.Command("go", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var paths []string
	dec := json.NewDecoder(stdout)
	for {
		var p Package
		if err := dec.Decode(&p); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		paths = append(paths, p.ImportPath)
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return paths, nil
}

// getPackageGraph returns a map of all dependencies
func getPackageGraph(test bool) (map[string]*Package, error) {
	args := []string{"list", "-json", "-deps"}
	if test {
		args = append(args, "-test")
	}
	args = append(args, "./...")

	cmd := exec.Command("go", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	// Increase buffer size for large outputs if needed, but Decoder handles stream.
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	graph := make(map[string]*Package)
	dec := json.NewDecoder(stdout)
	for {
		var p Package
		if err := dec.Decode(&p); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		// In case of duplicates (which shouldn't happen with standard go list -deps?), last one wins or ignore.
		// go list output is unique per ImportPath usually.
		// Note: When using -test, go list might output the package and its test binary variant.
		// The test binary variant usually has a distinct ImportPath or is just processed differently.
		// Standard `go list -json` for `p` includes `TestImports`.
		// `go list -deps -test` might output `p` and `p [p.test]`.
		// We care about `p`.
		// If `p` depends on `q`.
		// We store everything.
		graph[p.ImportPath] = &p
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return graph, nil
}
