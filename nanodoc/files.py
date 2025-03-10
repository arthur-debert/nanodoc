################################################################################
# Argument Expansion - Functions that turn arguments into a verified list of
# paths
################################################################################


import logging
import os
from typing import List, Tuple

from nanodoc.data import ContentItem, LineRange
from nanodoc.data import get_content as get_item_content
from nanodoc.data import validate_content_item

logger = logging.getLogger("nanodoc")
logger.setLevel(logging.CRITICAL)  # Start with logging disabled


def parse_line_reference(line_ref: str) -> List[LineRange]:
    """Parse a line reference string into a list of LineRange objects.

    Supports 'X' as end marker for the last line of a file.

    Args:
        line_ref (str): The line reference string (e.g., "L5", "L10-20", "L5-X")

    Returns:
        list: A list of LineRange objects

    Raises:
        ValueError: If the line reference is invalid
    """
    if not line_ref:
        raise ValueError("Empty line reference")

    # Check for invalid characters in the line reference
    valid_chars = set("L0123456789,-X")
    for char in line_ref:
        if char not in valid_chars:
            raise ValueError(f"Invalid character in line reference: '{char}'")

    ranges = []
    for part in line_ref.split(","):
        if not part.startswith("L"):
            raise ValueError(f"Invalid line reference format: {part}")

        # Remove the 'L' prefix
        num_part = part[1:]

        if "-" in num_part:
            # Range reference
            try:
                range_parts = num_part.split("-")
                if len(range_parts) != 2:
                    raise ValueError(f"Invalid range format: {part}")

                start = int(range_parts[0])

                # Handle 'X' as end marker
                if range_parts[1].upper() == "X":
                    end = "X"
                else:
                    end = int(range_parts[1])

                if start <= 0:
                    raise ValueError(f"Line numbers must be positive: {part}")
                if isinstance(end, int) and (end <= 0 or start > end):
                    raise ValueError(
                        f"Start line must be less than or equal to end line: {part}"
                    )
                ranges.append(LineRange(start, end))
            except ValueError as e:
                if "must be positive" in str(e) or "must be less than" in str(e):
                    raise
                raise ValueError(f"Invalid line range format: {part}")
        else:
            # Single line reference
            try:
                # Check if num_part contains only digits
                if not num_part.isdigit():
                    raise ValueError(f"Invalid line number format: {part}")
                line_num = int(num_part)
                if line_num <= 0:
                    raise ValueError(f"Line number must be positive: {part}")
                ranges.append(LineRange(line_num, line_num))
            except ValueError as e:
                if "must be positive" in str(e):
                    raise
                raise ValueError(f"Invalid line number format: {part}")

    return ranges


def convert_line_ranges_to_tuples(
    line_ranges: List[LineRange], max_lines: int = None
) -> List[Tuple[int, int]]:
    """Convert a list of LineRange objects to a list of (start, end) tuples.

    This function is used for backward compatibility with code that expects
    the old format of line ranges.

    Args:
        line_ranges (List[LineRange]): A list of LineRange objects
        max_lines (int, optional): The maximum number of lines in the file

    Returns:
        List[Tuple[int, int]]: A list of (start, end) tuples
    """
    if not line_ranges:
        return None

    result = []
    for range_obj in line_ranges:
        if range_obj.end == "X":
            if max_lines is not None:
                result.append((range_obj.start, max_lines))
            else:
                # If max_lines is not provided, just use a large number
                result.append((range_obj.start, 1000000))
        else:
            result.append((range_obj.start, range_obj.end))
    return result


def create_content_item(arg: str) -> ContentItem:
    """Create a ContentItem from a file path and optional line reference.

    Args:
        arg (str): The file path with optional line reference

    Returns:
        ContentItem: A ContentItem object

    Raises:
        ValueError: If the argument format is invalid
    """
    # Check if the path includes a line reference
    file_path = arg
    line_ref = None

    if ":L" in arg:
        parts = arg.split(":", 1)
        if len(parts) != 2:
            raise ValueError(f"Invalid path format: {arg}")

        file_path = parts[0]
        line_spec = parts[1]

        # Ensure the line spec starts with 'L'
        if not line_spec.startswith("L"):
            raise ValueError(
                f"Invalid line reference format: {line_spec} (must start with 'L')"
            )

        line_ref = line_spec

    # Parse line reference or create a full file reference
    if line_ref:
        ranges = parse_line_reference(line_ref)
    else:
        # Full file is represented as L1-X
        ranges = [LineRange(1, "X")]

    return ContentItem(original_arg=arg, file_path=file_path, ranges=ranges)


def verify_content(content_item: ContentItem) -> ContentItem:
    """Verify that a ContentItem is valid.

    Args:
        content_item (ContentItem): The ContentItem to verify

    Returns:
        ContentItem: The verified ContentItem

    Raises:
        FileNotFoundError: If the file does not exist
        PermissionError: If the file is not readable
        IsADirectoryError: If the path is a directory
        ValueError: If a line reference is invalid or out of range
    """
    validate_content_item(content_item)
    return content_item


def get_file_content(file_path, line=None, start=None, end=None, parts=None):
    """Get content from a file, optionally selecting specific lines or ranges.

    Args:
        file_path (str): The path to the file.
        line (int, optional): A specific line number to get (1-indexed).
        start (int, optional): The start line of a range (1-indexed).
        end (int, optional): The end line of a range (1-indexed).
        parts (list, optional): A list of (start, end) tuples representing
                               line ranges.

    Returns:
        str: The selected content from the file.

    Raises:
        FileNotFoundError: If the file does not exist.
        ValueError: If a line reference is out of range.
    """
    # If parts is a list of LineRange objects, convert it to tuples
    if parts and isinstance(parts[0], LineRange):
        with open(file_path, "r") as f:
            max_lines = len(f.readlines())
        parts = convert_line_ranges_to_tuples(parts, max_lines)

    # If we have a ContentItem, use its get_content method
    if isinstance(file_path, ContentItem):
        return get_item_content(file_path)

    try:
        with open(file_path, "r") as f:
            lines = f.readlines()
    except FileNotFoundError:
        raise FileNotFoundError(f"File not found: {file_path}")

    # If no specific parts are requested, return the entire file
    if (
        line is None
        and start is None
        and end is None
        and (parts is None or len(parts) == 0)
    ):
        return "".join(lines)

    # Convert to 0-indexed for internal use
    max_line = len(lines)

    # Handle single line
    if line is not None:
        if line <= 0 or line > max_line:
            raise ValueError(
                "Line reference out of range: " f"{line} (file has {max_line} lines)"
            )
        return lines[line - 1].rstrip("\n")

    # Handle range
    if start is not None and end is not None:
        if start <= 0 or end <= 0 or start > max_line or end > max_line:
            raise ValueError(
                "Line reference out of range: "
                f"{start}-{end} (file has {max_line} lines)"
            )
        # Join the lines and remove the trailing newline if present
        content = "".join(lines[start - 1 : end])
        return content.rstrip("\n")

    # Handle multiple parts
    if parts is not None:
        result = []
        for start, end in parts:
            if start <= 0 or end <= 0 or start > max_line or end > max_line:
                raise ValueError(
                    "Line reference out of range: "
                    f"{start}-{end} (file has {max_line} lines)"
                )
            result.extend(lines[start - 1 : end])
        # Join the lines and remove the trailing newline if present
        content = "".join(result)
        return content.rstrip("\n")

    # Default case
    return "".join(lines)


def expand_directory(directory, extensions=[".txt", ".md"]):
    """Find all files in a directory with specified extensions.

    This function expands a directory path into a list of file paths.

    Args:
        directory (str): The directory path to search.
        extensions (list): List of file extensions to include.

    Returns:
        list: A sorted list of file paths matching the extensions (not validated).
    """
    logger.debug(
        f"Expanding directory with directory='{directory}', "
        f"extensions='{extensions}'"
    )
    matches = []
    for root, _, filenames in os.walk(directory):
        for filename in filenames:
            if any(filename.endswith(ext) for ext in extensions):
                matches.append(os.path.join(root, filename))
    return sorted(matches)


def expand_bundles(bundle_file):
    """Extract list of files from a bundle file.

    This function expands a bundle file into a list of file paths.

    Args:
        bundle_file (str): Path to the bundle file.

    Returns:
        list: A list of file paths contained in the bundle (not validated).

    Raises:
        FileNotFoundError: If bundle file not found or contains no valid files.
    """
    # Extract file path if there's a line reference
    file_path = bundle_file
    line_ref = None
    if ":L" in bundle_file:
        file_path, line_ref = bundle_file.split(":L", 1)
        line_ref = "L" + line_ref

    logger.debug(f"Expanding bundles from file: {bundle_file}")

    try:
        # If there's a line reference, only read the specified lines
        if line_ref:
            try:
                ranges = parse_line_reference(line_ref)
                parts = convert_line_ranges_to_tuples(ranges)
                content = get_file_content(file_path, parts=parts)
                lines = [line.strip() for line in content.splitlines() if line.strip()]
            except ValueError as e:
                raise FileNotFoundError(
                    "Invalid line reference in bundle file: " f"{str(e)}"
                )
        else:
            # Read the entire file
            content = get_file_content(file_path)
            lines = [line.strip() for line in content.splitlines() if line.strip()]
    except FileNotFoundError:
        raise FileNotFoundError(f"Bundle file not found: {file_path}")

    expanded_files = []
    for line in [line for line in lines if line and not line.startswith("#")]:
        expanded_files.append(line)

    # Note: validation is now done separately

    return expanded_files


def is_bundle_file(file_path):
    """Determine if a file is a bundle file by checking its contents.

    A file is considered a bundle if its first non-empty, non-comment line
    points to an existing file.

    Args:
        file_path (str): The path to the file to check.

    Returns:
        bool: True if the file appears to be a bundle file, False otherwise.
    """
    logger.debug(f"Checking if {file_path} is a bundle file")
    try:
        with open(file_path, "r") as f:
            # Check the first few non-empty lines
            for _ in range(5):  # Check up to 5 lines
                line = f.readline().strip()
                if not line:
                    continue
                if line.startswith("#"):  # Skip comment lines
                    continue
                # If this line exists as a file, assume it's a bundle file
                if os.path.isfile(line):
                    return True
                else:
                    # Not a bundle file if a line is not a valid file
                    return False
            # Not a bundle file if none of the first 5 lines are valid files
            return False
    except FileNotFoundError:
        return False
    except Exception as e:
        logger.error(f"Error checking bundle file: {e}")
        return False


def expand_args(args):
    """Expand a list of arguments into a flattened list of file paths.

    This function expands a list of arguments (file paths, directory paths, or
    bundle files) into a flattened list of file paths by calling the appropriate
    expander for each argument.

    Args:
        args (list): A list of file paths, directory paths, or bundle files.

    Returns:
        list: A flattened list of file paths (not validated).
    """
    logger.debug(f"Expanding arguments: {args}")

    def expand_single_arg(arg):
        """Helper function to expand a single argument."""
        logger.debug(f"Expanding argument: {arg}")

        # Extract file path if there's a line reference
        file_path = arg
        if ":L" in arg:
            file_path = arg.split(":L", 1)[0]

        if os.path.isdir(arg):  # Directory path
            return expand_directory(arg)
        elif is_bundle_file(file_path):  # Bundle file
            return expand_bundles(arg)
        else:
            return [arg]  # Regular file path

    # Use list comprehension with sum to flatten the list of lists
    return sum([expand_single_arg(arg) for arg in args], [])


def verify_path(path):
    """Verify that a given path exists, is readable, and is not a directory.

    If the path includes a line reference (e.g., file.txt:L10 or
    file.txt:L5-10), the line reference is validated against the file content.

    Args:
        path (str): The file path to verify.

    Returns:
        str or tuple: If called from older code, returns just the verified path.
                     If called from newer code that expects line parts, returns
                     a tuple (str, list or None) with the verified path and line parts.
                     Line parts is a list of (start, end) tuples or None if no line reference.

    Raises:
        FileNotFoundError: If the path does not exist.
        PermissionError: If the file is not readable.
        IsADirectoryError: If the path is a directory.
        ValueError: If the line reference is invalid or out of range.
    """
    logger.debug(f"Verifying file path: {path}")

    # Check if the path includes a line reference
    file_path = path
    line_ref = None
    line_parts = None

    if ":L" in path:
        parts = path.split(":", 1)
        if len(parts) != 2:
            raise ValueError(f"Invalid path format: {path}")

        file_path = parts[0]
        line_spec = parts[1]

        # Ensure the line spec starts with 'L'
        if not line_spec.startswith("L"):
            raise ValueError(
                f"Invalid line reference format: {line_spec} (must start with 'L')"
            )

        line_ref = line_spec

    # Verify the file path
    if not os.path.exists(file_path):
        raise FileNotFoundError(f"Error: Path does not exist: {file_path}")
    if not os.access(file_path, os.R_OK):
        raise PermissionError(f"Error: File is not readable: {file_path}")
    if os.path.isdir(file_path):
        raise IsADirectoryError(
            "Error: Path is a directory, not a file: " f"{file_path}"
        )

    # Validate line reference if present
    if line_ref:
        try:
            ranges = parse_line_reference(line_ref)
            line_parts = convert_line_ranges_to_tuples(ranges)
            # Validate that all referenced lines exist in the file
            get_file_content(file_path, parts=line_parts)
        except ValueError as e:
            raise ValueError(f"Invalid line reference in {path}: {str(e)}")

    return file_path, line_parts


def file_sort_key(path):
    """Key function for sorting files by name then extension priority."""
    if isinstance(path, ContentItem):
        path = path.file_path
    base_name = os.path.splitext(os.path.basename(path))[0]
    ext = os.path.splitext(path)[1]
    # This ensures test_file.txt comes before test_file.md
    ext_priority = 0 if ext == ".txt" else 1 if ext == ".md" else 2
    return (base_name, ext_priority)


def get_files_from_args(srcs):
    """Process the sources and return a list of ContentItems.

    Args:
        srcs (list): List of source file paths, directories, or bundle files.

    Returns:
        list: A list of ContentItem objects.

    Raises:
        FileNotFoundError: If a file path does not exist.
        PermissionError: If a file is not readable.
        IsADirectoryError: If a path is a directory, not a file.
        ValueError: If a line reference is invalid or out of range.
    """
    # Phase 1: Expand all arguments into a flat list of file paths
    expanded_files = expand_args(srcs)
    if not expanded_files:
        return []

    # Phase 2: Create and validate ContentItems
    content_items = []
    for file_path in expanded_files:
        # Create a ContentItem
        content_item = create_content_item(file_path)
        # Validate the ContentItem
        verify_content(content_item)
        content_items.append(content_item)

    # Sort the content items
    content_items.sort(key=lambda x: file_sort_key(x.file_path))
    return content_items
