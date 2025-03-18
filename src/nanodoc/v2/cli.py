"""CLI interface for nanodoc v2.

This module provides the bridge between the CLI and the v2 implementation.
"""

import logging
from typing import Optional

from nanodoc.v2.document import build_document
from nanodoc.v2.extractor import gather_content, resolve_files
from nanodoc.v2.formatter import apply_theme_to_document
from nanodoc.v2.renderer import render_document
from nanodoc.v2.resolver import resolve_paths

# Initialize logger
logger = logging.getLogger("nanodoc")


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
    """
    logger.info(f"Processing with v2 implementation: {sources}")

    # Filter sources to remove any CLI options that might have been passed
    # accidentally. This is a safeguard in case the main CLI didn't properly
    # separate options from sources.
    filtered_sources = [
        source
        for source in sources
        if not source.startswith("--") and not source.startswith("-")
    ]

    if len(filtered_sources) != len(sources):
        invalid_sources = set(sources) - set(filtered_sources)
        logger.warning(
            f"Filtered out {len(sources) - len(filtered_sources)} "
            f"invalid sources: {invalid_sources}"
        )

    if not filtered_sources:
        raise ValueError("No valid source files provided")

    # Stage 1: Resolve Paths
    resolved_paths = resolve_paths(filtered_sources)
    logger.debug(f"Resolved paths: {resolved_paths}")

    # Stage 2: Resolve Files
    file_contents = resolve_files(resolved_paths)
    logger.debug(f"Resolved files: {len(file_contents)} content items")

    # Stage 3: Gather Content
    content_items = gather_content(file_contents)
    logger.debug(f"Gathered content: {len(content_items)} content items")

    # Stage 4: Build Document
    document = build_document(content_items)
    logger.debug(f"Built document with {len(document.content_items)} content items")

    # Stage 5: Apply Formatting
    use_rich_formatting = theme is not None
    document = apply_theme_to_document(
        document, theme_name=theme, use_rich_formatting=use_rich_formatting
    )

    # Stage 6: Render Document
    include_line_numbers = line_number_mode is not None
    include_toc = generate_toc
    rendered_content = render_document(
        document, include_toc=include_toc, include_line_numbers=include_line_numbers
    )
    logger.debug(f"Rendered document, length: {len(rendered_content)}")

    return rendered_content
