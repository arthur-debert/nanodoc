package nanodoc

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DryRunInfo contains information about what would be processed
type DryRunInfo struct {
	// Files that would be processed
	Files []FileInfo
	// Bundle files detected
	Bundles []string
	// Total count of files
	TotalFiles int
	// Total line count across all files
	TotalLines int
	// Files requiring additional extensions
	RequiresExtension map[string]string
	// Active formatting options
	Options FormattingOptions
}

// FileInfo contains dry run information about a file
type FileInfo struct {
	Path      string
	Source    string // Where it came from (directory, bundle, etc.)
	Extension string
	LineCount int    // Number of lines that will be processed
	RangeSpec string // Range specification if any (e.g., "L10-20")
}

// GenerateDryRunInfo analyzes what files would be processed without actually processing them
func GenerateDryRunInfo(pathInfos []PathInfo, opts FormattingOptions) (*DryRunInfo, error) {
	info := &DryRunInfo{
		Files:             make([]FileInfo, 0),
		Bundles:           make([]string, 0),
		RequiresExtension: make(map[string]string),
		Options:           opts,
	}

	// Process each path
	for _, pathInfo := range pathInfos {
		switch pathInfo.Type {
		case "file":
			ext := filepath.Ext(pathInfo.Absolute)
			// Extract range spec from original path
			_, rangeSpec := parsePathWithRange(pathInfo.Original)
			
			fileInfo := FileInfo{
				Path:      pathInfo.Absolute,
				Source:    "direct argument",
				Extension: ext,
				RangeSpec: rangeSpec,
			}
			
			// Count lines in the file
			lineCount, err := countFileLines(pathInfo.Original)
			if err != nil {
				return nil, err
			}
			fileInfo.LineCount = lineCount
			info.TotalLines += lineCount
			
			// Check if file needs additional extension
			if !isTextFileWithExtensions(pathInfo.Absolute, opts.AdditionalExtensions) {
				info.RequiresExtension[pathInfo.Absolute] = ext
			}
			
			info.Files = append(info.Files, fileInfo)
			
		case "directory":
			for _, file := range pathInfo.Files {
				fileInfo := FileInfo{
					Path:      file,
					Source:    fmt.Sprintf("directory: %s", pathInfo.Original),
					Extension: filepath.Ext(file),
				}
				
				// Count lines in the file
				lineCount, err := countFileLines(file)
				if err != nil {
					return nil, err
				}
				fileInfo.LineCount = lineCount
				info.TotalLines += lineCount
				
				info.Files = append(info.Files, fileInfo)
			}
			
		case "glob":
			for _, file := range pathInfo.Files {
				fileInfo := FileInfo{
					Path:      file,
					Source:    fmt.Sprintf("glob: %s", pathInfo.Original),
					Extension: filepath.Ext(file),
				}
				
				// Count lines in the file
				lineCount, err := countFileLines(file)
				if err != nil {
					return nil, err
				}
				fileInfo.LineCount = lineCount
				info.TotalLines += lineCount
				
				info.Files = append(info.Files, fileInfo)
			}
			
		case "bundle":
			info.Bundles = append(info.Bundles, pathInfo.Absolute)
			// For dry run, we need to process bundle contents to count lines
			bp := NewBundleProcessor()
			bundlePaths, err := bp.ProcessBundleFile(pathInfo.Absolute)
			if err != nil {
				return nil, err
			}
			
			// Count lines in each file referenced by the bundle
			for _, bundlePath := range bundlePaths {
				fileInfo := FileInfo{
					Path:      bundlePath,
					Source:    fmt.Sprintf("bundle: %s", filepath.Base(pathInfo.Absolute)),
					Extension: filepath.Ext(bundlePath),
				}
				
				// Extract range spec from bundle path
				_, rangeSpec := parsePathWithRange(bundlePath)
				fileInfo.RangeSpec = rangeSpec
				
				// Count lines in the file
				lineCount, err := countFileLines(bundlePath)
				if err != nil {
					// Skip files that can't be read
					continue
				}
				fileInfo.LineCount = lineCount
				info.TotalLines += lineCount
				
				info.Files = append(info.Files, fileInfo)
			}
		}
	}
	
	// Remove duplicates and sort
	info.Bundles = uniqueStrings(info.Bundles)
	sort.Strings(info.Bundles)
	
	// Count total files
	info.TotalFiles = len(info.Files)
	
	// Add lines for filenames if enabled
	if opts.ShowFilenames {
		// Each file gets a filename line
		info.TotalLines += info.TotalFiles
	}
	
	// Add lines for TOC if enabled
	if opts.ShowTOC {
		// Estimate TOC lines: title + separator + one line per file
		info.TotalLines += 2 + info.TotalFiles
	}
	
	return info, nil
}

// FormatDryRunOutput formats the dry run information for display
func FormatDryRunOutput(info *DryRunInfo) string {
	var output strings.Builder
	
	output.WriteString("Would process the following files:\n")
	
	// Show TOC line count if enabled
	if info.Options.ShowTOC {
		tocLines := 2 + info.TotalFiles // title + separator + entries
		output.WriteString(fmt.Sprintf("\nTable of Contents (%d lines)\n", tocLines))
	}
	
	// Group files by source
	filesBySource := make(map[string][]FileInfo)
	for _, file := range info.Files {
		filesBySource[file.Source] = append(filesBySource[file.Source], file)
	}
	
	// Sort sources for consistent output
	sources := make([]string, 0, len(filesBySource))
	for source := range filesBySource {
		sources = append(sources, source)
	}
	sort.Strings(sources)
	
	// Display files grouped by source
	fileNum := 1
	for _, source := range sources {
		files := filesBySource[source]
		output.WriteString(fmt.Sprintf("\nFrom %s:\n", source))
		
		// Sort files within each source
		sort.Slice(files, func(i, j int) bool {
			return files[i].Path < files[j].Path
		})
		
		for _, file := range files {
			relPath := filepath.Base(file.Path)
			if file.RangeSpec != "" {
				relPath = fmt.Sprintf("%s:%s", relPath, file.RangeSpec)
			}
			output.WriteString(fmt.Sprintf("%d. %s (%d lines)\n", fileNum, relPath, file.LineCount))
			fileNum++
		}
	}
	
	// Show bundle information
	if len(info.Bundles) > 0 {
		output.WriteString("\nBundle files detected:\n")
		for _, bundle := range info.Bundles {
			output.WriteString(fmt.Sprintf("  - %s\n", filepath.Base(bundle)))
		}
	}
	
	// Show files requiring extensions
	if len(info.RequiresExtension) > 0 {
		output.WriteString("\nFiles requiring --ext flag:\n")
		for file, ext := range info.RequiresExtension {
			output.WriteString(fmt.Sprintf("  - %s (requires --ext=%s)\n", 
				filepath.Base(file), strings.TrimPrefix(ext, ".")))
		}
	}
	
	// Summary
	output.WriteString(fmt.Sprintf("\nTotal files to process: %d (%d lines)\n", info.TotalFiles, info.TotalLines))
	
	// Show active options
	var activeOptions []string
	if info.Options.ShowTOC {
		activeOptions = append(activeOptions, "--toc")
	}
	if info.Options.LineNumbers != LineNumberNone {
		lineNumMode := "file"
		if info.Options.LineNumbers == LineNumberGlobal {
			lineNumMode = "global"
		}
		activeOptions = append(activeOptions, fmt.Sprintf("--linenum %s", lineNumMode))
	}
	if info.Options.Theme != "classic" {
		activeOptions = append(activeOptions, fmt.Sprintf("--theme %s", info.Options.Theme))
	}
	if !info.Options.ShowFilenames {
		activeOptions = append(activeOptions, "--filenames=false")
	}
	if info.Options.SequenceStyle != "numerical" {
		activeOptions = append(activeOptions, fmt.Sprintf("--file-numbering %s", info.Options.SequenceStyle))
	}
	if info.Options.FilenameStyle != "nice" {
		activeOptions = append(activeOptions, fmt.Sprintf("--file-style %s", info.Options.FilenameStyle))
	}
	
	if len(activeOptions) > 0 {
		output.WriteString("\nOptions:\n")
		for _, opt := range activeOptions {
			output.WriteString(fmt.Sprintf("  %s\n", opt))
		}
	}
	
	return output.String()
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Helper function to get unique strings
func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	
	return result
}

// isTextFileWithExtensions checks if a file is a text file considering additional extensions
func isTextFileWithExtensions(path string, additionalExtensions []string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	
	// Check default extensions
	for _, validExt := range DefaultTextExtensions {
		if ext == validExt {
			return true
		}
	}
	
	// Then check additional extensions
	for _, addExt := range additionalExtensions {
		// Normalize extension (add leading dot if missing)
		if !strings.HasPrefix(addExt, ".") {
			addExt = "." + addExt
		}
		if ext == strings.ToLower(addExt) {
			return true
		}
	}
	
	return false
}

// formatFileSize formats a file size in bytes to a human-readable string
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// countFileLines counts the number of lines in a file, respecting line ranges
func countFileLines(pathWithRange string) (int, error) {
	path, rangeSpec := parsePathWithRange(pathWithRange)
	
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = file.Close()
	}()
	
	// If no range specified, count all lines
	if rangeSpec == "" {
		scanner := bufio.NewScanner(file)
		lineCount := 0
		for scanner.Scan() {
			lineCount++
		}
		if err := scanner.Err(); err != nil {
			return 0, err
		}
		return lineCount, nil
	}
	
	// Parse range specification
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	
	ranges, err := parseRanges(rangeSpec, len(lines))
	if err != nil {
		return 0, err
	}
	
	// Count lines in all ranges
	totalLines := 0
	for _, r := range ranges {
		totalLines += r.End - r.Start + 1
	}
	
	return totalLines, nil
}