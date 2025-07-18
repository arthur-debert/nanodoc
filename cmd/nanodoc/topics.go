package main

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed docs
var docsFS embed.FS

var topicsCmd = &cobra.Command{
	Use:   "topics [topic-name]",
	Short: TopicsShort,
	Long:  TopicsLong,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listTopics(cmd)
		}
		return showTopic(cmd, args[0])
	},
}

func init() {
	rootCmd.AddCommand(topicsCmd)
}

// listTopics lists all available documentation topics
func listTopics(cmd *cobra.Command) error {
	topics, err := getAvailableTopics()
	if err != nil {
		return fmt.Errorf(ErrFailedToGetTopics, err)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), AvailableTopics)
	_, _ = fmt.Fprintln(cmd.OutOrStdout())
	
	// Group topics by directory
	rootTopics := []string{}
	groupedTopics := make(map[string][]string)
	
	for _, topic := range topics {
		if strings.Contains(topic, "/") {
			parts := strings.Split(topic, "/")
			group := parts[0]
			name := parts[1]
			groupedTopics[group] = append(groupedTopics[group], name)
		} else {
			rootTopics = append(rootTopics, topic)
		}
	}
	
	// Print root topics first
	if len(rootTopics) > 0 {
		for _, topic := range rootTopics {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", topic)
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout())
	}
	
	// Print grouped topics
	groups := make([]string, 0, len(groupedTopics))
	for group := range groupedTopics {
		groups = append(groups, group)
	}
	sort.Strings(groups)
	
	for _, group := range groups {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  %s:\n", group)
		sort.Strings(groupedTopics[group])
		for _, topic := range groupedTopics[group] {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "    %s\n", topic)
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout())
	}
	
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), `Run "nanodoc help <topic>" or "nanodoc topics <topic>" for more information.`)
	return nil
}

// showTopic displays the content of a specific topic
func showTopic(cmd *cobra.Command, topicName string) error {
	// Try to find the topic file
	content, err := findAndReadTopic(topicName)
	if err != nil {
		return fmt.Errorf(ErrTopicNotFound, topicName)
	}
	
	_, _ = fmt.Fprint(cmd.OutOrStdout(), content)
	return nil
}

// getAvailableTopics returns a sorted list of all available topics
func getAvailableTopics() ([]string, error) {
	topics := []string{}
	
	err := fs.WalkDir(docsFS, "docs", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if !d.IsDir() && strings.HasSuffix(path, ".txt") {
			// Skip example files and internal files
			if strings.Contains(path, "/examples/") || strings.Contains(path, "/internal/") {
				return nil
			}
			
			// Convert path to topic name
			topic := strings.TrimPrefix(path, "docs/")
			topic = strings.TrimSuffix(topic, ".txt")
			topics = append(topics, topic)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	sort.Strings(topics)
	return topics, nil
}

// findAndReadTopic attempts to find and read a topic file
func findAndReadTopic(topicName string) (string, error) {
	// Clean the topic name
	topicName = strings.ToLower(strings.ReplaceAll(topicName, "-", "_"))
	
	// Try direct match first
	possiblePaths := []string{
		fmt.Sprintf("docs/%s.txt", topicName),
		fmt.Sprintf("docs/%s/%s.txt", path.Dir(topicName), path.Base(topicName)),
	}
	
	// Also try with hyphens converted to underscores
	if strings.Contains(topicName, "_") {
		hyphenated := strings.ReplaceAll(topicName, "_", "-")
		possiblePaths = append(possiblePaths, 
			fmt.Sprintf("docs/%s.txt", hyphenated),
			fmt.Sprintf("docs/%s/%s.txt", path.Dir(hyphenated), path.Base(hyphenated)),
		)
	}
	
	for _, p := range possiblePaths {
		content, err := docsFS.ReadFile(p)
		if err == nil {
			return string(content), nil
		}
	}
	
	return "", fmt.Errorf("%s", TopicNotFoundMsg)
}