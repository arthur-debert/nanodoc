"""Tests for file title formatting in the renderer."""

from unittest.mock import patch

from nanodoc.renderer import render_document
from nanodoc.structures import Document, FileContent


def _extract_file_title(rendered_content: str, expected_title_part: str) -> str:
    """Helper to extract the file title from rendered content.

    This extracts the entire title line containing the expected part.
    """
    for line in rendered_content.split("\n"):
        if expected_title_part in line:
            return line.strip()
    return ""


def test_simple_word():
    """Test formatting a filename that's a single word with no extension."""
    with patch("os.getcwd") as mock_getcwd, patch("os.path.abspath") as mock_abspath:
        mock_getcwd.return_value = "/base"
        # Make abspath return the input to simplify testing
        mock_abspath.side_effect = lambda x: x

        file1 = FileContent(
            filepath="/base/word",
            ranges=[(1, None)],
            content="Some content",
            is_bundle=False,
        )
        document = Document(content_items=[file1])

        result = render_document(document)

        # Should be title-cased with relative path
        title = _extract_file_title(result, "word")
        assert title == "Word (word)"


def test_word_with_extension():
    """Test formatting a filename with a simple word and extension."""
    with patch("os.getcwd") as mock_getcwd, patch("os.path.abspath") as mock_abspath:
        mock_getcwd.return_value = "/base"
        # Make abspath return the input to simplify testing
        mock_abspath.side_effect = lambda x: x

        file1 = FileContent(
            filepath="/base/word.txt",
            ranges=[(1, None)],
            content="Some content",
            is_bundle=False,
        )
        document = Document(content_items=[file1])

        result = render_document(document)

        # Should be title-cased with extension in parentheses
        title = _extract_file_title(result, "word.txt")
        assert title == "Word (word.txt)"


def test_camel_case():
    """Test formatting a camel-case filename."""
    with patch("os.getcwd") as mock_getcwd, patch("os.path.abspath") as mock_abspath:
        mock_getcwd.return_value = "/base"
        # Make abspath return the input to simplify testing
        mock_abspath.side_effect = lambda x: x

        file1 = FileContent(
            filepath="/base/wordNice.txt",
            ranges=[(1, None)],
            content="Some content",
            is_bundle=False,
        )
        document = Document(content_items=[file1])

        result = render_document(document)

        # Should break camel case into separate words
        title = _extract_file_title(result, "wordNice.txt")
        assert title == "Word Nice (wordNice.txt)"


def test_dash_separated():
    """Test formatting a dash-separated filename."""
    with patch("os.getcwd") as mock_getcwd, patch("os.path.abspath") as mock_abspath:
        mock_getcwd.return_value = "/base"
        # Make abspath return the input to simplify testing
        mock_abspath.side_effect = lambda x: x

        file1 = FileContent(
            filepath="/base/word-nice.txt",
            ranges=[(1, None)],
            content="Some content",
            is_bundle=False,
        )
        document = Document(content_items=[file1])

        result = render_document(document)

        # Should replace dashes with spaces and title-case
        title = _extract_file_title(result, "word-nice.txt")
        assert title == "Word Nice (word-nice.txt)"


def test_underscore_separated():
    """Test formatting an underscore-separated filename."""
    with patch("os.getcwd") as mock_getcwd, patch("os.path.abspath") as mock_abspath:
        mock_getcwd.return_value = "/base"
        # Make abspath return the input to simplify testing
        mock_abspath.side_effect = lambda x: x

        file1 = FileContent(
            filepath="/base/word_nice.txt",
            ranges=[(1, None)],
            content="Some content",
            is_bundle=False,
        )
        document = Document(content_items=[file1])

        result = render_document(document)

        # Should replace underscores with spaces and title-case
        title = _extract_file_title(result, "word_nice.txt")
        assert title == "Word Nice (word_nice.txt)"


def test_space_separated():
    """Test formatting a space-separated filename."""
    with patch("os.getcwd") as mock_getcwd, patch("os.path.abspath") as mock_abspath:
        mock_getcwd.return_value = "/base"
        # Make abspath return the input to simplify testing
        mock_abspath.side_effect = lambda x: x

        file1 = FileContent(
            filepath="/base/word nice.txt",
            ranges=[(1, None)],
            content="Some content",
            is_bundle=False,
        )
        document = Document(content_items=[file1])

        result = render_document(document)

        # Should maintain spaces and title-case
        title = _extract_file_title(result, "word nice.txt")
        assert title == "Word Nice (word nice.txt)"


def test_multiple_word_separators():
    """Test formatting a filename with multiple types of word separators."""
    with patch("os.getcwd") as mock_getcwd, patch("os.path.abspath") as mock_abspath:
        mock_getcwd.return_value = "/base"
        # Make abspath return the input to simplify testing
        mock_abspath.side_effect = lambda x: x

        file1 = FileContent(
            filepath="/base/word-nice_very good.txt",
            ranges=[(1, None)],
            content="Some content",
            is_bundle=False,
        )
        document = Document(content_items=[file1])

        result = render_document(document)

        # Should replace all separators with spaces
        title = _extract_file_title(result, "word-nice_very good.txt")
        assert title == "Word Nice Very Good (word-nice_very good.txt)"
