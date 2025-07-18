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
	
	// Check options were parsed
	if len(result.Options.IncludePatterns) != 1 || result.Options.IncludePatterns[0] != "**/api/*.md" {
		t.Errorf("Expected include pattern '**/api/*.md', got %v", result.Options.IncludePatterns)
	}
	
	if len(result.Options.ExcludePatterns) != 1 || result.Options.ExcludePatterns[0] != "**/*test*" {
		t.Errorf("Expected exclude pattern '**/*test*', got %v", result.Options.ExcludePatterns)
	}
	
	// Test merging with command-line options
	cmdOpts := FormattingOptions{
		Theme: "dark",
		ExcludePatterns: []string{"**/README.md"},
	}
	
	mergedOpts := MergeFormattingOptions(result.Options, cmdOpts)
	
	// Check that patterns were merged
	if len(mergedOpts.IncludePatterns) != 1 {
		t.Errorf("Expected 1 include pattern, got %d", len(mergedOpts.IncludePatterns))
	}
	
	// Command-line excludes should be added to bundle excludes
	if len(mergedOpts.ExcludePatterns) != 2 {
		t.Errorf("Expected 2 exclude patterns, got %d: %v", len(mergedOpts.ExcludePatterns), mergedOpts.ExcludePatterns)
	}
}

func TestParseOptionWithPatterns(t *testing.T) {
	tests := []struct {
		name       string
		optionLine string
		wantInclude []string
		wantExclude []string
		wantErr    bool
	}{
		{
			name:       "include pattern",
			optionLine: "--include **/api/*.md",
			wantInclude: []string{"**/api/*.md"},
		},
		{
			name:       "exclude pattern",
			optionLine: "--exclude **/test/**",
			wantExclude: []string{"**/test/**"},
		},
		{
			name:       "include without value",
			optionLine: "--include",
			wantErr:    true,
		},
		{
			name:       "exclude without value",
			optionLine: "--exclude",
			wantErr:    true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options BundleOptions
			err := parseOption(tt.optionLine, &options)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if len(options.IncludePatterns) != len(tt.wantInclude) {
				t.Errorf("Expected %d include patterns, got %d", len(tt.wantInclude), len(options.IncludePatterns))
			}
			
			if len(options.ExcludePatterns) != len(tt.wantExclude) {
				t.Errorf("Expected %d exclude patterns, got %d", len(tt.wantExclude), len(options.ExcludePatterns))
			}
		})
	}
}