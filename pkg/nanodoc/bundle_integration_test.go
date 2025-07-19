package nanodoc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBundleOptionsCompleteIntegration tests the complete bundle options feature from issue #17
func TestBundleOptionsCompleteIntegration(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-issue17-integration-*")
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
	file2 := filepath.Join(tempDir, "file2.md")
	file3 := filepath.Join(tempDir, "config.go")
	
	if err := os.WriteFile(file1, []byte("File 1 content\nLine 2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("# Header\nMarkdown content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file3, []byte("package main\n\nfunc main() {}"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create bundle file with comprehensive options as described in issue #17
	bundleFile := filepath.Join(tempDir, "test.bundle.txt")
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
		"--ext go",
		"",
		"# --- Content ---",
		"# Files, directories, and glob patterns are listed below.",
		"",
		"file1.txt",
		"file2.md",
		"config.go",
	}
	if err := os.WriteFile(bundleFile, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// Test 1: Bundle options are correctly parsed and applied
	t.Run("bundle_options_applied", func(t *testing.T) {
		pathInfos := []PathInfo{
			{
				Original: bundleFile,
				Absolute: bundleFile,
				Type:     "bundle",
			},
		}

		// Options that would be the result of parsing and merging bundle options in CLI
		// (In real usage, the CLI layer would parse bundle options and merge them)
		mergedOptions := FormattingOptions{
			Theme:         "classic-dark",
			LineNumbers:   LineNumberGlobal,
			ShowFilenames:   true,
			FilenameStyle:   FilenameStyleNice,
			SequenceStyle: SequenceRoman,
			ShowTOC:       true,
			AdditionalExtensions: []string{"go"},
		}

		// Empty explicit flags (not used in new architecture)
		explicitFlags := make(map[string]bool)

		doc, err := BuildDocumentWithExplicitFlags(pathInfos, mergedOptions, explicitFlags)
		if err != nil {
			t.Fatalf("BuildDocumentWithExplicitFlags() error = %v", err)
		}

		// Check that bundle options were applied
		if doc.FormattingOptions.Theme != "classic-dark" {
			t.Errorf("Expected theme 'classic-dark', got '%s'", doc.FormattingOptions.Theme)
		}
		if doc.FormattingOptions.LineNumbers != LineNumberGlobal {
			t.Errorf("Expected LineNumberGlobal, got %v", doc.FormattingOptions.LineNumbers)
		}
		if doc.FormattingOptions.SequenceStyle != SequenceRoman {
			t.Errorf("Expected SequenceRoman, got %v", doc.FormattingOptions.SequenceStyle)
		}
		if !doc.FormattingOptions.ShowTOC {
			t.Error("Expected ShowTOC to be true")
		}
		if doc.FormattingOptions.FilenameStyle != FilenameStyleNice {
			t.Errorf("Expected FilenameStyleNice, got %v", doc.FormattingOptions.FilenameStyle)
		}
		
		// Check that txt-ext was applied
		found := false
		for _, ext := range doc.FormattingOptions.AdditionalExtensions {
			if ext == "go" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'go' extension to be in AdditionalExtensions")
		}

		// Check that all files were processed (including .go file due to --ext)
		if len(doc.ContentItems) != 3 {
			t.Errorf("Expected 3 content items, got %d", len(doc.ContentItems))
		}
	})

	// Test 2: CLI options override bundle options (issue #17 requirement)
	t.Run("cli_options_override_bundle", func(t *testing.T) {
		pathInfos := []PathInfo{
			{
				Original: bundleFile,
				Absolute: bundleFile,
				Type:     "bundle",
			},
		}

		// Options that would be the result of CLI flags overriding bundle options
		// (In real usage, the CLI layer would handle the merging based on explicit flags)
		mergedOptions := FormattingOptions{
			Theme:         "classic-light",  // CLI override
			LineNumbers:   LineNumberFile,   // CLI override
			ShowFilenames:   true,
			FilenameStyle:   FilenameStyleFilename,  // CLI override
			SequenceStyle: SequenceNumerical,    // CLI override
			ShowTOC:       false,                // CLI override
			AdditionalExtensions: []string{"go"}, // From bundle (not overridden)
		}

		// Explicit flags (not used in new architecture)
		explicitFlags := map[string]bool{}

		doc, err := BuildDocumentWithExplicitFlags(pathInfos, mergedOptions, explicitFlags)
		if err != nil {
			t.Fatalf("BuildDocumentWithExplicitFlags() error = %v", err)
		}

		// Check that CLI options overrode bundle options
		if doc.FormattingOptions.Theme != "classic-light" {
			t.Errorf("Expected CLI theme 'classic-light', got '%s'", doc.FormattingOptions.Theme)
		}
		if doc.FormattingOptions.LineNumbers != LineNumberFile {
			t.Errorf("Expected CLI LineNumberFile, got %v", doc.FormattingOptions.LineNumbers)
		}
		if doc.FormattingOptions.FilenameStyle != FilenameStyleFilename {
			t.Errorf("Expected CLI FilenameStyleFilename, got %v", doc.FormattingOptions.FilenameStyle)
		}
		if doc.FormattingOptions.SequenceStyle != SequenceNumerical {
			t.Errorf("Expected CLI SequenceNumerical, got %v", doc.FormattingOptions.SequenceStyle)
		}
		if doc.FormattingOptions.ShowTOC {
			t.Error("Expected CLI ShowTOC to be false")
		}
		
		// Bundle's txt-ext should still be applied since it wasn't overridden
		found := false
		for _, ext := range doc.FormattingOptions.AdditionalExtensions {
			if ext == "go" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected bundle 'go' extension to still be applied")
		}
	})

	// Test 3: Test end-to-end rendering with bundle options
	t.Run("end_to_end_rendering", func(t *testing.T) {
		pathInfos := []PathInfo{
			{
				Original: bundleFile,
				Absolute: bundleFile,
				Type:     "bundle",
			},
		}

		// Options that would be the result of parsing bundle options in CLI
		mergedOptions := FormattingOptions{
			Theme:         "classic-dark",
			LineNumbers:   LineNumberGlobal,
			ShowFilenames:   true,
			FilenameStyle:   FilenameStyleNice,
			SequenceStyle: SequenceRoman,
			ShowTOC:       true,
			AdditionalExtensions: []string{"go"},
		}

		explicitFlags := make(map[string]bool)

		doc, err := BuildDocumentWithExplicitFlags(pathInfos, mergedOptions, explicitFlags)
		if err != nil {
			t.Fatalf("BuildDocumentWithExplicitFlags() error = %v", err)
		}

		// Create formatting context and render
		ctx, err := NewFormattingContext(doc.FormattingOptions)
		if err != nil {
			t.Fatalf("NewFormattingContext() error = %v", err)
		}

		output, err := RenderDocument(doc, ctx)
		if err != nil {
			t.Fatalf("RenderDocument() error = %v", err)
		}

		// Check that output contains expected elements from bundle options
		if !strings.Contains(output, "Table of Contents") {
			t.Error("Expected TOC to be present in output")
		}
		if !strings.Contains(output, "1 |") {
			t.Error("Expected global line numbers to be present in output")
		}
		if !strings.Contains(output, "i. File1") {
			t.Error("Expected roman sequence style in filenames")
		}
		if !strings.Contains(output, "package main") {
			t.Error("Expected .go file content to be included due to --ext")
		}
	})
}

// TestBundleOptionsEdgeCases tests edge cases for bundle options
func TestBundleOptionsEdgeCases(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-edge-cases-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Test 1: Bundle with only options, no content
	t.Run("options_only_bundle", func(t *testing.T) {
		bundleFile := filepath.Join(tempDir, "options-only.bundle.txt")
		bundleContent := []string{
			"# Options only bundle",
			"--toc",
			"--theme classic-dark",
			"# No content files listed",
		}
		if err := os.WriteFile(bundleFile, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
			t.Fatal(err)
		}

		bp := NewBundleProcessor()
		result, err := bp.ProcessBundleFileWithOptions(bundleFile)
		if err != nil {
			t.Fatalf("ProcessBundleFileWithOptions() error = %v", err)
		}

		// Should have no paths but option lines should be collected
		if len(result.Paths) != 0 {
			t.Errorf("Expected 0 paths, got %d", len(result.Paths))
		}
		// Check that option lines were collected
		expectedOptions := []string{"--toc", "--theme classic-dark"}
		if len(result.OptionLines) != len(expectedOptions) {
			t.Errorf("Expected %d option lines, got %d", len(expectedOptions), len(result.OptionLines))
		}
		for i, expected := range expectedOptions {
			if i < len(result.OptionLines) && result.OptionLines[i] != expected {
				t.Errorf("Expected option line %d to be %q, got %q", i, expected, result.OptionLines[i])
			}
		}
	})

	// Test 2: Bundle with invalid options
	t.Run("invalid_options", func(t *testing.T) {
		bundleFile := filepath.Join(tempDir, "invalid-options.bundle.txt")
		bundleContent := []string{
			"# Bundle with invalid options",
			"--invalid-option",
			"--theme", // Missing value
			"file1.txt",
		}
		if err := os.WriteFile(bundleFile, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
			t.Fatal(err)
		}

		bp := NewBundleProcessor()
		result, err := bp.ProcessBundleFileWithOptions(bundleFile)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// We should collect both option lines (invalid ones too)
		expectedOptions := []string{"--invalid-option", "--theme"}
		if len(result.OptionLines) != len(expectedOptions) {
			t.Errorf("Expected %d option lines, got %d", len(expectedOptions), len(result.OptionLines))
		}
		for i, expected := range expectedOptions {
			if i < len(result.OptionLines) && result.OptionLines[i] != expected {
				t.Errorf("Expected option line %d to be %q, got %q", i, expected, result.OptionLines[i])
			}
		}
	})

	// Test 3: Multiple txt-ext options
	t.Run("multiple_txt_ext", func(t *testing.T) {
		bundleFile := filepath.Join(tempDir, "multiple-ext.bundle.txt")
		bundleContent := []string{
			"# Multiple txt-ext options",
			"--ext go",
			"--ext py",
			"--ext js",
			"file1.txt",
		}
		if err := os.WriteFile(bundleFile, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
			t.Fatal(err)
		}

		bp := NewBundleProcessor()
		result, err := bp.ProcessBundleFileWithOptions(bundleFile)
		if err != nil {
			t.Fatalf("ProcessBundleFileWithOptions() error = %v", err)
		}

		// Should have all three extension option lines
		expectedOptions := []string{"--ext go", "--ext py", "--ext js"}
		if len(result.OptionLines) != len(expectedOptions) {
			t.Errorf("Expected %d option lines, got %d", len(expectedOptions), len(result.OptionLines))
		}
		for i, expected := range expectedOptions {
			if i < len(result.OptionLines) && result.OptionLines[i] != expected {
				t.Errorf("Expected option line %d to be %q, got %q", i, expected, result.OptionLines[i])
			}
		}
	})
}

// TestBundleOptionsDocumentationExample tests the exact example from the issue #17 documentation
func TestBundleOptionsDocumentationExample(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-doc-example-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test files mentioned in the documentation
	readmeFile := filepath.Join(tempDir, "README.md")
	docsDir := filepath.Join(tempDir, "docs")
	pkgDir := filepath.Join(tempDir, "pkg", "nanodoc")
	
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatal(err)
	}
	
	if err := os.WriteFile(readmeFile, []byte("# My Project\nDocumentation"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "design.md"), []byte("# Design\nArchitecture"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create the exact bundle file from the issue #17 documentation
	bundleFile := filepath.Join(tempDir, "bundle.txt")
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
		"README.md",
		"docs/",
		"pkg/nanodoc/*.go",
	}
	if err := os.WriteFile(bundleFile, []byte(strings.Join(bundleContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// Test that the bundle processes correctly
	bp := NewBundleProcessor()
	result, err := bp.ProcessBundleFileWithOptions(bundleFile)
	if err != nil {
		t.Fatalf("ProcessBundleFileWithOptions() error = %v", err)
	}

	// Check that all option lines were collected (in order they appear in bundle)
	expectedOptions := []string{
		"--toc",
		"--linenum global",
		"--file-style nice",
		"--file-numbering roman",
		"--theme classic-dark",
	}
	if len(result.OptionLines) != len(expectedOptions) {
		t.Errorf("Expected %d option lines, got %d", len(expectedOptions), len(result.OptionLines))
	}
	for i, expected := range expectedOptions {
		if i < len(result.OptionLines) && result.OptionLines[i] != expected {
			t.Errorf("Expected option line %d to be %q, got %q", i, expected, result.OptionLines[i])
		}
	}

	// Check that paths were processed
	if len(result.Paths) != 3 {
		t.Errorf("Expected 3 paths, got %d", len(result.Paths))
	}
	
	// Expected paths (relative to bundle file)
	expectedPaths := []string{
		filepath.Join(tempDir, "README.md"),
		filepath.Join(tempDir, "docs"),
		filepath.Join(tempDir, "pkg/nanodoc/*.go"),
	}
	
	for i, expected := range expectedPaths {
		if result.Paths[i] != expected {
			t.Errorf("Expected path %d to be '%s', got '%s'", i, expected, result.Paths[i])
		}
	}
}