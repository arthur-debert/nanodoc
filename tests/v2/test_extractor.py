"""Tests for file extraction in Nanodoc v2."""

import pytest

from nanodoc.v2.extractor import (
    _apply_ranges,
    _parse_path_and_ranges,
    gather_content,
    resolve_files,
)
from nanodoc.v2.structures import FileContent


def test_parse_path_and_ranges_no_ranges():
    """Test parsing a path with no range specifiers."""
    path, ranges = _parse_path_and_ranges("file.txt")
    assert path == "file.txt"
    assert ranges == [(1, None)]


def test_parse_path_and_ranges_with_single_range():
    """Test parsing a path with a single range specifier."""
    path, ranges = _parse_path_and_ranges("file.txt:10-20")
    assert path == "file.txt"
    assert ranges == [(10, 20)]


def test_parse_path_and_ranges_with_eof_range():
    """Test parsing a path with a range to end of file."""
    path, ranges = _parse_path_and_ranges("file.txt:10-")
    assert path == "file.txt"
    assert ranges == [(10, None)]


def test_parse_path_and_ranges_with_single_line():
    """Test parsing a path with a single line range."""
    path, ranges = _parse_path_and_ranges("file.txt:10")
    assert path == "file.txt"
    assert ranges == [(10, 10)]


def test_parse_path_and_ranges_with_multiple_ranges():
    """Test parsing a path with multiple range specifiers."""
    path, ranges = _parse_path_and_ranges("file.txt:10-20,30-40")
    assert path == "file.txt"
    assert ranges == [(10, 20), (30, 40)]


def test_parse_path_and_ranges_with_mixed_ranges():
    """Test parsing a path with mixed range types."""
    path, ranges = _parse_path_and_ranges("file.txt:10-20,30,40-")
    assert path == "file.txt"
    assert ranges == [(10, 20), (30, 30), (40, None)]


def test_parse_path_and_ranges_with_invalid_range():
    """Test parsing a path with an invalid range."""
    with pytest.raises(ValueError):
        _parse_path_and_ranges("file.txt:abc")

    with pytest.raises(ValueError):
        _parse_path_and_ranges("file.txt:10-abc")

    with pytest.raises(ValueError):
        _parse_path_and_ranges("file.txt:abc-10")

    with pytest.raises(ValueError):
        _parse_path_and_ranges("file.txt:0-10")

    with pytest.raises(ValueError):
        _parse_path_and_ranges("file.txt:10-5")


def test_parse_path_and_ranges_with_spaces():
    """Test parsing a path with spaces in the range."""
    path, ranges = _parse_path_and_ranges("file.txt:10 - 20, 30 - 40")
    assert path == "file.txt"
    assert ranges == [(10, 20), (30, 40)]


def test_apply_ranges_entire_file():
    """Test applying ranges to get the entire file content."""
    lines = ["Line 1\n", "Line 2\n", "Line 3\n", "Line 4\n", "Line 5\n"]
    content = _apply_ranges(lines, [(1, None)])
    assert content == "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n"


def test_apply_ranges_single_range():
    """Test applying a single range."""
    lines = ["Line 1\n", "Line 2\n", "Line 3\n", "Line 4\n", "Line 5\n"]
    content = _apply_ranges(lines, [(2, 4)])
    assert content == "Line 2\nLine 3\n"


def test_apply_ranges_multiple_ranges():
    """Test applying multiple ranges."""
    lines = ["Line 1\n", "Line 2\n", "Line 3\n", "Line 4\n", "Line 5\n"]
    content = _apply_ranges(lines, [(1, 2), (4, 5)])
    assert content == "Line 1\nLine 4\n"


def test_apply_ranges_eof_range():
    """Test applying a range to the end of the file."""
    lines = ["Line 1\n", "Line 2\n", "Line 3\n", "Line 4\n", "Line 5\n"]
    content = _apply_ranges(lines, [(3, None)])
    assert content == "Line 3\nLine 4\nLine 5\n"


def test_apply_ranges_out_of_bounds():
    """Test applying ranges that are out of bounds."""
    lines = ["Line 1\n", "Line 2\n", "Line 3\n"]

    # Start beyond end of file
    content = _apply_ranges(lines, [(10, 20)])
    assert content == ""

    # End beyond end of file
    content = _apply_ranges(lines, [(2, 10)])
    assert content == "Line 2\nLine 3\n"


def test_resolve_files_basic():
    """Test resolving a list of files with no range specifiers."""
    file_paths = ["/path/to/file1.txt", "/path/to/file2.md"]
    result = resolve_files(file_paths)

    assert len(result) == 2

    assert result[0].filepath == "/path/to/file1.txt"
    assert result[0].ranges == [(1, None)]
    assert result[0].content == ""
    assert result[0].is_bundle is False

    assert result[1].filepath == "/path/to/file2.md"
    assert result[1].ranges == [(1, None)]
    assert result[1].content == ""
    assert result[1].is_bundle is False


def test_resolve_files_with_ranges():
    """Test resolving files with range specifiers."""
    file_paths = ["/path/to/file1.txt:10-20", "/path/to/file2.md:5"]
    result = resolve_files(file_paths)

    assert len(result) == 2

    assert result[0].filepath == "/path/to/file1.txt"
    assert result[0].ranges == [(10, 20)]
    assert result[0].content == ""
    assert result[0].is_bundle is False

    assert result[1].filepath == "/path/to/file2.md"
    assert result[1].ranges == [(5, 5)]
    assert result[1].content == ""
    assert result[1].is_bundle is False


def test_resolve_files_with_bundles():
    """Test resolving files with custom bundle extensions."""
    file_paths = ["/path/to/file1.ndoc", "/path/to/file2.bundle"]
    result = resolve_files(file_paths, bundle_extensions=[".ndoc", ".bundle"])

    assert len(result) == 2

    assert result[0].filepath == "/path/to/file1.ndoc"
    assert result[0].is_bundle is True

    assert result[1].filepath == "/path/to/file2.bundle"
    assert result[1].is_bundle is True


def test_gather_content(tmp_path):
    """Test gathering content from files."""
    # Create test files
    file1_path = tmp_path / "file1.txt"
    file1_content = "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n"
    file1_path.write_text(file1_content)

    file2_path = tmp_path / "file2.txt"
    file2_content = "File 2 Line 1\nFile 2 Line 2\nFile 2 Line 3\n"
    file2_path.write_text(file2_content)

    # Create FileContent objects
    file_contents = [
        FileContent(
            filepath=str(file1_path),
            ranges=[(2, 4)],  # Lines 2-3
            is_bundle=False,
        ),
        FileContent(
            filepath=str(file2_path),
            ranges=[(1, None)],  # Entire file
            is_bundle=True,
        ),
    ]

    # Gather content
    result = gather_content(file_contents)

    assert len(result) == 2

    # Check first file content (with range)
    assert result[0].filepath == str(file1_path)
    assert result[0].ranges == [(2, 4)]
    assert result[0].is_bundle is False
    assert result[0].content == "Line 2\nLine 3\n"

    # Check second file content (entire file)
    assert result[1].filepath == str(file2_path)
    assert result[1].ranges == [(1, None)]
    assert result[1].is_bundle is True
    assert result[1].content == file2_content


def test_gather_content_nonexistent_file():
    """Test gathering content from a nonexistent file."""
    file_contents = [
        FileContent(
            filepath="/path/to/nonexistent/file.txt",
            ranges=[(1, None)],
            is_bundle=False,
        )
    ]

    with pytest.raises(FileNotFoundError):
        gather_content(file_contents)


def test_gather_content_multiple_ranges(tmp_path):
    """Test gathering content with multiple ranges."""
    # Create test file
    file_path = tmp_path / "file.txt"
    file_content = "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n"
    file_path.write_text(file_content)

    # Create FileContent object with multiple ranges
    file_contents = [
        FileContent(
            filepath=str(file_path),
            ranges=[(1, 2), (4, 5)],  # Lines 1, 4
            is_bundle=False,
        )
    ]

    # Gather content
    result = gather_content(file_contents)

    assert len(result) == 1
    assert result[0].filepath == str(file_path)
    assert result[0].ranges == [(1, 2), (4, 5)]
    assert result[0].is_bundle is False
    assert result[0].content == "Line 1\nLine 4\n"


def test_gather_content_with_original_source(tmp_path):
    """Test gathering content with original_source set."""
    # Create test file
    file_path = tmp_path / "file.txt"
    file_content = "Line 1\nLine 2\nLine 3\n"
    file_path.write_text(file_content)

    # Create FileContent object with original_source
    file_contents = [
        FileContent(
            filepath=str(file_path),
            ranges=[(1, None)],
            is_bundle=False,
            original_source="/path/to/source.txt",
        )
    ]

    # Gather content
    result = gather_content(file_contents)

    assert len(result) == 1
    assert result[0].filepath == str(file_path)
    assert result[0].ranges == [(1, None)]
    assert result[0].is_bundle is False
    assert result[0].original_source == "/path/to/source.txt"
    assert result[0].content == file_content


def test_resolve_files_with_custom_bundle_extensions():
    """Test resolving files with custom bundle extensions."""
    file_paths = ["/path/to/file1.ndoc", "/path/to/file2.txt"]
    result = resolve_files(file_paths, bundle_extensions=[".ndoc", ".bundle"])

    assert len(result) == 2
    assert result[0].is_bundle is True
    assert result[1].is_bundle is False


def test_resolve_files_with_bundle_extensions():
    """Test resolving files with various bundle extensions."""
    file_paths = [
        "/path/to/file1.ndoc",  # Default bundle extension
        "/path/to/file2.bundle",  # Direct .bundle extension
        "/path/to/file3.bundle.txt",  # .bundle.* pattern
        "/path/to/file4.txt",  # Regular file
    ]

    # Use default bundle extensions
    result = resolve_files(file_paths)

    assert len(result) == 4
    assert result[0].is_bundle is True  # .ndoc file
    assert result[1].is_bundle is True  # .bundle file
    assert result[2].is_bundle is True  # .bundle.txt file
    assert result[3].is_bundle is False  # .txt file
