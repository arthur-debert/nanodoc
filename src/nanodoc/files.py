################################################################################
# Argument Expansion - Functions that turn arguments into a verified list of
# paths
################################################################################


import logging
import os
import re

from .data import ContentItem, LineRange, validate_content_item

# Define text file extensions
TXT_EXTENSIONS = [".txt", ".md", "txxt"]

# Error message templates
ERR_INVALID_LINE_REF_FORMAT = "Invalid line reference format: {} (must start with 'L')"
ERR_INVALID_RANGE = "Invalid range {}-{} (1-{})"
ERR_INLINE_FILE_READ = "Error reading inline file {}: {}"

logger = logging.getLogger("nanodoc")
logger.setLevel(logging.CRITICAL)  # Start with logging disabled


def parse_line_reference(line_ref: str) -> list[LineRange]:
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

    # Split by commas and process each part
    parts = line_ref.split(",")
    for part in parts:
        # Use a single regex pattern with an optional group for the range part
        # Format: L<number>(-<number or X>)?
        match = re.match(r"^L(\d+)(?:-(\d+|X))?$", part, re.IGNORECASE)
        if not match:
            raise ValueError(f"Invalid line reference format: {part}")

        # Extract start line number (always present)
        start = int(match.group(1))
        if start <= 0:
            raise ValueError(f"Line numbers must be positive: {part}")

        # Extract end line number (optional)
        if match.group(2):
            # This is a range reference
            end_str = match.group(2).upper()
            end = end_str if end_str == "X" else int(end_str)

            # Validate end line number if it's an integer
            if isinstance(end, int) and (end <= 0 or start > end):
                raise ValueError(
                    f"Start line must be less than or equal to end line: {part}"
                )
        else:
            # This is a single line reference
            end = start

        ranges.append(LineRange(start, end))

    return ranges


def convert_line_ranges_to_tuples(
    line_ranges: list[LineRange], max_lines: int = None
) -> list[tuple[int, int]]:
    """Convert a list of LineRange objects to a list of (start, end) tuples.

    This function is used for backward compatibility with code that expects
    the old format of line ranges.

    Args:
        line_ranges (list[LineRange]): A list of LineRange objects
        max_lines (int, optional): The maximum number of lines in the file

    Returns:
        list[tuple[int, int]]: A list of (start, end) tuples
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
            raise ValueError(ERR_INVALID_LINE_REF_FORMAT.format(line_spec))

        line_ref = line_spec

    # Parse line reference or create a full file reference
    ranges = parse_line_reference(line_ref) if line_ref else [LineRange(1, "X")]

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

    This function handles getting content from regular files with options to select
    specific lines or ranges. It can handle single lines, line ranges, or multiple
    parts of a file.

    Args:
        file_path (str): The path to the file.
        line (int, optional): A specific line number to get (1-indexed).
        start (int, optional): The start line of a range (1-indexed).
        end (int, optional): The end line of a range (1-indexed).
        parts (list, optional): A list of (start, end) tuples or LineRange
            objects representing line ranges.

    Examples:
        get_file_content("file.txt", line=5)  # Get line 5
        get_file_content("file.txt", start=10, end=20)  # Get lines 10-20
    Returns:
        str: The selected content from the file.

    Raises:
        FileNotFoundError: If the file does not exist.
        ValueError: If line references are out of range.
    """
    # Parameter validation
    if not file_path:
        raise ValueError("File path cannot be empty")

    # Check for conflicting parameters
    param_count = sum(x is not None for x in [line, start and end, parts])
    if param_count > 1:
        raise ValueError(
            "Cannot specify multiple line selection methods simultaneously"
        )

    # Open the file once and handle potential errors
    with open(file_path) as f:
        # For single line or small range requests, we can optimize by not
        # reading the entire file at once
        if line is not None or (start is not None and end is not None):
            lines = f.readlines()
            num_lines = len(lines)

            if line is not None:  # Single line
                if not 1 <= line <= num_lines:
                    raise ValueError(f"Line {line} out of range (1-{num_lines})")
                return lines[line - 1].rstrip("\n")

            elif start is not None and end is not None:  # Line range
                if (
                    not 1 <= start <= num_lines
                    or not 1 <= end <= num_lines
                    or start > end
                ):
                    raise ValueError(ERR_INVALID_RANGE.format(start, end, num_lines))
                return "".join(lines[start - 1 : end]).rstrip("\n")

        # Handle parts parameter
        elif parts:
            # Handle conversion of LineRange objects to tuples
            if parts and isinstance(parts[0], LineRange):
                # Get line count without reading all lines into memory if possible
                f.seek(0)
                num_lines = sum(1 for _ in f)
                parts = convert_line_ranges_to_tuples(parts, num_lines)
                f.seek(0)  # Reset file position

            # Read all lines for parts processing
            lines = f.readlines()
            num_lines = len(lines)

            result = []
            for start_line, end_line in parts:
                if (
                    not 1 <= start_line <= num_lines
                    or not 1 <= end_line <= num_lines
                    or start_line > end_line
                ):
                    raise ValueError(
                        ERR_INVALID_RANGE.format(start_line, end_line, num_lines)
                    )
                result.extend(lines[start_line - 1 : end_line])
            return "".join(result).rstrip("\n")

        else:  # Entire file
            return f.read().rstrip("\n")


def expand_directory(directory, extensions=TXT_EXTENSIONS):
    """Find all files in a directory with specified extensions.

    This function expands a directory path into a list of file paths.

    Args:
        directory (str): The directory path to search.
        extensions (list): list of file extensions to include.

    Returns:
        list: A sorted list of file paths matching the extensions
        (not validated).
    """
    logger.debug(
        f"Expanding directory with directory='{directory}', extensions='{extensions}'"
    )
    matches = []
    for root, _, filenames in os.walk(directory):
        for filename in filenames:
            if any(filename.endswith(ext) for ext in extensions):
                matches.append(os.path.join(root, filename))
    return sorted(matches)


def is_file_path_line(line):
    """Determine if a line contains only a file path that exists.

    Args:
        line (str): The line to check.

    Returns:
        bool: True if the line is a valid file path, False otherwise.
    """
    if not line:  # Handle empty strings explicitly
        return False

    stripped_line = line.strip()
    if not stripped_line:  # Handle whitespace-only strings
        return False

    return not stripped_line.startswith("#") and os.path.isfile(stripped_line)


def is_mixed_content_bundle(lines):
    """Determine if a bundle contains mixed content (text and file paths).

    Args:
        lines (list): list of lines from the bundle file.

    Returns:
        bool: True if the bundle contains mixed content, False otherwise.
    """
    has_file_path = False
    has_non_path_line = False

    for line in lines:
        stripped_line = line.strip()
        if not stripped_line or stripped_line.startswith("#"):
            continue  # Skip empty lines and comments

        if os.path.isfile(stripped_line):
            has_file_path = True
        else:
            has_non_path_line = True

        # If we found both types, we can return early
        if has_file_path and has_non_path_line:
            return True

    # If we have only file paths, it's a traditional bundle
    return False


def process_mixed_content_bundle(lines):
    """Process a mixed content bundle by replacing file paths with their content.

    Args:
        lines (list): list of lines from the bundle file.

    Returns:
        str: The processed content with file paths replaced by their content.
    """
    result = []
    for line in lines:
        stripped_line = line.strip()
        if is_file_path_line(stripped_line):
            # This is a file path - substitute with file content
            try:
                file_content = get_file_content(stripped_line)
                result.append(file_content)
            except Exception as e:
                logger.warning(f"Error reading file {stripped_line}: {e}")
                # Keep the original line if file can't be read
                result.append(line)
        else:
            # Check for inline file references @[file path]
            inline_pattern = r"@\[(.*?)\]"
            matches = re.findall(inline_pattern, line)

            if matches:
                processed_line = line
                for file_path in matches:
                    if os.path.isfile(file_path):
                        try:
                            # Get file content and remove line breaks
                            file_content = get_file_content(file_path)
                            inline_content = file_content.replace("\n", " ")
                            inline_content = inline_content.strip()
                            # Replace the @[file path] with the inline content
                            processed_line = processed_line.replace(
                                f"@[{file_path}]", inline_content
                            )
                        except Exception as e:
                            logger.warning(ERR_INLINE_FILE_READ.format(file_path, e))
                            # Keep the original reference if file can't be read
                result.append(processed_line)
            else:
                # Regular text line - keep as is
                result.append(line)

    # Join all lines with newlines
    joined_result = "\n".join(result)

    # Process any inline file references that might span multiple lines
    inline_pattern = r"@\[(.*?)\]"
    matches = re.findall(inline_pattern, joined_result)

    if matches:
        processed_result = joined_result
        for file_path in matches:
            if os.path.isfile(file_path):
                try:
                    # Get file content and remove line breaks
                    file_content = get_file_content(file_path)
                    inline_content = file_content.replace("\n", " ").strip()
                    # Replace the @[file path] with the inline content
                    processed_result = processed_result.replace(
                        f"@[{file_path}]", inline_content
                    )
                except Exception as e:
                    logger.warning(ERR_INLINE_FILE_READ.format(file_path, e))
                    # Keep the original reference if file can't be read
        return processed_result

    return joined_result


def process_traditional_bundle(lines):
    """Process a traditional bundle (list of file paths).

    Args:
        lines (list): list of lines from the bundle file.

    Returns:
        list: A list of file paths.
    """
    return [
        line.strip()
        for line in lines
        if line.strip() and not line.strip().startswith("#")
    ]


def expand_bundles(bundle_file):
    """Extract list of files from a bundle file or process mixed content.

    This function handles two types of bundle files:
    1. Traditional bundles: Each line is a file path to include
    2. Mixed content bundles: Text mixed with file paths, where paths are
       replaced with file content

    Args:
        bundle_file (str): Path to the bundle file.

    Returns:
        list or str:
            - For traditional bundles: A list of file paths
            - For mixed content bundles: A string with file paths replaced by
              their content

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
                lines = content.splitlines()
            except ValueError as e:
                raise FileNotFoundError(
                    f"Invalid line reference in bundle file: {str(e)}"
                ) from e
        else:
            # Read the entire file
            content = get_file_content(file_path)
            lines = content.splitlines()
    except FileNotFoundError:
        raise FileNotFoundError(f"Bundle file not found: {file_path}") from None

    # Check if this is a mixed content bundle
    if is_mixed_content_bundle(lines):
        # Process as mixed content bundle
        return process_mixed_content_bundle(lines)
    else:
        # Process as traditional bundle (list of files)
        return process_traditional_bundle(lines)


def is_bundle_file(file_path):
    """Determine if a file is a bundle file by checking its contents.

    A file is considered a bundle if:
    1. Its first non-empty, non-comment line points to an existing file, or
    2. It contains at least one line that is a valid file path

    Args:
        file_path (str): The path to the file to check.

    Returns:
        bool: True if the file appears to be a bundle file, False otherwise.
    """
    logger.debug(f"Checking if {file_path} is a bundle file")
    try:
        with open(file_path) as f:
            lines = f.readlines()

            # First check: traditional bundle detection
            for line in lines:
                line = line.strip()
                if not line or line.startswith("#"):
                    continue
                if os.path.isfile(line):
                    return True
                else:
                    # If we find a non-file line, check for mixed content
                    break

            # Second check: mixed content bundle detection
            has_file_path = False
            for line in lines:
                line = line.strip()
                if not line or line.startswith("#"):
                    continue
                if os.path.isfile(line):
                    has_file_path = True
                    break

            return has_file_path
    except FileNotFoundError:
        return False
    except Exception as e:
        logger.error(f"Error checking bundle file: {e}")
        return False


def expand_single_arg(arg, extensions=None):
    """Helper function to expand a single argument.

    Args:
        arg (str): The argument to expand.
        extensions (list, optional): list of file extensions to include when
                                    expanding directories.

    Returns:
        list: A list of expanded file paths.
    """
    logger.debug(f"Expanding argument: {arg}")

    # Extract file path if there's a line reference
    file_path = arg
    if ":L" in arg:
        file_path = arg.split(":L", 1)[0]

    if os.path.isdir(arg):  # Directory path
        return expand_directory(arg, extensions=extensions)
    elif is_bundle_file(file_path):  # Bundle file
        bundle_result = expand_bundles(arg)
        if isinstance(bundle_result, str):
            # This is a mixed content bundle - create a temporary file
            import tempfile

            with tempfile.NamedTemporaryFile(
                mode="w", delete=False, suffix=".txt"
            ) as temp:
                temp.write(bundle_result)
                temp_name = temp.name
            return [temp_name]
        else:
            # Traditional bundle - list of files
            return bundle_result
    else:
        return [arg]  # Regular file path


def expand_args(args, extensions=None):
    """Expand a list of arguments into a flattened list of file paths.

    This function expands a list of arguments (file paths, directory paths, or
    bundle files) into a flattened list of file paths by calling the appropriate
    expander for each argument.

    Args:
        args (list): A list of file paths, directory paths, or bundle files.
        extensions (list, optional): list of file extensions to include when
            expanding directories. Defaults to TXT_EXTENSIONS.

    Returns:
        list: A flattened list of file paths (not validated).
    """
    logger.debug(f"Expanding arguments: {args}")

    # Use default extensions if none provided
    if extensions is None:
        extensions = TXT_EXTENSIONS

    # Use list comprehension with sum to flatten the list of lists
    return sum([expand_single_arg(arg, extensions) for arg in args], [])


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
            Line parts is a list of (start, end) tuples or None if no line
            reference.

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
            raise ValueError(ERR_INVALID_LINE_REF_FORMAT.format(line_spec))

        line_ref = line_spec

    # Verify the file path
    if not os.path.exists(file_path):
        raise FileNotFoundError(f"Error: Path does not exist: {file_path}")
    if not os.access(file_path, os.R_OK):
        raise PermissionError(f"Error: File is not readable: {file_path}")
    if os.path.isdir(file_path):
        raise IsADirectoryError(f"Error: Path is a directory, not a file: {file_path}")

    # Validate line reference if present
    if line_ref:
        try:
            ranges = parse_line_reference(line_ref)
            line_parts = convert_line_ranges_to_tuples(ranges)
            # Validate that all referenced lines exist in the file
            get_file_content(file_path, parts=line_parts)
        except ValueError as e:
            raise ValueError(f"Invalid line reference in {path}: {str(e)}") from e

    return file_path, line_parts


def file_sort_key(path):
    """Key function for sorting files by name then extension priority."""
    if isinstance(path, ContentItem):
        path = path.file_path
    base_name = os.path.splitext(os.path.basename(path))[0]
    ext = os.path.splitext(path)[1]
    # This ensures test_file.txt comes before test_file.md
    ext_priority = (
        0 if ext == TXT_EXTENSIONS[0] else 1 if ext == TXT_EXTENSIONS[1] else 2
    )
    return (base_name, ext_priority)


def get_files_from_args(srcs, extensions=None):
    """Process the sources and return a list of ContentItems.

    Args:
        srcs (list): list of source file paths, directories, or bundle files.
        extensions (list, optional): list of file extensions to include when
            expanding directories. Defaults to TXT_EXTENSIONS.

    Returns:
        list: A list of ContentItem objects.

    Raises:
        FileNotFoundError: If a file path does not exist.
        PermissionError: If a file is not readable.
        IsADirectoryError: If a path is a directory, not a file.
        ValueError: If a line reference is invalid or out of range.
    """
    # Use default extensions if none provided
    if extensions is None:
        extensions = TXT_EXTENSIONS

    # Phase 1: Expand all arguments into a flat list of file paths
    expanded_files = expand_args(srcs, extensions=extensions)
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
