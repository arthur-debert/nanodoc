package nanodoc

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BundleOptions holds the formatting options parsed from a bundle file
type BundleOptions struct {
	// Theme name to use
	Theme *string

	// Line numbering mode
	LineNumbers *LineNumberMode

	// Whether to show headers
	ShowHeaders *bool

	// Header style
	HeaderStyle *HeaderStyle

	// Header sequence type
	SequenceStyle *SequenceStyle

	// Whether to show table of contents
	ShowTOC *bool

	// Additional file extensions to process
	AdditionalExtensions []string
}

// BundleResult holds both the options and file paths parsed from a bundle file
type BundleResult struct {
	// File paths from the bundle
	Paths []string
	// Options parsed from the bundle
	Options BundleOptions
}

// MergeFormattingOptions merges bundle options with command-line options
// Command-line options override bundle options
func MergeFormattingOptions(bundleOpts BundleOptions, cmdOpts FormattingOptions) FormattingOptions {
	result := cmdOpts // Start with command-line options

	// Only use bundle options if command-line options are at default values
	if bundleOpts.Theme != nil && cmdOpts.Theme == "classic" {
		result.Theme = *bundleOpts.Theme
	}
	if bundleOpts.LineNumbers != nil && cmdOpts.LineNumbers == LineNumberNone {
		result.LineNumbers = *bundleOpts.LineNumbers
	}
	if bundleOpts.ShowHeaders != nil && cmdOpts.ShowHeaders == true {
		result.ShowHeaders = *bundleOpts.ShowHeaders
	}
	if bundleOpts.HeaderStyle != nil && cmdOpts.HeaderStyle == HeaderStyleNice {
		result.HeaderStyle = *bundleOpts.HeaderStyle
	}
	if bundleOpts.SequenceStyle != nil && cmdOpts.SequenceStyle == SequenceNumerical {
		result.SequenceStyle = *bundleOpts.SequenceStyle
	}
	if bundleOpts.ShowTOC != nil && cmdOpts.ShowTOC == false {
		result.ShowTOC = *bundleOpts.ShowTOC
	}
	
	// For additional extensions, append bundle extensions if not already present
	if len(bundleOpts.AdditionalExtensions) > 0 {
		extensionMap := make(map[string]bool)
		for _, ext := range cmdOpts.AdditionalExtensions {
			extensionMap[ext] = true
		}
		for _, ext := range bundleOpts.AdditionalExtensions {
			if !extensionMap[ext] {
				result.AdditionalExtensions = append(result.AdditionalExtensions, ext)
			}
		}
	}

	return result
}

// MergeFormattingOptionsWithDefaults merges bundle options with command-line options
// This function needs to know which command-line options were explicitly set
func MergeFormattingOptionsWithDefaults(bundleOpts BundleOptions, cmdOpts FormattingOptions, explicitFlags map[string]bool) FormattingOptions {
	result := cmdOpts // Start with command-line options

	// Only use bundle options if command-line options were not explicitly set
	if bundleOpts.Theme != nil && !explicitFlags["theme"] {
		result.Theme = *bundleOpts.Theme
	}
	if bundleOpts.LineNumbers != nil && !explicitFlags["line-numbers"] {
		result.LineNumbers = *bundleOpts.LineNumbers
	}
	if bundleOpts.ShowHeaders != nil && !explicitFlags["no-header"] {
		result.ShowHeaders = *bundleOpts.ShowHeaders
	}
	if bundleOpts.HeaderStyle != nil && !explicitFlags["header-style"] {
		result.HeaderStyle = *bundleOpts.HeaderStyle
	}
	if bundleOpts.SequenceStyle != nil && !explicitFlags["sequence"] {
		result.SequenceStyle = *bundleOpts.SequenceStyle
	}
	if bundleOpts.ShowTOC != nil && !explicitFlags["toc"] {
		result.ShowTOC = *bundleOpts.ShowTOC
	}
	
	// For additional extensions, append bundle extensions if not already present
	if len(bundleOpts.AdditionalExtensions) > 0 {
		extensionMap := make(map[string]bool)
		for _, ext := range cmdOpts.AdditionalExtensions {
			extensionMap[ext] = true
		}
		for _, ext := range bundleOpts.AdditionalExtensions {
			if !extensionMap[ext] {
				result.AdditionalExtensions = append(result.AdditionalExtensions, ext)
			}
		}
	}

	return result
}

// BundleProcessor handles bundle file processing and circular dependency detection
type BundleProcessor struct {
	// Track visited bundles to detect circular dependencies
	visitedBundles map[string]bool
	// Track the current path for circular dependency error reporting
	bundlePath []string
}

// NewBundleProcessor creates a new bundle processor
func NewBundleProcessor() *BundleProcessor {
	return &BundleProcessor{
		visitedBundles: make(map[string]bool),
		bundlePath:     make([]string, 0),
	}
}

// ProcessBundleFile reads and processes a bundle file, returning the list of paths it contains
func (bp *BundleProcessor) ProcessBundleFile(bundlePath string) ([]string, error) {
	result, err := bp.ProcessBundleFileWithOptions(bundlePath)
	if err != nil {
		return nil, err
	}
	return result.Paths, nil
}

// ProcessBundleFileWithOptions reads and processes a bundle file, returning both paths and options
func (bp *BundleProcessor) ProcessBundleFileWithOptions(bundlePath string) (*BundleResult, error) {
	// Get absolute path for consistent tracking
	absBundlePath, err := filepath.Abs(bundlePath)
	if err != nil {
		return nil, &FileError{Path: bundlePath, Err: err}
	}

	// Check for circular dependency
	if bp.visitedBundles[absBundlePath] {
		return nil, &CircularDependencyError{
			Path:  absBundlePath,
			Chain: append(bp.bundlePath, absBundlePath),
		}
	}

	// Mark as visited and add to path
	bp.visitedBundles[absBundlePath] = true
	bp.bundlePath = append(bp.bundlePath, absBundlePath)
	defer func() {
		// Remove from path when done
		bp.bundlePath = bp.bundlePath[:len(bp.bundlePath)-1]
	}()

	// Read the bundle file
	file, err := os.Open(bundlePath)
	if err != nil {
		return nil, &FileError{Path: bundlePath, Err: err}
	}
	defer func() {
		_ = file.Close()
	}()

	var paths []string
	var options BundleOptions
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check if this line is a command-line option
		if strings.HasPrefix(line, "-") {
			if err := parseOption(line, &options); err != nil {
				return nil, &FileError{
					Path: bundlePath,
					Err:  fmt.Errorf("error parsing option on line %d: %w", lineNum, err),
				}
			}
			continue
		}

		// Handle file paths - make them relative to the bundle file's directory
		resolvedPath := line
		if !filepath.IsAbs(line) {
			bundleDir := filepath.Dir(bundlePath)
			resolvedPath = filepath.Join(bundleDir, line)
		}

		paths = append(paths, resolvedPath)
	}

	if err := scanner.Err(); err != nil {
		return nil, &FileError{Path: bundlePath, Err: err}
	}

	return &BundleResult{
		Paths:   paths,
		Options: options,
	}, nil
}

// parseOption parses a single command-line option and updates the BundleOptions struct
func parseOption(optionLine string, options *BundleOptions) error {
	// Split the option line into parts
	parts := strings.Fields(optionLine)
	if len(parts) == 0 {
		return fmt.Errorf("empty option line")
	}

	flag := parts[0]
	
	switch flag {
	case "--toc":
		options.ShowTOC = &[]bool{true}[0]
		
	case "--no-header":
		options.ShowHeaders = &[]bool{false}[0]
		
	case "--line-numbers", "-n":
		options.LineNumbers = &[]LineNumberMode{LineNumberFile}[0]
		
	case "--global-line-numbers", "-N":
		options.LineNumbers = &[]LineNumberMode{LineNumberGlobal}[0]
		
	case "--theme":
		if len(parts) < 2 {
			return fmt.Errorf("--theme requires a value")
		}
		options.Theme = &parts[1]
		
	case "--header-style":
		if len(parts) < 2 {
			return fmt.Errorf("--header-style requires a value")
		}
		style := HeaderStyle(parts[1])
		// Validate header style
		switch style {
		case HeaderStyleNice, HeaderStyleFilename, HeaderStylePath:
			options.HeaderStyle = &style
		default:
			return fmt.Errorf("invalid header style: %s (must be one of: nice, filename, path)", parts[1])
		}
		
	case "--sequence":
		if len(parts) < 2 {
			return fmt.Errorf("--sequence requires a value")
		}
		sequence := SequenceStyle(parts[1])
		// Validate sequence style
		switch sequence {
		case SequenceNumerical, SequenceLetter, SequenceRoman:
			options.SequenceStyle = &sequence
		default:
			return fmt.Errorf("invalid sequence style: %s (must be one of: numerical, letter, roman)", parts[1])
		}
		
	case "--txt-ext":
		if len(parts) < 2 {
			return fmt.Errorf("--txt-ext requires a value")
		}
		options.AdditionalExtensions = append(options.AdditionalExtensions, parts[1])
		
	default:
		return fmt.Errorf("unknown option: %s", flag)
	}
	
	return nil
}

// ProcessPaths takes a list of paths and expands any bundle files recursively
func (bp *BundleProcessor) ProcessPaths(paths []string) ([]string, error) {
	var expandedPaths []string

	for _, path := range paths {
		// Check if it's a bundle file
		if isBundleFile(path) {
			// Process the bundle file recursively
			bundlePaths, err := bp.ProcessBundleFile(path)
			if err != nil {
				return nil, err
			}

			// Recursively process the paths from the bundle
			expandedBundlePaths, err := bp.ProcessPaths(bundlePaths)
			if err != nil {
				return nil, err
			}

			expandedPaths = append(expandedPaths, expandedBundlePaths...)
		} else {
			// Regular file, add as-is
			expandedPaths = append(expandedPaths, path)
		}
	}

	return expandedPaths, nil
}

// BuildDocument creates a Document from resolved paths with bundle support
func BuildDocument(pathInfos []PathInfo, options FormattingOptions) (*Document, error) {
	// First, extract bundle options and merge with command-line options
	mergedOptions, err := ExtractAndMergeBundleOptions(pathInfos, options)
	if err != nil {
		return nil, err
	}
	
	return BuildDocumentWithOptions(pathInfos, mergedOptions)
}

// BuildDocumentWithExplicitFlags creates a Document from resolved paths with bundle support and explicit flag tracking
func BuildDocumentWithExplicitFlags(pathInfos []PathInfo, options FormattingOptions, explicitFlags map[string]bool) (*Document, error) {
	// First, extract bundle options and merge with command-line options using explicit flags
	mergedOptions, err := ExtractAndMergeBundleOptionsWithDefaults(pathInfos, options, explicitFlags)
	if err != nil {
		return nil, err
	}
	
	return BuildDocumentWithOptions(pathInfos, mergedOptions)
}

// BuildDocumentWithOptions creates a Document from resolved paths with already-merged options
func BuildDocumentWithOptions(pathInfos []PathInfo, options FormattingOptions) (*Document, error) {
	bp := NewBundleProcessor()
	var allPaths []string

	// First, collect all paths from PathInfo
	for _, info := range pathInfos {
		switch info.Type {
		case "file":
			allPaths = append(allPaths, info.Original)
		case "directory", "glob":
			allPaths = append(allPaths, info.Files...)
		case "bundle":
			allPaths = append(allPaths, info.Absolute)
		}
	}

	// Process paths to expand bundles
	expandedPaths, err := bp.ProcessPaths(allPaths)
	if err != nil {
		return nil, err
	}

	// Create PathInfo objects for expanded paths, treating them all as files
	var resolvedInfos []PathInfo
	for _, path := range expandedPaths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, &FileError{Path: path, Err: err}
		}
		
		// Treat all expanded paths as files, not bundles
		resolvedInfos = append(resolvedInfos, PathInfo{
			Original: path,
			Absolute: absPath,
			Type:     "file",
		})
	}

	// Extract content from all files
	contents, err := ResolveAndExtractFiles(resolvedInfos, options.AdditionalExtensions)
	if err != nil {
		return nil, err
	}

	// Gather content with range support
	gatheredContents, err := GatherContentWithRanges(contents)
	if err != nil {
		return nil, err
	}

	// Create the document
	doc := NewDocument()
	doc.ContentItems = gatheredContents
	doc.FormattingOptions = options

	// Process live bundles - integrate both approaches
	if err := ProcessLiveBundles(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// ExtractAndMergeBundleOptions extracts options from bundle files and merges them with command-line options
func ExtractAndMergeBundleOptions(pathInfos []PathInfo, cmdOptions FormattingOptions) (FormattingOptions, error) {
	bp := NewBundleProcessor()
	var bundleOptions BundleOptions

	// Extract options from all bundle files
	for _, info := range pathInfos {
		if info.Type == "bundle" {
			result, err := bp.ProcessBundleFileWithOptions(info.Absolute)
			if err != nil {
				return cmdOptions, err
			}
			
			// Merge bundle options (first bundle file wins for conflicting options)
			bundleOptions = mergeBundleOptions(bundleOptions, result.Options)
		}
	}

	// Merge bundle options with command-line options
	return MergeFormattingOptions(bundleOptions, cmdOptions), nil
}

// ExtractAndMergeBundleOptionsWithDefaults extracts options from bundle files and merges them with command-line options using explicit flags
func ExtractAndMergeBundleOptionsWithDefaults(pathInfos []PathInfo, cmdOptions FormattingOptions, explicitFlags map[string]bool) (FormattingOptions, error) {
	bp := NewBundleProcessor()
	var bundleOptions BundleOptions

	// Extract options from all bundle files
	for _, info := range pathInfos {
		if info.Type == "bundle" {
			result, err := bp.ProcessBundleFileWithOptions(info.Absolute)
			if err != nil {
				return cmdOptions, err
			}
			
			// Merge bundle options (first bundle file wins for conflicting options)
			bundleOptions = mergeBundleOptions(bundleOptions, result.Options)
		}
	}

	// Merge bundle options with command-line options using explicit flags
	return MergeFormattingOptionsWithDefaults(bundleOptions, cmdOptions, explicitFlags), nil
}

// mergeBundleOptions merges two BundleOptions structures
// The first one takes precedence for conflicting options
func mergeBundleOptions(first, second BundleOptions) BundleOptions {
	result := first
	
	if result.Theme == nil && second.Theme != nil {
		result.Theme = second.Theme
	}
	if result.LineNumbers == nil && second.LineNumbers != nil {
		result.LineNumbers = second.LineNumbers
	}
	if result.ShowHeaders == nil && second.ShowHeaders != nil {
		result.ShowHeaders = second.ShowHeaders
	}
	if result.HeaderStyle == nil && second.HeaderStyle != nil {
		result.HeaderStyle = second.HeaderStyle
	}
	if result.SequenceStyle == nil && second.SequenceStyle != nil {
		result.SequenceStyle = second.SequenceStyle
	}
	if result.ShowTOC == nil && second.ShowTOC != nil {
		result.ShowTOC = second.ShowTOC
	}
	
	// For additional extensions, merge them
	result.AdditionalExtensions = append(result.AdditionalExtensions, second.AdditionalExtensions...)
	
	return result
}

// ProcessLiveBundles iterates through document content and processes inline bundles.
func ProcessLiveBundles(doc *Document) error {
	for i := range doc.ContentItems {
		processedContent, err := ProcessLiveBundle(doc.ContentItems[i].Content)
		if err != nil {
			return err
		}
		doc.ContentItems[i].Content = processedContent
	}
	return nil
}

// ProcessLiveBundle handles inline bundle processing
// It looks for directives like [[file:path/to/file.txt]] or [[file:path/to/file.txt:L10-20]]
// and replaces them with the actual file content
func ProcessLiveBundle(content string) (string, error) {
	return processLiveBundleRecursive(content, 0, make(map[string]bool))
}

func processLiveBundleRecursive(content string, depth int, visited map[string]bool) (string, error) {
	// Prevent infinite recursion
	const maxDepth = 10
	if depth > maxDepth {
		return "", &CircularDependencyError{
			Path:  "live bundle",
			Chain: []string{"Maximum nesting depth exceeded"},
		}
	}

	// Process all directives in the content
	result := content
	startPos := 0
	
	for {
		// Find the next directive
		loc := strings.Index(result[startPos:], "[[file:")
		if loc == -1 {
			break
		}
		
		// Adjust location to absolute position
		loc += startPos
		
		// Find the closing ]]
		endLoc := strings.Index(result[loc:], "]]")
		if endLoc == -1 {
			// Malformed directive, skip it
			startPos = loc + 7 // len("[[file:")
			continue
		}
		endLoc += loc + 2 // Include the ]]
		
		// Parse the file path (and optional range)
		pathStart := loc + 7 // len("[[file:")
		pathEnd := endLoc - 2 // Before ]]
		pathWithRange := result[pathStart:pathEnd]
		
		// Check for circular references
		if visited[pathWithRange] {
			return "", &CircularDependencyError{
				Path:  pathWithRange,
				Chain: mapKeysToSlice(visited),
			}
		}
		
		// Mark as visited
		visited[pathWithRange] = true
		
		// Extract the file content
		fileContent, err := ExtractFileContent(pathWithRange)
		if err != nil {
			// On error, leave the directive as-is and continue
			startPos = endLoc
			continue
		}
		
		// Process nested directives in the included content
		processedContent, err := processLiveBundleRecursive(fileContent.Content, depth+1, visited)
		if err != nil {
			return "", err
		}
		
		// Replace the directive with the content
		result = result[:loc] + processedContent + result[endLoc:]
		
		// Update start position
		startPos = loc + len(processedContent)
		
		// Remove from visited after processing
		delete(visited, pathWithRange)
	}
	
	return result, nil
}

// Helper function to convert map keys to slice
func mapKeysToSlice(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}