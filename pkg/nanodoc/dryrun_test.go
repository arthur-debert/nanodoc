package nanodoc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateDryRunInfo(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-dryrun-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Create test files
	file1 := filepath.Join(tempDir, "test1.txt")
	file2 := filepath.Join(tempDir, "test2.md")
	file3 := filepath.Join(tempDir, "test3.go")
	bundle := filepath.Join(tempDir, "test.bundle.txt")
	
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file3, []byte("content3"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bundle, []byte(file1 + "\n" + file2), 0644); err != nil {
		t.Fatal(err)
	}

	// Create subdirectory with files
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	subFile := filepath.Join(subDir, "sub.txt")
	if err := os.WriteFile(subFile, []byte("sub content"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name                 string
		pathInfos            []PathInfo
		additionalExtensions []string
		wantTotalFiles       int
		wantBundles          int
		wantRequiresExt      int
	}{
		{
			name: "single file",
			pathInfos: []PathInfo{
				{
					Original: file1,
					Absolute: file1,
					Type:     "file",
				},
			},
			wantTotalFiles: 1,
			wantBundles:    0,
		},
		{
			name: "directory",
			pathInfos: []PathInfo{
				{
					Original: tempDir,
					Absolute: tempDir,
					Type:     "directory",
					Files:    []string{file1, file2},
				},
			},
			wantTotalFiles: 2,
			wantBundles:    0,
		},
		{
			name: "glob pattern",
			pathInfos: []PathInfo{
				{
					Original: filepath.Join(tempDir, "*.txt"),
					Absolute: filepath.Join(tempDir, "*.txt"),
					Type:     "glob",
					Files:    []string{file1},
				},
			},
			wantTotalFiles: 1,
			wantBundles:    0,
		},
		{
			name: "bundle file",
			pathInfos: []PathInfo{
				{
					Original: bundle,
					Absolute: bundle,
					Type:     "bundle",
				},
			},
			wantTotalFiles: 1, // Just the bundle file itself, not expanded
			wantBundles:    1,
		},
		{
			name: "file requiring extension",
			pathInfos: []PathInfo{
				{
					Original: file3,
					Absolute: file3,
					Type:     "file",
				},
			},
			wantTotalFiles:  1,
			wantBundles:     0,
			wantRequiresExt: 1,
		},
		{
			name: "file with additional extension",
			pathInfos: []PathInfo{
				{
					Original: file3,
					Absolute: file3,
					Type:     "file",
				},
			},
			additionalExtensions: []string{"go"},
			wantTotalFiles:       1,
			wantBundles:          0,
			wantRequiresExt:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := GenerateDryRunInfo(tt.pathInfos, tt.additionalExtensions)
			if err != nil {
				t.Fatalf("GenerateDryRunInfo() error = %v", err)
			}

			if info.TotalFiles != tt.wantTotalFiles {
				t.Errorf("TotalFiles = %d, want %d", info.TotalFiles, tt.wantTotalFiles)
			}

			if len(info.Bundles) != tt.wantBundles {
				t.Errorf("Bundles count = %d, want %d", len(info.Bundles), tt.wantBundles)
			}

			if len(info.RequiresExtension) != tt.wantRequiresExt {
				t.Errorf("RequiresExtension count = %d, want %d", len(info.RequiresExtension), tt.wantRequiresExt)
			}
		})
	}
}

func TestFormatDryRunOutput(t *testing.T) {
	info := &DryRunInfo{
		Files: []FileInfo{
			{Path: "/tmp/file1.txt", Source: "direct argument", Extension: ".txt"},
			{Path: "/tmp/file2.md", Source: "directory: /tmp", Extension: ".md"},
			{Path: "/tmp/script.py", Source: "glob: *.py", Extension: ".py"},
		},
		Bundles:           []string{"/tmp/test.bundle.txt"},
		TotalFiles:        3,
		RequiresExtension: map[string]string{"/tmp/script.py": ".py"},
	}

	output := FormatDryRunOutput(info)

	// Check that output contains expected elements
	expectedStrings := []string{
		"Would process the following files:",
		"From direct argument:",
		"file1.txt",
		"From directory: /tmp:",
		"file2.md",
		"From glob: *.py:",
		"script.py",
		"Bundle files detected:",
		"test.bundle.txt",
		"Files requiring --txt-ext flag:",
		"script.py (requires --txt-ext=py)",
		"Total files to process: 3",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Output missing expected string: %q\nGot:\n%s", expected, output)
		}
	}
}

func TestDryRunWithCircularBundle(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-dryrun-circular-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Create circular bundle references
	bundle1 := filepath.Join(tempDir, "bundle1.bundle.txt")
	bundle2 := filepath.Join(tempDir, "bundle2.bundle.txt")
	
	if err := os.WriteFile(bundle1, []byte(bundle2), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bundle2, []byte(bundle1), 0644); err != nil {
		t.Fatal(err)
	}

	pathInfos := []PathInfo{
		{
			Original: bundle1,
			Absolute: bundle1,
			Type:     "bundle",
		},
	}

	// Dry run should not fail on circular dependencies (since we don't read bundles now)
	info, err := GenerateDryRunInfo(pathInfos, nil)
	if err != nil {
		t.Fatalf("GenerateDryRunInfo() should not fail, got error: %v", err)
	}

	// Should have the bundle file in files list
	if len(info.Files) != 1 {
		t.Errorf("Expected 1 file entry, got %d", len(info.Files))
	}

	// Should have the bundle recorded
	if len(info.Bundles) != 1 {
		t.Errorf("Expected 1 bundle entry, got %d", len(info.Bundles))
	}
}

func TestDryRunHelperFunctions(t *testing.T) {
	// Test contains
	slice := []string{"go", "py", "js"}
	if !contains(slice, "go") {
		t.Error("contains() should return true for existing item")
	}
	if contains(slice, "rs") {
		t.Error("contains() should return false for non-existing item")
	}

	// Test uniqueStrings
	input := []string{"a", "b", "a", "c", "b", "d"}
	result := uniqueStrings(input)
	
	// Check that we have 4 unique strings
	if len(result) != 4 {
		t.Errorf("uniqueStrings() returned %d items, want 4", len(result))
	}
	
	// Check that each unique string is present
	expected := map[string]bool{"a": true, "b": true, "c": true, "d": true}
	for _, s := range result {
		if !expected[s] {
			t.Errorf("Unexpected string in result: %s", s)
		}
		delete(expected, s)
	}
	if len(expected) > 0 {
		t.Error("Not all expected strings were found in result")
	}
}

func TestIsTextFileWithExtensions(t *testing.T) {
	tests := []struct {
		name                 string
		path                 string
		additionalExtensions []string
		want                 bool
	}{
		{
			name: "default text extension",
			path: "/tmp/file.txt",
			additionalExtensions: []string{},
			want: true,
		},
		{
			name: "default markdown extension",
			path: "/tmp/file.md",
			additionalExtensions: []string{},
			want: true,
		},
		{
			name: "non-text file without additional extensions",
			path: "/tmp/file.go",
			additionalExtensions: []string{},
			want: false,
		},
		{
			name: "non-text file with matching additional extension",
			path: "/tmp/file.go",
			additionalExtensions: []string{"go", "py"},
			want: true,
		},
		{
			name: "non-text file with non-matching additional extension",
			path: "/tmp/file.rs",
			additionalExtensions: []string{"go", "py"},
			want: false,
		},
		{
			name: "case insensitive extension matching",
			path: "/tmp/file.GO",
			additionalExtensions: []string{"go"},
			want: true,
		},
		{
			name: "extension with dot in additional extensions",
			path: "/tmp/file.py",
			additionalExtensions: []string{".py"},
			want: false, // Should not match because we trim the dot
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTextFileWithExtensions(tt.path, tt.additionalExtensions)
			if got != tt.want {
				t.Errorf("isTextFileWithExtensions(%q, %v) = %v, want %v", 
					tt.path, tt.additionalExtensions, got, tt.want)
			}
		})
	}
}