"""Tests for renderer.create_header function."""

from nanodoc.renderer import create_header


def test_create_header_basic():
    """Test creating a basic header."""
    header = create_header("Test Header")
    assert header == "Test Header"


def test_create_header_with_style():
    """Test creating a header with style (ignored in V2)."""
    header = create_header("Test Header", style="fancy")
    # Style should be ignored in V2
    assert header == "Test Header"


def test_create_header_empty():
    """Test creating an empty header."""
    header = create_header("")
    assert header == ""
