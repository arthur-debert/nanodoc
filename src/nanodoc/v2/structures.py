"""Core data structures for Nanodoc v2.

This module defines the core data structures used throughout the Nanodoc v2
implementation, including Range and FileContent.
"""

from dataclasses import dataclass, field
from typing import Optional

# A Range can be represented as a simple tuple (start, end) where
# start is inclusive and end is inclusive or None (for EOF)
Range = tuple[int, Optional[int]]


@dataclass
class FileContent:
    """Represents the content and metadata for a single file.

    Attributes:
        filepath: Path to the file
        ranges: Line ranges to include (list of (start, end) tuples)
        content: Content after applying ranges
        is_bundle: True if this represents a bundle file
        original_source: Source file if part of an inline bundle
    """

    filepath: str
    ranges: list[Range]
    content: str = ""
    is_bundle: bool = False
    original_source: Optional[str] = None


@dataclass
class Document:
    """Represents the entire document after processing bundles.

    Attributes:
        content_items: Ordered list of content blocks
        toc: Table of contents data
        theme_name: Name of the theme to use for styling
        use_rich_formatting: Whether to use Rich formatting
    """

    content_items: list[FileContent]
    toc: list = None  # Will be defined in more detail later
    theme_name: Optional[str] = None
    use_rich_formatting: bool = False
    formatting_options: dict = field(default_factory=dict)
