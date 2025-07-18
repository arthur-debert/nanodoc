# Bundle Options Feature Implementation Summary

## Overview
Successfully implemented the bundle options feature described in GitHub issue #17. This feature allows bundle files to store nanodoc processing options alongside file paths, enabling predictable and consistent output with a single command.

## Features Implemented

### 1. Bundle Options Parsing
- **Command-line flags in bundle files**: Lines starting with `--` are treated as command-line options
- **Supported options**: All major formatting options are supported:
  - `--toc` - Generate table of contents
  - `--line-numbers` / `-n` - Enable per-file line numbering  
  - `--global-line-numbers` / `-N` - Enable global line numbering
  - `--theme THEME` - Set theme (classic, classic-light, classic-dark)
  - `--no-header` - Disable file headers
  - `--header-style STYLE` - Set header style (nice, filename, path)
  - `--sequence STYLE` - Set sequence style (numerical, letter, roman)
  - `--txt-ext EXTENSION` - Add file extension to process

### 2. Precedence Rules
- **Command-line override**: Command-line arguments always override bundle options
- **Explicit flag tracking**: System tracks which flags were explicitly set via CLI
- **Merge logic**: Bundle options are only applied when CLI flags are not explicitly set

### 3. Syntax and Format
- **Comments**: Lines starting with `#` are ignored
- **Empty lines**: Empty lines are ignored
- **File paths**: Non-option lines are treated as file paths
- **Mixed content**: Options and file paths can be mixed in any order

## Example Usage

### Bundle File Format
```txt
# My Project Documentation Bundle
#
# This bundle defines both formatting options and the content to include.
# Lines starting with '#' are comments. Empty lines are ignored.

# --- Options ---
# Options are specified using the same flags as the command line.

--toc
--global-line-numbers
--header-style nice
--sequence roman
--theme classic-dark

# --- Content ---
# Files, directories, and glob patterns are listed below.

README.md
docs/
pkg/nanodoc/*.go
```

### Command Usage
```bash
# Apply all options from bundle automatically
nanodoc project.bundle.txt

# Override specific options from command line
nanodoc --theme classic-light project.bundle.txt  # CLI theme overrides bundle theme
```

## Implementation Details

### Core Components
1. **BundleOptions struct**: Holds parsed options from bundle files
2. **BundleResult struct**: Contains both file paths and options
3. **parseOption function**: Parses individual option lines with validation
4. **Merge functions**: Handle option precedence and merging logic
5. **CLI integration**: Tracks explicit flags and uses appropriate build functions

### Key Files Modified
- `pkg/nanodoc/bundle.go` - Core bundle processing logic
- `cmd/nanodoc/root.go` - CLI integration with explicit flag tracking
- `cmd/nanodoc/root_test.go` - Test updates to use new functions
- `docs/specifying_files.txt` - Documentation updates

### Test Coverage
- **Unit tests**: Comprehensive test suite for all bundle options
- **Integration tests**: End-to-end testing of CLI with bundle options
- **Edge cases**: Invalid options, circular dependencies, precedence rules
- **Validation**: Option value validation and error handling

## Benefits

1. **Predictable Output**: Bundle files ensure consistent formatting across runs
2. **Reusability**: Teams can share bundle files with standard formatting
3. **Flexibility**: CLI options can still override bundle settings when needed
4. **Backward Compatible**: Existing bundle files continue to work unchanged
5. **Simple Syntax**: Uses familiar command-line flag syntax

## Verification

The implementation has been verified to work correctly:
- All existing tests pass
- New bundle options tests pass
- CLI integration works as expected
- Documentation is updated and accurate
- Example bundle files work correctly

## Closes

This implementation fully addresses GitHub issue #17 - "Bundles should be able to store nanodoc processing options".