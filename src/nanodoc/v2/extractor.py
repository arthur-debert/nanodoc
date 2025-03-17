"""File extraction for Nanodoc v2.

This module handles the "Resolving Files" and "Gathering Content" stages of the
Nanodoc v2 pipeline. It takes a list of file paths, creates FileContent objects,
and extracts content based on line ranges.
"""

import os

from nanodoc.v2.structures import FileContent, Range


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
        bundle_extensions = [".ndoc"]

    result = []

    for file_path in file_paths:
        # Parse path and range specifier if present
        path, ranges = _parse_path_and_ranges(file_path)

        # Determine if this file is a bundle
        _, ext = os.path.splitext(path)
        is_bundle = ext.lower() in bundle_extensions

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
        content = _apply_ranges(lines, updated_content.ranges)
        updated_content.content = content

        result.append(updated_content)

    return result


def _parse_path_and_ranges(file_path: str) -> tuple[str, list[Range]]:
    """Parse a file path with potential range specifiers.

    When specifying a range like 10-20, the end index is exclusive. For example,
    a range of (1, 3) will include lines 1 and 2, but not line 3.

    When specifying a single line like :10, it is treated as just that line.

    Examples:
        "file.txt" -> ("file.txt", [(1, None)])
        "file.txt:10-20" -> ("file.txt", [(10, 20)])  # Lines 10-19 inclusive
        "file.txt:10-" -> ("file.txt", [(10, None)])  # From line 10 to EOF
        "file.txt:10" -> ("file.txt", [(10, 10)])     # Only line 10
        "file.txt:10-20,30-40" -> ("file.txt", [(10, 20), (30, 40)])

    Args:
        file_path: File path potentially containing range specifiers

    Returns:
        Tuple of (clean_path, list_of_ranges)

    Raises:
        ValueError: If a range specifier is invalid
    """
    # Default range (entire file)
    default_range = [(1, None)]

    # Check if there's a range specifier
    if ":" not in file_path:
        return file_path, default_range

    # Split path and range specifier
    parts = file_path.split(":", 1)
    if len(parts) != 2:
        return file_path, default_range

    path, range_spec = parts

    # If range specifier is empty, return default range
    if not range_spec:
        return path, default_range

    # Parse range specifier
    ranges = []
    for range_part in range_spec.split(","):
        range_part = range_part.strip()
        if not range_part:
            continue

        # Parse each range part (e.g., "10-20", "10-", "10")
        if "-" in range_part:
            start_end = range_part.split("-", 1)
            if len(start_end) != 2:
                raise ValueError(f"Invalid range specifier: {range_part}")

            start_str, end_str = start_end

            # Parse start line
            try:
                start = int(start_str.strip())
            except ValueError as e:
                raise ValueError(f"Invalid start line number: {start_str}") from e

            # Parse end line
            if end_str.strip():
                try:
                    end = int(end_str.strip())
                except ValueError as e:
                    raise ValueError(f"Invalid end line number: {end_str}") from e
            else:
                # If end is empty (e.g., "10-"), use None to indicate EOF
                end = None
        else:
            # Single line range (e.g., "10")
            try:
                start = int(range_part.strip())
                # For single line, the test expects (10, 10) not (10, 11)
                end = start
            except ValueError as e:
                raise ValueError(f"Invalid line number: {range_part}") from e

        # Validate range
        if start < 1:
            raise ValueError(f"Line numbers must be positive: {start}")
        if end is not None and end != start and end < start:
            raise ValueError(f"End line must be >= start line: {start}-{end}")

        ranges.append((start, end))

    # If no valid ranges were parsed, use default range
    if not ranges:
        return path, default_range

    return path, ranges


def _apply_ranges(lines: list[str], ranges: list[Range]) -> str:
    """Apply ranges to extract relevant lines from a file.

    Each range is a tuple (start, end) where:
    - start is the 1-indexed starting line (inclusive)
    - end is the 1-indexed ending line (exclusive) or None for EOF

    Args:
        lines: List of lines from the file
        ranges: List of line ranges to extract

    Returns:
        Extracted content as a string

    Raises:
        ValueError: If a range is invalid
    """
    if not lines:
        return ""

    result = []

    for start, end in ranges:
        # Validate range
        if start < 1:
            raise ValueError(f"Line numbers must be positive: {start}")

        # Convert to 0-based indexing
        start_idx = max(0, start - 1)

        # Use ternary operator to determine end_idx
        # For single line, use just that line. Otherwise for normal ranges:
        # - If end is None, use all lines to the end
        # - For normal ranges, subtract 1 to convert to 0-based indexing
        end_idx = start if start == end else (len(lines) if end is None else end - 1)

        end_idx = min(end_idx, len(lines))

        # Extract lines for this range
        if start_idx < end_idx:
            result.extend(lines[start_idx:end_idx])

    return "".join(result)
