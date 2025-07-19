package main

import (
	"os"
)

func main() {
	if err := Execute(); err != nil {
		// Don't print error message here since we handle it in root.go
		os.Exit(1)
	}
} 