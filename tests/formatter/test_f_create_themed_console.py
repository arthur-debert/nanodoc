"""Tests for formatter.create_themed_console function."""

from unittest.mock import patch

from rich.theme import Theme

from nanodoc.formatter import create_themed_console


def test_create_themed_console():
    """Test creating a themed console."""
    with patch("nanodoc.formatter.load_theme") as mock_load_theme:
        mock_load_theme.return_value = Theme({"heading": ""})
        console = create_themed_console("test-theme")
        mock_load_theme.assert_called_once_with("test-theme")
        assert console is not None


def test_create_themed_console_with_default():
    """Test creating a themed console with default theme."""
    with patch("nanodoc.formatter.load_theme") as mock_load_theme:
        mock_load_theme.return_value = Theme({})
        console = create_themed_console()
        mock_load_theme.assert_called_once_with("classic")
        assert console is not None
