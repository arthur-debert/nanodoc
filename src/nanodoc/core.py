import logging
import os
from typing import Optional

from .data import (
    ContentItem,
    get_content,
    is_full_file,
    line_range_to_string,
    normalize_line_range,
)
from .formatting import create_header
from .toc import generate_table_of_contents

logger = logging.getLogger("nanodoc")
logger.setLevel(logging.CRITICAL)  # Start with logging disabled


def process_file(
    content_item: ContentItem,
    line_number_mode: Optional[str],
    line_counter: int,
    show_header: bool = True,
    sequence: Optional[str] = None,
    seq_index: int = 0,
    style: Optional[str] = None,
) -> tuple[str, int]:
    """Process a single ContentItem and format its content.

    Args:
        content_item (ContentItem): The ContentItem to process.
        line_number_mode (str): The line numbering mode ('file', 'all',
                               or None).
        line_counter (int): The current global line counter.
        show_header (bool): Whether to show the header.
        sequence (str): The header sequence type (numerical, letter, roman,
                        or None).
        seq_index (int): The index of the file in the sequence.
        style (str): The header style (filename, path, nice, or None).

    Returns:
        tuple: (str, int) Processed file content with header and line
               numbers, and the number of lines in the file.
    """
    logger.debug(
        f"Processing file: {content_item.file_path}, "
        f"line_number_mode: {line_number_mode}, "
        f"line_counter: {line_counter}, "
        f"ranges: {[line_range_to_string(r) for r in content_item.ranges]}"
    )
    try:
        # Get the content from the ContentItem
        get_content(content_item)

        # We need to get all lines to determine the actual line numbers
        with open(content_item.file_path) as f:
            all_lines = f.readlines()

        # Create a list of lines to include with their original line numbers
        lines_with_numbers = []
        for range_obj in content_item.ranges:
            max_lines = len(all_lines)
            start, end = normalize_line_range(range_obj, max_lines)
            for i in range(start - 1, end):
                if i < len(all_lines):
                    lines_with_numbers.append((i + 1, all_lines[i]))

        # Sort by line number to maintain order
        lines_with_numbers.sort(key=lambda x: x[0])
    except FileNotFoundError:
        return f"Error: File not found: {content_item.file_path}\n", 0

    output = ""
    if show_header:
        header = (
            "\n"
            + create_header(
                os.path.basename(content_item.file_path),
                sequence=sequence,
                seq_index=seq_index,
                style=style,
                original_path=content_item.file_path,
            )
            + "\n\n"
        )
        output = header

    for i, (line_num, line) in enumerate(lines_with_numbers):
        line_number = ""
        if line_number_mode == "all":
            line_number = f"{line_counter + i + 1:4d}: "
        elif line_number_mode == "file":
            line_number = f"{line_num:4d}: "
        output += line_number + line

    # Add a blank line if this is a partial content item (not a full file)
    if not (len(content_item.ranges) == 1 and is_full_file(content_item.ranges[0])):
        output += "\n"

    return output, len(lines_with_numbers)


def process_all(
    content_items: list[ContentItem],
    line_number_mode: Optional[str] = None,
    generate_toc: bool = False,
    show_header: bool = True,
    sequence: Optional[str] = None,
    style: Optional[str] = None,
) -> str:
    """Process all ContentItems and combine them into a single document.

    This is the main entry point for both command-line usage and testing.

    Args:
        content_items (list): list of ContentItem objects.
        line_number_mode (str): Line numbering mode ('file', 'all', or None).
        generate_toc (bool): Whether to generate a table of contents.
        show_header (bool): Whether to show headers.
        sequence (str): The header sequence type (numerical, letter, roman,
                        or None).
        style (str): The header style (filename, path, nice, or None).

    Returns:
        str: The combined content of all files with formatting.
    """
    logger.debug(
        f"Processing all files, line_number_mode: {line_number_mode}, "
        f"generate_toc: {generate_toc}"
    )
    output_buffer = ""
    line_counter = 0

    # Group ContentItems by file path
    file_groups = {}
    for item in content_items:
        if item.file_path not in file_groups:
            file_groups[item.file_path] = []
        file_groups[item.file_path].append(item)

    # Generate table of contents if needed
    toc = ""
    if generate_toc:
        toc, _ = generate_table_of_contents(content_items, style)

    # Reset line counter for actual file processing
    line_counter = 0

    # Process each file group
    for file_index, (_, items) in enumerate(file_groups.items()):
        # Process each ContentItem for this file
        for item in items:
            if line_number_mode == "file":
                line_counter = 0

            file_output, num_lines = process_file(
                item,
                line_number_mode,
                line_counter,
                show_header,
                sequence,
                file_index,
                style,
            )
            output_buffer += file_output
            line_counter += num_lines

    if generate_toc:
        output_buffer = toc + output_buffer

    return output_buffer
