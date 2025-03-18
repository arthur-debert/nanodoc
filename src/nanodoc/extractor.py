"""File extraction for Nanodoc v2.

This module handles the "Resolving Files" and "Gathering Content" stages of the
Nanodoc v2 pipeline. It takes a list of file paths, creates FileContent objects,
and extracts content based on line ranges.
"""

import os

from nanodoc.structures import FileContent
from nanodoc.utils import apply_ranges, parse_path_and_ranges


def resolve_files(
    file_paths: list[str], bundle_extensions: list[str] = None
) -> list[FileContent]:
    """Resolve a list of file paths to FileContent objects.

    This function:
    - Creates a FileContent object for each file path
    - Determines if a file is a bundle based on its extension
    - Parses range specifiers from file paths (e.g. file.txt:10-20)

    Args:
        file_paths: List of absolute file paths
        bundle_extensions: List of file extensions that should be considered
                           bundles (default: [".ndoc"])

    Returns:
        List of FileContent objects (with empty content at this stage)

    Raises:
        ValueError: If a range specifier is invalid
    """
    if bundle_extensions is None:
        # Default bundle extensions to recognize:
        # - .ndoc (primary extension for v2)
        # - .bundle (any file ending in .bundle)
        # - .bundle.* (any file ending in .bundle.something)
        bundle_extensions = [".ndoc", ".bundle"]

    result = []

    for file_path in file_paths:
        # Parse path and range specifier if present
        path, ranges = parse_path_and_ranges(file_path)

        # Determine if this file is a bundle
        is_bundle = False

        # Check direct extension match
        _, ext = os.path.splitext(path)
        if ext.lower() in bundle_extensions:
            is_bundle = True

        # Check for .bundle.* pattern (e.g., .bundle.txt, .bundle.md)
        basename = os.path.basename(path)
        if ".bundle." in basename.lower():
            is_bundle = True

        # Create FileContent object
        file_content = FileContent(filepath=path, ranges=ranges, is_bundle=is_bundle)

        result.append(file_content)

    return result


def gather_content(file_contents: list[FileContent]) -> list[FileContent]:
    """Read file contents and apply line ranges.

    This function:
    - Reads the content of each file
    - Applies the specified ranges to extract the relevant lines
    - Populates the content field of each FileContent object

    Args:
        file_contents: List of FileContent objects with empty content

    Returns:
        List of FileContent objects with populated content

    Raises:
        FileNotFoundError: If a file cannot be read
        ValueError: If a range is invalid
    """
    result = []

    for file_content in file_contents:
        # Make a copy to avoid modifying the original
        updated_content = FileContent(
            filepath=file_content.filepath,
            ranges=file_content.ranges.copy(),
            is_bundle=file_content.is_bundle,
            original_source=file_content.original_source,
        )

        # Read file content
        try:
            with open(updated_content.filepath, encoding="utf-8") as f:
                lines = f.readlines()
        except FileNotFoundError as e:
            raise FileNotFoundError(
                f"File not found: {updated_content.filepath}"
            ) from e

        # Apply ranges to extract relevant lines
        content = apply_ranges(lines, updated_content.ranges)
        updated_content.content = content

        result.append(updated_content)

    return result
