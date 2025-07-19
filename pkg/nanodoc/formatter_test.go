package nanodoc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAvailableThemes(t *testing.T) {
	themes, err := GetAvailableThemes()
	if err != nil {
		t.Fatalf("Failed to get available themes: %v", err)
	}

	// Check that we have at least the default themes
	expectedThemes := []string{"classic", "classic-dark", "classic-light"}
	for _, expected := range expectedThemes {
		found := false
		for _, theme := range themes {
			if theme == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected theme %q not found in available themes", expected)
		}
	}

	// Ensure we have at least the expected number of themes
	if len(themes) < len(expectedThemes) {
		t.Errorf("Expected at least %d themes, got %d", len(expectedThemes), len(themes))
	}
}

func TestLoadTheme(t *testing.T) {
	tests := []struct {
		name      string
		themeName string
		wantErr   bool
	}{
		{
			name:      "load default theme",
			themeName: "classic",
			wantErr:   false,
		},
		{
			name:      "load dark theme",
			themeName: "classic-dark",
			wantErr:   false,
		},
		{
			name:      "load light theme",
			themeName: "classic-light",
			wantErr:   false,
		},
		{
			name:      "empty theme name loads default",
			themeName: "",
			wantErr:   false,
		},
		{
			name:      "non-existent theme falls back to default",
			themeName: "non-existent-theme",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme, err := LoadTheme(tt.themeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTheme() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if theme == nil {
					t.Error("Expected theme to be non-nil")
					return
				}

				// Verify theme has required styles
				requiredStyles := []string{"heading", "error"}
				for _, style := range requiredStyles {
					if _, ok := theme.Styles[style]; !ok {
						t.Errorf("Theme missing required style: %s", style)
					}
				}

				// If non-existent theme was requested, it should fall back to default
				if tt.themeName == "non-existent-theme" || tt.themeName == "" {
					if theme.Name != DefaultTheme {
						t.Errorf("Expected theme name to be %q, got %q", DefaultTheme, theme.Name)
					}
				}
			}
		})
	}
}

func TestLoadCustomTheme(t *testing.T) {
	// Create a temporary theme file
	tmpDir := t.TempDir()
	customThemePath := filepath.Join(tmpDir, "custom.yaml")
	
	customThemeContent := `heading: "green bold"
error: "yellow italic"
title: "cyan underline"`

	if err := os.WriteFile(customThemePath, []byte(customThemeContent), 0644); err != nil {
		t.Fatalf("Failed to create custom theme file: %v", err)
	}

	// Test loading the custom theme
	theme, err := LoadCustomTheme(customThemePath)
	if err != nil {
		t.Fatalf("Failed to load custom theme: %v", err)
	}

	if theme.Name != "custom" {
		t.Errorf("Expected theme name to be 'custom', got %q", theme.Name)
	}

	// Verify custom styles
	expectedStyles := map[string]string{
		"heading": "green bold",
		"error":   "yellow italic",
		"title":   "cyan underline",
	}

	for key, expected := range expectedStyles {
		if got, ok := theme.Styles[key]; !ok || got != expected {
			t.Errorf("Style %q = %q, want %q", key, got, expected)
		}
	}

	// Test loading non-existent file
	_, err = LoadCustomTheme("/non/existent/theme.yaml")
	if err == nil {
		t.Error("Expected error when loading non-existent theme file")
	}
}

func TestNewFormattingContext(t *testing.T) {
	tests := []struct {
		name    string
		options FormattingOptions
		wantErr bool
	}{
		{
			name: "default options",
			options: FormattingOptions{
				Theme:         "classic",
				LineNumbers:   LineNumberNone,
				ShowFilenames:   true,
				FilenameStyle:   FilenameStyleNice,
				SequenceStyle: SequenceNumerical,
				ShowTOC:       false,
			},
			wantErr: false,
		},
		{
			name: "with line numbers",
			options: FormattingOptions{
				Theme:       "classic-dark",
				LineNumbers: LineNumberFile,
				ShowFilenames: true,
			},
			wantErr: false,
		},
		{
			name: "with TOC",
			options: FormattingOptions{
				Theme:   "classic-light",
				ShowTOC: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := NewFormattingContext(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFormattingContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if ctx == nil {
					t.Error("Expected context to be non-nil")
					return
				}

				// Verify context properties
				if ctx.Theme == nil {
					t.Error("Expected theme to be non-nil")
				}
				if ctx.LineNumbers != tt.options.LineNumbers {
					t.Errorf("LineNumbers = %v, want %v", ctx.LineNumbers, tt.options.LineNumbers)
				}
				if ctx.ShowFilenames != tt.options.ShowFilenames {
					t.Errorf("ShowFilenames = %v, want %v", ctx.ShowFilenames, tt.options.ShowFilenames)
				}
				if ctx.ShowTOC != tt.options.ShowTOC {
					t.Errorf("ShowTOC = %v, want %v", ctx.ShowTOC, tt.options.ShowTOC)
				}
			}
		})
	}
}