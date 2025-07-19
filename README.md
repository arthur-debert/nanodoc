# Nanodoc

A minimalist document bundler that scales down.

## Why nanodoc?

Most documentation tools focus on large, monolithic documents. Nanodoc takes a different approach: it encourages small, focused documentation fragments that live alongside your code. When you need a complete document, nanodoc stitches these fragments together.

## Quick Start

Nanodoc is simple - you pass it files, and it bundles them:

```bash
# Bundle specific files
$ nanodoc formats/json/about.txt formats/yaml/about.txt

# Extract specific lines
$ nanodoc formats/json/about.txt:L10-30 formats/yaml/about.txt:L4-10,15-30,40

# Add headers and table of contents
$ nanodoc --filenames --toc formats/json/*.txt formats/yaml/*.txt

# Bundle all text files in a directory
$ nanodoc formats/

# Include other file types
$ nanodoc formats/ --ext rst --ext go

# Fine-grained control with patterns
$ nanodoc src/ --ext go --include "**/http/*.go" --exclude "**/*_test.go"
```

## Reusable Bundles

Instead of remembering complex commands, save them as bundles:

```bash
# Save your command as a bundle
$ nanodoc src/ --ext go --toc --save-to-bundle api-docs.bundle.txt

# Later, just run the bundle
$ nanodoc api-docs.bundle.txt
```

A bundle file looks like this:

```
# API Documentation Bundle
--toc
--linenum file

docs/intro.txt
docs/design.txt:L20-200
src/**/*.go
```

## Key Features

- **Line ranges**: Extract specific lines with `file.txt:L10-20`
- **Glob patterns**: Use `**/*.go` to match files recursively
- **Live includes**: Embed files inline with `[[file:path]]` syntax
- **Themes**: Built-in themes for different contexts
- **Line numbering**: Per-file or global line numbers
- **Pattern filtering**: Include/exclude files with gitignore-style patterns

## Installation

```bash
brew install nanodoc
```

Or download a `.deb` package from the [releases](https://github.com/arthur-debert/nanodoc/releases).

## Learn More

For detailed documentation on any topic:

```bash
$ nanodoc help <topic>
```

Available topics include:
- `bundles` - Creating and managing bundle files
- `content` - File selection, patterns, and line ranges
- `filenames` - Customizing file headers and formatting
- `themes` - Available themes and styling options
- `toc` - Generating tables of contents

## Contributing

Contributions are welcome! From feedback, to bug reports or actual code, are all very welcomed.
Use [Github issues](https://github.com/arthur-debert/nanodoc/issues)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

*Nanodoc: Because sometimes you just need to stitch files together, simply.*