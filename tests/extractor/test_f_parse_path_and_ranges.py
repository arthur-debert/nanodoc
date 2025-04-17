"""Tests for extractor.parse_path_and_ranges function."""

import pytest

from nanodoc.extractor import parse_path_and_ranges


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
