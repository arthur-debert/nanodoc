"""Tests for relative path handling in file headers."""

from unittest.mock import patch

from nanodoc.renderer import get_relative_path, render_document
from nanodoc.structures import Document, FileContent


def test_get_relative_path_same_dir():
    """Test getting relative path when file is in the same directory."""
    with patch("os.getcwd") as mock_getcwd:
        mock_getcwd.return_value = "/home/user/projects"

        # File in the same directory
        filepath = "/home/user/projects/file.txt"
        rel_path = get_relative_path(filepath)
        assert rel_path == "file.txt"


def test_get_relative_path_subdir():
    """Test getting relative path when file is in a subdirectory."""
    with patch("os.getcwd") as mock_getcwd:
        mock_getcwd.return_value = "/home/user/projects"

        # File in a subdirectory
        filepath = "/home/user/projects/proj-a/file.txt"
        rel_path = get_relative_path(filepath)
        assert rel_path == "proj-a/file.txt"


def test_get_relative_path_parent_dir():
    """Test getting relative path when file is in a parent directory."""
    with patch("os.getcwd") as mock_getcwd:
        mock_getcwd.return_value = "/home/user/projects/subdir"

        # File in a parent directory
        filepath = "/home/user/projects/file.txt"
        rel_path = get_relative_path(filepath)
        # For files outside cwd, we should return absolute path
        assert rel_path == "/home/user/projects/file.txt"


def test_get_relative_path_different_branch():
    """Test getting relative path when file is in a completely different path."""
    with patch("os.getcwd") as mock_getcwd:
        mock_getcwd.return_value = "/home/user/projects"

        # File in an unrelated directory
        filepath = "/var/log/file.txt"
        rel_path = get_relative_path(filepath)
        # For files outside cwd, we should return absolute path
        assert rel_path == "/var/log/file.txt"


def test_relative_path_in_header():
    """Test that file headers use relative paths."""
    with patch("os.getcwd") as mock_getcwd, patch("os.path.abspath") as mock_abspath:
        mock_getcwd.return_value = "/home/user/projects"

        # Make abspath return the input to simplify testing
        mock_abspath.side_effect = lambda x: x

        # Create test files that would be in different directories
        file1 = FileContent(
            filepath="/home/user/projects/proj-a/README.txt",
            ranges=[(1, None)],
            content="Content of README in proj-a",
            is_bundle=False,
        )
        file2 = FileContent(
            filepath="/home/user/projects/proj-b/README.txt",
            ranges=[(1, None)],
            content="Content of README in proj-b",
            is_bundle=False,
        )
        document = Document(content_items=[file1, file2])

        # Render document
        result = render_document(document)

        # Check that headers include relative paths
        assert "Readme (proj-a/README.txt)" in result
        assert "Readme (proj-b/README.txt)" in result


def test_relative_path_same_name_files():
    """Test that file headers properly distinguish between files with same name."""
    with patch("os.getcwd") as mock_getcwd, patch("os.path.abspath") as mock_abspath:
        mock_getcwd.return_value = "/home/user/projects"

        # Make abspath return the input to simplify testing
        mock_abspath.side_effect = lambda x: x

        # Create test files that have the same name but in different directories
        file1 = FileContent(
            filepath="/home/user/projects/dir1/test.txt",
            ranges=[(1, None)],
            content="Content of test.txt in dir1",
            is_bundle=False,
        )
        file2 = FileContent(
            filepath="/home/user/projects/dir2/test.txt",
            ranges=[(1, None)],
            content="Content of test.txt in dir2",
            is_bundle=False,
        )
        file3 = FileContent(
            filepath="/home/user/projects/test.txt",
            ranges=[(1, None)],
            content="Content of test.txt in root",
            is_bundle=False,
        )
        document = Document(content_items=[file1, file2, file3])

        # Render document
        result = render_document(document)

        # Check that headers include proper relative paths
        assert "Test (dir1/test.txt)" in result
        assert "Test (dir2/test.txt)" in result
        assert "Test (test.txt)" in result


def _extract_header(rendered_content: str, path_part: str) -> str:
    """Helper to extract the file header from rendered content."""
    for line in rendered_content.split("\n"):
        if path_part in line:
            return line.strip()
    return ""
