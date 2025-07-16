package nanodoc

import (
	"errors"
	"testing"
)

func TestFileError(t *testing.T) {
	tests := []struct {
		name     string
		err      *FileError
		wantMsg  string
		wantWrap error
	}{
		{
			name: "file not found",
			err: &FileError{
				Path: "/path/to/file.txt",
				Err:  ErrFileNotFound,
			},
			wantMsg:  "/path/to/file.txt: file not found",
			wantWrap: ErrFileNotFound,
		},
		{
			name: "permission denied",
			err: &FileError{
				Path: "/etc/passwd",
				Err:  errors.New("permission denied"),
			},
			wantMsg:  "/etc/passwd: permission denied",
			wantWrap: errors.New("permission denied"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantMsg {
				t.Errorf("FileError.Error() = %v, want %v", got, tt.wantMsg)
			}
			if tt.err.Unwrap().Error() != tt.wantWrap.Error() {
				t.Errorf("FileError.Unwrap() = %v, want %v", tt.err.Unwrap(), tt.wantWrap)
			}
		})
	}
}

func TestCircularDependencyError(t *testing.T) {
	err := &CircularDependencyError{
		Path:  "bundle1.bundle.txt",
		Chain: []string{"bundle2.bundle.txt", "bundle3.bundle.txt", "bundle1.bundle.txt"},
	}

	want := "circular dependency detected: bundle1.bundle.txt -> [bundle2.bundle.txt bundle3.bundle.txt bundle1.bundle.txt]"
	if got := err.Error(); got != want {
		t.Errorf("CircularDependencyError.Error() = %v, want %v", got, want)
	}
}

func TestRangeError(t *testing.T) {
	tests := []struct {
		name     string
		err      *RangeError
		wantMsg  string
		wantWrap error
	}{
		{
			name: "invalid syntax",
			err: &RangeError{
				Input: "L10-5",
				Err:   errors.New("end line before start line"),
			},
			wantMsg:  "invalid range 'L10-5': end line before start line",
			wantWrap: errors.New("end line before start line"),
		},
		{
			name: "negative line number",
			err: &RangeError{
				Input: "L-5",
				Err:   errors.New("negative line number"),
			},
			wantMsg:  "invalid range 'L-5': negative line number",
			wantWrap: errors.New("negative line number"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantMsg {
				t.Errorf("RangeError.Error() = %v, want %v", got, tt.wantMsg)
			}
			if tt.err.Unwrap().Error() != tt.wantWrap.Error() {
				t.Errorf("RangeError.Unwrap() = %v, want %v", tt.err.Unwrap(), tt.wantWrap)
			}
		})
	}
}