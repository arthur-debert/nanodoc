// This file contains the messages for the nanodoc help system. .
package main

// Command descriptions
const (
	RootShort = "A minimalist document bundler"

	TopicsShort = "Display help topics"
	TopicsLong  = `Display available help topics or show the content of a specific topic.

Running 'nanodoc topics' lists all available topics.
Running 'nanodoc topics <topic-name>' displays the content of that topic.`

	CompletionShort = "Generate completion script"

	ManShort = "Generate man page"
	ManLong  = `Generate a man page for nanodoc`
)

// Error messages
const (
	ErrMinArgs           = "requires at least 1 arg(s), only received %d"
	ErrResolvingPaths    = "error resolving paths: %w"
	ErrGeneratingDryRun  = "error generating dry run info: %w"
	ErrBuildingDocument  = "error building document: %w"
	ErrCreatingContext   = "error creating formatting context: %w"
	ErrRenderingDocument = "error rendering document: %w"
	ErrTopicNotFound     = "topic '%s' not found"
	ErrFailedToGetTopics = "failed to get available topics: %w"
	ErrFailedGenManPage  = "failed to generate man page: %w"
)

// Flag descriptions
const (
	FlagLineNum           = "Enable line numbering (file|global) (see: nanodoc topics line-numbering)"
	FlagTOC               = "Generate a table of contents (see: nanodoc topics toc)"
	FlagTheme             = "Set the theme for formatting (see: nanodoc topics themes)"
	FlagFilenames         = "Show filenames between concatenated files"
	FlagHeaderFormat         = "Set the header display style (see: nanodoc topics filenames)"
	FlagFileNumbering     = "Set the file numbering style (see: nanodoc topics filenames)"
	FlagExt               = "Additional file extensions to treat as text"
	FlagInclude           = "Include only files matching patterns (see: nanodoc topics content)"
	FlagExclude           = "Exclude files matching patterns (see: nanodoc topics content)"
	FlagDryRun            = "Show what files would be processed without actually processing them"
	FlagVersion           = "Print the version number"
)

// Output messages
const (
	VersionFormat    = "nanodoc version %s (commit: %s, built: %s)\n"
	AvailableTopics  = "Available help topics:"
	RunTopicHelp     = `Run "nanodoc topics <topic-name>" for more information.`
	TopicNotFoundMsg = "topic not found"
)

// Man page constants
const (
	ManTitle   = "NANODOC"
	ManSection = "1"
	ManManual  = "Nanodoc Manual"
)

// Help template
const HelpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

// Misc group name for flags
const MiscGroupName = "Misc"

// Topic descriptions
var TopicDescriptions = map[string]string{
	"bundles":                "Create and manage bundle files for complex document combinations",
	"circular-dependencies":  "Understanding and resolving circular dependency issues",
	"content":                "File selection, patterns, and line ranges",
	"design":                 "Architecture and design principles of nanodoc",
	"filenames":                "Customize file filenames and separators with formatting options",
	"line-numbering":         "Add line numbers to your bundled documents with various modes",
	"themes":                 "Available themes and styling options",
	"toc":                    "Generate table of contents for your documents with navigation aids",
}
