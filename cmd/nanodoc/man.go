package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

var manCmd = &cobra.Command{
	Use:   "man",
	Short: ManShort,
	Long:  ManLong,
	Hidden: true, // Hide from help as it's mainly for build process
	RunE: func(cmd *cobra.Command, args []string) error {
		header := &doc.GenManHeader{
			Title:   ManTitle,
			Section: ManSection,
			Source:  "Nanodoc " + version,
			Manual:  ManManual,
		}

		// Generate man page to stdout
		err := doc.GenMan(rootCmd, header, os.Stdout)
		if err != nil {
			return fmt.Errorf(ErrFailedGenManPage, err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(manCmd)
}