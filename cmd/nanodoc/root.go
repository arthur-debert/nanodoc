package main

import (
	"fmt"

	"github.com/arthur-debert/nanodoc-go/pkg/nanodoc"
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
	dryRun             bool
)

var rootCmd = &cobra.Command{
	Use:   "nanodoc [paths...]",
	Short: "A minimalist document bundler",
	Long: `Nanodoc is a minimalist document bundler designed for stitching hints, reminders and short docs.
Useful for prompts, personalized docs highlights for your teams or a note to your future self.

No config, nothing to learn nor remember. Short, simple, sweet.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Resolve Paths
		pathInfos, err := nanodoc.ResolvePaths(args)
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
		}

		// 3. Build Document
		doc, err := nanodoc.BuildDocument(pathInfos, opts)
		if err != nil {
			return fmt.Errorf("error building document: %w", err)
		}

		// 4. Create Formatting Context
		ctx, err := nanodoc.NewFormattingContext(opts)
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Line numbering flags
	rootCmd.Flags().BoolVarP(&lineNumbers, "line-numbers", "n", false, "Enable per-file line numbering")
	rootCmd.Flags().BoolVarP(&globalLineNumbers, "global-line-numbers", "N", false, "Enable global line numbering")
	rootCmd.MarkFlagsMutuallyExclusive("line-numbers", "global-line-numbers")

	// TOC flag
	rootCmd.Flags().BoolVar(&toc, "toc", false, "Generate a table of contents")

	// Theme flag
	rootCmd.Flags().StringVar(&theme, "theme", "classic", "Set the theme for formatting (e.g., classic, classic-dark)")

	// Header flags
	rootCmd.Flags().BoolVar(&noHeader, "no-header", false, "Suppress file headers")
	rootCmd.Flags().StringVar(&headerStyle, "header-style", "nice", "Set the header style (nice, filename, path)")
	rootCmd.Flags().StringVar(&sequence, "sequence", "numerical", "Set the sequence style for headers (numerical, letter, roman)")

	// Other flags
	rootCmd.Flags().StringSliceVar(&additionalExt, "txt-ext", []string{}, "Additional file extensions to treat as text (e.g., .log,.conf)")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what files would be processed without actually processing them")

	// Add version command
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of nanodoc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("nanodoc version %s (commit: %s, built: %s)\n", version, commit, date)
	},
} 