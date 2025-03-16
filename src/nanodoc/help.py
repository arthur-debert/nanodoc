"""
# nanodoc

nanodoc is an ultra-lightweight documentation generator. no frills: concat
multiples files into a single document, adding a title separator.

## FEATURES

- Combine multiple text files
- Title Separator
- Flexible file selection
- [optional] Line Numbers: either per file or global (useful for addressing
  sections)
- [optional] Add table of contents

text files into a single document with formatted headers and optional line
numbering. It can process files provided as arguments or automatically find
`.txt` and `.md` files in the current directory.

## Usage

```bash
nanodoc [options] [file1.txt file2.txt ...]
```

## Specifying Files

nanodoc offers three ways to specify the files you want to bundle:

1. **Explicit File List:** Provide a list of files directly as arguments.

    ```bash
    nanodoc file1.txt file2.md chapter3.txt
    ```

2. **Directory:** Specify a directory, and nanodoc will include all `.txt` and
    `.md` files found within it.

    ```bash
    nanodoc docs/
    ```

3. **Bundle File:** Create a text file (e.g., `bundle.txt`) where each line
    contains the path to a file you want to include. nanodoc will read this
    file and bundle the listed files.

    ```text
    # bundle.txt
    file1.txt
    docs/chapter2.md
    /path/to/another_file.txt
    ```

    ```bash
    nanodoc bundle.txt
    ```

## Options

- `-v, --verbose`: Enable verbose output
- `-n`: Enable per-file line numbering (01, 02, etc.)
- `-nn`: Enable global line numbering (001, 002, etc.)
- `--toc`: Include a table of contents at the beginning
| - `--no-header`: Hide file headers completely
| - `--sequence`: Add sequence numbers to headers
|   - `numerical`: Use numbers (1., 2., etc.)
|   - `letter`: Use letters (a., b., etc.)
|   - `roman`: Use roman numerals (i., ii., etc.)
| - `--style`: Change how filenames are displayed
|   - `filename`: Just the filename
|   - `path`: Full file path
|   - `nice` (default): Formatted title (removes extension, replaces - and _
| |     with spaces, title case, adds original filename in parentheses)
- `-h, --help`: Show this help message

Between files, a separator line is inserted with the format:

```bash
########################## File Name  #########################################
```

The script will exit with an error if no files are found to bundle.

## Examples

```bash
nanodoc -n intro.txt chapter1.txt           # Bundle with per-file numbering
nanodoc -nn --toc                           # Bundle all files with TOC and global numbers
nanodoc --toc -v                            # Verbose bundle with TOC
nanodoc some_directory                      # Add all files in directory
| nanodoc --no-header file1.txt file2.txt   # Hide headers
| nanodoc --sequence=roman file1.txt        # Use roman numerals (i., ii., etc.)
| nanodoc --style=filename file1.txt        # Use filename style instead of nice (default)
nanodoc  bundle_file                        # bundle_file is a txt document with files paths on lines
```
"""

import argparse
import sys


def get_help_text():
    """Return the help text for nanodoc."""
    return __doc__


def print_help():
    """Print the help text for nanodoc."""
    print(get_help_text())
    sys.exit(0)


def print_usage():
    """Print the usage information for nanodoc."""
    parser = argparse.ArgumentParser(
        description="Generate documentation from source code.",
        prog="nanodoc",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.print_usage()
    sys.exit(0)


def check_help(args):
    """Check if help was requested and handle accordingly.

    Args:
        args: The parsed command-line arguments.
    """
    # Handle help command before any logging occurs
    if args.help == "help" or (len(sys.argv) == 2 and sys.argv[1] == "help"):
        print_help()

    if not args.sources and args.help is None:
        print_usage()
