"""Path resolution for Nanodoc v2.

This module handles the "Resolving Paths" stage of the Nanodoc v2 pipeline.
It takes user input (globs, paths, directories) and returns a list of absolute
file paths.
"""

import fnmatch
import glob
import os


def resolve_paths(
    inputs: list[str], recursive: bool = False, include_hidden: bool = False
) -> list[str]:
    """Resolve a list of input paths to absolute file paths.

    This function:
    - Expands globs (e.g., *.txt)
    - Handles directory inputs (recursively find files if recursive=True)
    - Validates that paths exist
    - Converts relative paths to absolute paths

    Args:
        inputs: List of input paths (globs, directories, files)
        recursive: Whether to recursively search directories
        include_hidden: Whether to include hidden files and directories

    Returns:
        List of absolute file paths

    Raises:
        FileNotFoundError: If any of the non-glob inputs do not exist
    """
    resolved_paths = []

    for input_path in inputs:
        # Check for glob patterns
        if any(char in input_path for char in ["*", "?", "["]):
            # Handle glob patterns
            matched_paths = glob.glob(input_path, recursive=recursive)

            # For recursive searches including hidden files, glob may not find
            # hidden files even with recursive=True, so handle manually
            if recursive and include_hidden and "**" in input_path:
                # Extract the base directory from the pattern
                parts = input_path.split("**", 1)
                base_dir = parts[0].rstrip(os.sep)

                # Extract the pattern after **
                pattern = parts[1].lstrip(os.sep) if len(parts) > 1 else "*"

                # If base dir exists, walk it manually to find all files
                # including hidden ones that match the pattern
                if os.path.isdir(base_dir):
                    # Get all files from directory (including hidden ones)
                    # and filter by pattern
                    all_files = []
                    for root, _dirs, files in os.walk(base_dir):
                        # Don't skip hidden directories for include_hidden=True
                        for file in files:
                            path = os.path.join(root, file)
                            rel_path = os.path.relpath(path, base_dir)
                            if fnmatch.fnmatch(rel_path, pattern):
                                all_files.append(os.path.abspath(path))

                    # Add the manually found files to the matched_paths
                    # but avoid duplicates
                    abs_paths = {os.path.abspath(p) for p in matched_paths}
                    for path in all_files:
                        if path not in abs_paths:
                            matched_paths.append(path)

            if not matched_paths:
                continue  # No matches found, skip to next input

            for path in matched_paths:
                # Skip hidden files/directories if not explicitly included
                basename = os.path.basename(path)
                if not include_hidden and basename.startswith("."):
                    continue

                if os.path.isfile(path):
                    resolved_paths.append(os.path.abspath(path))
                elif os.path.isdir(path) and recursive:
                    # If it's a directory and recursive is True, get all files
                    directory_paths = get_files_from_directory(
                        path, recursive=True, include_hidden=include_hidden
                    )
                    resolved_paths.extend(directory_paths)
        elif os.path.isdir(input_path):
            # Handle directory inputs
            directory_paths = get_files_from_directory(
                input_path, recursive=recursive, include_hidden=include_hidden
            )
            resolved_paths.extend(directory_paths)
        elif os.path.isfile(input_path):
            # Handle file inputs
            resolved_paths.append(os.path.abspath(input_path))
        else:
            # Handle non-existent paths
            raise FileNotFoundError(f"Path does not exist: {input_path}")

    # Remove duplicates while preserving order
    seen = set()
    return [p for p in resolved_paths if not (p in seen or seen.add(p))]


def get_files_from_directory(
    directory: str, recursive: bool = False, include_hidden: bool = False
) -> list[str]:
    """Get all files from a directory.

    Args:
        directory: Directory path
        recursive: Whether to recursively search subdirectories
        include_hidden: Whether to include hidden files and directories

    Returns:
        List of absolute file paths
    """
    result = []

    if recursive:
        for root, dirs, files in os.walk(directory):
            # Skip hidden directories if not explicitly included
            if not include_hidden:
                dirs[:] = [d for d in dirs if not d.startswith(".")]

            for file in files:
                # Skip hidden files if not explicitly included
                if not include_hidden and file.startswith("."):
                    continue

                result.append(os.path.abspath(os.path.join(root, file)))
    else:
        # Non-recursive directory listing
        for item in os.listdir(directory):
            # Skip hidden items if not explicitly included
            if not include_hidden and item.startswith("."):
                continue

            path = os.path.join(directory, item)
            if os.path.isfile(path):
                result.append(os.path.abspath(path))

    return result
