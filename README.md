# Nanodoc

Nanodoc is a minimalist document bundler designed for stitching hints, reminders, and short docs. It's useful for creating prompts, personalized documentation highlights for your teams, or a note to your future self.

No config, nothing to learn nor remember. Short, simple, sweet.

## Installation

### From source

To build from source, you'll need Go 1.23+ installed.

```bash
git clone https://github.com/arthur-debert/nanodoc-go.git
cd nanodoc-go
go build -o nanodoc ./cmd/nanodoc
```

## Usage

```bash
nanodoc [paths...] [flags]
```

### Basic Example

Combine two files into a single document:

```bash
nanodoc file1.txt file2.md
```

This will output the contents of `file1.txt` followed by `file2.md`, each with a "nice" header.

### Specifying Files

You can specify files, directories, or glob patterns.

- **Files**: `nanodoc file1.txt file2.md`
- **Directories**: `nanodoc ./docs/` (will include all `.txt` and `.md` files in the directory)
- **Globs**: `nanodoc "src/**/*.go"`

### Line Numbering

- **Per-file numbering**: `-n` or `--line-numbers`
- **Global numbering**: `-N` or `--global-line-numbers`

```bash
# Number lines for each file starting from 1
nanodoc -n file1.txt file2.txt

# Number lines continuously across all files
nanodoc -N file1.txt file2.txt
```

### Table of Contents

Generate a table of contents at the beginning of the document.

```bash
nanodoc --toc file1.md file2.md
```

The TOC is generated from Markdown headings (level 1 and 2) in `.md` files.

### Headers

By default, `nanodoc` adds a "nice" header before the content of each file.

- **Disable headers**: `--no-header`
- **Header style**: `--header-style [nice|filename|path]`
- **Sequence style**: `--sequence [numerical|letter|roman]`

```bash
# Get raw content with no headers
nanodoc --no-header file1.txt

# Use the full file path as the header
nanodoc --header-style path file1.txt

# Use roman numerals for sequence
nanodoc --sequence roman file1.txt file2.txt
```

### Themes

`nanodoc` supports themes for formatting (currently affects headers and other elements, with more to come).

- **Select a theme**: `--theme [classic|classic-dark|classic-light]`

```bash
nanodoc --theme classic-dark file1.txt
```

You can also provide a path to a custom theme file:

```bash
nanodoc --theme /path/to/my-theme.yaml file1.txt
```

### Bundle Files

Bundle files are text files that contain a list of other files to include. The bundle file itself is not included in the output. Bundle files are identified by having `.bundle.` in their name (e.g., `my.bundle.txt`).

**Example `my.bundle.txt`:**

```
# This is a comment
file1.txt
/path/to/another/file.md

# You can include other bundles
another.bundle.txt
```

To use a bundle file, just include it in the path list:

```bash
nanodoc my.bundle.txt
```

## Development

This project is built with Go and uses Cobra for the CLI.

- **Run tests**: `go test ./...`
- **Run linter**: `./scripts/lint`
- **Build**: `go build ./cmd/nanodoc`

## License

MIT