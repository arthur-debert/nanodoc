package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCompletionOutput(t *testing.T) {
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "/tmp/nanodoc-test", ".")
	buildCmd.Dir = "/Users/adebert/h/nanodoc/cmd/nanodoc"
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer func() { _ = os.Remove("/tmp/nanodoc-test") }()

	tests := []struct {
		name            string
		args            []string
		wantContains    []string
		wantNotContains []string
	}{
		{
			name: "flag completion",
			args: []string{"__complete", "nanodoc", "-"},
			wantContains: []string{
				"--theme",
				"--header-style",
				"--sequence",
				"--toc",
				"--line-numbers",
			},
		},
		{
			name: "theme value completion",
			args: []string{"__complete", "nanodoc", "--theme", ""},
			wantContains: []string{
				"classic",
				"classic-dark",
				"classic-light",
			},
		},
		{
			name: "header-style value completion",
			args: []string{"__complete", "nanodoc", "--header-style", ""},
			wantContains: []string{
				"nice",
				"simple",
				"path",
				"filename",
				"title",
			},
		},
		{
			name: "sequence value completion",
			args: []string{"__complete", "nanodoc", "--sequence", ""},
			wantContains: []string{
				"numerical",
				"alphabetical",
				"roman",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("/tmp/nanodoc-test", tt.args...)
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out

			if err := cmd.Run(); err != nil {
				// __complete returns exit code 1, which is expected
				if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 {
					t.Fatalf("Command failed with unexpected error: %v\nOutput: %s", err, out.String())
				}
			}

			output := out.String()

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Output does not contain %q.\nGot:\n%s", want, output)
				}
			}

			for _, notWant := range tt.wantNotContains {
				if strings.Contains(output, notWant) {
					t.Errorf("Output contains %q but should not.\nGot:\n%s", notWant, output)
				}
			}
		})
	}
}

func TestValidArgsFunction(t *testing.T) {
	// Test file argument completion
	suggestions, directive := rootCmd.ValidArgsFunction(rootCmd, []string{}, "")
	
	if suggestions != nil {
		t.Errorf("Expected nil suggestions for file completion, got %v", suggestions)
	}
	
	if directive != cobra.ShellCompDirectiveDefault {
		t.Errorf("Expected ShellCompDirectiveDefault, got %v", directive)
	}
}