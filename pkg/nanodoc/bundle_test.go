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
		"--header-style filename",
		"--sequence roman",
		"--global-line-numbers",
		"--txt-ext log",
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

	// Check that options were parsed correctly
	if result.Options.ShowTOC == nil || !*result.Options.ShowTOC {
		t.Error("Expected ShowTOC to be true")
	}
	if result.Options.Theme == nil || *result.Options.Theme != "classic-dark" {
		t.Error("Expected Theme to be 'classic-dark'")
	}
	if result.Options.HeaderStyle == nil || *result.Options.HeaderStyle != HeaderStyleFilename {
		t.Error("Expected HeaderStyle to be 'filename'")
	}
	if result.Options.SequenceStyle == nil || *result.Options.SequenceStyle != SequenceRoman {
		t.Error("Expected SequenceStyle to be 'roman'")
	}
	if result.Options.LineNumbers == nil || *result.Options.LineNumbers != LineNumberGlobal {
		t.Error("Expected LineNumbers to be LineNumberGlobal")
	}
	if len(result.Options.AdditionalExtensions) != 1 || result.Options.AdditionalExtensions[0] != "log" {
		t.Error("Expected AdditionalExtensions to contain 'log'")
	}

	// Check that file paths were parsed correctly
	if len(result.Paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(result.Paths))
	}

	// Test merging options with command-line defaults
	cmdOptions := FormattingOptions{
		Theme:         "classic",
		ShowTOC:       false,
		HeaderStyle:   HeaderStyleNice,
		SequenceStyle: SequenceNumerical,
		LineNumbers:   LineNumberNone,
		ShowHeaders:   true,
		AdditionalExtensions: []string{},
	}

	// Test merging with no explicit flags (bundle options should take precedence)
	explicitFlags := map[string]bool{}
	merged := MergeFormattingOptionsWithDefaults(result.Options, cmdOptions, explicitFlags)
	
	if merged.Theme != "classic-dark" {
		t.Errorf("Expected merged theme to be 'classic-dark', got %s", merged.Theme)
	}
	if !merged.ShowTOC {
		t.Error("Expected merged ShowTOC to be true")
	}
	if merged.HeaderStyle != HeaderStyleFilename {
		t.Errorf("Expected merged HeaderStyle to be 'filename', got %s", merged.HeaderStyle)
	}
	if merged.SequenceStyle != SequenceRoman {
		t.Errorf("Expected merged SequenceStyle to be 'roman', got %s", merged.SequenceStyle)
	}
	if merged.LineNumbers != LineNumberGlobal {
		t.Errorf("Expected merged LineNumbers to be LineNumberGlobal, got %v", merged.LineNumbers)
	}

	// Test command-line flags override bundle options
	explicitFlags["theme"] = true
	explicitFlags["toc"] = true
	merged = MergeFormattingOptionsWithDefaults(result.Options, cmdOptions, explicitFlags)
	
	if merged.Theme != "classic" {
		t.Errorf("Expected merged theme to be 'classic' (CLI override), got %s", merged.Theme)
	}
	if merged.ShowTOC != false {
		t.Error("Expected merged ShowTOC to be false (CLI override)")
	}
	// Header style should still come from bundle since it wasn't explicitly set
	if merged.HeaderStyle != HeaderStyleFilename {
		t.Errorf("Expected merged HeaderStyle to be 'filename', got %s", merged.HeaderStyle)
	}
}

func TestBundleOptionsValidation(t *testing.T) {
	tests := []struct {
		name    string
		option  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid header style",
			option:  "--header-style filename",
			wantErr: false,
		},
		{
			name:    "invalid header style",
			option:  "--header-style invalid",
			wantErr: true,
			errMsg:  "invalid header style: invalid",
		},
		{
			name:    "valid sequence style",
			option:  "--sequence roman",
			wantErr: false,
		},
		{
			name:    "invalid sequence style",
			option:  "--sequence invalid",
			wantErr: true,
			errMsg:  "invalid sequence style: invalid",
		},
		{
			name:    "missing theme value",
			option:  "--theme",
			wantErr: true,
			errMsg:  "--theme requires a value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options BundleOptions
			err := parseOption(tt.option, &options)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
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

func TestParseOption(t *testing.T) {
	tests := []struct {
		name          string
		optionLine    string
		wantError     bool
		expectedField string
		expectedValue interface{}
	}{
		{
			name:          "toc flag",
			optionLine:    "--toc",
			wantError:     false,
			expectedField: "ShowTOC",
			expectedValue: true,
		},
		{
			name:          "no-header flag",
			optionLine:    "--no-header",
			wantError:     false,
			expectedField: "ShowHeaders",
			expectedValue: false,
		},
		{
			name:          "line-numbers flag",
			optionLine:    "--line-numbers",
			wantError:     false,
			expectedField: "LineNumbers",
			expectedValue: LineNumberFile,
		},
		{
			name:          "line-numbers short flag",
			optionLine:    "-n",
			wantError:     false,
			expectedField: "LineNumbers",
			expectedValue: LineNumberFile,
		},
		{
			name:          "global-line-numbers flag",
			optionLine:    "--global-line-numbers",
			wantError:     false,
			expectedField: "LineNumbers",
			expectedValue: LineNumberGlobal,
		},
		{
			name:          "global-line-numbers short flag",
			optionLine:    "-N",
			wantError:     false,
			expectedField: "LineNumbers",
			expectedValue: LineNumberGlobal,
		},
		{
			name:          "theme with value",
			optionLine:    "--theme classic-dark",
			wantError:     false,
			expectedField: "Theme",
			expectedValue: "classic-dark",
		},
		{
			name:          "header-style with value",
			optionLine:    "--header-style path",
			wantError:     false,
			expectedField: "HeaderStyle",
			expectedValue: HeaderStylePath,
		},
		{
			name:          "sequence with value",
			optionLine:    "--sequence roman",
			wantError:     false,
			expectedField: "SequenceStyle",
			expectedValue: SequenceRoman,
		},
		{
			name:          "txt-ext with value",
			optionLine:    "--txt-ext go",
			wantError:     false,
			expectedField: "AdditionalExtensions",
			expectedValue: []string{"go"},
		},
		{
			name:       "theme without value",
			optionLine: "--theme",
			wantError:  true,
		},
		{
			name:       "header-style without value",
			optionLine: "--header-style",
			wantError:  true,
		},
		{
			name:       "sequence without value",
			optionLine: "--sequence",
			wantError:  true,
		},
		{
			name:       "txt-ext without value",
			optionLine: "--txt-ext",
			wantError:  true,
		},
		{
			name:       "unknown flag",
			optionLine: "--unknown-flag",
			wantError:  true,
		},
		{
			name:       "empty option line",
			optionLine: "",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := BundleOptions{}
			err := parseOption(tt.optionLine, &options)
			
			if (err != nil) != tt.wantError {
				t.Errorf("parseOption() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				// Check that the correct field was set
				switch tt.expectedField {
				case "ShowTOC":
					if options.ShowTOC == nil || *options.ShowTOC != tt.expectedValue.(bool) {
						t.Errorf("ShowTOC = %v, want %v", options.ShowTOC, tt.expectedValue)
					}
				case "ShowHeaders":
					if options.ShowHeaders == nil || *options.ShowHeaders != tt.expectedValue.(bool) {
						t.Errorf("ShowHeaders = %v, want %v", options.ShowHeaders, tt.expectedValue)
					}
				case "LineNumbers":
					if options.LineNumbers == nil || *options.LineNumbers != tt.expectedValue.(LineNumberMode) {
						t.Errorf("LineNumbers = %v, want %v", options.LineNumbers, tt.expectedValue)
					}
				case "Theme":
					if options.Theme == nil || *options.Theme != tt.expectedValue.(string) {
						t.Errorf("Theme = %v, want %v", options.Theme, tt.expectedValue)
					}
				case "HeaderStyle":
					if options.HeaderStyle == nil || *options.HeaderStyle != tt.expectedValue.(HeaderStyle) {
						t.Errorf("HeaderStyle = %v, want %v", options.HeaderStyle, tt.expectedValue)
					}
				case "SequenceStyle":
					if options.SequenceStyle == nil || *options.SequenceStyle != tt.expectedValue.(SequenceStyle) {
						t.Errorf("SequenceStyle = %v, want %v", options.SequenceStyle, tt.expectedValue)
					}
				case "AdditionalExtensions":
					expected := tt.expectedValue.([]string)
					if len(options.AdditionalExtensions) != len(expected) {
						t.Errorf("AdditionalExtensions length = %d, want %d", len(options.AdditionalExtensions), len(expected))
					} else {
						for i, ext := range expected {
							if options.AdditionalExtensions[i] != ext {
								t.Errorf("AdditionalExtensions[%d] = %s, want %s", i, options.AdditionalExtensions[i], ext)
							}
						}
					}
				}
			}
		})
	}
}

func TestMergeFormattingOptions(t *testing.T) {
	// Test bundle options merging with command-line options
	bundleOpts := BundleOptions{
		Theme:         &[]string{"classic-dark"}[0],
		LineNumbers:   &[]LineNumberMode{LineNumberGlobal}[0],
		ShowTOC:       &[]bool{true}[0],
		HeaderStyle:   &[]HeaderStyle{HeaderStylePath}[0],
		SequenceStyle: &[]SequenceStyle{SequenceRoman}[0],
		AdditionalExtensions: []string{"go", "py"},
	}

	// Test with default command-line options (should use bundle options)
	cmdOpts := FormattingOptions{
		Theme:         "classic",
		LineNumbers:   LineNumberNone,
		ShowTOC:       false,
		HeaderStyle:   HeaderStyleNice,
		SequenceStyle: SequenceNumerical,
		ShowHeaders:   true,
		AdditionalExtensions: []string{"js"},
	}

	merged := MergeFormattingOptions(bundleOpts, cmdOpts)

	// When command-line options are at default values, bundle options should be used
	if merged.Theme != "classic-dark" {
		t.Errorf("Expected Theme to be 'classic-dark' (from bundle), got %s", merged.Theme)
	}
	if merged.LineNumbers != LineNumberGlobal {
		t.Errorf("Expected LineNumbers to be LineNumberGlobal (from bundle), got %v", merged.LineNumbers)
	}
	if merged.ShowTOC != true {
		t.Errorf("Expected ShowTOC to be true (from bundle), got %v", merged.ShowTOC)
	}
	if merged.HeaderStyle != HeaderStylePath {
		t.Errorf("Expected HeaderStyle to be HeaderStylePath (from bundle), got %v", merged.HeaderStyle)
	}
	if merged.SequenceStyle != SequenceRoman {
		t.Errorf("Expected SequenceStyle to be SequenceRoman (from bundle), got %v", merged.SequenceStyle)
	}
	
	// Additional extensions should be merged
	expectedExtensions := []string{"js", "go", "py"}
	if len(merged.AdditionalExtensions) != len(expectedExtensions) {
		t.Errorf("Expected %d additional extensions, got %d", len(expectedExtensions), len(merged.AdditionalExtensions))
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
		"--global-line-numbers",
		"--header-style nice",
		"--sequence roman",
		"--theme classic-dark",
		"--txt-ext go",
		"--txt-ext py",
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

	// Check options
	if result.Options.ShowTOC == nil || !*result.Options.ShowTOC {
		t.Error("Expected ShowTOC to be true")
	}
	if result.Options.LineNumbers == nil || *result.Options.LineNumbers != LineNumberGlobal {
		t.Error("Expected LineNumbers to be LineNumberGlobal")
	}
	if result.Options.HeaderStyle == nil || *result.Options.HeaderStyle != HeaderStyleNice {
		t.Error("Expected HeaderStyle to be HeaderStyleNice")
	}
	if result.Options.SequenceStyle == nil || *result.Options.SequenceStyle != SequenceRoman {
		t.Error("Expected SequenceStyle to be SequenceRoman")
	}
	if result.Options.Theme == nil || *result.Options.Theme != "classic-dark" {
		t.Error("Expected Theme to be 'classic-dark'")
	}
	if len(result.Options.AdditionalExtensions) != 2 {
		t.Errorf("Expected 2 additional extensions, got %d", len(result.Options.AdditionalExtensions))
	}
	expectedExtensions := []string{"go", "py"}
	for i, ext := range expectedExtensions {
		if result.Options.AdditionalExtensions[i] != ext {
			t.Errorf("AdditionalExtensions[%d] = %s, want %s", i, result.Options.AdditionalExtensions[i], ext)
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
		"--line-numbers",
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

	// Test that bundle options are extracted and merged
	cmdOpts := FormattingOptions{
		Theme:         "classic",
		ShowTOC:       false,
		LineNumbers:   LineNumberNone,
		ShowHeaders:   true,
		HeaderStyle:   HeaderStyleNice,
		SequenceStyle: SequenceNumerical,
	}

	mergedOpts, err := ExtractAndMergeBundleOptions(pathInfos, cmdOpts)
	if err != nil {
		t.Fatal(err)
	}

	// Verify merged options
	if mergedOpts.Theme != "classic-dark" {
		t.Error("Expected theme to be 'classic-dark' from bundle")
	}
	if mergedOpts.ShowTOC != true {
		t.Error("Expected ShowTOC to be true from bundle")
	}
	if mergedOpts.LineNumbers != LineNumberFile {
		t.Error("Expected LineNumbers to be LineNumberFile from bundle")
	}

	// Test that command-line options override bundle options when not at defaults
	cmdOptsWithOverride := FormattingOptions{
		Theme:         "classic-light",  // Not default, should override
		ShowTOC:       false,           // This is default, so bundle option should be used
		LineNumbers:   LineNumberGlobal, // Not default, should override
		ShowHeaders:   true,            // This is default, so bundle option should be used
		HeaderStyle:   HeaderStyleNice, // This is default, so bundle option should be used
		SequenceStyle: SequenceNumerical, // This is default, so bundle option should be used
	}

	mergedOptsWithOverride, err := ExtractAndMergeBundleOptions(pathInfos, cmdOptsWithOverride)
	if err != nil {
		t.Fatal(err)
	}

	// Verify command-line options override bundle options when not at defaults
	if mergedOptsWithOverride.Theme != "classic-light" {
		t.Error("Expected theme to be 'classic-light' from command line (not default)")
	}
	if mergedOptsWithOverride.ShowTOC != true {
		t.Error("Expected ShowTOC to be true from bundle (command line at default)")
	}
	if mergedOptsWithOverride.LineNumbers != LineNumberGlobal {
		t.Error("Expected LineNumbers to be LineNumberGlobal from command line (not default)")
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
		"--global-line-numbers",
		"--header-style nice",
		"--sequence roman",
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

	// Default CLI options (as they would be when no flags are specified)
	defaultCLIOpts := FormattingOptions{
		Theme:         "classic",
		ShowTOC:       false,
		LineNumbers:   LineNumberNone,
		ShowHeaders:   true,
		HeaderStyle:   HeaderStyleNice,
		SequenceStyle: SequenceNumerical,
	}

	// Build document with bundle options
	doc, err := BuildDocument(pathInfos, defaultCLIOpts)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that bundle options were applied
	if doc.FormattingOptions.Theme != "classic-dark" {
		t.Errorf("Expected theme 'classic-dark' from bundle, got %s", doc.FormattingOptions.Theme)
	}
	if doc.FormattingOptions.ShowTOC != true {
		t.Error("Expected ShowTOC to be true from bundle")
	}
	if doc.FormattingOptions.LineNumbers != LineNumberGlobal {
		t.Error("Expected LineNumbers to be LineNumberGlobal from bundle")
	}
	if doc.FormattingOptions.HeaderStyle != HeaderStyleNice {
		t.Error("Expected HeaderStyle to be HeaderStyleNice from bundle")
	}
	if doc.FormattingOptions.SequenceStyle != SequenceRoman {
		t.Error("Expected SequenceStyle to be SequenceRoman from bundle")
	}

	// Test 2: CLI options override bundle options
	overrideCLIOpts := FormattingOptions{
		Theme:         "classic-light", // Should override bundle
		ShowTOC:       false,           // Default, so bundle should win
		LineNumbers:   LineNumberFile,  // Should override bundle
		ShowHeaders:   true,            // Default, so bundle should win
		HeaderStyle:   HeaderStyleFilename, // Should override bundle
		SequenceStyle: SequenceNumerical, // Default, so bundle should win
	}

	doc2, err := BuildDocument(pathInfos, overrideCLIOpts)
	if err != nil {
		t.Fatal(err)
	}

	// Verify CLI overrides
	if doc2.FormattingOptions.Theme != "classic-light" {
		t.Errorf("Expected theme 'classic-light' from CLI override, got %s", doc2.FormattingOptions.Theme)
	}
	if doc2.FormattingOptions.ShowTOC != true {
		t.Error("Expected ShowTOC to be true from bundle (CLI was default)")
	}
	if doc2.FormattingOptions.LineNumbers != LineNumberFile {
		t.Error("Expected LineNumbers to be LineNumberFile from CLI override")
	}
	if doc2.FormattingOptions.HeaderStyle != HeaderStyleFilename {
		t.Error("Expected HeaderStyle to be HeaderStyleFilename from CLI override")
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
		t.Error("Expected output to contain 'i. Intro' due to --sequence roman option")
	}
	if !strings.Contains(output, "1 |") {
		t.Error("Expected output to contain line numbers due to --global-line-numbers option")
	}
}