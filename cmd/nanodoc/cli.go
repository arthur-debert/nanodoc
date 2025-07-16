package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/arthur-debert/nanodoc-go/pkg/nanodoc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// CLI flags
var (
	lineNumberMode    string
	showTOC           bool
	themeName         string
	noHeader          bool
	sequenceType      string
	headerStyle       string
	additionalExts    []string
	verboseMode       bool
)

func init() {
	// Configure logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Set up root command
	rootCmd.Run = runNanodoc
	rootCmd.Args = cobra.ArbitraryArgs
	rootCmd.Example = `  # Bundle all .txt and .md files in current directory
  nanodoc

  # Bundle specific files
  nanodoc file1.txt file2.md

  # Use a bundle file
  nanodoc project.bundle.txt

  # With line numbers per file
  nanodoc -n file1.txt file2.txt

  # With global line numbers
  nanodoc -nn *.txt

  # With table of contents
  nanodoc --toc chapter*.md

  # With dark theme
  nanodoc --theme=classic-dark *.md`

	// Add actual flags
	rootCmd.Flags().BoolVarP(&verboseMode, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().StringVarP(&lineNumberMode, "line-numbers", "n", "", "Line numbering mode: 'file' (-n) or 'global' (-nn)")
	rootCmd.Flags().BoolVar(&showTOC, "toc", false, "Generate table of contents")
	rootCmd.Flags().StringVar(&themeName, "theme", "classic", "Theme name (classic, classic-light, classic-dark)")
	rootCmd.Flags().BoolVar(&noHeader, "no-header", false, "Disable file headers")
	rootCmd.Flags().StringVar(&sequenceType, "sequence", "numerical", "Header sequence type (numerical, letter, roman)")
	rootCmd.Flags().StringVar(&headerStyle, "style", "nice", "Header style (nice, filename, path)")
	rootCmd.Flags().StringSliceVar(&additionalExts, "txt-ext", nil, "Additional file extensions to process")

	// Support -nn for global line numbers
	rootCmd.Flags().BoolP("nn", "", false, "Global line numbering (shorthand for -n global)")
	rootCmd.Flags().Lookup("nn").NoOptDefVal = "true"
}

func runNanodoc(cmd *cobra.Command, args []string) {
	// Handle verbose mode
	if verboseMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Verbose mode enabled")
	}

	// Handle -nn flag
	if cmd.Flags().Changed("nn") {
		lineNumberMode = "global"
	} else if cmd.Flags().Changed("n") && lineNumberMode == "" {
		lineNumberMode = "file"
	}

	// Validate flags
	if err := validateFlags(); err != nil {
		exitWithError(err, 2)
	}

	// If no arguments, use current directory
	sources := args
	if len(sources) == 0 {
		sources = []string{"."}
		log.Debug().Msg("No arguments provided, using current directory")
	}

	// Create formatting options
	options := nanodoc.FormattingOptions{
		LineNumberMode:       lineNumberMode,
		ShowHeader:           !noHeader,
		Sequence:             sequenceType,
		Style:                headerStyle,
		AdditionalExtensions: additionalExts,
	}

	log.Debug().
		Str("lineNumberMode", options.LineNumberMode).
		Bool("showHeader", options.ShowHeader).
		Str("sequence", options.Sequence).
		Str("style", options.Style).
		Strs("extensions", options.AdditionalExtensions).
		Msg("Formatting options")

	// Process the document
	output, err := processDocument(sources, options, showTOC, themeName)
	if err != nil {
		exitWithError(err, 1)
	}

	// Output the result
	fmt.Print(output)
}

func validateFlags() error {
	// Validate line number mode
	if lineNumberMode != "" && lineNumberMode != "file" && lineNumberMode != "global" {
		return fmt.Errorf("invalid line number mode '%s': must be 'file' or 'global'", lineNumberMode)
	}

	// Validate sequence type
	validSequences := []string{"numerical", "letter", "roman"}
	if !contains(validSequences, sequenceType) {
		return fmt.Errorf("invalid sequence type '%s': must be one of %s", sequenceType, strings.Join(validSequences, ", "))
	}

	// Validate header style
	validStyles := []string{"nice", "filename", "path"}
	if !contains(validStyles, headerStyle) {
		return fmt.Errorf("invalid header style '%s': must be one of %s", headerStyle, strings.Join(validStyles, ", "))
	}

	// Validate theme
	validThemes := []string{"classic", "classic-light", "classic-dark"}
	if !contains(validThemes, themeName) && themeName != "" {
		// Allow custom themes, just log a warning
		log.Warn().Str("theme", themeName).Msg("Using custom theme")
	}

	return nil
}

func processDocument(sources []string, options nanodoc.FormattingOptions, toc bool, theme string) (string, error) {
	// Step 1: Resolve paths
	log.Debug().Strs("sources", sources).Msg("Resolving paths")
	pathInfos, err := nanodoc.ResolvePaths(sources)
	if err != nil {
		return "", fmt.Errorf("failed to resolve paths: %w", err)
	}

	// Step 2: Build document (handles bundles)
	log.Debug().Int("pathCount", len(pathInfos)).Msg("Building document")
	doc, err := nanodoc.BuildDocument(pathInfos, options)
	if err != nil {
		return "", fmt.Errorf("failed to build document: %w", err)
	}

	// TODO: Implement theme loading and rendering
	// For now, return a simple concatenation of content
	var output strings.Builder
	for i, item := range doc.ContentItems {
		if i > 0 {
			output.WriteString("\n")
		}
		output.WriteString(item.Content)
	}

	return output.String(), nil
}

func exitWithError(err error, code int) {
	// Format error based on type
	switch e := err.(type) {
	case *nanodoc.FileError:
		fmt.Fprintf(os.Stderr, "Error: Cannot access file '%s': %v\n", e.Path, e.Err)
	case *nanodoc.CircularDependencyError:
		fmt.Fprintf(os.Stderr, "Error: Circular dependency detected in bundle '%s'\n", e.Path)
		if len(e.Chain) > 0 {
			fmt.Fprintf(os.Stderr, "Dependency chain: %s\n", strings.Join(e.Chain, " -> "))
		}
	case *nanodoc.RangeError:
		fmt.Fprintf(os.Stderr, "Error: Invalid range '%s': %v\n", e.Input, e.Err)
	default:
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	// Add help text for common errors
	if strings.Contains(err.Error(), "no such file") {
		fmt.Fprintf(os.Stderr, "\nTip: Check that the file path is correct and the file exists.\n")
	} else if strings.Contains(err.Error(), "permission denied") {
		fmt.Fprintf(os.Stderr, "\nTip: Check that you have read permissions for the file.\n")
	} else if strings.Contains(err.Error(), "invalid range") {
		fmt.Fprintf(os.Stderr, "\nTip: Use format 'file.txt:L10-20' for line ranges.\n")
	}

	os.Exit(code)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}