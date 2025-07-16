package nanodoc

// Default file extensions to process
var DefaultTextExtensions = []string{".txt", ".md"}

// Bundle file pattern
const BundlePattern = ".bundle."

// Line numbering modes
const (
	LineNumberModeNone   = ""
	LineNumberModeFile   = "file"
	LineNumberModeGlobal = "all"
)

// Header sequence types
const (
	SequenceNumerical = "numerical"
	SequenceLetter    = "letter"
	SequenceRoman     = "roman"
)

// Header styles
const (
	StyleNice     = "nice"
	StyleFilename = "filename"
	StylePath     = "path"
)

// Default theme names
const (
	ThemeClassic      = "classic"
	ThemeClassicLight = "classic-light"
	ThemeClassicDark  = "classic-dark"
)
