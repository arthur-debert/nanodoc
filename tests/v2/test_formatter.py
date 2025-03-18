"""Tests for the formatting stage of Nanodoc v2."""

from unittest.mock import mock_open, patch

from rich.theme import Theme

from nanodoc.v2.formatter import (
    apply_theme_to_document,
    create_themed_console,
    enhance_rendering,
    format_with_line_numbers,
    get_available_themes,
    load_theme,
)
from nanodoc.v2.structures import Document, FileContent


def test_get_available_themes():
    """Test getting available themes."""
    with patch("os.listdir") as mock_listdir:
        mock_listdir.return_value = ["neutral.yaml", "classic-dark.yaml"]
        themes = get_available_themes()
        assert "neutral" in themes
        assert "classic-dark" in themes
        assert len(themes) == 2


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


def test_create_themed_console():
    """Test creating a themed console."""
    with patch("nanodoc.v2.formatter.load_theme") as mock_load_theme:
        mock_load_theme.return_value = Theme({"heading": ""})
        console = create_themed_console("test-theme")
        mock_load_theme.assert_called_once_with("test-theme")
        assert console is not None


def test_create_themed_console_with_default():
    """Test creating a themed console with default theme."""
    with patch("nanodoc.v2.formatter.load_theme") as mock_load_theme:
        mock_load_theme.return_value = Theme({"heading": ""})
        console = create_themed_console()
        mock_load_theme.assert_called_once_with("neutral")
        assert console is not None


def test_apply_theme_to_document():
    """Test applying theme to a document."""
    # Create a document
    file1 = FileContent(
        filepath="/path/to/file1.txt",
        ranges=[(1, None)],
        content="Content of file 1",
        is_bundle=False,
    )
    document = Document(content_items=[file1])

    # Apply theme
    themed_document = apply_theme_to_document(
        document, theme_name="test-theme", use_rich_formatting=True
    )

    # Check that theme info was added
    assert themed_document.theme_name == "test-theme"
    assert themed_document.use_rich_formatting is True


def test_apply_theme_to_document_without_formatting():
    """Test applying theme to a document with formatting disabled."""
    # Create a document
    file1 = FileContent(
        filepath="/path/to/file1.txt",
        ranges=[(1, None)],
        content="Content of file 1",
        is_bundle=False,
    )
    document = Document(content_items=[file1])

    # Apply theme with formatting disabled
    themed_document = apply_theme_to_document(
        document, theme_name="test-theme", use_rich_formatting=False
    )

    # Check that theme info was not added
    assert themed_document.theme_name is None
    assert themed_document.use_rich_formatting is False


def test_format_with_line_numbers():
    """Test formatting content with line numbers."""
    content = "Line 1\nLine 2\nLine 3"
    result = format_with_line_numbers(content)

    # Check the result
    assert "   1 | Line 1" in result
    assert "   2 | Line 2" in result
    assert "   3 | Line 3" in result


def test_format_with_line_numbers_custom_start():
    """Test formatting content with line numbers starting from a custom number."""
    content = "Line 1\nLine 2\nLine 3"
    result = format_with_line_numbers(content, start_number=10)

    # Check the result
    assert "  10 | Line 1" in result
    assert "  11 | Line 2" in result
    assert "  12 | Line 3" in result


def test_format_with_line_numbers_custom_format():
    """Test formatting content with line numbers using a custom format."""
    content = "Line 1\nLine 2\nLine 3"
    result = format_with_line_numbers(content, number_format="Line {}: ")

    # Check the result
    assert "Line 1: Line 1" in result
    assert "Line 2: Line 2" in result
    assert "Line 3: Line 3" in result


def test_enhance_rendering():
    """Test enhancing rendered content with Rich formatting."""
    with (
        patch("nanodoc.v2.formatter.create_themed_console") as mock_console,
        patch("rich.console.Console.print") as mock_print,
    ):
        content = "# Heading 1\nContent\n## Heading 2\nMore content"
        # We call the function but don't need to check the actual result
        # since we're mocking the console output
        enhance_rendering(content, theme_name="test-theme")

        # Check that the console was created
        mock_console.assert_called_once()

        # Check that the content was printed
        assert mock_print.call_count == 4


def test_enhance_rendering_without_formatting():
    """Test enhancing rendered content with formatting disabled."""
    content = "# Heading 1\nContent\n## Heading 2\nMore content"
    result = enhance_rendering(content, use_rich_formatting=False)

    # Check that the content was returned unchanged
    assert result == content
