"""Tests for legacy.run_content function."""

from nanodoc.legacy import run_content


def test_run_content_basic():
    """Test the run_content function with basic input."""
    content = "Test content"
    result = run_content(content)
    # Legacy function should return content unchanged
    assert result == content


def test_run_content_with_options():
    """Test the run_content function with various options."""
    content = "Test content"
    result = run_content(
        content,
        line_number_mode="file",
        generate_toc=True,
        theme="test-theme",
        show_header=False,
    )
    # Legacy function should return content unchanged regardless of options
    assert result == content
