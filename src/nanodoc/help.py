"""Help module for nanodoc."""

import argparse
import glob
import pathlib
import sys
from typing import Dict, Tuple


def _get_help_file_path():
    """Return the path to the help markdown file."""
    # Get the directory where this module is located
    module_dir = pathlib.Path(__file__).parent.absolute()
    # The help file is in the docs subdirectory
    return module_dir / "docs" / "HELP.md"


def _get_guides_dir():
    """Return the path to the guides directory."""
    module_dir = pathlib.Path(__file__).parent.absolute()
    return module_dir / "docs" / "guides"


def get_available_guides() -> Dict[str, str]:
    """Return a dictionary of available guides with their descriptions.

    Returns:
        Dict[str, str]: A dictionary mapping guide names to their short descriptions.
    """
    guides = {}
    guides_dir = _get_guides_dir()

    # Look for .txt and .md files in the guides directory
    for ext in [".txt", ".md"]:
        for guide_path in glob.glob(str(guides_dir / f"*{ext}")):
            guide_name = pathlib.Path(guide_path).name.replace(ext, "")

            # Extract the first line as the title/description
            try:
                with open(guide_path, "r", encoding="utf-8") as f:
                    first_line = f.readline().strip()
                    # Remove markdown heading symbols if present
                    description = first_line.lstrip("#").strip()
                    guides[guide_name] = description
            except Exception:
                guides[guide_name] = f"Guide: {guide_name}"

    return guides


def get_guide_content(guide_name: str) -> Tuple[bool, str]:
    """Get the content of a specific guide.

    Args:
        guide_name: The name of the guide to retrieve.

    Returns:
        Tuple[bool, str]: A tuple containing:
            - Boolean indicating if the guide was found
            - The guide content if found, or an error message if not
    """
    guides_dir = _get_guides_dir()

    # Check for the guide with .txt extension
    txt_path = guides_dir / f"{guide_name}.txt"
    if txt_path.exists():
        with open(txt_path, "r", encoding="utf-8") as f:
            return True, f.read()

    # Check for the guide with .md extension
    md_path = guides_dir / f"{guide_name}.md"
    if md_path.exists():
        with open(md_path, "r", encoding="utf-8") as f:
            return True, f.read()

    # Guide not found, prepare error message with available guides
    available_guides = get_available_guides()
    guides_list = "\n".join(
        [f"- {name}: {desc}" for name, desc in available_guides.items()]
    )
    return False, f"Guide '{guide_name}' not found. Available guides:\n\n{guides_list}"


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
        description="Generate documentation from source code.", prog="nanodoc"
    )
    parser.print_usage()
    sys.exit(0)


def print_guide(guide_name: str):
    """Print a specific guide.

    Args:
        guide_name: The name of the guide to print.
    """
    found, content = get_guide_content(guide_name)
    print(content)

    # Exit with status 0 if the guide was found, 1 if not
    if found:
        sys.exit(0)
    sys.exit(1)


def check_help(args):
    """Check if help was requested and handle accordingly.

    Args:
        args: The parsed command-line arguments.
    """
    # Handle help command before any logging occurs
    if len(sys.argv) >= 3 and sys.argv[1] == "help":
        # Handle guide-specific help: nanodoc help <guide-name>
        guide_name = sys.argv[2]
        print_guide(guide_name)
        # This function will exit, so the code below won't be reached

    # Handle general help command
    elif args.help == "help" or (len(sys.argv) == 2 and sys.argv[1] == "help"):
        print_help()

    if not args.sources and args.help is None:
        print_usage()
