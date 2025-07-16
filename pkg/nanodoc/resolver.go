package nanodoc

import (
	"os"
	"path/filepath"
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
	if len(sources) == 0 {
		return nil, ErrEmptySource
	}

	results := make([]PathInfo, 0, len(sources))

	for _, source := range sources {
		pathInfo, err := resolveSinglePath(source)
		if err != nil {
			return nil, &FileError{Path: source, Err: err}
		}
		results = append(results, pathInfo)
	}

	return results, nil
}

// resolveSinglePath resolves a single path to PathInfo
func resolveSinglePath(path string) (PathInfo, error) {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return PathInfo{}, err
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return PathInfo{}, ErrFileNotFound
		}
		return PathInfo{}, err
	}

	pathInfo := PathInfo{
		Original: path,
		Absolute: absPath,
	}

	if info.IsDir() {
		// Handle directory
		pathInfo.Type = "directory"
		files, err := findTextFilesInDir(absPath)
		if err != nil {
			return PathInfo{}, err
		}
		pathInfo.Files = files
	} else {
		// Check if it's a bundle file
		if isBundleFile(absPath) {
			pathInfo.Type = "bundle"
		} else {
			pathInfo.Type = "file"
		}
	}

	return pathInfo, nil
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

// isTextFile checks if a file has a text extension
func isTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, validExt := range DefaultTextExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// sortPaths sorts paths alphabetically
func sortPaths(paths []string) {
	// Simple bubble sort for now (can optimize later if needed)
	n := len(paths)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if paths[j] > paths[j+1] {
				paths[j], paths[j+1] = paths[j+1], paths[j]
			}
		}
	}
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