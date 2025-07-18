package nanodoc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIRangeSupport(t *testing.T) {
	// Create test file with numbered lines
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	lines := []string{
		"Line 1",
		"Line 2", 
		"Line 3",
		"Line 4",
		"Line 5",
		"Line 6",
		"Line 7",
		"Line 8",
		"Line 9",
		"Line 10",
	}
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name            string
		pathWithRange   string
		expectedLines   []string
		expectedRange   Range
	}{
		{
			name:          "full file (no range)",
			pathWithRange: testFile,
			expectedLines: lines,
			expectedRange: Range{Start: 1, End: 10},
		},
		{
			name:          "line range L3-5",
			pathWithRange: testFile + ":L3-5",
			expectedLines: []string{"Line 3", "Line 4", "Line 5"},
			expectedRange: Range{Start: 3, End: 5},
		},
		{
			name:          "single line L7",
			pathWithRange: testFile + ":L7",
			expectedLines: []string{"Line 7"},
			expectedRange: Range{Start: 7, End: 7},
		},
		{
			name:          "negative range L$3-$1",
			pathWithRange: testFile + ":L$3-$1",
			expectedLines: []string{"Line 8", "Line 9", "Line 10"},
			expectedRange: Range{Start: 8, End: 10},
		},
		{
			name:          "mixed range L2-$2",
			pathWithRange: testFile + ":L2-$2",
			expectedLines: []string{"Line 2", "Line 3", "Line 4", "Line 5", "Line 6", "Line 7", "Line 8", "Line 9"},
			expectedRange: Range{Start: 2, End: 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through the full pipeline: ResolvePaths -> ExtractFileContent
			pathInfos, err := ResolvePaths([]string{tt.pathWithRange})
			if err != nil {
				t.Fatalf("ResolvePaths error: %v", err)
			}
			
			if len(pathInfos) != 1 {
				t.Fatalf("Expected 1 path info, got %d", len(pathInfos))
			}
			
			// Extract content using the original path (which includes range)
			content, err := ExtractFileContent(pathInfos[0].Original)
			if err != nil {
				t.Fatalf("ExtractFileContent error: %v", err)
			}
			
			// Check the extracted content
			extractedLines := strings.Split(strings.TrimSpace(content.Content), "\n")
			if len(extractedLines) != len(tt.expectedLines) {
				t.Errorf("Expected %d lines, got %d", len(tt.expectedLines), len(extractedLines))
			}
			
			for i, line := range extractedLines {
				if i < len(tt.expectedLines) && line != tt.expectedLines[i] {
					t.Errorf("Line %d: expected %q, got %q", i+1, tt.expectedLines[i], line)
				}
			}
			
			// Check the range
			if len(content.Ranges) != 1 {
				t.Errorf("Expected 1 range, got %d", len(content.Ranges))
			} else {
				if content.Ranges[0].Start != tt.expectedRange.Start || content.Ranges[0].End != tt.expectedRange.End {
					t.Errorf("Expected range %v, got %v", tt.expectedRange, content.Ranges[0])
				}
			}
		})
	}
}