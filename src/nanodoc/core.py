import logging
import os
from typing import List, Optional, Tuple

from .data import (
    ContentItem,
    get_content,
    is_full_file,
    line_range_to_string,
    normalize_line_range,
)
from .formatting import apply_style_to_filename, create_header

logger = logging.getLogger("nanodoc")
logger.setLevel(logging.CRITICAL)  # Start with logging disabled


def generate_table_of_contents(content_items: List[ContentItem], style=None):
    """Generate a table of contents for the given ContentItems.

    Args:
        content_items (list): List of ContentItem objects
        style (str): The header style (filename, path, nice, or None)

    Returns:
        tuple: (str, dict) The table of contents string and a dictionary
               mapping source files to their line numbers in the final document
    """
    logger.debug(f"Generating table of contents for {len(content_items)} items")

    # Calculate line numbers for TOC
    toc_line_numbers = {}
    current_line = 0

    # Calculate the size of the TOC header
    toc_header_lines = 2  # Header line + blank line

    # Group ContentItems by file path
    file_groups = {}
    for item in content_items:
        if item.file_path not in file_groups:
            file_groups[item.file_path] = []
        file_groups[item.file_path].append(item)

    # Calculate the size of each TOC entry (filename + line number)
    # Each file gets one entry, plus one subentry for each range if there are
    # several ranges.
    toc_entries_lines = 0
    for file_path, items in file_groups.items():
        toc_entries_lines += 1  # Main file entry
        if len(items) > 1:
            toc_entries_lines += len(items)  # Subentries for multiple ranges

    # Add blank line after TOC
    toc_footer_lines = 1

    # Total TOC size
    toc_size = toc_header_lines + toc_entries_lines + toc_footer_lines
    current_line = toc_size

    # Calculate line numbers for each file
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

        # Add file lines plus 3 for the header (1 for header, 2 for blank lines)
        current_line += total_lines + 3

    # Create TOC with line numbers
    toc = ""
    toc += "\n" + create_header("TOC", sequence=None, style=style) + "\n\n"

    # Format filenames according to header style
    formatted_filenames = {}
    for file_path in file_groups.keys():
        filename = os.path.basename(file_path)
        formatted_name = apply_style_to_filename(filename, style, file_path)
        formatted_filenames[file_path] = formatted_name

    max_filename_length = max(
        len(formatted_name) for formatted_name in formatted_filenames.values()
    )

    # Add TOC entries
    for file_path, items in file_groups.items():
        formatted_name = formatted_filenames[file_path]
        line_num = toc_line_numbers[file_path]

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

    return toc, toc_line_numbers


def process_file(
    content_item: ContentItem,
    line_number_mode: Optional[str],
    line_counter: int,
    show_header: bool = True,
    sequence: Optional[str] = None,
    seq_index: int = 0,
    style: Optional[str] = None,
) -> Tuple[str, int]:
    """Process a single ContentItem and format its content.

    Args:
        content_item (ContentItem): The ContentItem to process.
        line_number_mode (str): The line numbering mode ('file', 'all', or None).
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
        f"Processing file: {content_item.file_path}, line_number_mode: {line_number_mode}, "
        f"line_counter: {line_counter}, ranges: {[line_range_to_string(r) for r in content_item.ranges]}"
    )
    try:
        # Get the content from the ContentItem
        get_content(content_item)

        # We need to get all lines to determine the actual line numbers
        with open(content_item.file_path, "r") as f:
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
    content_items: List[ContentItem],
    line_number_mode: Optional[str] = None,
    generate_toc: bool = False,
    show_header: bool = True,
    sequence: Optional[str] = None,
    style: Optional[str] = None,
) -> str:
    """Process all ContentItems and combine them into a single document.

    This is the main entry point for both command-line usage and testing.

    Args:
        content_items (list): List of ContentItem objects.
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
    file_index = 0
    for file_path, items in file_groups.items():
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
        file_index += 1

    if generate_toc:
        output_buffer = toc + output_buffer

    return output_buffer
