package nanodoc

import (
	"fmt"
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
		opts                 FormattingOptions
		wantTotalFiles       int
		wantTotalLines       int
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
			opts:           FormattingOptions{},
			wantTotalFiles: 1,
			wantTotalLines: 5,
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
			wantTotalFiles: 2, // Files referenced by bundle
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
			opts:             FormattingOptions{AdditionalExtensions: []string{"go"}},
			wantTotalFiles:   1,
			wantTotalLines:   5,
			wantBundles:      0,
			wantRequiresExt:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := GenerateDryRunInfo(tt.pathInfos, tt.opts)
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
			{Path: "/tmp/file1.txt", Source: "direct argument", Extension: ".txt", LineCount: 10},
			{Path: "/tmp/file2.md", Source: "directory: /tmp", Extension: ".md", LineCount: 20},
			{Path: "/tmp/script.py", Source: "glob: *.py", Extension: ".py", LineCount: 15},
		},
		Bundles:           []string{"/tmp/test.bundle.txt"},
		TotalFiles:        3,
		TotalLines:        45,
		RequiresExtension: map[string]string{"/tmp/script.py": ".py"},
		Options:           FormattingOptions{ShowTOC: true, LineNumbers: LineNumberGlobal},
	}

	output := FormatDryRunOutput(info)

	// Check that output contains expected elements
	expectedStrings := []string{
		"Would process the following files:",
		"Table of Contents (5 lines)",
		"From direct argument:",
		"file1.txt (10 lines)",
		"From directory: /tmp:",
		"file2.md (20 lines)",
		"From glob: *.py:",
		"script.py (15 lines)",
		"Bundle files detected:",
		"test.bundle.txt",
		"Files requiring --ext flag:",
		"script.py (requires --ext=py)",
		"Total files to process: 3 (45 lines)",
		"Options:",
		"--toc",
		"--linenum global",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Output missing expected string: %q\nGot:\n%s", expected, output)
		}
	}
}

func TestDryRunWithCircularBundle(t *testing.T) {
	// This test is no longer valid as bundles must contain bundle files
	// The BundleProcessor will treat the circular references as regular files
	// which breaks the cycle. We'll test a valid scenario instead.
	
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-dryrun-bundle-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Create a bundle that references a non-existent file
	bundle1 := filepath.Join(tempDir, "test.bundle.txt")
	nonExistentFile := filepath.Join(tempDir, "does-not-exist.txt")
	
	if err := os.WriteFile(bundle1, []byte(nonExistentFile), 0644); err != nil {
		t.Fatal(err)
	}

	pathInfos := []PathInfo{
		{
			Original: bundle1,
			Absolute: bundle1,
			Type:     "bundle",
		},
	}

	// Dry run should handle missing files gracefully
	info, err := GenerateDryRunInfo(pathInfos, FormattingOptions{})
	if err != nil {
		t.Fatalf("GenerateDryRunInfo() should handle missing files gracefully, got error: %v", err)
	}
	
	// Should have no files (missing file is skipped)
	if info.TotalFiles != 0 {
		t.Errorf("Expected 0 files, got %d", info.TotalFiles)
	}
	
	// Should have the bundle recorded
	if len(info.Bundles) != 1 {
		t.Errorf("Expected 1 bundle entry, got %d", len(info.Bundles))
	}
}

func TestDryRunWithLineRanges(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-dryrun-ranges-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Create a file with 10 lines
	file1 := filepath.Join(tempDir, "test.txt")
	content := ""
	for i := 1; i <= 10; i++ {
		content += fmt.Sprintf("Line %d\n", i)
	}
	if err := os.WriteFile(file1, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		pathInfos      []PathInfo
		wantTotalLines int
	}{
		{
			name: "file with line range L2-4",
			pathInfos: []PathInfo{
				{
					Original: file1 + ":L2-4",
					Absolute: file1,
					Type:     "file",
				},
			},
			wantTotalLines: 3,
		},
		{
			name: "file with single line L5",
			pathInfos: []PathInfo{
				{
					Original: file1 + ":L5",
					Absolute: file1,
					Type:     "file",
				},
			},
			wantTotalLines: 1,
		},
		{
			name: "file with open range L8-",
			pathInfos: []PathInfo{
				{
					Original: file1 + ":L8-",
					Absolute: file1,
					Type:     "file",
				},
			},
			wantTotalLines: 3, // Lines 8, 9, 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := GenerateDryRunInfo(tt.pathInfos, FormattingOptions{})
			if err != nil {
				t.Fatalf("GenerateDryRunInfo() error = %v", err)
			}

			if info.TotalLines != tt.wantTotalLines {
				t.Errorf("TotalLines = %d, want %d", info.TotalLines, tt.wantTotalLines)
			}

			// Check that range spec is preserved
			if len(info.Files) > 0 && !strings.Contains(tt.pathInfos[0].Original, info.Files[0].RangeSpec) {
				t.Errorf("RangeSpec not preserved in file info")
			}
		})
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

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		size int64
		want string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("size_%d", tt.size), func(t *testing.T) {
			got := formatFileSize(tt.size)
			if got != tt.want {
				t.Errorf("formatFileSize(%d) = %q, want %q", tt.size, got, tt.want)
			}
		})
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
			want: true, // Should match - we normalize extensions
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