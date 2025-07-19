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
	headerFormat       string
	headerAlign        string
	headerStyle        string
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
		trackExplicitFlags(cmd)

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
			HeaderFormat:   nanodoc.HeaderFormat(headerFormat),
			HeaderAlignment: headerAlign,
			HeaderStyle:     headerStyle,
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
			bundleOpts, err := parseBundleOptions(bundleOptionLines)
			if err != nil {
				return fmt.Errorf("error parsing bundle options: %w", err)
			}
			// Merge options - command line takes precedence
			mergedOpts = mergeOptionsWithExplicitFlags(bundleOpts, opts, explicitFlags)
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

// parseBundleOptions parses bundle option lines using Cobra
func parseBundleOptions(optionLines []string) (nanodoc.FormattingOptions, error) {
	// Create a temporary command to parse options
	tempCmd := &cobra.Command{}
	
	// Set up the same flags as the root command
	var bundleLineNum string
	var bundleToc bool
	var bundleTheme string
	var bundleShowFilenames bool
	var bundleFileNumbering string
	var bundleHeaderFormat string
	var bundleHeaderAlign string
	var bundleHeaderStyle string
	var bundlePageWidth int
	var bundleAdditionalExt []string
	var bundleIncludePatterns []string
	var bundleExcludePatterns []string
	
	tempCmd.Flags().StringVarP(&bundleLineNum, "linenum", "l", "", "")
	tempCmd.Flags().BoolVar(&bundleToc, "toc", false, "")
	tempCmd.Flags().StringVar(&bundleTheme, "theme", "classic", "")
	tempCmd.Flags().BoolVar(&bundleShowFilenames, "filenames", true, "")
	tempCmd.Flags().StringVar(&bundleHeaderFormat, "header-format", "nice", "")
	tempCmd.Flags().StringVar(&bundleHeaderAlign, "header-align", "left", "")
	tempCmd.Flags().StringVar(&bundleHeaderStyle, "header-style", "none", "")
	tempCmd.Flags().IntVar(&bundlePageWidth, "page-width", nanodoc.OUTPUT_WIDTH, "")
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
		return nanodoc.FormattingOptions{}, err
	}
	
	// Convert to FormattingOptions
	lineNumberMode := nanodoc.LineNumberNone
	switch bundleLineNum {
	case "file":
		lineNumberMode = nanodoc.LineNumberFile
	case "global":
		lineNumberMode = nanodoc.LineNumberGlobal
	}
	
	return nanodoc.FormattingOptions{
		LineNumbers:          lineNumberMode,
		ShowTOC:              bundleToc,
		Theme:                bundleTheme,
		ShowFilenames:          bundleShowFilenames,
		SequenceStyle:        nanodoc.SequenceStyle(bundleFileNumbering),
		HeaderFormat:          nanodoc.HeaderFormat(bundleHeaderFormat),
		HeaderAlignment:      bundleHeaderAlign,
		HeaderStyle:          bundleHeaderStyle,
		PageWidth:            bundlePageWidth,
		AdditionalExtensions: bundleAdditionalExt,
		IncludePatterns:      bundleIncludePatterns,
		ExcludePatterns:      bundleExcludePatterns,
	}, nil
}

// mergeOptionsWithExplicitFlags merges bundle options with command options based on explicit flags
func mergeOptionsWithExplicitFlags(bundleOpts, cmdOpts nanodoc.FormattingOptions, explicitFlags map[string]bool) nanodoc.FormattingOptions {
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
	rootCmd.Flags().StringVar(&headerFormat, "header-format", "nice", FlagHeaderFormat)
	_ = rootCmd.RegisterFlagCompletionFunc("header-format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"nice", "simple", "path", "filename", "title"}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.Flags().StringVar(&headerAlign, "header-align", "left", "Header alignment (left, center, right)")
	_ = rootCmd.RegisterFlagCompletionFunc("header-align", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"left", "center", "right"}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.Flags().StringVar(&headerStyle, "header-style", "none", "Header style (none, dashed, solid, boxed)")
	_ = rootCmd.RegisterFlagCompletionFunc("header-style", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"none", "dashed", "solid", "boxed"}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.Flags().IntVar(&pageWidth, "page-width", nanodoc.OUTPUT_WIDTH, "Page width for alignment")
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