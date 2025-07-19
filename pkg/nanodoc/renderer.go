package nanodoc

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// RendererOptions controls how the document is rendered
type RendererOptions struct {
	IncludeTOC         bool
	IncludeLineNumbers bool
	UseRichFormatting  bool
}

// RenderDocument renders a Document object to a string
func RenderDocument(doc *Document, ctx *FormattingContext) (string, error) {
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
			tocParts = append(tocParts, fmt.Sprintf("- %s (%s)", entry.Title, filepath.Base(entry.Path)))
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
		baseName = fmt.Sprintf("%s. %s", seq, baseName)
	}

	// Apply banner style
	switch opts.HeaderStyle {
	case "dashed":
		line := strings.Repeat("-", len(baseName))
		return fmt.Sprintf("%s\n%s\n%s", line, baseName, line)
	case "solid":
		line := strings.Repeat("=", len(baseName))
		return fmt.Sprintf("%s\n%s\n%s", line, baseName, line)
	}

	// Apply alignment
	switch opts.HeaderAlignment {
	case "center":
		return fmt.Sprintf("%*s", len(baseName)+10, baseName)
	case "right":
		return fmt.Sprintf("%*s", 80, baseName)
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

// generateTOC generates a table of contents for the document
func generateTOC(doc *Document) {
	headingsByFile := extractHeadings(doc)
	if len(headingsByFile) == 0 {
		return
	}

	// Update document TOC entries
	doc.TOC = make([]TOCEntry, 0)
	
	// Sort file paths for consistent order
	var sortedPaths []string
	for path := range headingsByFile {
		sortedPaths = append(sortedPaths, path)
	}
	sort.Strings(sortedPaths)

	sequenceNum := 1
	for _, filePath := range sortedPaths {
		headings := headingsByFile[filePath]
		// Sort headings by line number
		sort.Slice(headings, func(i, j int) bool {
			return headings[i].LineNum < headings[j].LineNum
		})

		for _, heading := range headings {
			doc.TOC = append(doc.TOC, TOCEntry{
				Title:      heading.Text,
				Path:       filePath,
				Sequence:   generateSequence(sequenceNum, doc.FormattingOptions.SequenceStyle),
				LineNumber: heading.LineNum,
			})
			sequenceNum++
		}
	}
}

// HeadingInfo represents a heading found in the document
type HeadingInfo struct {
	Text     string
	Level    int
	LineNum  int
}

// extractHeadings extracts headings from document content
func extractHeadings(doc *Document) map[string][]HeadingInfo {
	headingByFile := make(map[string][]HeadingInfo)
	
	// Markdown heading regex
	headingPattern := regexp.MustCompile(`^(#+)\s+(.+)$`)
	
	for _, item := range doc.ContentItems {
		var fileHeadings []HeadingInfo
		
		// Use original source if available
		filePath := item.OriginalSource
		if filePath == "" {
			filePath = item.Filepath
		}
		
		// Only extract headings from markdown files
		if !strings.HasSuffix(filePath, ".md") && !strings.HasSuffix(filePath, ".markdown") {
			continue
		}

		lines := strings.Split(item.Content, "\n")
		hasMarkdownHeadings := false
		
		for i, line := range lines {
			if matches := headingPattern.FindStringSubmatch(line); matches != nil {
				level := len(matches[1])
				text := strings.TrimSpace(matches[2])
				
				// Only include level 1 and 2 headings
				if level <= 2 {
					hasMarkdownHeadings = true
					fileHeadings = append(fileHeadings, HeadingInfo{
						Text:	text,
						Level:	level,
						LineNum: i + 1,
					})
				}
			}
		}
		
		// If no markdown headings, use first non-empty line as title for markdown files
		if !hasMarkdownHeadings {
			for i, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					// Limit to 50 characters
					if len(line) > 50 {
						line = line[:50] + "..."
					}
					fileHeadings = append(fileHeadings, HeadingInfo{
						Text:	line,
						Level:	1,
						LineNum: i + 1,
					})
					break
				}
			}
		}
		
		// Store headings if any were found
		if len(fileHeadings) > 0 {
			if existing, ok := headingByFile[filePath]; ok {
				headingByFile[filePath] = append(existing, fileHeadings...)
			} else {
				headingByFile[filePath] = fileHeadings
			}
		}
	}
	
	return headingByFile
}