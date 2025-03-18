"""Formatting and theming for Nanodoc v2.

This module handles the "Formatting" stage of the Nanodoc v2 pipeline.
It provides theming capabilities and formatting options for the rendered output.
"""

import logging
import os
import pathlib
from typing import Optional

import yaml
from rich.console import Console
from rich.style import Style
from rich.theme import Theme

from nanodoc.v2.structures import Document

# Default theme name
DEFAULT_THEME = "neutral"

# Initialize logger
logger = logging.getLogger("nanodoc")


def _get_themes_dir():
    """Return the path to the themes directory."""
    module_dir = pathlib.Path(__file__).parent.parent.absolute()
    return module_dir / "themes"


def get_available_themes() -> list[str]:
    """Get a list of available theme names.

    Returns:
        list[str]: A list of available theme names (without .yaml extension).
    """
    themes_dir = _get_themes_dir()
    themes = []

    if themes_dir.exists():
        for file in os.listdir(themes_dir):
            if file.endswith(".yaml"):
                themes.append(file.replace(".yaml", ""))

    logger.debug(f"Available themes: {themes}")
    return themes


def load_theme(theme_name=DEFAULT_THEME) -> Theme:
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
        logger.warning(f"Theme '{theme_name}' not found, using default theme")
        theme_path = themes_dir / f"{DEFAULT_THEME}.yaml"

    # Load the theme from YAML
    try:
        with open(theme_path, encoding="utf-8") as f:
            theme_data = yaml.safe_load(f)

        # Convert the YAML data to a Rich Theme
        styles = {}
        for key, value in theme_data.items():
            styles[key] = Style.parse(value)

        logger.debug(f"Theme '{theme_name}' loaded successfully")
        return Theme(styles)
    except Exception as e:
        logger.error(f"Error loading theme: {e}")
        # Return a minimal default theme if there's an error
        return Theme(
            {
                "heading": Style(color="blue", bold=True),
                "error": Style(color="red", bold=True),
            }
        )


def create_themed_console(theme_name=None) -> Console:
    """Create a Rich console with the specified theme.

    Args:
        theme_name: The name of the theme to use. If None, uses default theme.

    Returns:
        Console: A Rich Console object with the specified theme.
    """
    if theme_name is None:
        theme_name = DEFAULT_THEME

    theme = load_theme(theme_name)
    return Console(theme=theme)


def apply_theme_to_document(
    document: Document,
    theme_name: Optional[str] = None,
    use_rich_formatting: bool = True,
) -> Document:
    """Apply theme styling to a document.

    This function adds styling information to the document for later rendering.
    If Rich formatting is not used, the document is returned unchanged.

    Args:
        document: The document to apply theming to
        theme_name: The name of the theme to use, or None for default
        use_rich_formatting: Whether to use Rich for formatting

    Returns:
        Document: The document with theming information
    """
    if not use_rich_formatting:
        return document

    # Store theme info in the document for later use
    document.theme_name = theme_name
    document.use_rich_formatting = use_rich_formatting

    return document


def format_with_line_numbers(
    content: str, start_number: int = 1, number_format: str = "{:4d} | "
) -> str:
    """Format content with line numbers.

    Args:
        content: The content to format
        start_number: The starting line number
        number_format: The format string for line numbers

    Returns:
        str: Content with line numbers added
    """
    lines = content.split("\n")
    numbered_lines = []

    for i, line in enumerate(lines):
        line_num = start_number + i
        numbered_lines.append(f"{number_format.format(line_num)}{line}")

    return "\n".join(numbered_lines)


def enhance_rendering(
    plain_content: str,
    theme_name: Optional[str] = None,
    use_rich_formatting: bool = True,
) -> str:
    """Enhance rendered content with Rich formatting.

    Args:
        plain_content: Plain text content to enhance
        theme_name: Theme to use for styling
        use_rich_formatting: Whether to use Rich for formatting

    Returns:
        str: Enhanced content with Rich formatting
    """
    if not use_rich_formatting:
        return plain_content

    # Create a console with the specified theme
    console = create_themed_console(theme_name)

    # Create a string buffer to capture the output
    from io import StringIO

    buffer = StringIO()
    console_buffer = Console(file=buffer, theme=console.theme)

    # Process the content line by line to apply styles
    lines = plain_content.split("\n")
    for line in lines:
        # Apply heading styles
        if line.startswith("# "):
            console_buffer.print(line, style="heading.1")
        elif line.startswith("## "):
            console_buffer.print(line, style="heading.2")
        # Add more styling rules as needed
        else:
            console_buffer.print(line)

    return buffer.getvalue()
