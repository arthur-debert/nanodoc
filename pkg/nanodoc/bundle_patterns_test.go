package nanodoc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBundleWithPatterns(t *testing.T) {
	// Create test directory
	tmpDir := t.TempDir()
	
	// Create test files
	testFiles := map[string]string{
		"api/users.md":      "# Users API",
		"api/auth.md":       "# Auth API", 
		"test/test.md":      "# Test",
		"internal/notes.md": "# Internal",
		"README.md":         "# README",
	}
	
	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	
	// Create bundle file with patterns
	bundleContent := `# Bundle with patterns
--include **/api/*.md
--exclude **/*test*
.
`
	bundlePath := filepath.Join(tmpDir, "docs.bundle.txt")
	if err := os.WriteFile(bundlePath, []byte(bundleContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	// Process bundle
	bp := NewBundleProcessor()
	result, err := bp.ProcessBundleFileWithOptions(bundlePath)
	if err != nil {
		t.Fatalf("ProcessBundleFileWithOptions() error = %v", err)
	}
	
	// Check option lines were collected
	expectedOptions := []string{"--include **/api/*.md", "--exclude **/*test*"}
	if len(result.OptionLines) != len(expectedOptions) {
		t.Errorf("Expected %d option lines, got %d", len(expectedOptions), len(result.OptionLines))
	}
	for i, expected := range expectedOptions {
		if i < len(result.OptionLines) && result.OptionLines[i] != expected {
			t.Errorf("Expected option line %d to be %q, got %q", i, expected, result.OptionLines[i])
		}
	}
}