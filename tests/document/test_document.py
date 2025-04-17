"""Tests for the document construction stage of Nanodoc v2."""

import pytest

from nanodoc.document import (
    CircularDependencyError,
    build_document,
    process_include_directive,
    process_inline_directive,
)
from nanodoc.structures import Document, FileContent


def test_build_document_with_regular_files():
    """Test building a document from regular (non-bundle) files."""
    # Create some file content objects
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

    # Build document
    document = build_document([file1, file2])

    # Check results
    assert len(document.content_items) == 2
    assert document.content_items[0].filepath == "/path/to/file1.txt"
    assert document.content_items[0].content == "Content of file 1"
    assert document.content_items[1].filepath == "/path/to/file2.txt"
    assert document.content_items[1].content == "Content of file 2"


def test_build_document_with_bundle_file(tmp_path):
    """Test building a document from a bundle file with inline directives."""
    # Create temporary files for testing
    bundle_dir = tmp_path / "bundle_test"
    bundle_dir.mkdir()

    # Create a file to be inlined
    inlined_file = bundle_dir / "inlined.txt"
    inlined_file.write_text("This is inlined content\n")

    # Create a bundle file
    bundle_file = bundle_dir / "bundle.ndoc"
    bundle_file.write_text("Bundle header\n@inline inlined.txt\nBundle footer\n")

    # Create FileContent objects
    bundle_content = FileContent(
        filepath=str(bundle_file),
        ranges=[(1, None)],
        content=bundle_file.read_text(),
        is_bundle=True,
    )

    # Build document
    document = build_document([bundle_content])

    # Check results - should have 3 content items:
    # 1. Bundle header
    # 2. Inlined content
    # 3. Bundle footer
    assert len(document.content_items) == 3
    assert document.content_items[0].content == "Bundle header\n"
    assert document.content_items[1].content == "This is inlined content\n"
    assert document.content_items[2].content == "Bundle footer\n"


def test_inline_directive(tmp_path):
    """Test processing an inline directive."""
    # Create temporary files for testing
    test_dir = tmp_path / "inline_test"
    test_dir.mkdir()

    # Create a file to be inlined
    inlined_file = test_dir / "inlined.txt"
    inlined_file.write_text("Inlined content\n")

    # Create a Document to update
    document = Document(content_items=[])

    # Process inline directive
    parent_path = str(test_dir / "parent.ndoc")
    process_inline_directive(
        inline_path=str(inlined_file.name),
        base_path=str(test_dir),
        document=document,
        processed_files=set(),
        parent_bundle=parent_path,
    )

    # Check results
    assert len(document.content_items) == 1
    assert document.content_items[0].content == "Inlined content\n"
    assert document.content_items[0].original_source == parent_path


def test_include_directive(tmp_path):
    """Test processing an include directive."""
    # Create temporary files for testing
    test_dir = tmp_path / "include_test"
    test_dir.mkdir()

    # Create a file to be included
    included_file = test_dir / "included.txt"
    included_file.write_text("Included content\n")

    # Create a Document to update
    document = Document(content_items=[])

    # Process include directive
    process_include_directive(
        include_path=str(included_file.name),
        base_path=str(test_dir),
        document=document,
        processed_files=set(),
        parent_bundle=str(test_dir / "parent.ndoc"),
    )

    # Check results
    assert len(document.content_items) == 1
    assert document.content_items[0].content == "Included content\n"
    # Included content retains its own identity (not marked with original_source)
    assert document.content_items[0].original_source is None


def test_circular_dependency_detection(tmp_path):
    """Test that circular dependencies are detected and reported."""
    # Create a file that includes itself (directly)
    test_dir = tmp_path / "circular_test"
    test_dir.mkdir()

    circular_file = test_dir / "circular.ndoc"
    circular_file.write_text("Header\n@include circular.ndoc\nFooter\n")

    # Create FileContent
    file_content = FileContent(
        filepath=str(circular_file),
        ranges=[(1, None)],
        content=circular_file.read_text(),
        is_bundle=True,
    )

    # Attempt to build document - should raise CircularDependencyError
    with pytest.raises(CircularDependencyError):
        build_document([file_content])


def test_nested_bundles(tmp_path):
    """Test processing nested bundles."""
    # Create test files
    test_dir = tmp_path / "nested_test"
    test_dir.mkdir()

    # Create a file to be included in nested bundle
    content_file = test_dir / "content.txt"
    content_file.write_text("Content text\n")

    # Create a nested bundle
    nested_bundle = test_dir / "nested.ndoc"
    nested_bundle.write_text(
        "Nested bundle header\n@inline content.txt\nNested bundle footer\n"
    )

    # Create a main bundle that includes the nested bundle
    main_bundle = test_dir / "main.ndoc"
    main_bundle.write_text(
        "Main bundle header\n@include nested.ndoc\nMain bundle footer\n"
    )

    # Create FileContent for main bundle
    main_content = FileContent(
        filepath=str(main_bundle),
        ranges=[(1, None)],
        content=main_bundle.read_text(),
        is_bundle=True,
    )

    # Build document
    document = build_document([main_content])

    # Check results
    assert len(document.content_items) == 5
    assert document.content_items[0].content == "Main bundle header\n"
    # Nested bundle - header, content, footer
    assert "Nested bundle header" in document.content_items[1].content
    assert "Content text" in document.content_items[2].content
    assert "Nested bundle footer" in document.content_items[3].content
    assert document.content_items[4].content == "Main bundle footer\n"


def test_missing_inline_file(tmp_path):
    """Test handling missing files in inline directives."""
    # Create test directory
    test_dir = tmp_path / "missing_test"
    test_dir.mkdir()

    # Create a bundle that references a non-existent file
    bundle_file = test_dir / "bundle.ndoc"
    bundle_file.write_text("Bundle header\n@inline non_existent.txt\nBundle footer\n")

    # Create FileContent
    bundle_content = FileContent(
        filepath=str(bundle_file),
        ranges=[(1, None)],
        content=bundle_file.read_text(),
        is_bundle=True,
    )

    # Build document
    document = build_document([bundle_content])

    # Check results - should contain error message in place of missing file
    assert len(document.content_items) == 3
    assert document.content_items[0].content == "Bundle header\n"
    assert "ERROR: Could not find inlined file" in document.content_items[1].content
    assert document.content_items[2].content == "Bundle footer\n"
