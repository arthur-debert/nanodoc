# Bundle Options Feature Implementation

## Summary
Implements GitHub issue #17 - "Bundles should be able to store nanodoc processing options"

This PR adds the ability to embed command-line options directly in bundle files, enabling predictable and consistent output without requiring users to remember complex command-line arguments.

## Key Features

### Bundle Options Support
- **Syntax**: Lines starting with `--` are treated as command-line options
- **All options supported**: `--toc`, `--theme`, `--line-numbers`, `--global-line-numbers`, `--no-header`, `--header-style`, `--sequence`, `--txt-ext`
- **Precedence**: Command-line flags override bundle options
- **Validation**: Full validation of option values with helpful error messages

### Example Usage
```txt
# project.bundle.txt
--toc
--global-line-numbers
--header-style nice
--sequence roman
--theme classic-dark

README.md
docs/
```

```bash
# All options applied automatically
nanodoc project.bundle.txt

# CLI can override bundle options
nanodoc --theme classic-light project.bundle.txt
```

## Implementation Details

### Core Changes
- **Bundle parsing**: Enhanced `parseOption()` function with comprehensive validation
- **CLI integration**: Added explicit flag tracking with `explicitFlags` map
- **Option merging**: Proper precedence handling in `MergeFormattingOptionsWithDefaults()`
- **Data structures**: New `BundleOptions` and `BundleResult` structs

### Testing
- **Unit tests**: Comprehensive test suite for all bundle options
- **Integration tests**: End-to-end CLI testing with bundle options
- **Edge cases**: Invalid options, precedence rules, validation
- **All existing tests pass**: Full backward compatibility maintained

### Documentation
- **Updated docs**: Complete documentation in `docs/specifying_files.txt`
- **Examples**: Real-world usage examples in README
- **Implementation summary**: Detailed technical documentation

## Benefits
- **Predictable output**: Consistent formatting across runs
- **Team collaboration**: Shareable bundle files with standard formatting
- **Flexibility**: CLI options can still override bundle settings
- **Backward compatible**: Existing bundle files continue to work unchanged
- **Simple syntax**: Uses familiar command-line flag syntax

## Verification
✅ All tests pass  
✅ Feature works as designed  
✅ Documentation updated  
✅ Backward compatibility maintained  
✅ Zero breaking changes  

## Closes
Closes #17