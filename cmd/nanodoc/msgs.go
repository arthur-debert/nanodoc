package main

// Command descriptions
const (
	RootShort = "A minimalist document bundler"
	
	TopicsShort = "Display help topics"
	TopicsLong  = `Display available help topics or show the content of a specific topic.

Running 'nanodoc topics' lists all available topics.
Running 'nanodoc topics <topic-name>' displays the content of that topic.`
	
	CompletionShort = "Generate completion script"
	CompletionLong  = `To load completions:

Bash:

  $ source <(nanodoc completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ nanodoc completion bash > /etc/bash_completion.d/nanodoc
  # macOS:
  $ nanodoc completion bash > $(brew --prefix)/etc/bash_completion.d/nanodoc

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ nanodoc completion zsh > "${fpath[1]}/_nanodoc"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ nanodoc completion fish | source

  # To load completions for each session, execute once:
  $ nanodoc completion fish > ~/.config/fish/completions/nanodoc.fish

PowerShell:

  PS> nanodoc completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> nanodoc completion powershell > nanodoc.ps1
  # and source this file from your PowerShell profile.
`
	
	ManShort = "Generate man page"
	ManLong  = `Generate a man page for nanodoc`
)

// Error messages
const (
	ErrMinArgs             = "requires at least 1 arg(s), only received %d"
	ErrResolvingPaths      = "error resolving paths: %w"
	ErrGeneratingDryRun    = "error generating dry run info: %w"
	ErrBuildingDocument    = "error building document: %w"
	ErrCreatingContext     = "error creating formatting context: %w"
	ErrRenderingDocument   = "error rendering document: %w"
	ErrTopicNotFound       = "topic '%s' not found"
	ErrFailedToGetTopics   = "failed to get available topics: %w"
	ErrFailedGenManPage    = "failed to generate man page: %w"
)

// Flag descriptions
const (
	FlagLineNumbers       = "Enable per-file line numbering (see: nanodoc topics line-numbering)"
	FlagGlobalLineNumbers = "Enable global line numbering (see: nanodoc topics line-numbering)"
	FlagTOC               = "Generate a table of contents (see: nanodoc topics toc)"
	FlagTheme             = "Set the theme for formatting (see: nanodoc topics themes)"
	FlagNoHeader          = "Suppress file headers"
	FlagHeaderStyle       = "Set the header style (see: nanodoc topics headers)"
	FlagSequence          = "Set the sequence style (see: nanodoc topics headers)"
	FlagTxtExt            = "Additional file extensions to treat as text"
	FlagInclude           = "Include only files matching patterns (see: nanodoc topics content)"
	FlagExclude           = "Exclude files matching patterns (see: nanodoc topics content)"
	FlagDryRun            = "Show what files would be processed without actually processing them"
	FlagVersion           = "Print the version number"
)

// Output messages
const (
	VersionFormat         = "nanodoc version %s (commit: %s, built: %s)\n"
	AvailableTopics       = "Available help topics:"
	RunTopicHelp          = `Run "nanodoc topics <topic-name>" for more information.`
	TopicNotFoundMsg      = "topic not found"
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