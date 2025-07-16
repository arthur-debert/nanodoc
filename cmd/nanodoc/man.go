package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

var manCmd = &cobra.Command{
	Use:   "man",
	Short: "Generate man page",
	Long:  `Generate a man page for nanodoc`,
	Hidden: true, // Hide from help as it's mainly for build process
	RunE: func(cmd *cobra.Command, args []string) error {
		header := &doc.GenManHeader{
			Title:   "NANODOC",
			Section: "1",
			Source:  "Nanodoc " + version,
			Manual:  "Nanodoc Manual",
		}

		// Generate man page to stdout
		err := doc.GenMan(rootCmd, header, os.Stdout)
		if err != nil {
			return fmt.Errorf("failed to generate man page: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(manCmd)
}