package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	
	// Reset all flag values to ensure clean state
	rootCmd.ResetFlags()
	// Re-initialize flags after reset
	rootCmd.Flags().StringVarP(&lineNum, "linenum", "l", "", FlagLineNum)
	rootCmd.Flags().BoolVar(&toc, "toc", false, FlagTOC)
	rootCmd.Flags().StringVar(&theme, "theme", "classic", FlagTheme)
	rootCmd.Flags().BoolVar(&showFilenames, "filenames", true, FlagFilenames)
	rootCmd.Flags().StringVar(&fileStyle, "file-style", "nice", FlagFileStyle)
	rootCmd.Flags().StringVar(&fileNumbering, "file-numbering", "numerical", FlagFileNumbering)
	rootCmd.Flags().StringSliceVar(&additionalExt, "ext", []string{}, FlagExt)
	rootCmd.Flags().StringSliceVar(&includePatterns, "include", []string{}, FlagInclude)
	rootCmd.Flags().StringSliceVar(&excludePatterns, "exclude", []string{}, FlagExclude)
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, FlagDryRun)
	rootCmd.Flags().StringVar(&saveToBundlePath, "save-to-bundle", "", "Save the current invocation as a bundle file")
	rootCmd.Flags().BoolP("version", "v", false, FlagVersion)
	
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
			args:       []string{"-l", "file", file1},
			wantOutput: []string{"1 | hello", "2 | world"},
			wantErr:    false,
		},
		{
			name:       "with global line numbers",
			args:       []string{"-l", "global", file1, file2},
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
			name:          "without filenames",
			args:          []string{"--filenames=false", file1},
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
			args: []string{"--file-style", "filename", bundleFile},
			wantOutput: []string{
				"Table of Contents",  // --toc from bundle (not overridden)
				"1 | line1",          // --linenum file from bundle (not overridden)
				"test.txt",           // --file-style filename overrides bundle's path
			},
			dontWantOutput: []string{
				tempDir,              // Should not show full path
			},
		},
		{
			name: "cli_no_filenames_overrides_bundle",
			args: []string{"--filenames=false", bundleFile},
			wantOutput: []string{
				"Table of Contents",  // --toc from bundle (not overridden)
				"1 | line1",          // --linenum file from bundle (not overridden)
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
	lineNum = ""
	toc = false
	theme = "classic"
	showFilenames = true
	fileNumbering = "numerical"
	fileStyle = "nice"
	additionalExt = []string{}
	includePatterns = []string{}
	excludePatterns = []string{}
	dryRun = false
	saveToBundlePath = ""
	explicitFlags = make(map[string]bool)
}
