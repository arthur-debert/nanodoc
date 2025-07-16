package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arthur-debert/nanodoc-go/pkg/nanodoc"
	"github.com/spf13/cobra"
)

// setupTest creates temporary files and returns a cleanup function.
func setupTest(t *testing.T) (string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "nanodoc-cli-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create test files
	if err := os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("hello\nworld"), 0644); err != nil {
		t.Fatalf("Failed to write file1.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "file2.md"), []byte("# Title\n\ncontent"), 0644); err != nil {
		t.Fatalf("Failed to write file2.md: %v", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(tempDir)
	}
	return tempDir, cleanup
}

func executeCommand(args ...string) (string, error) {
	var out bytes.Buffer
	// Create a new root command for each test to prevent flag pollution
	cmd, _ := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return out.String(), err
}

func TestRootCmd(t *testing.T) {
	tempDir, cleanup := setupTest(t)
	defer cleanup()

	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.md")

	tests := []struct {
		name          string
		args          []string
		wantOutput    []string // Substrings to check for in the output
		dontWantOutput []string // Substrings that should NOT be in the output
		wantErr       bool
	}{
		{
			name:       "basic execution with two files",
			args:       []string{file1, file2},
			wantOutput: []string{"File1", "hello", "world", "Title", "content"},
			wantErr:    false,
		},
		{
			name:       "with line numbers",
			args:       []string{"-n", file1},
			wantOutput: []string{"1 | hello", "2 | world"},
			wantErr:    false,
		},
		{
			name:       "with global line numbers",
			args:       []string{"-N", file1, file2},
			wantOutput: []string{"1 | hello", "2 | world", "3 | # Title", "5 | content"},
			wantErr:    false,
		},
		{
			name:       "with table of contents",
			args:       []string{"--toc", file2},
			wantOutput: []string{"Table of Contents", "file2.md", "- Title"},
			wantErr:    false,
		},
		{
			name:          "with no header",
			args:          []string{"--no-header", file1},
			wantOutput:    []string{"hello", "world"},
			dontWantOutput:[]string{"1. File1"},
			wantErr:       false,
		},
		{
			name:       "with dark theme",
			args:       []string{"--theme", "classic-dark", file1},
			wantOutput: []string{"hello", "world"}, // Theme doesn't change text content
			wantErr:    false,
		},
		{
			name:    "no arguments",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "non-existent file",
			args:    []string{"nonexistent.txt"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags before each run
			resetFlags()

			output, err := executeCommand(tt.args...)

			if (err != nil) != tt.wantErr {
				t.Errorf("executeCommand() error = %v, wantErr %v\nOutput:\n%s", err, tt.wantErr, output)
				return
			}

			for _, want := range tt.wantOutput {
				if !strings.Contains(output, want) {
					t.Errorf("Output does not contain %q.\nGot:\n%s", want, output)
				}
			}
			
			for _, dontWant := range tt.dontWantOutput {
				if strings.Contains(output, dontWant) {
					t.Errorf("Output contains %q, but should not.\nGot:\n%s", dontWant, output)
				}
			}
		})
	}
}

// resetFlags resets all persistent flags to their default values.
func resetFlags() {
	lineNumbers = false
	globalLineNumbers = false
	toc = false
	theme = "classic"
	noHeader = false
	sequence = "numerical"
	headerStyle = "nice"
	additionalExt = []string{}
}

func newRootCmd() (*cobra.Command, *nanodoc.FormattingOptions) {
	// Reset all flags
	resetFlags()

	var opts nanodoc.FormattingOptions
	cmd := &cobra.Command{
		Use:   "nanodoc [paths...]",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pathInfos, err := nanodoc.ResolvePaths(args)
			if err != nil {
				return fmt.Errorf("Error resolving paths: %w", err)
			}
			
			lineNumberMode := nanodoc.LineNumberNone
			if globalLineNumbers {
				lineNumberMode = nanodoc.LineNumberGlobal
			} else if lineNumbers {
				lineNumberMode = nanodoc.LineNumberFile
			}

			opts = nanodoc.FormattingOptions{
				LineNumbers:   lineNumberMode,
				ShowTOC:       toc,
				Theme:         theme,
				ShowHeaders:   !noHeader,
				SequenceStyle: nanodoc.SequenceStyle(sequence),
				HeaderStyle:   nanodoc.HeaderStyle(headerStyle),
				AdditionalExtensions: additionalExt,
			}

			doc, err := nanodoc.BuildDocument(pathInfos, opts)
			if err != nil {
				return fmt.Errorf("Error building document: %w", err)
			}

			ctx, err := nanodoc.NewFormattingContext(opts)
			if err != nil {
				return fmt.Errorf("Error creating formatting context: %w", err)
			}

			output, err := nanodoc.RenderDocument(doc, ctx)
			if err != nil {
				return fmt.Errorf("Error rendering document: %w", err)
			}

			_, _ = fmt.Fprint(cmd.OutOrStdout(), output)
			return nil
		},
	}

	// Add flags
	cmd.Flags().BoolVarP(&lineNumbers, "line-numbers", "n", false, "Enable per-file line numbering")
	cmd.Flags().BoolVarP(&globalLineNumbers, "global-line-numbers", "N", false, "Enable global line numbering")
	cmd.MarkFlagsMutuallyExclusive("line-numbers", "global-line-numbers")
	cmd.Flags().BoolVar(&toc, "toc", false, "Generate a table of contents")
	cmd.Flags().StringVar(&theme, "theme", "classic", "Set the theme for formatting")
	cmd.Flags().BoolVar(&noHeader, "no-header", false, "Suppress file headers")
	cmd.Flags().StringVar(&headerStyle, "header-style", "nice", "Set the header style")
	cmd.Flags().StringVar(&sequence, "sequence", "numerical", "Set the sequence style")
	cmd.Flags().StringSliceVar(&additionalExt, "txt-ext", []string{}, "Additional file extensions")

	return cmd, &opts
}
