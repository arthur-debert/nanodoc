package nanodoc

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
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

	// Process live bundles
	if err := ProcessLiveBundles(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// ProcessLiveBundles iterates through document content and processes inline bundles.
func ProcessLiveBundles(doc *Document) error {
	for i := range doc.ContentItems {
		// Pass the directory of the current file for resolving relative paths
		baseDir := filepath.Dir(doc.ContentItems[i].Filepath)
		processedContent, err := ProcessLiveBundle(doc.ContentItems[i].Content, baseDir)
		if err != nil {
			return err
		}
		doc.ContentItems[i].Content = processedContent
	}
	return nil
}

// ProcessLiveBundle handles inline bundle processing for a single content string.
func ProcessLiveBundle(content, baseDir string) (string, error) {
	// Regex to find !bundle(path)
	re := regexp.MustCompile(`!bundle\(([^)]+)\)`)

	// Find all matches
	matches := re.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		return content, nil
	}

	// Process each match
	for _, match := range matches {
		fullMatch := match[0]
		bundlePath := match[1]

		// Resolve path relative to the file being processed
		if !filepath.IsAbs(bundlePath) {
			bundlePath = filepath.Join(baseDir, bundlePath)
		}

		// Read the content of the bundled file
		bundleContent, err := os.ReadFile(bundlePath)
		if err != nil {
			// If file not found, leave the directive as is
			if os.IsNotExist(err) {
				continue
			}
			return "", &FileError{Path: bundlePath, Err: err}
		}

		// Recursively process the content of the bundled file
		processedBundleContent, err := ProcessLiveBundle(string(bundleContent), filepath.Dir(bundlePath))
		if err != nil {
			return "", err
		}

		// Replace the directive with the file content
		content = strings.Replace(content, fullMatch, processedBundleContent, 1)
	}

	return content, nil
}
