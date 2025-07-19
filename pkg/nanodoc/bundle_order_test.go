package nanodoc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBundlePreservesFileOrder(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-bundle-order-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Create test files with names that would sort differently alphabetically
	fileZ := filepath.Join(tempDir, "z-file.txt")
	fileA := filepath.Join(tempDir, "a-file.txt")
	fileM := filepath.Join(tempDir, "m-file.txt")
	
	// Write distinct content to each file
	if err := os.WriteFile(fileZ, []byte("Content Z"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fileA, []byte("Content A"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fileM, []byte("Content M"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create bundle file with specific order: Z, A, M
	bundleFile := filepath.Join(tempDir, "test.bundle.txt")
	bundleContent := strings.Join([]string{fileZ, fileA, fileM}, "\n")
	if err := os.WriteFile(bundleFile, []byte(bundleContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Process the bundle
	bp := NewBundleProcessor()
	paths, err := bp.ProcessBundleFile(bundleFile)
	if err != nil {
		t.Fatalf("ProcessBundleFile() error = %v", err)
	}

	// Check that the order is preserved
	expectedOrder := []string{fileZ, fileA, fileM}
	if len(paths) != len(expectedOrder) {
		t.Fatalf("ProcessBundleFile() returned %d paths, want %d", len(paths), len(expectedOrder))
	}

	for i, path := range paths {
		if path != expectedOrder[i] {
			t.Errorf("ProcessBundleFile() path[%d] = %s, want %s", i, path, expectedOrder[i])
		}
	}

	// Now test the full document building to ensure order is preserved through the pipeline
	pathInfos := []PathInfo{
		{
			Original: bundleFile,
			Absolute: bundleFile,
			Type:     "bundle",
		},
	}

	doc, err := BuildDocument(pathInfos, FormattingOptions{ShowHeaders: false})
	if err != nil {
		t.Fatalf("BuildDocument() error = %v", err)
	}

	// Check that content appears in the correct order
	if len(doc.ContentItems) != 3 {
		t.Fatalf("BuildDocument() created %d content items, want 3", len(doc.ContentItems))
	}

	// Verify content order
	expectedContents := []string{"Content Z", "Content A", "Content M"}
	for i, item := range doc.ContentItems {
		if strings.TrimSpace(item.Content) != expectedContents[i] {
			t.Errorf("BuildDocument() content[%d] = %q, want %q", i, strings.TrimSpace(item.Content), expectedContents[i])
		}
	}
}

func TestBundleWithMixedSourcesPreservesOrder(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "nanodoc-bundle-mixed-order-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Create test files
	file1 := filepath.Join(tempDir, "1-file.txt")
	file2 := filepath.Join(tempDir, "2-file.txt")
	file3 := filepath.Join(tempDir, "3-file.txt")
	file4 := filepath.Join(tempDir, "4-file.txt")
	
	// Write distinct content
	files := map[string]string{
		file1: "Content 1",
		file2: "Content 2",
		file3: "Content 3",
		file4: "Content 4",
	}
	
	for file, content := range files {
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create bundle with files 3 and 1
	bundleFile := filepath.Join(tempDir, "test.bundle.txt")
	bundleContent := strings.Join([]string{file3, file1}, "\n")
	if err := os.WriteFile(bundleFile, []byte(bundleContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Test with mixed sources: file4, bundle (containing 3,1), file2
	pathInfos := []PathInfo{
		{
			Original: file4,
			Absolute: file4,
			Type:     "file",
		},
		{
			Original: bundleFile,
			Absolute: bundleFile,
			Type:     "bundle",
		},
		{
			Original: file2,
			Absolute: file2,
			Type:     "file",
		},
	}

	doc, err := BuildDocument(pathInfos, FormattingOptions{ShowHeaders: false})
	if err != nil {
		t.Fatalf("BuildDocument() error = %v", err)
	}

	// Expected order: 4, 3, 1, 2
	expectedContents := []string{"Content 4", "Content 3", "Content 1", "Content 2"}
	
	if len(doc.ContentItems) != len(expectedContents) {
		t.Fatalf("BuildDocument() created %d content items, want %d", len(doc.ContentItems), len(expectedContents))
	}

	for i, item := range doc.ContentItems {
		if strings.TrimSpace(item.Content) != expectedContents[i] {
			t.Errorf("BuildDocument() content[%d] = %q, want %q", i, strings.TrimSpace(item.Content), expectedContents[i])
		}
	}
}