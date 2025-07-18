package nanodoc

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ExtractFileContent reads a file and extracts content based on optional range specification
// The path can include a range suffix like "file.txt:L10-20" or "file.txt:L5"
func ExtractFileContent(pathWithRange string) (*FileContent, error) {
	// Parse the path and range
	path, rangeSpec := parsePathWithRange(pathWithRange)

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &FileError{Path: path, Err: ErrFileNotFound}
		}
		return nil, &FileError{Path: path, Err: err}
	}
	defer func() {
		_ = file.Close()
	}()

	// Read all lines first to determine total line count
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, &FileError{Path: path, Err: err}
	}

	// Parse range if specified
	var r *Range
	if rangeSpec != "" {
		parsedRange, err := parseRange(rangeSpec, len(lines))
		if err != nil {
			return nil, err
		}
		r = parsedRange
	} else {
		// Default to full file
		r = &Range{Start: 1, End: 0}
	}

	// Extract content based on range
	content := extractLinesInRange(lines, r)

	return &FileContent{
		Filepath: path,
		Content:  content,
		Ranges:   []Range{*r},
	}, nil
}

// parsePathWithRange splits a path specification into path and optional range
// Examples: "file.txt" -> ("file.txt", "")
//
//	"file.txt:L10-20" -> ("file.txt", "L10-20")
//	"file.txt:L5" -> ("file.txt", "L5")
func parsePathWithRange(pathWithRange string) (path, rangeSpec string) {
	// Look for the last colon followed by 'L' (to avoid issues with Windows paths)
	idx := strings.LastIndex(pathWithRange, ":L")
	if idx == -1 {
		return pathWithRange, ""
	}

	return pathWithRange[:idx], pathWithRange[idx+1:]
}

// parseRange parses a range specification like "L10-20", "L5", or "L$5-$1" (negative indices)
// Negative indices use $ notation: $1 is last line, $2 is second-to-last, etc.
// A single negative index like "$3" means from 3rd-to-last to end (equivalent to "L$3-$1")
func parseRange(spec string, totalLines int) (*Range, error) {
	if !strings.HasPrefix(spec, "L") {
		return nil, &RangeError{Input: spec, Err: fmt.Errorf("range must start with 'L'")}
	}

	spec = spec[1:] // Remove the 'L' prefix

	// Helper function to parse a line number (positive or negative with $)
	parseLine := func(s string) (int, bool, error) { // returns (line, isNegative, error)
		s = strings.TrimSpace(s)
		if s == "" {
			return 0, false, nil // Empty means EOF
		}
		
		if strings.HasPrefix(s, "$") {
			// Negative index
			if s == "$" {
				return 0, true, fmt.Errorf("$ must be followed by a number")
			}
			negIdx, err := strconv.Atoi(s[1:])
			if err != nil {
				return 0, true, fmt.Errorf("invalid negative index: %w", err)
			}
			if negIdx == 0 {
				return 0, true, fmt.Errorf("$0 is not valid")
			}
			if negIdx < 0 {
				return 0, true, fmt.Errorf("negative index cannot be negative")
			}
			return negIdx, true, nil
		} else {
			// Positive index
			line, err := strconv.Atoi(s)
			if err != nil {
				return 0, false, fmt.Errorf("invalid line number: %w", err)
			}
			return line, false, nil
		}
	}

	// Check if it's a single line or a range
	if strings.Contains(spec, "-") {
		parts := strings.Split(spec, "-")
		if len(parts) != 2 {
			return nil, &RangeError{Input: spec, Err: fmt.Errorf("invalid range format")}
		}

		startNum, startIsNeg, err := parseLine(parts[0])
		if err != nil {
			return nil, &RangeError{Input: spec, Err: fmt.Errorf("invalid start line: %w", err)}
		}

		endNum, endIsNeg, err := parseLine(parts[1])
		if err != nil {
			return nil, &RangeError{Input: spec, Err: fmt.Errorf("invalid end line: %w", err)}
		}

		// Convert negative indices to positive
		start := startNum
		if startIsNeg {
			start = totalLines - startNum + 1
			if start < 1 {
				start = 1
			}
		}

		end := endNum
		if endIsNeg {
			end = totalLines - endNum + 1
			if end < 1 {
				end = 1
			}
		}

		r, err := NewRange(start, end)
		if err != nil {
			return nil, &RangeError{Input: spec, Err: err}
		}
		return &r, nil
	} else {
		// Single line - could be positive or negative
		lineNum, isNeg, err := parseLine(spec)
		if err != nil {
			return nil, &RangeError{Input: spec, Err: err}
		}

		if isNeg {
			// Single negative index like "$3" means from that line to end
			start := totalLines - lineNum + 1
			if start < 1 {
				start = 1
			}
			end := totalLines
			
			r, err := NewRange(start, end)
			if err != nil {
				return nil, &RangeError{Input: spec, Err: err}
			}
			return &r, nil
		} else {
			// Regular single line
			r, err := NewRange(lineNum, lineNum)
			if err != nil {
				return nil, &RangeError{Input: spec, Err: err}
			}
			return &r, nil
		}
	}
}

// extractLinesInRange extracts lines from the slice based on the range
func extractLinesInRange(lines []string, r *Range) string {
	if len(lines) == 0 {
		return ""
	}

	// Adjust range boundaries
	start := r.Start - 1 // Convert to 0-based index
	if start < 0 {
		start = 0
	}
	if start >= len(lines) {
		return ""
	}

	end := r.End
	if end == 0 || end > len(lines) {
		end = len(lines)
	}
	if end < start {
		return ""
	}

	// Extract lines
	extractedLines := lines[start:end]
	return strings.Join(extractedLines, "\n")
}

// ResolveAndExtractFiles takes a list of resolved paths and extracts their content
func ResolveAndExtractFiles(pathInfos []PathInfo, additionalExtensions []string) ([]FileContent, error) {
	var contents []FileContent

	for _, info := range pathInfos {
		switch info.Type {
		case "file":
			// Single file - check if it has range specification in original path
			content, err := ExtractFileContent(info.Original)
			if err != nil {
				return nil, err
			}
			contents = append(contents, *content)

		case "directory", "glob":
			// Multiple files from directory or glob
			for _, filePath := range info.Files {
				content, err := ExtractFileContent(filePath)
				if err != nil {
					return nil, err
				}
				contents = append(contents, *content)
			}

		case "bundle":
			// Bundle files will be handled in step 6
			return nil, fmt.Errorf("bundle files not yet supported")
		}
	}

	return contents, nil
}

// MergeRanges merges overlapping or adjacent ranges
func MergeRanges(ranges []Range) []Range {
	if len(ranges) <= 1 {
		return ranges
	}

	// Sort ranges by start position
	sorted := make([]Range, len(ranges))
	copy(sorted, ranges)
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].Start > sorted[j+1].Start {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	// Merge overlapping ranges
	merged := []Range{sorted[0]}
	for i := 1; i < len(sorted); i++ {
		last := &merged[len(merged)-1]
		current := sorted[i]

		// Check if ranges overlap or are adjacent
		if last.End == 0 || current.Start <= last.End+1 {
			// Merge ranges
			if current.End == 0 || (last.End != 0 && current.End > last.End) {
				last.End = current.End
			}
		} else {
			merged = append(merged, current)
		}
	}

	return merged
}

// GatherContentWithRanges processes multiple file contents and applies range merging
func GatherContentWithRanges(contents []FileContent) ([]FileContent, error) {
	// Group content by file path
	fileMap := make(map[string]*FileContent)

	for _, content := range contents {
		if existing, exists := fileMap[content.Filepath]; exists {
			// Merge ranges for the same file
			existing.Ranges = append(existing.Ranges, content.Ranges...)
		} else {
			// Create a copy to avoid modifying the original
			newContent := content
			fileMap[content.Filepath] = &newContent
		}
	}

	// Process each file to merge ranges and re-extract content
	var result []FileContent
	for _, content := range fileMap {
		// Merge overlapping ranges
		content.Ranges = MergeRanges(content.Ranges)

		// Re-read the file and apply merged ranges
		file, err := os.Open(content.Filepath)
		if err != nil {
			return nil, &FileError{Path: content.Filepath, Err: err}
		}

		// Read all lines
		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		_ = file.Close()

		if err := scanner.Err(); err != nil {
			return nil, &FileError{Path: content.Filepath, Err: err}
		}

		// Apply ranges and gather content
		var contentParts []string
		for _, r := range content.Ranges {
			part := extractLinesInRange(lines, &r)
			if part != "" {
				contentParts = append(contentParts, part)
			}
		}

		content.Content = strings.Join(contentParts, "\n")
		result = append(result, *content)
	}

	return result, nil
}
