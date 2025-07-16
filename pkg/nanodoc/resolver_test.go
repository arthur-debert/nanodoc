package nanodoc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePaths(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "nanodoc-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test files
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testFile2 := filepath.Join(tempDir, "test2.md")
	testBundle := filepath.Join(tempDir, "test.bundle.txt")
	subDir := filepath.Join(tempDir, "subdir")
	testFile3 := filepath.Join(subDir, "test3.txt")

	// Create files and directories
	if err := os.WriteFile(testFile1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testFile2, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testBundle, []byte("bundle"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testFile3, []byte("content3"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		sources []string
		wantErr bool
		check   func(t *testing.T, results []PathInfo)
	}{
		{
			name:    "empty sources",
			sources: []string{},
			wantErr: true,
		},
		{
			name:    "single file",
			sources: []string{testFile1},
			check: func(t *testing.T, results []PathInfo) {
				if len(results) != 1 {
					t.Errorf("expected 1 result, got %d", len(results))
				}
				if results[0].Type != "file" {
					t.Errorf("expected type 'file', got %s", results[0].Type)
				}
			},
		},
		{
			name:    "bundle file",
			sources: []string{testBundle},
			check: func(t *testing.T, results []PathInfo) {
				if len(results) != 1 {
					t.Errorf("expected 1 result, got %d", len(results))
				}
				if results[0].Type != "bundle" {
					t.Errorf("expected type 'bundle', got %s", results[0].Type)
				}
			},
		},
		{
			name:    "directory",
			sources: []string{tempDir},
			check: func(t *testing.T, results []PathInfo) {
				if len(results) != 1 {
					t.Errorf("expected 1 result, got %d", len(results))
				}
				if results[0].Type != "directory" {
					t.Errorf("expected type 'directory', got %s", results[0].Type)
				}
				// Should find test1.txt, test2.md, and test.bundle.txt
				if len(results[0].Files) != 3 {
					t.Errorf("expected 3 files, got %d", len(results[0].Files))
				}
			},
		},
		{
			name:    "non-existent file",
			sources: []string{filepath.Join(tempDir, "nonexistent.txt")},
			wantErr: true,
		},
		{
			name:    "multiple sources",
			sources: []string{testFile1, testFile2},
			check: func(t *testing.T, results []PathInfo) {
				if len(results) != 2 {
					t.Errorf("expected 2 results, got %d", len(results))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := ResolvePaths(tt.sources)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolvePaths() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, results)
			}
		})
	}
}

func TestResolveGlobPath(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "nanodoc-glob-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test files
	testFile1 := filepath.Join(tempDir, "file1.txt")
	testFile2 := filepath.Join(tempDir, "file2.txt")
	testFile3 := filepath.Join(tempDir, "doc.md")
	testFile4 := filepath.Join(tempDir, "script.py")

	// Create files
	for _, file := range []string{testFile1, testFile2, testFile3, testFile4} {
		if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name    string
		pattern string
		wantErr bool
		check   func(t *testing.T, result PathInfo)
	}{
		{
			name:    "glob txt files",
			pattern: filepath.Join(tempDir, "*.txt"),
			check: func(t *testing.T, result PathInfo) {
				if result.Type != "glob" {
					t.Errorf("expected type 'glob', got %s", result.Type)
				}
				if len(result.Files) != 2 {
					t.Errorf("expected 2 files, got %d", len(result.Files))
				}
			},
		},
		{
			name:    "glob all files",
			pattern: filepath.Join(tempDir, "*"),
			check: func(t *testing.T, result PathInfo) {
				// Should only match .txt and .md files
				if len(result.Files) != 3 {
					t.Errorf("expected 3 files, got %d", len(result.Files))
				}
			},
		},
		{
			name:    "no matches",
			pattern: filepath.Join(tempDir, "*.nonexistent"),
			wantErr: true,
		},
		{
			name:    "single file glob",
			pattern: filepath.Join(tempDir, "file1.*"),
			check: func(t *testing.T, result PathInfo) {
				if len(result.Files) != 1 {
					t.Errorf("expected 1 file, got %d", len(result.Files))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveGlobPath(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveGlobPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, result)
			}
		})
	}
}

func TestIsTextFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"file.txt", true},
		{"file.md", true},
		{"file.TXT", true},
		{"file.MD", true},
		{"file.py", false},
		{"file.go", false},
		{"file", false},
		{"file.txt.bak", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isTextFile(tt.path); got != tt.want {
				t.Errorf("isTextFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsBundleFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"test.bundle.txt", true},
		{"my.bundle.md", true},
		{"/path/to/file.bundle.yaml", true},
		{"bundle.txt", false},
		{"test.txt", false},
		{"test_bundle.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isBundleFile(tt.path); got != tt.want {
				t.Errorf("isBundleFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestSortPaths(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  []string
	}{
		{
			name:  "already sorted",
			paths: []string{"a.txt", "b.txt", "c.txt"},
			want:  []string{"a.txt", "b.txt", "c.txt"},
		},
		{
			name:  "reverse order",
			paths: []string{"c.txt", "b.txt", "a.txt"},
			want:  []string{"a.txt", "b.txt", "c.txt"},
		},
		{
			name:  "mixed order",
			paths: []string{"b.txt", "c.txt", "a.txt"},
			want:  []string{"a.txt", "b.txt", "c.txt"},
		},
		{
			name:  "with paths",
			paths: []string{"/z/file.txt", "/a/file.txt", "/m/file.txt"},
			want:  []string{"/a/file.txt", "/m/file.txt", "/z/file.txt"},
		},
		{
			name:  "empty",
			paths: []string{},
			want:  []string{},
		},
		{
			name:  "single",
			paths: []string{"file.txt"},
			want:  []string{"file.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := make([]string, len(tt.paths))
			copy(paths, tt.paths)
			sortPaths(paths)

			if len(paths) != len(tt.want) {
				t.Errorf("sortPaths() length = %d, want %d", len(paths), len(tt.want))
			}

			for i := range paths {
				if paths[i] != tt.want[i] {
					t.Errorf("sortPaths()[%d] = %s, want %s", i, paths[i], tt.want[i])
				}
			}
		})
	}
}

func TestSymlinkHandling(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "nanodoc-symlink-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create a real file
	realFile := filepath.Join(tempDir, "real.txt")
	if err := os.WriteFile(realFile, []byte("real content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a symlink to the real file
	symlinkFile := filepath.Join(tempDir, "link.txt")
	if err := os.Symlink(realFile, symlinkFile); err != nil {
		t.Skip("Symlinks not supported on this platform")
	}

	// Test resolving the symlink
	result, err := resolveSinglePath(symlinkFile)
	if err != nil {
		t.Fatalf("resolveSinglePath() error = %v", err)
	}

	if result.Type != "file" {
		t.Errorf("expected type 'file', got %s", result.Type)
	}

	// The absolute path should be the symlink path, but it should resolve correctly
	if result.Absolute != symlinkFile {
		t.Errorf("expected absolute path to be %s, got %s", symlinkFile, result.Absolute)
	}
}

func TestGetFilesFromDirectory(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "nanodoc-getfiles-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test files
	files := map[string]string{
		"file1.txt": "content1",
		"file2.md":  "content2",
		"file3.py":  "content3",
		"file4.go":  "content4",
		"file5.rst": "content5",
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tempDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name       string
		extensions []string
		wantCount  int
	}{
		{
			name:       "default extensions",
			extensions: nil,
			wantCount:  2, // .txt and .md
		},
		{
			name:       "custom extensions",
			extensions: []string{".py", ".go"},
			wantCount:  2,
		},
		{
			name:       "single extension",
			extensions: []string{".rst"},
			wantCount:  1,
		},
		{
			name:       "no matches",
			extensions: []string{".java"},
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := GetFilesFromDirectory(tempDir, tt.extensions)
			if err != nil {
				t.Fatalf("GetFilesFromDirectory() error = %v", err)
			}
			if len(files) != tt.wantCount {
				t.Errorf("GetFilesFromDirectory() returned %d files, want %d", len(files), tt.wantCount)
			}
		})
	}
}
