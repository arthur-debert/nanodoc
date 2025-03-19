"""Tests for formatter.load_theme function."""

from unittest.mock import mock_open, patch

from rich.theme import Theme

from nanodoc.formatter import load_theme


def test_load_theme():
    """Test loading a theme from a YAML file."""
    theme_data = """
    heading: "blue bold"
    heading.1: "bright_blue bold"
    code: "bright_green dim"
    """

    with (
        patch("pathlib.Path.exists") as mock_exists,
        patch("builtins.open", mock_open(read_data=theme_data)),
    ):
        mock_exists.return_value = True
        theme = load_theme("test-theme")

        # Check that the theme is a Rich Theme object
        assert isinstance(theme, Theme)

        # Check that styles were parsed correctly
        assert "heading" in theme.styles
        assert "heading.1" in theme.styles
        assert "code" in theme.styles


def test_load_theme_with_fallback():
    """Test loading a theme with fallback to default."""
    theme_data = """
    heading: "blue bold"
    """

    with (
        patch("pathlib.Path.exists") as mock_exists,
        patch("builtins.open", mock_open(read_data=theme_data)),
    ):
        # First file doesn't exist, second does
        mock_exists.side_effect = [False, True]
        theme = load_theme("nonexistent-theme")

        # Check that the theme is a Rich Theme object
        assert isinstance(theme, Theme)

        # Check that styles were parsed correctly
        assert "heading" in theme.styles


def test_load_theme_with_error():
    """Test loading a theme with an error."""
    with (
        patch("pathlib.Path.exists") as mock_exists,
        patch("builtins.open", mock_open(read_data="invalid: yaml: :")),
    ):
        mock_exists.return_value = True
        # This should not raise an exception but return a minimal theme
        theme = load_theme("invalid-theme")

        # Verify minimal default theme is returned
        assert isinstance(theme, Theme)
        assert "heading" in theme.styles
        assert "error" in theme.styles
