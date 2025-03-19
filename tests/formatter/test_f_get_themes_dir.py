"""Tests for formatter._get_themes_dir function."""

import os
import pathlib

from nanodoc.formatter import _get_themes_dir


def test_get_themes_dir():
    """Test getting themes directory."""
    # Call the function directly without mocking
    result = _get_themes_dir()

    # Verify it's a Path object
    assert isinstance(result, pathlib.Path)

    # Verify it's named 'themes'
    assert result.name == "themes"

    # Verify it's a subdirectory of the nanodoc module
    assert result.parent.name == "nanodoc"

    # Verify the path actually exists
    assert os.path.exists(result)


def test_get_themes_dir_returns_subdirectory():
    """Test that _get_themes_dir returns a subdirectory of the module dir."""
    # No mocks here - test the actual behavior
    result = _get_themes_dir()

    # Verify it's a path object
    assert isinstance(result, pathlib.Path)

    # Verify it's a themes directory
    assert result.name == "themes"

    # Verify the parent is the module directory
    parent = result.parent
    assert parent.name == "nanodoc"
