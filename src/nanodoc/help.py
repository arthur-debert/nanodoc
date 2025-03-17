"""Help module for nanodoc."""

import argparse
import glob
import pathlib
import re
import sys
from typing import Dict, Tuple

import yaml
from rich.console import Console
from rich.markdown import Markdown
from rich.panel import Panel
from rich.style import Style
from rich.theme import Theme

from .files import TXT_EXTENSIONS

# Default theme name
DEFAULT_THEME = "neutral"


# Custom help action to use our custom help format
class CustomHelpAction(argparse.Action):
    """Custom action for --help flag to use our custom help format."""

    def __init__(
        self,
        option_strings,
        dest=argparse.SUPPRESS,
        default=argparse.SUPPRESS,
        help=None,
    ):
        super().__init__(
            option_strings=option_strings,
            dest=dest,
            default=default,
            nargs=0,
            help=help,
        )

    def __call__(self, parser, namespace, values, option_string=None):
        print_help()
        parser.exit()


def _get_themes_dir():
    """Return the path to the themes directory."""
    module_dir = pathlib.Path(__file__).parent.absolute()
    return module_dir / "themes"


def _load_theme(theme_name=DEFAULT_THEME):
    """Load a theme from a YAML file.

    Args:
        theme_name: The name of the theme to load.

    Returns:
        Theme: A Rich Theme object.
    """
    themes_dir = _get_themes_dir()
    theme_path = themes_dir / f"{theme_name}.yaml"

    # Fall back to default theme if the requested theme doesn't exist
    if not theme_path.exists():
        theme_path = themes_dir / f"{DEFAULT_THEME}.yaml"

    # Load the theme from YAML
    try:
        with open(theme_path, "r", encoding="utf-8") as f:
            theme_data = yaml.safe_load(f)

        # Convert the YAML data to a Rich Theme
        styles = {}
        for key, value in theme_data.items():
            styles[key] = Style.parse(value)

        return Theme(styles)
    except Exception as e:
        print(f"Error loading theme: {e}")
        # Return a minimal default theme if there's an error
        return Theme(
            {
                "heading": Style(color="blue", bold=True),
                "error": Style(color="red", bold=True),
            }
        )


# Initialize Rich console with the default theme
console = Console(theme=_load_theme())


def _get_docs_dir():
    """Return the path to the docs directory."""
    # Get the directory where this module is located
    module_dir = pathlib.Path(__file__).parent.absolute()
    # The docs directory
    return module_dir / "docs"


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

    # Look for files with extensions from TXT_EXTENSIONS in the guides directory
    for ext in TXT_EXTENSIONS:
        for guide_path in glob.glob(str(guides_dir / f"*{ext}")):
            guide_name = pathlib.Path(guide_path).name.replace(ext, "")

            # Extract the first line as the title/description
            try:
                with open(guide_path, "r", encoding="utf-8") as f:
                    first_line = f.readline().strip()

                    # Remove markdown heading symbols or numbered list markers if present
                    if first_line.startswith("#"):
                        # Remove markdown heading
                        first_line = first_line.lstrip("#").strip()
                    elif re.match(r"^\d+\.", first_line):
                        # Remove numbered list marker (e.g., "1.")
                        first_line = re.sub(r"^\d+\.\s*", "", first_line)

                    # Extract description starting from the first word
                    words = first_line.split()
                    if words:
                        description = " ".join(words)
                    else:
                        description = f"Guide: {guide_name}"

                    guides[guide_name] = description
            except Exception:
                guides[guide_name] = f"Guide: {guide_name}"

    return guides


def get_options_section():
    """Generate the OPTIONS section of the help content.

    Returns:
        str: The formatted OPTIONS section.
    """
    # Import here to avoid circular imports
    from . import nanodoc

    # Generate OPTIONS section
    options_content = ""
    # Get command line options from nanodoc
    for option_str, help_text in nanodoc.get_command_line_options():
        # Add to options content
        options_content += f"{option_str}:{' ' * (20 - len(option_str))} {help_text}\n"

    return options_content


def get_topics_section():
    """Generate the HELP TOPICS section of the help content.

    Returns:
        str: The formatted HELP TOPICS section.
    """
    # Generate HELP TOPICS section
    guides = get_available_guides()
    topics_content = ""
    for name, description in guides.items():
        topics_content += f"{name}:{' ' * (20 - len(name))} {description}\n"

    return topics_content


def _is_rich_content(content: str) -> bool:
    """Check if the content contains Rich markup or should be rendered as Rich.

    Args:
        content: The content to check.

    Returns:
        bool: True if the content contains Rich markup or has a Rich render directive.
    """
    # Check for Rich render directive
    if re.search(r"<!--\s*RENDER:\s*rich\s*-?->", content, re.IGNORECASE):
        return True
    # Check for Rich markup tags like [bold], [italic], etc.
    return bool(
        re.search(
            r"\[(?:bold|italic|red|green|blue|yellow|cyan|magenta|dim|underline)\]",
            content,
        )
    )


def _is_markdown_content(content: str, file_extension: str = None) -> bool:
    """Check if the content should be rendered as Markdown.

    Args:
        content: The content to check.
        file_extension: The file extension, if available.

    Returns:
        bool: True if the content should be rendered as Markdown.
    """
    if file_extension and file_extension.lower() in [".md", ".markdown"]:
        return True
    # Check for Markdown headings
    if re.search(r"^#+ ", content, re.MULTILINE):
        return True
    return False


def _render_content(content: str, guide_name: str = None):
    """Render content using the appropriate Rich formatter.

    Args:
        content: The content to render.
        guide_name: The name of the guide, if applicable.
    """
    if _is_rich_content(content):
        console.print(content)
    elif guide_name and (_is_markdown_content(content) or guide_name.endswith(".md")):
        # Use custom styles for markdown rendering
        md = Markdown(content, code_theme="monokai")
        console.print(md)
    else:
        # For plain text with structure (like manifesto.txt), we'll enhance it with some basic formatting
        # This will make numbered sections, code blocks, and lists look better
        # First, look for section headers (numbered or not)
        # Convert the plain text to a more Markdown-friendly format

        # Convert section headers to Markdown headers
        content = re.sub(
            r"^(\d+(\.\d+)*)\.\s+(.+)$", r"## \1. \3", content, flags=re.MULTILINE
        )

        # Convert indented code blocks to Markdown code blocks
        content = re.sub(
            r"(?m)^( {4}|\t)(.+(?:\n(?:    |\t).+)*)", r"```\n\2\n```", content
        )

        # Convert bullet lists
        content = re.sub(r"^(\s*)-\s+(.+)$", r"\1* \2", content, flags=re.MULTILINE)

        # Render as Markdown
        # Use custom styles for markdown rendering
        md = Markdown(content, code_theme="monokai")
        console.print(md)


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

    # Check for the guide with extensions from TXT_EXTENSIONS
    for ext in TXT_EXTENSIONS:
        guide_path = guides_dir / f"{guide_name}{ext}"
        if guide_path.exists():
            with open(guide_path, "r", encoding="utf-8") as f:
                return True, f.read()

    # Guide not found, prepare error message with available guides
    available_guides = get_available_guides()
    guides_list = "\n".join(
        [f"- {name}: {desc}" for name, desc in available_guides.items()]
    )
    return False, f"Guide '{guide_name}' not found. Available guides:\n\n{guides_list}"


def get_help_content() -> Tuple[bool, str]:
    """Get the content of the help file.

    This function reads the hardcoded help file and dynamically inserts
    the OPTIONS and HELP TOPICS sections.

    Returns:
        Tuple[bool, str]: A tuple containing:
            - Boolean indicating if the help file was found
            - The help content if found, or an error message if not
    """
    docs_dir = _get_docs_dir()

    # Check for the help file with extensions from TXT_EXTENSIONS
    for ext in TXT_EXTENSIONS:
        help_path = docs_dir / f"HELP{ext}"
        if help_path.exists():
            try:
                with open(help_path, "r", encoding="utf-8") as f:
                    content = f.read()

                # Get dynamic content sections
                options_content = get_options_section()
                topics_content = get_topics_section()

                # Replace OPTIONS section
                options_pattern = r"(\[bold\]OPTIONS:\[/bold\]\n\n).*?(\n\n\[bold\]HELP TOPICS:\[/bold\])"
                content = re.sub(
                    options_pattern,
                    r"\1" + options_content + r"\2",
                    content,
                    flags=re.DOTALL,
                )

                # Replace HELP TOPICS section
                topics_pattern = r"(\[bold\]HELP TOPICS:\[/bold\]\n\n).*?(\n\n\n\n\[bold\]EXAMPLES:\[/bold\])"
                content = re.sub(
                    topics_pattern,
                    r"\1" + topics_content + r"\2",
                    content,
                    flags=re.DOTALL,
                )

                return True, content
            except Exception as e:
                return False, f"Error processing help file: {e}"

    return False, "nanodoc help file not found. Please refer to the documentation."


def print_help():
    """Print the help text for nanodoc."""
    found, content = get_help_content()

    if found:
        _render_content(content)
    sys.exit(0)


def print_usage():
    """Print the usage information for nanodoc."""
    parser = argparse.ArgumentParser(
        description="Generate documentation from source code.", prog="nanodoc"
    )
    console.print(parser.format_usage())
    sys.exit(0)


def print_guide(guide_name: str):
    """Print a specific guide.

    Args:
        guide_name: The name of the guide to print.
    """
    found, content = get_guide_content(guide_name)

    if found:
        _render_content(content, guide_name)
        # Exit with status 0 if the guide was found
        sys.exit(0)
    else:
        # Format the error message with Rich
        console.print(
            Panel(
                content,
                title="Guide Not Found",
                border_style="error",
                title_align="center",
                padding=(1, 2),
            )
        )
        # Exit with status 1 if the guide was not found
    sys.exit(1)


def check_help(args):
    """Check if help was requested and handle accordingly.

    Args:
        args: The parsed command-line arguments.
    """
    # Check for theme selection
    theme_name = DEFAULT_THEME
    if hasattr(args, "theme") and args.theme:
        theme_name = args.theme
        # Update the console theme
        console.theme = _load_theme(theme_name)

    # Handle help command before any logging occurs
    if len(sys.argv) >= 2 and sys.argv[-1] == "help":
        # Handle guide-specific help: nanodoc help <guide-name>
        print_help()
        sys.exit(0)
    elif len(sys.argv) >= 3 and "help" in sys.argv:
        # Find the position of "help" in the arguments
        help_index = sys.argv.index("help")
        # Check if there's an argument after "help"
        if help_index < len(sys.argv) - 1:
            guide_name = sys.argv[help_index + 1]
            print_guide(guide_name)
            # This function will exit, so the code below won't be reached

    # Handle general help command
    elif args.help == "help" or (len(sys.argv) == 2 and sys.argv[1] == "help"):
        print_help()

    if not args.sources and args.help is None:
        print_usage()
