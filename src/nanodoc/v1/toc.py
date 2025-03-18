import logging
import os

from .data import ContentItem, get_content, line_range_to_string
from .formatting import apply_style_to_filename, create_header

logger = logging.getLogger("nanodoc")


def group_content_items_by_file(
    content_items: list[ContentItem],
) -> dict[str, list[ContentItem]]:
    """Group ContentItems by file path.

    Args:
        content_items (list[ContentItem]): list of ContentItem objects

    Returns:
        dict[str, list[ContentItem]]: Dictionary mapping file paths to lists of
            ContentItems
    """
    file_groups = {}
    for item in content_items:
        if item.file_path not in file_groups:
            file_groups[item.file_path] = []
        file_groups[item.file_path].append(item)
    return file_groups


def calculate_toc_size(file_groups: dict[str, list[ContentItem]]) -> int:
    """Calculate the total size of the table of contents.

    Args:
        file_groups (dict[str, list[ContentItem]]): Dictionary mapping file
            paths to ContentItems

    Returns:
        int: Total number of lines in the TOC
    """
    # Header line + blank line
    toc_header_lines = 2

    # Calculate the size of each TOC entry
    toc_entries_lines = 0
    for _, items in file_groups.items():
        toc_entries_lines += 1  # Main file entry
        if len(items) > 1:
            toc_entries_lines += len(items)  # Subentries for multiple ranges

    # Add blank line after TOC
    toc_footer_lines = 1

    # Total TOC size
    return toc_header_lines + toc_entries_lines + toc_footer_lines


def calculate_line_numbers(
    file_groups: dict[str, list[ContentItem]], toc_size: int
) -> dict[str, int]:
    """Calculate line numbers for each file in the final document.

    Args:
        file_groups (dict[str, list[ContentItem]]): Dictionary mapping file
            paths to ContentItems
        toc_size (int): Size of the table of contents in lines

    Returns:
        dict[str, int]: Dictionary mapping file paths to their line numbers
    """
    toc_line_numbers = {}
    current_line = toc_size

    for file_path, items in file_groups.items():
        # Add 3 for the file header (1 for the header line, 2 for blank lines)
        toc_line_numbers[file_path] = current_line + 3

        # Calculate total content lines
        total_lines = 0
        for item in items:
            content = get_content(item)
            file_lines = len(content.splitlines())
            total_lines += file_lines
            # Add a blank line between ranges if there are multiple ranges
            if len(items) > 1:
                total_lines += 1

        # Add file lines plus header (header line and two blank lines)
        current_line += total_lines + 3

    return toc_line_numbers


def format_filenames(
    file_groups: dict[str, list[ContentItem]], style=None
) -> dict[str, str]:
    """Format filenames according to the specified style.

    Args:
        file_groups (dict[str, list[ContentItem]]): Dictionary mapping file
            paths to ContentItems
        style (str): The header style (filename, path, nice, or None)

    Returns:
        dict[str, str]: Dictionary mapping file paths to formatted filenames
    """
    formatted_filenames = {}
    for file_path in file_groups:
        filename = os.path.basename(file_path)
        formatted_name = apply_style_to_filename(filename, style, file_path)
        formatted_filenames[file_path] = formatted_name
    return formatted_filenames


def create_toc_content(
    file_groups: dict[str, list[ContentItem]],
    formatted_filenames: dict[str, str],
    line_numbers: dict[str, int],
    style=None,
) -> str:
    """Create the table of contents content.

    Args:
        file_groups (dict[str, list[ContentItem]]): Dictionary mapping file
            paths to ContentItems
        formatted_filenames (dict[str, str]): Dictionary mapping file paths to
            formatted filenames
        line_numbers (dict[str, int]): Dictionary mapping file paths to line
            numbers
        style (str): The header style (filename, path, nice, or None)

    Returns:
        str: The table of contents string
    """
    toc = ""
    toc += (
        "\n" + create_header("Table of Contents", sequence=None, style=style) + "\n\n"
    )

    max_filename_length = max(len(name) for name in formatted_filenames.values())

    # Add TOC entries
    for file_path, items in file_groups.items():
        formatted_name = formatted_filenames[file_path]
        line_num = line_numbers[file_path]

        # Format the TOC entry with dots aligning the line numbers
        dots = "." * (max_filename_length - len(formatted_name) + 5)
        toc += f"{formatted_name} {dots} {line_num}\n"

        # Add subentries for ranges if there are multiple ranges
        if len(items) > 1:
            for i, item in enumerate(items):
                range_info = []
                for range_obj in item.ranges:
                    range_info.append(line_range_to_string(range_obj))
                range_str = ", ".join(range_info)

                # Indent the subentry and use a letter index (a, b, c, ...)
                toc += f"    {chr(97 + i)}. {range_str}\n"

    toc += "\n"
    return toc


def generate_table_of_contents(content_items: list[ContentItem], style=None):
    """Generate a table of contents for the given ContentItems.

    Args:
        content_items (list): list of ContentItem objects
        style (str): The header style (filename, path, nice, or None)

    Returns:
        tuple: (str, dict) The table of contents string and a dictionary
               mapping source files to their line numbers in the final document
    """
    logger.debug(f"Generating table of contents for {len(content_items)} items")

    # Group ContentItems by file path
    file_groups = group_content_items_by_file(content_items)

    # Calculate the size of the TOC
    toc_size = calculate_toc_size(file_groups)

    # Calculate line numbers for each file
    line_numbers = calculate_line_numbers(file_groups, toc_size)

    # Format filenames according to header style
    formatted_filenames = format_filenames(file_groups, style)

    # Create TOC with line numbers
    toc = create_toc_content(file_groups, formatted_filenames, line_numbers, style)

    return toc, line_numbers
