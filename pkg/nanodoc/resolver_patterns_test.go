package nanodoc

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestResolvePathsWithPatterns(t *testing.T) {
	// Create test directory structure
	tmpDir := t.TempDir()
	
	// Create test files
	testFiles := map[string]string{
		"api/users.md":           "# Users API",
		"api/auth.go":            "package auth",
		"api/test/users_test.go": "package auth_test",
		"docs/README.md":         "# README",
		"docs/api-guide.txt":     "API Guide",
		"internal/notes.md":      "# Internal",
		"examples/sample.go":     "package main",
		"test/integration.md":    "# Tests",
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
		source          string
		options         *FormattingOptions
		wantFiles       []string // relative paths from tmpDir
		wantErr         bool
	}{
		{
			name:   "directory without patterns",
			source: filepath.Join(tmpDir, "api"),
			wantFiles: []string{
				"api/users.md",
			},
		},
		{
			name:   "directory with include pattern",
			source: filepath.Join(tmpDir),
			options: &FormattingOptions{
				IncludePatterns: []string{"**/api/*.md"},
			},
			wantFiles: []string{
				"api/users.md",
			},
		},
		{
			name:   "directory with exclude pattern",
			source: filepath.Join(tmpDir),
			options: &FormattingOptions{
				ExcludePatterns: []string{"**/README.md"},
			},
			wantFiles: []string{
				"api/users.md",
				"docs/api-guide.txt", // .txt is a default extension
				"internal/notes.md",
				"test/integration.md",
			},
		},
		{
			name:   "include and exclude patterns",
			source: filepath.Join(tmpDir),
			options: &FormattingOptions{
				IncludePatterns: []string{"**/*.md"},
				ExcludePatterns: []string{"**/test/**", "**/README.md"},
			},
			wantFiles: []string{
				"api/users.md",
				"internal/notes.md",
			},
		},
		{
			name:   "additional extensions with patterns",
			source: filepath.Join(tmpDir),
			options: &FormattingOptions{
				AdditionalExtensions: []string{".go"},
				IncludePatterns:      []string{"**/api/**"},
			},
			wantFiles: []string{
				"api/auth.go",
				"api/test/users_test.go",
				"api/users.md",
			},
		},
		{
			name:   "glob with patterns not supported",
			source: filepath.Join(tmpDir, "**/*.md"),
			options: &FormattingOptions{
				ExcludePatterns: []string{"**/README.md"},
			},
			wantFiles: []string{
				"api/users.md",
				"docs/README.md", // Patterns not applied to glob results
				"internal/notes.md",
				"test/integration.md",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathInfos, err := ResolvePathsWithOptions([]string{tt.source}, tt.options)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ResolvePathsWithOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if tt.wantErr {
				return
			}
			
			// Collect all resolved files
			var gotFiles []string
			for _, info := range pathInfos {
				switch info.Type {
				case "file":
					rel, _ := filepath.Rel(tmpDir, info.Absolute)
					gotFiles = append(gotFiles, rel)
				case "directory", "glob":
					for _, f := range info.Files {
						rel, _ := filepath.Rel(tmpDir, f)
						gotFiles = append(gotFiles, rel)
					}
				}
			}
			
			// Sort for consistent comparison
			sort.Strings(gotFiles)
			sort.Strings(tt.wantFiles)
			
			if len(gotFiles) != len(tt.wantFiles) {
				t.Errorf("Got %d files, want %d files", len(gotFiles), len(tt.wantFiles))
				t.Errorf("Got files: %v", gotFiles)
				t.Errorf("Want files: %v", tt.wantFiles)
				return
			}
			
			for i, got := range gotFiles {
				if got != tt.wantFiles[i] {
					t.Errorf("File[%d] = %v, want %v", i, got, tt.wantFiles[i])
				}
			}
		})
	}
}

func TestDirectoryTraversalWithPatterns(t *testing.T) {
	// Create nested directory structure
	tmpDir := t.TempDir()
	
	// Create files at various depths
	testFiles := map[string]string{
		"README.md":                        "# Root",
		"api/v1/users.md":                  "# Users V1",
		"api/v2/users.md":                  "# Users V2",
		"api/v2/internal/schema.md":        "# Schema",
		"docs/api/reference.md":            "# Reference",
		"docs/guides/quickstart.md":        "# Quickstart",
		"vendor/github.com/foo/README.md":  "# Vendor",
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
		name      string
		patterns  *FormattingOptions
		wantCount int
		wantFiles []string // Sample files to check
	}{
		{
			name:      "no patterns - non-recursive",
			patterns:  &FormattingOptions{},
			wantCount: 1, // Only README.md at root
			wantFiles: []string{"README.md"},
		},
		{
			name: "** pattern enables recursion",
			patterns: &FormattingOptions{
				IncludePatterns: []string{"**/*.md"},
			},
			wantCount: 7, // All .md files
		},
		{
			name: "exclude vendor directory",
			patterns: &FormattingOptions{
				IncludePatterns: []string{"**/*.md"},
				ExcludePatterns: []string{"**/vendor/**"},
			},
			wantCount: 6, // All except vendor
		},
		{
			name: "include only api directories",
			patterns: &FormattingOptions{
				IncludePatterns: []string{"**/api/**/*.md"},
			},
			wantCount: 4, // Only files under api directories
			wantFiles: []string{
				"api/v1/users.md",
				"api/v2/users.md",
				"api/v2/internal/schema.md",
				"docs/api/reference.md",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathInfos, err := ResolvePathsWithOptions([]string{tmpDir}, tt.patterns)
			if err != nil {
				t.Fatalf("ResolvePathsWithOptions() error = %v", err)
			}
			
			var files []string
			for _, info := range pathInfos {
				if info.Type == "directory" {
					files = append(files, info.Files...)
				}
			}
			
			if len(files) != tt.wantCount {
				t.Errorf("Got %d files, want %d", len(files), tt.wantCount)
				for _, f := range files {
					rel, _ := filepath.Rel(tmpDir, f)
					t.Logf("  %s", rel)
				}
			}
			
			// Check specific files if provided
			if len(tt.wantFiles) > 0 {
				fileMap := make(map[string]bool)
				for _, f := range files {
					rel, _ := filepath.Rel(tmpDir, f)
					fileMap[rel] = true
				}
				
				for _, want := range tt.wantFiles {
					if !fileMap[want] {
						t.Errorf("Expected file %s not found", want)
					}
				}
			}
		})
	}
}