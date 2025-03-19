"""Tests for resolver.get_files_from_directory function."""

from unittest.mock import patch

from nanodoc.resolver import get_files_from_directory


def test_get_files_from_directory_non_recursive():
    """Test getting files from a directory (non-recursive)."""
    with (
        patch("os.walk") as mock_walk,
        patch("os.listdir") as mock_listdir,
        patch("os.path.isfile") as mock_isfile,
    ):
        # Mock os.listdir to return a list of files and directories
        mock_listdir.return_value = ["file1.txt", "file2.py", ".hidden", "subdir"]

        # Mock os.path.isfile to return True for files
        mock_isfile.side_effect = lambda path: not path.endswith("subdir")

        # Call the function
        result = get_files_from_directory("/path/to/dir", recursive=False)

        # Check that os.walk was not called
        mock_walk.assert_not_called()

        # Check that os.listdir was called with the correct path
        mock_listdir.assert_called_once_with("/path/to/dir")

        # Check the result
        # Should only include files (not directories) and exclude hidden files
        assert len(result) == 2
        assert any(path.endswith("file1.txt") for path in result)
        assert any(path.endswith("file2.py") for path in result)
        assert not any(path.endswith(".hidden") for path in result)
        assert not any(path.endswith("subdir") for path in result)


def test_get_files_from_directory_recursive():
    """Test getting files from a directory recursively."""
    with patch("os.walk") as mock_walk:
        # Mock os.walk to return a list of directories and files
        mock_walk.return_value = [
            (
                "/path/to/dir",
                ["subdir1", ".hidden_dir"],
                ["file1.txt", ".hidden"],
            ),
            ("/path/to/dir/subdir1", [], ["file2.py"]),
        ]

        # Call the function
        result = get_files_from_directory("/path/to/dir", recursive=True)

        # Check the result
        assert len(result) == 2
        assert any(path.endswith("file1.txt") for path in result)
        assert any(path.endswith("file2.py") for path in result)
        assert not any(path.endswith(".hidden") for path in result)


def test_get_files_from_directory_include_hidden():
    """Test getting files from a directory including hidden files."""
    with (
        patch("os.walk") as mock_walk,
        patch("os.listdir") as mock_listdir,
        patch("os.path.isfile") as mock_isfile,
    ):
        # Mock for non-recursive case
        mock_listdir.return_value = ["file1.txt", ".hidden"]
        mock_isfile.return_value = True

        # Mock for recursive case
        mock_walk.return_value = [
            ("/path/to/dir", [".hidden_dir"], ["file1.txt", ".hidden"]),
            ("/path/to/dir/.hidden_dir", [], [".hidden_file"]),
        ]

        # Test non-recursive with include_hidden=True
        result1 = get_files_from_directory(
            "/path/to/dir", recursive=False, include_hidden=True
        )
        assert len(result1) == 2
        assert any(path.endswith("file1.txt") for path in result1)
        assert any(path.endswith(".hidden") for path in result1)

        # Test recursive with include_hidden=True
        result2 = get_files_from_directory(
            "/path/to/dir", recursive=True, include_hidden=True
        )
        assert len(result2) == 3
        assert any(path.endswith("file1.txt") for path in result2)
        assert any(path.endswith(".hidden") for path in result2)
        assert any(path.endswith(".hidden_file") for path in result2)
