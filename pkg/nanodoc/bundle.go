package nanodoc

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)


// BundleResult holds both the raw option lines and file paths from a bundle file
type BundleResult struct {
	// File paths from the bundle
	Paths []string
	// Raw option lines from the bundle (unparsed)
	OptionLines []string
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
	var optionLines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check if this line is a command-line option
		if strings.HasPrefix(line, "-") {
			optionLines = append(optionLines, line)
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
		Paths:       paths,
		OptionLines: optionLines,
	}, nil
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

// BuildDocument creates a Document from resolved paths
// Note: Bundle option processing has been moved to the CLI layer
func BuildDocument(pathInfos []PathInfo, options FormattingOptions) (*Document, error) {
	return BuildDocumentWithOptions(pathInfos, options)
}

// BuildDocumentWithExplicitFlags creates a Document from resolved paths
// Note: Bundle option processing has been moved to the CLI layer
func BuildDocumentWithExplicitFlags(pathInfos []PathInfo, options FormattingOptions, explicitFlags map[string]bool) (*Document, error) {
	// The explicit flags are now handled in the CLI layer
	return BuildDocumentWithOptions(pathInfos, options)
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

	// Create the document
	doc := NewDocument()
	doc.ContentItems = contents
	doc.FormattingOptions = options

	// Process live bundles - integrate both approaches
	if err := ProcessLiveBundles(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// ExtractBundleOptionLines extracts raw option lines from all bundle files
func ExtractBundleOptionLines(pathInfos []PathInfo) ([]string, error) {
	bp := NewBundleProcessor()
	var allOptionLines []string

	// Extract option lines from all bundle files
	for _, info := range pathInfos {
		if info.Type == "bundle" {
			result, err := bp.ProcessBundleFileWithOptions(info.Absolute)
			if err != nil {
				return nil, err
			}
			
			// Collect all option lines
			allOptionLines = append(allOptionLines, result.OptionLines...)
		}
	}

	return allOptionLines, nil
}



// ProcessLiveBundles iterates through document content and processes inline bundles.
func ProcessLiveBundles(doc *Document) error {
	for i := range doc.ContentItems {
		// Skip processing for common documentation files to avoid processing
		// [[file:]] examples as actual directives
		if shouldSkipLiveBundleProcessing(doc.ContentItems[i].Filepath) {
			continue
		}
		
		processedContent, err := ProcessLiveBundle(doc.ContentItems[i].Content)
		if err != nil {
			return err
		}
		doc.ContentItems[i].Content = processedContent
	}
	return nil
}

// shouldSkipLiveBundleProcessing determines if a file should be skipped for live bundle processing
func shouldSkipLiveBundleProcessing(filepath string) bool {
	// Skip common documentation files that might contain [[file:]] examples
	filename := strings.ToLower(filepath)
	return strings.Contains(filename, "readme") ||
		strings.Contains(filename, "changelog") ||
		strings.Contains(filename, "troubleshooting") ||
		strings.Contains(filename, "contributing") ||
		strings.Contains(filename, "license")
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