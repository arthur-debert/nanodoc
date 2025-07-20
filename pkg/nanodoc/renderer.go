package nanodoc

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/arthur-debert/nanodoc/pkg/markdown"
)

// RendererOptions controls how the document is rendered
type RendererOptions struct {
	IncludeTOC         bool
	IncludeLineNumbers bool
	UseRichFormatting  bool
}

// RenderDocument renders a Document object to a string
func RenderDocument(doc *Document, ctx *FormattingContext) (string, error) {
	// For markdown output, use enhanced renderer with all features
	if doc.FormattingOptions.OutputFormat == "markdown" {
		return renderMarkdownEnhanced(doc, ctx)
	}

	// For plain output, concatenate without any formatting
	if doc.FormattingOptions.OutputFormat == "plain" {
		return renderPlainText(doc)
	}

	var parts []string

	// Generate TOC first, as it's used for filenames
	if ctx.ShowTOC || ctx.HeaderFormat == HeaderFormatNice {
		slog.Debug("Generating table of contents for filenames/TOC")
		generateTOC(doc)
	}

	// Render TOC if requested
	if ctx.ShowTOC {
		var tocParts []string
		tocParts = append(tocParts, "Table of Contents")
		tocParts = append(tocParts, "=================")
		tocParts = append(tocParts, "")
		for _, entry := range doc.TOC {
			// Indent based on heading level, assuming Level 1 is the base
			indent := strings.Repeat("  ", entry.Level-1)
			tocParts = append(tocParts, fmt.Sprintf("%s- %s (%s)", indent, entry.Title, filepath.Base(entry.Path)))
		}
		tocParts = append(tocParts, "")
		parts = append(parts, strings.Join(tocParts, "\n"))
		parts = append(parts, "\n")
	}

	// Render each content item
	prevOriginalSource := ""
	sequenceNumber := 0
	globalLineNumber := 1

	for _, item := range doc.ContentItems {
		// Check if we need a file separator
		isNotInlined := item.OriginalSource == ""
		differentSource := item.Filepath != prevOriginalSource

		if isNotInlined && differentSource && ctx.ShowFilenames {
			// Add separator if not first item
			if len(parts) > 0 && !strings.HasSuffix(parts[len(parts)-1], "\n\n") {
				parts = append(parts, "\n")
			}

			// Generate filename
			sequenceNumber++
			filename := generateFilename(item.Filepath, &doc.FormattingOptions, sequenceNumber, doc)
			parts = append(parts, filename)
			parts = append(parts, "\n\n")
		}

		// Add content with optional line numbers
		content := item.Content
		
		// Handle empty files
		if content == "" {
			content = "(empty file)"
		}
		
		if ctx.LineNumbers != LineNumberNone {
			numberedContent, newGlobalLineNum := addLineNumbers(content, ctx.LineNumbers, globalLineNumber)
			content = numberedContent
			if ctx.LineNumbers == LineNumberGlobal {
				globalLineNumber = newGlobalLineNum
			}
		}

		parts = append(parts, content)

		// Ensure content ends with newline
		if len(parts) > 0 && !strings.HasSuffix(parts[len(parts)-1], "\n") {
			parts = append(parts, "\n")
		}

		// Track source for next iteration
		if item.OriginalSource != "" {
			prevOriginalSource = item.OriginalSource
		} else {
			prevOriginalSource = item.Filepath
		}
	}

	result := strings.Join(parts, "")
	return result, nil
}

func generateFilename(filePath string, opts *FormattingOptions, seqNum int, doc *Document) string {
	headerText := generateFileHeaderText(filePath, opts, seqNum, doc)

	// Get banner style from registry
	style, exists := GetBannerStyle(opts.HeaderStyle)
	if !exists {
		// Fallback to none style if not found
		style, _ = GetBannerStyle("none")
	}

	// Apply the banner style
	return style.Apply(headerText, opts)
}

// generateFileHeaderText generates the text content for a file header
func generateFileHeaderText(filePath string, opts *FormattingOptions, seqNum int, doc *Document) string {
	// Find the primary title for this file from the TOC
	var title string
	for _, entry := range doc.TOC {
		if entry.Path == filePath {
			title = entry.Title
			break
		}
	}

	var baseName string
	switch opts.HeaderFormat {
	case HeaderFormatFilename:
		baseName = filepath.Base(filePath)
	case HeaderFormatPath:
		baseName = filePath
	case HeaderFormatNice:
		fallthrough
	default:
		// Use title from TOC if available, otherwise generate from filename
		niceName := title
		if niceName == "" {
			filename := filepath.Base(filePath)
			nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
			niceName = strings.ReplaceAll(nameWithoutExt, "_", " ")
			niceName = strings.ReplaceAll(niceName, "-", " ")
			niceName = splitCamelCase(niceName)
			niceName = toTitleCase(niceName)
		}
		baseName = niceName
	}

	// Add sequence number
	seq := generateSequence(seqNum, opts.SequenceStyle)
	if seq != "" {
		return fmt.Sprintf("%s. %s", seq, baseName)
	}
	return baseName
}

// generateSequence generates a sequence number in the specified style
func generateSequence(num int, style SequenceStyle) string {
	switch style {
	case SequenceNumerical:
		return strconv.Itoa(num)
	case SequenceLetter:
		if num <= 26 {
			return string(rune('a' + num - 1))
		}
		// For numbers > 26, use aa, ab, ac, etc.
		return string(rune('a'+((num-1)/26)-1)) + string(rune('a'+((num-1)%26)))
	case SequenceRoman:
		return toRoman(num)
	default:
		return strconv.Itoa(num)
	}
}

// splitCamelCase splits a camelCase string into words
func splitCamelCase(s string) string {
	// Add space before capital letters preceded by lowercase
	re1 := regexp.MustCompile("([a-z])([A-Z])")
	s = re1.ReplaceAllString(s, "$1 $2")
	
	// Handle consecutive uppercase followed by lowercase (e.g., HTMLFile -> HTML File)
	re2 := regexp.MustCompile("([A-Z])([A-Z][a-z])")
	s = re2.ReplaceAllString(s, "$1 $2")
	
	return s
}

// toTitleCase converts a string to title case (first letter of each word uppercase)
func toTitleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

// toRoman converts a number to Roman numerals (simplified version)
func toRoman(num int) string {
	values := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	symbols := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	
	result := ""
	for i := 0; i < len(values); i++ {
		for num >= values[i] {
			num -= values[i]
			result += symbols[i]
		}
	}
	return strings.ToLower(result)
}

// addLineNumbers adds line numbers to content
func addLineNumbers(content string, mode LineNumberMode, startNum int) (string, int) {
	lines := strings.Split(content, "\n")
	
	// Calculate the width needed for line numbers
	maxLineNum := startNum + len(lines) - 1
	if mode == LineNumberFile {
		maxLineNum = len(lines)
	}
	width := len(strconv.Itoa(maxLineNum))
	
	var result []string
	lineNum := startNum
	if mode == LineNumberFile {
		lineNum = 1
	}
	
	for _, line := range lines {
		// Don't add line numbers to empty lines at the end
		if line == "" && lineNum == len(lines) {
			result = append(result, line)
		} else {
			numberedLine := fmt.Sprintf("%*d | %s", width, lineNum, line)
			result = append(result, numberedLine)
		}
		lineNum++
	}
	
	return strings.Join(result, "\n"), lineNum
}

// generateTOC generates a table of contents for the document using the markdown parser.
func generateTOC(doc *Document) {
	doc.TOC = make([]TOCEntry, 0)
	parser := markdown.NewParser()
	tocGen := markdown.NewTOCGenerator()

	var allHeadings []TOCEntry
	sequenceNum := 1

	for _, item := range doc.ContentItems {
		// Only extract headings from markdown files
		if !strings.HasSuffix(item.Filepath, ".md") && !strings.HasSuffix(item.Filepath, ".markdown") {
			continue
		}

		mdDoc, err := parser.Parse([]byte(item.Content))
		if err != nil {
			slog.Warn("failed to parse markdown for TOC generation", "file", item.Filepath, "error", err)
			continue
		}

		entries := tocGen.ExtractTOC(mdDoc)
		for _, entry := range entries {
			allHeadings = append(allHeadings, TOCEntry{
				Title:    entry.Text,
				Level:    entry.Level,
				Path:     item.Filepath,
				Sequence: generateSequence(sequenceNum, doc.FormattingOptions.SequenceStyle),
				// LineNumber is not available from the new parser, which is acceptable.
			})
			sequenceNum++
		}
	}
	doc.TOC = allHeadings
}



// renderMarkdownBasic performs basic concatenation of markdown files without any modifications
// This is kept for backward compatibility and fallback
func renderMarkdownBasic(doc *Document) (string, error) {
	var parts []string

	for _, item := range doc.ContentItems {
		// Simply append the content as-is
		parts = append(parts, item.Content)
		
		// Ensure content ends with newline
		if len(parts) > 0 && !strings.HasSuffix(parts[len(parts)-1], "\n") {
			parts = append(parts, "\n")
		}
	}

	result := strings.Join(parts, "")
	return result, nil
}

// renderMarkdownEnhanced uses the markdown package to provide rich markdown output
func renderMarkdownEnhanced(doc *Document, ctx *FormattingContext) (string, error) {
	// Phase 2.1: POC - Demonstrate all capabilities
	parser := markdown.NewParser()
	transformer := markdown.NewTransformer()
	renderer := markdown.NewRenderer()
	tocGen := markdown.NewTOCGenerator()
	headerFormatter := markdown.NewHeaderFormatter()

	var processedDocs []*markdown.Document

	// Generate TOC first if needed, so it's available for all renderers
	if ctx.ShowTOC {
		slog.Debug("Generating table of contents for markdown output")
		generateTOC(doc)
	}

	// Process each content item
	for i, item := range doc.ContentItems {
		mdDoc, err := parser.Parse([]byte(item.Content))
		if err != nil {
			return "", fmt.Errorf("failed to parse content for file %s: %w", item.Filepath, err)
		}

		isMarkdown := strings.HasSuffix(item.Filepath, ".md") || strings.HasSuffix(item.Filepath, ".markdown")

		if isMarkdown {
			// Perform markdown-specific transformations

			// Adjust header levels for subsequent documents to maintain hierarchy
			if i > 0 && transformer.HasH1(mdDoc) {
				if err := transformer.AdjustHeaderLevels(mdDoc, 1); err != nil {
					return "", fmt.Errorf("failed to adjust header levels for %s: %w", item.Filepath, err)
				}
			}

			// Insert file headers if requested
			if ctx.ShowFilenames {
				sequenceNum := i + 1
				headerText := generateFileHeaderText(item.Filepath, &doc.FormattingOptions, sequenceNum, doc)

				// Format as a markdown header. H2 is chosen as a sensible default
				// to avoid conflicting with a potential H1 title in the first document.
				const headerLevel = 2
				mdHeaderText := headerFormatter.FormatFileHeader(headerText, "", headerLevel)

				if err := transformer.InsertFileHeader(mdDoc, mdHeaderText, headerLevel); err != nil {
					return "", fmt.Errorf("failed to insert file header for %s: %w", item.Filepath, err)
				}
			}
		}

		processedDocs = append(processedDocs, mdDoc)
	}

	// Build final output
	var output strings.Builder

	// Phase 2.4: Add TOC if requested
	if ctx.ShowTOC && len(doc.TOC) > 0 {
		// Convert nanodoc.TOCEntry to markdown.TOCEntry
		mdTOCEntries := make([]markdown.TOCEntry, len(doc.TOC))
		for i, entry := range doc.TOC {
			mdTOCEntries[i] = markdown.TOCEntry{
				Text:  fmt.Sprintf("%s - %s", filepath.Base(entry.Path), entry.Title),
				Level: entry.Level,
			}
		}
		tocMarkdown := tocGen.GenerateTOCMarkdown(mdTOCEntries)
		output.WriteString(tocMarkdown)
		output.WriteString("\n\n")
	}

	// Render all processed documents
	for i, mdDoc := range processedDocs {
		if i > 0 {
			output.WriteString("\n")
		}

		rendered, err := renderer.Render(mdDoc)
		if err != nil {
			return "", fmt.Errorf("failed to render markdown: %w", err)
		}

		output.Write(rendered)
	}

	return output.String(), nil
}

// renderPlainText performs basic concatenation without any formatting
func renderPlainText(doc *Document) (string, error) {
	var parts []string

	for _, item := range doc.ContentItems {
		// Simply append the content as-is
		parts = append(parts, item.Content)
		
		// Ensure content ends with newline
		if len(parts) > 0 && !strings.HasSuffix(parts[len(parts)-1], "\n") {
			parts = append(parts, "\n")
		}
	}

	result := strings.Join(parts, "")
	return result, nil
}