# Troubleshooting Guide

This guide helps you resolve common issues when using nanodoc.

## Common Issues

### Line ranges not working

**Issue**: Running `nanodoc file.txt:L10-20` returns "file not found" or an error about line range syntax.

**Solution**: Line ranges only work within bundle files or live bundles, not as direct command-line arguments. To use line ranges:

1. Create a bundle file:
```bash
echo "file.txt:L10-20" > selection.bundle.txt
nanodoc selection.bundle.txt
```

2. Or use a live bundle with inline syntax:
```bash
echo "Here are lines 10-20: [[file:file.txt:L10-20]]" > doc.txt
nanodoc doc.txt
```

### Files not being processed

**Issue**: `.go`, `.py`, or other code files are skipped when processing directories.

**Solution**: By default, nanodoc only processes `.txt` and `.md` files. Use the `--txt-ext` flag to include additional extensions:

```bash
# Include Go files
nanodoc --txt-ext=go src/

# Include multiple extensions
nanodoc --txt-ext=go --txt-ext=py --txt-ext=js project/
```

### Bundle file not recognized

**Issue**: A bundle file is being treated as regular text instead of a list of files.

**Solution**: Bundle files must follow the `.bundle.*` naming pattern. Examples:
- ✅ `project.bundle.txt`
- ✅ `files.bundle.md`
- ❌ `bundle.txt` (missing prefix)
- ❌ `project_bundle.txt` (underscore instead of dot)

### Empty output or missing content

**Issue**: No content appears in the output or some files are missing.

**Solution**: Check the following:

1. **File permissions**: Ensure you have read access to all files
2. **File paths**: Use absolute paths or ensure relative paths are correct
3. **File extensions**: Verify files have the correct extensions or use `--txt-ext`
4. **Empty files**: Empty files will only show headers (this will be improved in future versions)

### Circular dependency errors

**Issue**: Error message about circular dependencies when using bundles.

**Solution**: This occurs when bundle files include each other in a loop. For example:
- `bundle1.txt` includes `bundle2.txt`
- `bundle2.txt` includes `bundle1.txt`

To fix:
1. Review your bundle files to identify the cycle
2. Check the error message which shows the exact dependency chain
3. Reorganize to remove the circular reference
4. Consider using a single master bundle file that includes all others

For detailed information, see the [Circular Dependencies Guide](docs/circular_dependencies.md).

### Live bundle syntax not working

**Issue**: `[[file:path]]` syntax appears as plain text instead of including file content.

**Solution**: 
1. Ensure you're using the correct syntax: `[[file:path/to/file.txt]]`
2. The file path must be valid and the file must exist
3. For line ranges: `[[file:path/to/file.txt:L10-20]]`
4. Note: The older `@[file]` syntax mentioned in some docs has been replaced with `[[file:]]`

### Performance issues with large files

**Issue**: Slow processing or high memory usage with large files or many files.

**Solution**:
1. Use line ranges to include only needed portions of large files
2. Process files in smaller batches
3. Consider splitting large documents into multiple outputs
4. Use bundle files to organize complex file selections

## Getting Help

If you encounter issues not covered here:

1. Check the [README](README.md) for updated usage information
2. Report bugs at https://github.com/arthur-debert/nanodoc-go/issues
3. Include the following in bug reports:
   - Nanodoc version (`nanodoc version`)
   - Operating system
   - Command that caused the issue
   - Error message (if any)
   - Sample files to reproduce (if possible)