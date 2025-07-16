package nanodoc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCircularDependencyScenarios(t *testing.T) {
	tests := []struct {
		name        string
		setupFiles  map[string]string
		startFile   string
		wantErrMsg  string
		description string
	}{
		{
			name: "simple_circular_reference",
			setupFiles: map[string]string{
				"a.bundle.txt": "b.bundle.txt",
				"b.bundle.txt": "a.bundle.txt",
			},
			startFile:   "a.bundle.txt",
			wantErrMsg:  "circular dependency detected",
			description: "A includes B, B includes A",
		},
		{
			name: "three_way_circular_reference",
			setupFiles: map[string]string{
				"a.bundle.txt": "b.bundle.txt",
				"b.bundle.txt": "c.bundle.txt",
				"c.bundle.txt": "a.bundle.txt",
			},
			startFile:   "a.bundle.txt",
			wantErrMsg:  "circular dependency detected",
			description: "A -> B -> C -> A",
		},
		{
			name: "self_reference",
			setupFiles: map[string]string{
				"self.bundle.txt": "self.bundle.txt",
			},
			startFile:   "self.bundle.txt",
			wantErrMsg:  "circular dependency detected",
			description: "Bundle includes itself",
		},
		{
			name: "nested_circular_reference",
			setupFiles: map[string]string{
				"main.bundle.txt": "sub/a.bundle.txt\nfile.txt",
				"sub/a.bundle.txt": "b.bundle.txt",
				"sub/b.bundle.txt": "../main.bundle.txt",
				"file.txt": "content",
			},
			startFile:   "main.bundle.txt",
			wantErrMsg:  "circular dependency detected",
			description: "Circular reference through subdirectories",
		},
		{
			name: "valid_diamond_pattern",
			setupFiles: map[string]string{
				"top.bundle.txt":    "left.bundle.txt\nright.bundle.txt",
				"left.bundle.txt":   "bottom.txt",
				"right.bundle.txt":  "bottom.txt",
				"bottom.txt":        "shared content",
			},
			startFile:   "top.bundle.txt",
			wantErrMsg:  "", // Should succeed
			description: "Diamond pattern without circular reference",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir, err := os.MkdirTemp("", "nanodoc-circular-test-*")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = os.RemoveAll(tempDir)
			}()

			// Create subdirectory if needed
			if strings.Contains(tt.name, "nested") {
				subDir := filepath.Join(tempDir, "sub")
				if err := os.Mkdir(subDir, 0755); err != nil {
					t.Fatal(err)
				}
			}

			// Set up test files
			for path, content := range tt.setupFiles {
				fullPath := filepath.Join(tempDir, path)
				dir := filepath.Dir(fullPath)
				if dir != tempDir && dir != "." {
					if err := os.MkdirAll(dir, 0755); err != nil {
						t.Fatal(err)
					}
				}
				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Change to temp directory for relative path tests
			oldWd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = os.Chdir(oldWd)
			}()
			if err := os.Chdir(tempDir); err != nil {
				t.Fatal(err)
			}

			// Test bundle processing
			bp := NewBundleProcessor()
			_, err = bp.ProcessPaths([]string{tt.startFile})

			if tt.wantErrMsg != "" {
				// Expecting an error
				if err == nil {
					t.Fatalf("Expected error containing %q, got nil", tt.wantErrMsg)
				}
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("Error message %q doesn't contain %q", err.Error(), tt.wantErrMsg)
				}
				// Verify it's a CircularDependencyError
				if _, ok := err.(*CircularDependencyError); !ok {
					t.Errorf("Expected CircularDependencyError, got %T", err)
				}
			} else {
				// Should succeed
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestLiveBundleCircularScenarios(t *testing.T) {
	tests := []struct {
		name        string
		setupFiles  map[string]string
		content     string
		wantErrMsg  string
		description string
	}{
		{
			name: "simple_live_circular",
			setupFiles: map[string]string{
				"a.txt": "A content\n[[file:b.txt]]",
				"b.txt": "B content\n[[file:a.txt]]",
			},
			content:     "Start\n[[file:a.txt]]\nEnd",
			wantErrMsg:  "circular dependency detected",
			description: "Live bundle circular reference",
		},
		{
			name: "mixed_circular",
			setupFiles: map[string]string{
				"doc.txt": "Doc with [[file:include.txt]]",
				"include.txt": "Include with [[file:doc.txt]]",
			},
			content:     "Main [[file:doc.txt]]",
			wantErrMsg:  "circular dependency detected",
			description: "Mixed content circular reference",
		},
		{
			name: "deep_nesting",
			setupFiles: map[string]string{
				"1.txt": "Level 1 [[file:2.txt]]",
				"2.txt": "Level 2 [[file:3.txt]]",
				"3.txt": "Level 3 [[file:4.txt]]",
				"4.txt": "Level 4 [[file:5.txt]]",
				"5.txt": "Level 5",
			},
			content:     "Start [[file:1.txt]] End",
			wantErrMsg:  "", // Should succeed with reasonable depth
			description: "Deep nesting without circular reference",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir, err := os.MkdirTemp("", "nanodoc-live-circular-*")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = os.RemoveAll(tempDir)
			}()

			// Set up test files
			for path, content := range tt.setupFiles {
				fullPath := filepath.Join(tempDir, path)
				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Change to temp directory
			oldWd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = os.Chdir(oldWd)
			}()
			if err := os.Chdir(tempDir); err != nil {
				t.Fatal(err)
			}

			// Test live bundle processing
			_, err = ProcessLiveBundle(tt.content)

			if tt.wantErrMsg != "" {
				// Expecting an error
				if err == nil {
					t.Fatalf("Expected error containing %q, got nil", tt.wantErrMsg)
				}
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("Error message %q doesn't contain %q", err.Error(), tt.wantErrMsg)
				}
			} else {
				// Should succeed
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCircularDependencyErrorMessage(t *testing.T) {
	// Test that error messages are helpful
	tempDir, err := os.MkdirTemp("", "nanodoc-error-msg-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Create circular reference
	bundle1 := filepath.Join(tempDir, "project.bundle.txt")
	bundle2 := filepath.Join(tempDir, "includes.bundle.txt")
	
	if err := os.WriteFile(bundle1, []byte("includes.bundle.txt\nREADME.md"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bundle2, []byte("project.bundle.txt\nutils.txt"), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(oldWd)
	}()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	bp := NewBundleProcessor()
	_, err = bp.ProcessPaths([]string{"project.bundle.txt"})

	if err == nil {
		t.Fatal("Expected circular dependency error")
	}

	// Check error message is informative
	errMsg := err.Error()
	if !strings.Contains(errMsg, "circular dependency detected") {
		t.Errorf("Error message should mention circular dependency")
	}
	
	// Verify error type
	circErr, ok := err.(*CircularDependencyError)
	if !ok {
		t.Fatalf("Expected CircularDependencyError, got %T", err)
	}
	
	// Check that we have path and chain information
	if circErr.Path == "" {
		t.Errorf("CircularDependencyError should have a Path")
	}
	if len(circErr.Chain) == 0 {
		t.Errorf("CircularDependencyError should have a Chain")
	}
}