package main

import (
	_ "embed"
	"os"

	"github.com/spf13/cobra"
)

//go:embed completion.txt
var completionLong string

var completionCmd = &cobra.Command{
	Use:    "completion [bash|zsh|fish|powershell]",
	Short:  CompletionShort,
	Hidden: true, // Hide from help output but still accessible for build/packaging
	Long:   completionLong,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}