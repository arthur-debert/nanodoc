package nanodoc

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// setupTestFile creates a temporary file with the given content.
func setupTestFile(t *testing.T, name string, lines int) (string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "nanodoc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	filePath := filepath.Join(tempDir, name)
	var content []string
	for i := 1; i <= lines; i++ {
		content = append(content, "line "+strconv.Itoa(i))
	}

	if err := os.WriteFile(filePath, []byte(strings.Join(content, "\n")), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(tempDir)
	}
	return filePath, cleanup
}

func TestParseRanges(t *testing.T) {
	tests := []struct {
		name       string
		spec       string
		totalLines int
		want       []Range
		wantErr    bool
	}{
		{"single line", "L5", 10, []Range{{5, 5}}, false},
		{"single range", "L2-4", 10, []Range{{2, 4}}, false},
		{"multiple ranges", "L2-3,L5-6", 10, []Range{{2, 3}, {5, 6}}, false},
		{"mixed single and range", "L1,L3-4,L6", 10, []Range{{1, 1}, {3, 4}, {6, 6}}, false},
		{"unordered declaration", "L10,L1-2", 10, []Range{{10, 10}, {1, 2}}, false},
		{"overlapping ranges", "L1-5,L3-7", 10, []Range{{1, 5}, {3, 7}}, false},
		{"negative indices", "L$3-$1,L1", 10, []Range{{8, 10}, {1, 1}}, false},
		{"open-ended range", "L8-", 10, []Range{{8, 10}}, false},
		{"invalid spec", "L1,L-", 10, nil, true},
		{"no L prefix", "1-2", 10, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRanges(tt.spec, tt.totalLines)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRanges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !equalRanges(got, tt.want) {
				t.Errorf("parseRanges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractFileContent_MultiRange(t *testing.T) {
	filePath, cleanup := setupTestFile(t, "test.txt", 10)
	defer cleanup()

	tests := []struct {
		name        string
		pathWithRange string
		wantContent string
		wantErr     bool
	}{
		{
			"multiple ranges",
			filePath + ":L1,L3-4,L6",
			"line 1\nline 3\nline 4\nline 6",
			false,
		},
		{
			"unordered declaration",
			filePath + ":L5,L1-2",
			"line 5\nline 1\nline 2",
			false,
		},
		{
			"overlapping ranges",
			filePath + ":L1-3,L2-4",
			"line 1\nline 2\nline 3\nline 2\nline 3\nline 4",
			false,
		},
		{
			"full file if no range",
			filePath,
			"line 1\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7\nline 8\nline 9\nline 10",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc, err := ExtractFileContent(tt.pathWithRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractFileContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && fc.Content != tt.wantContent {
				t.Errorf("ExtractFileContent() content = %q, want %q", fc.Content, tt.wantContent)
			}
		})
	}
}

// equalRanges is a helper to compare two slices of Range.
func equalRanges(a, b []Range) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
