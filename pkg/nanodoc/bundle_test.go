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
		ShowFilenames:   true,
		FilenameStyle:   FilenameStyleNice,
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
	if doc.FormattingOptions.ShowFilenames != options.ShowFilenames {
		t.Errorf("FormattingOptions.ShowFilenames not set correctly")
	}
	if doc.FormattingOptions.FilenameStyle != options.FilenameStyle {
		t.Errorf("FormattingOptions.FilenameStyle not set correctly")
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

func TestBundleOptions(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-bundle-options-test-*")
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
	
	if err := os.WriteFile(file1, []byte("Content of file 1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("Content of file 2"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create bundle file with options
	bundleFile := filepath.Join(tempDir, "test.bundle.txt")
	bundleContent := []string{
		"# Bundle with options",
		"--toc",
		"--theme classic-dark",
		"--file-style filename",
		"--file-numbering roman",
		"--linenum global",
		"--ext log",
		"",
		"# Files to include",
		"file1.txt",
		"file2.txt",
	}
	if err := os.WriteFile(bundleFile, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// Test parsing bundle options
	bp := NewBundleProcessor()
	result, err := bp.ProcessBundleFileWithOptions(bundleFile)
	if err != nil {
		t.Fatalf("ProcessBundleFileWithOptions() error = %v", err)
	}

	// Check that option lines were collected correctly
	expectedOptions := []string{
		"--toc",
		"--theme classic-dark",
		"--file-style filename",
		"--file-numbering roman",
		"--linenum global",
		"--ext log",
	}
	if len(result.OptionLines) != len(expectedOptions) {
		t.Errorf("Expected %d option lines, got %d", len(expectedOptions), len(result.OptionLines))
	}
	for i, expected := range expectedOptions {
		if i < len(result.OptionLines) && result.OptionLines[i] != expected {
			t.Errorf("Expected option line %d to be %q, got %q", i, expected, result.OptionLines[i])
		}
	}

	// Check that file paths were parsed correctly
	if len(result.Paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(result.Paths))
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


func TestProcessBundleFileWithOptions(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-bundle-options-test-*")
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
	
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create bundle file with options and file paths
	bundleFile := filepath.Join(tempDir, "test.bundle.txt")
	bundleContent := []string{
		"# My Project Documentation Bundle",
		"#",
		"# This bundle defines both formatting options and the content to include.",
		"",
		"# --- Options ---",
		"--toc",
		"--linenum global",
		"--file-style nice",
		"--file-numbering roman",
		"--theme classic-dark",
		"--ext go",
		"--ext py",
		"",
		"# --- Content ---",
		"file1.txt",
		"file2.txt",
	}
	if err := os.WriteFile(bundleFile, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// Test processing the bundle file with options
	bp := NewBundleProcessor()
	result, err := bp.ProcessBundleFileWithOptions(bundleFile)
	if err != nil {
		t.Fatalf("ProcessBundleFileWithOptions() error = %v", err)
	}

	// Check paths
	if len(result.Paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(result.Paths))
	}

	// Check option lines
	expectedOptions := []string{
		"--toc",
		"--linenum global",
		"--file-style nice",
		"--file-numbering roman",
		"--theme classic-dark",
		"--ext go",
		"--ext py",
	}
	if len(result.OptionLines) != len(expectedOptions) {
		t.Errorf("Expected %d option lines, got %d", len(expectedOptions), len(result.OptionLines))
	}
	for i, expected := range expectedOptions {
		if i < len(result.OptionLines) && result.OptionLines[i] != expected {
			t.Errorf("Expected option line %d to be %q, got %q", i, expected, result.OptionLines[i])
		}
	}
}

func TestBundleOptionsIntegration(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-integration-test-*")
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
	if err := os.WriteFile(file1, []byte("hello\nworld"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create bundle file with options
	bundleFile := filepath.Join(tempDir, "test.bundle.txt")
	bundleContent := []string{
		"# Test bundle with options",
		"--toc",
		"--theme classic-dark",
		"--linenum file",
		"",
		"file1.txt",
	}
	if err := os.WriteFile(bundleFile, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// Resolve paths
	pathInfos, err := ResolvePaths([]string{bundleFile})
	if err != nil {
		t.Fatal(err)
	}

	// Test that bundle option lines are extracted
	optionLines, err := ExtractBundleOptionLines(pathInfos)
	if err != nil {
		t.Fatal(err)
	}

	// Verify option lines were extracted
	expectedOptions := []string{
		"--toc",
		"--theme classic-dark",
		"--linenum file",
	}
	if len(optionLines) != len(expectedOptions) {
		t.Errorf("Expected %d option lines, got %d", len(expectedOptions), len(optionLines))
	}
	for i, expected := range expectedOptions {
		if i < len(optionLines) && optionLines[i] != expected {
			t.Errorf("Expected option line %d to be %q, got %q", i, expected, optionLines[i])
		}
	}
}

func TestEndToEndBundleOptions(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-e2e-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test files
	file1 := filepath.Join(tempDir, "intro.txt")
	file2 := filepath.Join(tempDir, "chapter1.md")
	file3 := filepath.Join(tempDir, "conclusion.txt")
	
	if err := os.WriteFile(file1, []byte("Introduction to the project\nThis is the intro"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("# Chapter 1\n\nThis is chapter 1 content\nMore content here"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file3, []byte("Final thoughts\nEnd of document"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create bundle file with comprehensive options (as described in the GitHub issue)
	bundleFile := filepath.Join(tempDir, "project.bundle.txt")
	bundleContent := []string{
		"# My Project Documentation Bundle",
		"#",
		"# This bundle defines both formatting options and the content to include.", 
		"# Lines starting with '#' are comments. Empty lines are ignored.",
		"",
		"# --- Options ---",
		"# Options are specified using the same flags as the command line.",
		"",
		"--toc",
		"--linenum global",
		"--file-style nice",
		"--file-numbering roman",
		"--theme classic-dark",
		"",
		"# --- Content ---",
		"# Files, directories, and glob patterns are listed below.",
		"",
		"intro.txt",
		"chapter1.md",
		"conclusion.txt",
	}
	if err := os.WriteFile(bundleFile, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// Test 1: Process bundle file with options only (no CLI overrides)
	pathInfos, err := ResolvePaths([]string{bundleFile})
	if err != nil {
		t.Fatal(err)
	}

	// Options that would result from parsing bundle options in CLI
	// (In real usage, the CLI layer would parse and merge these)
	mergedOpts := FormattingOptions{
		Theme:         "classic-dark",
		ShowTOC:       true,
		LineNumbers:   LineNumberGlobal,
		ShowFilenames:   true,
		FilenameStyle:   FilenameStyleNice,
		SequenceStyle: SequenceRoman,
	}

	// Build document with merged options
	doc, err := BuildDocument(pathInfos, mergedOpts)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that merged options were applied
	if doc.FormattingOptions.Theme != "classic-dark" {
		t.Errorf("Expected theme 'classic-dark' from bundle, got %s", doc.FormattingOptions.Theme)
	}
	if doc.FormattingOptions.ShowTOC != true {
		t.Error("Expected ShowTOC to be true from bundle")
	}
	if doc.FormattingOptions.LineNumbers != LineNumberGlobal {
		t.Error("Expected LineNumbers to be LineNumberGlobal from bundle")
	}
	if doc.FormattingOptions.FilenameStyle != FilenameStyleNice {
		t.Error("Expected FilenameStyle to be FilenameStyleNice from bundle")
	}
	if doc.FormattingOptions.SequenceStyle != SequenceRoman {
		t.Error("Expected SequenceStyle to be SequenceRoman from bundle")
	}

	// Test 2: Simulate CLI options overriding bundle options
	// (In real usage, the CLI layer would handle merging based on explicit flags)
	overrideMergedOpts := FormattingOptions{
		Theme:         "classic-light",     // CLI override
		ShowTOC:       true,                // From bundle (CLI was default)
		LineNumbers:   LineNumberFile,      // CLI override
		ShowFilenames:   true,
		FilenameStyle:   FilenameStyleFilename, // CLI override
		SequenceStyle: SequenceRoman,       // From bundle (CLI was default)
	}

	doc2, err := BuildDocument(pathInfos, overrideMergedOpts)
	if err != nil {
		t.Fatal(err)
	}

	// Verify merged options
	if doc2.FormattingOptions.Theme != "classic-light" {
		t.Errorf("Expected theme 'classic-light' from CLI override, got %s", doc2.FormattingOptions.Theme)
	}
	if doc2.FormattingOptions.ShowTOC != true {
		t.Error("Expected ShowTOC to be true from bundle (CLI was default)")
	}
	if doc2.FormattingOptions.LineNumbers != LineNumberFile {
		t.Error("Expected LineNumbers to be LineNumberFile from CLI override")
	}
	if doc2.FormattingOptions.FilenameStyle != FilenameStyleFilename {
		t.Error("Expected FilenameStyle to be FilenameStyleFilename from CLI override")
	}
	if doc2.FormattingOptions.SequenceStyle != SequenceRoman {
		t.Error("Expected SequenceStyle to be SequenceRoman from bundle (CLI was default)")
	}

	// Test 3: Verify document content is correctly processed
	if len(doc.ContentItems) != 3 {
		t.Errorf("Expected 3 content items, got %d", len(doc.ContentItems))
	}

	// Test 4: Verify rendering works with bundle options
	ctx, err := NewFormattingContext(doc.FormattingOptions)
	if err != nil {
		t.Fatal(err)
	}

	output, err := RenderDocument(doc, ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Verify output contains expected elements based on bundle options
	if !strings.Contains(output, "Table of Contents") {
		t.Error("Expected output to contain 'Table of Contents' due to --toc option")
	}
	if !strings.Contains(output, "i. Intro") {
		t.Error("Expected output to contain 'i. Intro' due to --file-numbering roman option")
	}
	if !strings.Contains(output, "1 |") {
		t.Error("Expected output to contain line numbers due to --linenum global option")
	}
}

func TestBuildDocumentWithExplicitFlags(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-explicit-flags-test-*")
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
	if err := os.WriteFile(file1, []byte("line1\nline2"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create bundle file with options
	bundleFile := filepath.Join(tempDir, "test.bundle.txt")
	bundleContent := []string{
		"# Bundle with options",
		"--toc",
		"--theme classic-dark",
		"--file-style path",
		"--linenum file",
		"",
		"file1.txt",
	}
	if err := os.WriteFile(bundleFile, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// Resolve paths
	pathInfos, err := ResolvePaths([]string{bundleFile})
	if err != nil {
		t.Fatal(err)
	}

	// Test 1: Options that would result from bundle options when no CLI flags are set
	// (In real usage, the CLI layer would parse bundle options and use them)
	mergedOpts := FormattingOptions{
		Theme:         "classic-dark",  // from bundle
		ShowTOC:       true,            // from bundle
		LineNumbers:   LineNumberFile,  // from bundle
		FilenameStyle:   FilenameStylePath, // from bundle
		ShowFilenames:   true,
		SequenceStyle: SequenceNumerical,
	}
	
	explicitFlags := map[string]bool{} // Not used in new architecture
	
	doc, err := BuildDocumentWithExplicitFlags(pathInfos, mergedOpts, explicitFlags)
	if err != nil {
		t.Fatal(err)
	}

	// Verify options were applied
	if doc.FormattingOptions.Theme != "classic-dark" {
		t.Errorf("Expected theme 'classic-dark' from bundle, got %s", doc.FormattingOptions.Theme)
	}
	if doc.FormattingOptions.ShowTOC != true {
		t.Errorf("Expected ShowTOC true from bundle, got %t", doc.FormattingOptions.ShowTOC)
	}
	if doc.FormattingOptions.LineNumbers != LineNumberFile {
		t.Errorf("Expected LineNumbers LineNumberFile from bundle, got %v", doc.FormattingOptions.LineNumbers)
	}
	if doc.FormattingOptions.FilenameStyle != FilenameStylePath {
		t.Errorf("Expected FilenameStyle path from bundle, got %v", doc.FormattingOptions.FilenameStyle)
	}

	// Test 2: Options that would result from CLI overriding some bundle options
	// (In real usage, the CLI layer would merge based on explicit flags)
	mergedOptsWithOverride := FormattingOptions{
		Theme:         "classic-light",     // CLI override
		ShowTOC:       true,                // from bundle (not overridden)
		LineNumbers:   LineNumberGlobal,    // CLI override
		FilenameStyle:   FilenameStyleFilename, // CLI override
		ShowFilenames:   true,
		SequenceStyle: SequenceNumerical,
	}
	
	// Explicit flags not used in new architecture
	explicitFlags = map[string]bool{}
	
	docWithOverride, err := BuildDocumentWithExplicitFlags(pathInfos, mergedOptsWithOverride, explicitFlags)
	if err != nil {
		t.Fatal(err)
	}

	// Verify merged options
	if docWithOverride.FormattingOptions.Theme != "classic-light" {
		t.Errorf("Expected theme 'classic-light' from CLI override, got %s", docWithOverride.FormattingOptions.Theme)
	}
	if docWithOverride.FormattingOptions.LineNumbers != LineNumberGlobal {
		t.Errorf("Expected LineNumbers LineNumberGlobal from CLI override, got %v", docWithOverride.FormattingOptions.LineNumbers)
	}
	if docWithOverride.FormattingOptions.FilenameStyle != FilenameStyleFilename {
		t.Errorf("Expected FilenameStyle filename from CLI override, got %v", docWithOverride.FormattingOptions.FilenameStyle)
	}
	
	// Options from bundle (not overridden)
	if docWithOverride.FormattingOptions.ShowTOC != true {
		t.Errorf("Expected ShowTOC true from bundle (not overridden), got %t", docWithOverride.FormattingOptions.ShowTOC)
	}
}