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

func TestParseRange(t *testing.T) {
	tests := []struct {
		name       string
		spec       string
		totalLines int
		wantStart  int
		wantEnd    int
		wantErr    bool
	}{
		{
			name:       "single line",
			spec:       "L10",
			totalLines: 100,
			wantStart:  10,
			wantEnd:    10,
		},
		{
			name:       "line range",
			spec:       "L10-20",
			totalLines: 100,
			wantStart:  10,
			wantEnd:    20,
		},
		{
			name:       "open-ended range",
			spec:       "L10-",
			totalLines: 100,
			wantStart:  10,
			wantEnd:    0,
		},
		{
			name:       "invalid format - no L prefix",
			spec:       "10-20",
			totalLines: 100,
			wantErr:    true,
		},
		{
			name:       "invalid format - multiple dashes",
			spec:       "L10-20-30",
			totalLines: 100,
			wantErr:    true,
		},
		{
			name:       "invalid start line",
			spec:       "Labc-20",
			totalLines: 100,
			wantErr:    true,
		},
		{
			name:       "invalid end line",
			spec:       "L10-abc",
			totalLines: 100,
			wantErr:    true,
		},
		{
			name:       "negative line number",
			spec:       "L-5",
			totalLines: 100,
			wantErr:    true,
		},
		{
			name:       "end before start",
			spec:       "L20-10",
			totalLines: 100,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRange(tt.spec, tt.totalLines)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Start != tt.wantStart {
					t.Errorf("parseRange() Start = %v, want %v", got.Start, tt.wantStart)
				}
				if got.End != tt.wantEnd {
					t.Errorf("parseRange() End = %v, want %v", got.End, tt.wantEnd)
				}
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
			wantRange:     Range{Start: 1, End: 0},
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
			wantRange:     Range{Start: 8, End: 0},
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

func TestMergeRanges(t *testing.T) {
	tests := []struct {
		name   string
		ranges []Range
		want   []Range
	}{
		{
			name:   "empty ranges",
			ranges: []Range{},
			want:   []Range{},
		},
		{
			name:   "single range",
			ranges: []Range{{Start: 1, End: 10}},
			want:   []Range{{Start: 1, End: 10}},
		},
		{
			name:   "non-overlapping ranges",
			ranges: []Range{{Start: 1, End: 5}, {Start: 10, End: 15}},
			want:   []Range{{Start: 1, End: 5}, {Start: 10, End: 15}},
		},
		{
			name:   "overlapping ranges",
			ranges: []Range{{Start: 1, End: 10}, {Start: 5, End: 15}},
			want:   []Range{{Start: 1, End: 15}},
		},
		{
			name:   "adjacent ranges",
			ranges: []Range{{Start: 1, End: 5}, {Start: 6, End: 10}},
			want:   []Range{{Start: 1, End: 10}},
		},
		{
			name:   "multiple overlapping ranges",
			ranges: []Range{{Start: 1, End: 5}, {Start: 3, End: 8}, {Start: 7, End: 12}},
			want:   []Range{{Start: 1, End: 12}},
		},
		{
			name:   "unsorted ranges",
			ranges: []Range{{Start: 10, End: 15}, {Start: 1, End: 5}, {Start: 20, End: 25}},
			want:   []Range{{Start: 1, End: 5}, {Start: 10, End: 15}, {Start: 20, End: 25}},
		},
		{
			name:   "range with EOF",
			ranges: []Range{{Start: 10, End: 20}, {Start: 15, End: 0}},
			want:   []Range{{Start: 10, End: 0}},
		},
		{
			name:   "multiple ranges with EOF",
			ranges: []Range{{Start: 1, End: 5}, {Start: 10, End: 0}, {Start: 20, End: 0}},
			want:   []Range{{Start: 1, End: 5}, {Start: 10, End: 0}},
		},
		{
			name:   "contained ranges",
			ranges: []Range{{Start: 1, End: 20}, {Start: 5, End: 10}, {Start: 12, End: 15}},
			want:   []Range{{Start: 1, End: 20}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeRanges(tt.ranges)
			if len(got) != len(tt.want) {
				t.Errorf("MergeRanges() returned %d ranges, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("MergeRanges()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestGatherContentWithRanges(t *testing.T) {
	// Create temp directory and test files
	tempDir, err := os.MkdirTemp("", "nanodoc-gather-test-*")
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

	// Create another test file
	testFile2 := filepath.Join(tempDir, "test2.txt")
	testContent2 := []string{
		"File2 Line 1",
		"File2 Line 2",
		"File2 Line 3",
	}
	if err := os.WriteFile(testFile2, []byte(strings.Join(testContent2, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		contents     []FileContent
		wantCount    int
		checkContent func(t *testing.T, results []FileContent)
		wantErr      bool
	}{
		{
			name: "single file single range",
			contents: []FileContent{
				{
					Filepath: testFile,
					Ranges:   []Range{{Start: 1, End: 3}},
				},
			},
			wantCount: 1,
			checkContent: func(t *testing.T, results []FileContent) {
				if results[0].Content != "Line 1\nLine 2\nLine 3" {
					t.Errorf("Expected content doesn't match")
				}
			},
		},
		{
			name: "single file multiple ranges",
			contents: []FileContent{
				{
					Filepath: testFile,
					Ranges:   []Range{{Start: 1, End: 3}},
				},
				{
					Filepath: testFile,
					Ranges:   []Range{{Start: 5, End: 7}},
				},
			},
			wantCount: 1,
			checkContent: func(t *testing.T, results []FileContent) {
				want := "Line 1\nLine 2\nLine 3\nLine 5\nLine 6\nLine 7"
				if results[0].Content != want {
					t.Errorf("Expected content %q, got %q", want, results[0].Content)
				}
				if len(results[0].Ranges) != 2 {
					t.Errorf("Expected 2 ranges, got %d", len(results[0].Ranges))
				}
			},
		},
		{
			name: "single file overlapping ranges",
			contents: []FileContent{
				{
					Filepath: testFile,
					Ranges:   []Range{{Start: 1, End: 5}},
				},
				{
					Filepath: testFile,
					Ranges:   []Range{{Start: 3, End: 7}},
				},
			},
			wantCount: 1,
			checkContent: func(t *testing.T, results []FileContent) {
				// Should merge to a single range 1-7
				want := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7"
				if results[0].Content != want {
					t.Errorf("Expected content %q, got %q", want, results[0].Content)
				}
				if len(results[0].Ranges) != 1 {
					t.Errorf("Expected 1 merged range, got %d", len(results[0].Ranges))
				}
			},
		},
		{
			name: "multiple files",
			contents: []FileContent{
				{
					Filepath: testFile,
					Ranges:   []Range{{Start: 1, End: 3}},
				},
				{
					Filepath: testFile2,
					Ranges:   []Range{{Start: 2, End: 3}},
				},
			},
			wantCount: 2,
			checkContent: func(t *testing.T, results []FileContent) {
				// Find each file in results
				for _, r := range results {
					switch r.Filepath {
					case testFile:
						if r.Content != "Line 1\nLine 2\nLine 3" {
							t.Errorf("File1 content doesn't match")
						}
					case testFile2:
						if r.Content != "File2 Line 2\nFile2 Line 3" {
							t.Errorf("File2 content doesn't match")
						}
					}
				}
			},
		},
		{
			name: "non-existent file",
			contents: []FileContent{
				{
					Filepath: filepath.Join(tempDir, "nonexistent.txt"),
					Ranges:   []Range{{Start: 1, End: 3}},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GatherContentWithRanges(tt.contents)
			if (err != nil) != tt.wantErr {
				t.Errorf("GatherContentWithRanges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != tt.wantCount {
					t.Errorf("GatherContentWithRanges() returned %d items, want %d", len(got), tt.wantCount)
				}
				if tt.checkContent != nil {
					tt.checkContent(t, got)
				}
			}
		})
	}
}
