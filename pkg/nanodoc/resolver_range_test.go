package nanodoc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePathsWithRanges(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8\nLine 9\nLine 10"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		source   string
		wantType string
		wantErr  bool
	}{
		{
			name:     "file without range",
			source:   testFile,
			wantType: "file",
			wantErr:  false,
		},
		{
			name:     "file with line range",
			source:   testFile + ":L5-10",
			wantType: "file",
			wantErr:  false,
		},
		{
			name:     "file with single line",
			source:   testFile + ":L5",
			wantType: "file",
			wantErr:  false,
		},
		{
			name:     "file with negative range",
			source:   testFile + ":L$5-$1",
			wantType: "file",
			wantErr:  false,
		},
		{
			name:     "file with mixed range",
			source:   testFile + ":L3-$2",
			wantType: "file",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := ResolvePaths([]string{tt.source})
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolvePaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(results) != 1 {
					t.Errorf("ResolvePaths() returned %d results, want 1", len(results))
					return
				}
				if results[0].Type != tt.wantType {
					t.Errorf("ResolvePaths() type = %v, want %v", results[0].Type, tt.wantType)
				}
				// Verify the original path is preserved with range
				if results[0].Original != tt.source {
					t.Errorf("ResolvePaths() original = %v, want %v", results[0].Original, tt.source)
				}
			}
		})
	}
}

func TestResolveNonGlobPathWithRange(t *testing.T) {
	// Create a test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test that the range is preserved in the Original field but stripped for file operations
	pathWithRange := testFile + ":L1-5"
	result, err := resolveNonGlobPath(pathWithRange)
	if err != nil {
		t.Fatalf("resolveNonGlobPath() error = %v", err)
	}

	// Original should contain the range
	if result.Original != pathWithRange {
		t.Errorf("Original = %v, want %v", result.Original, pathWithRange)
	}

	// Absolute should not contain the range
	expectedAbs, _ := filepath.Abs(testFile)
	if result.Absolute != expectedAbs {
		t.Errorf("Absolute = %v, want %v", result.Absolute, expectedAbs)
	}
}