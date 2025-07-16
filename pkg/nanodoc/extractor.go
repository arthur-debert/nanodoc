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

// parseRange parses a range specification like "L10-20" or "L5"
func parseRange(spec string, totalLines int) (*Range, error) {
	if !strings.HasPrefix(spec, "L") {
		return nil, &RangeError{Input: spec, Err: fmt.Errorf("range must start with 'L'")}
	}

	spec = spec[1:] // Remove the 'L' prefix

	// Check if it's a single line or a range
	if strings.Contains(spec, "-") {
		parts := strings.Split(spec, "-")
		if len(parts) != 2 {
			return nil, &RangeError{Input: spec, Err: fmt.Errorf("invalid range format")}
		}

		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, &RangeError{Input: spec, Err: fmt.Errorf("invalid start line: %w", err)}
		}

		// Handle open-ended range (e.g., "L10-")
		var end int
		if parts[1] == "" {
			end = 0 // EOF
		} else {
			end, err = strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, &RangeError{Input: spec, Err: fmt.Errorf("invalid end line: %w", err)}
			}
		}

		r, err := NewRange(start, end)
		if err != nil {
			return nil, &RangeError{Input: spec, Err: err}
		}
		return &r, nil
	} else {
		// Single line
		line, err := strconv.Atoi(strings.TrimSpace(spec))
		if err != nil {
			return nil, &RangeError{Input: spec, Err: fmt.Errorf("invalid line number: %w", err)}
		}

		r, err := NewRange(line, line)
		if err != nil {
			return nil, &RangeError{Input: spec, Err: err}
		}
		return &r, nil
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
