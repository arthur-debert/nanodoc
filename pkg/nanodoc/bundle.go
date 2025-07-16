package nanodoc

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

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
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle relative paths - make them relative to the bundle file's directory
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

	return paths, nil
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

	// Now resolve the expanded paths
	resolvedInfos, err := ResolvePaths(expandedPaths)
	if err != nil {
		return nil, err
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

	return doc, nil
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
