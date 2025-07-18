package main

import (
	_ "embed"
	"fmt"

	"github.com/arthur-debert/nanodoc/pkg/nanodoc"
	"github.com/spf13/cobra"
)

var (
	// Flags
	lineNumbers        bool
	globalLineNumbers  bool
	toc                bool
	theme              string
	noHeader           bool
	sequence           string
	headerStyle        string
	additionalExt      []string
	includePatterns    []string
	excludePatterns    []string
	dryRun             bool
	explicitFlags      map[string]bool
	
	// Version information - set by ldflags during build
	version = "dev"      // Set by goreleaser: -X main.version={{.Version}}
	commit  = "unknown"  // Set by goreleaser: -X main.commit={{.Commit}}
	date    = "unknown"  // Set by goreleaser: -X main.date={{.Date}}
)

//go:embed help/root-long.txt
var rootLongHelp string

//go:embed help/root-examples.txt
var rootExamples string

var rootCmd = &cobra.Command{
	Use:     "nanodoc [paths...]",
	Short:   "A minimalist document bundler",
	Long:    rootLongHelp,
	Example: rootExamples,
	Args:    cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check version flag first
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Printf("nanodoc version %s (commit: %s, built: %s)\n", version, commit, date)
			return nil
		}
		
		// Check args only if not printing version
		if len(args) < 1 {
			return fmt.Errorf("requires at least 1 arg(s), only received %d", len(args))
		}
		// Track explicitly set flags
		trackExplicitFlags(cmd)

		// 1. Resolve Paths with pattern options
		pathOpts := &nanodoc.FormattingOptions{
			AdditionalExtensions: additionalExt,
			IncludePatterns: includePatterns,
			ExcludePatterns: excludePatterns,
		}
		pathInfos, err := nanodoc.ResolvePathsWithOptions(args, pathOpts)
		if err != nil {
			return fmt.Errorf("error resolving paths: %w", err)
		}

		// If dry run, show what would be processed and exit
		if dryRun {
			dryRunInfo, err := nanodoc.GenerateDryRunInfo(pathInfos, additionalExt)
			if err != nil {
				return fmt.Errorf("error generating dry run info: %w", err)
			}
			
			output := nanodoc.FormatDryRunOutput(dryRunInfo)
			_, _ = fmt.Fprint(cmd.OutOrStdout(), output)
			return nil
		}

		// 2. Set up Formatting Options
		lineNumberMode := nanodoc.LineNumberNone
		if globalLineNumbers {
			lineNumberMode = nanodoc.LineNumberGlobal
		} else if lineNumbers {
			lineNumberMode = nanodoc.LineNumberFile
		}

		opts := nanodoc.FormattingOptions{
			LineNumbers:   lineNumberMode,
			ShowTOC:       toc,
			Theme:         theme,
			ShowHeaders:   !noHeader,
			SequenceStyle: nanodoc.SequenceStyle(sequence),
			HeaderStyle:   nanodoc.HeaderStyle(headerStyle),
			AdditionalExtensions: additionalExt,
			IncludePatterns: includePatterns,
			ExcludePatterns: excludePatterns,
		}

		// 3. Build Document with explicit flags
		doc, err := nanodoc.BuildDocumentWithExplicitFlags(pathInfos, opts, explicitFlags)
		if err != nil {
			return fmt.Errorf("error building document: %w", err)
		}

		// 4. Create Formatting Context
		ctx, err := nanodoc.NewFormattingContext(doc.FormattingOptions)
		if err != nil {
			return fmt.Errorf("error creating formatting context: %w", err)
		}

		// 5. Render Document
		output, err := nanodoc.RenderDocument(doc, ctx)
		if err != nil {
			return fmt.Errorf("error rendering document: %w", err)
		}

		// 6. Print to stdout
		_, _ = fmt.Fprint(cmd.OutOrStdout(), output)

		return nil
	},
}

// trackExplicitFlags determines which flags were explicitly set by the user
func trackExplicitFlags(cmd *cobra.Command) {
	explicitFlags = make(map[string]bool)
	
	// Check if each flag was explicitly set
	if cmd.Flags().Changed("toc") {
		explicitFlags["toc"] = true
	}
	if cmd.Flags().Changed("theme") {
		explicitFlags["theme"] = true
	}
	if cmd.Flags().Changed("line-numbers") || cmd.Flags().Changed("global-line-numbers") {
		explicitFlags["line-numbers"] = true
	}
	if cmd.Flags().Changed("no-header") {
		explicitFlags["no-header"] = true
	}
	if cmd.Flags().Changed("header-style") {
		explicitFlags["header-style"] = true
	}
	if cmd.Flags().Changed("sequence") {
		explicitFlags["sequence"] = true
	}
	if cmd.Flags().Changed("txt-ext") {
		explicitFlags["txt-ext"] = true
	}
	if cmd.Flags().Changed("include") {
		explicitFlags["include"] = true
	}
	if cmd.Flags().Changed("exclude") {
		explicitFlags["exclude"] = true
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Line numbering flags
	rootCmd.Flags().BoolVarP(&lineNumbers, "line-numbers", "n", false, "Enable per-file line numbering (see: nanodoc topics line-numbering)")
	rootCmd.Flags().BoolVarP(&globalLineNumbers, "global-line-numbers", "N", false, "Enable global line numbering (see: nanodoc topics line-numbering)")
	rootCmd.MarkFlagsMutuallyExclusive("line-numbers", "global-line-numbers")
	_ = rootCmd.Flags().SetAnnotation("line-numbers", "group", []string{"Formatting"})
	_ = rootCmd.Flags().SetAnnotation("global-line-numbers", "group", []string{"Formatting"})

	// TOC flag
	rootCmd.Flags().BoolVar(&toc, "toc", false, "Generate a table of contents (see: nanodoc topics table-of-contents)")
	_ = rootCmd.Flags().SetAnnotation("toc", "group", []string{"Features"})

	// Theme flag
	rootCmd.Flags().StringVar(&theme, "theme", "classic", "Set the theme for formatting (see: nanodoc topics themes)")
	_ = rootCmd.Flags().SetAnnotation("theme", "group", []string{"Formatting"})

	// Header flags
	rootCmd.Flags().BoolVar(&noHeader, "no-header", false, "Suppress file headers")
	rootCmd.Flags().StringVar(&headerStyle, "header-style", "nice", "Set the header style (see: nanodoc topics headers-and-sequencing)")
	rootCmd.Flags().StringVar(&sequence, "sequence", "numerical", "Set the sequence style (see: nanodoc topics headers-and-sequencing)")
	_ = rootCmd.Flags().SetAnnotation("no-header", "group", []string{"Features"})
	_ = rootCmd.Flags().SetAnnotation("header-style", "group", []string{"Formatting"})
	_ = rootCmd.Flags().SetAnnotation("sequence", "group", []string{"Features"})

	// File filtering flags
	rootCmd.Flags().StringSliceVar(&additionalExt, "txt-ext", []string{}, "Additional file extensions to treat as text")
	rootCmd.Flags().StringSliceVar(&includePatterns, "include", []string{}, "Include only files matching patterns (see: nanodoc topics specifying-files)")
	rootCmd.Flags().StringSliceVar(&excludePatterns, "exclude", []string{}, "Exclude files matching patterns (see: nanodoc topics specifying-files)")
	_ = rootCmd.Flags().SetAnnotation("txt-ext", "group", []string{"File Selection"})
	_ = rootCmd.Flags().SetAnnotation("include", "group", []string{"File Selection"})
	_ = rootCmd.Flags().SetAnnotation("exclude", "group", []string{"File Selection"})
	
	// Other flags
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what files would be processed without actually processing them")
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
	_ = rootCmd.Flags().SetAnnotation("dry-run", "group", []string{"Misc"})
	_ = rootCmd.Flags().SetAnnotation("version", "group", []string{"Misc"})
	_ = rootCmd.Flags().SetAnnotation("help", "group", []string{"Misc"})
	
	// Initialize custom help system
	initHelpSystem()
} 