"""Tests for renderer._extract_headings function."""

from nanodoc.renderer import _extract_headings
from nanodoc.structures import Document, FileContent


def test_extract_headings_markdown():
    """Test extracting markdown headings from a document."""
    # Create a document with markdown headings
    file1 = FileContent(
        filepath="/path/to/file1.md",
        ranges=[],
        content="# Heading 1\nContent\n## Heading 2\nMore content",
        is_bundle=False,
    )
    document = Document(content_items=[file1])

    # Extract headings
    headings = _extract_headings(document)

    # Check results
    assert "/path/to/file1.md" in headings
    file1_headings = headings["/path/to/file1.md"]
    assert len(file1_headings) == 2
    assert file1_headings[0][0] == "Heading 1"
    assert file1_headings[1][0] == "Heading 2"


def test_extract_headings_plain_text():
    """Test extracting headings from plain text document."""
    # Create a document with plain text
    file1 = FileContent(
        filepath="/path/to/file1.txt",
        ranges=[],
        content="This is a plain text file\nWith multiple lines\n"
        "But no markdown headings",
        is_bundle=False,
    )
    document = Document(content_items=[file1])

    # Extract headings
    headings = _extract_headings(document)

    # Check results - should create a pseudo-heading from first line
    assert "/path/to/file1.txt" in headings
    file1_headings = headings["/path/to/file1.txt"]
    assert len(file1_headings) == 1
    assert file1_headings[0][0] == "This is a plain text file"


def test_extract_headings_empty_document():
    """Test extracting headings from an empty document."""
    # Create an empty document
    document = Document(content_items=[])

    # Extract headings
    headings = _extract_headings(document)

    # Check results
    assert len(headings) == 0


def test_extract_headings_with_original_source():
    """Test extracting headings with original_source attribute."""
    # Create a document with an inlined content item
    file1 = FileContent(
        filepath="/path/to/inline_content.md",
        ranges=[],
        content="# Heading from inline",
        is_bundle=False,
        original_source="/path/to/original.md",
    )
    document = Document(content_items=[file1])

    # Extract headings
    headings = _extract_headings(document)

    # Check results - should use original source as the key
    assert "/path/to/original.md" in headings
    assert "/path/to/inline_content.md" not in headings
