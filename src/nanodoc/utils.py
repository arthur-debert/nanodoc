"""Utility functions for Nanodoc v2.

This module contains utility functions that are used across the Nanodoc v2
implementation, providing common functionality for path handling, range
processing, and other generic operations.
"""

import fnmatch
import os
from typing import Optional


def matches_pattern(path: str, pattern: str) -> bool:
    """Simple pattern matching for file paths.

    Args:
        path: File path to check
        pattern: Pattern to match against (e.g., "*.txt")

    Returns:
        True if the path matches the pattern, False otherwise
    """
    # For *.txt pattern, check if file ends with .txt
    if pattern.startswith("*."):
        ext = pattern[1:]  # Extract extension including the dot
        return path.endswith(ext)
    # For specific filename, check exact match
    elif "*" not in pattern and "?" not in pattern and "[" not in pattern:
        return os.path.basename(path) == pattern
    # For more complex patterns, use glob
    else:
        return fnmatch.fnmatch(path, pattern)


def parse_path_and_ranges(
    file_path: str,
) -> tuple[str, list[tuple[int, Optional[int]]]]:
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
                msg = f"Invalid start line number: {start_str}"
                raise ValueError(msg) from e

            # Parse end line
            if end_str.strip():
                try:
                    end = int(end_str.strip())
                except ValueError as e:
                    msg = f"Invalid end line number: {end_str}"
                    raise ValueError(msg) from e
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
            msg = f"End line must be >= start line: {start}-{end}"
            raise ValueError(msg)

        ranges.append((start, end))

    # If no valid ranges were parsed, use default range
    if not ranges:
        return path, default_range

    return path, ranges


def apply_ranges(lines: list[str], ranges: list[tuple[int, Optional[int]]]) -> str:
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

        # Special handling for single line ranges (when start == end)
        # For single line ranges, extract just that line
        if start == end:
            # Check if the line exists
            if start_idx < len(lines):
                # Extract the single line
                result.append(lines[start_idx])
        else:
            # Normal range processing
            # If end is None, use all lines to the end
            # Otherwise, convert end to 0-based indexing and make exclusive
            end_idx = len(lines) if end is None else end - 1
            end_idx = min(end_idx, len(lines))

            # Extract lines for this range
            if start_idx < end_idx:
                result.extend(lines[start_idx:end_idx])

    return "".join(result)
