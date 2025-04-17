"""Tests for legacy.run_inline_directive function."""

from nanodoc.legacy import run_inline_directive


def test_run_inline_directive_basic():
    """Test the run_inline_directive function with basic input."""
    content = "Test content"
    directive = "@inline file.txt"
    result = run_inline_directive(content, directive)
    # Legacy function should return content unchanged
    assert result == content


def test_run_inline_directive_with_range():
    """Test run_inline_directive with range specifier."""
    content = "Test content"
    directive = "@inline file.txt:10-20"
    result = run_inline_directive(content, directive)
    # Legacy function should return content unchanged
    assert result == content
