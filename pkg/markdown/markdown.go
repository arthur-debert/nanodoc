// Package markdown provides markdown parsing, transformation, and rendering capabilities
// for nanodoc. It supports header level adjustment, TOC generation, and file header
// insertion to enable rich markdown output formatting.
package markdown

import (
	"bytes"
	"fmt"
	"strings"

	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Document represents a parsed markdown document with its AST and source
type Document struct {
	AST    ast.Node
	Source []byte
}

// Parser wraps goldmark parser functionality
type Parser struct {
	gm goldmark.Markdown
}

// NewParser creates a new markdown parser
func NewParser() *Parser {
	return &Parser{
		gm: goldmark.New(
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
		),
	}
}

// Parse converts markdown content into a Document
func (p *Parser) Parse(content []byte) (*Document, error) {
	reader := text.NewReader(content)
	node := p.gm.Parser().Parse(reader)
	
	return &Document{
		AST:    node,
		Source: content,
	}, nil
}

// Transformer provides methods to modify markdown AST
type Transformer struct{}

// NewTransformer creates a new markdown transformer
func NewTransformer() *Transformer {
	return &Transformer{}
}

// AdjustHeaderLevels increases all header levels by the specified amount
func (t *Transformer) AdjustHeaderLevels(doc *Document, increment int) error {
	if increment == 0 {
		return nil
	}

	return ast.Walk(doc.AST, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if heading, ok := n.(*ast.Heading); ok {
				newLevel := heading.Level + increment
				if newLevel > 6 {
					newLevel = 6 // Max header level in markdown
				}
				heading.Level = newLevel
			}
		}
		return ast.WalkContinue, nil
	})
}

// HasH1 checks if the document contains any H1 headers
func (t *Transformer) HasH1(doc *Document) bool {
	hasH1 := false
	_ = ast.Walk(doc.AST, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if heading, ok := n.(*ast.Heading); ok && heading.Level == 1 {
				hasH1 = true
				return ast.WalkStop, nil
			}
		}
		return ast.WalkContinue, nil
	})
	return hasH1
}

// InsertFileHeader adds a header at the beginning of the document
func (t *Transformer) InsertFileHeader(doc *Document, headerText string, level int) error {
	// Create new header node
	header := ast.NewHeading(level)
	header.AppendChild(header, ast.NewString([]byte(headerText)))

	// Insert at the beginning
	if doc.AST.FirstChild() != nil {
		doc.AST.InsertBefore(doc.AST, doc.AST.FirstChild(), header)
	} else {
		doc.AST.AppendChild(doc.AST, header)
	}

	return nil
}

// Renderer converts markdown AST back to markdown text
type Renderer struct {
	gm goldmark.Markdown
}

// NewRenderer creates a new markdown renderer
func NewRenderer() *Renderer {
	return &Renderer{
		gm: goldmark.New(
			goldmark.WithRenderer(
				markdown.NewRenderer(),
			),
		),
	}
}

// Render converts a Document back to markdown text
func (r *Renderer) Render(doc *Document) ([]byte, error) {
	var buf bytes.Buffer
	err := r.gm.Renderer().Render(&buf, doc.Source, doc.AST)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// TOCGenerator extracts and generates table of contents
type TOCGenerator struct{}

// NewTOCGenerator creates a new TOC generator
func NewTOCGenerator() *TOCGenerator {
	return &TOCGenerator{}
}

// TOCEntry represents a single entry in the table of contents
type TOCEntry struct {
	Level int
	Text  string
	ID    string
}

// ExtractTOC extracts all headers from the document for TOC generation
func (tg *TOCGenerator) ExtractTOC(doc *Document) []TOCEntry {
	var entries []TOCEntry
	
	_ = ast.Walk(doc.AST, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if heading, ok := n.(*ast.Heading); ok {
				text := tg.extractHeadingText(heading, doc.Source)
				id := tg.generateAnchorID(text)
				
				entries = append(entries, TOCEntry{
					Level: heading.Level,
					Text:  text,
					ID:    id,
				})
			}
		}
		return ast.WalkContinue, nil
	})
	
	return entries
}

// GenerateTOCMarkdown creates a markdown formatted table of contents
func (tg *TOCGenerator) GenerateTOCMarkdown(entries []TOCEntry) string {
	if len(entries) == 0 {
		return ""
	}
	
	var builder strings.Builder
	builder.WriteString("## Table of Contents\n\n")
	
	for _, entry := range entries {
		// Create indentation based on header level
		indent := strings.Repeat("  ", entry.Level-1)
		builder.WriteString(fmt.Sprintf("%s- [%s](#%s)\n", indent, entry.Text, entry.ID))
	}
	
	return builder.String()
}

// extractHeadingText extracts the text content from a heading node
func (tg *TOCGenerator) extractHeadingText(heading *ast.Heading, source []byte) string {
	var text strings.Builder
	
	_ = ast.Walk(heading, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch node := n.(type) {
			case *ast.Text:
				text.Write(node.Segment.Value(source))
			case *ast.String:
				text.Write(node.Value)
			}
		}
		return ast.WalkContinue, nil
	})
	
	return text.String()
}

// generateAnchorID creates a GitHub-style anchor ID from heading text
func (tg *TOCGenerator) generateAnchorID(text string) string {
	// Convert to lowercase
	id := strings.ToLower(text)
	
	// Replace spaces with hyphens
	id = strings.ReplaceAll(id, " ", "-")
	
	// Remove special characters except hyphens and underscores
	var cleaned strings.Builder
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			cleaned.WriteRune(r)
		}
	}
	
	return cleaned.String()
}

// HeaderFormatter formats file headers according to nanodoc options
type HeaderFormatter struct{}

// NewHeaderFormatter creates a new header formatter
func NewHeaderFormatter() *HeaderFormatter {
	return &HeaderFormatter{}
}

// FormatFileHeader creates a formatted header for a file
func (hf *HeaderFormatter) FormatFileHeader(filename string, sequence string, level int) string {
	header := filename
	
	if sequence != "" {
		header = fmt.Sprintf("%s. %s", sequence, header)
	}
	
	return header
}

// Helper function to normalize multiple consecutive newlines
func normalizeNewlines(content string) string {
	// Replace multiple newlines with double newline
	for strings.Contains(content, "\n\n\n") {
		content = strings.ReplaceAll(content, "\n\n\n", "\n\n")
	}
	return strings.TrimSpace(content) + "\n"
}