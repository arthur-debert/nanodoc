package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/arthur-debert/nanodoc/pkg/nanodoc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	// Flags
	lineNum            string
	toc                bool
	theme              string
	showFilenames      bool
	fileNumbering      string
	filenameFormat     string  // renamed from headerFormat
	filenameAlign      string  // renamed from headerAlign
	filenameBanner     string  // renamed from headerStyle
	pageWidth          int
	additionalExt      []string
	includePatterns    []string
	excludePatterns    []string
	dryRun             bool
	saveToBundlePath   string
	explicitFlags      map[string]bool

	// Version information - set by ldflags during build
	version = "dev"     // Set by goreleaser: -X main.version={{.Version}}
	commit  = "unknown" // Set by goreleaser: -X main.commit={{.Commit}}
	date    = "unknown" // Set by goreleaser: -X main.date={{.Date}}
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
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "Missing paths to bundle: $ nanodoc <path...>")
			_, _ = fmt.Fprintln(cmd.ErrOrStderr())
			cmd.SilenceUsage = false
			return fmt.Errorf("")
		}
		// Track explicitly set flags
		explicitFlags = nanodoc.TrackExplicitFlags(cmd)

		// 1. Set up Formatting Options first
		lineNumberMode := nanodoc.LineNumberNone
		switch lineNum {
		case "file":
			lineNumberMode = nanodoc.LineNumberFile
		case "global":
			lineNumberMode = nanodoc.LineNumberGlobal
		case "":
			// Default is none
		default:
			return fmt.Errorf("invalid --linenum value: %s (must be 'file' or 'global')", lineNum)
		}

		opts := nanodoc.FormattingOptions{
			LineNumbers:   lineNumberMode,
			ShowTOC:       toc,
			Theme:         theme,
			ShowFilenames:   showFilenames,
			SequenceStyle: nanodoc.SequenceStyle(fileNumbering),
			HeaderFormat:   nanodoc.HeaderFormat(filenameFormat),
			HeaderAlignment: filenameAlign,
			HeaderStyle:     filenameBanner,
			PageWidth:       pageWidth,
			AdditionalExtensions: additionalExt,
			IncludePatterns: includePatterns,
			ExcludePatterns: excludePatterns,
		}

		// 2. Resolve Paths with pattern options
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
			dryRunInfo, err := nanodoc.GenerateDryRunInfo(pathInfos, opts)
			if err != nil {
				return fmt.Errorf(ErrGeneratingDryRun, err)
			}
			
			output := nanodoc.FormatDryRunOutput(dryRunInfo)
			_, _ = fmt.Fprint(cmd.OutOrStdout(), output)
			return nil
		}

		// 3. Extract bundle option lines and merge with command options
		bundleOptionLines, err := nanodoc.ExtractBundleOptionLines(pathInfos)
		if err != nil {
			return fmt.Errorf("error extracting bundle options: %w", err)
		}
		
		// Parse bundle options using Cobra if there are any
		mergedOpts := opts
		if len(bundleOptionLines) > 0 {
			bundleOpts, err := nanodoc.ParseBundleOptions(bundleOptionLines)
			if err != nil {
				return fmt.Errorf("error parsing bundle options: %w", err)
			}
			// Merge options - command line takes precedence
			mergedOpts = nanodoc.MergeOptionsWithExplicitFlags(bundleOpts, opts, explicitFlags)
		}
		
		// 4. Build Document with merged options
		doc, err := nanodoc.BuildDocument(pathInfos, mergedOpts)
		if err != nil {
			return fmt.Errorf(ErrBuildingDocument, err)
		}

		// 5. Create Formatting Context
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

		// 7. Save to bundle if requested
		if saveToBundlePath != "" {
			if err := saveBundleFile(saveToBundlePath, args, opts, cmd); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "\n\nBundle saved to %s\n", saveToBundlePath)
		}

		return nil
	},
}


// saveBundleFile saves the current invocation as a bundle file
func saveBundleFile(path string, args []string, opts nanodoc.FormattingOptions, cmd *cobra.Command) error {
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("bundle file already exists: %s", path)
	}

	// Create the bundle content
	var content strings.Builder
	content.WriteString("# Bundle generated by nanodoc\n")
	content.WriteString("# Command: nanodoc " + reconstructCommand(cmd, args) + "\n\n")

	// Write options section
	content.WriteString("# --- Options ---\n")

	// Always write all options including defaults
	if opts.ShowTOC {
		content.WriteString("--toc\n")
	}

	// Line numbering
	switch opts.LineNumbers {
	case nanodoc.LineNumberFile:
		content.WriteString("--linenum=file\n")
	case nanodoc.LineNumberGlobal:
		content.WriteString("--linenum=global\n")
	}

	// Theme
	content.WriteString(fmt.Sprintf("--theme=%s\n", opts.Theme))

	// File filenames
	if !opts.ShowFilenames {
		content.WriteString("--filenames=false\n")
	}

	// File header format
	content.WriteString(fmt.Sprintf("--header-format=%s\n", string(opts.HeaderFormat)))
	content.WriteString(fmt.Sprintf("--header-align=%s\n", opts.HeaderAlignment))
	content.WriteString(fmt.Sprintf("--header-style=%s\n", opts.HeaderStyle))
	content.WriteString(fmt.Sprintf("--page-width=%d\n", opts.PageWidth))

	// File numbering
	content.WriteString(fmt.Sprintf("--file-numbering=%s\n", string(opts.SequenceStyle)))

	// Additional extensions
	for _, ext := range opts.AdditionalExtensions {
		content.WriteString(fmt.Sprintf("--ext=%s\n", ext))
	}

	// Include patterns
	for _, pattern := range opts.IncludePatterns {
		content.WriteString(fmt.Sprintf("--include=%q\n", pattern))
	}

	// Exclude patterns
	for _, pattern := range opts.ExcludePatterns {
		content.WriteString(fmt.Sprintf("--exclude=%q\n", pattern))
	}

	// Write content section
	content.WriteString("\n# --- Content ---\n")
	for _, arg := range args {
		content.WriteString(arg + "\n")
	}

	// Write to file
	return os.WriteFile(path, []byte(content.String()), 0644)
}



// reconstructCommand reconstructs the command-line invocation from cobra flags and args
func reconstructCommand(cmd *cobra.Command, args []string) string {
	var parts []string
	cmd.Flags().Visit(func(f *pflag.Flag) {
		if f.Name != "save-to-bundle" {
			// Handle boolean flags specially
			if f.Value.Type() == "bool" {
				if f.Value.String() == "true" {
					parts = append(parts, "--"+f.Name)
				} else {
					parts = append(parts, fmt.Sprintf("--%s=false", f.Name))
				}
			} else if f.Value.String() != "" {
				// For non-boolean flags, check if value contains spaces or special chars
				val := f.Value.String()
				if strings.ContainsAny(val, " \t*?") {
					parts = append(parts, fmt.Sprintf("--%s=%q", f.Name, val))
				} else {
					parts = append(parts, fmt.Sprintf("--%s=%s", f.Name, val))
				}
			} else {
				parts = append(parts, "--"+f.Name)
			}
		}
	})
	parts = append(parts, args...)
	return strings.Join(parts, " ")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Line numbering flag
	rootCmd.Flags().StringVarP(&lineNum, "linenum", "l", "", FlagLineNum)
	_ = rootCmd.RegisterFlagCompletionFunc("linenum", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"file", "global"}, cobra.ShellCompDirectiveNoFileComp
	})
	_ = rootCmd.Flags().SetAnnotation("linenum", "group", []string{"Formatting"})

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

	// File name flags
	rootCmd.Flags().BoolVar(&showFilenames, "filenames", true, FlagFilenames)
	rootCmd.Flags().StringVar(&filenameFormat, "header-format", "nice", FlagHeaderFormat)
	_ = rootCmd.RegisterFlagCompletionFunc("header-format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"nice", "simple", "path", "filename", "title"}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.Flags().StringVar(&filenameAlign, "header-align", "left", "Header alignment (left, center, right)")
	_ = rootCmd.RegisterFlagCompletionFunc("header-align", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"left", "center", "right"}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.Flags().StringVar(&filenameBanner, "header-style", "none", "Header style")
	_ = rootCmd.RegisterFlagCompletionFunc("header-style", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Dynamically get banner styles from registry
		return nanodoc.GetBannerStyleNames(), cobra.ShellCompDirectiveNoFileComp
	})
	// Auto-detect terminal width as default for page width
	defaultPageWidth := nanodoc.GetTerminalWidth()
	rootCmd.Flags().IntVar(&pageWidth, "page-width", defaultPageWidth, "Page width for alignment (auto-detected)")
	rootCmd.Flags().StringVar(&fileNumbering, "file-numbering", "numerical", FlagFileNumbering)
	_ = rootCmd.RegisterFlagCompletionFunc("file-numbering", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"numerical", "alphabetical", "roman"}, cobra.ShellCompDirectiveNoFileComp
	})
	_ = rootCmd.Flags().SetAnnotation("filenames", "group", []string{"Features"})
	_ = rootCmd.Flags().SetAnnotation("header-format", "group", []string{"Formatting"})
	_ = rootCmd.Flags().SetAnnotation("header-align", "group", []string{"Formatting"})
	_ = rootCmd.Flags().SetAnnotation("header-style", "group", []string{"Formatting"})
	_ = rootCmd.Flags().SetAnnotation("page-width", "group", []string{"Formatting"})
	_ = rootCmd.Flags().SetAnnotation("file-numbering", "group", []string{"Features"})

	// File filtering flags
	rootCmd.Flags().StringSliceVar(&additionalExt, "ext", []string{}, FlagExt)
	rootCmd.Flags().StringSliceVar(&includePatterns, "include", []string{}, FlagInclude)
	rootCmd.Flags().StringSliceVar(&excludePatterns, "exclude", []string{}, FlagExclude)
	_ = rootCmd.Flags().SetAnnotation("ext", "group", []string{"File Selection"})
	_ = rootCmd.Flags().SetAnnotation("include", "group", []string{"File Selection"})
	_ = rootCmd.Flags().SetAnnotation("exclude", "group", []string{"File Selection"})
	
	// Other flags
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, FlagDryRun)
	rootCmd.Flags().StringVar(&saveToBundlePath, "save-to-bundle", "", "Save the current invocation as a bundle file")
	rootCmd.Flags().BoolP("version", "v", false, FlagVersion)
	_ = rootCmd.Flags().SetAnnotation("dry-run", "group", []string{"Misc"})
	_ = rootCmd.Flags().SetAnnotation("save-to-bundle", "group", []string{"Features"})
	_ = rootCmd.Flags().SetAnnotation("version", "group", []string{"Misc"})
	_ = rootCmd.Flags().SetAnnotation("help", "group", []string{"Misc"})
	
	// Initialize custom help system
	initHelpSystem()
}