"""CLI interface for nanodoc v2.

This module provides the bridge between the CLI and the v2 implementation.
"""

import logging
import os
import sys
from typing import Optional

from nanodoc.v2.document import CircularDependencyError, build_document
from nanodoc.v2.extractor import gather_content, resolve_files
from nanodoc.v2.formatter import apply_theme_to_document
from nanodoc.v2.renderer import render_document
from nanodoc.v2.resolver import resolve_paths

# Initialize logger
logger = logging.getLogger("cli")


def configure_logging(verbose: bool = False):
    """Configure logging based on environment and CLI flags."""
    log_level = logging.WARNING  # Default level

    # Check environment variable or verbose flag
    if os.environ.get("NANODOC_VERBOSE") or verbose:
        log_level = logging.DEBUG

    # Configure root logger
    logging.basicConfig(
        level=log_level,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )

    # Ensure all our module loggers are set to the same level
    module_names = ["cli", "document", "formatter", "renderer", "resolver", "extractor"]
    for name in module_names:
        module_logger = logging.getLogger(name)
        module_logger.setLevel(log_level)

    logger.debug(
        "Logging configured with level: %s",
        "DEBUG" if log_level == logging.DEBUG else "WARNING",
    )


def process_v2(
    sources: list[str],
    line_number_mode: Optional[str] = None,
    generate_toc: bool = False,
    theme: Optional[str] = None,
    show_header: bool = True,
) -> str:
    """Process files using the v2 implementation.

    Args:
        sources: List of file paths or globs to process
        line_number_mode: Line numbering mode ("file", "all", or None)
        generate_toc: Whether to generate a table of contents
        theme: Theme name to use for rendering
        show_header: Whether to show file headers

    Returns:
        str: The processed document content

    Raises:
        CircularDependencyError: If a circular dependency is detected
        FileNotFoundError: If a file cannot be found
        ValueError: If there are invalid arguments or parameters
    """
    logger.debug("Entering process_v2")
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
