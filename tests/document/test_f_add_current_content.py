"""Tests for document._add_current_content function."""

from nanodoc.document import _add_current_content
from nanodoc.structures import Document, FileContent


def test_add_current_content_basic():
    """Test adding current content to a document."""
    # Create test data
    file_content = FileContent(
        filepath="/path/to/file.txt",
        ranges=[],
        content="Original content",
        is_bundle=False,
    )
    current_content = ["Line 1", "Line 2", "Line 3"]
    document = Document(content_items=[])

    # Add current content to document
    _add_current_content(current_content, file_content, document)

    # Check result
    assert len(document.content_items) == 1
    added_content = document.content_items[0]
    assert added_content.filepath == "/path/to/file.txt"
    assert added_content.content == "Line 1\nLine 2\nLine 3\n"
    assert not added_content.is_bundle
    assert added_content.original_source == "/path/to/file.txt"


def test_add_current_content_empty():
    """Test adding empty content to a document."""
    # Create test data
    file_content = FileContent(
        filepath="/path/to/file.txt",
        ranges=[],
        content="Original content",
        is_bundle=False,
    )
    current_content = []
    document = Document(content_items=[])

    # Add current content to document
    _add_current_content(current_content, file_content, document)

    # Check result
    assert len(document.content_items) == 1
    added_content = document.content_items[0]
    assert added_content.filepath == "/path/to/file.txt"
    assert added_content.content == "\n"  # Empty content with newline
    assert not added_content.is_bundle


def test_add_current_content_multiple_calls():
    """Test adding content multiple times to a document."""
    # Create test data
    file_content = FileContent(
        filepath="/path/to/file.txt",
        ranges=[],
        content="Original content",
        is_bundle=False,
    )
    document = Document(content_items=[])

    # Add content twice
    _add_current_content(["First line"], file_content, document)
    _add_current_content(["Second line"], file_content, document)

    # Check result
    assert len(document.content_items) == 2
    assert document.content_items[0].content == "First line\n"
    assert document.content_items[1].content == "Second line\n"
