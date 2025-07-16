package nanodoc

import (
	"errors"
	"fmt"
)

// Common errors
var (
	// ErrFileNotFound is returned when a file cannot be found
	ErrFileNotFound = errors.New("file not found")

	// ErrCircularDependency is returned when a circular dependency is detected in bundles
	ErrCircularDependency = errors.New("circular dependency detected")

	// ErrInvalidRange is returned when a line range is invalid
	ErrInvalidRange = errors.New("invalid line range")

	// ErrEmptySource is returned when no source files are provided
	ErrEmptySource = errors.New("no source files provided")

	// ErrInvalidTheme is returned when a theme cannot be loaded
	ErrInvalidTheme = errors.New("invalid or missing theme")
)

// FileError represents an error related to a specific file
type FileError struct {
	Path string
	Err  error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%s: %v", e.Path, e.Err)
}

func (e *FileError) Unwrap() error {
	return e.Err
}

// CircularDependencyError represents a circular dependency in bundle files
type CircularDependencyError struct {
	Path  string
	Chain []string
}

func (e *CircularDependencyError) Error() string {
	return fmt.Sprintf("circular dependency detected: %s -> %s", e.Path, e.Chain)
}

// RangeError represents an error in line range specification
type RangeError struct {
	Input string
	Err   error
}

func (e *RangeError) Error() string {
	return fmt.Sprintf("invalid range '%s': %v", e.Input, e.Err)
}

func (e *RangeError) Unwrap() error {
	return e.Err
}
