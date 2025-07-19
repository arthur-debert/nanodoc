package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCmdWithPatterns(t *testing.T) {
	tempDir, cleanup := setupTest(t)
	defer cleanup()
	
	// Create test directory structure
	testFiles := map[string]string{
		"api/users.md":       "# Users API\nUser endpoints",
		"api/auth.md":        "# Auth API\nAuthentication",
		"api/test/test.md":   "# Test API\nTest file",
		"docs/README.md":     "# README\nDocumentation",
		"docs/guide.md":      "# Guide\nUser guide",
		"internal/notes.md":  "# Internal\nInternal notes",
		"test/unit.md":       "# Unit Tests\nUnit test docs",
	}
	
	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	
	tests := []struct {
		name          string
		args          []string
		wantContains  []string
		dontWant      []string
		wantErr       bool
	}{
		{
			name: "include pattern",
			args: []string{tempDir, "--include", "**/api/*.md"},
			wantContains: []string{
				"Users API",
				"Auth API",
			},
			dontWant: []string{
				"README",
				"Guide",
				"Internal",
				"Test API", // in api/test/
			},
		},
		{
			name: "exclude pattern",
			args: []string{tempDir, "--exclude", "**/README.md", "--exclude", "**/test/**"},
			wantContains: []string{
				"Users API",
				"Auth API",
				"Guide",
				"Internal",
			},
			dontWant: []string{
				"README",
				"Test API",
				"Unit Tests",
			},
		},
		{
			name: "include and exclude combined",
			args: []string{tempDir, "--include", "**/*.md", "--exclude", "**/internal/**", "--exclude", "**/test/**"},
			wantContains: []string{
				"Users API",
				"Auth API",
				"README",
				"Guide",
			},
			dontWant: []string{
				"Internal",
				"Test API",
				"Unit Tests",
			},
		},
		{
			name: "patterns with additional extensions",
			args: []string{tempDir, "--ext", "go", "--include", "**/api/**"},
			wantContains: []string{
				"Users API",
				"Auth API",
				"Test API", // Now included because it's under api/
			},
			dontWant: []string{
				"README",
				"Guide",
				"Internal",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags()
			
			output, err := executeCommand(tt.args...)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("executeCommand() error = %v, wantErr %v\nOutput:\n%s", err, tt.wantErr, output)
				return
			}
			
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Output does not contain %q.\nGot:\n%s", want, output)
				}
			}
			
			for _, dontWant := range tt.dontWant {
				if strings.Contains(output, dontWant) {
					t.Errorf("Output contains %q, but should not.\nGot:\n%s", dontWant, output)
				}
			}
		})
	}
}