"""Tests for formatter.apply_theme_to_document function."""

from nanodoc.formatter import apply_theme_to_document
from nanodoc.structures import Document, FileContent


def test_apply_theme_to_document():
    """Test applying theme to a document."""
    # Create a document
    file1 = FileContent(
        filepath="/path/to/file1.txt",
        ranges=[(1, None)],
        content="Content of file 1",
        is_bundle=False,
    )
    document = Document(content_items=[file1])

    # Apply theme
    themed_document = apply_theme_to_document(
        document, theme_name="test-theme", use_rich_formatting=True
    )

    # Check that theme info was added
    assert themed_document.theme_name == "test-theme"
    assert themed_document.use_rich_formatting is True


def test_apply_theme_to_document_without_formatting():
    """Test applying theme to a document with formatting disabled."""
    # Create a document
    file1 = FileContent(
        filepath="/path/to/file1.txt",
        ranges=[(1, None)],
        content="Content of file 1",
        is_bundle=False,
    )
    document = Document(content_items=[file1])

    # Apply theme with formatting disabled
    themed_document = apply_theme_to_document(
        document, theme_name="test-theme", use_rich_formatting=False
    )

    # Check that theme info was not added
    assert themed_document.theme_name is None
    assert themed_document.use_rich_formatting is False
