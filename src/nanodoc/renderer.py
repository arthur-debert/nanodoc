"""Rendering for Nanodoc v2.

This module handles the "Rendering" stage of the Nanodoc v2 pipeline.
It takes a Document object and generates the final output string,
including basic concatenation of FileContent content and TOC generation.
"""

import logging
import os
import re

from nanodoc.formatter import (
    enhance_rendering,
    format_with_line_numbers,
)
from nanodoc.structures import Document

# Initialize logger
logger = logging.getLogger("nanodoc")


def create_header(text: str, style: str = None) -> str:
    """Create a header for a file.

    A simpler V2 implementation that just returns the text as is,
    since we handle styling through the formatter module.

    Args:
        text: The text to use for the header
        style: The style to apply (ignored in V2)

    Returns:
        The header text
    """
    return text


def split_camel_case(s: str) -> str:
    """Split a camel case string into words with spaces.

    Examples:
        wordNice -> word Nice
        WordNice -> Word Nice
        myHTMLFile -> my HTML File

    Args:
        s: Input string possibly in camel case

    Returns:
        String with spaces inserted between camel case words
    """
    # Add space before capital letters that are preceded by lowercase
    s = re.sub(r"([a-z])([A-Z])", r"\1 \2", s)

    # Handle consecutive uppercase followed by lowercase
    # (e.g., HTMLFile -> HTML File)
    s = re.sub(r"([A-Z])([A-Z][a-z])", r"\1 \2", s)

    return s


def get_relative_path(filepath: str) -> str:
    """Get the path relative to the current working directory.

    Args:
        filepath: Absolute or relative file path

    Returns:
        Path relative to current working directory, or the original path
        if it couldn't be made relative
    """
    try:
        # Get current working directory
        cwd = os.getcwd()

        # Convert to absolute paths for comparison
        abs_filepath = os.path.abspath(filepath)

        # Get relative path from cwd to file
        # os.path.relpath handles the case where the file
        # is not under the cwd appropriately
        rel_path = os.path.relpath(abs_filepath, cwd)

        # If rel_path starts with '..' it means the file is not under cwd
        # In that case, return the original absolute path
        if rel_path.startswith(".."):
            return abs_filepath

        return rel_path
    except Exception as e:
        # If there's any error, return the original path
        logger.warning(f"Error getting relative path for {filepath}: {e}")
        return filepath


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
        logger.debug(f"Generating TOC because include_toc={include_toc}")
        toc = generate_toc(document)
        if toc:
            logger.debug(f"TOC generated with length: {len(toc)}")
            rendered_parts.append(toc)
        else:
            logger.debug("No TOC was generated (empty return)")
    else:
        logger.debug(f"Skipping TOC generation because include_toc={include_toc}")

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
            # Use relative path instead of just the basename
            rel_path = get_relative_path(item.filepath)
            file_basename = os.path.basename(item.filepath)

            # Create a nicer format like "Filename (path/to/filename.ext)"
            file_name_without_ext = os.path.splitext(file_basename)[0]
            # Handle word separators: dashes, underscores, and camel case
            nice_name = file_name_without_ext.replace("_", " ").replace("-", " ")
            nice_name = split_camel_case(nice_name)
            nice_name = nice_name.title()
            nice_header = f"{nice_name} ({rel_path})"
            rendered_parts.append(f"\n{nice_header}\n\n")

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
    if hasattr(document, "use_rich_formatting") and (document.use_rich_formatting):
        return enhance_rendering(
            plain_content,
            theme_name=document.theme_name,
            use_rich_formatting=document.use_rich_formatting,
        )

    return plain_content


def generate_toc(document: Document) -> str:
    """Generate a table of contents for a Document.

    Args:
        document: Document object to generate TOC for

    Returns:
        str: Table of contents as a string
    """
    toc_parts = []

    # Add TOC header
    toc_parts.append(create_header("Table of Contents", style="filename"))
    toc_parts.append("\n")

    # Extract headings from content
    headings = _extract_headings(document)

    logger.debug(f"Extracted headings for TOC: {headings}")

    if not headings:
        logger.debug("No headings found for TOC")
        return ""

    for file_path, file_headings in headings.items():
        # Add file entry
        filename = os.path.basename(file_path)
        toc_parts.append(f"- {filename}\n")

        # Add headings for this file
        for heading, _ in file_headings:
            # Indent heading
            toc_parts.append(f"  - {heading}\n")

    toc_parts.append("\n")  # Add blank line after TOC

    # Store TOC data in the document for future reference
    document.toc = headings

    toc_content = "".join(toc_parts)
    logger.debug(f"Generated TOC content: {toc_content}")
    return toc_content


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

        # Check if the file has markdown headings
        has_markdown_headings = False
        for i, line in enumerate(lines):
            match = heading_pattern.match(line)
            if match:
                heading_level = len(match.group(1))  # Number of # characters
                heading_text = match.group(2).strip()

                # Only include level 1 and 2 headings
                if heading_level <= 2:
                    has_markdown_headings = True
                    file_headings.append((heading_text, i + 1))

        # If no markdown headings found, create a pseudo-heading
        # from first non-empty line. Helps with plain text files.
        if not has_markdown_headings and lines:
            # Find the first non-empty line
            for i, line in enumerate(lines):
                line = line.strip()
                if line:
                    # Use the first non-empty line as a heading
                    file_headings.append((line[:40], i + 1))  # Limit to 40 chars
                    break

        # Store headings for this file if any were found
        if file_headings:
            if file_path in headings_by_file:
                headings_by_file[file_path].extend(file_headings)
            else:
                headings_by_file[file_path] = file_headings

    return headings_by_file
