package markdown

import (
	"strings"
	"testing"
)

// Test Parser functionality
func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "simple markdown",
			content: "# Hello\n\nThis is a paragraph.",
			wantErr: false,
		},
		{
			name:    "empty content",
			content: "",
			wantErr: false,
		},
		{
			name:    "complex markdown",
			content: "# Title\n\n## Subtitle\n\n- Item 1\n- Item 2\n\n```go\ncode\n```",
			wantErr: false,
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.Parse([]byte(tt.content))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && doc == nil {
				t.Error("Parse() returned nil document")
			}
		})
	}
}

// Test Transformer header level adjustment
func TestTransformer_AdjustHeaderLevels(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		increment int
		want      string
	}{
		{
			name:      "increase by 1",
			content:   "# H1\n\n## H2\n\n### H3",
			increment: 1,
			want:      "## H1\n\n### H2\n\n#### H3",
		},
		{
			name:      "increase by 2",
			content:   "# H1\n\n## H2",
			increment: 2,
			want:      "### H1\n\n#### H2",
		},
		{
			name:      "no change",
			content:   "# H1\n\n## H2",
			increment: 0,
			want:      "# H1\n\n## H2",
		},
		{
			name:      "max level cap",
			content:   "##### H5\n\n###### H6",
			increment: 2,
			want:      "###### H5\n\n###### H6",
		},
	}

	parser := NewParser()
	transformer := NewTransformer()
	renderer := NewRenderer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.Parse([]byte(tt.content))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			err = transformer.AdjustHeaderLevels(doc, tt.increment)
			if err != nil {
				t.Fatalf("AdjustHeaderLevels() error = %v", err)
			}

			result, err := renderer.Render(doc)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			got := normalizeNewlines(string(result))
			want := normalizeNewlines(tt.want)
			if got != want {
				t.Errorf("AdjustHeaderLevels() got = %q, want %q", got, want)
			}
		})
	}
}

// Test H1 detection
func TestTransformer_HasH1(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "has H1",
			content: "# Title\n\n## Subtitle",
			want:    true,
		},
		{
			name:    "no H1",
			content: "## Subtitle\n\n### Sub-subtitle",
			want:    false,
		},
		{
			name:    "multiple H1",
			content: "# Title 1\n\n# Title 2",
			want:    true,
		},
		{
			name:    "empty document",
			content: "",
			want:    false,
		},
	}

	parser := NewParser()
	transformer := NewTransformer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.Parse([]byte(tt.content))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			got := transformer.HasH1(doc)
			if got != tt.want {
				t.Errorf("HasH1() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test file header insertion
func TestTransformer_InsertFileHeader(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		headerText string
		level      int
		want       string
	}{
		{
			name:       "insert H2 header",
			content:    "Some content",
			headerText: "file.md",
			level:      2,
			want:       "## file.md\n\nSome content",
		},
		{
			name:       "insert into empty doc",
			content:    "",
			headerText: "empty.md",
			level:      1,
			want:       "# empty.md",
		},
		{
			name:       "insert with sequence",
			content:    "# Existing\n\nContent",
			headerText: "1. test.md",
			level:      2,
			want:       "## 1. test.md\n\n# Existing\n\nContent",
		},
	}

	parser := NewParser()
	transformer := NewTransformer()
	renderer := NewRenderer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.Parse([]byte(tt.content))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			err = transformer.InsertFileHeader(doc, tt.headerText, tt.level)
			if err != nil {
				t.Fatalf("InsertFileHeader() error = %v", err)
			}

			result, err := renderer.Render(doc)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			got := normalizeNewlines(string(result))
			want := normalizeNewlines(tt.want)
			if got != want {
				t.Errorf("InsertFileHeader() got = %q, want %q", got, want)
			}
		})
	}
}

// Test TOC extraction
func TestTOCGenerator_ExtractTOC(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []TOCEntry
	}{
		{
			name: "multi-level headers",
			content: `# Main Title

## Section 1

### Subsection 1.1

## Section 2`,
			want: []TOCEntry{
				{Level: 1, Text: "Main Title", ID: "main-title"},
				{Level: 2, Text: "Section 1", ID: "section-1"},
				{Level: 3, Text: "Subsection 1.1", ID: "subsection-11"},
				{Level: 2, Text: "Section 2", ID: "section-2"},
			},
		},
		{
			name:    "no headers",
			content: "Just some text without headers",
			want:    []TOCEntry{},
		},
		{
			name: "special characters",
			content: `# Hello, World!

## What's New?

### Code & Testing`,
			want: []TOCEntry{
				{Level: 1, Text: "Hello, World!", ID: "hello-world"},
				{Level: 2, Text: "What's New?", ID: "whats-new"},
				{Level: 3, Text: "Code & Testing", ID: "code--testing"},
			},
		},
	}

	parser := NewParser()
	tocGen := NewTOCGenerator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.Parse([]byte(tt.content))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			got := tocGen.ExtractTOC(doc)
			if len(got) != len(tt.want) {
				t.Errorf("ExtractTOC() got %d entries, want %d", len(got), len(tt.want))
				return
			}

			for i, entry := range got {
				if entry.Level != tt.want[i].Level ||
					entry.Text != tt.want[i].Text ||
					entry.ID != tt.want[i].ID {
					t.Errorf("ExtractTOC()[%d] = {%d, %q, %q}, want {%d, %q, %q}",
						i, entry.Level, entry.Text, entry.ID,
						tt.want[i].Level, tt.want[i].Text, tt.want[i].ID)
				}
			}
		})
	}
}

// Test TOC markdown generation
func TestTOCGenerator_GenerateTOCMarkdown(t *testing.T) {
	tests := []struct {
		name    string
		entries []TOCEntry
		want    string
	}{
		{
			name: "nested TOC",
			entries: []TOCEntry{
				{Level: 1, Text: "Title", ID: "title"},
				{Level: 2, Text: "Section 1", ID: "section-1"},
				{Level: 3, Text: "Subsection", ID: "subsection"},
				{Level: 2, Text: "Section 2", ID: "section-2"},
			},
			want: `## Table of Contents

- [Title](#title)
  - [Section 1](#section-1)
    - [Subsection](#subsection)
  - [Section 2](#section-2)`,
		},
		{
			name:    "empty TOC",
			entries: []TOCEntry{},
			want:    "",
		},
		{
			name: "single entry",
			entries: []TOCEntry{
				{Level: 1, Text: "Only Title", ID: "only-title"},
			},
			want: `## Table of Contents

- [Only Title](#only-title)`,
		},
	}

	tocGen := NewTOCGenerator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tocGen.GenerateTOCMarkdown(tt.entries)
			got = strings.TrimSpace(got)
			want := strings.TrimSpace(tt.want)
			if got != want {
				t.Errorf("GenerateTOCMarkdown() got = %q, want %q", got, want)
			}
		})
	}
}

// Test header formatting
func TestHeaderFormatter_FormatFileHeader(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		sequence string
		level    int
		want     string
	}{
		{
			name:     "simple filename",
			filename: "test.md",
			sequence: "",
			level:    2,
			want:     "test.md",
		},
		{
			name:     "with numeric sequence",
			filename: "main.go",
			sequence: "1",
			level:    2,
			want:     "1. main.go",
		},
		{
			name:     "with letter sequence",
			filename: "readme.md",
			sequence: "a",
			level:    3,
			want:     "a. readme.md",
		},
		{
			name:     "with roman sequence",
			filename: "config.yaml",
			sequence: "iv",
			level:    2,
			want:     "iv. config.yaml",
		},
	}

	formatter := NewHeaderFormatter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatter.FormatFileHeader(tt.filename, tt.sequence, tt.level)
			if got != tt.want {
				t.Errorf("FormatFileHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test complex rendering scenarios
func TestRenderer_ComplexMarkdown(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name: "full document",
			content: `# Title

This is a paragraph with **bold** and *italic* text.

## Lists

### Unordered
- Item 1
- Item 2
  - Nested item
- Item 3

### Ordered
1. First
2. Second
3. Third

## Code

Inline ` + "`code`" + ` example.

` + "```go" + `
func main() {
    fmt.Println("Hello")
}
` + "```" + `

## Links

Check out [this link](https://example.com "Example").

> This is a blockquote
> with multiple lines`,
		},
	}

	parser := NewParser()
	renderer := NewRenderer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse -> Render -> Parse again to ensure consistency
			doc1, err := parser.Parse([]byte(tt.content))
			if err != nil {
				t.Fatalf("First parse error = %v", err)
			}

			rendered, err := renderer.Render(doc1)
			if err != nil {
				t.Fatalf("Render error = %v", err)
			}

			// The rendered output should be valid markdown
			doc2, err := parser.Parse(rendered)
			if err != nil {
				t.Fatalf("Second parse error = %v", err)
			}

			// Basic check: both should have content
			if doc2 == nil || doc2.AST == nil {
				t.Error("Rendered markdown could not be parsed")
			}
		})
	}
}

// Integration test for full workflow
func TestIntegration_FullWorkflow(t *testing.T) {
	// Simulate processing multiple markdown files as per phases 2.1-2.4
	
	// File 1: Has H1, should have headers adjusted
	file1 := `# Project Title

## Introduction

This is the intro.`

	// File 2: No H1, headers stay the same
	file2 := `## Features

- Feature 1
- Feature 2`

	parser := NewParser()
	transformer := NewTransformer()
	renderer := NewRenderer()
	tocGen := NewTOCGenerator()
	headerFormatter := NewHeaderFormatter()

	// Process file 1
	doc1, _ := parser.Parse([]byte(file1))
	if transformer.HasH1(doc1) {
		_ = transformer.AdjustHeaderLevels(doc1, 1)
	}
	_ = transformer.InsertFileHeader(doc1, headerFormatter.FormatFileHeader("readme.md", "1", 2), 2)

	// Process file 2
	doc2, _ := parser.Parse([]byte(file2))
	if transformer.HasH1(doc2) {
		_ = transformer.AdjustHeaderLevels(doc2, 1)
	}
	_ = transformer.InsertFileHeader(doc2, headerFormatter.FormatFileHeader("features.md", "2", 2), 2)

	// Render both
	rendered1, err := renderer.Render(doc1)
	if err != nil {
		t.Fatalf("Failed to render doc1: %v", err)
	}
	rendered2, err := renderer.Render(doc2)
	if err != nil {
		t.Fatalf("Failed to render doc2: %v", err)
	}

	// Debug output
	t.Logf("Rendered file 1:\n%s", string(rendered1))
	t.Logf("Rendered file 2:\n%s", string(rendered2))

	// Combine and generate TOC
	combined := string(rendered1) + "\n" + string(rendered2)
	combinedDoc, _ := parser.Parse([]byte(combined))
	tocEntries := tocGen.ExtractTOC(combinedDoc)
	toc := tocGen.GenerateTOCMarkdown(tocEntries)

	// Verify we have a complete document with TOC
	if toc == "" {
		t.Error("No TOC generated")
	}
	if !strings.Contains(string(rendered1), "## 1. readme.md") {
		t.Error("File header not added to file 1")
	}
	if !strings.Contains(string(rendered1), "## Project Title") {
		t.Error("H1 not adjusted in file 1")
	}
}