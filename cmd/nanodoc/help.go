package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

//go:embed help.txt
var usageTemplate string

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
	rootCmd.SetHelpTemplate(HelpTemplate)
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
		
		// Get group from annotation, default to "Misc"
		group := MiscGroupName
		if ann := flag.Annotations["group"]; len(ann) > 0 {
			group = ann[0]
		}
		
		if _, exists := groups[group]; !exists {
			groupOrder = append(groupOrder, group)
		}
		groups[group] = append(groups[group], flag)
	})
	
	// Sort group order, ensuring Misc is last
	sort.Slice(groupOrder, func(i, j int) bool {
		if groupOrder[i] == MiscGroupName {
			return false
		}
		if groupOrder[j] == MiscGroupName {
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