package nanodoc

import (
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// PatternMatcher handles include/exclude pattern matching for files
type PatternMatcher struct {
	includePatterns []string
	excludePatterns []string
	baseDir         string
	needsRecursion  bool
}

// NewPatternMatcher creates a new pattern matcher
func NewPatternMatcher(baseDir string, includePatterns, excludePatterns []string) *PatternMatcher {
	pm := &PatternMatcher{
		includePatterns: includePatterns,
		excludePatterns: excludePatterns,
		baseDir:         baseDir,
	}
	
	// Check if any pattern requires recursion
	pm.needsRecursion = pm.hasRecursivePattern()
	
	return pm
}

// hasRecursivePattern checks if any pattern contains ** which requires recursive traversal
func (pm *PatternMatcher) hasRecursivePattern() bool {
	for _, pattern := range pm.includePatterns {
		if strings.Contains(pattern, "**") {
			return true
		}
	}
	for _, pattern := range pm.excludePatterns {
		if strings.Contains(pattern, "**") {
			return true
		}
	}
	return false
}

// NeedsRecursion returns true if recursive directory traversal is needed
func (pm *PatternMatcher) NeedsRecursion() bool {
	return pm.needsRecursion
}

// ShouldInclude determines if a file should be included based on patterns
func (pm *PatternMatcher) ShouldInclude(filePath string) (bool, error) {
	// Get relative path from base directory
	relPath, err := filepath.Rel(pm.baseDir, filePath)
	if err != nil {
		// If we can't get relative path, use the full path
		relPath = filePath
	}
	
	// Normalize path separators for pattern matching
	relPath = filepath.ToSlash(relPath)
	
	// Check include patterns
	included := true
	if len(pm.includePatterns) > 0 {
		included = false
		for _, pattern := range pm.includePatterns {
			match, err := doublestar.Match(pattern, relPath)
			if err != nil {
				return false, err
			}
			if match {
				included = true
				break
			}
		}
	}
	
	// If not included, no need to check excludes
	if !included {
		return false, nil
	}
	
	// Check exclude patterns - they take precedence
	for _, pattern := range pm.excludePatterns {
		match, err := doublestar.Match(pattern, relPath)
		if err != nil {
			return false, err
		}
		if match {
			return false, nil
		}
	}
	
	return true, nil
}

// HasPatterns returns true if any include or exclude patterns are specified
func (pm *PatternMatcher) HasPatterns() bool {
	return len(pm.includePatterns) > 0 || len(pm.excludePatterns) > 0
}