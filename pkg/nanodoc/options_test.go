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