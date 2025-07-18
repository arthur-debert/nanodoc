package nanodoc

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PathInfo represents information about a resolved path
type PathInfo struct {
	// Original path as provided by user
	Original string

	// Absolute resolved path
	Absolute string

	// Type of path: "file", "directory", or "bundle"
	Type string

	// If directory, the files found within
	Files []string
}

// ResolvePaths takes a list of source paths and resolves them to absolute paths
// It handles files, directories, and bundle files
func ResolvePaths(sources []string) ([]PathInfo, error) {
	return ResolvePathsWithOptions(sources, nil)
}

// ResolvePathsWithOptions resolves paths with optional pattern filtering
func ResolvePathsWithOptions(sources []string, options *FormattingOptions) ([]PathInfo, error) {
	if len(sources) == 0 {
		return nil, ErrEmptySource
	}

	results := make([]PathInfo, 0, len(sources))

	for _, source := range sources {
		pathInfo, err := resolveSinglePathWithOptions(source, options)
		if err != nil {
			return nil, &FileError{Path: source, Err: err}
		}
		results = append(results, pathInfo)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Absolute < results[j].Absolute
	})

	return results, nil
}

// resolveSinglePath resolves a single path to PathInfo
func resolveSinglePath(path string) (PathInfo, error) {
	return resolveSinglePathWithOptions(path, nil)
}

// resolveSinglePathWithOptions resolves a single path with optional pattern filtering
func resolveSinglePathWithOptions(path string, options *FormattingOptions) (PathInfo, error) {
	if strings.ContainsAny(path, "*?[") {
		return resolveGlobPathWithOptions(path, options)
	}
	return resolveNonGlobPathWithOptions(path, options)
}

// resolveNonGlobPath handles resolving a path that is not a glob pattern.
func resolveNonGlobPath(path string) (PathInfo, error) {
	return resolveNonGlobPathWithOptions(path, nil)
}

// resolveNonGlobPathWithOptions handles resolving a path with optional pattern filtering
func resolveNonGlobPathWithOptions(path string, options *FormattingOptions) (PathInfo, error) {
	// Parse out any range specification for file system operations
	// but keep the original path with range for later processing
	basePath := path
	if idx := strings.LastIndex(path, ":L"); idx > 0 {
		basePath = path[:idx]
	}

	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return PathInfo{}, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return PathInfo{}, ErrFileNotFound
		}
		return PathInfo{}, err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		realPath, err := filepath.EvalSymlinks(absPath)
		if err != nil {
			return PathInfo{}, err
		}
		absPath = realPath
		info, err = os.Stat(absPath)
		if err != nil {
			return PathInfo{}, err
		}
	}

	pathInfo := PathInfo{
		Original: path,
		Absolute: absPath,
	}

	if info.IsDir() {
		return handleDirectoryWithOptions(pathInfo, options)
	}
	return handleFile(pathInfo)
}


// handleDirectoryWithOptions processes a directory path with optional pattern filtering
func handleDirectoryWithOptions(pathInfo PathInfo, options *FormattingOptions) (PathInfo, error) {
	pathInfo.Type = "directory"
	
	var files []string
	var err error
	
	// Check if we need pattern-based filtering
	if options != nil && (len(options.IncludePatterns) > 0 || len(options.ExcludePatterns) > 0) {
		matcher := NewPatternMatcher(pathInfo.Absolute, options.IncludePatterns, options.ExcludePatterns)
		
		if matcher.NeedsRecursion() {
			files, err = findTextFilesRecursive(pathInfo.Absolute, options.AdditionalExtensions, matcher)
		} else {
			files, err = findTextFilesWithMatcher(pathInfo.Absolute, options.AdditionalExtensions, matcher)
		}
	} else {
		// No patterns, use existing behavior
		files, err = findTextFilesInDir(pathInfo.Absolute)
	}
	
	if err != nil {
		return PathInfo{}, err
	}
	pathInfo.Files = files
	return pathInfo, nil
}

// handleFile processes a file path.
func handleFile(pathInfo PathInfo) (PathInfo, error) {
	if isBundleFile(pathInfo.Absolute) {
		pathInfo.Type = "bundle"
	} else {
		pathInfo.Type = "file"
	}
	return pathInfo, nil
}

// resolveGlobPath resolves a glob pattern to matching files
func resolveGlobPath(pattern string) (PathInfo, error) {
	return resolveGlobPathWithOptions(pattern, nil)
}

// resolveGlobPathWithOptions resolves a glob pattern with optional additional filtering
func resolveGlobPathWithOptions(pattern string, options *FormattingOptions) (PathInfo, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return PathInfo{}, err
	}

	if len(matches) == 0 {
		return PathInfo{}, ErrFileNotFound
	}

	// Filter to only include files (not directories)
	var files []string
	for _, match := range matches {
		absPath, err := filepath.Abs(match)
		if err != nil {
			continue
		}

		info, err := os.Stat(absPath)
		if err != nil {
			continue
		}

		if !info.IsDir() {
			isText := isTextFile(absPath)
			if options != nil && len(options.AdditionalExtensions) > 0 {
				isText = isTextFileWithExtensions(absPath, options.AdditionalExtensions)
			}
			if isText {
				files = append(files, absPath)
			}
		}
	}

	if len(files) == 0 {
		return PathInfo{}, ErrFileNotFound
	}

	sortPaths(files)

	return PathInfo{
		Original: pattern,
		Type:     "glob",
		Files:    files,
	}, nil
}

// isBundleFile checks if a file is a bundle file based on naming convention
func isBundleFile(path string) bool {
	base := filepath.Base(path)
	return strings.Contains(base, BundlePattern)
}

// findTextFilesInDir finds all text files in a directory
func findTextFilesInDir(dir string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(dir, entry.Name())
		if isTextFile(fullPath) {
			files = append(files, fullPath)
		}
	}

	// Sort files for consistent ordering
	sortPaths(files)

	return files, nil
}

// findTextFilesWithMatcher finds text files in a directory with pattern matching
func findTextFilesWithMatcher(dir string, additionalExtensions []string, matcher *PatternMatcher) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(dir, entry.Name())
		if isTextFileWithExtensions(fullPath, additionalExtensions) {
			shouldInclude, err := matcher.ShouldInclude(fullPath)
			if err != nil {
				return nil, err
			}
			if shouldInclude {
				files = append(files, fullPath)
			}
		}
	}

	sortPaths(files)
	return files, nil
}

// findTextFilesRecursive recursively finds text files with pattern matching
func findTextFilesRecursive(dir string, additionalExtensions []string, matcher *PatternMatcher) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if isTextFileWithExtensions(path, additionalExtensions) {
			shouldInclude, err := matcher.ShouldInclude(path)
			if err != nil {
				return err
			}
			if shouldInclude {
				files = append(files, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	sortPaths(files)
	return files, nil
}

// isTextFile checks if a file has a text extension
func isTextFile(path string) bool {
	return isTextFileWithExtensions(path, nil)
}

// sortPaths sorts paths alphabetically
func sortPaths(paths []string) {
	sort.Strings(paths)
}

// GetFilesFromDirectory is a helper that returns all text files from a directory
func GetFilesFromDirectory(dir string, extensions []string) ([]string, error) {
	if extensions == nil {
		extensions = DefaultTextExtensions
	}

	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(dir, entry.Name())
		ext := strings.ToLower(filepath.Ext(fullPath))

		for _, validExt := range extensions {
			if ext == validExt {
				files = append(files, fullPath)
				break
			}
		}
	}

	sortPaths(files)
	return files, nil
}
