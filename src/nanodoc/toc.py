import logging
import os
from typing import List

from .data import ContentItem, get_content, line_range_to_string
from .formatting import apply_style_to_filename, create_header

logger = logging.getLogger("nanodoc")


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
        current_line += total_lines + 3  # 3 for header line and two blank lines

    # Create TOC with line numbers
    toc = ""
    toc += "\n" + create_header("TOC", sequence=None, style=style) + "\n\n"

    # Format filenames according to header style
    formatted_filenames = {}
    for file_path in file_groups.keys():
        filename = os.path.basename(file_path)
        formatted_name = apply_style_to_filename(filename, style, file_path)
        formatted_filenames[file_path] = formatted_name

    max_filename_length = max(len(name) for name in formatted_filenames.values())

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
