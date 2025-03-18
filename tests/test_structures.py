"""Tests for core data structures in Nanodoc v2."""

from nanodoc.structures import Document, FileContent


def test_range_type():
    """Test that Range is a tuple of (int, Optional[int])."""
    # This is just testing our type alias, so we create a valid Range
    # and verify it behaves as expected
    range_with_end = (1, 10)
    range_to_eof = (1, None)

    # Check that we can access elements as expected
    assert range_with_end[0] == 1
    assert range_with_end[1] == 10

    assert range_to_eof[0] == 1
    assert range_to_eof[1] is None


def test_file_content_creation():
    """Test creating a FileContent object."""
    filepath = "/path/to/file.txt"
    ranges = [(1, 10), (20, 30)]

    # Create with required parameters
    fc = FileContent(filepath=filepath, ranges=ranges)
    assert fc.filepath == filepath
    assert fc.ranges == ranges
    assert fc.content == ""
    assert fc.is_bundle is False
    assert fc.original_source is None

    # Create with all parameters
    fc = FileContent(
        filepath=filepath,
        ranges=ranges,
        content="Some content",
        is_bundle=True,
        original_source="/path/to/source.txt",
    )
    assert fc.filepath == filepath
    assert fc.ranges == ranges
    assert fc.content == "Some content"
    assert fc.is_bundle is True
    assert fc.original_source == "/path/to/source.txt"

    # Test that we can modify fields
    fc.content = "Updated content"
    assert fc.content == "Updated content"


def test_document_creation():
    """Test creating a Document object."""
    file_content1 = FileContent(
        filepath="/path/to/file1.txt", ranges=[(1, 10)], content="Content 1"
    )
    file_content2 = FileContent(
        filepath="/path/to/file2.txt", ranges=[(1, None)], content="Content 2"
    )

    # Create document with content items
    doc = Document(content_items=[file_content1, file_content2])
    assert len(doc.content_items) == 2
    assert doc.content_items[0].content == "Content 1"
    assert doc.content_items[1].content == "Content 2"
    assert doc.toc is None

    # Test with TOC
    toc = [
        ("/path/to/file1.txt", "Heading 1", 1),
        ("/path/to/file2.txt", "Heading 2", 1),
    ]
    doc = Document(content_items=[file_content1, file_content2], toc=toc)
    assert doc.toc == toc
