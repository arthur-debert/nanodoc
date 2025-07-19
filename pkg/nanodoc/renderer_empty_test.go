package nanodoc

import (
	"strings"
	"testing"
)

func TestRenderEmptyFiles(t *testing.T) {
	tests := []struct {
		name        string
		doc         *Document
		ctx         *FormattingContext
		wantContain string
	}{
		{
			name: "empty file without line numbers",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "empty.txt",
						Content:  "",
					},
				},
				FormattingOptions: FormattingOptions{
					ShowFilenames:   true,
					HeaderFormat:   HeaderFormatNice,
					SequenceStyle: SequenceNumerical,
				},
			},
			ctx: &FormattingContext{
				ShowFilenames:   true,
				HeaderFormat:   HeaderFormatNice,
				SequenceStyle: SequenceNumerical,
				LineNumbers:   LineNumberNone,
			},
			wantContain: "(empty file)",
		},
		{
			name: "empty file with file line numbers",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "empty.txt",
						Content:  "",
					},
				},
				FormattingOptions: FormattingOptions{
					ShowFilenames:   true,
					HeaderFormat:   HeaderFormatNice,
					SequenceStyle: SequenceNumerical,
				},
			},
			ctx: &FormattingContext{
				ShowFilenames:   true,
				HeaderFormat:   HeaderFormatNice,
				SequenceStyle: SequenceNumerical,
				LineNumbers:   LineNumberFile,
			},
			wantContain: "1 | (empty file)",
		},
		{
			name: "empty file with global line numbers",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "first.txt",
						Content:  "Some content\nLine 2",
					},
					{
						Filepath: "empty.txt",
						Content:  "",
					},
				},
				FormattingOptions: FormattingOptions{
					ShowFilenames:   true,
					HeaderFormat:   HeaderFormatNice,
					SequenceStyle: SequenceNumerical,
				},
			},
			ctx: &FormattingContext{
				ShowFilenames:   true,
				HeaderFormat:   HeaderFormatNice,
				SequenceStyle: SequenceNumerical,
				LineNumbers:   LineNumberGlobal,
			},
			wantContain: "3 | (empty file)",
		},
		{
			name: "multiple empty files",
			doc: &Document{
				ContentItems: []FileContent{
					{
						Filepath: "empty1.txt",
						Content:  "",
					},
					{
						Filepath: "empty2.txt",
						Content:  "",
					},
				},
				FormattingOptions: FormattingOptions{
					ShowFilenames:   true,
					HeaderFormat:   HeaderFormatNice,
					SequenceStyle: SequenceNumerical,
				},
			},
			ctx: &FormattingContext{
				ShowFilenames:   true,
				HeaderFormat:   HeaderFormatNice,
				SequenceStyle: SequenceNumerical,
				LineNumbers:   LineNumberNone,
			},
			wantContain: "(empty file)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderDocument(tt.doc, tt.ctx)
			if err != nil {
				t.Fatalf("RenderDocument() error = %v", err)
			}

			if !strings.Contains(got, tt.wantContain) {
				t.Errorf("RenderDocument() output doesn't contain %q\nGot:\n%s", tt.wantContain, got)
			}

			// Count occurrences for multiple empty files test
			if tt.name == "multiple empty files" {
				count := strings.Count(got, "(empty file)")
				if count != 2 {
					t.Errorf("Expected 2 occurrences of '(empty file)', got %d", count)
				}
			}
		})
	}
}

func TestEmptyFileIntegration(t *testing.T) {
	// Test that empty files are handled correctly through the full pipeline
	doc := &Document{
		ContentItems: []FileContent{
			{
				Filepath: "normal.txt",
				Content:  "This file has content",
			},
			{
				Filepath: "empty.txt",
				Content:  "",
			},
			{
				Filepath: "another.txt",
				Content:  "More content here",
			},
		},
		FormattingOptions: FormattingOptions{
			ShowFilenames:   true,
			HeaderFormat:   HeaderFormatNice,
			SequenceStyle: SequenceNumerical,
		},
	}

	ctx := &FormattingContext{
		ShowFilenames:   true,
		HeaderFormat:   HeaderFormatNice,
		SequenceStyle: SequenceNumerical,
		LineNumbers:   LineNumberNone,
	}

	got, err := RenderDocument(doc, ctx)
	if err != nil {
		t.Fatalf("RenderDocument() error = %v", err)
	}

	// Check structure
	expectedParts := []string{
		"1. Normal",
		"This file has content",
		"2. Empty",
		"(empty file)",
		"3. Another",
		"More content here",
	}

	for _, part := range expectedParts {
		if !strings.Contains(got, part) {
			t.Errorf("Output missing expected part: %q\nGot:\n%s", part, got)
		}
	}
}

func TestEmptyDocument(t *testing.T) {
	// Test rendering a document with no content items
	doc := &Document{
		ContentItems: []FileContent{},
		FormattingOptions: FormattingOptions{
			ShowFilenames:   true,
			HeaderFormat:   HeaderFormatNice,
			SequenceStyle: SequenceNumerical,
		},
	}

	ctx := &FormattingContext{
		ShowFilenames:   true,
		HeaderFormat:   HeaderFormatNice,
		SequenceStyle: SequenceNumerical,
		LineNumbers:   LineNumberNone,
	}

	got, err := RenderDocument(doc, ctx)
	if err != nil {
		t.Fatalf("RenderDocument() error = %v", err)
	}

	if got != "" {
		t.Errorf("Expected empty string for empty document, got %q", got)
	}
}

func TestDocumentWithOnlyEmptyFiles(t *testing.T) {
	// Test rendering a document with only empty files
	doc := &Document{
		ContentItems: []FileContent{
			{
				Filepath: "empty1.txt",
				Content:  "",
			},
			{
				Filepath: "empty2.txt",
				Content:  "",
			},
		},
		FormattingOptions: FormattingOptions{
			ShowFilenames:   true,
			HeaderFormat:   HeaderFormatNice,
			SequenceStyle: SequenceNumerical,
		},
	}

	ctx := &FormattingContext{
		ShowFilenames:   true,
		HeaderFormat:   HeaderFormatNice,
		SequenceStyle: SequenceNumerical,
		LineNumbers:   LineNumberNone,
	}

	got, err := RenderDocument(doc, ctx)
	if err != nil {
		t.Fatalf("RenderDocument() error = %v", err)
	}

	// Check structure
	expectedParts := []string{
		"1. Empty1",
		"(empty file)",
		"2. Empty2",
		"(empty file)",
	}

	for _, part := range expectedParts {
		if !strings.Contains(got, part) {
			t.Errorf("Output missing expected part: %q\nGot:\n%s", part, got)
		}
	}
}

func TestEmptyFileWithTOC(t *testing.T) {
	// Test that empty files don't break TOC generation
	doc := &Document{
		ContentItems: []FileContent{
			{
				Filepath: "empty.md",
				Content:  "",
			},
			{
				Filepath: "content.md",
				Content:  "# Title\nSome content",
			},
		},
		FormattingOptions: FormattingOptions{
			ShowFilenames:   true,
			HeaderFormat:   HeaderFormatNice,
			SequenceStyle: SequenceNumerical,
			ShowTOC:       true,
		},
	}

	ctx := &FormattingContext{
		ShowFilenames:   true,
		HeaderFormat:   HeaderFormatNice,
		SequenceStyle: SequenceNumerical,
		LineNumbers:   LineNumberNone,
		ShowTOC:       true,
	}

	got, err := RenderDocument(doc, ctx)
	if err != nil {
		t.Fatalf("RenderDocument() error = %v", err)
	}

	// Check that TOC is generated correctly
	if !strings.Contains(got, "Table of Contents") {
		t.Error("Output missing TOC")
	}
	if !strings.Contains(got, "- Title (content.md)") {
		t.Error("Output missing TOC entry for content.md")
	}
	if strings.Contains(got, "empty.md") {
		t.Error("Output should not contain TOC entry for empty.md")
	}
}
