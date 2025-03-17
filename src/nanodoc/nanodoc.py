#! /usr/bin/env python3
"""Main module for nanodoc application (legacy wrapper)."""

from .cli import main


def get_command_line_options():
    """Return a list of tuples with option strings and help text.

    Returns:
        list: A list of tuples (option_str, help_text)
    """
    return [
        ("-v", "Enable verbose mode"),
        ("-n", "Enable line number mode (one -n for file, two for all)"),
        ("--toc", "Generate table of contents"),
        ("--no-header", "Hide file headers"),
        ("--sequence", "Add sequence numbers to headers " "(numerical, letter, roman)"),
        ("--style", "Header style: nice, filename or path"),
        ("--txt-ext", "Add additional file extensions to search for"),
        ("--version", "Show the version and exit"),
        ("-h, --help", "Show this message and exit"),
    ]


if __name__ == "__main__":
    main()
