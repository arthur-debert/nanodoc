# Include/Exclude Patterns

Nanodoc supports gitignore-style patterns for fine-grained control over which files are processed when working with directories.

## Basic Usage

### Include Patterns
Process only files matching specific patterns:
```bash
# Include only markdown files in api directories
nanodoc docs/ --include="**/api/*.md"

# Include multiple patterns
nanodoc . --include="**/*.md" --include="**/*.txt"
```

### Exclude Patterns
Exclude files matching specific patterns:
```bash
# Exclude README files
nanodoc docs/ --exclude="**/README.md"

# Exclude test files and draft content
nanodoc . --exclude="**/*_test.go" --exclude="**/draft-*"
```

### Combining Include and Exclude
When both are specified, files must match include patterns and NOT match exclude patterns:
```bash
# Process all markdown files except those in test directories
nanodoc . --include="**/*.md" --exclude="**/test/**"

# Process API docs but exclude internal directories
nanodoc docs/ --include="**/api/**" --exclude="**/internal/**"
```

## Pattern Syntax

Patterns follow gitignore-style syntax:

- `*` matches any string within a path segment
- `**` matches zero or more directories
- `?` matches any single character
- `[abc]` matches any character in the set

### Examples

| Pattern | Matches | Doesn't Match |
|---------|---------|---------------|
| `*.md` | `file.md` | `dir/file.md` |
| `**/*.md` | `file.md`, `dir/file.md`, `a/b/c/file.md` | `file.txt` |
| `api/*.md` | `api/users.md` | `api/v1/users.md` |
| `**/api/*.md` | `api/users.md`, `docs/api/users.md` | `api/v1/users.md` |
| `**/test/**` | `test/file.md`, `src/test/file.go` | `test.md` |
| `README.*` | `README.md`, `README.txt` | `README/file.md` |

## Directory Traversal

- Without patterns: Only processes files in the specified directory (non-recursive)
- With patterns containing `**`: Recursively traverses subdirectories
- With patterns not containing `**`: Non-recursive (current behavior)

## Bundle File Support

Include/exclude patterns can be specified in bundle files:

```
# docs.bundle.txt
--include **/api/*.md
--exclude **/internal/**
--exclude **/test/**

# Additional files can be listed
important.md
```

## Integration with Other Features

Patterns work seamlessly with other nanodoc features:

```bash
# With additional extensions
nanodoc src/ --txt-ext=go --include="**/*.go" --exclude="**/*_test.go"

# With formatting options
nanodoc docs/ --include="**/api/**" --theme=dark --toc

# With line ranges
nanodoc . --include="**/*.md" --exclude="**/README.md" src/file.go:L10-20
```

## Use Cases

### Documentation Generation
```bash
# Include only public API documentation
nanodoc docs/ --include="**/api/**" --exclude="**/internal/**"
```

### Code Review Preparation
```bash
# Include source files but exclude tests and examples
nanodoc src/ --txt-ext=go --include="**/*.go" --exclude="**/*_test.go" --exclude="**/examples/**"
```

### Project Overview
```bash
# Include all markdown files except generated ones
nanodoc . --include="**/*.md" --exclude="**/node_modules/**" --exclude="**/vendor/**" --exclude="**/.git/**"
```

## Precedence Rules

1. If no patterns are specified, all text files in the directory are included
2. If include patterns are specified, only matching files are considered
3. Exclude patterns are always applied last and take precedence
4. Command-line patterns override bundle file patterns