package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help [topic]",
	Short: "Help about any command or topic",
	Long: `Help provides help for any command in the application.
Simply type nanodoc help [path to command] for full details.
You can also type nanodoc help [topic] to read about a specific topic.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// No topic provided, show root help
			_ = rootCmd.Help()
			return
		}

		// Check if it's a topic
		topic := args[0]
		content, err := findAndReadTopic(topic)
		if err == nil {
			// It's a valid topic, show it
			_, _ = fmt.Fprint(cmd.OutOrStdout(), content)
			return
		}

		// Not a topic, try to find it as a command
		cmd, _, err = rootCmd.Find(args)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Unknown help topic %#q\n", args[0])
			_ = rootCmd.Usage()
			return
		}

		// Show help for the found command
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(helpCmd)
	// Disable the default help command
	rootCmd.SetHelpCommand(helpCmd)
}