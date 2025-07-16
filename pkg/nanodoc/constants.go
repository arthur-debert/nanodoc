package nanodoc

// Default file extensions to process
var DefaultTextExtensions = []string{".txt", ".md"}

// Bundle file pattern
const BundlePattern = ".bundle."

// LineNumberMode represents different line numbering modes
type LineNumberMode int

const (
	// LineNumberNone - no line numbers
	LineNumberNone LineNumberMode = iota
	// LineNumberFile - restart numbering for each file
	LineNumberFile
	// LineNumberGlobal - continuous numbering across all files
	LineNumberGlobal
)

// HeaderStyle represents different header styles
type HeaderStyle string

const (
	// HeaderStyleNice - formatted headers with decorations
	HeaderStyleNice HeaderStyle = "nice"
	// HeaderStyleFilename - just the filename
	HeaderStyleFilename HeaderStyle = "filename"
	// HeaderStylePath - full file path
	HeaderStylePath HeaderStyle = "path"
)

// SequenceStyle represents different sequence numbering styles
type SequenceStyle string

const (
	// SequenceNumerical - 1, 2, 3...
	SequenceNumerical SequenceStyle = "numerical"
	// SequenceLetter - a, b, c...
	SequenceLetter SequenceStyle = "letter"
	// SequenceRoman - i, ii, iii...
	SequenceRoman SequenceStyle = "roman"
)

// Default theme names
const (
	ThemeClassic      = "classic"
	ThemeClassicLight = "classic-light"
	ThemeClassicDark  = "classic-dark"
)

// FilePatterns are the default file patterns to match when scanning directories
var FilePatterns = []string{"*.txt", "*.md"}