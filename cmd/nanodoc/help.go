package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// helpTemplate is the custom template for the root command help
const helpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

// usageTemplate is the custom template for usage
const usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags | groupedFlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

// trimTrailingWhitespaces removes trailing whitespace
func trimTrailingWhitespaces(s string) string {
	return strings.TrimRightFunc(s, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
}

// initHelpSystem sets up the custom help system
func initHelpSystem() {
	// Add custom template functions
	cobra.AddTemplateFunc("trimTrailingWhitespaces", trimTrailingWhitespaces)
	cobra.AddTemplateFunc("groupedFlagUsages", groupedFlagUsages)
	
	// Set custom templates
	rootCmd.SetHelpTemplate(helpTemplate)
	rootCmd.SetUsageTemplate(usageTemplate)
}

// groupedFlagUsages returns flags grouped by their annotations
func groupedFlagUsages(fs *pflag.FlagSet) string {
	if fs == nil {
		return ""
	}
	
	// Group flags by their group annotation
	groups := make(map[string][]*pflag.Flag)
	var groupOrder []string
	
	fs.VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}
		
		// Get group from annotation, default to "MISC"
		group := "MISC"
		if ann := flag.Annotations["group"]; len(ann) > 0 {
			group = ann[0]
		}
		
		if _, exists := groups[group]; !exists {
			groupOrder = append(groupOrder, group)
		}
		groups[group] = append(groups[group], flag)
	})
	
	// Sort group order, ensuring MISC is last
	sort.Slice(groupOrder, func(i, j int) bool {
		if groupOrder[i] == "MISC" {
			return false
		}
		if groupOrder[j] == "MISC" {
			return true
		}
		return groupOrder[i] < groupOrder[j]
	})
	
	// Build output
	var buf bytes.Buffer
	for i, groupName := range groupOrder {
		if i > 0 {
			buf.WriteString("\n")
		}
		
		// Write group header
		buf.WriteString(fmt.Sprintf("%s:\n", groupName))
		
		// Write flags in this group
		for _, flag := range groups[groupName] {
			buf.WriteString(fmt.Sprintf("  %s\n", flagUsage(flag)))
		}
	}
	
	return buf.String()
}

// flagUsage returns the usage string for a single flag
func flagUsage(f *pflag.Flag) string {
	var buf bytes.Buffer
	
	// Build flag name part
	if f.Shorthand != "" && f.ShorthandDeprecated == "" {
		buf.WriteString(fmt.Sprintf("-%s, --%s", f.Shorthand, f.Name))
	} else {
		buf.WriteString(fmt.Sprintf("    --%s", f.Name))
	}
	
	// Add type if not bool
	if f.Value.Type() != "bool" {
		buf.WriteString(fmt.Sprintf(" %s", f.Value.Type()))
	}
	
	// Pad to align descriptions
	padding := 30 - buf.Len()
	if padding > 0 {
		buf.WriteString(strings.Repeat(" ", padding))
	} else {
		buf.WriteString("  ")
	}
	
	// Add usage text
	buf.WriteString(f.Usage)
	
	// Add default value if not empty
	if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "[]" {
		buf.WriteString(fmt.Sprintf(" (default %q)", f.DefValue))
	}
	
	return buf.String()
}