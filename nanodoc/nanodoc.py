#! /usr/bin/env python3
""" 
# nanodocs

nanodocs is an ultra-lightweight documentation generator. no frills: concat
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
nanodocs [options] [file1.txt file2.txt ...]
```

## Specifying Files

nanodocs offers three ways to specify the files you want to bundle:

1. **Explicit File List:** Provide a list of files directly as arguments.

    ```bash
    nanodocs file1.txt file2.md chapter3.txt
    ```

2. **Directory:** Specify a directory, and nanodocs will include all `.txt` and
    `.md` files found within it.

    ```bash
    nanodocs docs/
    ```

3. **Bundle File:** Create a text file (e.g., `bundle.txt`) where each line
    contains the path to a file you want to include. nanodocs will read this
    file and bundle the listed files.

    ```text
    # bundle.txt
    file1.txt
    docs/chapter2.md
    /path/to/another_file.txt
    ```

    ```bash
    nanodocs bundle.txt
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
|   - `nice` (default): Formatted title (removes extension, replaces - and _ with spaces, title case, adds original filename in parentheses)
- `-h, --help`: Show this help message

Between files, a separator line is inserted with the format:

```bash
########################## File Name  #########################################
```

The script will exit with an error if no files are found to bundle.

## Examples

```bash
nanodocs -n intro.txt chapter1.txt           # Bundle with per-file numbering
nanodocs -nn --toc                           # Bundle all files with TOC and global numbers
nanodocs --toc -v                            # Verbose bundle with TOC
nanodocs some_directory                      # Add all files in directory
| nanodocs --no-header file1.txt file2.txt     # Hide headers
| nanodocs --sequence=roman file1.txt        # Use roman numerals (i., ii., etc.)
| nanodocs --style=filename file1.txt   # Use filename style instead of nice (default)
nanodocs  bundle_file                         # bundle_file is a txt docuument with files paths on lines
```

"""
import argparse
import os
import sys
import logging

import re

# Version and configuration constants
VERSION = "0.1.0"
LINE_WIDTH = 80


# Custom exception for bundle file errors
class BundleError(Exception):
    """Custom exception for handling errors related to bundle files."""

    pass


# Initialize logger at the module level - disabled by default
logger = logging.getLogger("nanodoc")
logger.setLevel(logging.CRITICAL)  # Start with logging disabled

################################################################################
# Argument Expansion - Functions that turn arguments into a verified list of paths
################################################################################

def expand_directory(directory, extensions=[".txt", ".md"]):
    """Find all files in a directory with specified extensions.

    This function expands a directory path into a list of file paths.

    Args:
        directory (str): The directory path to search.
        extensions (list): List of file extensions to include.

    Returns:
        list: A sorted list of file paths matching the extensions (not validated).
    """
    logger.debug(
        f"Expanding directory with directory='{directory}', extensions='{extensions}'"
    )
    matches = []
    for root, _, filenames in os.walk(directory):
        for filename in filenames:
            if any(filename.endswith(ext) for ext in extensions):
                matches.append(os.path.join(root, filename))
    return sorted(matches)


def expand_bundles(bundle_file):
    """Extract list of files from a bundle file.

    This function expands a bundle file into a list of file paths.

    Args:
        bundle_file (str): Path to the bundle file.

    Returns:
        list: A list of file paths contained in the bundle (not validated).

    Raises:
        BundleError: If bundle file not found or contains no valid files.
    """
    logger.debug(f"Expanding bundles from file: {bundle_file}")
    try:
        with open(bundle_file, "r") as f:
            lines = [line.strip() for line in f if line.strip()]  # Skip empty lines
    except FileNotFoundError:
        raise BundleError(f"Bundle file not found: {bundle_file}")

    expanded_files = []
    for line in [l for l in lines if l and not l.startswith('#')]:
        expanded_files.append(line)
    
    # Note: validation is now done separately

    return expanded_files


def expand_args(args):
    """Expand a list of arguments into a flattened list of file paths.

    This function expands a list of arguments (file paths, directory paths, or bundle files)
    into a flattened list of file paths by calling the appropriate expander for each argument.

    Args:
        args (list): A list of file paths, directory paths, or bundle files.

    Returns:
        list: A flattened list of file paths (not validated).
    """
    logger.debug(f"Expanding arguments: {args}")

    def expand_single_arg(arg):
        """Helper function to expand a single argument."""
        logger.debug(f"Expanding argument: {arg}")
        if os.path.isdir(arg):  # Directory path
            return expand_directory(arg)
        elif is_bundle_file(arg):  # Bundle file
            return expand_bundles(arg)
        else:
            return [arg]  # Regular file path
    
    # Use list comprehension with sum to flatten the list of lists
    return sum([expand_single_arg(arg) for arg in args], [])

def verify_path(path):
    """Verify that a given path exists, is readable, and is not a directory.

    Args:
        path (str): The file path to verify.

    Returns:
        str: The verified path.

    Raises:
        FileNotFoundError: If the path does not exist.
        PermissionError: If the file is not readable.
        IsADirectoryError: If the path is a directory.
    """
    logger.debug(f"Verifying file path: {path}")
    if not os.path.exists(path):
        raise FileNotFoundError(f"Error: Path does not exist: {path}")
    if not os.access(path, os.R_OK):
        raise PermissionError(f"Error: File is not readable: {path}")
    if os.path.isdir(path):
        raise IsADirectoryError(f"Error: Path is a directory, not a file: {path}")
    return path


def is_bundle_file(file_path):
    """Determine if a file is a bundle file by checking its contents.

    A file is considered a bundle if its first non-empty, non-comment line
    points to an existing file.

    Args:
        file_path (str): The path to the file to check.

    Returns:
        bool: True if the file appears to be a bundle file, False otherwise.
    """
    logger.debug(f"Checking if {file_path} is a bundle file")
    try:
        with open(file_path, "r") as f:
            # Check the first few non-empty lines
            for _ in range(5):  # Check up to 5 lines
                line = f.readline().strip()
                if not line:
                    continue
                if line.startswith("#"):  # Skip comment lines
                    continue
                # If this line exists as a file, assume it's a bundle file
                if os.path.isfile(line):
                    return True
                else:
                    return False  # Not a bundle file if a line is not a valid file
            return (
                False  # Not a bundle file if none of the first 5 lines are valid files
            )
    except FileNotFoundError:
        return False
    except Exception as e:
        logger.error(f"Error checking bundle file: {e}")
        return False


################################################################################
# Formatting - Functions related to headers, line numbers, and table of contents
################################################################################

def apply_style_to_filename(filename, style, original_path=None):
    """Apply the specified style to a filename.

    Args:
        filename (str): The filename to style.
        style (str): The style to apply (filename, path, nice, or None).
        original_path (str, optional): The original file path (used for path and nice styles).

    Returns:
        str: The styled filename.
    """
    logger.debug(f"Applying style '{style}' to filename '{filename}'")

    if not style or style == "filename" or not original_path:
        return filename

    if style == "path":
        # Use the full file path
        return original_path
    elif style == "nice":
        # Remove extension, replace - and _ with spaces, title case, then add filename in parentheses
        basename = os.path.splitext(filename)[0]  # Remove extension

        # Replace - and _ with spaces
        nice_name = re.sub(r"[-_]", " ", basename)

        # Title case
        nice_name = nice_name.title()

        # Add filename in parentheses
        return f"{nice_name} ({filename})"

    # Default to filename if style is not recognized
    return filename


def to_roman(num):
    """Convert integer to roman numeral.

    Args:
        num (int): A positive integer to convert.

    Returns:
        str: Roman numeral representation of the input.
    """
    if not isinstance(num, int) or num <= 0:
        raise ValueError("Input must be a positive integer")

    val = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1]
    syms = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"]

    roman_num = ""
    i = 0
    while num > 0:
        for _ in range(num // val[i]):
            roman_num += syms[i]
            num -= val[i]
        i += 1
    return roman_num.lower()


def format_pos(style, position):
    """Format the sequence prefix based on the sequence type.

    Args:
        style (str): The sequence type (numerical, letter, roman).
        position (int): The index of the item in the sequence.

    Returns:
        str: The formatted sequence prefix.
    """
    if not style:
        return ""

    # Calculate one-indexed number first
    pos_one_indexed = position + 1

    # Dictionary mapping styles to formatting functions
    style_formatters = {
        "numerical": lambda n: f"{int(n)}. ",
        "letter": lambda n: f"{chr(96 + ((n - 1) % 26) + 1)}. ",
        "roman": lambda n: f"{to_roman(n)}. ",
    }

    # Use the appropriate formatter or return empty string if style not found
    return style_formatters.get(style, lambda _: "")(pos_one_indexed)


def apply_sequence_to_text(text, sequence, seq_index):
    """Apply the specified sequence to text."""
    prefix = format_pos(sequence, seq_index)
    return prefix + text if prefix else text


def create_header(
    text, char="#", sequence=None, seq_index=0, style=None, original_path=None
):
    """Create a formatted header with the given text.

    Args:
        text (str): The text to include in the header.
        char (str): The character to use for the header border.
        sequence (str): The header sequence type (numerical, letter, roman, or None).
        seq_index (int): The index of the file in the sequence.
        style (str): The header style (filename, path, nice, or None).
        original_path (str): The original file path (used for path and nice styles).

    Returns:
        str: A formatted header string with the text centered.
    """
    # Apply style to the text if original_path is provided
    if original_path:
        filename = os.path.basename(original_path)
        styled_text = apply_style_to_filename(filename, style, original_path)
    else:
        styled_text = text

    # Apply sequence to the styled text
    header = apply_sequence_to_text(styled_text, sequence, seq_index)
    logger.debug(
        f"Creating header with text='{text}', char='{char}', final: '{header}'"
    )

    return header


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


def to_stds(
    srcs,
    verbose=False,
    line_number_mode=None,
    generate_toc=False,
    show_header=True,
    sequence=None,
    style=None,
):
    """Process sources and return the result as a string.

    This function handles setting up logging and error handling.

    Args:
        srcs (list): List of source file paths, directories, or bundle files.
        verbose (bool): Whether to enable verbose logging.
        line_number_mode (str): Line numbering mode ('file', 'all', or None).
        generate_toc (bool): Whether to generate a table of contents.
        show_header (bool): Whether to show headers.
        sequence (str): The header sequence type (numerical, letter, roman, or None).
        style (str): The header style (filename, path, nice, or None).

    Returns:
        str: The processed output.

    Raises:
        Exception: Any error encountered during processing.
    """
    # Enable logging only when verbose is True
    setup_logging(to_stderr=True, enabled=verbose)
    try:
        result = init(
            srcs,
            verbose,
            line_number_mode,
            generate_toc,
            show_header,
            sequence,
            style,
        )
    except Exception as e:
        raise e

    # Always print the result to stdout
    return result


################################################################################
# Main Processing - Core processing functions
################################################################################

def process_file(
    file_path,
    line_number_mode,
    line_counter,
    show_header=True,
    sequence=None,
    seq_index=0,
    style=None,
):
    """Process a single file and format its content.

    Args:
        file_path (str): The path of the file to process.
        line_number_mode (str): The line numbering mode ('file', 'all', or None).
        line_counter (int): The current global line counter.

        show_header (bool): Whether to show the header.
        sequence (str): The header sequence type (numerical, letter, roman, or None).
        seq_index (int): The index of the file in the sequence.
        style (str): The header style (filename, path, nice, or None).
    Returns:
        tuple: (str, int) Processed file content with header and line numbers,
               and the number of lines in the file.
    """
    logger.debug(
        f"Processing file: {file_path}, line_number_mode: {line_number_mode}, line_counter: {line_counter}"
    )
    try:
        with open(file_path, "r") as f:
            lines = f.readlines()
    except FileNotFoundError:
        return f"Error: File not found: {file_path}\n", 0

    output = ""
    if show_header:
        header = (
            "\n"
            + create_header(
                os.path.basename(file_path),
                sequence=sequence,
                seq_index=seq_index,
                style=style,
                original_path=file_path,
            )
            + "\n\n"
        )
        output = header

    for i, line in enumerate(lines):
        line_number = ""
        if line_number_mode == "all":
            line_number = f"{line_counter + i + 1:4d}: "
        elif line_number_mode == "file":
            line_number = f"{i + 1:4d}: "
        output += line_number + line
    return output, len(lines)


def process_all(
    verified_sources,
    line_number_mode,
    generate_toc,
    show_header=True,
    sequence=None,
    style=None,
):
    """Process all source files and combine them.

    Args:
        verified_sources (list): List of verified source file paths.
        line_number_mode (str): Line numbering mode ('file', 'all', or None).
        generate_toc (bool): Whether to generate a table of contents.

        show_header (bool): Whether to show headers.
        sequence (str): The header sequence type (numerical, letter, roman, or None).
        style (str): The header style (filename, path, nice, or None).
    Returns:
        str: The combined content of all files with formatting.
    """
    logger.debug(
        f"Processing all files, line_number_mode: {line_number_mode}, generate_toc: {generate_toc}"
    )
    output_buffer = ""
    line_counter = 0

    # Custom sort to ensure .txt files come before .md files when base names match
    def file_sort_key(path):
        """Key function for sorting files by name then extension priority."""
        base_name = os.path.splitext(os.path.basename(path))[0]
        ext = os.path.splitext(path)[1]
        # This ensures test_file.txt comes before test_file.md
        ext_priority = 0 if ext == ".txt" else 1 if ext == ".md" else 2
        return (base_name, ext_priority)

    # Sort the verified sources with custom sorting
    verified_sources = sorted(verified_sources, key=file_sort_key)

    # Pre-calculate line numbers for TOC if needed
    toc_line_numbers = {}
    current_line = 0

    if generate_toc:
        # Calculate the size of the TOC header
        toc_header_lines = 2  # Header line + blank line

        # Calculate the size of each TOC entry (filename + line number)
        toc_entries_lines = len(verified_sources)

        # Add blank line after TOC
        toc_footer_lines = 1

        # Total TOC size
        toc_size = toc_header_lines + toc_entries_lines + toc_footer_lines
        current_line = toc_size

        # Calculate line numbers for each file
        for source_file in verified_sources:
            # Add 3 for the file header (1 for the header line, 2 for the blank lines)
            toc_line_numbers[source_file] = current_line + 3
            with open(source_file, "r") as f:
                file_lines = len(f.readlines())
            # Add file lines plus 3 for the header (1 for header line, 2 for blank lines)
            current_line += file_lines + 3

    # Create TOC with line numbers
    toc = ""
    if generate_toc:
        toc += "\n" + create_header("TOC", sequence=None, style=style) + "\n\n"

        # Format filenames according to header style
        formatted_filenames = {}
        for source_file in verified_sources:
            filename = os.path.basename(source_file)
            formatted_filenames[source_file] = apply_style_to_filename(
                filename, style, source_file
            )

        max_filename_length = max(
            len(formatted_name) for formatted_name in formatted_filenames.values()
        )

        for source_file in verified_sources:
            formatted_name = formatted_filenames[source_file]
            line_num = toc_line_numbers[source_file]
            # Format the TOC entry with dots aligning the line numbers
            dots = "." * (max_filename_length - len(formatted_name) + 5)
            toc += f"{formatted_name} {dots} {line_num}\n"

        toc += "\n"

    # Reset line counter for actual file processing
    line_counter = 0

    # Process each file
    for i, source_file in enumerate(verified_sources):
        if line_number_mode == "file":
            line_counter = 0
        file_output, num_lines = process_file(
            source_file,
            line_number_mode,
            line_counter,
            show_header,
            sequence,
            i,
            style,
        )
        output_buffer += file_output
        line_counter += num_lines

    if generate_toc:
        output_buffer = toc + output_buffer

    return output_buffer


def init(
    srcs,
    verbose=False,
    line_number_mode=None,
    generate_toc=False,
    show_header=True,
    sequence=None,
    style=None,
):
    """Initialize and process the sources.

    Args:
        srcs (list): List of source file paths, directories, or bundle files.
        verbose (bool): Whether to enable verbose logging.
        line_number_mode (str): Line numbering mode ('file', 'all', or None).
        generate_toc (bool): Whether to generate a table of contents.
        show_header (bool): Whether to show headers.
        sequence (str): The header sequence type (numerical, letter, roman, or None).
        style (str): The header style (filename, path, nice, or None).

    Returns:
        str: The processed output.
    """
    logger.debug(
        f"Initializing with sources: {srcs}, verbose: {verbose}, line_number_mode: {line_number_mode}, generate_toc: {generate_toc}"
    )
    
    # Phase 1: Expand all arguments into a flat list of file paths
    expanded_files = expand_args(srcs)
    
    if not expanded_files:
        return "Error: No source files found."
    
    # Phase 2: Validate all file paths
    verified_sources = [verify_path(file_path) for file_path in expanded_files]
    if not verified_sources:
        return "Error: No valid source files found."

    output = process_all(
        verified_sources,
        line_number_mode,
        generate_toc,
        show_header,
        sequence,
        style,
    )
    return output


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Generate documentation from source code.",
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
        help="Header style: nice (default, formatted title), filename (just filename), or path (full path)",
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

    # Handle help command before any logging occurs
    if args.help == "help" or (len(sys.argv) == 2 and sys.argv[1] == "help"):
        print(__doc__)
        sys.exit(0)

    if not args.sources and args.help is None:
        parser.print_usage()
        sys.exit(0)

    line_number_mode = None
    if args.n == 1:
        line_number_mode = "file"
    elif args.n >= 2:
        line_number_mode = "all"

    try:
        output = to_stds(
            srcs=args.sources,
            verbose=args.v,
            line_number_mode=line_number_mode,
            generate_toc=args.toc,
            show_header=not args.no_header,
            sequence=args.sequence,
            style=args.style,
        )
        print(output)
    except Exception as e:
        print(f"An error occurred: {e}")
        sys.exit(1)
