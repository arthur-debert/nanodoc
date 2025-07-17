package nanodoc

import (
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
	// Files requiring additional extensions
	RequiresExtension map[string]string
}

// FileInfo contains dry run information about a file
type FileInfo struct {
	Path      string
	Source    string // Where it came from (directory, bundle, etc.)
	Extension string
	Size      int64  // File size in bytes
	ModTime   string // Modification time
}

// GenerateDryRunInfo analyzes what files would be processed without actually processing them
func GenerateDryRunInfo(pathInfos []PathInfo, additionalExtensions []string) (*DryRunInfo, error) {
	info := &DryRunInfo{
		Files:             make([]FileInfo, 0),
		Bundles:           make([]string, 0),
		RequiresExtension: make(map[string]string),
	}

	// Process each path
	for _, pathInfo := range pathInfos {
		switch pathInfo.Type {
		case "file":
			ext := filepath.Ext(pathInfo.Absolute)
			fileInfo := FileInfo{
				Path:      pathInfo.Absolute,
				Source:    "direct argument",
				Extension: ext,
			}
			
			// Get file metadata
			if stat, err := os.Stat(pathInfo.Absolute); err == nil {
				fileInfo.Size = stat.Size()
				fileInfo.ModTime = stat.ModTime().Format("2006-01-02 15:04:05")
			}
			
			// Check if file needs additional extension
			if !isTextFileWithExtensions(pathInfo.Absolute, additionalExtensions) {
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
				
				// Get file metadata
				if stat, err := os.Stat(file); err == nil {
					fileInfo.Size = stat.Size()
					fileInfo.ModTime = stat.ModTime().Format("2006-01-02 15:04:05")
				}
				
				info.Files = append(info.Files, fileInfo)
			}
			
		case "glob":
			for _, file := range pathInfo.Files {
				fileInfo := FileInfo{
					Path:      file,
					Source:    fmt.Sprintf("glob: %s", pathInfo.Original),
					Extension: filepath.Ext(file),
				}
				
				// Get file metadata
				if stat, err := os.Stat(file); err == nil {
					fileInfo.Size = stat.Size()
					fileInfo.ModTime = stat.ModTime().Format("2006-01-02 15:04:05")
				}
				
				info.Files = append(info.Files, fileInfo)
			}
			
		case "bundle":
			info.Bundles = append(info.Bundles, pathInfo.Absolute)
			// For dry run, we don't read bundle contents to avoid file I/O
			// Just add the bundle itself as a file entry
			fileInfo := FileInfo{
				Path:      pathInfo.Absolute,
				Source:    "bundle file",
				Extension: filepath.Ext(pathInfo.Absolute),
			}
			
			// Get file metadata
			if stat, err := os.Stat(pathInfo.Absolute); err == nil {
				fileInfo.Size = stat.Size()
				fileInfo.ModTime = stat.ModTime().Format("2006-01-02 15:04:05")
			}
			
			info.Files = append(info.Files, fileInfo)
		}
	}
	
	// Remove duplicates and sort
	info.Bundles = uniqueStrings(info.Bundles)
	sort.Strings(info.Bundles)
	
	// Count total files
	info.TotalFiles = len(info.Files)
	
	return info, nil
}

// FormatDryRunOutput formats the dry run information for display
func FormatDryRunOutput(info *DryRunInfo) string {
	var output strings.Builder
	
	output.WriteString("Would process the following files:\n")
	
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
			// Format size in human-readable format
			sizeStr := formatFileSize(file.Size)
			output.WriteString(fmt.Sprintf("%d. %s (%s, %s)\n", fileNum, relPath, sizeStr, file.ModTime))
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
		output.WriteString("\nFiles requiring --txt-ext flag:\n")
		for file, ext := range info.RequiresExtension {
			output.WriteString(fmt.Sprintf("  - %s (requires --txt-ext=%s)\n", 
				filepath.Base(file), strings.TrimPrefix(ext, ".")))
		}
	}
	
	// Summary
	output.WriteString(fmt.Sprintf("\nTotal files to process: %d\n", info.TotalFiles))
	
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
	// First check default extensions
	if isTextFile(path) {
		return true
	}
	
	// Then check additional extensions
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	for _, addExt := range additionalExtensions {
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