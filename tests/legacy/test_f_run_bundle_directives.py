"""Tests for legacy.run_bundle_directives function."""

from nanodoc.legacy import run_bundle_directives


def test_run_bundle_directives_basic():
    """Test the run_bundle_directives function with basic input."""
    content = "Line 1\nLine 2\nLine 3"
    result = run_bundle_directives(content)
    # Without directives, should return content unchanged
    assert result == content


def test_run_bundle_directives_with_directives():
    """Test run_bundle_directives with directives."""
    content = "Line 1\n@inline file.txt\nLine 3\n@include other.txt"
    result = run_bundle_directives(content)
    # Should skip the directive lines
    expected = "Line 1\nLine 3"
    assert result == expected
