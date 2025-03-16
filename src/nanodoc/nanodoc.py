#! /usr/bin/env python3
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
import logging
import sys

from .core import process_all
from .files import TXT_EXTENSIONS, get_files_from_args
from .version import VERSION

LINE_WIDTH = 80


# Custom exception for bundle file errors
class BundleError(Exception):
    """Custom exception for handling errors related to bundle files."""


# Initialize logger at the module level - disabled by default
logger = logging.getLogger("nanodoc")
logger.setLevel(logging.CRITICAL)  # Start with logging disabled

################################################################################
# Sys - System-level functions for logging and output
################################################################################


def setup_logging(to_stderr=False, enabled=False):
    """Configure logging based on requirements.

    Args:
        to_stderr (bool): If True, logs to stderr instead of stdout.
        enabled (bool): If True, sets logging level to DEBUG, otherwise CRITICAL.

    Returns:
        logger: Configured logging object.
    """
    global logger
    if not logger.hasHandlers():  # Only set up logging once
        # Set initial log level
        level = logging.DEBUG if enabled else logging.CRITICAL
        logger.setLevel(level)

        # Create handler to the appropriate stream
        stream = sys.stderr if to_stderr else sys.stdout
        handler = logging.StreamHandler(stream)
        formatter = logging.Formatter(
            "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
        )
        handler.setFormatter(formatter)
        logger.addHandler(handler)
    else:
        # If handlers are already set, just adjust the level
        level = logging.DEBUG if enabled else logging.CRITICAL
        logger.setLevel(level)
    return logger


################################################################################
# Main Processing - Core processing functions
################################################################################


# For backward compatibility with tests


def parse_args():
    """Parse command-line arguments.

    Returns:
        argparse.Namespace: The parsed arguments with processed values.
    """
    parser = argparse.ArgumentParser(
        description="Generate documentation from source code.",
        prog="nanodoc",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument("-v", action="store_true", help="Enable verbose mode")
    parser.add_argument(
        "-n",
        action="count",
        default=0,
        help="Enable line number mode (one -n for file, two for all)",
    )
    parser.add_argument("--toc", action="store_true", help="Generate table of contents")
    parser.add_argument("--no-header", action="store_true", help="Hide file headers")
    parser.add_argument(
        "--sequence",
        choices=["numerical", "letter", "roman"],
        help="Add sequence numbers to headers (numerical, letter, or roman)",
    )
    parser.add_argument(
        "--style",
        choices=["filename", "path", "nice"],
        default="nice",
        help="Header style: nice (default, formatted title), filename (just filename), "
        "or path (full path)",
    )

    parser.add_argument(
        "--txt-ext",
        action="append",
        help="Add additional file extensions to search for (can be used multiple times)",
        metavar="EXT",
    )

    parser.add_argument("sources", nargs="*", help="Source file(s)")
    parser.add_argument("--version", action="version", version=f"%(prog)s {VERSION}")
    parser.add_argument(
        "help",
        nargs="?",
        help="Show program's manual",
        default=None,
        choices=["help"],
    )

    args = parser.parse_args()

    # Process line numbering mode
    if args.n == 0:
        args.line_number_mode = None
    elif args.n == 1:
        args.line_number_mode = "file"
    else:  # args.n >= 2
        args.line_number_mode = "all"

    return args


def _check_help(args):
    # Handle help command before any logging occurs
    if args.help == "help" or (len(sys.argv) == 2 and sys.argv[1] == "help"):
        print(__doc__)
        sys.exit(0)

    if not args.sources and args.help is None:
        parser = argparse.ArgumentParser(
            description="Generate documentation from source code.",
            prog="nanodoc",
            formatter_class=argparse.RawTextHelpFormatter,
        )
        parser.print_usage()
        sys.exit(0)


def main():
    """Main entry point for the nanodoc application."""
    args = parse_args()

    # short circuit for help
    _check_help(args)

    try:
        # Set up logging based on verbose flag
        setup_logging(to_stderr=True, enabled=args.v)

        # Process additional file extensions if provided
        extensions = list(TXT_EXTENSIONS)  # Create a copy of the default extensions
        if args.txt_ext:
            for ext in args.txt_ext:
                # Add a leading dot if not present
                if not ext.startswith("."):
                    ext = "." + ext
                # Only add if not already in the list
                if ext not in extensions:
                    extensions.append(ext)

        # Get verified content items from arguments
        if args.txt_ext:
            # Only pass extensions if custom extensions were provided
            content_items = get_files_from_args(args.sources, extensions=extensions)
        else:
            # Use default extensions
            content_items = get_files_from_args(args.sources)

        # Process the files and print the result
        if not content_items:
            print("Error: No valid source files found.", file=sys.stderr)
            sys.exit(1)

        output = process_all(
            content_items,
            args.line_number_mode,
            args.toc,
            not args.no_header,
            args.sequence,
            args.style,
        )
        print(output)
    except Exception as e:
        print(f"An error occurred: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
