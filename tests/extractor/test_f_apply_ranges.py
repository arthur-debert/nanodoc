"""Tests for extractor.apply_ranges function."""

import pytest

from nanodoc.extractor import apply_ranges


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
