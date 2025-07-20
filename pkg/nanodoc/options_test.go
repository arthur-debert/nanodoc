package nanodoc

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestTrackExplicitFlags(t *testing.T) {
	tests := []struct {
		name          string
		setupFlags    func(*cobra.Command)
		expectedFlags map[string]bool
	}{
		{
			name: "no_flags_set",
			setupFlags: func(cmd *cobra.Command) {
				// Don't mark any flags as changed
			},
			expectedFlags: map[string]bool{},
		},
		{
			name: "some_flags_set",
			setupFlags: func(cmd *cobra.Command) {
				// Mark specific flags as changed
				_ = cmd.Flags().Set("toc", "true")
				_ = cmd.Flags().Set("theme", "dark")
			},
			expectedFlags: map[string]bool{
				"toc":   true,
				"theme": true,
			},
		},
		{
			name: "all_tracked_flags_set",
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().Set("toc", "true")
				_ = cmd.Flags().Set("theme", "dark")
				_ = cmd.Flags().Set("linenum", "global")
				_ = cmd.Flags().Set("filenames", "false")
				_ = cmd.Flags().Set("header-format", "path")
			},
			expectedFlags: map[string]bool{
				"toc":           true,
				"theme":         true,
				"line-numbers":  true,
				"no-header":     true,
				"header-format": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test command with necessary flags
			cmd := &cobra.Command{}
			cmd.Flags().Bool("toc", false, "")
			cmd.Flags().String("theme", "classic", "")
			cmd.Flags().String("linenum", "", "")
			cmd.Flags().Bool("filenames", true, "")
			cmd.Flags().String("header-format", "nice", "")
			cmd.Flags().String("header-align", "left", "")
			cmd.Flags().String("header-style", "none", "")
			cmd.Flags().Int("page-width", 80, "")
			cmd.Flags().String("file-numbering", "numerical", "")
			cmd.Flags().StringSlice("ext", []string{}, "")
			cmd.Flags().StringSlice("include", []string{}, "")
			cmd.Flags().StringSlice("exclude", []string{}, "")

			// Setup the flags as per test case
			tt.setupFlags(cmd)

			// Track explicit flags
			result := TrackExplicitFlags(cmd)

			// Compare lengths first
			if len(result) != len(tt.expectedFlags) {
				t.Errorf("Expected %d explicit flags, got %d", len(tt.expectedFlags), len(result))
			}

			// Check each expected flag
			for flag, expected := range tt.expectedFlags {
				if actual, exists := result[flag]; !exists || actual != expected {
					t.Errorf("Expected flag %q to be %v, got %v (exists: %v)", flag, expected, actual, exists)
				}
			}
		})
	}
}

func TestBuildFormattingOptions(t *testing.T) {
	tests := []struct {
		name            string
		lineNum         string
		toc             bool
		theme           string
		showFilenames   bool
		fileNumbering   string
		filenameFormat  string
		filenameAlign   string
		filenameBanner  string
		pageWidth       int
		additionalExt   []string
		includePatterns []string
		excludePatterns []string
		outputFormat    string
		wantOpts        FormattingOptions
		wantErr         bool
	}{
		{
			name:           "default_options",
			lineNum:        "",
			toc:            false,
			theme:          "classic",
			showFilenames:  true,
			fileNumbering:  "numerical",
			filenameFormat: "nice",
			filenameAlign:  "left",
			filenameBanner: "none",
			pageWidth:      80,
			outputFormat:   "term",
			wantOpts: FormattingOptions{
				LineNumbers:     LineNumberNone,
				ShowTOC:         false,
				Theme:           "classic",
				ShowFilenames:   true,
				SequenceStyle:   "numerical",
				HeaderFormat:    "nice",
				HeaderAlignment: "left",
				HeaderStyle:     "none",
				PageWidth:       80,
				OutputFormat:    "term",
			},
			wantErr: false,
		},
		{
			name:           "file_line_numbers",
			lineNum:        "file",
			toc:            true,
			theme:          "dark",
			showFilenames:  false,
			fileNumbering:  "alphabetical",
			filenameFormat: "path",
			filenameAlign:  "center",
			filenameBanner: "dashed",
			pageWidth:      120,
			outputFormat:   "term",
			wantOpts: FormattingOptions{
				LineNumbers:     LineNumberFile,
				ShowTOC:         true,
				Theme:           "dark",
				ShowFilenames:   false,
				SequenceStyle:   "alphabetical",
				HeaderFormat:    "path",
				HeaderAlignment: "center",
				HeaderStyle:     "dashed",
				PageWidth:       120,
				OutputFormat:    "term",
			},
			wantErr: false,
		},
		{
			name:           "global_line_numbers",
			lineNum:        "global",
			toc:            false,
			theme:          "classic",
			showFilenames:  true,
			fileNumbering:  "roman",
			filenameFormat: "filename",
			filenameAlign:  "right",
			filenameBanner: "boxed",
			pageWidth:      100,
			outputFormat:   "term",
			wantOpts: FormattingOptions{
				LineNumbers:     LineNumberGlobal,
				ShowTOC:         false,
				Theme:           "classic",
				ShowFilenames:   true,
				SequenceStyle:   "roman",
				HeaderFormat:    "filename",
				HeaderAlignment: "right",
				HeaderStyle:     "boxed",
				PageWidth:       100,
				OutputFormat:    "term",
			},
			wantErr: false,
		},
		{
			name:           "with_patterns",
			lineNum:        "",
			toc:            false,
			theme:          "classic",
			showFilenames:  true,
			fileNumbering:  "numerical",
			filenameFormat: "nice",
			filenameAlign:  "left",
			filenameBanner: "none",
			pageWidth:      80,
			additionalExt:  []string{".txt", ".log"},
			outputFormat:   "term",
			includePatterns: []string{"*.go", "*.md"},
			excludePatterns: []string{"*_test.go", "vendor/*"},
			wantOpts: FormattingOptions{
				LineNumbers:          LineNumberNone,
				ShowTOC:              false,
				Theme:                "classic",
				ShowFilenames:        true,
				SequenceStyle:        "numerical",
				HeaderFormat:         "nice",
				HeaderAlignment:      "left",
				HeaderStyle:          "none",
				PageWidth:            80,
				AdditionalExtensions: []string{".txt", ".log"},
				IncludePatterns:      []string{"*.go", "*.md"},
				ExcludePatterns:      []string{"*_test.go", "vendor/*"},
				OutputFormat:         "term",
			},
			wantErr: false,
		},
		{
			name:    "invalid_line_num",
			lineNum: "invalid",
			wantErr: true,
		},
		{
			name:           "with_output_format_term",
			lineNum:        "",
			toc:            false,
			theme:          "classic",
			showFilenames:  true,
			fileNumbering:  "numerical",
			filenameFormat: "nice",
			filenameAlign:  "left",
			filenameBanner: "none",
			pageWidth:      80,
			outputFormat:   "term",
			wantOpts: FormattingOptions{
				LineNumbers:     LineNumberNone,
				ShowTOC:         false,
				Theme:           "classic",
				ShowFilenames:   true,
				SequenceStyle:   "numerical",
				HeaderFormat:    "nice",
				HeaderAlignment: "left",
				HeaderStyle:     "none",
				PageWidth:       80,
				OutputFormat:    "term",
			},
			wantErr: false,
		},
		{
			name:           "with_output_format_markdown",
			lineNum:        "",
			toc:            false,
			theme:          "classic",
			showFilenames:  true,
			fileNumbering:  "numerical",
			filenameFormat: "nice",
			filenameAlign:  "left",
			filenameBanner: "none",
			pageWidth:      80,
			outputFormat:   "markdown",
			wantOpts: FormattingOptions{
				LineNumbers:     LineNumberNone,
				ShowTOC:         false,
				Theme:           "classic",
				ShowFilenames:   true,
				SequenceStyle:   "numerical",
				HeaderFormat:    "nice",
				HeaderAlignment: "left",
				HeaderStyle:     "none",
				PageWidth:       80,
				OutputFormat:    "markdown",
			},
			wantErr: false,
		},
		{
			name:         "invalid_output_format",
			lineNum:      "",
			outputFormat: "invalid",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := BuildFormattingOptions(
				tt.lineNum,
				tt.toc,
				tt.theme,
				tt.showFilenames,
				tt.fileNumbering,
				tt.filenameFormat,
				tt.filenameAlign,
				tt.filenameBanner,
				tt.pageWidth,
				tt.additionalExt,
				tt.includePatterns,
				tt.excludePatterns,
				tt.outputFormat,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuildFormattingOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Compare the options
				if opts.LineNumbers != tt.wantOpts.LineNumbers {
					t.Errorf("LineNumbers = %v, want %v", opts.LineNumbers, tt.wantOpts.LineNumbers)
				}
				if opts.ShowTOC != tt.wantOpts.ShowTOC {
					t.Errorf("ShowTOC = %v, want %v", opts.ShowTOC, tt.wantOpts.ShowTOC)
				}
				if opts.Theme != tt.wantOpts.Theme {
					t.Errorf("Theme = %v, want %v", opts.Theme, tt.wantOpts.Theme)
				}
				if opts.ShowFilenames != tt.wantOpts.ShowFilenames {
					t.Errorf("ShowFilenames = %v, want %v", opts.ShowFilenames, tt.wantOpts.ShowFilenames)
				}
				if opts.SequenceStyle != tt.wantOpts.SequenceStyle {
					t.Errorf("SequenceStyle = %v, want %v", opts.SequenceStyle, tt.wantOpts.SequenceStyle)
				}
				if opts.HeaderFormat != tt.wantOpts.HeaderFormat {
					t.Errorf("HeaderFormat = %v, want %v", opts.HeaderFormat, tt.wantOpts.HeaderFormat)
				}
				if opts.HeaderAlignment != tt.wantOpts.HeaderAlignment {
					t.Errorf("HeaderAlignment = %v, want %v", opts.HeaderAlignment, tt.wantOpts.HeaderAlignment)
				}
				if opts.HeaderStyle != tt.wantOpts.HeaderStyle {
					t.Errorf("HeaderStyle = %v, want %v", opts.HeaderStyle, tt.wantOpts.HeaderStyle)
				}
				if opts.PageWidth != tt.wantOpts.PageWidth {
					t.Errorf("PageWidth = %v, want %v", opts.PageWidth, tt.wantOpts.PageWidth)
				}
				// Compare slices
				if len(opts.AdditionalExtensions) != len(tt.wantOpts.AdditionalExtensions) {
					t.Errorf("AdditionalExtensions length = %v, want %v", len(opts.AdditionalExtensions), len(tt.wantOpts.AdditionalExtensions))
				}
				if len(opts.IncludePatterns) != len(tt.wantOpts.IncludePatterns) {
					t.Errorf("IncludePatterns length = %v, want %v", len(opts.IncludePatterns), len(tt.wantOpts.IncludePatterns))
				}
				if len(opts.ExcludePatterns) != len(tt.wantOpts.ExcludePatterns) {
					t.Errorf("ExcludePatterns length = %v, want %v", len(opts.ExcludePatterns), len(tt.wantOpts.ExcludePatterns))
				}
				if opts.OutputFormat != tt.wantOpts.OutputFormat {
					t.Errorf("OutputFormat = %v, want %v", opts.OutputFormat, tt.wantOpts.OutputFormat)
				}
			}
		})
	}
}