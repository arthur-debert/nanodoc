# Bundle Options Implementation Summary

## Issue #17: Bundle Options Feature

This document summarizes the complete implementation of the bundle options feature requested in [GitHub issue #17](https://github.com/arthur-debert/nanodoc-go/issues/17).

## Feature Overview

The bundle options feature allows users to embed command-line flags directly in bundle files, enabling predictable and consistent output without needing to remember and retype command-line arguments.

### Before
```txt
# my-bundle.txt
README.md
src/main.go
docs/api.md
```

```bash
nanodoc --toc --theme classic-dark --global-line-numbers --header-style nice --sequence roman my-bundle.txt
```

### After
```txt
# my-bundle.txt
--toc
--theme classic-dark
--global-line-numbers
--header-style nice
--sequence roman

README.md
src/main.go
docs/api.md
```

```bash
nanodoc my-bundle.txt  # All options applied automatically!
```

## Implementation Details

### Core Components

1. **Bundle Option Parsing** (`pkg/nanodoc/bundle.go`)
   - Extended `parseOption()` function to handle all CLI flags
   - Added validation for header-style and sequence values
   - Support for all available options: `--toc`, `--theme`, `--line-numbers`, `--global-line-numbers`, `--no-header`, `--header-style`, `--sequence`, `--txt-ext`

2. **Explicit Flag Tracking** (`cmd/nanodoc/root.go`)
   - Added `explicitFlags` map to track which CLI flags were explicitly set
   - Implemented proper precedence: CLI flags override bundle options
   - Added `BuildDocumentWithExplicitFlags()` function

3. **Option Merging** (`pkg/nanodoc/bundle.go`)
   - Enhanced `MergeFormattingOptionsWithDefaults()` to use explicit flag tracking
   - Multiple bundle files merge options (first bundle wins for conflicts)
   - Command-line options always take precedence over bundle options

4. **Data Structures** (`pkg/nanodoc/bundle.go`)
   - `BundleOptions` struct with pointers for optional values
   - `BundleResult` struct containing both paths and options
   - Proper handling of additional extensions merging

### Precedence Rules

The implementation follows a clear precedence hierarchy:

1. **Command-line flags** (highest priority)
2. **Bundle file options** (middle priority)
3. **Default values** (lowest priority)

Example:
```bash
# Bundle file contains: --theme classic-dark
# Command-line overrides: --theme classic-light
# Result: classic-light theme is used
```

### Supported Options

All nanodoc command-line options are supported in bundle files:

| Option | Description | Example |
|--------|-------------|---------|
| `--toc` | Generate table of contents | `--toc` |
| `--theme` | Set theme | `--theme classic-dark` |
| `--line-numbers` / `-n` | Per-file line numbers | `--line-numbers` |
| `--global-line-numbers` / `-N` | Global line numbers | `--global-line-numbers` |
| `--no-header` | Disable file headers | `--no-header` |
| `--header-style` | Header style | `--header-style filename` |
| `--sequence` | Sequence style | `--sequence roman` |
| `--txt-ext` | Additional extensions | `--txt-ext log` |

### Validation

The implementation includes comprehensive validation:

- **Header styles**: Must be one of `nice`, `filename`, `path`
- **Sequence styles**: Must be one of `numerical`, `letter`, `roman`
- **Required values**: Options like `--theme` must have a value
- **Option format**: Lines must start with `--` to be treated as options

### Error Handling

Clear error messages for invalid options:
```
Error: invalid header style: invalid (must be one of: nice, filename, path)
Error: --theme requires a value
Error: unknown option: --invalid-flag
```

## Testing

### Comprehensive Test Suite

1. **Bundle Option Parsing Tests**
   - Valid option parsing
   - Invalid option handling
   - Multiple options in single bundle
   - Option validation

2. **Precedence Tests**
   - Bundle options applied when no CLI flags
   - CLI flags override bundle options
   - Multiple bundle files merge correctly

3. **Integration Tests**
   - End-to-end CLI testing with bundle options
   - Real-world usage scenarios
   - Error handling validation

4. **Edge Cases**
   - Empty bundle files
   - Options-only bundles
   - Multiple `--txt-ext` values
   - Invalid option combinations

### Test Coverage

```bash
$ go test ./pkg/nanodoc -v -run TestBundleOptions
=== RUN   TestBundleOptions
--- PASS: TestBundleOptions (0.00s)
=== RUN   TestBundleOptionsValidation
--- PASS: TestBundleOptionsValidation (0.00s)
=== RUN   TestBundleOptionsCompleteIntegration
--- PASS: TestBundleOptionsCompleteIntegration (0.00s)
PASS
ok      github.com/arthur-debert/nanodoc-go/pkg/nanodoc 0.003s
```

## Documentation Updates

### Updated Documentation Files

1. **`docs/specifying_files.txt`** - Comprehensive bundle options documentation
2. **`README.md`** - Examples and usage patterns
3. **`docs/live_bundles.txt`** - Integration with live bundles
4. **`docs/options/`** - Individual option documentation

### Example Documentation

```txt
# Bundle Files with Options

Bundle files can contain formatting options mixed with file paths.
Lines starting with '--' are treated as command-line options:

# My Project Bundle
--toc
--theme classic-dark
--line-numbers
--header-style nice
--sequence roman

chapter1.txt
images/diagram.md
/absolute/path/to/notes.txt

Available options in bundle files:
- --toc                     Generate table of contents
- --theme THEME             Set theme (classic, classic-dark, classic-light)
- --line-numbers / -n       Enable per-file line numbering
- --global-line-numbers / -N Enable global line numbering
- --no-header               Disable file headers
- --header-style STYLE      Set header style (nice, filename, path)
- --sequence STYLE          Set sequence style (numerical, letter, roman)
- --txt-ext EXTENSION       Add file extension to process

Command-line options override bundle options when both are specified.
```

## Real-World Usage Examples

### 1. Project Documentation Bundle

```txt
# project-docs.bundle.txt
# Complete project documentation with consistent formatting

--toc
--global-line-numbers
--header-style nice
--sequence roman
--theme classic-dark

README.md
docs/architecture.md
docs/api.md
CHANGELOG.md
```

### 2. Code Review Bundle

```txt
# code-review.bundle.txt
# Code review bundle with line numbers for easy reference

--line-numbers
--header-style path
--sequence numerical
--txt-ext go
--txt-ext js

src/main.go
src/utils.go
frontend/app.js
```

### 3. Release Notes Bundle

```txt
# release-notes.bundle.txt
# Clean release notes without headers

--no-header
--theme classic-light

CHANGELOG.md:L1-50
docs/migration-guide.md
```

## Performance Impact

The bundle options feature has minimal performance impact:

- **Parsing**: Options are parsed during bundle file processing (already part of the pipeline)
- **Memory**: Minimal additional memory usage for option storage
- **Processing**: No impact on content processing or rendering

## Backward Compatibility

The implementation maintains full backward compatibility:

- **Existing bundle files**: Continue to work without modification
- **CLI behavior**: Unchanged for non-bundle usage
- **File format**: Comments and file paths work exactly as before

## Future Enhancements

Potential future improvements identified during implementation:

1. **Option validation**: More sophisticated validation rules
2. **Option inheritance**: Nested bundle option inheritance
3. **Custom options**: User-defined option aliases
4. **Option profiles**: Named option configurations

## Conclusion

The bundle options feature successfully addresses the requirements from issue #17:

✅ **Bundle files can store nanodoc processing options**
✅ **Predictable and consistent output**
✅ **Simple syntax using existing CLI flag format**
✅ **Proper precedence rules (CLI overrides bundle)**
✅ **Comprehensive validation and error handling**
✅ **Full backward compatibility**
✅ **Extensive test coverage**
✅ **Complete documentation**

The implementation follows the original design specification closely while adding robust error handling, validation, and testing that ensures reliability in production use.

## Related Files

### Core Implementation
- `pkg/nanodoc/bundle.go` - Bundle processing and option parsing
- `cmd/nanodoc/root.go` - CLI flag tracking and precedence
- `pkg/nanodoc/structures.go` - Data structures

### Tests
- `pkg/nanodoc/bundle_test.go` - Bundle option tests
- `cmd/nanodoc/root_test.go` - CLI integration tests

### Documentation
- `docs/specifying_files.txt` - Primary bundle options documentation
- `README.md` - Usage examples and feature overview
- `docs/options/` - Individual option documentation

This implementation fully satisfies the requirements of GitHub issue #17 and provides a solid foundation for future enhancements.