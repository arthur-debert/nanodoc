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