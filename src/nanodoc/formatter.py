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

from nanodoc.structures import Document

# Default theme name
DEFAULT_THEME = "classic"

# Initialize logger
logger = logging.getLogger("formatter")


def _get_themes_dir():
    """Return the path to the themes directory."""
    module_dir = pathlib.Path(__file__).parent
    themes_dir = module_dir / "themes"
    logger.debug("Using themes directory: %s", themes_dir)
    return themes_dir


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

    logger.debug("Found available themes: %s", themes)
    return themes


def load_theme(theme_name=DEFAULT_THEME) -> Theme:
    """Load a theme from a YAML file.

    Args:
        theme_name: The name of the theme to load.

    Returns:
        Theme: A Rich Theme object.
    """
    logger.debug("Loading theme: %s", theme_name)
    themes_dir = _get_themes_dir()
    theme_path = themes_dir / f"{theme_name}.yaml"

    # Fall back to default theme if the requested theme doesn't exist
    if not theme_path.exists():
        logger.warning(
            "Theme '%s' not found, falling back to default theme", theme_name
        )
        theme_path = themes_dir / f"{DEFAULT_THEME}.yaml"

    # Load the theme from YAML
    try:
        with open(theme_path, encoding="utf-8") as f:
            theme_data = yaml.safe_load(f)
            logger.debug("Loaded theme data from %s", theme_path)

        # Convert the YAML data to a Rich Theme
        styles = {}
        for key, value in theme_data.items():
            styles[key] = Style.parse(value)
            logger.debug("Parsed style %s: %s", key, value)

        logger.debug("Theme '%s' loaded successfully", theme_name)
        return Theme(styles)
    except Exception as e:
        logger.error("Error loading theme: %s", e)
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
    logger.debug("Creating console with theme: %s", theme_name)

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
    logger.debug(
        "Applying theme to document (theme: %s, rich: %s)",
        theme_name or "default",
        use_rich_formatting,
    )

    if not use_rich_formatting:
        logger.debug("Rich formatting disabled, skipping theme application")
        return document

    # Store theme info in the document for later use
    document.theme_name = theme_name
    document.use_rich_formatting = use_rich_formatting
    logger.debug("Theme information stored in document")

    return document


def format_with_line_numbers(
    content: str, start_number: int = 1, number_format: str = "{:4d}: "
) -> str:
    """Format content with line numbers.

    Args:
        content: The content to format
        start_number: The starting line number
        number_format: The format string for line numbers

    Returns:
        str: Content with line numbers added
    """
    logger.debug(
        "Adding line numbers (start: %d, format: %s)", start_number, number_format
    )

    # Handle empty content case
    if not content:
        return ""

    lines = content.split("\n")
    numbered_lines = []

    for i, line in enumerate(lines):
        line_num = start_number + i
        numbered_lines.append(f"{number_format.format(line_num)}{line}")

    logger.debug("Added line numbers to %d lines", len(lines))
    return "\n".join(numbered_lines)


def enhance_rendering(
    plain_content: str,
    theme_name: Optional[str] = None,
    use_rich_formatting: bool = True,
) -> str:
    """Enhance rendered content with Rich formatting."""
    logger.debug(
        "Enhancing content with theme: %s (rich: %s)",
        theme_name or "default",
        use_rich_formatting,
    )

    if not use_rich_formatting:
        logger.debug("Rich formatting disabled, returning plain content")
        return plain_content

    # Create a string buffer to capture the output
    from io import StringIO

    buffer = StringIO()

    # Create a new console with the theme
    theme = load_theme(theme_name or DEFAULT_THEME)
    console_buffer = Console(file=buffer, theme=theme)

    # Process the content line by line to apply styles
    lines = plain_content.split("\n")
    styled_count = 0

    for line in lines:
        # Apply heading styles
        if line.startswith("# "):
            console_buffer.print(line, style="heading.1")
            styled_count += 1
        elif line.startswith("## "):
            console_buffer.print(line, style="heading.2")
            styled_count += 1
        else:
            console_buffer.print(line)

    logger.debug("Applied styling to %d lines", styled_count)
    return buffer.getvalue()
