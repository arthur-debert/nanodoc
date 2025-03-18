"""Tests for the rendering stage of Nanodoc v2."""

from nanodoc.v2.renderer import (
    _add_line_numbers,
    _extract_headings,
    generate_toc,
    render_document,
)
from nanodoc.v2.structures import Document, FileContent


def test_render_document_basic():
    """Test basic document rendering with regular files."""
    # Create test document
    file1 = FileContent(
        filepath="/path/to/file1.txt",
        ranges=[(1, None)],
        content="Content of file 1",
        is_bundle=False,
    )
    file2 = FileContent(
        filepath="/path/to/file2.txt",
        ranges=[(1, None)],
        content="Content of file 2",
        is_bundle=False,
    )
    document = Document(content_items=[file1, file2])

    # Render document
    result = render_document(document)

    # Check result
    assert "# file1.txt" in result
    assert "Content of file 1" in result
    assert "# file2.txt" in result
    assert "Content of file 2" in result


def test_render_document_with_inline_content():
    """Test document rendering with inlined content."""
    # Create parent bundle file
    parent = FileContent(
        filepath="/path/to/bundle.txt",
        ranges=[(1, None)],
        content="Bundle content",
        is_bundle=True,
    )
    # Create inlined content
    inlined = FileContent(
        filepath="/path/to/inlined.txt",
        ranges=[(1, None)],
        content="Inlined content",
        is_bundle=False,
        original_source="/path/to/bundle.txt",  # Indicates this was inlined
    )
    # Create regular file
    regular = FileContent(
        filepath="/path/to/regular.txt",
        ranges=[(1, None)],
        content="Regular content",
        is_bundle=False,
    )
    document = Document(content_items=[parent, inlined, regular])

    # Render document
    result = render_document(document)

    # Check result
    assert "# bundle.txt" in result
    assert "Bundle content" in result
    # Inlined content should not have its own header
    assert "# inlined.txt" not in result
    assert "Inlined content" in result
    assert "# regular.txt" in result
    assert "Regular content" in result


def test_render_document_with_line_numbers():
    """Test document rendering with line numbers."""
    # Create test document
    file1 = FileContent(
        filepath="/path/to/file1.txt",
        ranges=[(1, None)],
        content="Line 1\nLine 2\nLine 3",
        is_bundle=False,
    )
    document = Document(content_items=[file1])

    # Render document with line numbers
    result = render_document(document, include_line_numbers=True)

    # Check result
    assert "# file1.txt" in result
    assert "   1 | Line 1" in result
    assert "   2 | Line 2" in result
    assert "   3 | Line 3" in result


def test_render_document_with_toc():
    """Test document rendering with table of contents."""
    # Create test document with headings
    file1 = FileContent(
        filepath="/path/to/file1.md",
        ranges=[(1, None)],
        content="# Heading 1\nContent\n## Subheading\nMore content",
        is_bundle=False,
    )
    file2 = FileContent(
        filepath="/path/to/file2.md",
        ranges=[(1, None)],
        content="# Heading 2\nOther content",
        is_bundle=False,
    )
    document = Document(content_items=[file1, file2])

    # Render document with TOC
    result = render_document(document, include_toc=True)

    # Check result
    assert "# Table of Contents" in result
    assert "- file1.md" in result
    assert "  - Heading 1" in result
    assert "  - Subheading" in result
    assert "- file2.md" in result
    assert "  - Heading 2" in result


def test_generate_toc():
    """Test generating table of contents."""
    # Create test document with headings
    file1 = FileContent(
        filepath="/path/to/file1.md",
        ranges=[(1, None)],
        content="# Heading 1\nContent\n## Subheading\nMore content",
        is_bundle=False,
    )
    file2 = FileContent(
        filepath="/path/to/file2.md",
        ranges=[(1, None)],
        content="# Heading 2\nOther content",
        is_bundle=False,
    )
    document = Document(content_items=[file1, file2])

    # Generate TOC
    toc = generate_toc(document)

    # Check result
    assert "# Table of Contents" in toc
    assert "- file1.md" in toc
    assert "  - Heading 1" in toc
    assert "  - Subheading" in toc
    assert "- file2.md" in toc
    assert "  - Heading 2" in toc

    # Check that document.toc is populated
    assert document.toc is not None
    assert "/path/to/file1.md" in document.toc
    assert len(document.toc["/path/to/file1.md"]) == 2  # Two headings


def test_generate_toc_empty():
    """Test generating TOC with no headings."""
    # Create test document with no headings
    file1 = FileContent(
        filepath="/path/to/file1.txt",
        ranges=[(1, None)],
        content="No headings here",
        is_bundle=False,
    )
    document = Document(content_items=[file1])

    # Generate TOC
    toc = generate_toc(document)

    # Check result
    assert toc == ""


def test_extract_headings():
    """Test extracting headings from document content."""
    # Create test document with headings
    file1 = FileContent(
        filepath="/path/to/file1.md",
        ranges=[(1, None)],
        content="# Heading 1\nContent\n## Subheading\nMore content\n### H3\n",
        is_bundle=False,
    )
    document = Document(content_items=[file1])

    # Extract headings
    headings = _extract_headings(document)

    # Check result
    assert "/path/to/file1.md" in headings
    assert len(headings["/path/to/file1.md"]) == 2  # H1 and H2, but not H3
    assert headings["/path/to/file1.md"][0] == ("Heading 1", 1)
    assert headings["/path/to/file1.md"][1] == ("Subheading", 3)


def test_extract_headings_with_inline():
    """Test extracting headings from document with inlined content."""
    # Create parent bundle file
    parent = FileContent(
        filepath="/path/to/bundle.md",
        ranges=[(1, None)],
        content="# Bundle Heading",
        is_bundle=True,
    )
    # Create inlined content
    inlined = FileContent(
        filepath="/path/to/inlined.md",
        ranges=[(1, None)],
        content="# Inlined Heading",
        is_bundle=False,
        original_source="/path/to/bundle.md",  # Indicates this was inlined
    )
    document = Document(content_items=[parent, inlined])

    # Extract headings
    headings = _extract_headings(document)

    # Check result
    assert "/path/to/bundle.md" in headings
    assert len(headings["/path/to/bundle.md"]) == 2
    assert headings["/path/to/bundle.md"][0] == ("Bundle Heading", 1)
    assert headings["/path/to/bundle.md"][1] == ("Inlined Heading", 1)


def test_add_line_numbers():
    """Test adding line numbers to content."""
    content = "Line 1\nLine 2\nLine 3"
    result = _add_line_numbers(content)
    assert result == "   1 | Line 1\n   2 | Line 2\n   3 | Line 3"


def test_add_line_numbers_empty():
    """Test adding line numbers to empty content."""
    content = ""
    result = _add_line_numbers(content)
    assert result == "   1 | "
