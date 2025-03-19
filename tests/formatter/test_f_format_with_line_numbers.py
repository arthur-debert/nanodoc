"""Tests for formatter.format_with_line_numbers function."""

from nanodoc.formatter import format_with_line_numbers


def test_format_with_line_numbers():
    """Test formatting content with line numbers."""
    content = "Line 1\nLine 2\nLine 3"
    result = format_with_line_numbers(content)

    # Check the result
    assert "   1: Line 1" in result
    assert "   2: Line 2" in result
    assert "   3: Line 3" in result


def test_format_with_line_numbers_custom_start():
    """Test formatting content with custom start number."""
    content = "Line 1\nLine 2\nLine 3"
    result = format_with_line_numbers(content, start_number=10)

    # Check the result
    assert "  10: Line 1" in result
    assert "  11: Line 2" in result
    assert "  12: Line 3" in result


def test_format_with_line_numbers_custom_format():
    """Test formatting content with line numbers using a custom format."""
    content = "Line 1\nLine 2\nLine 3"
    result = format_with_line_numbers(content, number_format="Line {}: ")

    # Check the result
    assert "Line 1: Line 1" in result
    assert "Line 2: Line 2" in result
    assert "Line 3: Line 3" in result


def test_format_with_line_numbers_empty_content():
    """Test formatting empty content with line numbers."""
    content = ""
    result = format_with_line_numbers(content)

    # Empty content should return empty string
    assert result == ""


def test_format_with_line_numbers_trailing_newline():
    """Test formatting content with trailing newline."""
    content = "Line 1\nLine 2\n"
    result = format_with_line_numbers(content)

    # Check that all lines have numbers
    assert "   1: Line 1" in result
    assert "   2: Line 2" in result
    # The function treats the trailing newline as an empty line,
    # so there should be 3 lines total (line 3 is empty)
    assert len(result.strip().split("\n")) == 3
    assert result.strip().split("\n")[2] == "   3:"
