"""Tests for formatter.enhance_rendering function."""

from unittest.mock import patch

from nanodoc.formatter import enhance_rendering


def test_enhance_rendering():
    """Test enhancing rendered content with Rich formatting."""
    with (
        patch("nanodoc.formatter.load_theme") as mock_load_theme,
        patch("rich.console.Console.print") as mock_print,
    ):
        content = "# Heading 1\nContent\n## Heading 2\nMore content"
        # We call the function but don't need to check the actual result
        # since we're mocking the console output
        enhance_rendering(content, theme_name="test-theme")

        # Check that the theme was loaded
        mock_load_theme.assert_called_once()

        # Check that the content was printed
        assert mock_print.call_count == 4


def test_enhance_rendering_without_formatting():
    """Test enhancing rendered content with formatting disabled."""
    content = "# Heading 1\nContent\n## Heading 2\nMore content"
    result = enhance_rendering(content, use_rich_formatting=False)

    # Check that the content was returned unchanged
    assert result == content
