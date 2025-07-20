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

func TestGenerateFilename(t *testing.T) {
	doc := &Document{
		TOC: []TOCEntry{
			{Path: "/path/to/myTestFile.txt", Title: "My Test File"},
		},
	}

	tests := []struct {
		name     string
		filepath string
		opts     *FormattingOptions
		seqNum   int
		doc      *Document
		want     string
	}{
		{
			name:     "nice style with sequence",
			filepath: "/path/to/test_file.txt",
			opts: &FormattingOptions{
				HeaderFormat:    HeaderFormatNice,
				SequenceStyle:   SequenceNumerical,
				HeaderAlignment: "left",
				HeaderStyle:     "none",
			},
			seqNum: 1,
			doc:    &Document{},
			want:   "1. Test File",
		},
		{
			name:     "filename style",
			filepath: "/path/to/test_file.txt",
			opts: &FormattingOptions{
				HeaderFormat:    HeaderFormatFilename,
				SequenceStyle:   SequenceNumerical,
				HeaderAlignment: "left",
				HeaderStyle:     "none",
			},
			seqNum: 1,
			doc:    &Document{},
			want:   "1. test_file.txt",
		},
		{
			name:     "path style",
			filepath: "/path/to/test_file.txt",
			opts: &FormattingOptions{
				HeaderFormat:    HeaderFormatPath,
				SequenceStyle:   SequenceNumerical,
				HeaderAlignment: "left",
				HeaderStyle:     "none",
			},
			seqNum: 1,
			doc:    &Document{},
			want:   "1. /path/to/test_file.txt",
		},
		{
			name:     "camelCase file with TOC title",
			filepath: "/path/to/myTestFile.txt",
			opts: &FormattingOptions{
				HeaderFormat:    HeaderFormatNice,
				SequenceStyle:   SequenceRoman,
				HeaderAlignment: "left",
				HeaderStyle:     "none",
			},
			seqNum: 2,
			doc:    doc,
			want:   "ii. My Test File",
		},
		{
			name:     "dashed style",
			filepath: "/path/to/test_file.txt",
			opts: &FormattingOptions{
				HeaderFormat:    HeaderFormatNice,
				SequenceStyle:   SequenceNumerical,
				HeaderAlignment: "left",
				HeaderStyle:     "dashed",
			},
			seqNum: 1,
			doc:    &Document{},
			want:   "------------\n1. Test File\n------------",
		},
		{
			name:     "right align with page width",
			filepath: "/path/to/test_file.txt",
			opts: &FormattingOptions{
				HeaderFormat:    HeaderFormatNice,
				SequenceStyle:   SequenceNumerical,
				HeaderAlignment: "right",
				HeaderStyle:     "none",
				PageWidth:       50,
			},
			seqNum: 1,
			doc:    &Document{},
			want:   "                                      1. Test File",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateFilename(tt.filepath, tt.opts, tt.seqNum, tt.doc)
			if got != tt.want {
				t.Errorf("generateFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBannerStyles(t *testing.T) {
	tests := []struct {
		name            string
		headerStyle     string
		headerAlignment string
		pageWidth       int
		fileName        string
		wantContains    []string
	}{
		{
			name:            "dashed_style",
			headerStyle:     "dashed",
			headerAlignment: "left",
			pageWidth:       80,
			fileName:        "test.txt",
			wantContains:    []string{"-----------", "1. test.txt", "-----------"},
		},
		{
			name:            "solid_style",
			headerStyle:     "solid",
			headerAlignment: "left",
			pageWidth:       80,
			fileName:        "test.txt",
			wantContains:    []string{"===========", "1. test.txt", "==========="},
		},
		{
			name:            "boxed_style_left",
			headerStyle:     "boxed",
			headerAlignment: "left",
			pageWidth:       80,
			fileName:        "test.txt",
			wantContains:    []string{"########", "### 1. test.txt", "###"},
		},
		{
			name:            "boxed_style_center",
			headerStyle:     "boxed",
			headerAlignment: "center",
			pageWidth:       80,
			fileName:        "test.txt",
			wantContains:    []string{"########", "###", "1. test.txt", "###"},
		},
		{
			name:            "boxed_style_right",
			headerStyle:     "boxed",
			headerAlignment: "right",
			pageWidth:       80,
			fileName:        "test.txt",
			wantContains:    []string{"########", "###", "1. test.txt ###"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{
				ContentItems: []FileContent{
					{
						Filepath: "/test/" + tt.fileName,
						Content:  "Test content",
					},
				},
				FormattingOptions: FormattingOptions{
					ShowFilenames:   true,
					HeaderStyle:     tt.headerStyle,
					HeaderAlignment: tt.headerAlignment,
					PageWidth:       tt.pageWidth,
					SequenceStyle:   SequenceNumerical,
					HeaderFormat:    HeaderFormatFilename, // Use filename format for predictable output
				},
			}

			ctx, err := NewFormattingContext(doc.FormattingOptions)
			if err != nil {
				t.Fatalf("NewFormattingContext() error = %v", err)
			}

			output, err := RenderDocument(doc, ctx)
			if err != nil {
				t.Fatalf("RenderDocument() error = %v", err)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Output does not contain %q\nGot:\n%s", want, output)
				}
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
		t.Fatalf("Expected 2 TOC entries, got %d", len(doc.TOC))
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
			name: "basic rendering with filenames",
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
				ShowFilenames:   true,
				HeaderFormat:   HeaderFormatNice,
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
				ShowFilenames: false,
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
				ShowFilenames: true,
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

func TestRenderMarkdownBasic(t *testing.T) {
	tests := []struct {
		name     string
		doc      *Document
		expected string
	}{
		{
			name: "single_markdown_file",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "test.md",
						Content:  "# Header\n\nThis is content.\n",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat: "markdown",
				},
			},
			expected: "# Header\n\nThis is content.\n",
		},
		{
			name: "multiple_markdown_files",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "file1.md",
						Content:  "# First File\n\nContent 1\n",
					},
					{
						Filepath: "file2.md",
						Content:  "# Second File\n\nContent 2\n",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat: "markdown",
				},
			},
			expected: "# First File\n\nContent 1\n# Second File\n\nContent 2\n",
		},
		{
			name: "file_without_trailing_newline",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "test.md",
						Content:  "# Header\n\nNo trailing newline",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat: "markdown",
				},
			},
			expected: "# Header\n\nNo trailing newline\n",
		},
		{
			name: "empty_file",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "empty.md",
						Content:  "",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat: "markdown",
				},
			},
			expected: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renderMarkdownBasic(tt.doc)
			if err != nil {
				t.Errorf("renderMarkdownBasic() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("renderMarkdownBasic() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestRenderPlainText(t *testing.T) {
	tests := []struct {
		name     string
		doc      *Document
		expected string
	}{
		{
			name: "single_text_file",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "test.txt",
						Content:  "Hello World\nThis is plain text.\n",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat: "plain",
				},
			},
			expected: "Hello World\nThis is plain text.\n",
		},
		{
			name: "multiple_text_files",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "file1.txt",
						Content:  "File 1 content\n",
					},
					{
						Filepath: "file2.txt",
						Content:  "File 2 content\n",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat: "plain",
				},
			},
			expected: "File 1 content\nFile 2 content\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renderPlainText(tt.doc)
			if err != nil {
				t.Errorf("renderPlainText() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("renderPlainText() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestRenderDocumentWithOutputFormat(t *testing.T) {
	tests := []struct {
		name         string
		doc          *Document
		ctx          *FormattingContext
		outputFormat string
		checkFunc    func(t *testing.T, result string)
	}{
		{
			name: "markdown_format_uses_context",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "test.md",
						Content:  "# Test\n",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat:  "markdown",
					ShowFilenames: true,
					ShowTOC:       true,
				},
			},
			ctx: &FormattingContext{
				ShowFilenames: true,
				ShowTOC:       true,
				LineNumbers:   LineNumberFile,
			},
			outputFormat: "markdown",
			checkFunc: func(t *testing.T, result string) {
				// Enhanced markdown should now respect context flags
				if !strings.Contains(result, "Table of Contents") {
					t.Error("Markdown output should contain TOC when requested")
				}
				if !strings.Contains(result, "## 1. test.md") {
					t.Error("Markdown output should contain file header when requested")
				}
				// Should not have line numbers (markdown doesn't support this)
				if strings.Contains(result, "1 |") {
					t.Error("Markdown output should not contain line numbers")
				}
			},
		},
		{
			name: "plain_format_ignores_context",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "test.txt",
						Content:  "Plain text content\n",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat:  "plain",
					ShowFilenames: true,
					ShowTOC:       true,
				},
			},
			ctx: &FormattingContext{
				ShowFilenames: true,
				ShowTOC:       true,
				LineNumbers:   LineNumberFile,
			},
			outputFormat: "plain",
			checkFunc: func(t *testing.T, result string) {
				// Should not contain any formatting
				if strings.Contains(result, "Table of Contents") {
					t.Error("Plain output should not contain TOC")
				}
				if strings.Contains(result, "1 |") {
					t.Error("Plain output should not contain line numbers")
				}
				// Should only contain the raw content
				if result != "Plain text content\n" {
					t.Errorf("Expected raw plain content, got %q", result)
				}
			},
		},
		{
			name: "term_format_uses_context",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "/tmp/test.txt",
						Content:  "Terminal content\n",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat:  "term",
					ShowFilenames: true,
					ShowTOC:       false,
					HeaderFormat:  HeaderFormatFilename,
					SequenceStyle: SequenceNumerical,
					HeaderAlignment: "left",
					HeaderStyle:    "none",
				},
			},
			ctx: &FormattingContext{
				ShowFilenames: true,
				ShowTOC:       false,
				LineNumbers:   LineNumberNone,
				HeaderFormat:  HeaderFormatFilename,
			},
			outputFormat: "term",
			checkFunc: func(t *testing.T, result string) {
				// Should contain filename header
				if !strings.Contains(result, "test.txt") {
					t.Errorf("Term output should contain filename, got: %q", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RenderDocument(tt.doc, tt.ctx)
			if err != nil {
				t.Errorf("RenderDocument() error = %v", err)
				return
			}
			tt.checkFunc(t, result)
		})
	}
}

// TestRenderMarkdownEnhanced tests the enhanced markdown rendering with all phases
func TestRenderMarkdownEnhanced(t *testing.T) {
	tests := []struct {
		name      string
		doc       *Document
		ctx       *FormattingContext
		checkFunc func(t *testing.T, result string)
	}{
		{
			name: "phase_2.1_poc_basic_parsing",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "/tmp/test.md",
						Content:  "# Hello World\n\nThis is a test.",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat: "markdown",
				},
			},
			ctx: &FormattingContext{
				ShowFilenames: false,
				ShowTOC:       false,
			},
			checkFunc: func(t *testing.T, result string) {
				if !strings.Contains(result, "# Hello World") {
					t.Errorf("Should preserve markdown header, got: %q", result)
				}
				if !strings.Contains(result, "This is a test.") {
					t.Errorf("Should preserve content, got: %q", result)
				}
			},
		},
		{
			name: "phase_2.2_header_adjustment",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "/tmp/first.md",
						Content:  "# First Doc\n\n## Section",
					},
					{
						Filepath: "/tmp/second.md",
						Content:  "# Second Doc\n\n## Another Section",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat: "markdown",
				},
			},
			ctx: &FormattingContext{
				ShowFilenames: false,
				ShowTOC:       false,
			},
			checkFunc: func(t *testing.T, result string) {
				// First doc should keep H1
				if !strings.Contains(result, "# First Doc") {
					t.Errorf("First doc should keep H1, got: %q", result)
				}
				// Second doc should have H1 adjusted to H2
				if !strings.Contains(result, "## Second Doc") {
					t.Errorf("Second doc should have H1 adjusted to H2, got: %q", result)
				}
				// Second doc's H2 should become H3
				if !strings.Contains(result, "### Another Section") {
					t.Errorf("Second doc's H2 should become H3, got: %q", result)
				}
			},
		},
		{
			name: "phase_2.3_file_headers",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "/tmp/doc1.md",
						Content:  "Content 1",
					},
					{
						Filepath: "/tmp/doc2.md",
						Content:  "Content 2",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat:  "markdown",
					ShowFilenames: true,
					SequenceStyle: SequenceNumerical,
				},
			},
			ctx: &FormattingContext{
				ShowFilenames: true,
			},
			checkFunc: func(t *testing.T, result string) {
				if !strings.Contains(result, "## 1. doc1.md") {
					t.Errorf("Should contain numbered file header for doc1, got: %q", result)
				}
				if !strings.Contains(result, "## 2. doc2.md") {
					t.Errorf("Should contain numbered file header for doc2, got: %q", result)
				}
			},
		},
		{
			name: "phase_2.4_table_of_contents",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "/tmp/guide.md",
						Content:  "# Main Title\n\n## Section 1\n\n### Subsection\n\n## Section 2",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat: "markdown",
					ShowTOC:      true,
				},
			},
			ctx: &FormattingContext{
				ShowTOC: true,
			},
			checkFunc: func(t *testing.T, result string) {
				if !strings.Contains(result, "## Table of Contents") {
					t.Errorf("Should contain TOC header, got: %q", result)
				}
				if !strings.Contains(result, "- [guide.md - Main Title]") {
					t.Errorf("Should contain TOC entry for main title, got: %q", result)
				}
				if !strings.Contains(result, "  - [guide.md - Section 1]") {
					t.Errorf("Should contain nested TOC entry, got: %q", result)
				}
			},
		},
		{
			name: "all_phases_combined",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "/tmp/readme.md",
						Content:  "# Project\n\n## Overview\n\nIntroduction.",
					},
					{
						Filepath: "/tmp/guide.md", 
						Content:  "# Guide\n\n## Getting Started\n\nInstructions.",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat:  "markdown",
					ShowFilenames: true,
					ShowTOC:       true,
					SequenceStyle: SequenceLetter,
				},
			},
			ctx: &FormattingContext{
				ShowFilenames: true,
				ShowTOC:       true,
			},
			checkFunc: func(t *testing.T, result string) {
				// Check TOC
				if !strings.Contains(result, "## Table of Contents") {
					t.Errorf("Should contain TOC, got: %q", result)
				}
				// Check file headers
				if !strings.Contains(result, "## a. readme.md") {
					t.Errorf("Should contain letter sequence file header, got: %q", result)
				}
				if !strings.Contains(result, "## b. guide.md") {
					t.Errorf("Should contain second file header, got: %q", result)
				}
				// Check header adjustment
				if !strings.Contains(result, "## Guide") {
					t.Errorf("Second file H1 should be adjusted to H2, got: %q", result)
				}
			},
		},
		{
			name: "non_markdown_files_passthrough",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "/tmp/code.go",
						Content:  "package main\n\nfunc main() {}",
					},
				},
				FormattingOptions: FormattingOptions{
					OutputFormat: "markdown",
				},
			},
			ctx: &FormattingContext{},
			checkFunc: func(t *testing.T, result string) {
				if !strings.Contains(result, "package main") {
					t.Errorf("Non-markdown content should pass through, got: %q", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renderMarkdownEnhanced(tt.doc, tt.ctx)
			if err != nil {
				t.Fatalf("renderMarkdownEnhanced() error = %v", err)
			}
			tt.checkFunc(t, result)
		})
	}
}
