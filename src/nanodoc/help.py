"""Help module for nanodoc."""

import argparse
import pathlib
import sys


def _get_help_file_path():
    """Return the path to the help markdown file."""
    # Get the directory where this module is located
    module_dir = pathlib.Path(__file__).parent
    # The help file is in the same directory
    return module_dir / "HELP.md"


def get_help_text():
    """Return the help text for nanodoc."""
    help_file_path = _get_help_file_path()
    if help_file_path.exists():
        with open(help_file_path, "r", encoding="utf-8") as f:
            return f.read()
    else:
        # Fallback in case the file is not found
        return "nanodoc help file not found. Please refer to the documentation."


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
