"""Legacy functions for Nanodoc v2.

This module contains functions that are maintained for backward compatibility
with previous versions of Nanodoc.
"""

import logging

# Initialize logger
logger = logging.getLogger("document")


def run_content(
    content: str,
    line_number_mode: str = None,
    generate_toc: bool = False,
    theme: str = None,
    show_header: bool = True,
) -> str:
    """Process content and apply formatting.

    Args:
        content (str): Content to process.
        line_number_mode (str): Line numbering mode ('file', 'all', or None).
        generate_toc (bool): Whether to generate a table of contents.
        theme (str): Theme to use for output.
        show_header (bool): Whether to show headers.

    Returns:
        str: The processed content with formatting applied.
    """
    logger.warning("run_content is deprecated and will be removed in the future")
    # Placeholder for backward compatibility
    return content


def run_bundle_directives(content: str) -> str:
    """Process bundle directives in content.

    Args:
        content (str): Content to process.

    Returns:
        str: Content with bundle directives processed.
    """
    logger.warning(
        "run_bundle_directives is deprecated and will be removed in the future"
    )

    # Legacy implementation - placeholder
    lines = content.splitlines()
    result_lines = []

    for line in lines:
        if line.strip().startswith("@inline") or line.strip().startswith("@include"):
            logger.warning(f"Skipping directive: {line.strip()}")
            continue
        result_lines.append(line)

    return "\n".join(result_lines)


def run_inline_directive(content: str, directive: str) -> str:
    """Process an inline directive.

    Args:
        content (str): Content to process.
        directive (str): Directive to process.

    Returns:
        str: Content with directive processed.
    """
    logger.warning(
        "run_inline_directive is deprecated and will be removed in the future"
    )
    # Legacy implementation - placeholder
    return content


def run_include_directive(content: str, directive: str) -> str:
    """Process an include directive.

    Args:
        content (str): Content to process.
        directive (str): Directive to process.

    Returns:
        str: Content with directive processed.
    """
    logger.warning(
        "run_include_directive is deprecated and will be removed in the future"
    )
    # Legacy implementation - placeholder
    return content
