package nanodoc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParsePathWithRange(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantPath  string
		wantRange string
	}{
		{
			name:      "path without range",
			input:     "file.txt",
			wantPath:  "file.txt",
			wantRange: "",
		},
		{
			name:      "path with single line",
			input:     "file.txt:L10",
			wantPath:  "file.txt",
			wantRange: "L10",
		},
		{
			name:      "path with range",
			input:     "file.txt:L10-20",
			wantPath:  "file.txt",
			wantRange: "L10-20",
		},
		{
			name:      "path with open-ended range",
			input:     "file.txt:L10-",
			wantPath:  "file.txt",
			wantRange: "L10-",
		},
		{
			name:      "absolute path with range",
			input:     "/path/to/file.txt:L5-15",
			wantPath:  "/path/to/file.txt",
			wantRange: "L5-15",
		},
		{
			name:      "Windows path with range",
			input:     "C:\\path\\to\\file.txt:L1-5",
			wantPath:  "C:\\path\\to\\file.txt",
			wantRange: "L1-5",
		},
		{
			name:      "path with colon but no range",
			input:     "C:\\file.txt",
			wantPath:  "C:\\file.txt",
			wantRange: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotRange := parsePathWithRange(tt.input)
			if gotPath != tt.wantPath {
				t.Errorf("parsePathWithRange() path = %v, want %v", gotPath, tt.wantPath)
			}
			if gotRange != tt.wantRange {
				t.Errorf("parsePathWithRange() range = %v, want %v", gotRange, tt.wantRange)
			}
		})
	}
}



func TestExtractLinesInRange(t *testing.T) {
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

	tests := []struct {
		name  string
		lines []string
		r     *Range
		want  string
	}{
		{
			name:  "extract middle range",
			lines: lines,
			r:     &Range{Start: 3, End: 5},
			want:  "Line 3\nLine 4\nLine 5",
		},
		{
			name:  "extract single line",
			lines: lines,
			r:     &Range{Start: 5, End: 5},
			want:  "Line 5",
		},
		{
			name:  "extract to end",
			lines: lines,
			r:     &Range{Start: 8, End: 0},
			want:  "Line 8\nLine 9\nLine 10",
		},
		{
			name:  "extract full file",
			lines: lines,
			r:     &Range{Start: 1, End: 0},
			want:  strings.Join(lines, "\n"),
		},
		{
			name:  "start beyond file",
			lines: lines,
			r:     &Range{Start: 20, End: 25},
			want:  "",
		},
		{
			name:  "empty lines",
			lines: []string{},
			r:     &Range{Start: 1, End: 10},
			want:  "",
		},
		{
			name:  "start at 0 (should be treated as 1)",
			lines: lines,
			r:     &Range{Start: 0, End: 3},
			want:  "Line 1\nLine 2\nLine 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractLinesInRange(tt.lines, tt.r)
			if got != tt.want {
				t.Errorf("extractLinesInRange() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractFileContent(t *testing.T) {
	// Create temp directory and test file
	tempDir, err := os.MkdirTemp("", "nanodoc-extract-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test file with known content
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := []string{
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
	if err := os.WriteFile(testFile, []byte(strings.Join(testContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		pathWithRange string
		wantContent   string
		wantRange     Range
		wantErr       bool
	}{
		{
			name:          "full file",
			pathWithRange: testFile,
			wantContent:   strings.Join(testContent, "\n"),
			wantRange:     Range{Start: 1, End: 10}, // Now uses actual line count instead of 0
		},
		{
			name:          "single line",
			pathWithRange: testFile + ":L5",
			wantContent:   "Line 5",
			wantRange:     Range{Start: 5, End: 5},
		},
		{
			name:          "line range",
			pathWithRange: testFile + ":L3-7",
			wantContent:   "Line 3\nLine 4\nLine 5\nLine 6\nLine 7",
			wantRange:     Range{Start: 3, End: 7},
		},
		{
			name:          "open-ended range",
			pathWithRange: testFile + ":L8-",
			wantContent:   "Line 8\nLine 9\nLine 10",
			wantRange:     Range{Start: 8, End: 10},
		},
		{
			name:          "non-existent file",
			pathWithRange: filepath.Join(tempDir, "nonexistent.txt"),
			wantErr:       true,
		},
		{
			name:          "invalid range",
			pathWithRange: testFile + ":L20-10",
			wantErr:       true,
		},
		// New tests for $ notation
		{
			name:          "last line with $1",
			pathWithRange: testFile + ":L$1",
			wantContent:   "Line 10",
			wantRange:     Range{Start: 10, End: 10},
		},
		{
			name:          "last 3 lines with $3",
			pathWithRange: testFile + ":L$3",
			wantContent:   "Line 8\nLine 9\nLine 10",
			wantRange:     Range{Start: 8, End: 10},
		},
		{
			name:          "range with negative end",
			pathWithRange: testFile + ":L2-$2",
			wantContent:   "Line 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8\nLine 9",
			wantRange:     Range{Start: 2, End: 9},
		},
		{
			name:          "range with negative start and end",
			pathWithRange: testFile + ":L$5-$2",
			wantContent:   "Line 6\nLine 7\nLine 8\nLine 9",
			wantRange:     Range{Start: 6, End: 9},
		},
		{
			name:          "full file with $1 notation",
			pathWithRange: testFile + ":L1-$1",
			wantContent:   strings.Join(testContent, "\n"),
			wantRange:     Range{Start: 1, End: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractFileContent(tt.pathWithRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractFileContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Content != tt.wantContent {
					t.Errorf("ExtractFileContent() content = %q, want %q", got.Content, tt.wantContent)
				}
				if len(got.Ranges) != 1 || got.Ranges[0] != tt.wantRange {
					t.Errorf("ExtractFileContent() range = %v, want %v", got.Ranges[0], tt.wantRange)
				}
			}
		})
	}
}

func TestResolveAndExtractFiles(t *testing.T) {
	// Create temp directory and test files
	tempDir, err := os.MkdirTemp("", "nanodoc-resolve-extract-test-*")
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
	if err := os.WriteFile(file1, []byte("File 1 content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("File 2 content"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		pathInfos []PathInfo
		wantCount int
		wantErr   bool
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
			wantCount: 1,
		},
		{
			name: "multiple files from directory",
			pathInfos: []PathInfo{
				{
					Original: tempDir,
					Absolute: tempDir,
					Type:     "directory",
					Files:    []string{file1, file2},
				},
			},
			wantCount: 2,
		},
		{
			name: "glob pattern files",
			pathInfos: []PathInfo{
				{
					Original: filepath.Join(tempDir, "*.txt"),
					Type:     "glob",
					Files:    []string{file1, file2},
				},
			},
			wantCount: 2,
		},
		{
			name: "bundle file (not supported yet)",
			pathInfos: []PathInfo{
				{
					Original: filepath.Join(tempDir, "test.bundle.txt"),
					Type:     "bundle",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveAndExtractFiles(tt.pathInfos, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveAndExtractFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("ResolveAndExtractFiles() returned %d items, want %d", len(got), tt.wantCount)
			}
		})
	}
}


