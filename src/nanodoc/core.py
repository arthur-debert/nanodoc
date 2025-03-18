"""CLI interface for nanodoc v2.

This module provides the bridge between the CLI and the v2 implementation.
"""

import logging
import sys
from typing import Optional

from nanodoc.boot import configure_logging
from nanodoc.document import CircularDependencyError, build_document
from nanodoc.extractor import gather_content, resolve_files
from nanodoc.formatter import apply_theme_to_document
from nanodoc.renderer import render_document
from nanodoc.resolver import resolve_paths

# Initialize logger
logger = logging.getLogger("cli")


def run(
    sources: list[str],
    line_number_mode: Optional[str] = None,
    generate_toc: bool = False,
    theme: Optional[str] = None,
    show_header: bool = True,
    sequence: Optional[str] = None,
    style: Optional[str] = "nice",
    txt_ext: Optional[list[str]] = None,
) -> str:
    """Process source files and generate documentation using v2 implementation.

    Args:
        sources (List[str]): List of source file paths.
        line_number_mode (str): Line numbering mode ('file', 'all', or None).
        generate_toc (bool): Whether to generate a table of contents.
        theme (str): Theme to use for output.
        show_header (bool): Whether to show headers.
        sequence (str): The header sequence type (numerical, letter, roman).
        style (str): The header style (filename, path, nice).
        txt_ext (List[str]): Additional file extensions to process.

    Returns:
        str: The combined content of all files with formatting.

    Raises:
        CircularDependencyError: If a circular dependency is detected
        FileNotFoundError: If a file cannot be found
        ValueError: If there are invalid arguments or parameters
    """
    # Configure logging with default settings
    configure_logging()

    logger.debug("Entering v2 implementation")
    logger.info(f"Processing with v2 implementation: {sources}")

    # Stage 1: Resolve Paths
    logger.debug("Starting path resolution")
    resolved_paths = resolve_paths(sources)
    logger.debug(f"Resolved paths: {resolved_paths}")

    # Stage 2: Resolve Files
    logger.debug("Starting file resolution")
    file_contents = resolve_files(resolved_paths)
    logger.debug(f"Resolved {len(file_contents)} files")

    # Stage 3: Gather Content
    logger.debug("Starting content gathering")
    content_items = gather_content(file_contents)
    logger.debug(f"Gathered {len(content_items)} content items")

    try:
        # Stage 4: Build Document
        logger.debug("Building document")
        document = build_document(content_items)
        logger.debug(f"Built document with {len(document.content_items)} content items")

        # Stage 5: Apply Formatting
        logger.debug(f"Applying theme: {theme}")
        use_rich_formatting = theme is not None
        document = apply_theme_to_document(
            document, theme_name=theme, use_rich_formatting=use_rich_formatting
        )

        # Stage 6: Render Document
        logger.debug("Starting document rendering")
        include_line_numbers = line_number_mode is not None
        include_toc = generate_toc
        rendered_content = render_document(
            document,
            include_toc=include_toc,
            include_line_numbers=include_line_numbers,
        )
        logger.debug(f"Rendered document: {len(rendered_content)} characters")

        return rendered_content
    except CircularDependencyError as e:
        logger.error(f"Circular dependency detected: {str(e)}")
        print(str(e), file=sys.stderr)
        sys.exit(1)
