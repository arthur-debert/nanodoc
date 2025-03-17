#! /usr/bin/env python3
"""Main module for nanodoc application."""
import argparse
import logging
import os
import pathlib
import sys

from .core import process_all
from .files import TXT_EXTENSIONS, get_files_from_args
from .help import CustomHelpAction
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


# For backward compatibility with tests


def get_available_themes():
    """Get a list of available theme names.

    Returns:
        list: A list of available theme names (without the .yaml extension).
    """
    # Get the directory where this module is located
    module_dir = pathlib.Path(__file__).parent.absolute()
    themes_dir = module_dir / "themes"
    themes = []

    if themes_dir.exists():
        for file in os.listdir(themes_dir):
            if file.endswith(".yaml"):
                themes.append(file.replace(".yaml", ""))

    return themes


def get_command_line_options():
    """Get a list of command line options with their descriptions.

    Returns:
        list: A list of tuples containing (option_string, help_text).
    """
    parser = argparse.ArgumentParser(
        description="Generate documentation from source code.",
        prog="nanodoc",
        formatter_class=argparse.RawTextHelpFormatter,
        add_help=False,  # Disable the default help flag
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
        help="How to format session numbers (numerical, letter, or roman)",
    )
    parser.add_argument(
        "--style",
        choices=["filename", "path", "nice"],
        default="nice",
        help="Header style: nice (default) filename or path",
    )

    parser.add_argument(
        "--txt-ext",
        action="append",
        help="Add extensions to expand for (can be used multiple times)",
        metavar="EXT",
    )
    parser.add_argument("sources", nargs="*", help="Source file(s)")
    parser.add_argument("--version", action="version", version=f"%(prog)s {VERSION}")

    # Add theme selection option
    available_themes = get_available_themes()
    parser.add_argument(
        "--theme",
        choices=available_themes if available_themes else ["neutral"],
        help="Select a theme for rendering help and guide content",
    )

    # Get all actions from the parser
    options = []
    for action in parser._actions:
        # Skip positional arguments and version action
        if not action.option_strings or action.dest == "version":
            continue

        # Format the option string
        option_str = ", ".join(action.option_strings)

        # Get the help text
        help_text = action.help if action.help else ""

        # Add to options list
        options.append((option_str, help_text))

    return options


def parse_args():
    """Parse command-line arguments.

    Returns:
        argparse.Namespace: The parsed arguments with processed values.
    """
    parser = argparse.ArgumentParser(
        description="Generate documentation from source code.",
        prog="nanodoc",
        formatter_class=argparse.RawTextHelpFormatter,
        add_help=False,  # Disable the default help flag
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
        help="Header style: nice  (formatted title), filename or path.",
    )

    parser.add_argument(
        "--txt-ext",
        action="append",
        help="Add additional file extensions to search for (can be used multiple times)",
        metavar="EXT",
    )

    parser.add_argument("sources", nargs="*", help="Source file(s)")
    parser.add_argument("--version", action="version", version=f"%(prog)s {VERSION}")

    # Add theme selection option
    available_themes = get_available_themes()
    parser.add_argument(
        "--theme",
        choices=available_themes if available_themes else ["neutral"],
        help="Select a theme for rendering help and guide content",
    )

    parser.add_argument(
        "-h", "--help", action=CustomHelpAction, help="Show help for command"
    )
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


def main():
    """Main entry point for the nanodoc application."""
    args = parse_args()

    # Import here to avoid circular imports
    from .help import check_help

    try:
        # short circuit for help
        check_help(args)

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
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
