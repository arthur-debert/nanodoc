package nanodoc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessBundleFile(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-bundle-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test files
	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")
	subDir := filepath.Join(tempDir, "subdir")
	file3 := filepath.Join(subDir, "file3.txt")

	// Create files and directories
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file3, []byte("content3"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create bundle file with various types of entries
	bundleFile := filepath.Join(tempDir, "test.bundle.txt")
	bundleContent := []string{
		"# This is a comment",
		"",
		"file1.txt",
		"file2.txt",
		"  # Another comment with spaces",
		"subdir/file3.txt",
		"",
		"# Absolute path",
		file1,
	}
	if err := os.WriteFile(bundleFile, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// Test processing the bundle file
	bp := NewBundleProcessor()
	paths, err := bp.ProcessBundleFile(bundleFile)
	if err != nil {
		t.Fatalf("ProcessBundleFile() error = %v", err)
	}

	// Should have 4 paths (3 relative + 1 absolute)
	if len(paths) != 4 {
		t.Errorf("Expected 4 paths, got %d", len(paths))
	}

	// Check that relative paths were resolved correctly
	expectedPaths := []string{
		filepath.Join(tempDir, "file1.txt"),
		filepath.Join(tempDir, "file2.txt"),
		filepath.Join(tempDir, "subdir/file3.txt"),
		file1, // Absolute path should remain as-is
	}

	for i, expected := range expectedPaths {
		if i < len(paths) && paths[i] != expected {
			t.Errorf("Path[%d] = %q, want %q", i, paths[i], expected)
		}
	}
}

func TestCircularDependencyDetection(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-circular-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create bundle files that reference each other
	bundle1 := filepath.Join(tempDir, "bundle1.bundle.txt")
	bundle2 := filepath.Join(tempDir, "bundle2.bundle.txt")
	bundle3 := filepath.Join(tempDir, "bundle3.bundle.txt")

	// bundle1 -> bundle2
	if err := os.WriteFile(bundle1, []byte("bundle2.bundle.txt"), 0644); err != nil {
		t.Fatal(err)
	}

	// bundle2 -> bundle3
	if err := os.WriteFile(bundle2, []byte("bundle3.bundle.txt"), 0644); err != nil {
		t.Fatal(err)
	}

	// bundle3 -> bundle1 (creates cycle)
	if err := os.WriteFile(bundle3, []byte("bundle1.bundle.txt"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test circular dependency detection
	bp := NewBundleProcessor()
	_, err = bp.ProcessPaths([]string{bundle1})

	if err == nil {
		t.Fatal("Expected circular dependency error, got nil")
	}

	// Check that it's a CircularDependencyError
	if _, ok := err.(*CircularDependencyError); !ok {
		t.Errorf("Expected CircularDependencyError, got %T", err)
	}
}

func TestProcessPaths(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-processpaths-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test files
	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")
	file3 := filepath.Join(tempDir, "file3.txt")

	for _, file := range []string{file1, file2, file3} {
		if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create a bundle file
	bundle := filepath.Join(tempDir, "test.bundle.txt")
	bundleContent := []string{
		"file2.txt",
		"file3.txt",
	}
	if err := os.WriteFile(bundle, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// Test processing mixed paths
	bp := NewBundleProcessor()
	input := []string{file1, bundle}
	expanded, err := bp.ProcessPaths(input)
	if err != nil {
		t.Fatalf("ProcessPaths() error = %v", err)
	}

	// Should have 3 files total (file1 + file2 + file3 from bundle)
	if len(expanded) != 3 {
		t.Errorf("Expected 3 paths, got %d", len(expanded))
	}
}

func TestNestedBundles(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-nested-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test files
	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")
	file3 := filepath.Join(tempDir, "file3.txt")

	for _, file := range []string{file1, file2, file3} {
		if err := os.WriteFile(file, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create nested bundles
	// inner.bundle.txt contains file2.txt and file3.txt
	innerBundle := filepath.Join(tempDir, "inner.bundle.txt")
	innerContent := []string{
		"file2.txt",
		"file3.txt",
	}
	if err := os.WriteFile(innerBundle, []byte(strings.Join(innerContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// outer.bundle.txt contains file1.txt and inner.bundle.txt
	outerBundle := filepath.Join(tempDir, "outer.bundle.txt")
	outerContent := []string{
		"file1.txt",
		"inner.bundle.txt",
	}
	if err := os.WriteFile(outerBundle, []byte(strings.Join(outerContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// Test processing nested bundles
	bp := NewBundleProcessor()
	expanded, err := bp.ProcessPaths([]string{outerBundle})
	if err != nil {
		t.Fatalf("ProcessPaths() error = %v", err)
	}

	// Should have all 3 files
	if len(expanded) != 3 {
		t.Errorf("Expected 3 paths, got %d", len(expanded))
	}
}

func TestBuildDocument(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-builddoc-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test files
	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")

	if err := os.WriteFile(file1, []byte("File 1 content\nLine 2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("File 2 content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a bundle file
	bundle := filepath.Join(tempDir, "test.bundle.txt")
	if err := os.WriteFile(bundle, []byte("file2.txt"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test building document with mixed inputs
	pathInfos := []PathInfo{
		{
			Original: file1,
			Absolute: file1,
			Type:     "file",
		},
		{
			Original: bundle,
			Absolute: bundle,
			Type:     "bundle",
		},
	}

	options := FormattingOptions{
		ShowHeaders:   true,
		HeaderStyle:   HeaderStyleNice,
		LineNumbers:   LineNumberNone,
		SequenceStyle: SequenceNumerical,
	}

	doc, err := BuildDocument(pathInfos, options)
	if err != nil {
		t.Fatalf("BuildDocument() error = %v", err)
	}

	// Should have 2 content items
	if len(doc.ContentItems) != 2 {
		t.Errorf("Expected 2 content items, got %d", len(doc.ContentItems))
	}

	// Check formatting options were applied
	if doc.FormattingOptions.ShowHeaders != options.ShowHeaders {
		t.Errorf("FormattingOptions.ShowHeaders not set correctly")
	}
	if doc.FormattingOptions.HeaderStyle != options.HeaderStyle {
		t.Errorf("FormattingOptions.HeaderStyle not set correctly")
	}
}

func TestEmptyBundle(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-empty-bundle-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create empty bundle file
	emptyBundle := filepath.Join(tempDir, "empty.bundle.txt")
	if err := os.WriteFile(emptyBundle, []byte("# Just comments\n\n# Nothing else"), 0644); err != nil {
		t.Fatal(err)
	}

	bp := NewBundleProcessor()
	paths, err := bp.ProcessBundleFile(emptyBundle)
	if err != nil {
		t.Fatalf("ProcessBundleFile() error = %v", err)
	}

	// Should return empty list, not error
	if len(paths) != 0 {
		t.Errorf("Expected 0 paths from empty bundle, got %d", len(paths))
	}
}

func TestBundleWithMissingFile(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-missing-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create bundle referencing non-existent file
	bundle := filepath.Join(tempDir, "test.bundle.txt")
	if err := os.WriteFile(bundle, []byte("nonexistent.txt"), 0644); err != nil {
		t.Fatal(err)
	}

	// This should succeed - the bundle processor only expands paths
	// The actual file checking happens later in the pipeline
	bp := NewBundleProcessor()
	paths, err := bp.ProcessBundleFile(bundle)
	if err != nil {
		t.Fatalf("ProcessBundleFile() error = %v", err)
	}

	if len(paths) != 1 {
		t.Errorf("Expected 1 path, got %d", len(paths))
	}
}

func TestProcessLiveBundle(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		setupFunc func(string) error
		want      string
		wantErr   bool
	}{
		{
			name: "simple_live_bundle",
			content: `This is a test document.
Here we include a file: [[file:test.txt]]
And here is more text.`,
			setupFunc: func(tempDir string) error {
				return os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("Included content"), 0644)
			},
			want: `This is a test document.
Here we include a file: Included content
And here is more text.`,
			wantErr: false,
		},
		{
			name: "nested_live_bundle",
			content: `Main document
[[file:file1.txt]]
End of main`,
			setupFunc: func(tempDir string) error {
				// file1.txt includes file2.txt
				if err := os.WriteFile(filepath.Join(tempDir, "file1.txt"), 
					[]byte("File 1 start\n[[file:file2.txt]]\nFile 1 end"), 0644); err != nil {
					return err
				}
				// file2.txt has final content
				return os.WriteFile(filepath.Join(tempDir, "file2.txt"), 
					[]byte("File 2 content"), 0644)
			},
			want: `Main document
File 1 start
File 2 content
File 1 end
End of main`,
			wantErr: false,
		},
		{
			name: "non-existent_bundle",
			content: `Document with missing file: [[file:missing.txt]]`,
			setupFunc: func(tempDir string) error {
				return nil
			},
			want:    `Document with missing file: [[file:missing.txt]]`,
			wantErr: false, // Should leave directive as-is
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir, err := os.MkdirTemp("", "nanodoc-live-test-*")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := os.RemoveAll(tempDir); err != nil {
					t.Logf("Failed to remove temp dir: %v", err)
				}
			}()

			// Change to temp directory to resolve relative paths
			oldDir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Chdir(tempDir); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := os.Chdir(oldDir); err != nil {
					t.Logf("Failed to change back to original dir: %v", err)
				}
			}()

			// Run setup function
			if tt.setupFunc != nil {
				if err := tt.setupFunc(tempDir); err != nil {
					t.Fatal(err)
				}
			}

			got, err := ProcessLiveBundle(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessLiveBundle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProcessLiveBundle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProcessLiveBundleWithRanges(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-live-range-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create a file with multiple lines
	testFile := filepath.Join(tempDir, "multiline.txt")
	content := []string{
		"Line 1",
		"Line 2",
		"Line 3",
		"Line 4",
		"Line 5",
	}
	if err := os.WriteFile(testFile, []byte(strings.Join(content, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Logf("Failed to change back to original dir: %v", err)
		}
	}()

	// Test with range
	input := `Include lines 2-4: [[file:multiline.txt:L2-4]]`
	want := `Include lines 2-4: Line 2
Line 3
Line 4`

	got, err := ProcessLiveBundle(input)
	if err != nil {
		t.Fatalf("ProcessLiveBundle() error = %v", err)
	}
	if got != want {
		t.Errorf("ProcessLiveBundle() = %q, want %q", got, want)
	}
}

func TestProcessLiveBundleCircularReference(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-live-circular-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create files that reference each other
	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")

	// file1 includes file2
	if err := os.WriteFile(file1, []byte("File 1\n[[file:file2.txt]]"), 0644); err != nil {
		t.Fatal(err)
	}

	// file2 includes file1 (circular)
	if err := os.WriteFile(file2, []byte("File 2\n[[file:file1.txt]]"), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Logf("Failed to change back to original dir: %v", err)
		}
	}()

	// Test circular reference detection
	_, err = ProcessLiveBundle("Start\n[[file:file1.txt]]\nEnd")
	if err == nil {
		t.Fatal("Expected circular dependency error, got nil")
	}

	// Check that it's a CircularDependencyError
	if _, ok := err.(*CircularDependencyError); !ok {
		t.Errorf("Expected CircularDependencyError, got %T", err)
	}
}