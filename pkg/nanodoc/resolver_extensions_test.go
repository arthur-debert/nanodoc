package nanodoc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDirectoryWithAdditionalExtensions(t *testing.T) {
	// Create a temporary directory with various files
	tmpDir, err := os.MkdirTemp("", "nanodoc-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	// Create test files with different extensions
	testFiles := map[string]string{
		"file1.txt":  "Content of txt file",
		"file2.md":   "Content of md file",
		"file3.txxt": "Content of txxt file",
		"file4.go":   "Content of go file",
		"file5.py":   "Content of py file",
	}

	for name, content := range testFiles {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name                 string
		options              *FormattingOptions
		expectedFileCount    int
		expectedExtensions   []string
	}{
		{
			name:                 "Default extensions only",
			options:              nil,
			expectedFileCount:    2, // .txt and .md
			expectedExtensions:   []string{".txt", ".md"},
		},
		{
			name: "With txxt extension",
			options: &FormattingOptions{
				AdditionalExtensions: []string{"txxt"},
			},
			expectedFileCount:    3, // .txt, .md, and .txxt
			expectedExtensions:   []string{".txt", ".md", ".txxt"},
		},
		{
			name: "With multiple additional extensions",
			options: &FormattingOptions{
				AdditionalExtensions: []string{"txxt", "go", "py"},
			},
			expectedFileCount:    5, // all files
			expectedExtensions:   []string{".txt", ".md", ".txxt", ".go", ".py"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Resolve the directory with options
			paths, err := ResolvePathsWithOptions([]string{tmpDir}, tt.options)
			if err != nil {
				t.Fatalf("Failed to resolve paths: %v", err)
			}

			if len(paths) != 1 {
				t.Fatalf("Expected 1 path, got %d", len(paths))
			}

			pathInfo := paths[0]
			if pathInfo.Type != "directory" {
				t.Errorf("Expected type 'directory', got '%s'", pathInfo.Type)
			}

			// Check the number of files found
			if len(pathInfo.Files) != tt.expectedFileCount {
				t.Errorf("Expected %d files, got %d", tt.expectedFileCount, len(pathInfo.Files))
				t.Logf("Found files: %v", pathInfo.Files)
			}

			// Verify that files with expected extensions are included
			foundExtensions := make(map[string]bool)
			for _, file := range pathInfo.Files {
				ext := filepath.Ext(file)
				foundExtensions[ext] = true
			}

			for _, expectedExt := range tt.expectedExtensions {
				if !foundExtensions[expectedExt] {
					t.Errorf("Expected to find files with extension %s, but didn't", expectedExt)
				}
			}
		})
	}
}

func TestResolveDirectoryWithAdditionalExtensionsNoOptions(t *testing.T) {
	// Test the specific bug case: directory expansion should respect additional extensions
	tmpDir, err := os.MkdirTemp("", "nanodoc-txxt-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	// Create a .txxt file
	txxtFile := filepath.Join(tmpDir, "spec.txxt")
	if err := os.WriteFile(txxtFile, []byte("Txxt content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a .txt file for comparison
	txtFile := filepath.Join(tmpDir, "spec.txt")
	if err := os.WriteFile(txtFile, []byte("Txt content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test with --txt-ext txxt option
	options := &FormattingOptions{
		AdditionalExtensions: []string{"txxt"},
	}

	paths, err := ResolvePathsWithOptions([]string{tmpDir}, options)
	if err != nil {
		t.Fatalf("Failed to resolve paths: %v", err)
	}

	if len(paths) != 1 {
		t.Fatalf("Expected 1 path, got %d", len(paths))
	}

	// Should find both .txt and .txxt files
	if len(paths[0].Files) != 2 {
		t.Errorf("Expected 2 files (.txt and .txxt), got %d", len(paths[0].Files))
		t.Logf("Found files: %v", paths[0].Files)
	}

	// Check that .txxt file is included
	foundTxxt := false
	for _, file := range paths[0].Files {
		if filepath.Ext(file) == ".txxt" {
			foundTxxt = true
			break
		}
	}

	if !foundTxxt {
		t.Error("Expected to find .txxt file when using --txt-ext txxt, but didn't")
	}
}