"""Tests for formatter.get_available_themes function."""

from unittest.mock import patch

from nanodoc.formatter import get_available_themes


def test_get_available_themes():
    """Test getting available themes."""
    with patch("os.listdir") as mock_listdir:
        mock_listdir.return_value = ["classic.yaml", "classic-dark.yaml"]
        themes = get_available_themes()
        assert "classic" in themes
        assert "classic-dark" in themes
        assert len(themes) == 2


def test_get_available_themes_with_non_yaml():
    """Test getting available themes with non-yaml files."""
    with patch("os.listdir") as mock_listdir:
        mock_listdir.return_value = ["classic.yaml", "notes.txt", "dark.yaml"]
        themes = get_available_themes()
        assert "classic" in themes
        assert "dark" in themes
        assert "notes" not in themes
        assert len(themes) == 2


def test_get_available_themes_empty_dir():
    """Test getting available themes from an empty directory."""
    with patch("os.listdir") as mock_listdir:
        mock_listdir.return_value = []
        themes = get_available_themes()
        assert len(themes) == 0


def test_get_available_themes_dir_not_exists():
    """Test getting available themes when directory doesn't exist."""
    with (
        patch("os.listdir") as mock_listdir,
        patch("pathlib.Path.exists") as mock_exists,
    ):
        mock_exists.return_value = False
        mock_listdir.side_effect = FileNotFoundError
        themes = get_available_themes()
        assert len(themes) == 0
