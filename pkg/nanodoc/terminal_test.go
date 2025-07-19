package nanodoc

import (
	"os"
	"testing"
)

func TestGetTerminalWidth(t *testing.T) {
	// This test is challenging because it depends on the environment
	// We'll test that it returns a reasonable value
	width := GetTerminalWidth()
	
	// Should return a positive width
	if width <= 0 {
		t.Errorf("GetTerminalWidth() returned %d, expected positive value", width)
	}
	
	// In non-terminal environments (like CI), it should return the default
	if !IsTerminal(int(os.Stdout.Fd())) && !IsTerminal(int(os.Stderr.Fd())) {
		if width != OUTPUT_WIDTH {
			t.Errorf("GetTerminalWidth() returned %d in non-terminal environment, expected %d", width, OUTPUT_WIDTH)
		}
	}
}

func TestIsTerminal(t *testing.T) {
	// Test with stdout
	result := IsTerminal(int(os.Stdout.Fd()))
	
	// This will be true when running interactively, false in CI
	// We just verify it doesn't panic
	t.Logf("IsTerminal(stdout) = %v", result)
	
	// Test with stderr
	result = IsTerminal(int(os.Stderr.Fd()))
	t.Logf("IsTerminal(stderr) = %v", result)
	
	// Test with invalid FD
	result = IsTerminal(-1)
	if result {
		t.Error("IsTerminal(-1) returned true, expected false")
	}
}