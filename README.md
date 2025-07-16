# Nanodoc

A minimalist document bundler designed for stitching hints, reminders and short docs. Useful for prompts, personalized docs highlights for your teams or a note to your future self.

No config, nothing to learn nor remember. Short, simple, sweet.

## 🚀 Installation

### From Source

To build from source, you'll need Go 1.23+ installed.

```bash
git clone https://github.com/arthur-debert/nanodoc-go.git
cd nanodoc-go
go build -o nanodoc ./cmd/nanodoc
```

### Using Go Install

```bash
go install github.com/arthur-debert/nanodoc-go/cmd/nanodoc@latest
```

## ✨ Features

- **Simple file bundling**: Combine multiple text files into one document
- **Bundle files**: Create reusable file lists with `.bundle.*` files
- **Live bundles**: Include files inline using `[[file:path]]` syntax
- **Line ranges**: Extract specific lines with `file.txt:L10-20` syntax
- **Line numbering**: Add line numbers per-file (`-n`) or globally (`-nn`)
- **Table of contents**: Generate TOC with `--toc`
- **Multiple themes**: Built-in themes (classic, classic-light, classic-dark)
- **Smart path resolution**: Handles files, directories, and glob patterns
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
```

### Advanced Usage

```bash
# With line numbers per file
nanodoc -n file1.txt file2.txt

# With global line numbers
nanodoc -nn *.txt

# With table of contents
nanodoc --toc chapter*.md

# With dark theme and no headers
nanodoc --theme=classic-dark --no-header *.md

# Include specific line ranges
nanodoc README.md:L1-10 src/main.go:L20-50
```

### Bundle Files

Create a `.bundle.txt` file to define reusable file lists:

```txt
# My project bundle
README.md
src/main.go
src/utils.go
docs/api.md
```

Bundle files support:
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
- Circular reference detection
- Graceful handling of missing files

## 🛠️ Command Line Options

```
Usage:
  nanodoc [flags] [files/directories/patterns...]

Flags:
  -n, --line-numbers string   Line numbering mode: 'file' (-n) or 'global' (-nn)
      --nn                    Global line numbering (shorthand for -n global)
      --toc                   Generate table of contents
      --theme string          Theme name (classic, classic-light, classic-dark) (default "classic")
      --no-header             Disable file headers
      --sequence string       Header sequence type (numerical, letter, roman) (default "numerical")
      --style string          Header style (nice, filename, path) (default "nice")
      --txt-ext strings       Additional file extensions to process
  -v, --verbose               Enable verbose output
  -h, --help                  Help for nanodoc
      --version               Print version number
```

## 📝 Examples

### 1. Create a Project Overview

```bash
# Create a bundle file
cat > project.bundle.txt << EOF
# Project Overview Bundle
README.md
LICENSE
src/main.go:L1-50
docs/architecture.md
EOF

# Generate the overview
nanodoc --toc project.bundle.txt > project-overview.txt
```

### 2. Code Documentation

```bash
# Document all Go files with line numbers
nanodoc -n --style=filename src/*.go > code-docs.txt
```

### 3. Release Notes with Live Bundles

```txt
# Release v1.2.0

## New Features
[[file:docs/features/auth.md]]
[[file:docs/features/caching.md]]

## Bug Fixes
[[file:CHANGELOG.md:L15-30]]

## Migration Guide
[[file:docs/migration-v1.2.md]]
```

```bash
nanodoc --theme=classic-light release-notes.txt
```

## 🎨 Themes

Nanodoc supports multiple built-in themes:

- **classic**: Default theme with standard formatting
- **classic-light**: Light theme optimized for bright backgrounds
- **classic-dark**: Dark theme optimized for dark backgrounds

Custom themes can be specified by name (implementation pending).

## 🔧 Development

### Building

```bash
# Build all packages
go build -v ./...

# Build the CLI
go build -o nanodoc ./cmd/nanodoc
```

### Testing

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Run linting
golangci-lint run ./...
```

### Project Structure

```
.
├── cmd/nanodoc/          # CLI application
│   ├── main.go          # Main entry point
│   ├── cli.go           # CLI implementation and error handling
│   └── cli_test.go      # CLI tests
├── pkg/nanodoc/         # Core library
│   ├── bundle.go        # Bundle file processing and live bundles
│   ├── constants.go     # Constants and enums
│   ├── errors.go        # Error types and handling
│   ├── extractor.go     # File content extraction and ranges
│   ├── resolver.go      # Path resolution (files, dirs, globs)
│   └── structures.go    # Core data structures
└── docs/dev/            # Development documentation
```

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Run linting (`golangci-lint run ./...`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the need for simple document bundling workflows
- Built with [Cobra](https://github.com/spf13/cobra) for CLI functionality
- Uses [zerolog](https://github.com/rs/zerolog) for structured logging

---

*Nanodoc: Because sometimes you just need to stitch files together, simply.*