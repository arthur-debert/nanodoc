package nanodoc

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Theme represents a formatting theme with style definitions
type Theme struct {
	Name   string
	Styles map[string]string
}

// FormattingContext holds the state for formatting operations
type FormattingContext struct {
	Theme         *Theme
	LineNumbers   LineNumberMode
	ShowFilenames   bool
	FilenameStyle   FilenameStyle
	SequenceStyle SequenceStyle
	ShowTOC       bool
}

const (
	// DefaultTheme is the default theme name
	DefaultTheme = "classic"
)

//go:embed themes/*.yaml
var themesFS embed.FS

// GetAvailableThemes returns a list of available theme names
func GetAvailableThemes() ([]string, error) {
	entries, err := themesFS.ReadDir("themes")
	if err != nil {
		return nil, fmt.Errorf("failed to read themes directory: %w", err)
	}

	var themes []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			themeName := strings.TrimSuffix(entry.Name(), ".yaml")
			themes = append(themes, themeName)
		}
	}

	slog.Debug("Found available themes", "themes", themes)
	return themes, nil
}

// LoadTheme loads a theme from the embedded filesystem
func LoadTheme(themeName string) (*Theme, error) {
	if themeName == "" {
		themeName = DefaultTheme
	}

	slog.Debug("Loading theme", "name", themeName)

	// Try to load the requested theme
	themeData, err := loadThemeFile(themeName)
	if err != nil {
		// Fall back to default theme
		slog.Warn("Failed to load theme, falling back to default", 
			"theme", themeName, "error", err)
		
		themeData, err = loadThemeFile(DefaultTheme)
		if err != nil {
			return nil, fmt.Errorf("failed to load default theme: %w", err)
		}
		themeName = DefaultTheme
	}

	theme := &Theme{
		Name:   themeName,
		Styles: themeData,
	}

	slog.Debug("Theme loaded successfully", "name", themeName)
	return theme, nil
}

// LoadCustomTheme loads a theme from a custom file path
func LoadCustomTheme(themePath string) (*Theme, error) {
	slog.Debug("Loading custom theme", "path", themePath)

	data, err := os.ReadFile(themePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file: %w", err)
	}

	var styles map[string]string
	if err := yaml.Unmarshal(data, &styles); err != nil {
		return nil, fmt.Errorf("failed to parse theme YAML: %w", err)
	}

	themeName := strings.TrimSuffix(filepath.Base(themePath), ".yaml")
	theme := &Theme{
		Name:   themeName,
		Styles: styles,
	}

	return theme, nil
}

// loadThemeFile loads a theme from the embedded filesystem
func loadThemeFile(themeName string) (map[string]string, error) {
	themePath := fmt.Sprintf("themes/%s.yaml", themeName)
	
	data, err := themesFS.ReadFile(themePath)
	if err != nil {
		return nil, fmt.Errorf("theme not found: %s", themeName)
	}

	var styles map[string]string
	if err := yaml.Unmarshal(data, &styles); err != nil {
		return nil, fmt.Errorf("failed to parse theme YAML: %w", err)
	}

	return styles, nil
}

// NewFormattingContext creates a new formatting context with the given options
func NewFormattingContext(options FormattingOptions) (*FormattingContext, error) {
	theme, err := LoadTheme(options.Theme)
	if err != nil {
		return nil, fmt.Errorf("failed to load theme: %w", err)
	}

	return &FormattingContext{
		Theme:         theme,
		LineNumbers:   options.LineNumbers,
		ShowFilenames:   options.ShowFilenames,
		FilenameStyle:   options.FilenameStyle,
		SequenceStyle: options.SequenceStyle,
		ShowTOC:       options.ShowTOC,
	}, nil
}

// ApplyTheme applies the theme to a document (placeholder for now)
func (fc *FormattingContext) ApplyTheme(doc *Document) error {
	// This is a placeholder - actual formatting will be implemented
	// when we have the rendering pipeline
	slog.Debug("Applying theme to document", "theme", fc.Theme.Name)
	return nil
}