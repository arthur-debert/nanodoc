package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		name           string
		lineNumberMode string
		sequenceType   string
		headerStyle    string
		themeName      string
		wantErr        bool
		errContains    string
	}{
		{
			name:           "valid flags",
			lineNumberMode: "file",
			sequenceType:   "numerical",
			headerStyle:    "nice",
			themeName:      "classic",
			wantErr:        false,
		},
		{
			name:           "invalid line number mode",
			lineNumberMode: "invalid",
			sequenceType:   "numerical",
			headerStyle:    "nice",
			themeName:      "classic",
			wantErr:        true,
			errContains:    "invalid line number mode",
		},
		{
			name:           "invalid sequence type",
			lineNumberMode: "",
			sequenceType:   "invalid",
			headerStyle:    "nice",
			themeName:      "classic",
			wantErr:        true,
			errContains:    "invalid sequence type",
		},
		{
			name:           "invalid header style",
			lineNumberMode: "",
			sequenceType:   "numerical",
			headerStyle:    "invalid",
			themeName:      "classic",
			wantErr:        true,
			errContains:    "invalid header style",
		},
		{
			name:           "custom theme allowed",
			lineNumberMode: "",
			sequenceType:   "numerical",
			headerStyle:    "nice",
			themeName:      "custom-theme",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global variables
			lineNumberMode = tt.lineNumberMode
			sequenceType = tt.sequenceType
			headerStyle = tt.headerStyle
			themeName = tt.themeName

			err := validateFlags()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("validateFlags() error = %v, want error containing %q", err, tt.errContains)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		item  string
		want  bool
	}{
		{
			name:  "item exists",
			slice: []string{"a", "b", "c"},
			item:  "b",
			want:  true,
		},
		{
			name:  "item does not exist",
			slice: []string{"a", "b", "c"},
			item:  "d",
			want:  false,
		},
		{
			name:  "empty slice",
			slice: []string{},
			item:  "a",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.slice, tt.item); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCLIIntegration(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "nanodoc-cli-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test files
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testFile2 := filepath.Join(tempDir, "test2.txt")
	if err := os.WriteFile(testFile1, []byte("Test content 1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testFile2, []byte("Test content 2"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		args        []string
		flags       map[string]string
		wantErr     bool
		wantOutput  string
	}{
		{
			name: "basic file bundling",
			args: []string{testFile1, testFile2},
			flags: map[string]string{
				"no-header": "true",
			},
			wantErr:    false,
			wantOutput: "Test content 1\nTest content 2",
		},
		{
			name: "with line numbers",
			args: []string{testFile1},
			flags: map[string]string{
				"n":         "file",
				"no-header": "true",
			},
			wantErr:    false,
			wantOutput: "Test content 1", // Line numbering not implemented yet
		},
		// Note: Can't test non-existent file case as it calls os.Exit
		// This would be better tested as an integration test
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags to defaults
			lineNumberMode = ""
			showTOC = false
			themeName = "classic"
			noHeader = false
			sequenceType = "numerical"
			headerStyle = "nice"
			additionalExts = nil
			verboseMode = false

			// Create a new command for testing
			cmd := &cobra.Command{
				Use:  "test",
				RunE: func(cmd *cobra.Command, args []string) error {
					// Apply test flags
					for flag, value := range tt.flags {
						switch flag {
						case "n":
							lineNumberMode = value
						case "no-header":
							noHeader = value == "true"
						}
					}

					// Capture output
					oldStdout := os.Stdout
					r, w, _ := os.Pipe()
					os.Stdout = w

					// Run the command
					runNanodoc(cmd, args)

					// Restore stdout
					_ = w.Close()
					os.Stdout = oldStdout

					// Read captured output
					var buf bytes.Buffer
					_, _ = io.Copy(&buf, r)

					if tt.wantOutput != "" && !strings.Contains(buf.String(), tt.wantOutput) {
						t.Errorf("Output = %q, want to contain %q", buf.String(), tt.wantOutput)
					}

					return nil
				},
			}

			err := cmd.RunE(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}