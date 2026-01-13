package main_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary once
	tmpDir, err := os.MkdirTemp("", "affectedpkgs-build")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	binaryPath = filepath.Join(tmpDir, "affectedpkgs")
	if err := exec.Command("go", "build", "-o", binaryPath, ".").Run(); err != nil {
		panic("failed to build affectedpkgs: " + err.Error())
	}

	os.Exit(m.Run())
}

func runCLI(t *testing.T, dir string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run command: %v", err)
		}
	}
	return stdout.String(), exitCode
}

func TestAffectedPkgs_Basic(t *testing.T) {
	wd, _ := os.Getwd()
	testDataDir := filepath.Join(wd, "testdata")

	output, exitCode := runCLI(t, testDataDir, "github.com/sirupsen/logrus")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	expected := []string{
		"example.com/test",
		"example.com/test/lib",
	}
	
	// Split and trim output
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	// Check if all expected packages are present
	for _, exp := range expected {
		found := false
		for _, line := range lines {
			if strings.TrimSpace(line) == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected package %q not found in output:\n%s", exp, output)
		}
	}
}

func TestAffectedPkgs_Roots(t *testing.T) {
	wd, _ := os.Getwd()
	testDataDir := filepath.Join(wd, "testdata")

	output, exitCode := runCLI(t, testDataDir, "--roots", "github.com/sirupsen/logrus")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	// Only example.com/test should be present because example.com/test/lib is imported by example.com/test
	if !strings.Contains(output, "example.com/test") {
		t.Errorf("expected 'example.com/test' in output, got:\n%s", output)
	}
	if strings.Contains(output, "example.com/test/lib") {
		t.Errorf("did not expect 'example.com/test/lib' in output (it is not a root), got:\n%s", output)
	}
}

func TestAffectedPkgs_Test(t *testing.T) {
	wd, _ := os.Getwd()
	testDataDir := filepath.Join(wd, "testdata")

	output, exitCode := runCLI(t, testDataDir, "--test", "github.com/sirupsen/logrus")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	// Check for tools package which depends on logrus only in tests
	if !strings.Contains(output, "example.com/test/tools") {
		t.Errorf("expected 'example.com/test/tools' in output with --test, got:\n%s", output)
	}
}

func TestAffectedPkgs_JSON(t *testing.T) {
	wd, _ := os.Getwd()
	testDataDir := filepath.Join(wd, "testdata")

	output, exitCode := runCLI(t, testDataDir, "--json", "github.com/sirupsen/logrus")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	var pkgs []string
	if err := json.Unmarshal([]byte(output), &pkgs); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	foundLib := false
	for _, p := range pkgs {
		if p == "example.com/test/lib" {
			foundLib = true
			break
		}
	}
	if !foundLib {
		t.Error("expected example.com/test/lib in JSON output")
	}
}

func TestAffectedPkgs_NoMatch(t *testing.T) {
	wd, _ := os.Getwd()
	testDataDir := filepath.Join(wd, "testdata")

	output, exitCode := runCLI(t, testDataDir, "example.com/nonexistent")
	if exitCode != 1 {
		t.Errorf("expected exit code 1 for no match, got %d", exitCode)
	}
	if len(strings.TrimSpace(output)) > 0 {
		t.Errorf("expected empty output for no match, got: %s", output)
	}
}
