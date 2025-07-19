package nanodoc

import (
	"errors"
	"testing"
)

func TestNewRange(t *testing.T) {
	tests := []struct {
		name    string
		start   int
		end     int
		want    Range
		wantErr bool
	}{
		{
			name:  "valid range",
			start: 10,
			end:   20,
			want:  Range{Start: 10, End: 20},
		},
		{
			name:  "full file range",
			start: 1,
			end:   0,
			want:  Range{Start: 1, End: 0},
		},
		{
			name:  "single line",
			start: 5,
			end:   5,
			want:  Range{Start: 5, End: 5},
		},
		{
			name:    "negative start",
			start:   -1,
			end:     10,
			wantErr: true,
		},
		{
			name:    "zero start",
			start:   0,
			end:     10,
			wantErr: true,
		},
		{
			name:    "end before start",
			start:   10,
			end:     5,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRange(tt.start, tt.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("NewRange() = %v, want %v", got, tt.want)
			}
			if tt.wantErr && err != nil {
				var rangeErr *RangeError
				if !errors.As(err, &rangeErr) {
					t.Errorf("NewRange() error type = %T, want *RangeError", err)
				}
			}
		})
	}
}

func TestRange_Contains(t *testing.T) {
	tests := []struct {
		name string
		r    Range
		line int
		want bool
	}{
		{
			name: "line in range",
			r:    Range{Start: 10, End: 20},
			line: 15,
			want: true,
		},
		{
			name: "line at start",
			r:    Range{Start: 10, End: 20},
			line: 10,
			want: true,
		},
		{
			name: "line at end",
			r:    Range{Start: 10, End: 20},
			line: 20,
			want: true,
		},
		{
			name: "line before range",
			r:    Range{Start: 10, End: 20},
			line: 5,
			want: false,
		},
		{
			name: "line after range",
			r:    Range{Start: 10, End: 20},
			line: 25,
			want: false,
		},
		{
			name: "EOF range contains high line",
			r:    Range{Start: 10, End: 0},
			line: 1000,
			want: true,
		},
		{
			name: "EOF range doesn't contain line before start",
			r:    Range{Start: 10, End: 0},
			line: 5,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Contains(tt.line); got != tt.want {
				t.Errorf("Range.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRange_IsFullFile(t *testing.T) {
	tests := []struct {
		name string
		r    Range
		want bool
	}{
		{
			name: "full file range",
			r:    Range{Start: 1, End: 0},
			want: true,
		},
		{
			name: "not full file - different start",
			r:    Range{Start: 2, End: 0},
			want: false,
		},
		{
			name: "not full file - has end",
			r:    Range{Start: 1, End: 100},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.IsFullFile(); got != tt.want {
				t.Errorf("Range.IsFullFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDocument(t *testing.T) {
	doc := NewDocument()

	if doc == nil {
		t.Fatal("NewDocument() returned nil")
	}

	if doc.ContentItems == nil {
		t.Error("NewDocument() ContentItems is nil")
	}

	if len(doc.ContentItems) != 0 {
		t.Errorf("NewDocument() ContentItems length = %d, want 0", len(doc.ContentItems))
	}

	if doc.TOC == nil {
		t.Error("NewDocument() TOC is nil")
	}

	if len(doc.TOC) != 0 {
		t.Errorf("NewDocument() TOC length = %d, want 0", len(doc.TOC))
	}

	if !doc.FormattingOptions.ShowFilenames {
		t.Error("NewDocument() ShowFilenames = false, want true")
	}

	if doc.FormattingOptions.FilenameStyle != FilenameStyleNice {
		t.Errorf("NewDocument() FilenameStyle = %s, want %s", doc.FormattingOptions.FilenameStyle, FilenameStyleNice)
	}
}

func TestFileContent(t *testing.T) {
	// Test FileContent structure initialization
	fc := FileContent{
		Filepath: "/path/to/file.txt",
		Ranges: []Range{
			{Start: 1, End: 10},
			{Start: 20, End: 30},
		},
		Content:        "test content",
		IsBundle:       false,
		OriginalSource: "",
	}

	if fc.Filepath != "/path/to/file.txt" {
		t.Errorf("FileContent.Filepath = %s, want /path/to/file.txt", fc.Filepath)
	}

	if len(fc.Ranges) != 2 {
		t.Errorf("FileContent.Ranges length = %d, want 2", len(fc.Ranges))
	}

	if fc.Content != "test content" {
		t.Errorf("FileContent.Content = %s, want 'test content'", fc.Content)
	}
}

func TestFormattingOptions(t *testing.T) {
	opts := FormattingOptions{
		Theme:                ThemeClassic,
		LineNumbers:          LineNumberFile,
		ShowFilenames:          true,
		FilenameStyle:          FilenameStylePath,
		SequenceStyle:        SequenceNumerical,
		ShowTOC:              false,
		AdditionalExtensions: []string{".go", ".py"},
	}

	if opts.LineNumbers != LineNumberFile {
		t.Errorf("FormattingOptions.LineNumbers = %v, want %v", opts.LineNumbers, LineNumberFile)
	}

	if !opts.ShowFilenames {
		t.Error("FormattingOptions.ShowFilenames = false, want true")
	}

	if opts.SequenceStyle != SequenceNumerical {
		t.Errorf("FormattingOptions.SequenceStyle = %s, want %s", opts.SequenceStyle, SequenceNumerical)
	}

	if opts.FilenameStyle != FilenameStylePath {
		t.Errorf("FormattingOptions.FilenameStyle = %s, want %s", opts.FilenameStyle, FilenameStylePath)
	}

	if len(opts.AdditionalExtensions) != 2 {
		t.Errorf("FormattingOptions.AdditionalExtensions length = %d, want 2", len(opts.AdditionalExtensions))
	}
}

func TestTOCEntry(t *testing.T) {
	entry := TOCEntry{
		Title:      "Chapter 1",
		Path:       "/path/to/chapter1.txt",
		Sequence:   "1",
		LineNumber: 42,
	}

	if entry.Title != "Chapter 1" {
		t.Errorf("TOCEntry.Title = %s, want 'Chapter 1'", entry.Title)
	}

	if entry.Path != "/path/to/chapter1.txt" {
		t.Errorf("TOCEntry.Path = %s, want '/path/to/chapter1.txt'", entry.Path)
	}

	if entry.Sequence != "1" {
		t.Errorf("TOCEntry.Sequence = %s, want '1'", entry.Sequence)
	}

	if entry.LineNumber != 42 {
		t.Errorf("TOCEntry.LineNumber = %d, want 42", entry.LineNumber)
	}
}
