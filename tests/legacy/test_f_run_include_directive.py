"""Tests for legacy.run_include_directive function."""

from nanodoc.legacy import run_include_directive


def test_run_include_directive_basic():
    """Test the run_include_directive function with basic input."""
    content = "Test content"
    directive = "@include file.txt"
    result = run_include_directive(content, directive)
    # Legacy function should return content unchanged
    assert result == content


def test_run_include_directive_with_range():
    """Test run_include_directive with range specifier."""
    content = "Test content"
    directive = "@include file.txt:10-20"
    result = run_include_directive(content, directive)
    # Legacy function should return content unchanged
    assert result == content
