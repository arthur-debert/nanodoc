package nanodoc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePathsPreservesOrder(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-order-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Create test files with names that would sort differently alphabetically
	fileZ := filepath.Join(tempDir, "z-file.txt")
	fileA := filepath.Join(tempDir, "a-file.txt")
	fileM := filepath.Join(tempDir, "m-file.txt")
	
	for _, file := range []string{fileZ, fileA, fileM} {
		if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Test cases with different orderings
	tests := []struct {
		name     string
		sources  []string
		wantOrder []string
	}{
		{
			name:     "z-a-m order",
			sources:  []string{fileZ, fileA, fileM},
			wantOrder: []string{fileZ, fileA, fileM},
		},
		{
			name:     "m-z-a order",
			sources:  []string{fileM, fileZ, fileA},
			wantOrder: []string{fileM, fileZ, fileA},
		},
		{
			name:     "a-m-z order (alphabetical)",
			sources:  []string{fileA, fileM, fileZ},
			wantOrder: []string{fileA, fileM, fileZ},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := ResolvePaths(tt.sources)
			if err != nil {
				t.Fatalf("ResolvePaths() error = %v", err)
			}

			if len(results) != len(tt.wantOrder) {
				t.Fatalf("ResolvePaths() returned %d results, want %d", len(results), len(tt.wantOrder))
			}

			// Check that the order is preserved
			for i, result := range results {
				if result.Absolute != tt.wantOrder[i] {
					t.Errorf("ResolvePaths() result[%d] = %s, want %s", i, result.Absolute, tt.wantOrder[i])
				}
			}
		})
	}
}

func TestResolvePathsWithDirectoryPreservesOrder(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-order-dir-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Create subdirectories
	dirA := filepath.Join(tempDir, "a-dir")
	dirZ := filepath.Join(tempDir, "z-dir")
	
	for _, dir := range []string{dirA, dirZ} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create files in directories
	fileInA := filepath.Join(dirA, "file.txt")
	fileInZ := filepath.Join(dirZ, "file.txt")
	standaloneFile := filepath.Join(tempDir, "standalone.txt")
	
	for _, file := range []string{fileInA, fileInZ, standaloneFile} {
		if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Test with mixed directory and file arguments
	sources := []string{dirZ, standaloneFile, dirA}
	
	results, err := ResolvePaths(sources)
	if err != nil {
		t.Fatalf("ResolvePaths() error = %v", err)
	}

	// We expect 3 PathInfo entries in the order: dirZ, standaloneFile, dirA
	if len(results) != 3 {
		t.Fatalf("ResolvePaths() returned %d results, want 3", len(results))
	}

	// Check order - directories should be in the order specified
	expectedOrder := []string{dirZ, standaloneFile, dirA}
	for i, result := range results {
		if result.Absolute != expectedOrder[i] {
			t.Errorf("ResolvePaths() result[%d].Absolute = %s, want %s", i, result.Absolute, expectedOrder[i])
		}
	}
}