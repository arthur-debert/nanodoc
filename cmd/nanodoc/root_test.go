package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arthur-debert/nanodoc/pkg/nanodoc"
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
	// Reset flags before each test
	resetFlags()
	
	// Use the actual root command
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
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
			wantOutput: []string{"1. File1", "1 | hello", "2 | world", "2. Title", "3 | # Title", "5 | content"},
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
			name:       "no arguments",
			args:       []string{},
			wantOutput: []string{"Missing paths to bundle: $ nanodoc <path...>"},
			wantErr:    true,
		},
		{
			name:       "help flag",
			args:       []string{"--help"},
			wantOutput: []string{"a minimal document bundler"},
			wantErr:    false,
		},
		{
			name:       "help command",
			args:       []string{"help"},
			wantOutput: []string{"a minimal document bundler"},
			wantErr:    false,
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

func TestRootCmdBundleOptions(t *testing.T) {
	tempDir, cleanup := setupTest(t)
	defer cleanup()

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("line1\nline2\nline3"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create bundle file with options
	bundleFile := filepath.Join(tempDir, "test.bundle.txt")
	bundleContent := []string{
		"# Bundle with options",
		"--toc",
		"--theme classic-dark",
		"--header-style path",
		"--line-numbers",
		"",
		"test.txt",
	}
	if err := os.WriteFile(bundleFile, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		args          []string
		wantOutput    []string
		dontWantOutput []string
	}{
		{
			name: "bundle_options_applied",
			args: []string{bundleFile},
			wantOutput: []string{
				"Table of Contents",  // --toc from bundle
				"1 | line1",          // --line-numbers from bundle
				"2 | line2",          // line numbers continue
				"3 | line3",
				tempDir,              // --header-style path shows full path
			},
		},
		{
			name: "cli_overrides_bundle",
			args: []string{"--header-style", "filename", bundleFile},
			wantOutput: []string{
				"Table of Contents",  // --toc from bundle (not overridden)
				"1 | line1",          // --line-numbers from bundle (not overridden)
				"test.txt",           // --header-style filename overrides bundle's path
			},
			dontWantOutput: []string{
				tempDir,              // Should not show full path
			},
		},
		{
			name: "cli_no_header_overrides_bundle",
			args: []string{"--no-header", bundleFile},
			wantOutput: []string{
				"Table of Contents",  // --toc from bundle (not overridden)
				"1 | line1",          // --line-numbers from bundle (not overridden)
			},
			dontWantOutput: []string{
				"test.txt",           // No header should be shown
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags before each run
			resetFlags()

			output, err := executeCommand(tt.args...)
			if err != nil {
				t.Errorf("executeCommand() error = %v\nOutput:\n%s", err, output)
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
			// Track explicitly set flags
			explicitFlags := make(map[string]bool)
			if cmd.Flags().Changed("toc") {
				explicitFlags["toc"] = true
			}
			if cmd.Flags().Changed("theme") {
				explicitFlags["theme"] = true
			}
			if cmd.Flags().Changed("line-numbers") || cmd.Flags().Changed("global-line-numbers") {
				explicitFlags["line-numbers"] = true
			}
			if cmd.Flags().Changed("no-header") {
				explicitFlags["no-header"] = true
			}
			if cmd.Flags().Changed("header-style") {
				explicitFlags["header-style"] = true
			}
			if cmd.Flags().Changed("sequence") {
				explicitFlags["sequence"] = true
			}
			if cmd.Flags().Changed("txt-ext") {
				explicitFlags["txt-ext"] = true
			}
			if cmd.Flags().Changed("include") {
				explicitFlags["include"] = true
			}
			if cmd.Flags().Changed("exclude") {
				explicitFlags["exclude"] = true
			}

			// Resolve paths with pattern options
			pathOpts := &nanodoc.FormattingOptions{
				AdditionalExtensions: additionalExt,
				IncludePatterns: includePatterns,
				ExcludePatterns: excludePatterns,
			}
			pathInfos, err := nanodoc.ResolvePathsWithOptions(args, pathOpts)
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
				IncludePatterns: includePatterns,
				ExcludePatterns: excludePatterns,
			}

			doc, err := nanodoc.BuildDocumentWithExplicitFlags(pathInfos, opts, explicitFlags)
			if err != nil {
				return fmt.Errorf("Error building document: %w", err)
			}

			ctx, err := nanodoc.NewFormattingContext(doc.FormattingOptions)
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
	cmd.Flags().StringSliceVar(&includePatterns, "include", []string{}, "Include only files matching these patterns")
	cmd.Flags().StringSliceVar(&excludePatterns, "exclude", []string{}, "Exclude files matching these patterns")

	return cmd, &opts
}
