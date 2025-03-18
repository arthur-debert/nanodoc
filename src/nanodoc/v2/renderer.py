"""Rendering for Nanodoc v2.

This module handles the "Rendering" stage of the Nanodoc v2 pipeline.
It takes a Document object and generates the final output string,
including basic concatenation of FileContent content and TOC generation.
"""

import os
import re

from nanodoc.v2.formatter import (
    enhance_rendering,
    format_with_line_numbers,
)
from nanodoc.v2.structures import Document


def render_document(
    document: Document, include_toc: bool = False, include_line_numbers: bool = False
) -> str:
    """Render a Document object to a string.

    This function:
    - Generates a table of contents if requested
    - Concatenates the content of all FileContent objects
    - Adds file separators between non-inlined content
    - Optionally adds line numbers

    Args:
        document: Document object to render
        include_toc: Whether to include a table of contents
        include_line_numbers: Whether to include line numbers

    Returns:
        Rendered document as a string
    """
    rendered_parts = []

    # Generate TOC if requested
    if include_toc:
        toc = generate_toc(document)
        if toc:
            rendered_parts.append(toc)

    # Concatenate content
    prev_original_source = None
    for item in document.content_items:
        # Add file separator if needed
        # Check if not inlined and different from previous
        is_not_inlined = not item.original_source
        different_source = item.filepath != prev_original_source

        if is_not_inlined and different_source:
            # Add a separator if this isn't the first content item
            if rendered_parts and not rendered_parts[-1].endswith("\n\n"):
                rendered_parts.append("\n")

            # Add file header
            file_basename = os.path.basename(item.filepath)
            rendered_parts.append(f"# {file_basename}\n\n")

        # Add the content with optional line numbers
        content_to_add = item.content
        if include_line_numbers:
            # Use the formatter's line numbering function
            content_to_add = format_with_line_numbers(content_to_add)

        rendered_parts.append(content_to_add)

        # Ensure content ends with a newline
        if rendered_parts and not rendered_parts[-1].endswith("\n"):
            rendered_parts.append("\n")

        # Track the source for the next iteration
        prev_original_source = item.original_source or item.filepath

    # Join all parts to create the final content
    plain_content = "".join(rendered_parts)

    # Apply theming if requested
    if hasattr(document, "use_rich_formatting") and document.use_rich_formatting:
        return enhance_rendering(
            plain_content,
            theme_name=document.theme_name,
            use_rich_formatting=document.use_rich_formatting,
        )

    return plain_content


def generate_toc(document: Document) -> str:
    """Generate a table of contents from a Document.

    Args:
        document: Document object to generate TOC for

    Returns:
        Table of contents as a string
    """
    # Extract headings from content
    headings = _extract_headings(document)

    if not headings:
        return ""

    toc_lines = ["# Table of Contents\n\n"]

    for file_path, file_headings in headings.items():
        # Add file entry
        filename = os.path.basename(file_path)
        toc_lines.append(f"- {filename}\n")

        # Add headings for this file
        for heading, _ in file_headings:
            # Indent heading
            toc_lines.append(f"  - {heading}\n")

    toc_lines.append("\n")  # Add blank line after TOC

    # Store TOC data in the document for future reference
    document.toc = headings

    return "".join(toc_lines)


def _extract_headings(document: Document) -> dict[str, list[tuple[str, int]]]:
    """Extract headings from document content.

    Args:
        document: Document object to extract headings from

    Returns:
        Dictionary mapping file paths to lists of (heading_text,
        line_number) tuples
    """
    headings_by_file = {}

    # Markdown heading regex (# Heading)
    heading_pattern = re.compile(r"^(#+)\s+(.+)$", re.MULTILINE)

    for item in document.content_items:
        file_headings = []

        # Use the original source if available, otherwise use the filepath
        file_path = item.original_source or item.filepath

        # Extract headings with line numbers
        lines = item.content.split("\n")
        for i, line in enumerate(lines):
            match = heading_pattern.match(line)
            if match:
                heading_level = len(match.group(1))  # Number of # characters
                heading_text = match.group(2).strip()

                # Only include level 1 and 2 headings
                if heading_level <= 2:
                    file_headings.append((heading_text, i + 1))

        # Store headings for this file if any were found
        if file_headings:
            if file_path in headings_by_file:
                headings_by_file[file_path].extend(file_headings)
            else:
                headings_by_file[file_path] = file_headings

    return headings_by_file


# This function is no longer used directly - formatter.format_with_line_numbers
# is used instead, but keeping it here to avoid breaking tests
def _add_line_numbers(content: str) -> str:
    """Add line numbers to content.

    Args:
        content: Content to add line numbers to

    Returns:
        Content with line numbers
    """
    return format_with_line_numbers(content)
