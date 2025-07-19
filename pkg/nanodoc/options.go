package nanodoc

import (
	"strings"

	"github.com/spf13/cobra"
)

// ParseBundleOptions parses bundle option lines using Cobra
func ParseBundleOptions(optionLines []string) (FormattingOptions, error) {
	// Create a temporary command to parse options
	tempCmd := &cobra.Command{}
	
	// Set up the same flags as the root command
	var bundleLineNum string
	var bundleToc bool
	var bundleTheme string
	var bundleShowFilenames bool
	var bundleFileNumbering string
	var bundleFilenameFormat string
	var bundleFilenameAlign string
	var bundleFilenameBanner string
	var bundlePageWidth int
	var bundleAdditionalExt []string
	var bundleIncludePatterns []string
	var bundleExcludePatterns []string
	
	tempCmd.Flags().StringVarP(&bundleLineNum, "linenum", "l", "", "")
	tempCmd.Flags().BoolVar(&bundleToc, "toc", false, "")
	tempCmd.Flags().StringVar(&bundleTheme, "theme", "classic", "")
	tempCmd.Flags().BoolVar(&bundleShowFilenames, "filenames", true, "")
	tempCmd.Flags().StringVar(&bundleFilenameFormat, "header-format", "nice", "")
	tempCmd.Flags().StringVar(&bundleFilenameAlign, "header-align", "left", "")
	tempCmd.Flags().StringVar(&bundleFilenameBanner, "header-style", "none", "")
	tempCmd.Flags().IntVar(&bundlePageWidth, "page-width", OUTPUT_WIDTH, "")
	tempCmd.Flags().StringVar(&bundleFileNumbering, "file-numbering", "numerical", "")
	tempCmd.Flags().StringSliceVar(&bundleAdditionalExt, "ext", []string{}, "")
	tempCmd.Flags().StringSliceVar(&bundleIncludePatterns, "include", []string{}, "")
	tempCmd.Flags().StringSliceVar(&bundleExcludePatterns, "exclude", []string{}, "")
	
	// Parse the option lines
	// Need to split options that have values into separate elements
	var args []string
	for _, line := range optionLines {
		// Split by spaces to separate flag and value
		parts := strings.Fields(line)
		args = append(args, parts...)
	}
	
	if err := tempCmd.ParseFlags(args); err != nil {
		return FormattingOptions{}, err
	}
	
	// Convert to FormattingOptions
	lineNumberMode := LineNumberNone
	switch bundleLineNum {
	case "file":
		lineNumberMode = LineNumberFile
	case "global":
		lineNumberMode = LineNumberGlobal
	}
	
	return FormattingOptions{
		LineNumbers:          lineNumberMode,
		ShowTOC:              bundleToc,
		Theme:                bundleTheme,
		ShowFilenames:        bundleShowFilenames,
		SequenceStyle:        SequenceStyle(bundleFileNumbering),
		HeaderFormat:         HeaderFormat(bundleFilenameFormat),
		HeaderAlignment:      bundleFilenameAlign,
		HeaderStyle:          bundleFilenameBanner,
		PageWidth:            bundlePageWidth,
		AdditionalExtensions: bundleAdditionalExt,
		IncludePatterns:      bundleIncludePatterns,
		ExcludePatterns:      bundleExcludePatterns,
	}, nil
}

// TrackExplicitFlags determines which flags were explicitly set by the user
func TrackExplicitFlags(cmd *cobra.Command) map[string]bool {
	explicitFlags := make(map[string]bool)
	
	// Check if each flag was explicitly set
	if cmd.Flags().Changed("toc") {
		explicitFlags["toc"] = true
	}
	if cmd.Flags().Changed("theme") {
		explicitFlags["theme"] = true
	}
	if cmd.Flags().Changed("linenum") {
		explicitFlags["line-numbers"] = true
	}
	if cmd.Flags().Changed("filenames") {
		explicitFlags["no-header"] = true
	}
	if cmd.Flags().Changed("header-format") {
		explicitFlags["header-format"] = true
	}
	if cmd.Flags().Changed("header-align") {
		explicitFlags["header-align"] = true
	}
	if cmd.Flags().Changed("header-style") {
		explicitFlags["header-style"] = true
	}
	if cmd.Flags().Changed("page-width") {
		explicitFlags["page-width"] = true
	}
	if cmd.Flags().Changed("file-numbering") {
		explicitFlags["sequence"] = true
	}
	if cmd.Flags().Changed("ext") {
		explicitFlags["txt-ext"] = true
	}
	if cmd.Flags().Changed("include") {
		explicitFlags["include"] = true
	}
	if cmd.Flags().Changed("exclude") {
		explicitFlags["exclude"] = true
	}
	
	return explicitFlags
}

// MergeOptionsWithExplicitFlags merges bundle options with command options based on explicit flags
func MergeOptionsWithExplicitFlags(bundleOpts, cmdOpts FormattingOptions, explicitFlags map[string]bool) FormattingOptions {
	result := cmdOpts
	
	// Only use bundle options if command-line options were not explicitly set
	if !explicitFlags["theme"] {
		result.Theme = bundleOpts.Theme
	}
	if !explicitFlags["line-numbers"] {
		result.LineNumbers = bundleOpts.LineNumbers
	}
	if !explicitFlags["no-header"] {
		result.ShowFilenames = bundleOpts.ShowFilenames
	}
	if !explicitFlags["header-format"] {
		result.HeaderFormat = bundleOpts.HeaderFormat
	}
	if !explicitFlags["header-align"] {
		result.HeaderAlignment = bundleOpts.HeaderAlignment
	}
	if !explicitFlags["header-style"] {
		result.HeaderStyle = bundleOpts.HeaderStyle
	}
	if !explicitFlags["page-width"] {
		result.PageWidth = bundleOpts.PageWidth
	}
	if !explicitFlags["sequence"] {
		result.SequenceStyle = bundleOpts.SequenceStyle
	}
	if !explicitFlags["toc"] {
		result.ShowTOC = bundleOpts.ShowTOC
	}
	
	// Merge additional extensions (bundle + command line)
	if len(bundleOpts.AdditionalExtensions) > 0 && !explicitFlags["txt-ext"] {
		result.AdditionalExtensions = append(bundleOpts.AdditionalExtensions, result.AdditionalExtensions...)
	}
	
	// Merge patterns
	if len(bundleOpts.IncludePatterns) > 0 && !explicitFlags["include"] {
		result.IncludePatterns = append(bundleOpts.IncludePatterns, result.IncludePatterns...)
	}
	if len(bundleOpts.ExcludePatterns) > 0 && !explicitFlags["exclude"] {
		result.ExcludePatterns = append(bundleOpts.ExcludePatterns, result.ExcludePatterns...)
	}
	
	return result
}