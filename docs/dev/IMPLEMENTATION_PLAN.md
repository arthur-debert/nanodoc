# Nanodoc Implementation Plan

This document outlines the implementation plan for the Nanodoc project, a Go port of the Python-based document bundler.

**Project Goal**: Create a minimalist document bundler in Go that matches the Python version's functionality while leveraging Go's strengths.

## Core Architecture

The implementation follows a pipeline architecture:

```
CLI Args → Resolve Paths → Resolve Files → Gather Content → Build Document → Apply Formatting → Render Document → Output
```

How you will work:

* Make sure that each milestone step has at leas a commit for it. It's ok for a step to require multiple commits, but we should not mix various steps in the same commit, which makes review hard.
* As each milestone step is done, verify, tests , lint, commit and push
* Commit messages should follow the Milestone <number>: <step number> <description form>

## Implementation Milestones

### Milestone 1: Core Infrastructure (Steps 1-2)

**Goal**: Set up the basic project structure and data types

1. **Set up project structure and core types** (DONE)
   - Create package structure: `pkg/nanodoc/` with subpackages
   - Set up logging infrastructure
   - Create error types and constants
   - Update go.mod to use proper module name
   - Write tests for logging and error handling

2. **Implement structures (FileContent, Document, Range) with tests** (DONE)
   - Create `structures.go` with core data types:
     - `Range`: tuple of (start, end) line numbers
     - `FileContent`: filepath, ranges, content, metadata
     - `Document`: collection of FileContent with formatting options
   - Write unit tests for all data structures
   - Test validation and edge cases

### Milestone 2: File Processing Pipeline (Steps 3-6)

**Goal**: Implement the core file processing stages

3. **Implement path resolution (resolver.go) with tests** (DONE)
   - Resolve relative paths to absolute
   - Handle directories (expand to .txt/.md files)
   - Handle glob patterns
   - Detect bundle files (.bundle.*)
   - Write unit tests for all path resolution scenarios
   - Test edge cases (empty dirs, invalid paths, symlinks)

4. **Implement file resolution and content extraction with tests** (DONE)
   - Read file contents
   - Parse path:range syntax (e.g., "file.txt:L10-20")
   - Extract specified line ranges
   - Handle missing files gracefully
   - Write unit tests for file reading and range parsing
   - Test error cases (permissions, non-existent files)

5. **Implement content gathering with range support and tests** (DONE)
   - Apply line ranges to file content
   - Merge overlapping ranges
   - Preserve file metadata for rendering
   - Write unit tests for range operations
   - Test edge cases (empty ranges, out-of-bounds)

6. **Implement document building with bundle support and tests** (DONE)
   - Process bundle files (list of paths)
   - Handle circular dependency detection
   - Maintain document order
   - Write unit tests for bundle processing
   - Test circular dependency detection

### Milestone 3: Formatting and Rendering (Steps 7-8)

**Goal**: Implement the presentation layer

7. **Implement theme system and formatting with tests** (DONE)
   - Create theme structure (YAML-based)
   - Load built-in themes (classic, classic-light, classic-dark)
   - Apply syntax highlighting (if using rich formatting)
   - Support custom themes
   - Write unit tests for theme loading and application
   - Test invalid theme handling

8. **Implement document rendering with headers/TOC and tests** (DONE)
   - Generate file headers with different styles (nice, filename, path)
   - Support sequence numbering (numerical, letter, roman)
   - Generate table of contents
   - Render final output string
   - Write unit tests for all rendering options
   - Test TOC generation with various document structures

### Milestone 4: CLI and Features (Steps 9-10)

**Goal**: Create the command-line interface and main features

9. **Implement CLI with Cobra (flags, args parsing) with tests**
   * Create main command with all flags:
     * `-n/-nn`: line numbering modes
     * `--toc`: table of contents
     * `--theme`: theme selection
     * `--no-header`: disable headers
     * `--sequence`: numbering style
     * `--style`: header style
     * `--txt-ext`: additional extensions
   * Wire up to core pipeline
   * Write unit tests for CLI argument parsing
   * Test flag combinations and validation

10. **Implement line numbering modes (file/global) with tests**
    * Per-file numbering (-n): restart at 1 for each file
    * Global numbering (-nn): continuous across all files
    * Format with proper padding
    * Write unit tests for both numbering modes
    * Test edge cases (empty files, single lines)

### Milestone 5: Advanced Features (Steps 11-12) ✅ COMPLETED

**Goal**: Implement bundle and advanced file selection

11. **Implement bundle file support (.bundle.* files) with tests** ✅ DONE
    - Parse bundle files (one path per line)
    - Support comments in bundle files
    - Recursive bundle resolution
    - Write unit tests for bundle file parsing
    - Test recursive bundles and error cases

12. **Implement live bundle support (inline directives) with tests** ✅ DONE
    - Parse live bundle syntax
    - Replace file references with content inline
    - Support nested live bundles
    - Write unit tests for live bundle processing
    - Test nested bundle scenarios

### Milestone 6: Quality and Polish (Steps 13-14) ✅ COMPLETED

**Goal**: Ensure robustness and documentation

13. **Add comprehensive error handling** ✅ DONE
    - User-friendly error messages
    - Proper exit codes
    - Validation of inputs
    - Ensure all error paths are tested

14. **Update README and documentation** ✅ DONE
    - Installation instructions
    - Usage examples
    - API documentation for library use
    - Comprehensive feature documentation

## Technical Decisions

* Use `github.com/spf13/cobra` for CLI (already in go.mod)
* Use standard library for file operations
* Keep formatting simple initially (no syntax highlighting)
* Use YAML for theme files to match Python version
* Implement as both CLI tool and importable library
* Write tests alongside implementation (TDD approach)

## Testing Strategy

* **Unit Tests**: Write for each component as it's implemented
* **Integration Tests**: Test interactions between components
* **E2E Tests**: Test full CLI workflows
* **Coverage Goal**: Minimum 80% test coverage
* **Test Organization**: Place tests in `*_test.go` files alongside implementation

## Success Criteria

* Feature parity with Python version
* Similar performance or better
* Clean, idiomatic Go code
* Comprehensive test coverage (>80%)
* Clear documentation
* All tests passing before moving to next milestone