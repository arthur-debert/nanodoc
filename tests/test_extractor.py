"""Tests for file extraction in Nanodoc v2."""

import pytest

from nanodoc.extractor import (
    apply_ranges,
    gather_content,
    parse_path_and_ranges,
    resolve_files,
)
from nanodoc.structures import FileContent


def test_parse_path_and_ranges_no_ranges():
    """Test parsing a path with no range specifiers."""
    path, ranges = parse_path_and_ranges("file.txt")
    assert path == "file.txt"
    assert ranges == [(1, None)]


def test_parse_path_and_ranges_with_single_range():
    """Test parsing a path with a single range specifier."""
    path, ranges = parse_path_and_ranges("file.txt:10-20")
    assert path == "file.txt"
    assert ranges == [(10, 20)]


def test_parse_path_and_ranges_with_eof_range():
    """Test parsing a path with a range to end of file."""
    path, ranges = parse_path_and_ranges("file.txt:10-")
    assert path == "file.txt"
    assert ranges == [(10, None)]


def test_parse_path_and_ranges_with_single_line():
    """Test parsing a path with a single line range."""
    path, ranges = parse_path_and_ranges("file.txt:10")
    assert path == "file.txt"
    assert ranges == [(10, 10)]


def test_parse_path_and_ranges_with_multiple_ranges():
    """Test parsing a path with multiple range specifiers."""
    path, ranges = parse_path_and_ranges("file.txt:10-20,30-40")
    assert path == "file.txt"
    assert ranges == [(10, 20), (30, 40)]


def test_parse_path_and_ranges_with_invalid_range():
    """Test parsing a path with an invalid range."""
    with pytest.raises(ValueError):
        parse_path_and_ranges("file.txt:invalid")


def test_parse_path_and_ranges_with_negative_line():
    """Test parsing a path with a negative line number."""
    with pytest.raises(ValueError):
        parse_path_and_ranges("file.txt:-10")


def test_parse_path_and_ranges_with_invalid_range_order():
    """Test parsing a path with an invalid range order."""
    with pytest.raises(ValueError):
        parse_path_and_ranges("file.txt:20-10")


def test_parse_path_and_ranges_with_spaces():
    """Test parsing a path with spaces in the range."""
    path, ranges = parse_path_and_ranges("file.txt:10 - 20, 30 - 40")
    assert path == "file.txt"
    assert ranges == [(10, 20), (30, 40)]


def test_parse_path_and_ranges_with_empty_range():
    """Test parsing a path with an empty range."""
    path, ranges = parse_path_and_ranges("file.txt:")
    assert path == "file.txt"
    assert ranges == [(1, None)]


def test_apply_ranges_single_range():
    """Test applying a single range."""
    lines = ["Line 1\n", "Line 2\n", "Line 3\n", "Line 4\n", "Line 5\n"]
    content = apply_ranges(lines, [(2, 4)])
    assert content == "Line 2\nLine 3\n"


def test_apply_ranges_multiple_ranges():
    """Test applying multiple ranges."""
    lines = ["Line 1\n", "Line 2\n", "Line 3\n", "Line 4\n", "Line 5\n"]
    content = apply_ranges(lines, [(1, 2), (4, 5)])
    assert content == "Line 1\nLine 4\n"


def test_apply_ranges_with_eof():
    """Test applying a range to the end of file."""
    lines = ["Line 1\n", "Line 2\n", "Line 3\n", "Line 4\n", "Line 5\n"]
    content = apply_ranges(lines, [(3, None)])
    assert content == "Line 3\nLine 4\nLine 5\n"


def test_apply_ranges_with_single_line():
    """Test applying a single line range."""
    lines = ["Line 1\n", "Line 2\n", "Line 3\n", "Line 4\n", "Line 5\n"]
    content = apply_ranges(lines, [(3, 3)])
    # Per implementation, single line ranges include just that line
    assert content == "Line 3\n"


def test_apply_ranges_with_invalid_range():
    """Test applying an invalid range."""
    lines = ["Line 1\n", "Line 2\n", "Line 3\n"]
    with pytest.raises(ValueError):
        apply_ranges(lines, [(-1, 2)])


def test_apply_ranges_with_empty_lines():
    """Test applying ranges to empty lines."""
    content = apply_ranges([], [(1, 3)])
    assert content == ""


def test_apply_ranges_with_out_of_bounds():
    """Test applying ranges that are out of bounds."""
    lines = ["Line 1\n", "Line 2\n", "Line 3\n"]
    # This should not error but should just return nothing
    content = apply_ranges(lines, [(10, 20)])
    assert content == ""


def test_gather_content_basic(tmp_path):
    """Test gathering content from files."""
    # Create temporary files for testing
    file1 = tmp_path / "file1.txt"
    file1.write_text("Line 1\nLine 2\nLine 3\n")

    file2 = tmp_path / "file2.txt"
    file2.write_text("Line A\nLine B\nLine C\n")

    # Create FileContent objects
    file_contents = [
        FileContent(filepath=str(file1), ranges=[(1, None)], is_bundle=False),
        FileContent(filepath=str(file2), ranges=[(1, None)], is_bundle=False),
    ]

    # Gather content
    result = gather_content(file_contents)

    # Check result
    assert len(result) == 2
    assert result[0].content == "Line 1\nLine 2\nLine 3\n"
    assert result[1].content == "Line A\nLine B\nLine C\n"


def test_gather_content_with_ranges(tmp_path):
    """Test gathering content from files with ranges."""
    # Create temporary file for testing
    file1 = tmp_path / "file1.txt"
    file1.write_text("Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n")

    # Create FileContent objects with ranges
    file_contents = [
        FileContent(filepath=str(file1), ranges=[(2, 4)], is_bundle=False),
    ]

    # Gather content
    result = gather_content(file_contents)

    # Check result
    assert len(result) == 1
    assert result[0].content == "Line 2\nLine 3\n"


def test_gather_content_file_not_found():
    """Test gathering content from non-existent files."""
    # Create FileContent object for non-existent file
    file_contents = [
        FileContent(
            filepath="/path/to/nonexistent/file.txt",
            ranges=[(1, None)],
            is_bundle=False,
        ),
    ]

    # Should raise FileNotFoundError
    with pytest.raises(FileNotFoundError):
        gather_content(file_contents)


def test_resolve_files_basic():
    """Test resolving files with default settings."""
    file_paths = ["/path/to/file1.txt", "/path/to/file2.txt"]

    # Resolve files
    result = resolve_files(file_paths)

    # Check result
    assert len(result) == 2
    assert result[0].filepath == "/path/to/file1.txt"
    assert result[0].ranges == [(1, None)]
    assert result[0].is_bundle is False
    assert result[1].filepath == "/path/to/file2.txt"
    assert result[1].ranges == [(1, None)]
    assert result[1].is_bundle is False


def test_resolve_files_with_ranges():
    """Test resolving files with range specifiers."""
    file_paths = ["/path/to/file1.txt:10-20", "/path/to/file2.txt:5"]

    # Resolve files
    result = resolve_files(file_paths)

    # Check result
    assert len(result) == 2
    assert result[0].filepath == "/path/to/file1.txt"
    assert result[0].ranges == [(10, 20)]
    assert result[0].is_bundle is False
    assert result[1].filepath == "/path/to/file2.txt"
    assert result[1].ranges == [(5, 5)]
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


def test_resolve_files_with_custom_bundle_extensions():
    """Test resolving files with custom bundle extensions."""
    file_paths = [
        "/path/to/file1.ndoc",  # Standard extension
        "/path/to/file2.bundle",  # Standard extension
        "/path/to/file3.custom",  # Custom extension
        "/path/to/file4.txt",  # Regular file
    ]

    # Use custom bundle extensions
    result = resolve_files(file_paths, bundle_extensions=[".custom"])

    assert len(result) == 4
    assert result[0].is_bundle is False  # .ndoc file (not in custom list)
    assert result[1].is_bundle is False  # .bundle file (not in custom list)
    assert result[2].is_bundle is True  # .custom file (in custom list)
    assert result[3].is_bundle is False  # .txt file
