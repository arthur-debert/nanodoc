package nanodoc

import (
	"strings"
	"testing"
)

func TestGenerateSequence(t *testing.T) {
	tests := []struct {
		name  string
		num   int
		style SequenceStyle
		want  string
	}{
		{
			name:  "numerical 1",
			num:   1,
			style: SequenceNumerical,
			want:  "1",
		},
		{
			name:  "numerical 10",
			num:   10,
			style: SequenceNumerical,
			want:  "10",
		},
		{
			name:  "letter a",
			num:   1,
			style: SequenceLetter,
			want:  "a",
		},
		{
			name:  "letter z",
			num:   26,
			style: SequenceLetter,
			want:  "z",
		},
		{
			name:  "letter aa",
			num:   27,
			style: SequenceLetter,
			want:  "aa",
		},
		{
			name:  "roman i",
			num:   1,
			style: SequenceRoman,
			want:  "i",
		},
		{
			name:  "roman v",
			num:   5,
			style: SequenceRoman,
			want:  "v",
		},
		{
			name:  "roman ix",
			num:   9,
			style: SequenceRoman,
			want:  "ix",
		},
		{
			name:  "roman xiv",
			num:   14,
			style: SequenceRoman,
			want:  "xiv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateSequence(tt.num, tt.style)
			if got != tt.want {
				t.Errorf("generateSequence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple camelCase",
			input: "wordNice",
			want:  "word Nice",
		},
		{
			name:  "PascalCase",
			input: "WordNice",
			want:  "Word Nice",
		},
		{
			name:  "with acronym",
			input: "myHTMLFile",
			want:  "my HTML File",
		},
		{
			name:  "all lowercase",
			input: "lowercase",
			want:  "lowercase",
		},
		{
			name:  "all uppercase",
			input: "UPPERCASE",
			want:  "UPPERCASE",
		},
		{
			name:  "mixed complex",
			input: "XMLHttpRequest",
			want:  "XML Http Request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitCamelCase(tt.input)
			if got != tt.want {
				t.Errorf("splitCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddLineNumbers(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		mode     LineNumberMode
		startNum int
		wantHas  []string
	}{
		{
			name:     "file mode",
			content:  "line 1\nline 2\nline 3",
			mode:     LineNumberFile,
			startNum: 10, // Should be ignored in file mode
			wantHas:  []string{"1 | line 1", "2 | line 2", "3 | line 3"},
		},
		{
			name:     "global mode",
			content:  "line 1\nline 2",
			mode:     LineNumberGlobal,
			startNum: 10,
			wantHas:  []string{"10 | line 1", "11 | line 2"},
		},
		{
			name:     "proper padding",
			content:  strings.Repeat("line\n", 10),
			mode:     LineNumberFile,
			startNum: 1,
			wantHas:  []string{" 1 | line", "10 | line"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := addLineNumbers(tt.content, tt.mode, tt.startNum)
			for _, want := range tt.wantHas {
				if !strings.Contains(got, want) {
					t.Errorf("addLineNumbers() result doesn't contain %q", want)
				}
			}
		})
	}
}

func TestGenerateFileHeader(t *testing.T) {
	doc := &Document{
		TOC: []TOCEntry{
			{Path: "/path/to/myTestFile.txt", Title: "My Test File"},
		},
	}

	tests := []struct {
		name     string
		filepath string
		style    HeaderStyle
		seqStyle SequenceStyle
		seqNum   int
		doc      *Document
		want     string
	}{
		{
			name:     "nice style with sequence",
			filepath: "/path/to/test_file.txt",
			style:    HeaderStyleNice,
			seqStyle: SequenceNumerical,
			seqNum:   1,
			doc:      &Document{},
			want:     "1. Test File",
		},
		{
			name:     "filename style",
			filepath: "/path/to/test_file.txt",
			style:    HeaderStyleFilename,
			seqStyle: SequenceNumerical,
			seqNum:   1,
			doc:      &Document{},
			want:     "1. test_file.txt",
		},
		{
			name:     "path style",
			filepath: "/path/to/test_file.txt",
			style:    HeaderStylePath,
			seqStyle: SequenceNumerical,
			seqNum:   1,
			doc:      &Document{},
			want:     "1. /path/to/test_file.txt",
		},
		{
			name:     "camelCase file with TOC title",
			filepath: "/path/to/myTestFile.txt",
			style:    HeaderStyleNice,
			seqStyle: SequenceRoman,
			seqNum:   2,
			doc:      doc,
			want:     "ii. My Test File",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateFileHeader(tt.filepath, tt.style, tt.seqStyle, tt.seqNum, tt.doc)
			if got != tt.want {
				t.Errorf("generateFileHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractHeadings(t *testing.T) {
	doc := &Document{
		ContentItems: []FileContent{
			{
				Filepath: "/test/file1.md",
				Content: `# Main Title
Some content here

## Subsection
More content`,
			},
			{
				Filepath: "/test/file2.txt",
				Content: `This is a plain text file`,
			},
		},
		FormattingOptions: FormattingOptions{
			SequenceStyle: SequenceNumerical,
		},
	}

	generateTOC(doc)

	if len(doc.TOC) != 2 {
		t.Errorf("Expected 2 TOC entries, got %d", len(doc.TOC))
	}

	if doc.TOC[0].Title != "Main Title" {
		t.Errorf("Expected title 'Main Title', got %q", doc.TOC[0].Title)
	}
	if doc.TOC[1].Title != "Subsection" {
		t.Errorf("Expected title 'Subsection', got %q", doc.TOC[1].Title)
	}
}

func TestRenderDocument(t *testing.T) {
	tests := []struct {
		name string
		doc  *Document
		ctx  *FormattingContext
		want []string // List of strings that should appear in output
	}{
		{
			name: "basic rendering with headers",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "/path/to/file1.txt",
						Content:  "Content of file 1",
					},
					{
						Filepath: "/path/to/file2.md",
						Content:  "# Title of File 2",
					},
				},
			},
			ctx: &FormattingContext{
				ShowHeaders:   true,
				HeaderStyle:   HeaderStyleNice,
				SequenceStyle: SequenceNumerical,
				LineNumbers:   LineNumberNone,
			},
			want: []string{
				"1. File1",
				"Content of file 1",
				"2. Title of File 2",
			},
		},
		{
			name: "with line numbers",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "/path/to/test.txt",
						Content:  "Line 1\nLine 2\nLine 3",
					},
				},
			},
			ctx: &FormattingContext{
				ShowHeaders: false,
				LineNumbers: LineNumberFile,
			},
			want: []string{
				"1 | Line 1",
				"2 | Line 2",
				"3 | Line 3",
			},
		},
		{
			name: "with TOC",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "/path/to/doc.md",
						Content:  "# Title\n\n## Section 1\n\nContent",
					},
				},
				FormattingOptions: FormattingOptions{
					SequenceStyle: SequenceNumerical,
				},
			},
			ctx: &FormattingContext{
				ShowHeaders: true,
				ShowTOC:     true,
			},
			want: []string{
				"Table of Contents",
				"doc.md",
				"Title",
				"Section 1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderDocument(tt.doc, tt.ctx)
			if err != nil {
				t.Fatalf("RenderDocument() error = %v", err)
			}

			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("RenderDocument() output doesn't contain %q", want)
				}
			}
		})
	}
}

func TestGenerateTOC(t *testing.T) {
	doc := &Document{
		ContentItems: []FileContent{
			{
				Filepath: "/test/doc1.md",
				Content:  "# Main Title\n\n## Subsection\n\nContent",
			},
			{
				Filepath: "/test/doc2.txt",
				Content:  "Plain text content\nwith multiple lines",
			},
		},
		FormattingOptions: FormattingOptions{
			SequenceStyle: SequenceNumerical,
		},
	}

	generateTOC(doc)

	if len(doc.TOC) != 2 {
		t.Fatalf("Expected 2 TOC entries, got %d", len(doc.TOC))
	}

	expectedTitles := []string{
		"Main Title",
		"Subsection",
	}

	for i, title := range expectedTitles {
		if doc.TOC[i].Title != title {
			t.Errorf("Expected TOC title %q, got %q", title, doc.TOC[i].Title)
		}
	}
}