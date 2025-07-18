# Bundle Options Feature Implementation

## Overview

This PR implements the bundle options feature requested in [Issue #17](https://github.com/arthur-debert/nanodoc-go/issues/17) - "Bundles should be able to store nanodoc processing options".

Bundle files can now embed command-line options directly within the file, providing predictable and consistent output when processing bundles.

## What's New

### 🎯 Core Feature

Bundle files now support embedding command-line options alongside file paths:

```
# My project bundle
--toc
--global-line-numbers
--header-style nice
--sequence roman
--theme classic-dark
--txt-ext go

README.md
src/main.go
docs/api.md
```

### 🔧 Supported Options

All major command-line options are supported:
- `--toc` - Generate table of contents
- `--line-numbers` / `-n` - Per-file line numbering
- `--global-line-numbers` / `-N` - Global line numbering
- `--no-header` - Suppress file headers
- `--header-style` - Set header style (nice, filename, path)
- `--sequence` - Set sequence style (numerical, letter, roman)
- `--theme` - Set theme (classic, classic-light, classic-dark)
- `--txt-ext` - Add file extensions to process (supports multiple values)

### 📏 Precedence Rules

Command-line options always override bundle options, allowing:
- Bundle files to define standard, reproducible output
- Users to make ad-hoc adjustments from the command line

Example:
```bash
# Bundle contains --theme classic-dark
nanodoc --theme classic-light my-bundle.bundle.txt
# Result: classic-light theme is used (CLI wins)
```

## Changes Made

### Core Implementation

- **`BundleOptions` struct**: Holds formatting options parsed from bundle files
- **`BundleResult` struct**: Contains both file paths and options from a bundle
- **`ProcessBundleFileWithOptions()`**: Enhanced bundle processing to extract both files and options
- **`parseOption()`**: Parses individual command-line options from bundle files
- **`MergeFormattingOptions()`**: Merges bundle options with command-line options
- **`ExtractAndMergeBundleOptions()`**: Extracts and merges options from multiple bundle files

### Integration

- **`BuildDocument()`**: Modified to extract and merge bundle options before processing
- **Precedence handling**: Command-line options always override bundle options

### Bug Fixes

1. **Fixed circular dependency detection**: Removed duplicate checks that were causing false positives
2. **Fixed live bundle processing**: Added filtering to skip documentation files with `[[file:]]` examples

### Testing

- **`TestProcessBundleFileWithOptions()`**: Tests option parsing from bundle files
- **`TestBundleOptionsIntegration()`**: Tests end-to-end option merging
- **End-to-end testing**: Created `test-issue-17.bundle.txt` demonstrating the feature
- **All existing tests pass**: No regressions in existing functionality

### Documentation

- **`docs/options/bundle_options.txt`**: Comprehensive documentation covering syntax, usage, and best practices
- **Updated README**: Added bundle options to features list and examples
- **Implementation summary**: Detailed documentation of the implementation

## Demo

Created `test-issue-17.bundle.txt` to demonstrate the feature:

```
# Test Bundle for Issue #17 - Bundle Options Feature
--toc
--global-line-numbers
--header-style nice
--sequence roman
--theme classic-dark
--txt-ext go

README.md
pkg/nanodoc/constants.go
```

Run with:
```bash
nanodoc test-issue-17.bundle.txt
```

This generates output with:
- ✅ Table of contents
- ✅ Global line numbering
- ✅ Roman numeral sequences (i., ii., iii.)
- ✅ Nice header style
- ✅ Classic-dark theme
- ✅ Processes .go files as text

## Benefits

1. **Reproducible Output**: Every run produces the same formatting
2. **Team Consistency**: Share bundles that enforce consistent standards
3. **Convenience**: No need to remember complex command-line options
4. **Flexibility**: Command-line options can still override bundle options
5. **Backward Compatibility**: Existing bundle files continue to work unchanged

## Breaking Changes

None. This is a purely additive feature that maintains full backward compatibility.

## Testing

All tests pass:
```bash
$ go test ./pkg/nanodoc
ok      github.com/arthur-debert/nanodoc-go/pkg/nanodoc
```

## Migration

Existing bundle files continue to work unchanged. To add options to existing bundles:

1. Add options at the top of the bundle file
2. Ensure options start with `--` or `-`
3. Test with `--dry-run` to verify the configuration

## Related

Closes #17

## Type of Change

- [x] New feature (non-breaking change which adds functionality)
- [x] Bug fix (non-breaking change which fixes an issue)
- [x] Documentation update

## Checklist

- [x] Code follows the project's style guidelines
- [x] Self-review of code completed
- [x] Code is commented, particularly in hard-to-understand areas
- [x] Corresponding changes to documentation made
- [x] Tests added that prove the fix is effective or that the feature works
- [x] New and existing unit tests pass locally
- [x] Changes generate no new warnings
- [x] Dependent changes merged and published