package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nanodoc",
	Short: "A minimalist document bundler",
	Long: `Nanodoc is a minimalist document bundler designed for stitching hints, reminders and short docs.
Useful for prompts, personalized docs highlights for your teams or a note to your future self.

No config, nothing to learn nor remember. Short, simple, sweet.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.my-awesome-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	
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