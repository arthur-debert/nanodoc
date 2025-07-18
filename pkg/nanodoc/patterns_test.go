package nanodoc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPatternMatcher(t *testing.T) {
	// Create a test directory structure
	tmpDir := t.TempDir()
	
	// Create test files
	testFiles := map[string]string{
		"api/users.md":           "# Users API",
		"api/auth.md":            "# Auth API",
		"api/test/auth_test.md":  "# Auth API Tests",
		"docs/README.md":         "# README",
		"docs/api-guide.md":      "# API Guide",
		"internal/notes.md":      "# Internal Notes",
		"test/integration.md":    "# Integration Tests",
		"examples/sample.md":     "# Sample",
		"node_modules/pkg/doc.md": "# Package Doc",
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
	
	tests := []struct {
		name            string
		baseDir         string
		includePatterns []string
		excludePatterns []string
		file            string
		wantInclude     bool
		wantRecursion   bool
	}{
		{
			name:        "no patterns - include all",
			baseDir:     tmpDir,
			file:        filepath.Join(tmpDir, "api/users.md"),
			wantInclude: true,
		},
		{
			name:            "include pattern with **",
			baseDir:         tmpDir,
			includePatterns: []string{"**/api/*.md"},
			file:            filepath.Join(tmpDir, "api/users.md"),
			wantInclude:     true,
			wantRecursion:   true,
		},
		{
			name:            "include pattern without match",
			baseDir:         tmpDir,
			includePatterns: []string{"**/api/*.md"},
			file:            filepath.Join(tmpDir, "docs/README.md"),
			wantInclude:     false,
			wantRecursion:   true,
		},
		{
			name:            "exclude pattern",
			baseDir:         tmpDir,
			excludePatterns: []string{"**/README.md"},
			file:            filepath.Join(tmpDir, "docs/README.md"),
			wantInclude:     false,
			wantRecursion:   true,
		},
		{
			name:            "exclude takes precedence",
			baseDir:         tmpDir,
			includePatterns: []string{"**/*.md"},
			excludePatterns: []string{"**/test/*.md"},
			file:            filepath.Join(tmpDir, "api/test/auth_test.md"),
			wantInclude:     false,
			wantRecursion:   true,
		},
		{
			name:            "exclude entire directory",
			baseDir:         tmpDir,
			excludePatterns: []string{"**/node_modules/**"},
			file:            filepath.Join(tmpDir, "node_modules/pkg/doc.md"),
			wantInclude:     false,
			wantRecursion:   true,
		},
		{
			name:            "simple pattern without **",
			baseDir:         tmpDir,
			includePatterns: []string{"api/*.md"},
			file:            filepath.Join(tmpDir, "api/users.md"),
			wantInclude:     true,
			wantRecursion:   false,
		},
		{
			name:            "prefix pattern",
			baseDir:         tmpDir,
			includePatterns: []string{"**/api-*.md"},
			file:            filepath.Join(tmpDir, "docs/api-guide.md"),
			wantInclude:     true,
			wantRecursion:   true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewPatternMatcher(tt.baseDir, tt.includePatterns, tt.excludePatterns)
			
			if got := matcher.NeedsRecursion(); got != tt.wantRecursion {
				t.Errorf("NeedsRecursion() = %v, want %v", got, tt.wantRecursion)
			}
			
			got, err := matcher.ShouldInclude(tt.file)
			if err != nil {
				t.Fatalf("ShouldInclude() error = %v", err)
			}
			if got != tt.wantInclude {
				t.Errorf("ShouldInclude(%s) = %v, want %v", tt.file, got, tt.wantInclude)
			}
		})
	}
}

func TestPatternMatcherHasPatterns(t *testing.T) {
	tests := []struct {
		name            string
		includePatterns []string
		excludePatterns []string
		want            bool
	}{
		{
			name: "no patterns",
			want: false,
		},
		{
			name:            "has include patterns",
			includePatterns: []string{"*.md"},
			want:            true,
		},
		{
			name:            "has exclude patterns",
			excludePatterns: []string{"*.txt"},
			want:            true,
		},
		{
			name:            "has both patterns",
			includePatterns: []string{"*.md"},
			excludePatterns: []string{"*.txt"},
			want:            true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewPatternMatcher("/tmp", tt.includePatterns, tt.excludePatterns)
			if got := matcher.HasPatterns(); got != tt.want {
				t.Errorf("HasPatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}