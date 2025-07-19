package nanodoc

import (
	"fmt"
	"strings"
	"sync"
)

// BannerStyle defines the interface for banner style implementations
type BannerStyle interface {
	// Apply formats the filename with the banner style
	Apply(filename string, opts *FormattingOptions) string
	// Name returns the name of the banner style
	Name() string
	// Description returns a description of the banner style
	Description() string
}

// BannerRegistry manages banner style implementations
type BannerRegistry struct {
	mu     sync.RWMutex
	styles map[string]BannerStyle
}

// Global banner style registry
var globalBannerRegistry = &BannerRegistry{
	styles: make(map[string]BannerStyle),
}

// Register adds a new banner style to the registry
func (r *BannerRegistry) Register(style BannerStyle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	name := style.Name()
	if _, exists := r.styles[name]; exists {
		return fmt.Errorf("banner style %q already registered", name)
	}
	
	r.styles[name] = style
	return nil
}

// Get retrieves a banner style by name
func (r *BannerRegistry) Get(name string) (BannerStyle, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	style, exists := r.styles[name]
	return style, exists
}

// List returns all registered banner style names
func (r *BannerRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	names := make([]string, 0, len(r.styles))
	for name := range r.styles {
		names = append(names, name)
	}
	return names
}

// GetDescriptions returns a map of style names to descriptions
func (r *BannerRegistry) GetDescriptions() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	descriptions := make(map[string]string)
	for name, style := range r.styles {
		descriptions[name] = style.Description()
	}
	return descriptions
}

// RegisterBannerStyle registers a banner style in the global registry
func RegisterBannerStyle(style BannerStyle) error {
	return globalBannerRegistry.Register(style)
}

// GetBannerStyle retrieves a banner style from the global registry
func GetBannerStyle(name string) (BannerStyle, bool) {
	return globalBannerRegistry.Get(name)
}

// GetBannerStyleNames returns all registered banner style names
func GetBannerStyleNames() []string {
	return globalBannerRegistry.List()
}

// GetBannerStyleDescriptions returns descriptions of all banner styles
func GetBannerStyleDescriptions() map[string]string {
	return globalBannerRegistry.GetDescriptions()
}

// applyAlignment applies text alignment within the given width
func applyAlignment(text, alignment string, width int) string {
	textLen := len(text)
	if textLen >= width {
		return text
	}
	
	switch alignment {
	case "center":
		leftPadding := (width - textLen) / 2
		rightPadding := width - textLen - leftPadding
		return strings.Repeat(" ", leftPadding) + text + strings.Repeat(" ", rightPadding)
	case "right":
		padding := width - textLen
		return strings.Repeat(" ", padding) + text
	default: // left
		padding := width - textLen
		return text + strings.Repeat(" ", padding)
	}
}

// Built-in banner style implementations

// NoneBannerStyle displays just the filename with optional alignment
type NoneBannerStyle struct{}

func (n NoneBannerStyle) Name() string        { return "none" }
func (n NoneBannerStyle) Description() string { return "No banner decoration" }

func (n NoneBannerStyle) Apply(filename string, opts *FormattingOptions) string {
	// Apply alignment
	return applyAlignment(filename, opts.HeaderAlignment, opts.PageWidth)
}

// DashedBannerStyle uses dashed lines above and below
type DashedBannerStyle struct{}

func (d DashedBannerStyle) Name() string        { return "dashed" }
func (d DashedBannerStyle) Description() string { return "Dashed lines above and below" }

func (d DashedBannerStyle) Apply(filename string, opts *FormattingOptions) string {
	// For dashed/solid styles, we keep the line length matching the text
	// but apply alignment to the whole block
	line := strings.Repeat("-", len(filename))
	block := fmt.Sprintf("%s\n%s\n%s", line, filename, line)
	
	// For non-left alignment, we need to align each line
	if opts.HeaderAlignment != "left" && opts.HeaderAlignment != "" {
		lines := strings.Split(block, "\n")
		for i, l := range lines {
			lines[i] = applyAlignment(l, opts.HeaderAlignment, opts.PageWidth)
		}
		return strings.Join(lines, "\n")
	}
	
	return block
}

// SolidBannerStyle uses solid lines above and below
type SolidBannerStyle struct{}

func (s SolidBannerStyle) Name() string        { return "solid" }
func (s SolidBannerStyle) Description() string { return "Solid lines above and below" }

func (s SolidBannerStyle) Apply(filename string, opts *FormattingOptions) string {
	// For dashed/solid styles, we keep the line length matching the text
	// but apply alignment to the whole block
	line := strings.Repeat("=", len(filename))
	block := fmt.Sprintf("%s\n%s\n%s", line, filename, line)
	
	// For non-left alignment, we need to align each line
	if opts.HeaderAlignment != "left" && opts.HeaderAlignment != "" {
		lines := strings.Split(block, "\n")
		for i, l := range lines {
			lines[i] = applyAlignment(l, opts.HeaderAlignment, opts.PageWidth)
		}
		return strings.Join(lines, "\n")
	}
	
	return block
}

// BoxedBannerStyle creates a box around the filename
type BoxedBannerStyle struct{}

func (b BoxedBannerStyle) Name() string        { return "boxed" }
func (b BoxedBannerStyle) Description() string { return "Box with hash characters" }

func (b BoxedBannerStyle) Apply(filename string, opts *FormattingOptions) string {
	// Calculate padding for boxed style
	borderChar := "#"
	borderLength := opts.PageWidth
	if borderLength < len(filename)+8 { // Minimum space for "### text ###"
		borderLength = len(filename) + 8
	}
	
	topBottom := strings.Repeat(borderChar, borderLength)
	
	// Calculate padding based on alignment
	innerWidth := borderLength - 6 // Account for "### " and " ###"
	var middleLine string
	
	switch opts.HeaderAlignment {
	case "center":
		leftPadding := (innerWidth - len(filename)) / 2
		rightPadding := innerWidth - len(filename) - leftPadding
		middleLine = fmt.Sprintf("### %s%s%s ###", 
			strings.Repeat(" ", leftPadding),
			filename,
			strings.Repeat(" ", rightPadding))
	case "right":
		leftPadding := innerWidth - len(filename)
		middleLine = fmt.Sprintf("### %s%s ###", 
			strings.Repeat(" ", leftPadding),
			filename)
	default: // left
		rightPadding := innerWidth - len(filename)
		middleLine = fmt.Sprintf("### %s%s ###", 
			filename,
			strings.Repeat(" ", rightPadding))
	}
	
	return fmt.Sprintf("%s\n%s\n%s", topBottom, middleLine, topBottom)
}

// Initialize built-in banner styles
func init() {
	// Register built-in styles
	_ = RegisterBannerStyle(NoneBannerStyle{})
	_ = RegisterBannerStyle(DashedBannerStyle{})
	_ = RegisterBannerStyle(SolidBannerStyle{})
	_ = RegisterBannerStyle(BoxedBannerStyle{})
}