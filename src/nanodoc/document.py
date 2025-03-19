"""Document construction for Nanodoc v2.

This module handles the "Building the Content" stage of the Nanodoc v2 pipeline.
It processes bundle files (recursively), parses inline and include directives,
creates new FileContent objects for inlined content, flattens the structure into
a Document object, and detects circular dependencies.
"""

import logging
import os
import re
from typing import Optional

from nanodoc.extractor import gather_content, resolve_files
from nanodoc.structures import Document, FileContent

# Initialize logger
logger = logging.getLogger("document")


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
    logger.debug("Starting document build with %d file contents", len(file_contents))

    # Initialize document with empty content items
    document = Document(content_items=[])

    # Track processed files to detect circular dependencies
    processed_files: set[str] = set()

    # Process each file content
    for file_content in file_contents:
        logger.debug("Processing file: %s", file_content.filepath)
        process_content(
            file_content=file_content,
            document=document,
            processed_files=processed_files,
        )

    logger.debug(
        "Document build complete with %d content items", len(document.content_items)
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
    logger.debug(
        "Processing content from %s (parent: %s)", filepath, parent_bundle or "None"
    )

    # If this file is already being processed in the current branch, we have a cycle
    if filepath in processed_files:
        parent_info = f" from {parent_bundle}" if parent_bundle else ""
        msg = f"Circular dependency detected: {filepath}"
        msg += f" included{parent_info}"
        logger.error(msg)
        raise CircularDependencyError(msg)

    # If it's not a bundle file, add it directly to the document
    if not file_content.is_bundle:
        logger.debug("Adding non-bundle file directly: %s", filepath)
        document.content_items.append(file_content)
        return

    # Mark this file as being processed
    processed_files.add(filepath)
    logger.debug("Processing bundle file: %s", filepath)

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
        logger.debug("Completed processing bundle: %s", filepath)


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
    logger.debug("Processing directives in bundle: %s", file_content.filepath)

    # Parse lines looking for directives
    lines = file_content.content.splitlines()
    current_content = []

    for line_num, line in enumerate(lines, 1):
        # Check for inline directive: @inline <file_path>[:<range>]
        inline_match = re.match(r"@inline\s+(.+)", line.strip())
        if inline_match:
            logger.debug(
                "Found @inline directive at line %d: %s",
                line_num,
                inline_match.group(1),
            )
            if current_content:
                _add_current_content(current_content, file_content, document)
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
            logger.debug(
                "Found @include directive at line %d: %s",
                line_num,
                include_match.group(1),
            )
            if current_content:
                _add_current_content(current_content, file_content, document)
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
        _add_current_content(current_content, file_content, document)


def _add_current_content(
    current_content: list[str],
    file_content: FileContent,
    document: Document,
) -> None:
    """Helper to add accumulated content to document."""
    inline_content = FileContent(
        filepath=file_content.filepath,
        ranges=[],  # Not applicable for inline content
        content="\n".join(current_content) + "\n",
        is_bundle=False,
        original_source=file_content.filepath,
    )
    document.content_items.append(inline_content)
    logger.debug(
        "Added %d lines of content from %s", len(current_content), file_content.filepath
    )


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
    logger.debug("Processing @inline directive: %s from %s", inline_path, parent_bundle)

    # Resolve the path (handle relative paths)
    if not os.path.isabs(inline_path):
        resolved_path = os.path.join(base_path, inline_path)
    else:
        resolved_path = inline_path
    logger.debug("Resolved inline path: %s", resolved_path)

    # Create a FileContent object for the inlined file
    inlined_files = resolve_files([resolved_path])

    if not inlined_files:
        # Handle case where the file doesn't exist
        logger.error("Could not find inlined file: %s", inline_path)
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
        logger.debug("Gathered %d content items from inline", len(inlined_with_content))

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
        logger.error("Could not find or read inlined file: %s", inline_path)
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
    logger.debug(
        "Processing @include directive: %s from %s", include_path, parent_bundle
    )

    # Resolve the path (handle relative paths)
    if not os.path.isabs(include_path):
        resolved_path = os.path.join(base_path, include_path)
    else:
        resolved_path = include_path
    logger.debug("Resolved include path: %s", resolved_path)

    # Create a FileContent object for the included file
    included_files = resolve_files([resolved_path])

    if not included_files:
        # Handle case where the file doesn't exist
        logger.error("Could not find included file: %s", include_path)
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
        logger.debug(
            "Gathered %d content items from include", len(included_with_content)
        )

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
        logger.error("Could not find or read included file: %s", include_path)
        error_msg = f"ERROR: Could not find included file: {include_path}\n"
        error_content = FileContent(
            filepath=parent_bundle,
            ranges=[],
            content=error_msg,
            is_bundle=False,
            original_source=parent_bundle,
        )
        document.content_items.append(error_content)
