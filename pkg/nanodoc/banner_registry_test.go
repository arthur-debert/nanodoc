package nanodoc

import (
	"strings"
	"testing"
)

func TestBannerRegistry(t *testing.T) {
	// Test that built-in styles are registered
	t.Run("built_in_styles_registered", func(t *testing.T) {
		expectedStyles := []string{"none", "dashed", "solid", "boxed"}
		
		registeredStyles := GetBannerStyleNames()
		
		// Check all expected styles are present
		for _, expected := range expectedStyles {
			found := false
			for _, registered := range registeredStyles {
				if registered == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected style %q not found in registry", expected)
			}
		}
	})
	
	// Test getting a specific style
	t.Run("get_style", func(t *testing.T) {
		style, exists := GetBannerStyle("dashed")
		if !exists {
			t.Fatal("Expected dashed style to exist")
		}
		if style.Name() != "dashed" {
			t.Errorf("Expected style name 'dashed', got %q", style.Name())
		}
	})
	
	// Test style application
	t.Run("apply_styles", func(t *testing.T) {
		opts := &FormattingOptions{
			HeaderAlignment: "left",
			PageWidth:       80,
		}
		
		tests := []struct {
			styleName string
			filename  string
			contains  []string
		}{
			{
				styleName: "none",
				filename:  "test.txt",
				contains:  []string{"test.txt"},
			},
			{
				styleName: "dashed",
				filename:  "test.txt",
				contains:  []string{"--------", "test.txt"},
			},
			{
				styleName: "solid",
				filename:  "test.txt",
				contains:  []string{"========", "test.txt"},
			},
			{
				styleName: "boxed",
				filename:  "test.txt",
				contains:  []string{"########", "### test.txt"},
			},
		}
		
		for _, tt := range tests {
			t.Run(tt.styleName, func(t *testing.T) {
				style, exists := GetBannerStyle(tt.styleName)
				if !exists {
					t.Fatalf("Style %q not found", tt.styleName)
				}
				
				result := style.Apply(tt.filename, opts)
				
				for _, expected := range tt.contains {
					if !strings.Contains(result, expected) {
						t.Errorf("Result does not contain %q\nGot: %s", expected, result)
					}
				}
			})
		}
	})
}

// Test custom banner style registration
type CustomBannerStyle struct{}

func (c CustomBannerStyle) Name() string        { return "custom" }
func (c CustomBannerStyle) Description() string { return "Custom test style" }
func (c CustomBannerStyle) Apply(filename string, opts *FormattingOptions) string {
	return ">>> " + filename + " <<<"
}

func TestBannerStyleAlignment(t *testing.T) {
	tests := []struct {
		name      string
		styleName string
		alignment string
		filename  string
		pageWidth int
		check     func(t *testing.T, result string)
	}{
		{
			name:      "none_center",
			styleName: "none",
			alignment: "center",
			filename:  "test.txt",
			pageWidth: 40,
			check: func(t *testing.T, result string) {
				// Should be centered in 40 chars
				if len(result) != 40 {
					t.Errorf("Expected length 40, got %d", len(result))
				}
				trimmed := strings.TrimSpace(result)
				if trimmed != "test.txt" {
					t.Errorf("Expected 'test.txt', got %q", trimmed)
				}
			},
		},
		{
			name:      "dashed_center",
			styleName: "dashed",
			alignment: "center",
			filename:  "test.txt",
			pageWidth: 40,
			check: func(t *testing.T, result string) {
				lines := strings.Split(result, "\n")
				if len(lines) != 3 {
					t.Errorf("Expected 3 lines, got %d", len(lines))
				}
				// Each line should be centered
				for _, line := range lines {
					if len(line) != 40 {
						t.Errorf("Expected line length 40, got %d for line: %q", len(line), line)
					}
				}
			},
		},
		{
			name:      "boxed_right",
			styleName: "boxed",
			alignment: "right",
			filename:  "short.txt",
			pageWidth: 60,
			check: func(t *testing.T, result string) {
				lines := strings.Split(result, "\n")
				// Middle line should have text aligned right
				if !strings.Contains(lines[1], "short.txt ###") {
					t.Errorf("Expected right-aligned text in boxed style, got: %q", lines[1])
				}
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style, exists := GetBannerStyle(tt.styleName)
			if !exists {
				t.Fatalf("Style %q not found", tt.styleName)
			}
			
			opts := &FormattingOptions{
				HeaderAlignment: tt.alignment,
				PageWidth:       tt.pageWidth,
			}
			
			result := style.Apply(tt.filename, opts)
			tt.check(t, result)
		})
	}
}

func TestCustomBannerStyle(t *testing.T) {
	// Create a new registry for this test to avoid conflicts
	registry := &BannerRegistry{
		styles: make(map[string]BannerStyle),
	}
	
	// Register custom style
	err := registry.Register(CustomBannerStyle{})
	if err != nil {
		t.Fatalf("Failed to register custom style: %v", err)
	}
	
	// Try to register again - should fail
	err = registry.Register(CustomBannerStyle{})
	if err == nil {
		t.Error("Expected error when registering duplicate style")
	}
	
	// Get the style
	style, exists := registry.Get("custom")
	if !exists {
		t.Fatal("Custom style not found")
	}
	
	// Apply the style
	result := style.Apply("test.txt", &FormattingOptions{})
	expected := ">>> test.txt <<<"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}