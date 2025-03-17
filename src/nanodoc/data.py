import os
from dataclasses import dataclass
from typing import Optional, Union


@dataclass
class LineRange:
    """A data class representing a range of lines in a file.

    Note: Operations on LineRange objects should be performed using functions,
    not methods. This class is intended to be used as a data holder only.

    Attributes:
        start (int): The start line number (1-indexed).
        end (Union[int, str]): The end line number (1-indexed) or 'X' for end of file.
    """

    start: int
    end: Union[int, str]  # Can be an integer or 'X' for end of file


def is_single_line(line_range: LineRange) -> bool:
    """Check if this range represents a single line."""
    return line_range.start == line_range.end and isinstance(line_range.end, int)


def is_full_file(line_range: LineRange) -> bool:
    """Check if this range represents the entire file."""
    return line_range.start == 1 and line_range.end == "X"


def normalize_line_range(line_range: LineRange, max_lines: int) -> tuple[int, int]:
    """Convert to actual line numbers based on file length.

    Args:
        line_range (LineRange): The line range to normalize.
        max_lines (int): The total number of lines in the file.

    Returns:
        tuple: A tuple of (start, end) line numbers.
    """
    end = max_lines if line_range.end == "X" else line_range.end
    return (line_range.start, end)


def line_range_to_string(line_range: LineRange) -> str:
    """Convert to string representation for display.

    Args:
        line_range (LineRange): The line range to convert to string.

    Returns:
        str: String representation of the line range.
    """
    if is_single_line(line_range):
        return f"L{line_range.start}"
    elif line_range.end == "X":
        return f"L{line_range.start}-X"
    else:
        return f"L{line_range.start}-{line_range.end}"


@dataclass
class ContentItem:
    """A data class representing a file and its line ranges.

    Note: Operations on ContentItem objects should be performed using functions,
    not methods. This class is intended to be used as a data holder only.

    Attributes:
        original_arg (str): The original argument used to specify this content.
        file_path (str): The path to the file.
        ranges (list[LineRange]): A list of LineRange objects.
        content (Optional[str]): The cached content from the file.
    """

    original_arg: str
    file_path: str
    ranges: list[LineRange]
    content: Optional[str] = None


def validate_content_item(content_item: ContentItem) -> bool:
    """Validate that the file exists and ranges are valid.

    Args:
        content_item (ContentItem): The content item to validate.

    Returns:
        bool: True if the content item is valid.

    Raises:
        FileNotFoundError: If the file does not exist.
        PermissionError: If the file is not readable.
        IsADirectoryError: If the path is a directory.
        ValueError: If a line reference is invalid or out of range.
    """
    # Check file existence and readability
    if not os.path.exists(content_item.file_path):
        raise FileNotFoundError(f"File not found: {content_item.file_path}")
    if not os.access(content_item.file_path, os.R_OK):
        raise PermissionError(f"File is not readable: {content_item.file_path}")
    if os.path.isdir(content_item.file_path):
        raise IsADirectoryError(
            f"Path is a directory, not a file: {content_item.file_path}"
        )

    # Validate ranges against file content
    with open(content_item.file_path) as f:
        lines = f.readlines()

    max_lines = len(lines)
    for range_obj in content_item.ranges:
        start, end = normalize_line_range(range_obj, max_lines)
        if start <= 0 or end <= 0 or start > max_lines or end > max_lines:
            raise ValueError(
                f"Line reference out of range: {line_range_to_string(range_obj)} "
                f"(file has {max_lines} lines)"
            )
        if start > end:
            raise ValueError(
                f"Start line must be less than or equal to end line: "
                f"{line_range_to_string(range_obj)}"
            )

    return True


def is_content_item_valid(content_item: ContentItem) -> bool:
    """Check if the content item is valid without raising exceptions.

    Args:
        content_item (ContentItem): The content item to check.

    Returns:
        bool: True if the content item is valid, False otherwise.
    """
    try:
        return validate_content_item(content_item)
    except Exception:
        return False


def load_content(content_item: ContentItem) -> str:
    """Load and cache the content from the file.

    Args:
        content_item (ContentItem): The content item to load content for.

    Returns:
        str: The loaded content.
    """
    if content_item.content is not None:
        return content_item.content

    with open(content_item.file_path) as f:
        all_lines = f.readlines()

    max_lines = len(all_lines)
    result = []

    for range_obj in content_item.ranges:
        start, end = normalize_line_range(range_obj, max_lines)
        result.extend(all_lines[start - 1 : end])

    content_item.content = "".join(result).rstrip("\n")
    return content_item.content


def get_content(content_item: ContentItem) -> str:
    """Get the content, loading it if necessary.

    Args:
        content_item (ContentItem): The content item to get content for.

    Returns:
        str: The content of the file.
    """
    if content_item.content is None:
        return load_content(content_item)
    return content_item.content
