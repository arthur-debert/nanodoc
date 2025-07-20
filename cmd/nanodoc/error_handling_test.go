package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestCLIErrorDisplay ensures that CLI errors are properly displayed to users
// This test prevents regression of issue #66 where errors were silently swallowed
func TestCLIErrorDisplay(t *testing.T) {
	// Skip if not in CI or if explicitly requested
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test")
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", "test-nanodoc", ".")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer func() { _ = os.Remove("test-nanodoc") }()

	tests := []struct {
		name           string
		args           []string
		wantError      string
		wantExitCode   int
	}{
		{
			name:         "invalid flag",
			args:         []string{"--invalid-option"},
			wantError:    "unknown flag: --invalid-option",
			wantExitCode: 1,
		},
		{
			name:         "invalid linenum value",
			args:         []string{"--linenum", "invalid", "README.md"},
			wantError:    "invalid --linenum value: invalid (must be 'file' or 'global')",
			wantExitCode: 1,
		},
		{
			name:         "invalid output format",
			args:         []string{"--output-format", "wrongformat", "README.md"},
			wantError:    "invalid --output-format value: wrongformat (must be 'term', 'plain', or 'markdown')",
			wantExitCode: 1,
		},
		{
			name:         "missing required arguments",
			args:         []string{},
			wantError:    "Missing paths to bundle",
			wantExitCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./test-nanodoc", tt.args...)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			err := cmd.Run()

			// Check exit code
			if exitError, ok := err.(*exec.ExitError); ok {
				if exitError.ExitCode() != tt.wantExitCode {
					t.Errorf("Expected exit code %d, got %d", tt.wantExitCode, exitError.ExitCode())
				}
			} else if tt.wantExitCode != 0 {
				t.Errorf("Expected exit code %d, but command succeeded", tt.wantExitCode)
			}

			// Check error message
			stderrStr := stderr.String()
			if !strings.Contains(stderrStr, tt.wantError) {
				t.Errorf("Expected error containing %q, got %q", tt.wantError, stderrStr)
			}

			// Ensure error is not empty when exit code is non-zero
			if tt.wantExitCode != 0 && stderrStr == "" {
				t.Error("Expected error message but stderr was empty")
			}
		})
	}
}