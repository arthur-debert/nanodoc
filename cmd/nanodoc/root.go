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
	Short:   RootShort,
	Long:    rootLongHelp,
	Example: rootExamples,
	Args:    cobra.ArbitraryArgs,
	SilenceUsage: true,
	SilenceErrors: true,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Default to file completion
		return nil, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check version flag first
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Printf(VersionFormat, version, commit, date)
			return nil
		}
		
		// Check args only if not printing version
		if len(args) < 1 {
			fmt.Fprintln(cmd.ErrOrStderr(), "Missing paths to bundle: $ nanodoc <path...>")
			fmt.Fprintln(cmd.ErrOrStderr())
			cmd.SilenceUsage = false
			return fmt.Errorf("")
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
			return fmt.Errorf(ErrResolvingPaths, err)
		}

		// If dry run, show what would be processed and exit
		if dryRun {
			dryRunInfo, err := nanodoc.GenerateDryRunInfo(pathInfos, additionalExt)
			if err != nil {
				return fmt.Errorf(ErrGeneratingDryRun, err)
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
			return fmt.Errorf(ErrBuildingDocument, err)
		}

		// 4. Create Formatting Context
		ctx, err := nanodoc.NewFormattingContext(doc.FormattingOptions)
		if err != nil {
			return fmt.Errorf(ErrCreatingContext, err)
		}

		// 5. Render Document
		output, err := nanodoc.RenderDocument(doc, ctx)
		if err != nil {
			return fmt.Errorf(ErrRenderingDocument, err)
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
	rootCmd.Flags().BoolVarP(&lineNumbers, "line-numbers", "n", false, FlagLineNumbers)
	rootCmd.Flags().BoolVarP(&globalLineNumbers, "global-line-numbers", "N", false, FlagGlobalLineNumbers)
	rootCmd.MarkFlagsMutuallyExclusive("line-numbers", "global-line-numbers")
	_ = rootCmd.Flags().SetAnnotation("line-numbers", "group", []string{"Formatting"})
	_ = rootCmd.Flags().SetAnnotation("global-line-numbers", "group", []string{"Formatting"})

	// TOC flag
	rootCmd.Flags().BoolVar(&toc, "toc", false, FlagTOC)
	_ = rootCmd.Flags().SetAnnotation("toc", "group", []string{"Features"})

	// Theme flag
	rootCmd.Flags().StringVar(&theme, "theme", "classic", FlagTheme)
	_ = rootCmd.RegisterFlagCompletionFunc("theme", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		themes, err := nanodoc.GetAvailableThemes()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return themes, cobra.ShellCompDirectiveNoFileComp
	})
	_ = rootCmd.Flags().SetAnnotation("theme", "group", []string{"Formatting"})

	// Header flags
	rootCmd.Flags().BoolVar(&noHeader, "no-header", false, FlagNoHeader)
	rootCmd.Flags().StringVar(&headerStyle, "header-style", "nice", FlagHeaderStyle)
	_ = rootCmd.RegisterFlagCompletionFunc("header-style", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"nice", "simple", "path", "filename", "title"}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.Flags().StringVar(&sequence, "sequence", "numerical", FlagSequence)
	_ = rootCmd.RegisterFlagCompletionFunc("sequence", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"numerical", "alphabetical", "roman"}, cobra.ShellCompDirectiveNoFileComp
	})
	_ = rootCmd.Flags().SetAnnotation("no-header", "group", []string{"Features"})
	_ = rootCmd.Flags().SetAnnotation("header-style", "group", []string{"Formatting"})
	_ = rootCmd.Flags().SetAnnotation("sequence", "group", []string{"Features"})

	// File filtering flags
	rootCmd.Flags().StringSliceVar(&additionalExt, "txt-ext", []string{}, FlagTxtExt)
	rootCmd.Flags().StringSliceVar(&includePatterns, "include", []string{}, FlagInclude)
	rootCmd.Flags().StringSliceVar(&excludePatterns, "exclude", []string{}, FlagExclude)
	_ = rootCmd.Flags().SetAnnotation("txt-ext", "group", []string{"File Selection"})
	_ = rootCmd.Flags().SetAnnotation("include", "group", []string{"File Selection"})
	_ = rootCmd.Flags().SetAnnotation("exclude", "group", []string{"File Selection"})
	
	// Other flags
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, FlagDryRun)
	rootCmd.Flags().BoolP("version", "v", false, FlagVersion)
	_ = rootCmd.Flags().SetAnnotation("dry-run", "group", []string{"Misc"})
	_ = rootCmd.Flags().SetAnnotation("version", "group", []string{"Misc"})
	_ = rootCmd.Flags().SetAnnotation("help", "group", []string{"Misc"})
	
	// Initialize custom help system
	initHelpSystem()
} 