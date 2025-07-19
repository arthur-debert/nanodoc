# Nanodoc

A minimalist document bundler designed for stitching hints, reminders and short docs. Useful for prompts, personalized docs highlights for your teams or a note to your future self.

No config, nothing to learn nor remember. Short, simple, sweet.


## ✨ Features

- **Simple file bundling**: Combine multiple text files into one document
- **Bundle files**: Create reusable file lists with `.bundle.*` files
- **Live bundles**: Include files inline using `[[file:path]]` syntax
- **Line ranges**: Extract specific lines with `file.txt:L10-20` syntax
- **Pattern filtering**: Include/exclude files with gitignore-style patterns
- **Line numbering**: Add line numbers per-file (`-n`) or globally (`-nn`)
- **Table of contents**: Generate TOC with `--toc`
- **Multiple themes**: Built-in themes (classic, classic-light, classic-dark)
- **Smart path resolution**: Handles files, directories, and glob patterns
- **Dry run mode**: Preview files before processing with `--dry-run`
- **Error handling**: User-friendly error messages and proper exit codes

## 📖 Usage

### Basic Usage

```bash
# Bundle all .txt and .md files in current directory
nanodoc

# Bundle specific files
nanodoc file1.txt file2.md

# Use a bundle file
nanodoc project.bundle.txt

# Bundle files from a directory
nanodoc docs/

# Preview what files will be processed
nanodoc --dry-run docs/ *.md
```

### Advanced Usage

```bash
# With line numbers per file
nanodoc -n file1.txt file2.txt

# With global line numbers
nanodoc -nn *.txt

# With table of contents
nanodoc --toc chapter*.md

# With dark theme and no filenames
nanodoc --theme=classic-dark --filenames=false *.md

# Use glob patterns
nanodoc src/*.go docs/*.md

# Include/exclude patterns for directories
nanodoc docs/ --include="**/api/*.md" --exclude="**/internal/**"

# Process only Go source files, excluding tests
nanodoc src/ --ext=go --include="**/*.go" --exclude="**/*_test.go"
```

### Bundle Files

Create a `.bundle.txt` file to define reusable file lists and formatting options:

```txt
# My project bundle
--toc
--linenum global
--file-style nice
--file-numbering roman
--theme classic-dark

README.md
src/main.go
src/utils.go
docs/api.md
```

Bundle files support:
- **Bundle options**: Embed command-line flags directly in bundle files
  - Lines starting with `--` are treated as command-line options
  - Available options: `--toc`, `--theme`, `--linenum`, `--filenames`, `--file-style`, `--file-numbering`, `--ext`
  - Command-line options override bundle options when both are specified
- Comments (lines starting with `#`)
- Relative and absolute paths
- Recursive bundle inclusion
- Circular dependency detection

### Live Bundles

Include files inline using the `[[file:path]]` syntax:

```txt
# Project Overview

This is our main application:
[[file:src/main.go]]

Here are the utilities:
[[file:src/utils.go:L1-20]]

## Documentation
[[file:docs/README.md]]
```

Live bundles support:
- Nested includes (files can include other files)
- Line ranges with `[[file:path:L10-20]]`
- Circular reference detection (see [docs](docs/circular_dependencies.md))
- Graceful handling of missing files

## Installation

```bash
brew install nanodoc
```

Or download a `.deb` package from the [releases](https://github.com/arthur-debert/nanodoc/releases).

## 🤝 Contributing

Contributions are welcome! From feedback , to bug reports or actual code, are all very welcomed.
Use [Github issues](https://github.com/arthur-debert/nanodoc/issues)

##  License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
---

*Nanodoc: Because sometimes you just need to stitch files together, simply.*
