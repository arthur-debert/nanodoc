# Pull Request: Bundle Options Feature Implementation

## Overview

This PR implements the bundle options feature requested in **Issue #17**, allowing bundle files to store nanodoc processing options alongside file paths. This provides a way to achieve predictable and consistent output when using bundle files.

## Problem Statement

Previously, bundle files could only contain file paths. Users had to remember to specify the same command-line options each time they processed a bundle, leading to inconsistent output and a poor user experience.

## Solution

Bundle files can now contain command-line options that serve as defaults, with command-line arguments taking precedence when explicitly set. This allows users to create self-contained bundle files that produce consistent output.

## Key Features

✅ **Command-line Options in Bundle Files**: Bundle files can now contain options like `--toc`, `--line-numbers`, `--theme`, etc.

✅ **Same Syntax as CLI**: Options are specified using the same syntax as command-line flags

✅ **Proper Precedence**: Command-line options override bundle options when explicitly set

✅ **Comprehensive Coverage**: Supports all major formatting options (TOC, line numbers, themes, headers, sequences)

✅ **Backward Compatibility**: Existing bundle files continue to work unchanged

## Usage Example

### Before (Bundle files could only contain paths):
```
# project.bundle.txt
README.md
src/main.go
docs/api.md
```

### After (Bundle files can contain options + paths):
```
# project.bundle.txt
# Options - these serve as defaults
--toc
--line-numbers
--theme classic-dark
--header-style filename
--sequence roman

# Files
README.md
src/main.go
docs/api.md
```

### Command-line Override:
```bash
# Uses bundle options as defaults
nanodoc project.bundle.txt

# Command-line options override bundle options
nanodoc --theme classic-light project.bundle.txt
```

## Technical Implementation

### Core Changes

1. **Extended Bundle Processing**: 
   - Added `BundleOptions` struct to hold parsed options
   - Extended `ProcessBundleFileWithOptions` to parse both paths and options
   - Added `parseOption` function to handle individual option parsing

2. **Option Merging Logic**:
   - Implemented `MergeFormattingOptions` and `MergeFormattingOptionsWithDefaults`
   - Added explicit flag tracking in CLI to distinguish between default and user-set values
   - Proper precedence: CLI options > Bundle options > System defaults

3. **CLI Integration**:
   - Added `trackExplicitFlags` function to track which flags were explicitly set
   - Updated CLI to use `BuildDocumentWithExplicitFlags` for proper option merging
   - Maintained backward compatibility with existing CLI usage

### Bug Fixes

During implementation, discovered and fixed bugs in the renderer:

1. **Header Style Bug**: `HeaderStyleFilename` was returning full paths instead of just filenames
2. **Sequence Number Bug**: Sequence numbers weren't being applied to `HeaderStyleFilename` and `HeaderStylePath`

### Supported Options

All major nanodoc options are supported in bundle files:

- `--toc` - Generate table of contents
- `--line-numbers` / `-n` - Enable per-file line numbering  
- `--global-line-numbers` / `-N` - Enable global line numbering
- `--theme THEME` - Set theme (classic, classic-dark, classic-light)
- `--header-style STYLE` - Set header style (nice, filename, path)
- `--sequence STYLE` - Set sequence style (numerical, letter, roman)
- `--no-header` - Suppress file headers
- `--txt-ext EXT` - Additional file extensions to process

## Testing

### Comprehensive Test Coverage

- **Unit Tests**: All new functions have comprehensive unit tests
- **Integration Tests**: End-to-end tests covering the complete workflow
- **Edge Case Tests**: Invalid options, circular dependencies, complex merging scenarios
- **Regression Tests**: Ensuring existing functionality remains intact

### Manual Testing

Created comprehensive test scenarios to verify:
- Bundle options are correctly parsed and applied
- Command-line options properly override bundle options
- All formatting options work correctly (TOC, line numbers, themes, headers, sequences)
- Backward compatibility with existing bundle files

## Documentation Updates

- Updated `docs/specifying_files.txt` with comprehensive bundle options documentation
- README.md already contained bundle options documentation
- Added inline code comments explaining the new functionality
- Comprehensive test cases serve as usage examples

## Backward Compatibility

✅ **Fully Backward Compatible**: Existing bundle files continue to work unchanged
✅ **No Breaking Changes**: All existing CLI usage patterns remain supported
✅ **Graceful Degradation**: Invalid options in bundle files produce clear error messages

## Example Output

### Bundle File:
```
# example.bundle.txt
--toc
--line-numbers
--theme classic-dark
--header-style filename
--sequence roman

file1.txt
file2.txt
```

### Output:
```
Table of Contents
=================

- File1 (file1.txt)
- File2 (file2.txt)

i. file1.txt

1 | Content of file 1
2 | More content

ii. file2.txt

1 | Content of file 2
2 | More content
```

## Closes

Fixes #17

---

**Ready for Review**: This PR is complete and ready for review. The feature is fully implemented, tested, and documented.
