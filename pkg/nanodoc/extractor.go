package nanodoc

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ExtractFileContent reads a file and extracts content based on optional range specifications.
// The path can include a range suffix like "file.txt:L10-20,L30,L40-".
func ExtractFileContent(pathWithRange string) (*FileContent, error) {
	path, rangeSpec := parsePathWithRange(pathWithRange)

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

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, &FileError{Path: path, Err: err}
	}

	var ranges []Range
	if rangeSpec != "" {
		parsedRanges, err := parseRanges(rangeSpec, len(lines))
		if err != nil {
			return nil, err
		}
		ranges = parsedRanges
	} else {
		// Default to the full file range.
		ranges = []Range{{Start: 1, End: len(lines)}}
	}

	var contentParts []string
	for _, r := range ranges {
		contentPart := extractLinesInRange(lines, &r)
		contentParts = append(contentParts, contentPart)
	}
	content := strings.Join(contentParts, "\n")

	return &FileContent{
		Filepath: path,
		Content:  content,
		Ranges:   ranges,
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

// parseRanges parses a comma-separated list of range specifications.
func parseRanges(spec string, totalLines int) ([]Range, error) {
	rangeStrings := strings.Split(spec, ",")
	var ranges []Range

	for _, rangeStr := range rangeStrings {
		if !strings.HasPrefix(rangeStr, "L") {
			return nil, &RangeError{Input: spec, Err: fmt.Errorf("range specifier must start with 'L'")}
		}
		parsedRange, err := parseSingleRange(rangeStr, totalLines)
		if err != nil {
			return nil, err // Propagate error with original spec
		}
		ranges = append(ranges, *parsedRange)
	}

	return ranges, nil
}

// parseSingleRange parses a single range specification like "L10-20" or "L$5-$1".
func parseSingleRange(spec string, totalLines int) (*Range, error) {
	if !strings.HasPrefix(spec, "L") {
		return nil, &RangeError{Input: spec, Err: fmt.Errorf("range must start with 'L'")}
	}

	spec = spec[1:] // Remove 'L'

	parseLine := func(s string) (int, bool, error) {
		s = strings.TrimSpace(s)
		if s == "" {
			return 0, false, nil // Empty means EOF
		}
		if strings.HasPrefix(s, "$") {
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
			line, err := strconv.Atoi(s)
			if err != nil {
				return 0, false, fmt.Errorf("invalid line number: %w", err)
			}
			return line, false, nil
		}
	}

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
		} else if endNum == 0 && parts[1] == "" { // Handle open-ended range like L10-
			end = totalLines
		}

		if end < 1 {
			end = 1
		}

		r, err := NewRange(start, end)
		if err != nil {
			return nil, &RangeError{Input: spec, Err: err}
		}
		return &r, nil
	} else {
		lineNum, isNeg, err := parseLine(spec)
		if err != nil {
			return nil, &RangeError{Input: spec, Err: err}
		}

		if isNeg {
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


