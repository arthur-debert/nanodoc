"""Document construction for Nanodoc v2.

This module handles the "Building the Content" stage of the Nanodoc v2 pipeline.
It processes bundle files (recursively), parses inline and include directives,
creates new FileContent objects for inlined content, flattens the structure into
a Document object, and detects circular dependencies.
"""

import os
import re
from typing import Optional

from nanodoc.v2.extractor import gather_content, resolve_files
from nanodoc.v2.structures import Document, FileContent


class CircularDependencyError(Exception):
    """Raised when a circular dependency is detected in bundle processing."""

    pass


def build_document(file_contents: list[FileContent]) -> Document:
    """Build a document from a list of FileContent objects.

    This function:
    - Processes bundle files recursively
    - Parses inline and include directives
    - Creates new FileContent objects for inlined content
    - Flattens the structure into a Document object
    - Detects circular dependencies

    Args:
        file_contents: List of FileContent objects with content loaded

    Returns:
        Document object with all content processed

    Raises:
        CircularDependencyError: If a circular dependency is detected
    """
    # Initialize document with empty content items
    document = Document(content_items=[])

    # Track processed files to detect circular dependencies
    processed_files: set[str] = set()

    # Process each file content
    for file_content in file_contents:
        # Process content (handling bundles recursively)
        process_content(
            file_content=file_content,
            document=document,
            processed_files=processed_files,
        )

    return document


def process_content(
    file_content: FileContent,
    document: Document,
    processed_files: set[str],
    parent_bundle: Optional[str] = None,
) -> None:
    """Process content from a file, handling bundle directives recursively.

    Args:
        file_content: FileContent object to process
        document: Document object to update
        processed_files: Set of already processed file paths
        parent_bundle: Path of the parent bundle (for circular dependency detection)

    Raises:
        CircularDependencyError: If a circular dependency is detected
    """
    filepath = file_content.filepath

    # If this file is already being processed in the current branch, we have a cycle
    if filepath in processed_files:
        parent_info = f" from {parent_bundle}" if parent_bundle else ""
        msg = f"Circular dependency detected: {filepath}"
        msg += f" included{parent_info}"
        raise CircularDependencyError(msg)

    # If it's not a bundle file, add it directly to the document
    if not file_content.is_bundle:
        document.content_items.append(file_content)
        return

    # Mark this file as being processed
    processed_files.add(filepath)

    try:
        # Parse and process directives
        process_bundle_directives(
            file_content=file_content,
            document=document,
            processed_files=processed_files,
        )
    finally:
        # Remove from processed set after we're done with this branch
        processed_files.remove(filepath)


def process_bundle_directives(
    file_content: FileContent,
    document: Document,
    processed_files: set[str],
) -> None:
    """Parse and process bundle directives.

    Args:
        file_content: FileContent object containing bundle file content
        document: Document object to update
        processed_files: Set of already processed file paths
    """
    # Parse lines looking for directives
    lines = file_content.content.splitlines()
    current_content = []

    for line in lines:
        # Check for inline directive: @inline <file_path>[:<range>]
        inline_match = re.match(r"@inline\s+(.+)", line.strip())
        if inline_match:
            # If we have accumulated content, add it first
            if current_content:
                inline_content = FileContent(
                    filepath=file_content.filepath,
                    ranges=[],  # Not applicable for inline content
                    content="\n".join(current_content) + "\n",
                    is_bundle=False,
                    original_source=file_content.filepath,
                )
                document.content_items.append(inline_content)
                current_content = []

            # Process the inline directive
            process_inline_directive(
                inline_path=inline_match.group(1),
                base_path=os.path.dirname(file_content.filepath),
                document=document,
                processed_files=processed_files,
                parent_bundle=file_content.filepath,
            )
            continue

        # Check for include directive: @include <file_path>[:<range>]
        include_match = re.match(r"@include\s+(.+)", line.strip())
        if include_match:
            # If we have accumulated content, add it first
            if current_content:
                inline_content = FileContent(
                    filepath=file_content.filepath,
                    ranges=[],  # Not applicable for inline content
                    content="\n".join(current_content) + "\n",
                    is_bundle=False,
                    original_source=file_content.filepath,
                )
                document.content_items.append(inline_content)
                current_content = []

            # Process the include directive
            process_include_directive(
                include_path=include_match.group(1),
                base_path=os.path.dirname(file_content.filepath),
                document=document,
                processed_files=processed_files,
                parent_bundle=file_content.filepath,
            )
            continue

        # Regular line, add to current content
        current_content.append(line)

    # Add any remaining content
    if current_content:
        inline_content = FileContent(
            filepath=file_content.filepath,
            ranges=[],  # Not applicable for inline content
            content="\n".join(current_content) + "\n",
            is_bundle=False,
            original_source=file_content.filepath,
        )
        document.content_items.append(inline_content)


def process_inline_directive(
    inline_path: str,
    base_path: str,
    document: Document,
    processed_files: set[str],
    parent_bundle: str,
) -> None:
    """Process an inline directive by inlining content from another file.

    The inlined content is added as FileContent objects with original_source
    set to indicate that they are inlined.

    Args:
        inline_path: Path to the file to inline (with optional range specifier)
        base_path: Base directory path for resolving relative paths
        document: Document object to update
        processed_files: Set of already processed file paths
        parent_bundle: Path of the parent bundle
    """
    # Resolve the path (handle relative paths)
    if not os.path.isabs(inline_path):
        resolved_path = os.path.join(base_path, inline_path)
    else:
        resolved_path = inline_path

    # Create a FileContent object for the inlined file
    inlined_files = resolve_files([resolved_path])

    if not inlined_files:
        # Handle case where the file doesn't exist
        error_msg = f"ERROR: Could not find inlined file: {inline_path}\n"
        error_content = FileContent(
            filepath=parent_bundle,
            ranges=[],
            content=error_msg,
            is_bundle=False,
            original_source=parent_bundle,
        )
        document.content_items.append(error_content)
        return

    try:
        # Get content for the inlined file
        inlined_with_content = gather_content(inlined_files)

        # Mark inlined content as being from the original source
        for content in inlined_with_content:
            content.original_source = parent_bundle

        # Process the inlined content (handle nested bundles)
        for content in inlined_with_content:
            process_content(
                file_content=content,
                document=document,
                processed_files=processed_files,
                parent_bundle=parent_bundle,
            )
    except FileNotFoundError:
        # Handle case where the file doesn't exist or can't be read
        error_msg = f"ERROR: Could not find inlined file: {inline_path}\n"
        error_content = FileContent(
            filepath=parent_bundle,
            ranges=[],
            content=error_msg,
            is_bundle=False,
            original_source=parent_bundle,
        )
        document.content_items.append(error_content)


def process_include_directive(
    include_path: str,
    base_path: str,
    document: Document,
    processed_files: set[str],
    parent_bundle: str,
) -> None:
    """Process an include directive by including content from another file.

    Unlike inline, included content is treated as a separate file.

    Args:
        include_path: Path to the file to include (with optional range specifier)
        base_path: Base directory path for resolving relative paths
        document: Document object to update
        processed_files: Set of already processed file paths
        parent_bundle: Path of the parent bundle
    """
    # Resolve the path (handle relative paths)
    if not os.path.isabs(include_path):
        resolved_path = os.path.join(base_path, include_path)
    else:
        resolved_path = include_path

    # Create a FileContent object for the included file
    included_files = resolve_files([resolved_path])

    if not included_files:
        # Handle case where the file doesn't exist
        error_msg = f"ERROR: Could not find included file: {include_path}\n"
        error_content = FileContent(
            filepath=parent_bundle,
            ranges=[],
            content=error_msg,
            is_bundle=False,
            original_source=parent_bundle,
        )
        document.content_items.append(error_content)
        return

    try:
        # Get content for the included file
        included_with_content = gather_content(included_files)

        # Process the included content (handle nested bundles)
        for content in included_with_content:
            process_content(
                file_content=content,
                document=document,
                processed_files=processed_files,
                parent_bundle=parent_bundle,
            )
    except FileNotFoundError:
        # Handle case where the file doesn't exist or can't be read
        error_msg = f"ERROR: Could not find included file: {include_path}\n"
        error_content = FileContent(
            filepath=parent_bundle,
            ranges=[],
            content=error_msg,
            is_bundle=False,
            original_source=parent_bundle,
        )
        document.content_items.append(error_content)
