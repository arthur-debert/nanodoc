package nanodoc

// Range represents a line range in a file
// Start is 1-based inclusive, End is 1-based inclusive (or 0 for EOF)
type Range struct {
	Start int
	End   int // 0 means end of file
}

// FileContent represents the content and metadata for a single file
type FileContent struct {
	// Path to the file (absolute path)
	Filepath string

	// Line ranges to include
	Ranges []Range

	// Content after applying ranges
	Content string

	// True if this represents a bundle file
	IsBundle bool

	// Source file if part of an inline bundle
	OriginalSource string
}

// Document represents the entire document after processing bundles
type Document struct {
	// Ordered list of content blocks
	ContentItems []FileContent

	// Table of contents data
	TOC []TOCEntry

	// Name of the theme to use for styling
	ThemeName string

	// Whether to use Rich formatting
	UseRichFormatting bool

	// Formatting options
	FormattingOptions FormattingOptions
}

// TOCEntry represents an entry in the table of contents
type TOCEntry struct {
	// Display title for the TOC entry
	Title string

	// File path this entry refers to
	Path string

	// Sequence number/letter/roman numeral if applicable
	Sequence string

	// Line number in the final document
	LineNumber int
}

// FormattingOptions contains all formatting-related options
type FormattingOptions struct {
	// Line numbering mode: "", "file", or "all"
	LineNumberMode string

	// Whether to show headers
	ShowHeader bool

	// Header sequence type: "numerical", "letter", "roman"
	Sequence string

	// Header style: "nice", "filename", "path"
	Style string

	// Additional file extensions to process
	AdditionalExtensions []string
}

// NewRange creates a new Range with validation
func NewRange(start, end int) (Range, error) {
	if start < 1 {
		return Range{}, &RangeError{
			Input: "start line",
			Err:   ErrInvalidRange,
		}
	}
	if end != 0 && end < start {
		return Range{}, &RangeError{
			Input: "end before start",
			Err:   ErrInvalidRange,
		}
	}
	return Range{Start: start, End: end}, nil
}

// Contains checks if a line number is within this range
func (r Range) Contains(line int) bool {
	if line < r.Start {
		return false
	}
	if r.End == 0 {
		return true // EOF
	}
	return line <= r.End
}

// IsFullFile returns true if this range represents the entire file
func (r Range) IsFullFile() bool {
	return r.Start == 1 && r.End == 0
}

// NewDocument creates a new Document with default options
func NewDocument() *Document {
	return &Document{
		ContentItems: make([]FileContent, 0),
		TOC:          make([]TOCEntry, 0),
		FormattingOptions: FormattingOptions{
			ShowHeader: true,
			Style:      StyleNice,
		},
	}
}