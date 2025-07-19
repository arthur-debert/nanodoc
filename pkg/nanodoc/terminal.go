package nanodoc

import (
	"os"

	"golang.org/x/term"
)

// GetTerminalWidth returns the current terminal width, or the default if it cannot be determined
func GetTerminalWidth() int {
	// Try to get terminal width from stdout
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		// Try stderr as fallback
		width, _, err = term.GetSize(int(os.Stderr.Fd()))
		if err != nil || width <= 0 {
			// Return default if we can't determine terminal size
			return OUTPUT_WIDTH
		}
	}
	return width
}

// IsTerminal checks if the given file descriptor is a terminal
func IsTerminal(fd int) bool {
	return term.IsTerminal(fd)
}